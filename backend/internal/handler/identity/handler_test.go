package identity_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
	"time"

	app "github.com/umekikazuya/me/internal/app/identity"
	"github.com/umekikazuya/me/internal/handler/identity"
)

var (
	fixtureAccessToken  = "accessTokenValue"
	fixtureRefreshToken = "refreshTokenValue"
	baseTime            = time.Now()
)

type fakeInteractor struct {
	callBackDto app.InputCallbackFromGithubDto
	err         error
}

// CallbackFromGithub implements [identity.Interactor].
func (f *fakeInteractor) CallbackFromGithub(
	ctx context.Context,
	in app.InputCallbackFromGithubDto,
) (*app.OutputCallbackWithGithubDto, error) {
	if f.err != nil {
		return nil, f.err
	}
	return &app.OutputCallbackWithGithubDto{
		AT: fixtureAccessToken,
		RT: fixtureRefreshToken,
	}, nil
}

// LoginWithGithub implements [identity.Interactor].
func (f *fakeInteractor) LoginWithGithub(ctx context.Context, in app.InputLoginWithGithubDto) error {
	panic("unimplemented")
}

// Logout implements [identity.Interactor].
func (f *fakeInteractor) Logout(ctx context.Context, in app.InputLogoutDto) error {
	panic("unimplemented")
}

// RefreshTokens implements [identity.Interactor].
func (f *fakeInteractor) RefreshTokens(ctx context.Context, in app.InputRefreshTokensDto) (*app.OutputRefreshTokensDto, error) {
	panic("unimplemented")
}

// RevokeAllSessions implements [identity.Interactor].
func (f *fakeInteractor) RevokeAllSessions(ctx context.Context, in app.InputRevokeAllSessionsDto) error {
	panic("unimplemented")
}

var _ app.Interactor = (*fakeInteractor)(nil)

func TestHandler_LoginFromGithub(t *testing.T) {
	tests := []struct {
		name       string
		interactor app.Interactor
		tokenSrv   app.TokenService
		w          *httptest.ResponseRecorder
		r          *http.Request
		wantCode   int
	}{
		{
			name:       "ok",
			interactor: nil,
			tokenSrv:   nil,
			w:          httptest.NewRecorder(),
			r: httptest.NewRequest(
				http.MethodGet,
				"/auth/github/login",
				nil,
			),
			wantCode: 302,
		},
		{
			name:       "ok",
			interactor: nil,
			tokenSrv:   nil,
			w:          httptest.NewRecorder(),
			r: httptest.NewRequest(
				http.MethodGet,
				"/auth/github/login",
				nil,
			),
			wantCode: 302,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := identity.NewHandler(
				tt.interactor,
				tt.tokenSrv,
			)
			h.LoginFromGithub(tt.w, tt.r)
			if tt.w.Code != tt.wantCode {
				t.Errorf("w.Code = %v, want = %v", tt.w.Code, tt.wantCode)
			}
			if !slices.ContainsFunc(
				tt.w.Result().Cookies(),
				func(c *http.Cookie) bool {
					return c.Name == "authState"
				},
			) {
				t.Errorf("cookies = %v", tt.w.Result().Cookies())
			}
			redirectURL, _ := tt.w.Result().Location()
			if redirectURL.Query().Get("state") == "" {
				t.Errorf("redirectURL.Query() = %v", redirectURL.Query())
			}
		})
	}
}

func TestHandler_CallbackGithub(t *testing.T) {
	tests := []struct {
		name           string
		interactor     app.Interactor
		tokenSrv       app.TokenService
		w              *httptest.ResponseRecorder
		r              *http.Request
		seedCookieFn   func(t *testing.T, req *http.Request)
		wantCode       int
		assertCookieFn func(t *testing.T, res *http.Response)
	}{
		{
			name: "ok",
			interactor: &fakeInteractor{
				callBackDto: app.InputCallbackFromGithubDto{
					Code:     "1234",
					BaseTime: baseTime,
				},
				err: nil,
			},
			tokenSrv: nil,
			w:        httptest.NewRecorder(),
			r: httptest.NewRequest(
				http.MethodGet,
				"/auth/github/callback?state=testAbc",
				nil,
			),
			seedCookieFn: func(t *testing.T, req *http.Request) {
				t.Helper()
				req.AddCookie(
					&http.Cookie{
						Name:     "authState",
						Value:    "testAbc",
						Expires:  time.Now().Add(20 * time.Minute),
						Quoted:   false,
						HttpOnly: true,
					},
				)
			},
			wantCode: http.StatusFound,
			assertCookieFn: func(t *testing.T, res *http.Response) {
				t.Helper()
				var count int
				for _, c := range res.Cookies() {
					if c.Name == "meAccessToken" || c.Name == "meRefreshToken" {
						count++
						if c.Value == "" {
							t.Errorf("cookie: %v is empty", c.Name)
						}
					}
				}
				if count != 2 {
					t.Errorf("cookie count is %d", count)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := identity.NewHandler(tt.interactor, tt.tokenSrv)
			tt.seedCookieFn(t, tt.r)
			h.CallbackGithub(tt.w, tt.r)

			if tt.w.Code != tt.wantCode {
				t.Errorf(
					"tt.w.Code = %v, want = %v, body = %v",
					tt.w.Code,
					tt.wantCode,
					tt.w.Body,
				)
			}
			tt.assertCookieFn(t, tt.w.Result())
		})
	}
}
