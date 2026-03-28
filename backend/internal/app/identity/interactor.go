package identity

import (
	"context"

	domain "github.com/umekikazuya/me/internal/domain/identity"
)

var _ interactor = (*Interactor)(nil)

// Identity / Session のユースケース設計
type interactor interface {
	// メールアドレスを変更する
	ChangeEmail(ctx context.Context, input InputChangeEmailDto) error
	// ログインする
	Login(ctx context.Context, input InputLoginDto) error
	// ログアウトする
	Logout(ctx context.Context) error
	// パスワードをリセットする
	ResetPassword(ctx context.Context, input InputResetPasswordDto) error
	// トークンをリフレッシュする
	RefreshTokens(ctx context.Context) error
	// Identityを登録する
	Register(ctx context.Context, input InputRegisterDto) error
	// 全RTを失効させる
	RevokeAllSessions(ctx context.Context) error
}

type Interactor struct {
	identityRepo domain.IdentityRepo
	sessionRepo  domain.SessionRepo
}

func (i *Interactor) ChangeEmail(ctx context.Context, input InputChangeEmailDto) error {
	return nil
}

func (i *Interactor) Login(ctx context.Context, input InputLoginDto) error {
	return nil
}

func (i *Interactor) Logout(ctx context.Context) error {
	return nil
}

func (i *Interactor) ResetPassword(ctx context.Context, input InputResetPasswordDto) error {
	return nil
}

func (i *Interactor) RefreshTokens(ctx context.Context) error {
	return nil
}

func (i *Interactor) Register(ctx context.Context, input InputRegisterDto) error {

	return nil
}

func (i *Interactor) RevokeAllSessions(ctx context.Context) error {
	return nil
}
