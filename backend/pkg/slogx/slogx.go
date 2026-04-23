// Package slogx はアプリ全体で共有する基盤ロガーを提供する。
// 関心事は JSON 出力 / LOG_LEVEL env / context 由来の requestId 自動注入のみ。
// component (app/http/batch など) タグは呼び出し側で `slog.With("component", ...)` として付ける。
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

// New は基盤ロガーを返す。LOG_LEVEL env (debug/info/warn/error, 既定 info) を反映した
// JSON handler に、context の RequestID を全 entry に `requestId` として自動付与する
// contextHandler を積んだ構成。component は付けないので、app / http / batch 等の派生は
// 呼び出し側で `.With("component", "...")` して作る。
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
