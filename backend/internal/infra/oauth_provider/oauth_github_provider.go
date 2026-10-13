package oauthprovider

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/umekikazuya/me/internal/app/port"
	"github.com/umekikazuya/me/pkg/errs"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type OauthGithubProvider struct {
	conf       *oauth2.Config
	httpClient *http.Client
}

// AuthCodeURL implements [port.OauthProvider].
func (o *OauthGithubProvider) AuthCodeURL(state string) string {
	return o.conf.AuthCodeURL(state)
}

// GetClient implements [port.OauthProvider].
func (o *OauthGithubProvider) GetClient(ctx context.Context, code string) (*http.Client, error) {
	ctx = context.WithValue(
		ctx,
		oauth2.HTTPClient,
		o.httpClient,
	)
	t, err := o.conf.Exchange(ctx, code)
	if err != nil {
		return nil, errs.New(errs.ErrBadRequest, err.Error())
	}
	return o.conf.Client(ctx, t), nil
}

// GetResource implements [port.OauthProvider].
func (o *OauthGithubProvider) GetResource(
	ctx context.Context,
	client *http.Client,
) ([]byte, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://api.github.com/user",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("データ取得エラー: %w", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("github user 情報取得エラー: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK {
		return nil, fmt.Errorf("不正なリクエスト: code=%d", resp.StatusCode)
	}
	if http.StatusMultipleChoices <= resp.StatusCode {
		return nil, fmt.Errorf("不正なリクエスト: code=%d", resp.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, errs.WrapInternal("システムエラー", err)
	}
	return body, nil
}

func New(
	clientID, clientSecret, RedirectURI string,
) port.OauthProvider {
	return &OauthGithubProvider{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		conf: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Endpoint:     github.Endpoint,
			RedirectURL:  RedirectURI,
			Scopes:       []string{"read:user"},
		},
	}
}

var _ port.OauthProvider = (*OauthGithubProvider)(nil)
