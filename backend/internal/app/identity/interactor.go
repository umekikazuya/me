package identity

import (
	"context"
	"fmt"

	domain "github.com/umekikazuya/me/internal/domain/identity"
	"github.com/umekikazuya/me/pkg/errs"
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
	tokenSrv     TokenService
}

func (i *Interactor) ChangeEmail(ctx context.Context, input InputChangeEmailDto) error {
	return nil
}

func (i *Interactor) Login(ctx context.Context, input InputLoginDto) error {
	// メール検証
	email, err := domain.NewEmail(input.EmailAddress)
	if err != nil {
		return err
	}
	// 入力されたメールアドレスでアカウントを検索
	entity, err := i.identityRepo.FindByEmail(ctx, email)
	if err != nil {
		return err
	}
	if entity == nil {
		return fmt.Errorf("Login: %w", errs.ErrNotFound)
	}
	// 認証
	err = entity.Authenticate(input.Password)
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

// Register は認証プロファイルの登録処理
func (i *Interactor) Register(ctx context.Context, input InputRegisterDto) error {
	e, err := domain.NewIdentity(
		input.EmailAddress,
		input.Password,
	)
	if err != nil {
		return err
	}
	err = i.identityRepo.Save(ctx, e)
	if err != nil {
		return err
	}
	return nil
}

func (i *Interactor) RevokeAllSessions(ctx context.Context) error {
	return nil
}
