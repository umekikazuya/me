package obs

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"strings"

	"go.opentelemetry.io/otel/trace"
)

// ParseLevel は LOG_LEVEL 文字列を slog.Level に変換する純関数。
// 既定値は Info。main 境界から cfg.Level に詰める用途。
func ParseLevel(s string) slog.Level {
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

func buildLogger(cfg Config) *slog.Logger {
	var h slog.Handler = slog.NewJSONHandler(cfg.Writer, &slog.HandlerOptions{
		Level: cfg.Level,
	})
	h = &traceHandler{Handler: h, addSource: cfg.AddSource}
	if len(cfg.SensitiveKeys) > 0 {
		h = newRedactHandler(h, cfg.SensitiveKeys)
	}
	logger := slog.New(h).With(AttrServiceName, cfg.ServiceName)
	if cfg.ServiceVersion != "" {
		logger = logger.With(AttrServiceVersion, cfg.ServiceVersion)
	}
	return logger
}

// traceHandler は context 由来の属性 (request.id / trace_id / span_id) を自動で注入し、
// ERROR 以上のときのみ source を付与する slog.Handler デコレータ。
type traceHandler struct {
	slog.Handler
	addSource bool
}

func (h *traceHandler) Handle(ctx context.Context, r slog.Record) error {
	if id := RequestIDFromContext(ctx); id != "" {
		r.AddAttrs(slog.String(AttrRequestID, id))
	}
	if span := trace.SpanFromContext(ctx); span.SpanContext().IsValid() {
		sc := span.SpanContext()
		r.AddAttrs(
			slog.String(AttrTraceID, sc.TraceID().String()),
			slog.String(AttrSpanID, sc.SpanID().String()),
		)
	}
	if h.addSource && r.Level >= slog.LevelError && r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		if f.File != "" {
			r.AddAttrs(slog.Any(slog.SourceKey, &slog.Source{
				Function: f.Function,
				File:     f.File,
				Line:     f.Line,
			}))
		}
	}
	if err := h.Handler.Handle(ctx, r); err != nil {
		return fmt.Errorf("obs handler: %w", err)
	}
	return nil
}

func (h *traceHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &traceHandler{Handler: h.Handler.WithAttrs(attrs), addSource: h.addSource}
}

func (h *traceHandler) WithGroup(name string) slog.Handler {
	return &traceHandler{Handler: h.Handler.WithGroup(name), addSource: h.addSource}
}
