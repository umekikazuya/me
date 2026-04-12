package identity

// InputChangeEmailDto defines parameters for ChangeEmail.
type InputChangeEmailDto struct {
	ID              string `json:"-"`
	NewEmailAddress string `json:"newEmailAddress" validate:"required,email"`
}

// InputLoginDto defines parameters for Login.
type InputLoginDto struct {
	EmailAddress string `json:"emailAddress" validate:"required,email"`
	Password     string `json:"password" validate:"required,min=8,max=72"`
}

type InputLogoutDto struct {
	IdentityID string `json:"-"`
	RT         string `json:"-"`
}

type OutputLoginDto struct {
	AT string
	RT string
}

// InputResetPasswordDto defines parameters for ResetPassword.
type InputResetPasswordDto struct {
	ID          string `json:"-"`
	NewPassword string `json:"newPassword" validate:"required,min=8,max=72"`
}

type InputRefreshTokensDto struct {
	IdentityID string `json:"-"`
	RT         string `json:"-"`
}

type OutputRefreshTokensDto struct {
	AT string
	RT string
}

// InputRegisterDto defines parameters for Register.
type InputRegisterDto struct {
	EmailAddress string `json:"emailAddress" validate:"required,email"`
	Password     string `json:"password" validate:"required,min=8,max=72"`
}

type InputRevokeAllSessionsDto struct {
	IdentityID string `json:"-"`
}
