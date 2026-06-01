package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCORS(t *testing.T) {
	t.Parallel()

	cfg := CORSConfig{
		AllowedOrigins: []string{"https://www.example.com"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Accept",
			"Content-Type",
			"X-Request-ID",
			"X-Requested-With",
		},
		ExposedHeaders: []string{"X-Request-ID"},
	}

	tests := []struct {
		name            string
		buildRequest    func() *http.Request
		wantStatus      int
		wantCalled      bool
		wantAllowOrigin string
		wantCredentials string
		wantExpose      string
		wantProblemJSON bool
		wantVary        []string
	}{
		{
			name: "Origin なしはそのまま通す",
			buildRequest: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, "/articles", nil)
			},
			wantStatus: http.StatusOK,
			wantCalled: true,
		},
		{
			name: "許可された Origin の通常リクエストは CORS ヘッダーを付けて通す",
			buildRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/articles", nil)
				req.Header.Set("Origin", "https://www.example.com")
				return req
			},
			wantStatus:      http.StatusOK,
			wantCalled:      true,
			wantAllowOrigin: "https://www.example.com",
			wantCredentials: "true",
			wantExpose:      "X-Request-ID",
			wantVary:        []string{"Origin"},
		},
		{
			name: "許可された Origin の preflight は 204 を返す",
			buildRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodOptions, "/auth/login", nil)
				req.Header.Set("Origin", "https://www.example.com")
				req.Header.Set("Access-Control-Request-Method", http.MethodPost)
				req.Header.Set("Access-Control-Request-Headers", "content-type, x-requested-with")
				return req
			},
			wantStatus:      http.StatusNoContent,
			wantAllowOrigin: "https://www.example.com",
			wantCredentials: "true",
			wantExpose:      "X-Request-ID",
			wantVary: []string{
				"Origin",
				"Access-Control-Request-Method",
				"Access-Control-Request-Headers",
			},
		},
		{
			name: "許可されていない Origin の通常リクエストは通すが CORS を付けない",
			buildRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/articles", nil)
				req.Header.Set("Origin", "https://evil.example.net")
				return req
			},
			wantStatus: http.StatusOK,
			wantCalled: true,
			wantVary:   []string{"Origin"},
		},
		{
			name: "許可されていない Origin の preflight は 403",
			buildRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodOptions, "/auth/login", nil)
				req.Header.Set("Origin", "https://evil.example.net")
				req.Header.Set("Access-Control-Request-Method", http.MethodPost)
				return req
			},
			wantStatus:      http.StatusForbidden,
			wantProblemJSON: true,
			wantVary:        []string{"Origin"},
		},
		{
			name: "未許可ヘッダーを含む preflight は 403",
			buildRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodOptions, "/auth/login", nil)
				req.Header.Set("Origin", "https://www.example.com")
				req.Header.Set("Access-Control-Request-Method", http.MethodPost)
				req.Header.Set("Access-Control-Request-Headers", "x-requested-with, x-evil")
				return req
			},
			wantStatus:      http.StatusForbidden,
			wantProblemJSON: true,
			wantAllowOrigin: "https://www.example.com",
			wantCredentials: "true",
			wantExpose:      "X-Request-ID",
			wantVary: []string{
				"Origin",
				"Access-Control-Request-Method",
				"Access-Control-Request-Headers",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			called := false
			handler := CORS(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				called = true
				w.WriteHeader(http.StatusOK)
			}), cfg)

			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, tt.buildRequest())

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if called != tt.wantCalled {
				t.Fatalf("called = %v, want %v", called, tt.wantCalled)
			}
			if allowOrigin := rec.Header().Get("Access-Control-Allow-Origin"); allowOrigin != tt.wantAllowOrigin {
				t.Fatalf("Access-Control-Allow-Origin = %q, want %q", allowOrigin, tt.wantAllowOrigin)
			}
			if credentials := rec.Header().Get("Access-Control-Allow-Credentials"); credentials != tt.wantCredentials {
				t.Fatalf("Access-Control-Allow-Credentials = %q, want %q", credentials, tt.wantCredentials)
			}
			if exposeHeaders := rec.Header().Get("Access-Control-Expose-Headers"); exposeHeaders != tt.wantExpose {
				t.Fatalf("Access-Control-Expose-Headers = %q, want %q", exposeHeaders, tt.wantExpose)
			}
			if tt.wantProblemJSON {
				if contentType := rec.Header().Get("Content-Type"); contentType != "application/problem+json" {
					t.Fatalf("Content-Type = %q, want application/problem+json", contentType)
				}
			}

			vary := strings.Join(rec.Header().Values("Vary"), ",")
			for _, token := range tt.wantVary {
				if !strings.Contains(vary, token) {
					t.Fatalf("Vary = %q, want token %q", vary, token)
				}
			}
		})
	}
}
