package obs

import (
	"fmt"

	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func newTracerProvider(cfg Config, res *resource.Resource) (*sdktrace.TracerProvider, error) {
	exp, err := stdouttrace.New(stdouttrace.WithWriter(cfg.Writer))
	if err != nil {
		return nil, fmt.Errorf("stdouttrace: %w", err)
	}
	// SyncExport: dev/debug では span を即 stdout に流すため SimpleSpanProcessor を使う。
	// 本番 (false) では既定の BatchSpanProcessor でスループットを確保する。
	spanOpt := sdktrace.WithBatcher(exp)
	if cfg.SyncExport {
		spanOpt = sdktrace.WithSyncer(exp)
	}
	return sdktrace.NewTracerProvider(
		spanOpt,
		sdktrace.WithResource(res),
	), nil
}
