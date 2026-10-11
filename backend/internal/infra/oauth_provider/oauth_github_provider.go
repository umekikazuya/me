package oauthprovider

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/umekikazuya/me/internal/app/port"
	"github.com/umekikazuya/me/pkg/errs"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var conf = &oauth2.Config{
	ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
	ClientSecret: os.Getenv("GITHUB_CLIENT_SECRETS"),
	Endpoint:     github.Endpoint,
	RedirectURL:  os.Getenv("GITHUB_REDIRECT_URI"),
	Scopes:       []string{"read:user"},
}

type OauthGithubProvider struct{}

// GetClient implements [port.OauthProvider].
func (o *OauthGithubProvider) GetClient(ctx context.Context, code string) (*http.Client, error) {
	t, err := conf.Exchange(ctx, code)
	if err != nil {
		return nil, errs.New(errs.ErrBadRequest, err.Error())
	}
	return conf.Client(ctx, t), nil
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

func New() port.OauthProvider {
	return &OauthGithubProvider{}
}

var _ port.OauthProvider = (*OauthGithubProvider)(nil)
