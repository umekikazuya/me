package identity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Identity は認証・認可の集約
type Identity struct {
	id           identityID
	email        email
	passwordHash passwordHash
	createdAt    time.Time
	updatedAt    time.Time
}

// OptFuncIdentity はFunctionalOptionパターンを表現
type OptFuncIdentity func(*Identity) error

// NewIdentity はIdentity集約のファクトリー関数
func NewIdentity(
	inputEmail string, inputPassword string,
	hashFn func(plainPassword string) ([]byte, error),
) (*Identity, error) {
	id := newIdentityID(uuid.New())
	e, err := NewEmail(inputEmail)
	if err != nil {
		return nil, err
	}
	p, err := newPassword(inputPassword)
	if err != nil {
		return nil, err
	}
	hashedPassword, err := p.Hashed(hashFn)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	return &Identity{
		id:           id,
		email:        e,
		passwordHash: hashedPassword,
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

// --- Reconstruct ---

// ReconstructIdentityInput はReconstructIdentityの入力型
type ReconstructIdentityInput struct {
	ID           uuid.UUID
	Email        string
	PasswordHash []byte
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// ReconstructIdentity はDBから取得した信頼済みデータでIdentityを復元する
func ReconstructIdentity(input ReconstructIdentityInput) (*Identity, error) {
	return &Identity{
		id:           identityID{value: input.ID},
		email:        email{value: input.Email},
		passwordHash: passwordHash{value: input.PasswordHash},
		createdAt:    input.CreatedAt,
		updatedAt:    input.UpdatedAt,
	}, nil
}

// --- Getter ---

// ID はIdentity集約のIDを返却
func (e *Identity) ID() string {
	return e.id.Value()
}

// Email はIdentity集約のemailを返却
func (e *Identity) Email() email {
	return e.email
}

// PasswordHash はIdentityのPasswordHashを返却
func (e *Identity) PasswordHash() []byte {
	return e.passwordHash.Value()
}

// CreatedAt はIdentityのcreatedAtを返却
func (e *Identity) CreatedAt() time.Time {
	return e.createdAt
}

// UpdatedAt はIdentityのupdatedAtを返却
func (e *Identity) UpdatedAt() time.Time {
	return e.updatedAt
}

// --- 振る舞い---

// Register は認証プロファイルの発行を行う
func Register(
	email, password string,
	hashFn func(plainPassword string) ([]byte, error),
) (*Identity, error) {
	e, err := NewIdentity(
		email,
		password,
		hashFn,
	)
	if err != nil {
		return nil, err
	}

	return e, nil
}

// Authenticate はプロファイルの照合を行う
func (e *Identity) Authenticate(
	plainPassword string,
	verifyFn func(
		hashedPassword, plainPassword string,
	) error,
) error {
	err := verifyFn(string(e.PasswordHash()), plainPassword)
	if err != nil {
		return err
	}
	return nil
}

// ResetPassword はパスワード変更を行う
func (e *Identity) ResetPassword(
	inputNewPassword string,
	hashFn func(plainPassword string) ([]byte, error),
	verifyFn func(
		hashedPassword, plainPassword string,
	) error,
) error {
	p, err := newPassword(inputNewPassword)
	if err != nil {
		return err
	}
	err = verifyFn(
		string(e.PasswordHash()),
		inputNewPassword,
	)
	if err == nil {
		return errors.New("パスワードが以前と同じです")
	}
	hashed, err := p.Hashed(hashFn)
	if err != nil {
		return err
	}
	e.passwordHash = hashed
	e.updatedAt = time.Now()

	return nil
}

// ChangeEmail はメールアドレス変更を行う
func (e *Identity) ChangeEmail(input string) error {
	val, err := NewEmail(input)
	if err != nil {
		return err
	}
	e.email = val
	e.updatedAt = time.Now()

	return nil
}
