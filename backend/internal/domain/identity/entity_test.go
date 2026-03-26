package identity

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

// --- Test helpers ---

func mustNewIdentity(t *testing.T, email, password string) *Identity {
	t.Helper()
	i, err := NewIdentity(email, password)
	if err != nil {
		t.Fatalf("mustNewIdentity(%q, %q): %v", email, password, err)
	}
	return i
}

func mustNewSession(t *testing.T, id identityID) *Session {
	t.Helper()
	s, err := NewSession("a3f9b2c1d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1", id)
	if err != nil {
		t.Fatalf("mustNewSession: %v", err)
	}
	return s
}

func someIdentityID() identityID {
	return identityID{value: uuid.New()}
}

func assertSingleEvent(t *testing.T, events []DomainEvent, want EventType) {
	t.Helper()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d: %v", len(events), events)
	}
	if events[0].Type() != want {
		t.Errorf("event type = %v, want %v", events[0].Type(), want)
	}
}

func assertNoEvents(t *testing.T, events []DomainEvent) {
	t.Helper()
	if len(events) != 0 {
		t.Errorf("expected no events, got %d: %v", len(events), events)
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
		got, err := Register("user@example.com", "Password1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		assertSingleEvent(t, got.Events(), EventTypeRegistered)
	})

	t.Run("invalid email: no event published", func(t *testing.T) {
		t.Parallel()
		_, err := Register("notanemail", "Password1")
		if err == nil {
			t.Error("Register() with invalid email should fail")
		}
	})

	t.Run("invalid password: no event published", func(t *testing.T) {
		t.Parallel()
		_, err := Register("user@example.com", "weak")
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
		got, err := NewIdentity("user@example.com", "Password1")
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
			_, err := NewIdentity(tc.email, "Password1")
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
			_, err := NewIdentity("user@example.com", tc.password)
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

	t.Run("correct password: succeeds and publishes Authenticated event", func(t *testing.T) {
		t.Parallel()
		identity := mustNewIdentity(t, "user@example.com", "Password1")
		identity.ClearEvents()

		if err := identity.Authenticate("Password1"); err != nil {
			t.Fatalf("Authenticate() error = %v", err)
		}
		assertSingleEvent(t, identity.Events(), EventTypeAuthenticated)
	})

	t.Run("wrong password: rejected, no event, no state mutation", func(t *testing.T) {
		t.Parallel()
		identity := mustNewIdentity(t, "user@example.com", "Password1")
		hashBefore := string(identity.passwordHash.Value())
		updatedAtBefore := identity.updatedAt
		identity.ClearEvents()

		if err := identity.Authenticate("WrongPass1"); err == nil {
			t.Error("Authenticate() should fail with wrong password")
		}
		if string(identity.passwordHash.Value()) != hashBefore {
			t.Error("passwordHash must not change on failed authentication")
		}
		if !identity.updatedAt.Equal(updatedAtBefore) {
			t.Error("updatedAt must not change on failed authentication")
		}
		assertNoEvents(t, identity.Events())
	})

	t.Run("password is case-sensitive", func(t *testing.T) {
		t.Parallel()
		identity := mustNewIdentity(t, "user@example.com", "Password1")
		if err := identity.Authenticate("password1"); err == nil {
			t.Error("Authenticate() must be case-sensitive")
		}
		if err := identity.Authenticate("PASSWORD1"); err == nil {
			t.Error("Authenticate() must be case-sensitive")
		}
	})

	t.Run("empty password: rejected without panicking", func(t *testing.T) {
		t.Parallel()
		identity := mustNewIdentity(t, "user@example.com", "Password1")
		if err := identity.Authenticate(""); err == nil {
			t.Error("Authenticate() should fail with empty password")
		}
	})

	t.Run("authenticate is idempotent: multiple correct attempts all succeed", func(t *testing.T) {
		t.Parallel()
		identity := mustNewIdentity(t, "user@example.com", "Password1")
		for i := range 3 {
			if err := identity.Authenticate("Password1"); err != nil {
				t.Errorf("attempt %d: Authenticate() error = %v", i+1, err)
			}
		}
	})
}

// --- Identity.ResetPassword ---
// ドキュメント: resetPassword(newHash) → passwordHash を上書き更新 / PasswordReset イベント
// 前提条件「トークンが有効」はアプリ層の責務。

func TestIdentity_ResetPassword(t *testing.T) {
	t.Parallel()

	t.Run("valid new password: hash updated, not plaintext, createdAt immutable, updatedAt advances, PasswordReset event", func(t *testing.T) {
		t.Parallel()
		identity := mustNewIdentity(t, "user@example.com", "Password1")
		createdAtBefore := identity.createdAt
		updatedAtBefore := identity.updatedAt
		hashBefore := string(identity.passwordHash.Value())
		identity.ClearEvents()

		if err := identity.ResetPassword("NewPass99"); err != nil {
			t.Fatalf("ResetPassword() error = %v", err)
		}
		if string(identity.passwordHash.Value()) == hashBefore {
			t.Error("passwordHash must change after ResetPassword")
		}
		if string(identity.passwordHash.Value()) == "NewPass99" {
			t.Error("passwordHash must not store plaintext")
		}
		if !identity.createdAt.Equal(createdAtBefore) {
			t.Error("createdAt must be immutable")
		}
		if !identity.updatedAt.After(updatedAtBefore) {
			t.Error("updatedAt must advance after ResetPassword")
		}
		assertSingleEvent(t, identity.Events(), EventTypePasswordReset)
	})

	t.Run("old password rejected, new password authenticates after reset", func(t *testing.T) {
		t.Parallel()
		identity := mustNewIdentity(t, "user@example.com", "Password1")

		if err := identity.ResetPassword("NewPass99"); err != nil {
			t.Fatalf("ResetPassword() error = %v", err)
		}
		if err := identity.Authenticate("Password1"); err == nil {
			t.Error("old password must not authenticate after reset")
		}
		if err := identity.Authenticate("NewPass99"); err != nil {
			t.Errorf("new password must authenticate after reset: %v", err)
		}
	})

	t.Run("reset to same password is rejected (仕様: 現在と同じパスワードへの変更禁止)", func(t *testing.T) {
		t.Parallel()
		identity := mustNewIdentity(t, "user@example.com", "Password1")
		hashBefore := string(identity.passwordHash.Value())
		identity.ClearEvents()

		if err := identity.ResetPassword("Password1"); err == nil {
			t.Error("ResetPassword to same password should fail")
		}
		if string(identity.passwordHash.Value()) != hashBefore {
			t.Error("passwordHash must not change on same-password reset")
		}
		assertNoEvents(t, identity.Events())
	})

	invalidPasswordCases := []struct {
		name     string
		password string
	}{
		{"empty", ""},
		{"too short", "Pass1Ab"},
		{"no uppercase", "password1"},
	}
	for _, tc := range invalidPasswordCases {
		t.Run("invalid password: "+tc.name+" — no state mutation, no event", func(t *testing.T) {
			t.Parallel()
			identity := mustNewIdentity(t, "user@example.com", "Password1")
			hashBefore := string(identity.passwordHash.Value())
			updatedAtBefore := identity.updatedAt
			identity.ClearEvents()

			if err := identity.ResetPassword(tc.password); err == nil {
				t.Errorf("ResetPassword(%q) should fail", tc.password)
			}
			if string(identity.passwordHash.Value()) != hashBefore {
				t.Error("passwordHash must not change on failed ResetPassword")
			}
			if !identity.updatedAt.Equal(updatedAtBefore) {
				t.Error("updatedAt must not change on failed ResetPassword")
			}
			assertNoEvents(t, identity.Events())
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
		assertSingleEvent(t, identity.Events(), EventTypeEmailChanged)
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
		assertSingleEvent(t, s.Events(), EventTypeSessionRevoked)
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
		assertSingleEvent(t, old.Events(), EventTypeSessionRotated)
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
