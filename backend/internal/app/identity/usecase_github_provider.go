package identity

import "context"

type usecaseGithubProviderImpl struct{}

// Callback implements [githubProviderUsecase].
func (usecase *usecaseGithubProviderImpl) Callback(ctx context.Context, in InputCallbackFromGithubDto) error {
	panic("unimplemented")
}

// Login implements [githubProviderUsecase].
func (usecase *usecaseGithubProviderImpl) Login(ctx context.Context, input InputLoginWithGithubDto) (*OutputLoginWithGithubDto, error) {
	panic("unimplemented")
}

func newGithub() usecaseGithubProvider {
	return &usecaseGithubProviderImpl{}
}

var _ usecaseGithubProvider = (*usecaseGithubProviderImpl)(nil)
