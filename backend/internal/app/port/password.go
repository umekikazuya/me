package port

import "context"

type PasswordHasher interface {
	Hash(ctx context.Context, input string) ([]byte, error)
}

type PasswordVerifier interface {
	Verify(ctx context.Context, hashedPassword, plainPassword string) error
}

type PasswordManager interface {
	PasswordHasher
	PasswordVerifier
}
