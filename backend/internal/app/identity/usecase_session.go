package identity

import (
	"context"
	"fmt"

	"github.com/umekikazuya/me/internal/domain/identity"
	domain "github.com/umekikazuya/me/internal/domain/identity"
	"github.com/umekikazuya/me/pkg/errs"
)

type sessionUsecaseImpl struct {
	identityRepo domain.IdentityRepo
	sessionRepo  domain.SessionRepo
	tokenSrv     TokenService
}

// Logout implements [sessionUsecase].
func (usecase *sessionUsecaseImpl) Logout(ctx context.Context, in InputLogoutDto) error {
	idn, err := usecase.identityRepo.FindByID(ctx, in.IdentityID)
	if err != nil {
		return errs.WrapInternal("identity.identityRepo.FindByID", err)
	}
	if idn == nil {
		return fmt.Errorf("Logout: %w", errs.ErrNotFound)
	}
	hashedRT, err := usecase.tokenSrv.Hash(ctx, in.RT)
	if err != nil {
		return errs.WrapInternal("identity.tokenSrv.Hash", err)
	}
	ses, err := usecase.sessionRepo.FindByIdentityIdAndTokenHash(ctx, idn.ID(), hashedRT)
	if err != nil {
		return errs.WrapInternal("identity.sessionRepo.FindByIdentityIdAndTokenHash", err)
	}
	if ses == nil {
		return fmt.Errorf("Logout %w", errs.ErrNotFound)
	}
	err = ses.Revoke()
	if err != nil {
		return err
	}
	err = usecase.sessionRepo.Save(ctx, ses)
	if err != nil {
		return errs.WrapInternal("identity.sessionRepo.Save", err)
	}
	return nil
}

// RefreshTokens implements [sessionUsecase].
func (usecase *sessionUsecaseImpl) RefreshTokens(ctx context.Context, in InputRefreshTokensDto) (*OutputRefreshTokensDto, error) {
	idn, err := usecase.identityRepo.FindByID(ctx, in.IdentityID)
	if err != nil {
		return nil, errs.WrapInternal("identity.identityRepo.FindByID", err)
	}
	if idn == nil {
		return nil, fmt.Errorf("RefreshTokens: %w", errs.ErrNotFound)
	}
	hashedRT, err := usecase.tokenSrv.Hash(ctx, in.RT)
	if err != nil {
		return nil, errs.WrapInternal("identity.tokenSrv.Hash", err)
	}
	ses, err := usecase.sessionRepo.FindByIdentityIdAndTokenHash(ctx, in.IdentityID, hashedRT)
	if err != nil {
		return nil, errs.WrapInternal("identity.sessionRepo.FindByIdentityIdAndTokenHash", err)
	}
	if ses == nil {
		return nil, fmt.Errorf("RefreshTokens: sessionが存在しません %w", errs.ErrNotFound)
	}
	if !ses.IsActive() {
		return nil, errs.New(errs.ErrConflict, "RefreshTokens: RTが失効済みです")
	}

	newAT, err := usecase.tokenSrv.GenerateAT(ctx, *idn)
	if err != nil {
		return nil, errs.WrapInternal("identity.tokenSrv.GenerateAT", err)
	}
	newRT, err := usecase.tokenSrv.GenerateRT(ctx)
	if err != nil {
		return nil, errs.WrapInternal("identity.tokenSrv.GenerateRT", err)
	}
	newHashedRT, err := usecase.tokenSrv.Hash(ctx, newRT)
	if err != nil {
		return nil, errs.WrapInternal("identity.tokenSrv.Hash", err)
	}
	newSes, err := ses.Rotate(newHashedRT)
	if err != nil {
		return nil, err
	}
	err = usecase.sessionRepo.Save(ctx, ses)
	if err != nil {
		return nil, errs.WrapInternal("identity.sessionRepo.Save", err)
	}
	err = usecase.sessionRepo.Save(ctx, newSes)
	if err != nil {
		return nil, errs.WrapInternal("identity.sessionRepo.Save", err)
	}
	return &OutputRefreshTokensDto{
		AT: newAT,
		RT: newRT,
	}, nil
}

// RevokeAllSessions implements [sessionUsecase].
func (usecase *sessionUsecaseImpl) RevokeAllSessions(ctx context.Context, in InputRevokeAllSessionsDto) error {
	idn, err := usecase.identityRepo.FindByID(ctx, in.IdentityID)
	if err != nil {
		return errs.WrapInternal("identity.identityRepo.FindByID", err)
	}
	if idn == nil {
		return fmt.Errorf("RevokeAllSessions: %w", errs.ErrNotFound)
	}
	err = usecase.sessionRepo.RevokeAll(ctx, idn.ID())
	if err != nil {
		return errs.WrapInternal("identity.sessionRepo.RevokeAll", err)
	}
	return nil
}

func newSessionUsecase(
	identityRepo identity.IdentityRepo,
	sessionRepo identity.SessionRepo,
	tokenSrv TokenService,
) usecaseSession {
	return &sessionUsecaseImpl{
		identityRepo: identityRepo,
		sessionRepo:  sessionRepo,
		tokenSrv:     tokenSrv,
	}
}

var _ usecaseSession = (*sessionUsecaseImpl)(nil)
