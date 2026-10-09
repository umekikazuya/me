package port

import (
	"context"
	"net/http"
)

type PasswordHasher interface {
	Hash(ctx context.Context, input string) ([]byte, error)
}

type PasswordVerifier interface {
	Verify(ctx context.Context, hashedPassword, plainPassword string) error
}

type PasswordManager interface {
	PasswordHasher
	PasswordVerifier
}

type OauthProvider interface {
	GetClient(ctx context.Context, code string) (*http.Client, error)
	GetResource(ctx context.Context, client *http.Client) ([]byte, error)
}
