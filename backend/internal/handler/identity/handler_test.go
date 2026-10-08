package identity_test

import (
	"net/http"
	"net/http/httptest"
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
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := identity.NewHandler(tt.interactor, tt.tokenSrv)
			h.LoginFromGithub(tt.w, tt.r)
			if tt.w.Code != tt.wantCode {
				t.Errorf("w.Code = %v, want = %v", tt.w.Code, tt.wantCode)
			}
		})
	}
}
