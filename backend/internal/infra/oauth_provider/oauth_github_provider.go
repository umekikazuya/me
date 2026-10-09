package oauthprovider

import (
	"context"
	"io"
	"net/http"
	"os"

	"github.com/umekikazuya/me/internal/app/port"
	"github.com/umekikazuya/me/pkg/errs"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

const cookieName = "authState"

var conf = &oauth2.Config{
	ClientID:     os.Getenv("CLIENT_ID"),
	ClientSecret: os.Getenv("CLIENT_SECRETS"),
	Endpoint:     github.Endpoint,
	RedirectURL:  "http://localhost:8050/auth/github/callback",
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
func (o *OauthGithubProvider) GetResource(ctx context.Context, client *http.Client) ([]byte, error) {
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, errs.WrapInternal("データ取得エラー", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errs.WrapInternal("システムエラー", err)
	}
	return body, nil
}

func New() port.OauthProvider {
	return &OauthGithubProvider{}
}

var _ port.OauthProvider = (*OauthGithubProvider)(nil)
