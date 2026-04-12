package identity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/umekikazuya/me/pkg/domain"
	"golang.org/x/crypto/bcrypt"
)

// Identity は認証・認可の集約
type Identity struct {
	id           identityID
	email        email
	passwordHash passwordHash
	createdAt    time.Time
	updatedAt    time.Time
	events       []domain.DomainEvent
}

// Session はセッション管理集約
type Session struct {
	tokenHash  tokenHash
	identityID identityID
	status     status
	issuedAt   time.Time
	expiresAt  time.Time
	events     []domain.DomainEvent
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

// ReconstructSessionInput はReconstructSessionの入力型
type ReconstructSessionInput struct {
	IdentityID string
	TokenHash  string
	Status     string
	IssuedAt   time.Time
	ExpiresAt  time.Time
}

// ReconstructSession はDBから取得した信頼済みデータでSessionを復元する
func ReconstructSession(input ReconstructSessionInput) (*Session, error) {
	id, err := uuid.Parse(input.IdentityID)
	if err != nil {
		return nil, err
	}
	th, err := NewTokenHash(input.TokenHash)
	if err != nil {
		return nil, err
	}
	var s status
	switch input.Status {
	case statusActive.Value():
		s = statusActive
	case statusRevoked.Value():
		s = statusRevoked
	default:
		return nil, errors.New("不正なステータスです: " + input.Status)
	}
	return &Session{
		identityID: newIdentityID(id),
		tokenHash:  th,
		status:     s,
		issuedAt:   input.IssuedAt,
		expiresAt:  input.ExpiresAt,
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

	// イベントの発行
	e.events = append(e.events, RegisteredEvent{
		identityID: e.id.Value(),
		email:      e.Email().Value(),
		occurredAt: time.Now(),
	})
	return e, nil
}

// Authenticate はプロファイルの照合を行う
func (e *Identity) Authenticate(plainPassword string) error {
	err := e.comparePassword(plainPassword)
	if err != nil {
		return err
	}
	e.events = append(e.events, AuthenticatedEvent{
		identityID: e.ID(),
		email:      e.Email().Value(),
		occurredAt: time.Now(),
	})
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

	e.events = append(e.events, PasswordResetEvent{
		identityID: e.id.Value(),
		email:      e.Email().Value(),
		occurredAt: time.Now(),
	})
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

	e.events = append(e.events, EmailChangedEvent{
		identityID: e.id.Value(),
		email:      e.Email().Value(),
		occurredAt: time.Now(),
	})
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

// CreateSession はセッションを生成
func (e *Identity) CreateSession(tokenHash string) (*Session, error) {
	return NewSession(tokenHash, e.id)
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
	e.events = append(e.events, SessionRotateEvent{
		identityID: e.identityID.Value(),
		occurredAt: time.Now(),
	})
	return new, nil
}

// Revoke はセッションの無効化を行う
func (e *Session) Revoke() error {
	err := e.revoke()
	if err != nil {
		return err
	}
	e.events = append(e.events, SessionRevokedEvent{
		identityID: e.identityID.Value(),
		occurredAt: time.Now(),
	})
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

func (e *Identity) Events() []domain.DomainEvent {
	return e.events
}

func (e *Session) Events() []domain.DomainEvent {
	return e.events
}

func (e *Identity) ClearEvents() {
	e.events = nil
}

func (e *Session) ClearEvents() {
	e.events = nil
}
