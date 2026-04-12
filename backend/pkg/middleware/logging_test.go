package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/umekikazuya/me/pkg/errs"
)

func TestLogging(t *testing.T) {
	tests := []struct {
		name       string
		handler    http.Handler
		wantLevel  string
		wantStatus int
		wantError  string
	}{
		{
			name: "2xx は info で出る",
			handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
			wantLevel:  "INFO",
			wantStatus: http.StatusOK,
		},
		{
			name: "errs ベースの 4xx は warn と error を出す",
			handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				errs.WriteProblem(w, fmt.Errorf("decode request body: %w", errs.ErrBadRequest))
			}),
			wantLevel:  "WARN",
			wantStatus: http.StatusBadRequest,
			wantError:  "decode request body: bad request",
		},
		{
			name: "未知エラーの 5xx は error と元エラーを出す",
			handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				errs.WriteProblem(w, errors.New("database unavailable"))
			}),
			wantLevel:  "ERROR",
			wantStatus: http.StatusInternalServerError,
			wantError:  "database unavailable",
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
			t.Cleanup(func() {
				slog.SetDefault(prev)
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			rec := httptest.NewRecorder()

			Logging(tt.handler).ServeHTTP(rec, req)

			if got["level"] != tt.wantLevel {
				t.Fatalf("level = %v, want %q", got["level"], tt.wantLevel)
			}
			if int(got["status"].(float64)) != tt.wantStatus {
				t.Fatalf("status = %v, want %d", got["status"], tt.wantStatus)
			}
			if got["method"] != http.MethodGet {
				t.Fatalf("method = %v, want %q", got["method"], http.MethodGet)
			}
			if got["path"] != "/test" {
				t.Fatalf("path = %v, want %q", got["path"], "/test")
			}

			gotError, ok := got["error"]
			if tt.wantError == "" {
				if ok {
					t.Fatalf("unexpected error field = %v", gotError)
				}
				return
			}
			if !ok {
				t.Fatal("expected error field to be present")
			}
			if gotError != tt.wantError {
				t.Fatalf("error = %v, want %q", gotError, tt.wantError)
			}
		})
	}
}

type testWriter func([]byte) (int, error)

func (w testWriter) Write(p []byte) (int, error) {
	return w(p)
}
