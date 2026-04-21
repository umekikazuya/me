package slogx

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/umekikazuya/me/pkg/reqctx"
)

// New は LOG_LEVEL env (debug/info/warn/error, 既定 info) を反映した JSON slog.Logger を返す。
// context の RequestID は全ログ entry に requestId フィールドとして自動付与される。
func New(w io.Writer) *slog.Logger {
	if w == nil {
		w = os.Stdout
	}
	level := parseLevel(os.Getenv("LOG_LEVEL"))
	base := slog.NewJSONHandler(w, &slog.HandlerOptions{Level: level})
	return slog.New(&contextHandler{Handler: base})
}

func parseLevel(s string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

type contextHandler struct{ slog.Handler }

func (h *contextHandler) Handle(ctx context.Context, r slog.Record) error {
	if id := reqctx.RequestIDFromContext(ctx); id != "" {
		r.AddAttrs(slog.String("requestId", id))
	}
	if err := h.Handler.Handle(ctx, r); err != nil {
		return fmt.Errorf("slog handler: %w", err)
	}
	return nil
}

func (h *contextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &contextHandler{Handler: h.Handler.WithAttrs(attrs)}
}

func (h *contextHandler) WithGroup(name string) slog.Handler {
	return &contextHandler{Handler: h.Handler.WithGroup(name)}
}
