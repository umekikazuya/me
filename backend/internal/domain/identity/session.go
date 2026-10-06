package identity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

const SessionExpiresInDays = 30

// Session はセッション管理集約
type Session struct {
	tokenHash  tokenHash
	identityID identityID
	status     status
	issuedAt   time.Time
	expiresAt  time.Time
}

// OptFuncSession is Functional Option pattern.
type OptFuncSession func(*Session) error

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

// IsActive はセッションが「未失効 かつ 未期限切れ」なら true を返す。
// 呼び出し側は Status() == "active" を直接見ず、必ずこのメソッドで判定する
// (status が active のまま expiresAt を過ぎているケースを取りこぼさないため)。
func (e *Session) IsActive() bool {
	return e.status.Value() == statusActive.Value() && time.Now().Before(e.expiresAt)
}

// IsRevoked はセッションが revoked 状態なら true を返す。
func (e *Session) IsRevoked() bool {
	return e.status.Value() == statusRevoked.Value()
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
	return new, nil
}

// Revoke はセッションの無効化を行う
func (e *Session) Revoke() error {
	err := e.revoke()
	if err != nil {
		return err
	}
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
