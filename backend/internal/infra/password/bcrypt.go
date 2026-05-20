package password

import (
	"context"

	"github.com/umekikazuya/me/internal/app/port"
	"golang.org/x/crypto/bcrypt"
)

type BcryptPasswordManager struct{}

// Hash implements [port.PasswordManager].
func (b *BcryptPasswordManager) Hash(ctx context.Context, input string) ([]byte, error) {
	h, err := bcrypt.GenerateFromPassword(
		[]byte(input), bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}
	return h, nil
}

// Verify implements [port.PasswordManager].
func (b *BcryptPasswordManager) Verify(
	ctx context.Context,
	hashedPassword, plainPassword string,
) error {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(plainPassword),
	)
	if err != nil {
		return ErrNoMatchPassword
	}
	return nil
}

var _ port.PasswordManager = (*BcryptPasswordManager)(nil)
