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
			var got map[string]any
			logger := slog.New(slog.NewJSONHandler(testWriter(func(p []byte) (int, error) {
				return len(p), json.Unmarshal(p, &got)
			}), nil))
			prev := slog.Default()
			slog.SetDefault(logger)
			t.Cleanup(func() { slog.SetDefault(prev) })

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			rec := httptest.NewRecorder()

			Recover(tt.handler).ServeHTTP(rec, req)

			if !tt.wantLogged {
				if got != nil {
					t.Fatalf("想定外のログ: %v", got)
				}
				if rec.Code != http.StatusOK {
					t.Fatalf("status = %d, want 200", rec.Code)
				}
				return
			}

			if got == nil {
				t.Fatal("ログが出ていない")
			}
			if got["level"] != "ERROR" {
				t.Fatalf("level = %v, want ERROR", got["level"])
			}
			if got["msg"] != "unhandled panic" {
				t.Fatalf("msg = %v, want 'unhandled panic'", got["msg"])
			}
			if tt.wantStack {
				stack, _ := got["stack"].(string)
				if stack == "" || !strings.Contains(stack, "goroutine") {
					t.Fatalf("stack が想定通りでない: %q", stack)
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

type testWriter func([]byte) (int, error)

func (w testWriter) Write(p []byte) (int, error) {
	return w(p)
}
