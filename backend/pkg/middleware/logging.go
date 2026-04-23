package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status int
	err    error
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw *responseWriter) RecordError(err error) {
	rw.err = err
}

func logLevel(status int) slog.Level {
	switch {
	case status >= http.StatusInternalServerError:
		return slog.LevelError
	case status >= http.StatusBadRequest:
		return slog.LevelWarn
	default:
		return slog.LevelInfo
	}
}

// Logging は HTTP アクセスログを logger に出力する middleware を返す。
// logger が nil の場合は slog.Default() にフォールバックする。
func Logging(logger *slog.Logger) func(http.Handler) http.Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(rw, r)

			args := []any{
				"method", r.Method,
				"path", r.URL.Path,
				"status", rw.status,
				"duration", time.Since(start),
				"remoteAddr", r.RemoteAddr,
				"userAgent", r.UserAgent(),
			}
			if rw.err != nil {
				args = append(args, "error", rw.err.Error())
			}

			logger.Log(r.Context(), logLevel(rw.status), "request", args...)
		})
	}
}
