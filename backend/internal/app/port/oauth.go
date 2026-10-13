package port

import (
	"context"
	"net/http"
)

type OauthProvider interface {
	AuthCodeURL(state string) string
	GetClient(ctx context.Context, code string) (*http.Client, error)
	GetResource(ctx context.Context, client *http.Client) ([]byte, error)
}
