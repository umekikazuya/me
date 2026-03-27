package identity

// InputChangeEmailDto defines parameters for ChangeEmail.
type InputChangeEmailDto struct {
	NewEmailAddress string `json:"newEmailAddress"`
	Token           string `json:"token"`
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
	EmailAddress string `json:"emailAddress"`
	Password     string `json:"password"`
}

// LogoutParams defines parameters for Logout.
type LogoutParams struct {
	// XRequestedWith CSRF 対策のためのカスタムヘッダ
	XRequestedWith LogoutParamsXRequestedWith `json:"X-Requested-With"`
}

// LogoutParamsXRequestedWith defines parameters for Logout.
type LogoutParamsXRequestedWith string

// InputResetPasswordDto defines parameters for ResetPassword.
type InputResetPasswordDto struct {
	NewPassword string `json:"newPassword"`
	Token       string `json:"token"`
}

// RefreshTokensParams defines parameters for RefreshTokens.
type RefreshTokensParams struct {
	// XRequestedWith CSRF 対策のためのカスタムヘッダ
	XRequestedWith RefreshTokensParamsXRequestedWith `json:"X-Requested-With"`
}

// RefreshTokensParamsXRequestedWith defines parameters for RefreshTokens.
type RefreshTokensParamsXRequestedWith string

// InputRegisterDto defines parameters for Register.
type InputRegisterDto struct {
	EmailAddress string `json:"emailAddress"`
	Password     string `json:"password"`
}

// RevokeAllSessionsParams defines parameters for RevokeAllSessions.
type RevokeAllSessionsParams struct {
	// XRequestedWith CSRF 対策のためのカスタムヘッダ
	XRequestedWith RevokeAllSessionsParamsXRequestedWith `json:"X-Requested-With"`
}

// RevokeAllSessionsParamsXRequestedWith defines parameters for RevokeAllSessions.
type RevokeAllSessionsParamsXRequestedWith string
