// Package obs はアプリの観測性基盤 (logs / traces / metrics) を提供する。
//
// 設計の position は v1 (2026-04 時点、2 ヶ月後見直し前提) で、
// `docs/developments/observability.md` に明文化されている。
//
// 原則:
//   - アクセスログはインフラ層 (API Gateway / ALB / CloudFront) の責務、アプリでは出さない。
//   - 3 本柱は全て stdout に出力する (OTLP は v1 範囲外)。
//   - 属性名は OpenTelemetry Semantic Conventions に準拠する (定数は attr.go)。
//
// 推奨される使い方:
//
//	// プロセス起動時
//	prov, shutdown, err := obs.Bootstrap(ctx, obs.Config{
//	    ServiceName:   "api",
//	    Level:         obs.ParseLevel(os.Getenv("LOG_LEVEL")),
//	    EnableTraces:  true,
//	    EnableMetrics: true,
//	})
//	if err != nil { ... }
//	defer shutdown(ctx)
//	slog.SetDefault(prov.Logger)
//
//	// 業務ログ (ctx 必須)
//	slog.ErrorContext(ctx, "internal error",
//	    obs.AttrExceptionMessage, err.Error(),
//	)
package obs

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.39.0"
	"go.opentelemetry.io/otel/trace"
)

// Config は Bootstrap のパラメータ。
// ServiceName のみ必須、それ以外は実用的な既定値を持つ。
type Config struct {
	// ServiceName は OTel resource の service.name に直行する (必須)。
	ServiceName string

	// ServiceVersion は optional。空なら付与しない。
	ServiceVersion string

	// Writer は logs 出力先。nil なら os.Stdout。
	Writer io.Writer

	// Level は slog の最低出力レベル。nil なら LevelInfo。
	Level slog.Leveler

	// SensitiveKeys に列挙された attribute key (lowercase 比較) の値は "[REDACTED]" に置換される。
	SensitiveKeys []string

	// AddSource が true のとき、ERROR 以上のログに source (file:line) を付与する。
	AddSource bool

	// EnableTraces が true のとき stdouttrace に span を吐く。false なら NoOp。
	EnableTraces bool

	// EnableMetrics が true のとき stdoutmetric に測定値を吐く。false なら NoOp。
	EnableMetrics bool

	// SyncExport が true のとき tracer を BatchSpanProcessor ではなく SimpleSpanProcessor
	// で構成する (span ごとに即 stdout へ出す)。dev / debug 用途で flush 遅延を消す。
	// 本番では false のまま (= BatchSpanProcessor) にする。
	SyncExport bool
}

// Provider は初期化済みの Logger / Tracer / Meter を保持する。
// Bootstrap が返した shutdown 関数をプロセス終了時に必ず呼ぶこと。
type Provider struct {
	Logger *slog.Logger
	Tracer trace.Tracer
	Meter  metric.Meter

	tracerProvider *sdktrace.TracerProvider
	meterProvider  *sdkmetric.MeterProvider
}

// Bootstrap はログ/トレース/メトリクスを初期化して Provider と shutdown 関数を返す。
// 返り値の shutdown は defer で必ず呼び出すこと (traces/metrics の flush に必要)。
//
// ctx は stdout exporter のみを扱う v1 時点では未使用だが、将来の OTLP exporter 移行で
// 接続確立のキャンセル用に使う予定のためシグネチャに残す。
func Bootstrap(ctx context.Context, cfg Config) (*Provider, func(context.Context) error, error) {
	_ = ctx
	if cfg.ServiceName == "" {
		return nil, nil, errors.New("obs: ServiceName is required")
	}
	if cfg.Writer == nil {
		cfg.Writer = os.Stdout
	}
	// logs / traces / metrics はそれぞれ内部で直列化するが、互いに非協調で
	// 同一 fd に書くとバイト単位で interleave しうる。外側に mutex を被せて
	// 3 者の Write を跨って直列化する。
	cfg.Writer = newLockedWriter(cfg.Writer)
	if cfg.Level == nil {
		cfg.Level = slog.LevelInfo
	}

	res, err := newResource(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("obs: resource: %w", err)
	}

	p := &Provider{Logger: buildLogger(cfg)}

	if cfg.EnableTraces {
		tp, err := newTracerProvider(cfg, res)
		if err != nil {
			return nil, nil, fmt.Errorf("obs: tracer: %w", err)
		}
		otel.SetTracerProvider(tp)
		p.tracerProvider = tp
		p.Tracer = tp.Tracer(cfg.ServiceName)
	} else {
		p.Tracer = otel.Tracer(cfg.ServiceName)
	}

	if cfg.EnableMetrics {
		mp, err := newMeterProvider(cfg, res)
		if err != nil {
			return nil, nil, fmt.Errorf("obs: meter: %w", err)
		}
		otel.SetMeterProvider(mp)
		p.meterProvider = mp
		p.Meter = mp.Meter(cfg.ServiceName)
	} else {
		p.Meter = otel.Meter(cfg.ServiceName)
	}

	return p, p.shutdown, nil
}

func (p *Provider) shutdown(ctx context.Context) error {
	var merr []error
	if p.tracerProvider != nil {
		if err := p.tracerProvider.Shutdown(ctx); err != nil {
			merr = append(merr, fmt.Errorf("tracer shutdown: %w", err))
		}
	}
	if p.meterProvider != nil {
		if err := p.meterProvider.Shutdown(ctx); err != nil {
			merr = append(merr, fmt.Errorf("meter shutdown: %w", err))
		}
	}
	return errors.Join(merr...)
}

// lockedWriter は io.Writer を sync.Mutex で直列化するデコレータ。
// logs / traces / metrics の 3 系統が同一 fd に書く際、内部ロックだけでは
// 系統間で Write が interleave しうるため、外側でまとめて直列化する。
type lockedWriter struct {
	mu sync.Mutex
	w  io.Writer
}

func newLockedWriter(w io.Writer) *lockedWriter {
	return &lockedWriter{w: w}
}

func (l *lockedWriter) Write(p []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.w.Write(p)
}

func newResource(cfg Config) (*resource.Resource, error) {
	attrs := []attribute.KeyValue{semconv.ServiceName(cfg.ServiceName)}
	if cfg.ServiceVersion != "" {
		attrs = append(attrs, semconv.ServiceVersion(cfg.ServiceVersion))
	}
	// Schema URL は空にして resource.Default() 側 (SDK が宣言する sem-conv 版)
	// を採用させる。semconv パッケージのバージョンと SDK の schema URL は
	// ずれることがあり、明示すると resource.Merge が conflicting schema URL で落ちる。
	return resource.Merge(
		resource.Default(),
		resource.NewWithAttributes("", attrs...),
	)
}
