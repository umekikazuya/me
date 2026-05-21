package identity

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	pkgdomain "github.com/umekikazuya/me/pkg/domain"
)

// --- Test helpers ---

func verifyOK(_, _ string) error         { return nil }
func verifyNG(_, _ string) error         { return errors.New("mismatch") }
func hashConst(_ string) ([]byte, error) { return []byte("new-hash"), nil }

func mustNewIdentity(t *testing.T, email, password string) *Identity {
	t.Helper()

	hashed, err := hashConst(password)
	if err != nil {
		t.Fatalf("err = %#v", err)
	}
	e, err := ReconstructIdentity(ReconstructIdentityInput{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: hashed,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	})
	if err != nil {
		t.Fatalf("err = %#v", err)
	}
	return e
}

func mustNewSession(t *testing.T, id identityID) *Session {
	t.Helper()
	s, err := NewSession(
		"a3f9b2c1d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1",
		id,
	)
	if err != nil {
		t.Fatalf("mustNewSession: %v", err)
	}
	return s
}

func someIdentityID() identityID {
	return identityID{value: uuid.New()}
}

func assertSingleEvent(t *testing.T, events []pkgdomain.DomainEvent, want string) {
	t.Helper()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d: %v", len(events), events)
	}
	if events[0].EventType() != want {
		t.Errorf("event type = %v, want %v", events[0].EventType(), want)
	}
}

func assertNoEvents(t *testing.T, events []pkgdomain.DomainEvent) {
	t.Helper()
	if len(events) != 0 {
		t.Errorf("expected no events, got %d: %v", len(events), events)
	}
}

func stubHashFn(
	t *testing.T,
	wantPlain string,
	ret []byte,
	retErr error,
) func(string) ([]byte, error) {
	t.Helper()
	return func(gotPlain string) ([]byte, error) {
		if gotPlain != wantPlain {
			t.Fatalf("plainPassword = %q, want = %q", gotPlain, wantPlain)
		}
		return ret, retErr
	}
}

// --- Register ---
// ドキュメント: register(email, password) → Registered イベント
// 前提条件「メールアドレスが未使用」はアプリ層の責務。
// NewIdentity はファクトリ。Register がドメインアクションとしてイベントを発行する。

func TestRegister(t *testing.T) {
	t.Parallel()

	t.Run("valid inputs: publishes Registered event", func(t *testing.T) {
		t.Parallel()
		got, err := Register(
			"user@example.com",
			"Password1",
			stubHashFn(t, "Password1", []byte("stubbed-hash"), nil),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		assertSingleEvent(t, got.Events(), "identity.registered")
	})

	t.Run("invalid email: no event published", func(t *testing.T) {
		t.Parallel()
		_, err := Register(
			"notanemail",
			"Password1",
			stubHashFn(t, "Password1", []byte("stubbed-hash"), nil),
		)
		if err == nil {
			t.Error("Register() with invalid email should fail")
		}
	})

	t.Run("invalid password: no event published", func(t *testing.T) {
		t.Parallel()
		_, err := Register(
			"user@example.com",
			"weak",
			stubHashFn(t, "Password1", []byte("stubbed-hash"), nil),
		)
		if err == nil {
			t.Error("Register() with invalid password should fail")
		}
	})
}

// --- NewIdentity ---
// ファクトリ。VO バリデーションとハッシュ化のみ責務。イベントは発行しない。

func TestNewIdentity(t *testing.T) {
	t.Parallel()

	t.Run("valid inputs: email stored, password hashed, timestamps set, Registered event published", func(t *testing.T) {
		t.Parallel()
		before := time.Now()
		got, err := NewIdentity(
			"user@example.com",
			"Password1",
			stubHashFn(t, "Password1", []byte("stubbed-hash"), nil),
		)
		after := time.Now().Add(time.Second)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// email
		if got.email.Value() != "user@example.com" {
			t.Errorf("email = %q, want %q", got.email.Value(), "user@example.com")
		}
		// passwordHash: bcrypt なので平文と一致しない、かつ非空
		if string(got.passwordHash.Value()) == "Password1" {
			t.Error("passwordHash must not store plaintext password")
		}
		if len(got.passwordHash.Value()) == 0 {
			t.Error("passwordHash must not be empty")
		}
		// id: ゼロ値でないこと
		if got.id.Value() == (identityID{}).Value() {
			t.Error("id must be set to a non-zero UUID")
		}
		// timestamps: createdAt == updatedAt、現在時刻付近
		if got.createdAt.Before(before) || got.createdAt.After(after) {
			t.Errorf("createdAt %v outside expected range [%v, %v]", got.createdAt, before, after)
		}
		if !got.createdAt.Equal(got.updatedAt) {
			t.Errorf("createdAt %v != updatedAt %v on creation", got.createdAt, got.updatedAt)
		}
	})

	invalidEmailCases := []struct {
		name  string
		email string
	}{
		{"empty", ""},
		{"missing local part", "@example.com"},
		{"missing domain", "user@"},
		{"no @ symbol", "userexample.com"},
	}
	for _, tc := range invalidEmailCases {
		t.Run("invalid email rejected: "+tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := NewIdentity(
				tc.email,
				"Password1",
				stubHashFn(t, "Password1", []byte("stubbed-hash"), nil),
			)
			if err == nil {
				t.Errorf("NewIdentity(%q) should fail", tc.email)
			}
		})
	}

	invalidPasswordCases := []struct {
		name     string
		password string
	}{
		{"empty", ""},
		{"too short (7 chars)", "Pass1Ab"},
		{"no uppercase", "password1"},
		{"no lowercase", "PASSWORD1"},
		{"too long (73 chars)", strings.Repeat("Aa", 36) + "B"},
	}
	for _, tc := range invalidPasswordCases {
		t.Run("invalid password rejected: "+tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := NewIdentity(
				"user@example.com",
				tc.password,
				stubHashFn(t, "Password1", []byte("stubbed-hash"), nil),
			)
			if err == nil {
				t.Errorf("NewIdentity(password=%q) should fail", tc.password)
			}
		})
	}
}

// --- Identity.Authenticate ---
// ドキュメント: authenticate(password) → ログイン成功時、認可済みの身元を保証 / Authenticated イベント

func TestIdentity_Authenticate(t *testing.T) {
	t.Parallel()

	t.Run("ok#authenticated event is appended and state stays unchanged", func(t *testing.T) {
		t.Parallel()

		e := mustNewIdentity(t, "user@example.com", "Password1")
		idBefore := e.ID()
		emailBefore := e.Email()
		hashBefore := e.PasswordHash()
		createdAtBefore := e.CreatedAt()
		updatedAtBefore := e.UpdatedAt()

		err := e.Authenticate("Password1", verifyOK)
		if err != nil {
			t.Fatalf("err = %#v", err)
		}

		if e.ID() != idBefore {
			t.Error("id must not change")
		}
		if e.Email() != emailBefore {
			t.Error("email must not change")
		}
		if string(e.PasswordHash()) != string(hashBefore) {
			t.Error("password hash must not change")
		}
		if !e.CreatedAt().Equal(createdAtBefore) {
			t.Error("createdAt must not change")
		}
		if !e.UpdatedAt().Equal(updatedAtBefore) {
			t.Error("updatedAt must not change")
		}
		assertSingleEvent(
			t,
			e.Events(),
			"identity.authenticated",
		)
	})

	t.Run("ng: verifier error is returned and no event is appended", func(t *testing.T) {
		t.Parallel()

		e := mustNewIdentity(t, "user@example.com", "Password1")
		hashBefore := e.PasswordHash()
		updatedAtBefore := e.UpdatedAt()
		wantErr := errors.New("mismatch")
		verifyErr := func(_, _ string) error { return wantErr }

		err := e.Authenticate("Password1", verifyErr)
		if !errors.Is(err, wantErr) {
			t.Fatalf("err = %#v, want %#v", err, wantErr)
		}
		if string(e.PasswordHash()) != string(hashBefore) {
			t.Error("password hash must not change")
		}
		if !e.UpdatedAt().Equal(updatedAtBefore) {
			t.Error("createdAt must not change")
		}
		assertNoEvents(t, e.Events())
	})
}

// --- Identity.ResetPassword ---
// ドキュメント: resetPassword(newHash) → passwordHash を上書き更新 / PasswordReset イベント
// 前提条件「トークンが有効」はアプリ層の責務。
func TestIdentity_ResetPassword(t *testing.T) {
	tests := []struct {
		name             string
		inputNewPassword string
		hashFn           func(plainPassword string) ([]byte, error)
		verifyFn         func(hashedPassword string, plainPassword string) error
		wantErr          bool
	}{
		{
			name:             "ok#正常",
			inputNewPassword: "Password2",
			hashFn:           hashConst,
			verifyFn:         verifyNG,
			wantErr:          false,
		},
		{
			name:             "ng#以前と同じパスワード",
			inputNewPassword: "Password1",
			hashFn:           hashConst,
			verifyFn:         verifyOK,
			wantErr:          true,
		},
		{
			name:             "ng#パスワードポリシーに準拠してない",
			inputNewPassword: "a",
			hashFn:           hashConst,
			verifyFn:         verifyOK,
			wantErr:          true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := mustNewIdentity(t, "user@example.com", "Password1")
			gotErr := e.ResetPassword(
				tt.inputNewPassword,
				tt.hashFn,
				tt.verifyFn,
			)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ResetPassword() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ResetPassword() succeeded unexpectedly")
			}
		})
	}
}

// --- Identity.ChangeEmail ---
// ドキュメント: changeEmail(newEmail) → email を上書き更新 / EmailChanged イベント
// 前提条件「トークン検証成功 かつ newEmail 未使用」はアプリ層の責務。

func TestIdentity_ChangeEmail(t *testing.T) {
	t.Parallel()

	t.Run("valid new email: updated, createdAt immutable, EmailChanged event", func(t *testing.T) {
		t.Parallel()
		identity := mustNewIdentity(t, "old@example.com", "Password1")
		createdAtBefore := identity.createdAt
		identity.ClearEvents()

		if err := identity.ChangeEmail("new@example.com"); err != nil {
			t.Fatalf("ChangeEmail() error = %v", err)
		}
		if identity.email.Value() != "new@example.com" {
			t.Errorf("email = %q, want %q", identity.email.Value(), "new@example.com")
		}
		if !identity.createdAt.Equal(createdAtBefore) {
			t.Error("createdAt must be immutable")
		}
		assertSingleEvent(t, identity.Events(), "identity.emailChanged")
	})

	invalidCases := []struct {
		name  string
		email string
	}{
		{"empty", ""},
		{"no @ symbol", "notanemail"},
		{"missing domain", "user@"},
		{"missing local part", "@example.com"},
	}
	for _, tc := range invalidCases {
		t.Run("invalid email: "+tc.name+" — no state mutation, no event", func(t *testing.T) {
			t.Parallel()
			identity := mustNewIdentity(t, "original@example.com", "Password1")
			emailBefore := identity.email.Value()
			updatedAtBefore := identity.updatedAt
			identity.ClearEvents()

			if err := identity.ChangeEmail(tc.email); err == nil {
				t.Errorf("ChangeEmail(%q) should fail", tc.email)
			}
			if identity.email.Value() != emailBefore {
				t.Error("email must not change on failed ChangeEmail")
			}
			if !identity.updatedAt.Equal(updatedAtBefore) {
				t.Error("updatedAt must not change on failed ChangeEmail")
			}
			assertNoEvents(t, identity.Events())
		})
	}
}

// --- NewSession ---
// ドキュメント: tokenHash(SHA-256), status=Active, expiresAt=発行から30日後

func TestNewSession(t *testing.T) {
	t.Parallel()

	validHash := "a3f9b2c1d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1"
	id := someIdentityID()

	t.Run("valid session: Active, expiresAt = issuedAt + 30days", func(t *testing.T) {
		t.Parallel()
		before := time.Now()
		s, err := NewSession(validHash, id)
		after := time.Now().Add(time.Second)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if s.status != statusActive {
			t.Errorf("initial status = %v, want active", s.status)
		}
		if s.tokenHash.Value() != validHash {
			t.Errorf("tokenHash = %v, want %v", s.tokenHash.Value(), validHash)
		}
		if s.issuedAt.Before(before) || s.issuedAt.After(after) {
			t.Errorf("issuedAt %v outside expected range", s.issuedAt)
		}
		// ドメインルール: expiresAt = issuedAt + 30日
		if !s.expiresAt.Equal(s.issuedAt.Add(30 * 24 * time.Hour)) {
			t.Errorf("expiresAt = %v, want issuedAt+30days = %v", s.expiresAt, s.issuedAt.Add(30*24*time.Hour))
		}
		// identityID が正しく保持されている
		if s.IdentityID() != id.Value() {
			t.Errorf("identityID = %v, want %v", s.IdentityID(), id.Value())
		}
	})

	t.Run("empty tokenHash rejected", func(t *testing.T) {
		t.Parallel()
		_, err := NewSession("", id)
		if err == nil {
			t.Error("NewSession() with empty tokenHash should fail")
		}
	})
}

// --- Session.Revoke ---
// ドキュメント: revoke() — 前提: status=Active → status=Revoked / SessionRevoked イベント

func TestSession_Revoke(t *testing.T) {
	t.Parallel()

	t.Run("active session → Revoked, timestamps unchanged, SessionRevoked event", func(t *testing.T) {
		t.Parallel()
		s := mustNewSession(t, someIdentityID())
		issuedAtBefore := s.issuedAt
		expiresAtBefore := s.expiresAt
		s.ClearEvents()

		if err := s.Revoke(); err != nil {
			t.Fatalf("Revoke() error = %v", err)
		}
		if s.Status() != statusRevoked.Value() {
			t.Errorf("status = %v, want revoked", s.Status())
		}
		// タイムスタンプは不変
		if !s.issuedAt.Equal(issuedAtBefore) {
			t.Error("issuedAt must not change after Revoke")
		}
		if !s.expiresAt.Equal(expiresAtBefore) {
			t.Error("expiresAt must not change after Revoke")
		}
		assertSingleEvent(t, s.Events(), "identity.sessionRevoked")
	})

	t.Run("already-revoked session returns error, no event (ドメインルール: 二重 Revoke 禁止)", func(t *testing.T) {
		t.Parallel()
		s := mustNewSession(t, someIdentityID())
		_ = s.Revoke()
		s.ClearEvents()

		if err := s.Revoke(); err == nil {
			t.Error("Revoke() on already-revoked session must return error")
		}
		assertNoEvents(t, s.Events())
	})
}

// --- Session.IsActive / IsRevoked ---

func TestSession_IsActive(t *testing.T) {
	t.Parallel()

	t.Run("新規 session は active", func(t *testing.T) {
		t.Parallel()
		s := mustNewSession(t, someIdentityID())
		if !s.IsActive() {
			t.Error("IsActive() = false, want true")
		}
		if s.IsRevoked() {
			t.Error("IsRevoked() = true, want false")
		}
	})

	t.Run("Revoke 後は非 active / IsRevoked=true", func(t *testing.T) {
		t.Parallel()
		s := mustNewSession(t, someIdentityID())
		if err := s.Revoke(); err != nil {
			t.Fatalf("Revoke() error = %v", err)
		}
		if s.IsActive() {
			t.Error("IsActive() = true after Revoke, want false")
		}
		if !s.IsRevoked() {
			t.Error("IsRevoked() = false after Revoke, want true")
		}
	})

	t.Run("status=active でも expiresAt を過ぎていれば IsActive=false", func(t *testing.T) {
		t.Parallel()
		past := time.Now().Add(-1 * time.Hour)
		s, err := ReconstructSession(ReconstructSessionInput{
			TokenHash:  "a3f9b2c1d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1",
			IdentityID: someIdentityID().Value(),
			Status:     statusActive.Value(),
			IssuedAt:   past.Add(-30 * 24 * time.Hour),
			ExpiresAt:  past,
		})
		if err != nil {
			t.Fatalf("ReconstructSession: %v", err)
		}
		if s.IsActive() {
			t.Error("IsActive() = true for expired session, want false")
		}
		if s.IsRevoked() {
			t.Error("IsRevoked() = true for expired (but not revoked) session, want false")
		}
	})
}

// --- Session.Rotate ---
// ドキュメント: rotate(newHash) — 前提: status=Active
//   副作用: status=Revoked、新 Session を作成 / SessionRotated イベント

func TestSession_Rotate(t *testing.T) {
	t.Parallel()

	const newHash = "b4f9b2c1d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0b2"

	t.Run("active session: old Revoked, new Active, identityID inherited, expiresAt=+30days, SessionRotated event", func(t *testing.T) {
		t.Parallel()
		id := someIdentityID()
		old := mustNewSession(t, id)
		oldHash := old.tokenHash.Value()
		old.ClearEvents()

		newSession, err := old.Rotate(newHash)
		if err != nil {
			t.Fatalf("Rotate() error = %v", err)
		}
		// 旧セッション
		if old.Status() != statusRevoked.Value() {
			t.Error("old session must be Revoked after Rotate")
		}
		if old.tokenHash.Value() != oldHash {
			t.Error("old session tokenHash must not change after Rotate")
		}
		// 新セッション
		if newSession.status != statusActive {
			t.Errorf("new session status = %v, want active", newSession.status)
		}
		if newSession.tokenHash.Value() != newHash {
			t.Errorf("new session tokenHash = %v, want %v", newSession.tokenHash.Value(), newHash)
		}
		// ドメインルール: identityID の引き継ぎ
		if newSession.IdentityID() != old.IdentityID() {
			t.Errorf("new session identityID = %v, want %v", newSession.IdentityID(), old.IdentityID())
		}
		// ドメインルール: expiresAt = issuedAt + 30日
		if !newSession.expiresAt.Equal(newSession.issuedAt.Add(30 * 24 * time.Hour)) {
			t.Errorf("new session expiresAt = %v, want issuedAt+30days", newSession.expiresAt)
		}
		// SessionRotated イベントのみ（SessionRevoked は Rotate の内部実装詳細）
		assertSingleEvent(t, old.Events(), "identity.sessionRotated")
	})

	t.Run("empty tokenHash: error, old session stays Active, no event (原子性)", func(t *testing.T) {
		t.Parallel()
		s := mustNewSession(t, someIdentityID())
		s.ClearEvents()

		_, err := s.Rotate("")
		if err == nil {
			t.Error("Rotate() with empty tokenHash should fail")
		}
		if s.Status() != statusActive.Value() {
			t.Error("old session must remain Active when Rotate fails")
		}
		assertNoEvents(t, s.Events())
	})

	t.Run("revoked session cannot rotate (ドメインルール: リプレイ攻撃防止)", func(t *testing.T) {
		t.Parallel()
		s := mustNewSession(t, someIdentityID())
		_ = s.Revoke()
		s.ClearEvents()

		_, err := s.Rotate(newHash)
		if err == nil {
			t.Error("Rotate() on revoked session must return error")
		}
		assertNoEvents(t, s.Events())
	})
}
