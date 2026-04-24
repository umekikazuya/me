package middleware

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/umekikazuya/me/pkg/errs"
	"github.com/umekikazuya/me/pkg/obs"
	"github.com/umekikazuya/me/pkg/obs/obstest"
)

func TestRecover(t *testing.T) {
	tests := []struct {
		name       string
		handler    http.Handler
		wantLogged bool
		wantStack  bool
	}{
		{
			name: "正常系は何もログしない",
			handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
			wantLogged: false,
		},
		{
			name: "error panic は 500 + ERROR ログ",
			handler: http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
				panic(errors.New("boom"))
			}),
			wantLogged: true,
			wantStack:  true,
		},
		{
			name: "非 error panic も 500 + ERROR ログ",
			handler: http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
				panic("boom")
			}),
			wantLogged: true,
			wantStack:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cap := obstest.NewCapture(t)
			logger := slog.New(cap.Handler())

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			rec := httptest.NewRecorder()

			Recover(tt.handler, WithLogger(logger)).ServeHTTP(rec, req)

			records := cap.Records()
			if !tt.wantLogged {
				if len(records) > 0 {
					t.Fatalf("想定外のログ: %v", records)
				}
				if rec.Code != http.StatusOK {
					t.Fatalf("status = %d, want 200", rec.Code)
				}
				return
			}

			if len(records) == 0 {
				t.Fatal("ログが出ていない")
			}
			got := records[0]
			if got["level"] != "ERROR" {
				t.Fatalf("level = %v, want ERROR", got["level"])
			}
			if got["msg"] != "unhandled panic" {
				t.Fatalf("msg = %v, want 'unhandled panic'", got["msg"])
			}
			if tt.wantStack {
				stack, _ := got[obs.AttrExceptionStack].(string)
				if stack == "" || !strings.Contains(stack, "goroutine") {
					t.Fatalf("exception.stacktrace が想定通りでない: %q", stack)
				}
			}
			if rec.Code != http.StatusInternalServerError {
				t.Fatalf("status = %d, want 500", rec.Code)
			}
			if !strings.Contains(rec.Header().Get("Content-Type"), "problem+json") {
				t.Fatalf("Content-Type = %q, want problem+json", rec.Header().Get("Content-Type"))
			}

			var body errs.ProblemDetail
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("body decode: %v", err)
			}
			if body.Status != http.StatusInternalServerError {
				t.Fatalf("body.Status = %d, want 500", body.Status)
			}
		})
	}
}
