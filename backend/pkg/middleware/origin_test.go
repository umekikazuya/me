package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestOriginGuard(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		allowedOrigins []string
		method         string
		origin         string
		wantStatus     int
		wantCalled     bool
	}{
		{
			name:           "configured origins are optional when allowlist is empty",
			allowedOrigins: nil,
			method:         http.MethodPost,
			wantStatus:     http.StatusNoContent,
			wantCalled:     true,
		},
		{
			name:           "safe methods skip validation",
			allowedOrigins: []string{"https://www.example.com"},
			method:         http.MethodGet,
			wantStatus:     http.StatusNoContent,
			wantCalled:     true,
		},
		{
			name:           "allowed origin passes",
			allowedOrigins: []string{"https://www.example.com/"},
			method:         http.MethodPost,
			origin:         "https://www.example.com",
			wantStatus:     http.StatusNoContent,
			wantCalled:     true,
		},
		{
			name:           "missing origin is denied for unsafe methods",
			allowedOrigins: []string{"https://www.example.com"},
			method:         http.MethodPost,
			wantStatus:     http.StatusForbidden,
			wantCalled:     false,
		},
		{
			name:           "disallowed origin is denied",
			allowedOrigins: []string{"https://www.example.com"},
			method:         http.MethodDelete,
			origin:         "https://evil.example.net",
			wantStatus:     http.StatusForbidden,
			wantCalled:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			called := false
			handler := OriginGuard(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				called = true
				w.WriteHeader(http.StatusNoContent)
			}), tt.allowedOrigins)

			req := httptest.NewRequest(tt.method, "/resource", nil)
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}

			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if called != tt.wantCalled {
				t.Fatalf("called = %v, want %v", called, tt.wantCalled)
			}
		})
	}
}

func TestOriginGuardErrorResponseUsesProblemDetail(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodPost, "/resource", nil)
	rec := httptest.NewRecorder()

	OriginGuard(
		http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}),
		[]string{"https://www.example.com"},
	).ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusForbidden)
	}
	if contentType := rec.Header().Get("Content-Type"); contentType != "application/problem+json" {
		t.Fatalf("Content-Type = %q, want %q", contentType, "application/problem+json")
	}
	if !strings.Contains(rec.Body.String(), `"title":"Forbidden"`) {
		t.Fatalf("body = %q, want problem detail response", rec.Body.String())
	}
}
