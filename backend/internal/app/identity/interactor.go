package identity

import (
	"context"

	"github.com/umekikazuya/me/internal/app/port"
	domain "github.com/umekikazuya/me/internal/domain/identity"
)

// Identity / Session のユースケース設計
type (
	githubProviderUsecase interface {
		Login(ctx context.Context, input InputLoginWithGithubDto) (*OutputLoginWithGithubDto, error)
		Callback(ctx context.Context, in InputCallbackFromGithubDto) error
	}
	sessionUsecase interface {
		Logout(ctx context.Context, in InputLogoutDto) error
		RefreshTokens(ctx context.Context, in InputRefreshTokensDto) (*OutputRefreshTokensDto, error)
		RevokeAllSessions(ctx context.Context, in InputRevokeAllSessionsDto) error
	}
	Interactor interface {
		LoginWithGithub(ctx context.Context, in InputLoginWithGithubDto) (*OutputLoginWithGithubDto, error)
		CallbackFromGithub(ctx context.Context, in InputCallbackFromGithubDto) error
		Logout(ctx context.Context, in InputLogoutDto) error
		RefreshTokens(ctx context.Context, in InputRefreshTokensDto) (*OutputRefreshTokensDto, error)
		RevokeAllSessions(ctx context.Context, in InputRevokeAllSessionsDto) error
	}
	interactor struct {
		sessionUsecase        sessionUsecase
		githubProviderUsecase githubProviderUsecase
		identityRepo          domain.IdentityRepo
		sessionRepo           domain.SessionRepo
		tokenSrv              TokenService
		passwordManager       port.PasswordManager
	}
)

// CallbackFromGithub implements [Interactor].
func (i *interactor) CallbackFromGithub(ctx context.Context, in InputCallbackFromGithubDto) error {
	return i.githubProviderUsecase.Callback(ctx, in)
}

// LoginWithGithub implements [Interactor].
func (i *interactor) LoginWithGithub(ctx context.Context, in InputLoginWithGithubDto) (*OutputLoginWithGithubDto, error) {
	return i.githubProviderUsecase.Login(ctx, in)
}

func (i *interactor) Logout(ctx context.Context, in InputLogoutDto) error {
	return i.sessionUsecase.Logout(ctx, in)
}

func (i *interactor) RefreshTokens(ctx context.Context, in InputRefreshTokensDto) (*OutputRefreshTokensDto, error) {
	return i.sessionUsecase.RefreshTokens(ctx, in)
}

func (i *interactor) RevokeAllSessions(
	ctx context.Context,
	in InputRevokeAllSessionsDto,
) error {
	return i.sessionUsecase.RevokeAllSessions(ctx, in)
}

var _ Interactor = (*interactor)(nil)
