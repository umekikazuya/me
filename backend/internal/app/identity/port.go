package identity

import (
	"context"
	"errors"
	"time"

	domain "github.com/umekikazuya/me/internal/domain/identity"
)

var (
	ErrTokenInvalid = errors.New("token is invalid")
	ErrTokenExpired = errors.New("token is expired")
	ErrTokenRevoked = errors.New("token is revoked")
)

type TokenService interface {
	// アカウント情報から新しいアクセストークンを生成
	Generate(ctx context.Context, identity domain.Identity, expiresIn time.Duration) (string, error)

	// トークンのハッシュ化を行う
	Hash(ctx context.Context, token string) (string, error)

	// アクセストークンを検証、失効・署名不正等を検出
	Validate(ctx context.Context, token string) error
}
