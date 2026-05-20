package password

import "errors"

var (
	ErrNoMatchPassword           = errors.New("パスワードが一致していません")
	ErrInvalidHashFormat         = errors.New("invalid argon2 hash format")
	ErrIncompatibleVersion       = errors.New("incompatible argon2 version")
	ErrMismatchedHashAndPassword = errors.New(
		"hashed password is not the hash of the given password",
	)
)
