package identity

// InputChangeEmailDto defines parameters for ChangeEmail.
type InputChangeEmailDto struct {
	ID              string `json:"-"`
	NewEmailAddress string `json:"newEmailAddress" validate:"required,email"`
}

// ChangeEmailParams defines parameters for ChangeEmail.
type ChangeEmailParams struct {
	// XRequestedWith CSRF 対策のためのカスタムヘッダ
	XRequestedWith ChangeEmailParamsXRequestedWith `json:"X-Requested-With"`
}

// ChangeEmailParamsXRequestedWith defines parameters for ChangeEmail.
type ChangeEmailParamsXRequestedWith string

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

// InputRegisterDto defines parameters for Register.
type InputRegisterDto struct {
	EmailAddress string `json:"emailAddress" validate:"required,email"`
	Password     string `json:"password" validate:"required,min=8,max=72"`
}

// RevokeAllSessionsParams defines parameters for RevokeAllSessions.
type RevokeAllSessionsParams struct {
	// XRequestedWith CSRF 対策のためのカスタムヘッダ
	XRequestedWith RevokeAllSessionsParamsXRequestedWith `json:"X-Requested-With"`
}

// RevokeAllSessionsParamsXRequestedWith defines parameters for RevokeAllSessions.
type RevokeAllSessionsParamsXRequestedWith string
