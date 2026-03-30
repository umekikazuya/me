package identity

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	appidentity "github.com/umekikazuya/me/internal/app/identity"
	domain "github.com/umekikazuya/me/internal/domain/identity"
)

// --- モック ---

type mockInteractor struct{}

func (m *mockInteractor) ChangeEmail(ctx context.Context, input appidentity.InputChangeEmailDto) error {
	return nil
}

func (m *mockInteractor) Login(ctx context.Context, input appidentity.InputLoginDto) (*appidentity.OutputLoginDto, error) {
	return nil, nil
}

func (m *mockInteractor) Logout(ctx context.Context, input appidentity.InputLogoutDto) error {
	return nil
}

func (m *mockInteractor) ResetPassword(
	ctx context.Context,
	input appidentity.InputResetPasswordDto,
) error {
	return nil
}

func (m *mockInteractor) RefreshTokens(ctx context.Context, input appidentity.InputRefreshTokensDto) (*appidentity.OutputRefreshTokensDto, error) {
	return nil, nil
}

func (m *mockInteractor) Register(ctx context.Context, input appidentity.InputRegisterDto) error {
	return nil
}

func (m *mockInteractor) RevokeAllSessions(ctx context.Context, input appidentity.InputRevokeAllSessionsDto) error {
	return nil
}

type mockTokenSrv struct {
	validateATFn func(ctx context.Context, token string) (string, error)
}

func (m *mockTokenSrv) GenerateAT(_ context.Context, _ domain.Identity) (string, error) {
	return "", nil
}
func (m *mockTokenSrv) GenerateRT(_ context.Context) (string, error)     { return "", nil }
func (m *mockTokenSrv) Hash(_ context.Context, _ string) (string, error) { return "", nil }
func (m *mockTokenSrv) ValidateAT(ctx context.Context, token string) (string, error) {
	if m.validateATFn != nil {
		return m.validateATFn(ctx, token)
	}
	return "", nil
}

// --- ヘルパー ---

const mwTestSecret = "test-secret-key-32-bytes-minimum!"

func newTestHandler(validateATFn func(context.Context, string) (string, error)) *Handler {
	mockInteractor := &mockInteractor{}
	mockTokenSrv := &mockTokenSrv{validateATFn: validateATFn}
	return NewHandler(mockInteractor, mockTokenSrv)
}

// 有効な AT（HS256、mwTestSecret で署名）
func makeMWValidAT(t *testing.T, sub string) string {
	t.Helper()
	tok := gojwt.NewWithClaims(gojwt.SigningMethodHS256, gojwt.RegisteredClaims{
		Subject:   sub,
		ExpiresAt: gojwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
	})
	s, err := tok.SignedString([]byte(mwTestSecret))
	if err != nil {
		t.Fatalf("makeMWValidAT: %v", err)
	}
	return s
}

// 期限切れ AT（sleep なし）
func makeMWExpiredAT(t *testing.T, sub string) string {
	t.Helper()
	tok := gojwt.NewWithClaims(gojwt.SigningMethodHS256, gojwt.RegisteredClaims{
		Subject:   sub,
		ExpiresAt: gojwt.NewNumericDate(time.Now().Add(-time.Hour)),
	})
	s, err := tok.SignedString([]byte(mwTestSecret))
	if err != nil {
		t.Fatalf("makeMWExpiredAT: %v", err)
	}
	return s
}

// sub に UUID でない値を持つ AT
func makeMWATWithSub(t *testing.T, sub string) string {
	t.Helper()
	tok := gojwt.NewWithClaims(gojwt.SigningMethodHS256, gojwt.RegisteredClaims{
		Subject:   sub,
		ExpiresAt: gojwt.NewNumericDate(time.Now().Add(time.Hour)),
	})
	s, err := tok.SignedString([]byte(mwTestSecret))
	if err != nil {
		t.Fatalf("makeMWATWithSub: %v", err)
	}
	return s
}

// alg:none 攻撃トークン（署名なし）
func makeMWNoneAlgAT(sub string) string {
	h := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none","typ":"JWT"}`))
	p := base64.RawURLEncoding.EncodeToString(
		[]byte(fmt.Sprintf(`{"sub":%q,"exp":%d}`, sub, time.Now().Add(time.Hour).Unix())),
	)
	return h + "." + p + "."
}

// next が呼ばれたか + context の identityID を記録するハンドラ
type captureHandler struct {
	called bool
	id     string
}

func (c *captureHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.called = true
	c.id, _ = identityIDFromContext(r.Context())
	w.WriteHeader(http.StatusOK)
}

func withCookie(r *http.Request, name, value string) *http.Request {
	r.AddCookie(&http.Cookie{Name: name, Value: value})
	return r
}

// --- CSRFMiddleware ---

func TestCSRFMiddleware(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		header     string // 空文字 = ヘッダーなし
		wantStatus int
		wantCalled bool
	}{
		{
			name:       "X-Requested-With: XMLHttpRequest → next を呼ぶ",
			header:     "XMLHttpRequest",
			wantStatus: http.StatusOK,
			wantCalled: true,
		},
		{
			name:       "ヘッダーなし → 403",
			header:     "",
			wantStatus: http.StatusForbidden,
			wantCalled: false,
		},
		{
			name:       "不正な値（fetch）→ 403",
			header:     "fetch",
			wantStatus: http.StatusForbidden,
			wantCalled: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			cap := &captureHandler{}
			h := CSRFMiddleware(cap)

			r := httptest.NewRequest(http.MethodPost, "/", nil)
			if tc.header != "" {
				r.Header.Set("X-Requested-With", tc.header)
			}
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)

			if w.Code != tc.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tc.wantStatus)
			}
			if cap.called != tc.wantCalled {
				t.Errorf("next called = %v, want %v", cap.called, tc.wantCalled)
			}
		})
	}
}

// --- AuthMiddleware ---

func TestHandler_AuthMiddleware(t *testing.T) {
	t.Parallel()

	validSub := uuid.New().String()

	cases := []struct {
		name         string
		buildReq     func(t *testing.T) *http.Request
		validateATFn func(context.Context, string) (string, error)
		wantStatus   int
		wantCalled   bool
		wantID       string
	}{
		{
			name: "meAccessToken Cookie なし → 401",
			buildReq: func(t *testing.T) *http.Request {
				return httptest.NewRequest(http.MethodGet, "/", nil)
			},
			wantStatus: http.StatusUnauthorized,
			wantCalled: false,
		},
		{
			name: "ValidateAT が ErrTokenInvalid → 401",
			buildReq: func(t *testing.T) *http.Request {
				return withCookie(
					httptest.NewRequest(http.MethodGet, "/", nil),
					"meAccessToken", makeMWValidAT(t, validSub),
				)
			},
			validateATFn: func(_ context.Context, _ string) (string, error) {
				return "", appidentity.ErrTokenInvalid
			},
			wantStatus: http.StatusUnauthorized,
			wantCalled: false,
		},
		{
			name: "ValidateAT が ErrTokenExpired → 401",
			buildReq: func(t *testing.T) *http.Request {
				return withCookie(
					httptest.NewRequest(http.MethodGet, "/", nil),
					"meAccessToken", makeMWExpiredAT(t, validSub),
				)
			},
			validateATFn: func(_ context.Context, _ string) (string, error) {
				return "", appidentity.ErrTokenExpired
			},
			wantStatus: http.StatusUnauthorized,
			wantCalled: false,
		},
		{
			name: "alg:none 攻撃トークン → ValidateAT が拒否して 401",
			buildReq: func(t *testing.T) *http.Request {
				return withCookie(
					httptest.NewRequest(http.MethodGet, "/", nil),
					"meAccessToken", makeMWNoneAlgAT(validSub),
				)
			},
			validateATFn: func(_ context.Context, _ string) (string, error) {
				return "", appidentity.ErrTokenInvalid
			},
			wantStatus: http.StatusUnauthorized,
			wantCalled: false,
		},
		{
			name: "有効なトークン → next が呼ばれ identityID がコンテキストに注入される",
			buildReq: func(t *testing.T) *http.Request {
				return withCookie(
					httptest.NewRequest(http.MethodGet, "/", nil),
					"meAccessToken", makeMWValidAT(t, validSub),
				)
			},
			validateATFn: func(_ context.Context, _ string) (string, error) { return validSub, nil },
			wantStatus:   http.StatusOK,
			wantCalled:   true,
			wantID:       validSub,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			cap := &captureHandler{}
			h := newTestHandler(tc.validateATFn)
			mw := h.AuthMiddleware(cap)

			w := httptest.NewRecorder()
			mw.ServeHTTP(w, tc.buildReq(t))

			if w.Code != tc.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tc.wantStatus)
			}
			if cap.called != tc.wantCalled {
				t.Errorf("next called = %v, want %v", cap.called, tc.wantCalled)
			}
			if tc.wantID != "" && cap.id != tc.wantID {
				t.Errorf("identityID in context = %q, want %q", cap.id, tc.wantID)
			}
		})
	}
}

// TestHandler_AuthMiddleware_ValidateReceivesCookieToken は
// Validate が meAccessToken Cookie の値をそのまま受け取っているかを検証する。
// テーブルテストでは validateFn の引数を無視しているため、入力ソースの正しさは別途確認が必要。
func TestHandler_AuthMiddleware_ValidateReceivesCookieToken(t *testing.T) {
	t.Parallel()
	sub := uuid.New().String()
	token := makeMWValidAT(t, sub)

	var gotToken string
	h := newTestHandler(func(_ context.Context, tok string) (string, error) {
		gotToken = tok
		return sub, nil
	})

	cap := &captureHandler{}
	mw := h.AuthMiddleware(cap)

	r := withCookie(httptest.NewRequest(http.MethodGet, "/", nil), "meAccessToken", token)
	w := httptest.NewRecorder()
	mw.ServeHTTP(w, r)

	if gotToken != token {
		t.Errorf("Validate received %q, want cookie value %q", gotToken, token)
	}
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	if cap.id != sub {
		t.Errorf("identityID in context = %q, want %q", cap.id, sub)
	}
}

// --- RefreshMiddleware ---

// RefreshMiddleware は AT の署名検証を行わず ParseUnverified で sub を取得する。
// 期限切れ AT でも sub を抽出して next を呼ぶことで、RT ローテーション時の IdentityID 特定を担う。
func TestHandler_RefreshMiddleware(t *testing.T) {
	t.Parallel()

	validSub := uuid.New().String()

	cases := []struct {
		name       string
		buildReq   func(t *testing.T) *http.Request
		wantStatus int
		wantCalled bool
		wantID     string
	}{
		{
			name: "meAccessToken Cookie なし → 401",
			buildReq: func(t *testing.T) *http.Request {
				return httptest.NewRequest(http.MethodPost, "/", nil)
			},
			wantStatus: http.StatusUnauthorized,
			wantCalled: false,
		},
		{
			name: "JWT でない任意文字列 → 401",
			buildReq: func(t *testing.T) *http.Request {
				return withCookie(
					httptest.NewRequest(http.MethodPost, "/", nil),
					"meAccessToken", "not-a-jwt",
				)
			},
			wantStatus: http.StatusUnauthorized,
			wantCalled: false,
		},
		{
			name: "sub が UUID でない → 401",
			buildReq: func(t *testing.T) *http.Request {
				return withCookie(
					httptest.NewRequest(http.MethodPost, "/", nil),
					"meAccessToken", makeMWATWithSub(t, "not-a-uuid"),
				)
			},
			wantStatus: http.StatusUnauthorized,
			wantCalled: false,
		},
		{
			name: "有効な AT → sub がコンテキストに注入されて next が呼ばれる",
			buildReq: func(t *testing.T) *http.Request {
				return withCookie(
					httptest.NewRequest(http.MethodPost, "/", nil),
					"meAccessToken", makeMWValidAT(t, validSub),
				)
			},
			wantStatus: http.StatusOK,
			wantCalled: true,
			wantID:     validSub,
		},
		{
			name: "期限切れ AT でも sub を抽出して next を呼ぶ（AuthMiddleware との違い）",
			buildReq: func(t *testing.T) *http.Request {
				return withCookie(
					httptest.NewRequest(http.MethodPost, "/", nil),
					"meAccessToken", makeMWExpiredAT(t, validSub),
				)
			},
			wantStatus: http.StatusOK,
			wantCalled: true,
			wantID:     validSub,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			cap := &captureHandler{}
			h := newTestHandler(nil) // RefreshMiddleware は tokenSrv.ValidateAT を使用しない
			mw := h.RefreshMiddleware(cap)

			w := httptest.NewRecorder()
			mw.ServeHTTP(w, tc.buildReq(t))

			if w.Code != tc.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tc.wantStatus)
			}
			if cap.called != tc.wantCalled {
				t.Errorf("next called = %v, want %v", cap.called, tc.wantCalled)
			}
			if tc.wantID != "" && cap.id != tc.wantID {
				t.Errorf("identityID in context = %q, want %q", cap.id, tc.wantID)
			}
		})
	}
}

// --- identityIDFromContext ---

func TestIdentityIDFromContext(t *testing.T) {
	t.Parallel()

	t.Run("注入済みの identityID を取得できる", func(t *testing.T) {
		t.Parallel()
		want := uuid.New().String()
		ctx := context.WithValue(context.Background(), identityIDKey, want)
		got, ok := identityIDFromContext(ctx)
		if !ok {
			t.Fatal("expected ok = true")
		}
		if got != want {
			t.Errorf("id = %q, want %q", got, want)
		}
	})

	t.Run("キーが存在しない → 空文字と false", func(t *testing.T) {
		t.Parallel()
		got, ok := identityIDFromContext(context.Background())
		if ok {
			t.Error("expected ok = false")
		}
		if got != "" {
			t.Errorf("id = %q, want empty string", got)
		}
	})

	t.Run("値の型が string でない → 空文字と false", func(t *testing.T) {
		t.Parallel()
		ctx := context.WithValue(context.Background(), identityIDKey, 12345)
		got, ok := identityIDFromContext(ctx)
		if ok {
			t.Error("expected ok = false")
		}
		if got != "" {
			t.Errorf("id = %q, want empty string", got)
		}
	})
}
