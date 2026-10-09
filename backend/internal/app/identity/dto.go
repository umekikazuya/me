package identity

import "time"

type (
	InputLoginWithGithubDto    struct{}
	InputCallbackFromGithubDto struct {
		Code     string
		BaseTime time.Time
	}
	OutputCallbackWithGithubDto struct {
		AT string
		RT string
	}
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
