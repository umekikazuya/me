package identity_test

import (
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	app "github.com/umekikazuya/me/internal/app/identity"
	"github.com/umekikazuya/me/internal/handler/identity"
)

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
			r:          httptest.NewRequest(http.MethodGet, "/auth/github/callback", nil),
			wantCode:   http.StatusFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := identity.NewHandler(tt.interactor, tt.tokenSrv)
			h.CallbackGithub(tt.w, tt.r)

			if tt.w.Code == tt.wantCode {
				t.Errorf("tt.w.Code = %v, want = %v", tt.w.Code, tt.wantCode)
			}
		})
	}
}
