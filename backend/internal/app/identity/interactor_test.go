package identity

import (
	"context"
	"errors"
	"testing"

	domain "github.com/umekikazuya/me/internal/domain/identity"
	pkgdomain "github.com/umekikazuya/me/pkg/domain"
	"github.com/umekikazuya/me/pkg/errs"
)

// --- constants ---

const (
	validEmail     = "user@example.com"
	validPassword  = "Password1"
	validTokenHash = "a3f9b2c1d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1"
)

// --- mocks ---

type mockIdentityRepo struct {
	findByIDFn    func(ctx context.Context, id string) (*domain.Identity, error)
	findByEmailFn func(ctx context.Context, email string) (*domain.Identity, error)
	saveFn        func(ctx context.Context, identity *domain.Identity) error
}

func (m *mockIdentityRepo) FindByID(ctx context.Context, id string) (*domain.Identity, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *mockIdentityRepo) FindByEmail(ctx context.Context, email string) (*domain.Identity, error) {
	if m.findByEmailFn != nil {
		return m.findByEmailFn(ctx, email)
	}
	return nil, nil
}

func (m *mockIdentityRepo) Save(ctx context.Context, identity *domain.Identity) error {
	if m.saveFn != nil {
		return m.saveFn(ctx, identity)
	}
	return nil
}

type mockSessionRepo struct {
	findByIdentityIdAndTokenHashFn func(ctx context.Context, identityID string, tokenHash string) (*domain.Session, error)
	findActiveByIdentityFn         func(ctx context.Context, identityID string) ([]*domain.Session, error)
	saveFn                         func(ctx context.Context, session *domain.Session) error
	revokeAllFn                    func(ctx context.Context, id string) error
}

func (m *mockSessionRepo) FindByIdentityIdAndTokenHash(ctx context.Context, identityID, tokenHash string) (*domain.Session, error) {
	if m.findByIdentityIdAndTokenHashFn != nil {
		return m.findByIdentityIdAndTokenHashFn(ctx, identityID, tokenHash)
	}
	return nil, nil
}

func (m *mockSessionRepo) FindActiveByIdentity(ctx context.Context, identityID string) ([]*domain.Session, error) {
	if m.findActiveByIdentityFn != nil {
		return m.findActiveByIdentityFn(ctx, identityID)
	}
	return nil, nil
}

func (m *mockSessionRepo) Save(ctx context.Context, session *domain.Session) error {
	if m.saveFn != nil {
		return m.saveFn(ctx, session)
	}
	return nil
}

func (m *mockSessionRepo) RevokeAll(ctx context.Context, id string) error {
	if m.revokeAllFn != nil {
		return m.revokeAllFn(ctx, id)
	}
	return nil
}

type mockTokenSrv struct {
	generateATFn func(ctx context.Context, identity domain.Identity) (string, error)
	generateRTFn func(ctx context.Context) (string, error)
	hashFn       func(ctx context.Context, token string) (string, error)
	validateATFn func(ctx context.Context, token string) (string, error)
}

func (m *mockTokenSrv) GenerateAT(ctx context.Context, identity domain.Identity) (string, error) {
	if m.generateATFn != nil {
		return m.generateATFn(ctx, identity)
	}
	return "access-token", nil
}

func (m *mockTokenSrv) GenerateRT(ctx context.Context) (string, error) {
	if m.generateRTFn != nil {
		return m.generateRTFn(ctx)
	}
	return "refresh-token", nil
}

func (m *mockTokenSrv) Hash(ctx context.Context, token string) (string, error) {
	if m.hashFn != nil {
		return m.hashFn(ctx, token)
	}
	return validTokenHash, nil
}

func (m *mockTokenSrv) ValidateAT(ctx context.Context, token string) (string, error) {
	if m.validateATFn != nil {
		return m.validateATFn(ctx, token)
	}
	return "", nil
}

type mockEventPublisher struct {
	publishFn func(ctx context.Context, events []pkgdomain.DomainEvent) error
}

func (m *mockEventPublisher) Publish(ctx context.Context, events []pkgdomain.DomainEvent) error {
	if m.publishFn != nil {
		return m.publishFn(ctx, events)
	}
	return nil
}

// --- helpers ---

func newInteractor(ir *mockIdentityRepo, sr *mockSessionRepo, ts *mockTokenSrv) *interactor {
	return &interactor{
		identityRepo: ir,
		sessionRepo:  sr,
		tokenSrv:     ts,
		publisher:    &mockEventPublisher{},
	}
}

// freshIdentityFn is a mock fn that returns a new *domain.Identity on every call.
// Use inside mock closures to avoid sharing mutable state across subtests.
func freshIdentityFn(_ context.Context, _ string) (*domain.Identity, error) {
	return domain.NewIdentity(validEmail, validPassword)
}

// freshSessionFn is a mock fn that returns a new active *domain.Session on every call.
func freshSessionFn(idn *domain.Identity) func(context.Context, string, string) (*domain.Session, error) {
	return func(_ context.Context, _, _ string) (*domain.Session, error) {
		return idn.CreateSession(validTokenHash)
	}
}

// assertErr checks the error result matches the test expectation.
func assertErr(t *testing.T, err error, wantErr bool, errTarget error) {
	t.Helper()
	if (err != nil) != wantErr {
		t.Fatalf("error = %v, wantErr = %v", err, wantErr)
	}
	if errTarget != nil && !errors.Is(err, errTarget) {
		t.Errorf("error = %v, want errors.Is(%v)", err, errTarget)
	}
}

// --- tests ---

func TestInteractor_Register(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tests := []struct {
		name          string
		input         InputRegisterDto
		findByEmailFn func(context.Context, string) (*domain.Identity, error)
		saveFn        func(context.Context, *domain.Identity) error
		wantErr       bool
		errTarget     error
	}{
		{
			name:  "success: 新規メールで登録",
			input: InputRegisterDto{EmailAddress: validEmail, Password: validPassword},
		},
		{
			name:  "error: メール重複",
			input: InputRegisterDto{EmailAddress: validEmail, Password: validPassword},
			findByEmailFn: func(_ context.Context, _ string) (*domain.Identity, error) {
				idn, err := domain.NewIdentity(validEmail, validPassword)
				if err != nil {
					return nil, err
				}
				return idn, nil
			},
			wantErr:   true,
			errTarget: errs.ErrConflict,
		},
		{
			name:    "error: 無効なメール形式",
			input:   InputRegisterDto{EmailAddress: "not-an-email", Password: validPassword},
			wantErr: true,
		},
		{
			name:    "error: 無効なパスワード（短すぎる）",
			input:   InputRegisterDto{EmailAddress: validEmail, Password: "weak"},
			wantErr: true,
		},
		{
			name:  "error: FindByEmail インフラ障害",
			input: InputRegisterDto{EmailAddress: validEmail, Password: validPassword},
			findByEmailFn: func(_ context.Context, _ string) (*domain.Identity, error) {
				return nil, errs.ErrInternal
			},
			wantErr:   true,
			errTarget: errs.ErrInternal,
		},
		{
			name:  "error: Save インフラ障害",
			input: InputRegisterDto{EmailAddress: validEmail, Password: validPassword},
			saveFn: func(_ context.Context, _ *domain.Identity) error {
				return errs.ErrInternal
			},
			wantErr:   true,
			errTarget: errs.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := newInteractor(
				&mockIdentityRepo{findByEmailFn: tt.findByEmailFn, saveFn: tt.saveFn},
				&mockSessionRepo{},
				&mockTokenSrv{},
			)
			err := i.Register(ctx, tt.input)
			assertErr(t, err, tt.wantErr, tt.errTarget)
		})
	}
}

func TestInteractor_Login(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tests := []struct {
		name          string
		input         InputLoginDto
		findByEmailFn func(context.Context, string) (*domain.Identity, error)
		generateATFn  func(context.Context, domain.Identity) (string, error)
		sessionSaveFn func(context.Context, *domain.Session) error
		wantErr       bool
		errTarget     error
		check         func(*testing.T, *OutputLoginDto)
	}{
		{
			name:  "success: 正常ログイン",
			input: InputLoginDto{EmailAddress: validEmail, Password: validPassword},
			findByEmailFn: func(_ context.Context, _ string) (*domain.Identity, error) {
				return domain.NewIdentity(validEmail, validPassword)
			},
			check: func(t *testing.T, got *OutputLoginDto) {
				if got == nil {
					t.Fatal("expected non-nil output")
				}
				if got.AT == "" {
					t.Error("AT must not be empty")
				}
				if got.RT == "" {
					t.Error("RT must not be empty")
				}
			},
		},
		{
			name:  "error: メールに対応するIdentityなし",
			input: InputLoginDto{EmailAddress: validEmail, Password: validPassword},
			findByEmailFn: func(_ context.Context, _ string) (*domain.Identity, error) {
				return nil, nil
			},
			wantErr:   true,
			errTarget: errs.ErrNotFound,
		},
		{
			name:  "error: パスワード不一致",
			input: InputLoginDto{EmailAddress: validEmail, Password: "WrongPass1"},
			findByEmailFn: func(_ context.Context, _ string) (*domain.Identity, error) {
				return domain.NewIdentity(validEmail, validPassword)
			},
			wantErr: true,
		},
		{
			name:    "error: 無効なメール形式",
			input:   InputLoginDto{EmailAddress: "bad-email", Password: validPassword},
			wantErr: true,
		},
		{
			name:  "error: FindByEmail インフラ障害",
			input: InputLoginDto{EmailAddress: validEmail, Password: validPassword},
			findByEmailFn: func(_ context.Context, _ string) (*domain.Identity, error) {
				return nil, errs.ErrInternal
			},
			wantErr:   true,
			errTarget: errs.ErrInternal,
		},
		{
			name:  "error: GenerateAT 失敗",
			input: InputLoginDto{EmailAddress: validEmail, Password: validPassword},
			findByEmailFn: func(_ context.Context, _ string) (*domain.Identity, error) {
				return domain.NewIdentity(validEmail, validPassword)
			},
			generateATFn: func(_ context.Context, _ domain.Identity) (string, error) {
				return "", errs.ErrInternal
			},
			wantErr: true,
		},
		{
			name:  "error: sessionRepo.Save 失敗",
			input: InputLoginDto{EmailAddress: validEmail, Password: validPassword},
			findByEmailFn: func(_ context.Context, _ string) (*domain.Identity, error) {
				return domain.NewIdentity(validEmail, validPassword)
			},
			sessionSaveFn: func(_ context.Context, _ *domain.Session) error {
				return errs.ErrInternal
			},
			wantErr:   true,
			errTarget: errs.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := newInteractor(
				&mockIdentityRepo{findByEmailFn: tt.findByEmailFn},
				&mockSessionRepo{saveFn: tt.sessionSaveFn},
				&mockTokenSrv{generateATFn: tt.generateATFn},
			)
			got, err := i.Login(ctx, tt.input)
			assertErr(t, err, tt.wantErr, tt.errTarget)
			if !tt.wantErr && tt.check != nil {
				tt.check(t, got)
			}
		})
	}
}

func TestInteractor_Logout(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tests := []struct {
		name                           string
		input                          InputLogoutDto
		findByIDFn                     func(context.Context, string) (*domain.Identity, error)
		findByIdentityIdAndTokenHashFn func(context.Context, string, string) (*domain.Session, error)
		sessionSaveFn                  func(context.Context, *domain.Session) error
		wantErr                        bool
		errTarget                      error
	}{
		{
			name:       "success: アクティブセッションを失効",
			input:      InputLogoutDto{IdentityID: "id", RT: "raw-rt"},
			findByIDFn: freshIdentityFn,
			findByIdentityIdAndTokenHashFn: func(_ context.Context, _, _ string) (*domain.Session, error) {
				idn, err := domain.NewIdentity(validEmail, validPassword)
				if err != nil {
					return nil, err
				}
				return idn.CreateSession(validTokenHash)
			},
		},
		{
			name:  "error: Identityが存在しない",
			input: InputLogoutDto{IdentityID: "id", RT: "raw-rt"},
			findByIDFn: func(_ context.Context, _ string) (*domain.Identity, error) {
				return nil, nil
			},
			wantErr:   true,
			errTarget: errs.ErrNotFound,
		},
		{
			name:       "error: セッションが存在しない",
			input:      InputLogoutDto{IdentityID: "id", RT: "raw-rt"},
			findByIDFn: freshIdentityFn,
			findByIdentityIdAndTokenHashFn: func(_ context.Context, _, _ string) (*domain.Session, error) {
				return nil, nil
			},
			wantErr:   true,
			errTarget: errs.ErrNotFound,
		},
		{
			name:       "error: セッションが既に失効済み",
			input:      InputLogoutDto{IdentityID: "id", RT: "raw-rt"},
			findByIDFn: freshIdentityFn,
			findByIdentityIdAndTokenHashFn: func(_ context.Context, _, _ string) (*domain.Session, error) {
				idn, err := domain.NewIdentity(validEmail, validPassword)
				if err != nil {
					return nil, err
				}
				ses, err := idn.CreateSession(validTokenHash)
				if err != nil {
					return nil, err
				}
				_ = ses.Revoke()
				return ses, nil
			},
			wantErr: true,
		},
		{
			name:  "error: FindByID インフラ障害",
			input: InputLogoutDto{IdentityID: "id", RT: "raw-rt"},
			findByIDFn: func(_ context.Context, _ string) (*domain.Identity, error) {
				return nil, errs.ErrInternal
			},
			wantErr:   true,
			errTarget: errs.ErrInternal,
		},
		{
			name:       "error: sessionRepo.Save 失敗",
			input:      InputLogoutDto{IdentityID: "id", RT: "raw-rt"},
			findByIDFn: freshIdentityFn,
			findByIdentityIdAndTokenHashFn: func(_ context.Context, _, _ string) (*domain.Session, error) {
				idn, err := domain.NewIdentity(validEmail, validPassword)
				if err != nil {
					return nil, err
				}
				return idn.CreateSession(validTokenHash)
			},
			sessionSaveFn: func(_ context.Context, _ *domain.Session) error {
				return errs.ErrInternal
			},
			wantErr:   true,
			errTarget: errs.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := newInteractor(
				&mockIdentityRepo{findByIDFn: tt.findByIDFn},
				&mockSessionRepo{
					findByIdentityIdAndTokenHashFn: tt.findByIdentityIdAndTokenHashFn,
					saveFn:                         tt.sessionSaveFn,
				},
				&mockTokenSrv{},
			)
			err := i.Logout(ctx, tt.input)
			assertErr(t, err, tt.wantErr, tt.errTarget)
		})
	}
}

func TestInteractor_ChangeEmail(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tests := []struct {
		name          string
		input         InputChangeEmailDto
		findByIDFn    func(context.Context, string) (*domain.Identity, error)
		findByEmailFn func(context.Context, string) (*domain.Identity, error)
		saveFn        func(context.Context, *domain.Identity) error
		wantErr       bool
		errTarget     error
	}{
		{
			name:       "success: メール変更",
			input:      InputChangeEmailDto{ID: "id", NewEmailAddress: "new@example.com"},
			findByIDFn: freshIdentityFn,
		},
		{
			name:  "error: Identityが存在しない",
			input: InputChangeEmailDto{ID: "id", NewEmailAddress: "new@example.com"},
			findByIDFn: func(_ context.Context, _ string) (*domain.Identity, error) {
				return nil, nil
			},
			wantErr:   true,
			errTarget: errs.ErrNotFound,
		},
		{
			name:       "error: 新メールが既に使用済み",
			input:      InputChangeEmailDto{ID: "id", NewEmailAddress: "taken@example.com"},
			findByIDFn: freshIdentityFn,
			findByEmailFn: func(_ context.Context, _ string) (*domain.Identity, error) {
				idn, err := domain.NewIdentity("taken@example.com", validPassword)
				if err != nil {
					return nil, err
				}
				return idn, nil
			},
			wantErr:   true,
			errTarget: errs.ErrConflict,
		},
		{
			name:       "error: 無効なメール形式",
			input:      InputChangeEmailDto{ID: "id", NewEmailAddress: "bad-email"},
			findByIDFn: freshIdentityFn,
			wantErr:    true,
		},
		{
			name:  "error: FindByID インフラ障害",
			input: InputChangeEmailDto{ID: "id", NewEmailAddress: "new@example.com"},
			findByIDFn: func(_ context.Context, _ string) (*domain.Identity, error) {
				return nil, errs.ErrInternal
			},
			wantErr:   true,
			errTarget: errs.ErrInternal,
		},
		{
			name:       "error: Save インフラ障害",
			input:      InputChangeEmailDto{ID: "id", NewEmailAddress: "new@example.com"},
			findByIDFn: freshIdentityFn,
			saveFn: func(_ context.Context, _ *domain.Identity) error {
				return errs.ErrInternal
			},
			wantErr:   true,
			errTarget: errs.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := newInteractor(
				&mockIdentityRepo{
					findByIDFn:    tt.findByIDFn,
					findByEmailFn: tt.findByEmailFn,
					saveFn:        tt.saveFn,
				},
				&mockSessionRepo{},
				&mockTokenSrv{},
			)
			err := i.ChangeEmail(ctx, tt.input)
			assertErr(t, err, tt.wantErr, tt.errTarget)
		})
	}
}

func TestInteractor_ResetPassword(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tests := []struct {
		name        string
		input       InputResetPasswordDto
		findByIDFn  func(context.Context, string) (*domain.Identity, error)
		saveFn      func(context.Context, *domain.Identity) error
		revokeAllFn func(context.Context, string) error
		wantErr     bool
		errTarget   error
	}{
		{
			name:       "success: パスワードリセット＋全セッション失効",
			input:      InputResetPasswordDto{ID: "id", NewPassword: "NewPass1"},
			findByIDFn: freshIdentityFn,
		},
		{
			name:  "error: Identityが存在しない",
			input: InputResetPasswordDto{ID: "id", NewPassword: "NewPass1"},
			findByIDFn: func(_ context.Context, _ string) (*domain.Identity, error) {
				return nil, nil
			},
			wantErr:   true,
			errTarget: errs.ErrNotFound,
		},
		{
			name:       "error: 同じパスワードを指定",
			input:      InputResetPasswordDto{ID: "id", NewPassword: validPassword},
			findByIDFn: freshIdentityFn,
			wantErr:    true,
		},
		{
			name:       "error: Save インフラ障害",
			input:      InputResetPasswordDto{ID: "id", NewPassword: "NewPass1"},
			findByIDFn: freshIdentityFn,
			saveFn: func(_ context.Context, _ *domain.Identity) error {
				return errs.ErrInternal
			},
			wantErr:   true,
			errTarget: errs.ErrInternal,
		},
		{
			name:       "error: RevokeAll インフラ障害（Save成功後）",
			input:      InputResetPasswordDto{ID: "id", NewPassword: "NewPass1"},
			findByIDFn: freshIdentityFn,
			revokeAllFn: func(_ context.Context, _ string) error {
				return errs.ErrInternal
			},
			wantErr:   true,
			errTarget: errs.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := newInteractor(
				&mockIdentityRepo{findByIDFn: tt.findByIDFn, saveFn: tt.saveFn},
				&mockSessionRepo{revokeAllFn: tt.revokeAllFn},
				&mockTokenSrv{},
			)
			err := i.ResetPassword(ctx, tt.input)
			assertErr(t, err, tt.wantErr, tt.errTarget)
		})
	}
}

func TestInteractor_RefreshTokens(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tests := []struct {
		name                           string
		input                          InputRefreshTokensDto
		findByIDFn                     func(context.Context, string) (*domain.Identity, error)
		findByIdentityIdAndTokenHashFn func(context.Context, string, string) (*domain.Session, error)
		sessionSaveFn                  func(context.Context, *domain.Session) error
		wantErr                        bool
		errTarget                      error
		check                          func(*testing.T, *OutputRefreshTokensDto)
	}{
		{
			name:       "success: トークンローテーション",
			input:      InputRefreshTokensDto{IdentityID: "id", RT: "raw-rt"},
			findByIDFn: freshIdentityFn,
			findByIdentityIdAndTokenHashFn: func(_ context.Context, _, _ string) (*domain.Session, error) {
				idn, err := domain.NewIdentity(validEmail, validPassword)
				if err != nil {
					return nil, err
				}
				return idn.CreateSession(validTokenHash)
			},
			check: func(t *testing.T, got *OutputRefreshTokensDto) {
				if got == nil {
					t.Fatal("expected non-nil output")
				}
				if got.AT == "" {
					t.Error("AT must not be empty")
				}
				if got.RT == "" {
					t.Error("RT must not be empty")
				}
			},
		},
		{
			name:  "error: Identityが存在しない",
			input: InputRefreshTokensDto{IdentityID: "id", RT: "raw-rt"},
			findByIDFn: func(_ context.Context, _ string) (*domain.Identity, error) {
				return nil, nil
			},
			wantErr:   true,
			errTarget: errs.ErrNotFound,
		},
		{
			name:       "error: セッションが存在しない",
			input:      InputRefreshTokensDto{IdentityID: "id", RT: "raw-rt"},
			findByIDFn: freshIdentityFn,
			findByIdentityIdAndTokenHashFn: func(_ context.Context, _, _ string) (*domain.Session, error) {
				return nil, nil
			},
			wantErr:   true,
			errTarget: errs.ErrNotFound,
		},
		{
			name:       "error: セッションが失効済み",
			input:      InputRefreshTokensDto{IdentityID: "id", RT: "raw-rt"},
			findByIDFn: freshIdentityFn,
			findByIdentityIdAndTokenHashFn: func(_ context.Context, _, _ string) (*domain.Session, error) {
				idn, err := domain.NewIdentity(validEmail, validPassword)
				if err != nil {
					return nil, err
				}
				ses, err := idn.CreateSession(validTokenHash)
				if err != nil {
					return nil, err
				}
				_ = ses.Revoke()
				return ses, nil
			},
			wantErr:   true,
			errTarget: errs.ErrUnprocessable,
		},
		{
			name:  "error: FindByID インフラ障害",
			input: InputRefreshTokensDto{IdentityID: "id", RT: "raw-rt"},
			findByIDFn: func(_ context.Context, _ string) (*domain.Identity, error) {
				return nil, errs.ErrInternal
			},
			wantErr:   true,
			errTarget: errs.ErrInternal,
		},
		{
			name:       "error: sessionRepo.Save 失敗",
			input:      InputRefreshTokensDto{IdentityID: "id", RT: "raw-rt"},
			findByIDFn: freshIdentityFn,
			findByIdentityIdAndTokenHashFn: func(_ context.Context, _, _ string) (*domain.Session, error) {
				idn, err := domain.NewIdentity(validEmail, validPassword)
				if err != nil {
					return nil, err
				}
				return idn.CreateSession(validTokenHash)
			},
			sessionSaveFn: func(_ context.Context, _ *domain.Session) error {
				return errs.ErrInternal
			},
			wantErr:   true,
			errTarget: errs.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := newInteractor(
				&mockIdentityRepo{findByIDFn: tt.findByIDFn},
				&mockSessionRepo{
					findByIdentityIdAndTokenHashFn: tt.findByIdentityIdAndTokenHashFn,
					saveFn:                         tt.sessionSaveFn,
				},
				&mockTokenSrv{},
			)
			got, err := i.RefreshTokens(ctx, tt.input)
			assertErr(t, err, tt.wantErr, tt.errTarget)
			if !tt.wantErr && tt.check != nil {
				tt.check(t, got)
			}
		})
	}
}

func TestInteractor_RevokeAllSessions(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tests := []struct {
		name        string
		input       InputRevokeAllSessionsDto
		findByIDFn  func(context.Context, string) (*domain.Identity, error)
		revokeAllFn func(context.Context, string) error
		wantErr     bool
		errTarget   error
	}{
		{
			name:       "success: 全セッション失効",
			input:      InputRevokeAllSessionsDto{IdentityID: "id"},
			findByIDFn: freshIdentityFn,
		},
		{
			name:  "error: Identityが存在しない",
			input: InputRevokeAllSessionsDto{IdentityID: "id"},
			findByIDFn: func(_ context.Context, _ string) (*domain.Identity, error) {
				return nil, nil
			},
			wantErr:   true,
			errTarget: errs.ErrNotFound,
		},
		{
			name:  "error: FindByID インフラ障害",
			input: InputRevokeAllSessionsDto{IdentityID: "id"},
			findByIDFn: func(_ context.Context, _ string) (*domain.Identity, error) {
				return nil, errs.ErrInternal
			},
			wantErr:   true,
			errTarget: errs.ErrInternal,
		},
		{
			name:       "error: RevokeAll インフラ障害",
			input:      InputRevokeAllSessionsDto{IdentityID: "id"},
			findByIDFn: freshIdentityFn,
			revokeAllFn: func(_ context.Context, _ string) error {
				return errs.ErrInternal
			},
			wantErr:   true,
			errTarget: errs.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := newInteractor(
				&mockIdentityRepo{findByIDFn: tt.findByIDFn},
				&mockSessionRepo{revokeAllFn: tt.revokeAllFn},
				&mockTokenSrv{},
			)
			err := i.RevokeAllSessions(ctx, tt.input)
			assertErr(t, err, tt.wantErr, tt.errTarget)
		})
	}
}
