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

	// AT を検証し、成功時は identityID (UUID 文字列) を返す
	// 失効・署名不正・期限切れ等の場合はエラーを返す
	ValidateAT(ctx context.Context, token string) (string, error)
}
