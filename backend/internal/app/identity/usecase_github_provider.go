package identity

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/umekikazuya/me/internal/app/port"
	"github.com/umekikazuya/me/internal/domain/identity"
	"github.com/umekikazuya/me/pkg/errs"
)

type usecaseGithubProviderImpl struct {
	identityRepo  identity.IdentityRepo
	oauthProvider port.OauthProvider
}

// Callback implements [githubProviderUsecase].
func (usecase *usecaseGithubProviderImpl) Callback(
	ctx context.Context,
	in InputCallbackFromGithubDto,
) error {
	client, err := usecase.oauthProvider.GetClient(ctx, in.Code)
	if err != nil {
		return errs.New(errs.ErrBadRequest, err.Error())
	}
	body, err := usecase.oauthProvider.GetResource(ctx, client)
	var values map[string]int
	json.Unmarshal(body, values)
	if err != nil {
		return errs.New(errs.ErrBadRequest, err.Error())
	}
	githubID, ok := values["id"]
	if !ok {
		return errs.WrapInternal("システムエラー", errors.New("データの取得・解析に失敗"))
	}
	exist, err := usecase.identityRepo.FindByGithubID(ctx, strconv.Itoa(githubID))
	if err != nil {
		if !errors.Is(err, errs.ErrNotFound) {
			return errs.WrapInternal("システムエラー", err)
		}
	}
	// exist がある場合は登録をしないでセッション発行
	i, err := identity.RegisterWithGithub(
		strconv.Itoa(githubID),
		in.BaseTime,
	)
	err = usecase.identityRepo.Save(ctx, i)
	if err != nil {
		return errs.WrapInternal("データベースエラー", err)
	}
	return nil
}

// Login implements [githubProviderUsecase].
func (usecase *usecaseGithubProviderImpl) Login(ctx context.Context, input InputLoginWithGithubDto) (*OutputLoginWithGithubDto, error) {
	panic("unimplemented")
}

func newGithub(
	identityRepo identity.IdentityRepo,
	oauthProvider port.OauthProvider,
) usecaseGithubProvider {
	return &usecaseGithubProviderImpl{
		oauthProvider: oauthProvider,
	}
}

var _ usecaseGithubProvider = (*usecaseGithubProviderImpl)(nil)
