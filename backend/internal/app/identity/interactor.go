package identity

import (
	"context"

	"github.com/umekikazuya/me/internal/app/port"
	"github.com/umekikazuya/me/internal/domain/identity"
)

// Identity / Session のユースケース設計
type (
	usecaseGithubProvider interface {
		Callback(ctx context.Context, in InputCallbackFromGithubDto) (*OutputCallbackWithGithubDto, error)
	}
	usecaseSession interface {
		Logout(ctx context.Context, in InputLogoutDto) error
		RefreshTokens(ctx context.Context, in InputRefreshTokensDto) (*OutputRefreshTokensDto, error)
		RevokeAllSessions(ctx context.Context, in InputRevokeAllSessionsDto) error
	}
	Interactor interface {
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

func (i *interactor) Logout(ctx context.Context, in InputLogoutDto) error {
	return i.usecaseSession.Logout(ctx, in)
}

func (i *interactor) RefreshTokens(ctx context.Context, in InputRefreshTokensDto) (*OutputRefreshTokensDto, error) {
	return i.usecaseSession.RefreshTokens(ctx, in)
}

func (i *interactor) RevokeAllSessions(ctx context.Context, in InputRevokeAllSessionsDto) error {
	return i.usecaseSession.RevokeAllSessions(ctx, in)
}

// New is factory func.
//
// It receives oauth provider that like github, google.
func New(
	identityRepo identity.IdentityRepo,
	sessionRepo identity.SessionRepo,
	tokenSrv TokenService,
	oauthGithub port.OauthProvider,
) Interactor {
	return &interactor{
		usecaseSession:        newSessionUsecase(identityRepo, sessionRepo, tokenSrv),
		usecaseGithubProvider: newGithub(identityRepo, sessionRepo, oauthGithub, tokenSrv),
	}
}

var _ Interactor = (*interactor)(nil)
