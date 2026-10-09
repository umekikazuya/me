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
	var values map[string]int
	json.Unmarshal(body, values)
	if err != nil {
		return nil, errs.New(errs.ErrBadRequest, err.Error())
	}
	githubID, ok := values["id"]
	if !ok {
		return nil, errs.WrapInternal("システムエラー", errors.New("データの取得・解析に失敗"))
	}
	i, err := usecase.identityRepo.FindByGithubID(ctx, strconv.Itoa(githubID))
	if err != nil {
		if !errors.Is(err, errs.ErrNotFound) {
			return nil, errs.WrapInternal("システムエラー", err)
		}
		i, err := domain.RegisterWithGithub(
			strconv.Itoa(githubID),
			in.BaseTime,
		)
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

// Login implements [githubProviderUsecase].
func (usecase *usecaseGithubProviderImpl) Login(ctx context.Context, input InputLoginWithGithubDto) (*OutputCallbackWithGithubDto, error) {
	panic("unimplemented")
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
