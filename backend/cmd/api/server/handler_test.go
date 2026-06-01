package server

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/umekikazuya/me/pkg/middleware"
)

func TestSplitEnvList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		raw  string
		want []string
	}{
		{
			name: "empty string returns nil",
			raw:  "",
			want: nil,
		},
		{
			name: "trim and deduplicate values",
			raw:  " https://www.example.com ,https://api.example.com,https://www.example.com, ",
			want: []string{
				"https://www.example.com",
				"https://api.example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := splitEnvList(tt.raw); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("splitEnvList(%q) = %v, want %v", tt.raw, got, tt.want)
			}
		})
	}
}

func TestLoadCORSConfig(t *testing.T) {
	t.Setenv(corsAllowedOriginsEnv, " https://www.example.com ,https://admin.example.com,https://www.example.com ")

	got := LoadCORSConfig()
	want := middleware.CORSConfig{
		AllowedOrigins: []string{"https://www.example.com", "https://admin.example.com"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodHead,
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

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("LoadCORSConfig() = %#v, want %#v", got, want)
	}
}

func TestNewHandlerIncludesCORS(t *testing.T) {
	t.Parallel()

	handler := NewHandler(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
		middleware.CORSConfig{
			AllowedOrigins: []string{"https://www.example.com"},
			AllowedMethods: []string{http.MethodPost, http.MethodOptions},
			AllowedHeaders: []string{"Content-Type", "X-Request-ID", "X-Requested-With"},
			ExposedHeaders: []string{"X-Request-ID"},
		},
	)

	req := httptest.NewRequest(http.MethodOptions, "/auth/login", nil)
	req.Header.Set("Origin", "https://www.example.com")
	req.Header.Set("Access-Control-Request-Method", http.MethodPost)
	req.Header.Set("Access-Control-Request-Headers", "Content-Type, X-Request-ID")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNoContent)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "https://www.example.com" {
		t.Fatalf("Access-Control-Allow-Origin = %q, want %q", got, "https://www.example.com")
	}
	if got := rec.Header().Get("Access-Control-Allow-Credentials"); got != "true" {
		t.Fatalf("Access-Control-Allow-Credentials = %q, want %q", got, "true")
	}
}
