package password

import "errors"

var (
	ErrInvalidHashFormat         = errors.New("invalid argon2 hash format")
	ErrIncompatibleVersion       = errors.New("incompatible argon2 version")
	ErrMismatchedHashAndPassword = errors.New(
		"パスワードが一致していません",
	)
)
