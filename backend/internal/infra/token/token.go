package token

import (
	"context"
	"time"

	appidentity "github.com/umekikazuya/me/internal/app/identity"
	domain "github.com/umekikazuya/me/internal/domain/identity"
)

var _ appidentity.TokenService = (*JWTTokenService)(nil)

type JWTTokenService struct {
	secret   []byte
	atExpiry time.Duration
}

func NewJWTTokenService(secret string, atExpiry time.Duration) *JWTTokenService {
	return &JWTTokenService{
		secret:   []byte(secret),
		atExpiry: atExpiry,
	}
}

func (s *JWTTokenService) GenerateAT(_ context.Context, _ domain.Identity) (string, error) {
	return "", nil
}

func (s *JWTTokenService) GenerateRT(_ context.Context) (string, error) {
	return "", nil
}

func (s *JWTTokenService) Hash(_ context.Context, _ string) (string, error) {
	return "", nil
}

func (s *JWTTokenService) Validate(_ context.Context, _ string) error {
	return nil
}
