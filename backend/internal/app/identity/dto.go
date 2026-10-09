package identity

import "time"

type (
	InputLoginWithGithubDto  struct{}
	OutputLoginWithGithubDto struct {
		AT string
		RT string
	}
	InputCallbackFromGithubDto struct {
		Code     string
		BaseTime time.Time
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
