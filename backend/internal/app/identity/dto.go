package identity

type (
	InputLoginWithGithubDto  struct{}
	OutputLoginWithGithubDto struct {
		AT string
		RT string
	}
	InputCallbackFromGithubDto struct {
		State string
		Code  string
	}
	OutputCallbackFromGithubDto struct{}
	InputLogoutDto              struct {
		IdentityID string `json:"-"`
		RT         string `json:"-"`
	}
	InputRefreshTokensDto struct {
		IdentityID string `json:"-"`
		RT         string `json:"-"`
	}
	OutputRefreshTokensDto struct {
		AT string
		RT string
	}
	InputRevokeAllSessionsDto struct {
		IdentityID string `json:"-"`
	}
)
