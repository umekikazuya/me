package identity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Identity は認証・認可の集約
type Identity struct {
	id           identityID
	email        email
	passwordHash passwordHash
	createdAt    time.Time
	updatedAt    time.Time
	events       []DomainEvent
}

// Session はセッション管理集約
type Session struct {
	tokenHash  tokenHash
	identityID identityID
	status     status
	issuedAt   time.Time
	expiresAt  time.Time
	events     []DomainEvent
}

// OptFuncIdentity はFunctionalOptionパターンを表現
type OptFuncIdentity func(*Identity) error

// OptFuncSession はFunctionalOptionパターンを表現
type OptFuncSession func(*Session) error

// NewIdentity はIdentity集約のファクトリー関数
func NewIdentity(
	inputEmail string, inputPassword string,
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
	hashedPassword, err := p.HashPassword()
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

const SessionExpiresInDays = 30

// NewSession はSession集約のファクトリー関数
func NewSession(
	inputTokenHash string,
	inputIdentityID identityID,
) (*Session, error) {
	h, err := NewTokenHash(inputTokenHash)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	return &Session{
		tokenHash:  h,
		identityID: inputIdentityID,
		status:     statusActive,
		issuedAt:   now,
		expiresAt:  now.Add(SessionExpiresInDays * 24 * time.Hour),
	}, nil
}

// --- Getter ---

// ID はIdentity集約のIDを返却
func (e *Identity) ID() string {
	return e.id.Value()
}

// Email はIdentity集約のemailを返却
func (e *Identity) Email() string {
	return e.email.Value()
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

// TokenHash はSession集約のtokenHashを返却
func (e *Session) TokenHash() string {
	return e.tokenHash.Value()
}

// IdentityID はSession集約のIdentityIDを返却
func (e *Session) IdentityID() string {
	return e.identityID.Value()
}

// Status はSession集約のstatusを返却
func (e *Session) Status() string {
	return e.status.Value()
}

// IssuedAt はSession集約のissuedAtを返却
func (e *Session) IssuedAt() time.Time {
	return e.issuedAt
}

// ExpiresAt はSession集約のexpiresAtを返却
func (e *Session) ExpiresAt() time.Time {
	return e.expiresAt
}

// --- 振る舞い---

// Register は認証プロファイルの発行を行う
func Register(email, password string) (*Identity, error) {
	e, err := NewIdentity(email, password)
	if err != nil {
		return nil, err
	}
	e.events = append(e.events, EventTypeRegistered)
	return e, nil
}

// Authenticate はプロファイルの照合を行う
func (e *Identity) Authenticate(plainPassword string) error {
	err := e.comparePassword(plainPassword)
	if err != nil {
		return err
	}
	e.events = append(e.events, EventTypeAuthenticated)
	return nil
}

// ResetPassword はパスワード変更を行う
func (e *Identity) ResetPassword(inputNewPassword string) error {
	p, err := newPassword(inputNewPassword)
	if err != nil {
		return err
	}
	err = e.comparePassword(inputNewPassword)
	if err == nil {
		return errors.New("パスワードが以前と同じです")
	}
	hashed, err := p.HashPassword()
	if err != nil {
		return err
	}
	e.passwordHash = hashed
	e.updatedAt = time.Now()

	e.events = append(e.events, EventTypePasswordReset)
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

	e.events = append(e.events, EventTypeEmailChanged)
	return nil
}

// comparePassword はパスワードの照合を行う
func (e *Identity) comparePassword(plainPassword string) error {
	err := bcrypt.CompareHashAndPassword(
		e.PasswordHash(),
		[]byte(plainPassword),
	)
	if err != nil {
		return errors.New("パスワードが一致していません")
	}
	return nil
}

// Rotate はセッションの無効化を行い新しいセッションを生成する
func (e *Session) Rotate(newHash string) (*Session, error) {
	new, err := NewSession(
		newHash, e.identityID,
	)
	if err != nil {
		return nil, err
	}
	err = e.revoke()
	if err != nil {
		return nil, err
	}
	e.events = append(e.events, EventTypeSessionRotated)
	return new, nil
}

// Revoke はセッションの無効化を行う
func (e *Session) Revoke() error {
	err := e.revoke()
	if err != nil {
		return err
	}
	e.events = append(e.events, EventTypeSessionRevoked)
	return nil
}

// revoke はセッションの無効化のsetter
func (e *Session) revoke() error {
	if e.Status() == statusRevoked.Value() {
		return errors.New("既にトークンが無効化されています")
	}
	e.status = statusRevoked
	return nil
}

// --- イベント ---

func (e *Identity) Events() []DomainEvent {
	return e.events
}

func (e *Session) Events() []DomainEvent {
	return e.events
}

func (e *Identity) ClearEvents() {
	e.events = nil
}

func (e *Session) ClearEvents() {
	e.events = nil
}
