package identity

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// findCookie は Set-Cookie ヘッダから指定名の Cookie を返す
func findCookie(t *testing.T, w *httptest.ResponseRecorder, name string) *http.Cookie {
	t.Helper()
	for _, c := range w.Result().Cookies() {
		if c.Name == name {
			return c
		}
	}
	return nil
}

// --- setTokenCookies ---

func TestSetTokenCookies_AT(t *testing.T) {
	t.Parallel()
	w := httptest.NewRecorder()
	setTokenCookies(w, "at-value", "rt-value")

	c := findCookie(t, w, accessTokenCookieName)
	if c == nil {
		t.Fatalf("cookie %q not found in Set-Cookie header", accessTokenCookieName)
	}

	t.Run("value", func(t *testing.T) {
		if c.Value != "at-value" {
			t.Errorf("Value = %q, want %q", c.Value, "at-value")
		}
	})
	t.Run("HttpOnly", func(t *testing.T) {
		if !c.HttpOnly {
			t.Error("HttpOnly must be true")
		}
	})
	t.Run("Secure", func(t *testing.T) {
		if !c.Secure {
			t.Error("Secure must be true")
		}
	})
	t.Run("SameSite=Strict", func(t *testing.T) {
		if c.SameSite != http.SameSiteStrictMode {
			t.Errorf("SameSite = %v, want Strict", c.SameSite)
		}
	})
	t.Run("MaxAge", func(t *testing.T) {
		if c.MaxAge != atMaxAge {
			t.Errorf("MaxAge = %d, want %d", c.MaxAge, atMaxAge)
		}
	})
	t.Run("Path=/", func(t *testing.T) {
		if c.Path != "/" {
			t.Errorf("Path = %q, want %q", c.Path, "/")
		}
	})
}

func TestSetTokenCookies_RT(t *testing.T) {
	t.Parallel()
	w := httptest.NewRecorder()
	setTokenCookies(w, "at-value", "rt-value")

	c := findCookie(t, w, refreshTokenCookieName)
	if c == nil {
		t.Fatalf("cookie %q not found in Set-Cookie header", refreshTokenCookieName)
	}

	t.Run("value", func(t *testing.T) {
		if c.Value != "rt-value" {
			t.Errorf("Value = %q, want %q", c.Value, "rt-value")
		}
	})
	t.Run("HttpOnly", func(t *testing.T) {
		if !c.HttpOnly {
			t.Error("HttpOnly must be true")
		}
	})
	t.Run("Secure", func(t *testing.T) {
		if !c.Secure {
			t.Error("Secure must be true")
		}
	})
	t.Run("SameSite=Strict", func(t *testing.T) {
		if c.SameSite != http.SameSiteStrictMode {
			t.Errorf("SameSite = %v, want Strict", c.SameSite)
		}
	})
	t.Run("MaxAge", func(t *testing.T) {
		if c.MaxAge != rtMaxAge {
			t.Errorf("MaxAge = %d, want %d", c.MaxAge, rtMaxAge)
		}
	})
	t.Run("Path はリフレッシュエンドポイントに限定", func(t *testing.T) {
		if c.Path != "/identity/refresh" {
			t.Errorf("Path = %q, want %q", c.Path, "/identity/refresh")
		}
	})
}

// --- clearTokenCookies ---

func TestClearTokenCookies(t *testing.T) {
	t.Parallel()

	cases := []struct {
		cookieName string
	}{
		{accessTokenCookieName},
		{refreshTokenCookieName},
	}

	for _, tc := range cases {
		t.Run(tc.cookieName+" が削除される", func(t *testing.T) {
			t.Parallel()
			w := httptest.NewRecorder()
			clearTokenCookies(w)

			c := findCookie(t, w, tc.cookieName)
			if c == nil {
				t.Fatalf("cookie %q not found (must be present with MaxAge=-1 to delete)", tc.cookieName)
			}
			if c.MaxAge != -1 {
				t.Errorf("MaxAge = %d, want -1 (deletion)", c.MaxAge)
			}
			if c.Value != "" {
				t.Errorf("Value = %q, want empty on clear", c.Value)
			}
		})
	}
}
