package identity

import (
	"context"

	"github.com/umekikazuya/me/internal/domain/identity"
	oauthprovider "github.com/umekikazuya/me/internal/infra/oauth_provider"
)

// Identity / Session のユースケース設計
type (
	usecaseGithubProvider interface {
		Login(ctx context.Context, input InputLoginWithGithubDto) error
		Callback(ctx context.Context, in InputCallbackFromGithubDto) (*OutputCallbackWithGithubDto, error)
	}
	usecaseSession interface {
		Logout(ctx context.Context, in InputLogoutDto) error
		RefreshTokens(ctx context.Context, in InputRefreshTokensDto) (*OutputRefreshTokensDto, error)
		RevokeAllSessions(ctx context.Context, in InputRevokeAllSessionsDto) error
	}
	Interactor interface {
		LoginWithGithub(ctx context.Context, in InputLoginWithGithubDto) error
		CallbackFromGithub(ctx context.Context, in InputCallbackFromGithubDto) (*OutputCallbackWithGithubDto, error)
		Logout(ctx context.Context, in InputLogoutDto) error
		RefreshTokens(ctx context.Context, in InputRefreshTokensDto) (*OutputRefreshTokensDto, error)
		RevokeAllSessions(ctx context.Context, in InputRevokeAllSessionsDto) error
	}
	interactor struct {
		usecaseSession        usecaseSession
		usecaseGithubProvider usecaseGithubProvider
	}
)

// CallbackFromGithub implements [Interactor].
func (i *interactor) CallbackFromGithub(
	ctx context.Context,
	in InputCallbackFromGithubDto,
) (
	*OutputCallbackWithGithubDto,
	error,
) {
	return i.usecaseGithubProvider.Callback(ctx, in)
}

// LoginWithGithub implements [Interactor].
func (i *interactor) LoginWithGithub(ctx context.Context, in InputLoginWithGithubDto) error {
	return i.usecaseGithubProvider.Login(ctx, in)
}

func (i *interactor) Logout(ctx context.Context, in InputLogoutDto) error {
	return i.usecaseSession.Logout(ctx, in)
}

func (i *interactor) RefreshTokens(ctx context.Context, in InputRefreshTokensDto) (*OutputRefreshTokensDto, error) {
	return i.usecaseSession.RefreshTokens(ctx, in)
}

func (i *interactor) RevokeAllSessions(ctx context.Context, in InputRevokeAllSessionsDto) error {
	return i.usecaseSession.RevokeAllSessions(ctx, in)
}

func New(
	identityRepo identity.IdentityRepo,
	sessionRepo identity.SessionRepo,
	tokenSrv TokenService,
) Interactor {
	return &interactor{
		usecaseSession:        newSessionUsecase(identityRepo, sessionRepo, tokenSrv),
		usecaseGithubProvider: newGithub(identityRepo, sessionRepo, oauthprovider.New(), tokenSrv),
	}
}

var _ Interactor = (*interactor)(nil)
