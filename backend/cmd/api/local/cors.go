package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/umekikazuya/me/pkg/middleware"
)

const corsAllowedOriginsEnv = "CORS_ALLOWED_ORIGINS"

func loadCORSConfig() middleware.CORSConfig {
	return middleware.CORSConfig{
		AllowedOrigins: splitEnvList(os.Getenv(corsAllowedOriginsEnv)),
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
}

func splitEnvList(raw string) []string {
	if raw == "" {
		return nil
	}

	values := strings.Split(raw, ",")
	result := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		normalizedValue := strings.TrimSpace(value)
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
