package identity

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/umekikazuya/me/internal/app/port"
	domain "github.com/umekikazuya/me/internal/domain/identity"
	"github.com/umekikazuya/me/pkg/errs"
)

type usecaseGithubProviderImpl struct {
	identityRepo  domain.IdentityRepo
	sessionRepo   domain.SessionRepo
	oauthProvider port.OauthProvider
	tokenSrv      TokenService
}

// Callback implements [githubProviderUsecase].
func (usecase *usecaseGithubProviderImpl) Callback(
	ctx context.Context,
	in InputCallbackFromGithubDto,
) (*OutputCallbackWithGithubDto, error) {
	client, err := usecase.oauthProvider.GetClient(ctx, in.Code)
	if err != nil {
		return nil, errs.New(errs.ErrBadRequest, err.Error())
	}
	body, err := usecase.oauthProvider.GetResource(ctx, client)
	if err != nil {
		return nil, errs.New(errs.ErrBadRequest, err.Error())
	}
	var user struct {
		ID *int64 `json:"id"`
	}
	if err = json.Unmarshal(body, &user); err != nil || user.ID == nil {
		return nil, errs.New(errs.ErrBadRequest, err.Error())
	}
	githubID := strconv.FormatInt(*user.ID, 10)
	i, err := usecase.identityRepo.FindByGithubID(
		ctx,
		githubID,
	)
	if err != nil {
		if !errors.Is(err, errs.ErrNotFound) {
			return nil, errs.WrapInternal("システムエラー", err)
		}
		i, err = domain.RegisterWithGithub(
			githubID,
			in.BaseTime,
		)
		if err != nil {
			return nil, errs.New(errs.ErrConflict, err.Error())
		}
		err = usecase.identityRepo.Save(ctx, i)
		if err != nil {
			return nil, errs.WrapInternal("データベースエラー", err)
		}
	}
	return usecase.issueSession(ctx, i)
}

func (usecase *usecaseGithubProviderImpl) issueSession(
	ctx context.Context,
	i *domain.Account,
) (*OutputCallbackWithGithubDto, error) {
	at, err := usecase.tokenSrv.GenerateAT(ctx, *i)
	if err != nil {
		return &OutputCallbackWithGithubDto{}, errs.WrapInternal("identity.tokenSrv.GenerateAT", err)
	}
	rt, err := usecase.tokenSrv.GenerateRT(ctx)
	if err != nil {
		return &OutputCallbackWithGithubDto{}, errs.WrapInternal("identity.tokenSrv.GenerateRT", err)
	}
	hashedRT, err := usecase.tokenSrv.Hash(ctx, rt)
	if err != nil {
		return &OutputCallbackWithGithubDto{}, errs.WrapInternal("identity.tokenSrv.Hash", err)
	}
	ses, err := i.CreateSession(hashedRT)
	if err != nil {
		return &OutputCallbackWithGithubDto{}, err
	}
	err = usecase.sessionRepo.Save(ctx, ses) // TODO: アクティブセッション数の制限制御
	if err != nil {
		return &OutputCallbackWithGithubDto{}, errs.WrapInternal("identity.sessionRepo.Save", err)
	}
	return &OutputCallbackWithGithubDto{AT: at, RT: rt}, nil
}

func newGithub(
	identityRepo domain.IdentityRepo,
	sessionRepo domain.SessionRepo,
	oauthProvider port.OauthProvider,
	tokenSrv TokenService,
) usecaseGithubProvider {
	return &usecaseGithubProviderImpl{
		identityRepo:  identityRepo,
		sessionRepo:   sessionRepo,
		oauthProvider: oauthProvider,
		tokenSrv:      tokenSrv,
	}
}

var _ usecaseGithubProvider = (*usecaseGithubProviderImpl)(nil)
