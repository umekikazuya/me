package middleware

import (
	"net/http"

	"github.com/umekikazuya/me/pkg/errs"
)

func OriginGuard(next http.Handler, allowedOrigins []string) http.Handler {
	allowedOriginSet := newOriginAllowlist(allowedOrigins)
	if len(allowedOriginSet) == 0 {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !requiresOriginValidation(r.Method) {
			next.ServeHTTP(w, r)
			return
		}

		origin := normalizeOrigin(r.Header.Get("Origin"))
		if origin == "" {
			errs.WriteProblem(w, r, errs.ErrPermissionDenied)
			return
		}
		if _, ok := allowedOriginSet[origin]; !ok {
			errs.WriteProblem(w, r, errs.ErrPermissionDenied)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func newOriginAllowlist(origins []string) map[string]struct{} {
	allowedOrigins := make(map[string]struct{}, len(origins))
	for _, origin := range origins {
		normalizedOrigin := normalizeOrigin(origin)
		if normalizedOrigin == "" {
			continue
		}
		allowedOrigins[normalizedOrigin] = struct{}{}
	}
	return allowedOrigins
}

func requiresOriginValidation(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodTrace:
		return false
	default:
		return true
	}
}
