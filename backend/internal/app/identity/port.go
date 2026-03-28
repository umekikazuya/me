package identity

import (
	"context"
	"errors"

	domain "github.com/umekikazuya/me/internal/domain/identity"
)

var (
	ErrTokenInvalid = errors.New("token is invalid")
	ErrTokenExpired = errors.New("token is expired")
	ErrTokenRevoked = errors.New("token is revoked")
)

type TokenService interface {
	// アカウント情報から新しいアクセストークンを生成
	GenerateAT(ctx context.Context, identity domain.Identity) (string, error)

	GenerateRT(ctx context.Context) (string, error)

	// トークンのハッシュ化を行う
	// token が空文字の場合はエラーを返却
	Hash(ctx context.Context, token string) (string, error)

	// アクセストークンを検証、失効・署名不正等を検出
	Validate(ctx context.Context, token string) error
}
