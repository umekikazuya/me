package identity

import (
	"context"
	"fmt"

	appevent "github.com/umekikazuya/me/internal/app/event"
	domain "github.com/umekikazuya/me/internal/domain/identity"
	"github.com/umekikazuya/me/pkg/errs"
)

var _ Interactor = (*interactor)(nil)

// Identity / Session のユースケース設計
type Interactor interface {
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

type interactor struct {
	identityRepo domain.IdentityRepo
	sessionRepo  domain.SessionRepo
	tokenSrv     TokenService
	dispatcher   appevent.EventDispatcher
}

func NewInteractor(
	identityRepo domain.IdentityRepo,
	sessionRepo domain.SessionRepo,
	tokenSrv TokenService,
	dispatcher appevent.EventDispatcher,
) Interactor {
	return &interactor{
		identityRepo: identityRepo,
		sessionRepo:  sessionRepo,
		tokenSrv:     tokenSrv,
		dispatcher:   dispatcher,
	}
}

func (i *interactor) ChangeEmail(ctx context.Context, input InputChangeEmailDto) error {
	idn, err := i.identityRepo.FindByID(ctx, input.ID)
	if err != nil {
		return errs.WrapInternal("identity.identityRepo.FindByID", err)
	}
	if idn == nil {
		return fmt.Errorf("ChangeEmail: %w", errs.ErrNotFound)
	}
	// メールの重複チェック
	newEmail, err := domain.NewEmail(input.NewEmailAddress)
	if err != nil {
		return err
	}
	exists, err := i.identityRepo.FindByEmail(ctx, newEmail.Value())
	if err != nil {
		return errs.WrapInternal("identity.identityRepo.FindByEmail", err)
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
		return errs.WrapInternal("identity.identityRepo.Save", err)
	}
	if err = i.dispatcher.Dispatch(ctx, idn.Events()); err != nil {
		return errs.WrapInternal("identity.dispatcher.Dispatch", err)
	}
	idn.ClearEvents()
	return nil
}

func (i *interactor) Login(ctx context.Context, input InputLoginDto) (*OutputLoginDto, error) {
	// メール検証
	email, err := domain.NewEmail(input.EmailAddress)
	if err != nil {
		return nil, err
	}
	// 入力されたメールアドレスでアカウントを検索
	idn, err := i.identityRepo.FindByEmail(ctx, email.Value())
	if err != nil {
		return nil, errs.WrapInternal("identity.identityRepo.FindByEmail", err)
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
		return nil, errs.WrapInternal("identity.tokenSrv.GenerateAT", err)
	}
	rt, err := i.tokenSrv.GenerateRT(ctx)
	if err != nil {
		return nil, errs.WrapInternal("identity.tokenSrv.GenerateRT", err)
	}
	hashedRT, err := i.tokenSrv.Hash(ctx, rt)
	if err != nil {
		return nil, errs.WrapInternal("identity.tokenSrv.Hash", err)
	}
	ses, err := idn.CreateSession(hashedRT)
	if err != nil {
		return nil, err
	}
	err = i.sessionRepo.Save(ctx, ses) // TODO: アクティブセッション数の制限制御
	if err != nil {
		return nil, errs.WrapInternal("identity.sessionRepo.Save", err)
	}
	if err = i.dispatcher.Dispatch(ctx, idn.Events()); err != nil { // TODO: 原子性の対応
		return nil, errs.WrapInternal("identity.dispatcher.Dispatch", err)
	}
	idn.ClearEvents()

	return &OutputLoginDto{
		AT: at,
		RT: rt,
	}, nil
}

func (i *interactor) Logout(ctx context.Context, input InputLogoutDto) error {
	idn, err := i.identityRepo.FindByID(ctx, input.IdentityID)
	if err != nil {
		return errs.WrapInternal("identity.identityRepo.FindByID", err)
	}
	if idn == nil {
		return fmt.Errorf("Logout: %w", errs.ErrNotFound)
	}
	hashedRT, err := i.tokenSrv.Hash(ctx, input.RT)
	if err != nil {
		return errs.WrapInternal("identity.tokenSrv.Hash", err)
	}
	ses, err := i.sessionRepo.FindByIdentityIdAndTokenHash(ctx, idn.ID(), hashedRT)
	if err != nil {
		return errs.WrapInternal("identity.sessionRepo.FindByIdentityIdAndTokenHash", err)
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
		return errs.WrapInternal("identity.sessionRepo.Save", err)
	}
	if err = i.dispatcher.Dispatch(ctx, ses.Events()); err != nil {
		return errs.WrapInternal("identity.dispatcher.Dispatch", err)
	}
	ses.ClearEvents()
	return nil
}

func (i *interactor) ResetPassword(ctx context.Context, input InputResetPasswordDto) error {
	idn, err := i.identityRepo.FindByID(ctx, input.ID)
	if err != nil {
		return errs.WrapInternal("identity.identityRepo.FindByID", err)
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
		return errs.WrapInternal("identity.identityRepo.Save", err)
	}
	err = i.sessionRepo.RevokeAll(ctx, idn.ID())
	if err != nil {
		return errs.WrapInternal("identity.sessionRepo.RevokeAll", err)
	}
	if err = i.dispatcher.Dispatch(ctx, idn.Events()); err != nil {
		return errs.WrapInternal("identity.dispatcher.Dispatch", err)
	}
	idn.ClearEvents()
	return nil
}

func (i *interactor) RefreshTokens(ctx context.Context, input InputRefreshTokensDto) (*OutputRefreshTokensDto, error) {
	idn, err := i.identityRepo.FindByID(ctx, input.IdentityID)
	if err != nil {
		return nil, errs.WrapInternal("identity.identityRepo.FindByID", err)
	}
	if idn == nil {
		return nil, fmt.Errorf("RefreshTokens: %w", errs.ErrNotFound)
	}
	hashedRT, err := i.tokenSrv.Hash(ctx, input.RT)
	if err != nil {
		return nil, errs.WrapInternal("identity.tokenSrv.Hash", err)
	}
	ses, err := i.sessionRepo.FindByIdentityIdAndTokenHash(ctx, idn.ID(), hashedRT)
	if err != nil {
		return nil, errs.WrapInternal("identity.sessionRepo.FindByIdentityIdAndTokenHash", err)
	}
	if ses == nil {
		return nil, fmt.Errorf("RefreshTokens: sessionが存在しません %w", errs.ErrNotFound)
	}
	// TODO: IsActive or IsRevoke を実装する
	// TODO: return e.Status() == "active" && time.Now().Before(e.ExpiresAt())
	if ses.Status() != "active" {
		return nil, fmt.Errorf("RefreshTokens: RTが失効済みです %w", errs.ErrUnprocessable)
	}

	newAT, err := i.tokenSrv.GenerateAT(ctx, *idn)
	if err != nil {
		return nil, errs.WrapInternal("identity.tokenSrv.GenerateAT", err)
	}
	newRT, err := i.tokenSrv.GenerateRT(ctx)
	if err != nil {
		return nil, errs.WrapInternal("identity.tokenSrv.GenerateRT", err)
	}
	newHashedRT, err := i.tokenSrv.Hash(ctx, newRT)
	if err != nil {
		return nil, errs.WrapInternal("identity.tokenSrv.Hash", err)
	}
	newSes, err := ses.Rotate(newHashedRT)
	if err != nil {
		return nil, err
	}
	err = i.sessionRepo.Save(ctx, ses)
	if err != nil {
		return nil, errs.WrapInternal("identity.sessionRepo.Save", err)
	}
	err = i.sessionRepo.Save(ctx, newSes)
	if err != nil {
		return nil, errs.WrapInternal("identity.sessionRepo.Save", err)
	}
	if err = i.dispatcher.Dispatch(ctx, ses.Events()); err != nil {
		return nil, errs.WrapInternal("identity.dispatcher.Dispatch", err)
	}
	ses.ClearEvents()

	return &OutputRefreshTokensDto{
		AT: newAT,
		RT: newRT,
	}, nil
}

// Register は認証プロファイルの登録処理
func (i *interactor) Register(ctx context.Context, input InputRegisterDto) error {
	email, err := domain.NewEmail(input.EmailAddress)
	if err != nil {
		return err
	}
	exists, err := i.identityRepo.FindByEmail(ctx, email.Value())
	if err != nil {
		return errs.WrapInternal("identity.identityRepo.FindByEmail", err)
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
		return errs.WrapInternal("identity.identityRepo.Save", err)
	}
	if err = i.dispatcher.Dispatch(ctx, e.Events()); err != nil {
		return errs.WrapInternal("identity.dispatcher.Dispatch", err)
	}
	e.ClearEvents()
	return nil
}

func (i *interactor) RevokeAllSessions(ctx context.Context, input InputRevokeAllSessionsDto) error {
	idn, err := i.identityRepo.FindByID(ctx, input.IdentityID)
	if err != nil {
		return errs.WrapInternal("identity.identityRepo.FindByID", err)
	}
	if idn == nil {
		return fmt.Errorf("RevokeAllSessions: %w", errs.ErrNotFound)
	}
	err = i.sessionRepo.RevokeAll(ctx, idn.ID())
	if err != nil {
		return errs.WrapInternal("identity.sessionRepo.RevokeAll", err)
	}
	return nil
}
