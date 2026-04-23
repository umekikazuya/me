// Package slogx はアプリの業務ログ基盤を提供する。
//
// HTTP アクセスログ (method/path/status/duration) はインフラ層
// (API Gateway / ALB / CloudFront 等) の責務とし、アプリでは出さない。
// 本パッケージが担うのは業務ログのみ。
//
// 関心事:
//   - JSON 出力
//   - LOG_LEVEL env での level 制御 (debug/info/warn/error, 既定 info)
//   - context 由来の requestId 自動注入 (contextHandler)
//
// 推奨される使い方:
//
//	// プロセス起動時
//	slog.SetDefault(slogx.New(os.Stdout).With("service", "api"))
//
//	// 業務ログ (アプリ / ユースケース / インフラアダプタ層から)
//	slog.ErrorContext(ctx, "infra error",
//	    "component", "infra",
//	    "op", "dynamo.save",
//	    "error", err,
//	)
//
// 属性の軸:
//   - service   — プロセス識別 (api / batch / worker 等)。起動時に必ず付与
//   - component — 論理レイヤ (infra 等)。レイヤ識別が必要な時のみ付与 (任意)
//   - requestId — context から自動注入 (contextHandler 経由)
//
// アクセスログは出さないが、panic の安全網として middleware.Recover が ERROR を残す。
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

// New は業務ログ基盤の logger を返す。service 属性は呼び出し側で `.With("service", ...)`
// として付与すること (プロセス起動時に 1 度)。component はレイヤ識別が必要な場面のみ付ける。
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

