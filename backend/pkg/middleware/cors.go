package middleware

import (
	"net/http"
	"strings"

	"github.com/umekikazuya/me/pkg/errs"
)

type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
	ExposedHeaders []string
}

func CORS(next http.Handler, cfg CORSConfig) http.Handler {
	allowedOrigins := newOriginAllowlist(cfg.AllowedOrigins)

	allowedMethods := normalizeMethods(cfg.AllowedMethods)
	allowedMethodSet := make(map[string]struct{}, len(allowedMethods))
	for _, method := range allowedMethods {
		allowedMethodSet[method] = struct{}{}
	}

	allowedHeaders := normalizeHeaders(cfg.AllowedHeaders)
	allowedHeaderSet := make(map[string]struct{}, len(allowedHeaders))
	for _, header := range allowedHeaders {
		allowedHeaderSet[strings.ToLower(header)] = struct{}{}
	}

	exposedHeaders := normalizeHeaders(cfg.ExposedHeaders)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := normalizeOrigin(r.Header.Get("Origin"))
		if origin == "" {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Add("Vary", "Origin")

		if _, ok := allowedOrigins[origin]; !ok {
			if isPreflight(r) {
				errs.WriteProblem(w, r, errs.ErrPermissionDenied)
				return
			}
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if len(exposedHeaders) > 0 {
			w.Header().Set("Access-Control-Expose-Headers", strings.Join(exposedHeaders, ", "))
		}

		if !isPreflight(r) {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Add("Vary", "Access-Control-Request-Method")
		w.Header().Add("Vary", "Access-Control-Request-Headers")

		requestMethod := strings.ToUpper(strings.TrimSpace(r.Header.Get("Access-Control-Request-Method")))
		if _, ok := allowedMethodSet[requestMethod]; !ok {
			errs.WriteProblem(w, r, errs.ErrPermissionDenied)
			return
		}

		requestHeaders := parseHeaderList(r.Header.Get("Access-Control-Request-Headers"))
		for _, header := range requestHeaders {
			if _, ok := allowedHeaderSet[strings.ToLower(header)]; !ok {
				errs.WriteProblem(w, r, errs.ErrPermissionDenied)
				return
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ", "))
		if len(allowedHeaders) > 0 {
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ", "))
		}
		w.WriteHeader(http.StatusNoContent)
	})
}

func isPreflight(r *http.Request) bool {
	return r.Method == http.MethodOptions && strings.TrimSpace(r.Header.Get("Access-Control-Request-Method")) != ""
}

func normalizeOrigin(origin string) string {
	return strings.TrimRight(strings.TrimSpace(origin), "/")
}

func normalizeMethods(methods []string) []string {
	return uniqueNormalized(methods, strings.ToUpper)
}

func normalizeHeaders(headers []string) []string {
	result := make([]string, 0, len(headers))
	seen := make(map[string]struct{}, len(headers))
	for _, header := range headers {
		normalizedHeader := strings.TrimSpace(header)
		if normalizedHeader == "" {
			continue
		}

		comparisonKey := strings.ToLower(normalizedHeader)
		if _, ok := seen[comparisonKey]; ok {
			continue
		}
		seen[comparisonKey] = struct{}{}
		result = append(result, normalizedHeader)
	}
	return result
}

func uniqueNormalized(values []string, normalize func(string) string) []string {
	result := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		normalizedValue := normalize(strings.TrimSpace(value))
		if normalizedValue == "" {
			continue
		}
		if _, ok := seen[normalizedValue]; ok {
			continue
		}
		seen[normalizedValue] = struct{}{}
		result = append(result, normalizedValue)
	}
	return result
}

func parseHeaderList(raw string) []string {
	if raw == "" {
		return nil
	}

	return normalizeHeaders(strings.Split(raw, ","))
}
