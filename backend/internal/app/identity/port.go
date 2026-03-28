package identity

import (
	"context"
	"time"

	domain "github.com/umekikazuya/me/internal/domain/identity"
)

type TokenService interface {
	// アカウント情報から新しいアクセストークンを生成
	Generate(ctx context.Context, account *domain.Identity, expiresIn time.Duration) (string, error)

	// トークンのハッシュ化を行う
	Hash(ctx context.Context, token string) (string, error)

	// アクセストークンを検証、失効・署名不正等を検出
	Validate(ctx context.Context, token string) error
}
