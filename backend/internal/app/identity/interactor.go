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
	Login(ctx context.Context, input InputLoginDto) (*OutputLoginDto, error)
	// ログアウトする
	Logout(ctx context.Context, input InputLogoutDto) error
	// パスワードをリセットする
	ResetPassword(ctx context.Context, input InputResetPasswordDto) error
	// トークンをリフレッシュする
	RefreshTokens(ctx context.Context, input InputRefreshTokensDto) (*OutputRefreshTokensDto, error)
	// Identityを登録する
	Register(ctx context.Context, input InputRegisterDto) error
	// 全RTを失効させる
	RevokeAllSessions(ctx context.Context, input InputRevokeAllSessionsDto) error
}

type Interactor struct {
	identityRepo domain.IdentityRepo
	sessionRepo  domain.SessionRepo
	tokenSrv     TokenService
}

func NewInteractor(
	identityRepo domain.IdentityRepo,
	sessionRepo domain.SessionRepo,
	tokenSrv TokenService,
) interactor {
	return &Interactor{
		identityRepo: identityRepo,
		sessionRepo:  sessionRepo,
		tokenSrv:     tokenSrv,
	}
}

func (i *Interactor) ChangeEmail(ctx context.Context, input InputChangeEmailDto) error {
	idn, err := i.identityRepo.FindByID(ctx, input.ID)
	if err != nil {
		return err
	}
	if idn == nil {
		return fmt.Errorf("ChangeEmail: %w", errs.ErrNotFound)
	}
	// メールの重複チェック
	newEmail, err := domain.NewEmail(input.NewEmailAddress)
	if err != nil {
		return err
	}
	exists, err := i.identityRepo.FindByEmail(ctx, newEmail)
	if err != nil {
		return err
	}
	if exists != nil {
		return fmt.Errorf("ChangeEmail: %w", errs.ErrConflict)
	}
	err = idn.ChangeEmail(input.NewEmailAddress)
	if err != nil {
		return err
	}
	err = i.identityRepo.Save(ctx, idn)
	if err != nil {
		return err
	}
	return nil
}

func (i *Interactor) Login(ctx context.Context, input InputLoginDto) (*OutputLoginDto, error) {
	// メール検証
	email, err := domain.NewEmail(input.EmailAddress)
	if err != nil {
		return nil, err
	}
	// 入力されたメールアドレスでアカウントを検索
	idn, err := i.identityRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if idn == nil {
		return nil, fmt.Errorf("Login: %w", errs.ErrNotFound)
	}
	// 認証
	err = idn.Authenticate(input.Password)
	if err != nil {
		return nil, err
	}

	at, err := i.tokenSrv.GenerateAT(ctx, *idn)
	if err != nil {
		return nil, err
	}
	rt, err := i.tokenSrv.GenerateRT(ctx)
	if err != nil {
		return nil, err
	}
	hashedRT, err := i.tokenSrv.Hash(ctx, rt)
	if err != nil {
		return nil, err
	}
	ses, err := idn.CreateSession(hashedRT)
	if err != nil {
		return nil, err
	}
	err = i.sessionRepo.Save(ctx, ses) // TODO: アクティブセッション数の制限制御
	if err != nil {
		return nil, err
	}

	return &OutputLoginDto{
		AT: at,
		RT: rt,
	}, nil
}

func (i *Interactor) Logout(ctx context.Context, input InputLogoutDto) error {
	idn, err := i.identityRepo.FindByID(ctx, input.IdentityID)
	if err != nil {
		return err
	}
	if idn == nil {
		return fmt.Errorf("Logout: %w", errs.ErrNotFound)
	}
	hashedRT, err := i.tokenSrv.Hash(ctx, input.RT)
	if err != nil {
		return err
	}
	ses, err := i.sessionRepo.FindByIdentityIdAndTokenHash(ctx, idn.ID(), hashedRT)
	if err != nil {
		return err
	}
	if ses == nil {
		return fmt.Errorf("Logout %w", errs.ErrNotFound)
	}
	// セッションの無効化処理
	err = ses.Revoke()
	if err != nil {
		return err
	}
	err = i.sessionRepo.Save(ctx, ses)
	if err != nil {
		return err
	}

	return nil
}

func (i *Interactor) ResetPassword(ctx context.Context, input InputResetPasswordDto) error {
	idn, err := i.identityRepo.FindByID(ctx, input.ID)
	if err != nil {
		return err
	}
	if idn == nil {
		return fmt.Errorf("ResetPassword: %w", errs.ErrNotFound)
	}
	err = idn.ResetPassword(input.NewPassword)
	if err != nil {
		return err
	}
	err = i.identityRepo.Save(ctx, idn)
	if err != nil {
		return err
	}
	err = i.sessionRepo.RevokeAll(ctx, idn.ID())
	if err != nil {
		return err
	}
	return nil
}

func (i *Interactor) RefreshTokens(ctx context.Context, input InputRefreshTokensDto) (*OutputRefreshTokensDto, error) {
	idn, err := i.identityRepo.FindByID(ctx, input.IdentityID)
	if err != nil {
		return nil, err
	}
	if idn == nil {
		return nil, fmt.Errorf("RefreshTokens: %w", errs.ErrNotFound)
	}
	hashedRT, err := i.tokenSrv.Hash(ctx, input.RT)
	if err != nil {
		return nil, err
	}
	ses, err := i.sessionRepo.FindByIdentityIdAndTokenHash(ctx, idn.ID(), hashedRT)
	if err != nil {
		return nil, err
	}
	if ses == nil {
		return nil, fmt.Errorf("RefreshTokens: sessionが存在しません %w", errs.ErrNotFound)
	}
	if ses.Status() != "active" { // TODO: IsActive or IsRevoke を実装する
		return nil, fmt.Errorf("RefreshTokens: RTが失効済みです %w", errs.ErrUnprocessable)
	}

	newAT, err := i.tokenSrv.GenerateAT(ctx, *idn)
	if err != nil {
		return nil, fmt.Errorf("RefreshTokens: %w", err)
	}
	newRT, err := i.tokenSrv.GenerateRT(ctx)
	if err != nil {
		return nil, fmt.Errorf("RefreshTokens: %w", err)
	}
	newHashedRT, err := i.tokenSrv.Hash(ctx, newRT)
	if err != nil {
		return nil, fmt.Errorf("RefreshTokens: %w", err)
	}
	newSes, err := ses.Rotate(newHashedRT)
	if err != nil {
		return nil, err
	}
	err = i.sessionRepo.Save(ctx, ses)
	if err != nil {
		return nil, err
	}
	err = i.sessionRepo.Save(ctx, newSes)
	if err != nil {
		return nil, err
	}

	return &OutputRefreshTokensDto{
		AT: newAT,
		RT: newRT,
	}, nil
}

// Register は認証プロファイルの登録処理
func (i *Interactor) Register(ctx context.Context, input InputRegisterDto) error {
	email, err := domain.NewEmail(input.EmailAddress)
	if err != nil {
		return err
	}
	exists, err := i.identityRepo.FindByEmail(ctx, email)
	if err != nil {
		return err
	}
	if exists != nil {
		return fmt.Errorf("Register: %w", errs.ErrConflict)
	}
	e, err := domain.Register(
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

func (i *Interactor) RevokeAllSessions(ctx context.Context, input InputRevokeAllSessionsDto) error {
	idn, err := i.identityRepo.FindByID(ctx, input.IdentityID)
	if err != nil {
		return err
	}
	if idn == nil {
		return fmt.Errorf("RevokeAllSessions: %w", errs.ErrNotFound)
	}
	err = i.sessionRepo.RevokeAll(ctx, idn.ID())
	if err != nil {
		return err
	}
	return nil
}
