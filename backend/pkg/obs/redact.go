package obs

import (
	"context"
	"log/slog"
	"strings"
)

const redactedValue = "[REDACTED]"

// redactHandler は指定 key と完全一致 (lowercase) する attribute の値を [REDACTED] に置換する。
// ネストした slog.Group 配下の key も辿ってマスクする。
type redactHandler struct {
	slog.Handler
	keys map[string]struct{}
}

func newRedactHandler(inner slog.Handler, keys []string) slog.Handler {
	m := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		m[strings.ToLower(k)] = struct{}{}
	}
	return &redactHandler{Handler: inner, keys: m}
}

func (h *redactHandler) Handle(ctx context.Context, r slog.Record) error {
	attrs := make([]slog.Attr, 0, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, h.redact(a))
		return true
	})
	newR := slog.NewRecord(r.Time, r.Level, r.Message, r.PC)
	newR.AddAttrs(attrs...)
	return h.Handler.Handle(ctx, newR)
}

func (h *redactHandler) redact(a slog.Attr) slog.Attr {
	if _, ok := h.keys[strings.ToLower(a.Key)]; ok {
		return slog.String(a.Key, redactedValue)
	}
	// LogValuer は lazy resolve のため、resolve してから redact を再適用する。
	// Resolve しないと `password=xxx` を LogValue 内で返す型が素通りしてしまう。
	if a.Value.Kind() == slog.KindLogValuer {
		resolved := a.Value.Resolve()
		return h.redact(slog.Attr{Key: a.Key, Value: resolved})
	}
	if a.Value.Kind() == slog.KindGroup {
		group := a.Value.Group()
		redacted := make([]slog.Attr, len(group))
		for i, g := range group {
			redacted[i] = h.redact(g)
		}
		return slog.Attr{Key: a.Key, Value: slog.GroupValue(redacted...)}
	}
	return a
}

func (h *redactHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	redacted := make([]slog.Attr, len(attrs))
	for i, a := range attrs {
		redacted[i] = h.redact(a)
	}
	return &redactHandler{Handler: h.Handler.WithAttrs(redacted), keys: h.keys}
}

func (h *redactHandler) WithGroup(name string) slog.Handler {
	return &redactHandler{Handler: h.Handler.WithGroup(name), keys: h.keys}
}
