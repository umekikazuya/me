package obs

import (
	"fmt"

	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

func newMeterProvider(cfg Config, res *resource.Resource) (*sdkmetric.MeterProvider, error) {
	exp, err := stdoutmetric.New(stdoutmetric.WithWriter(cfg.Writer))
	if err != nil {
		return nil, fmt.Errorf("stdoutmetric: %w", err)
	}
	return sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exp)),
		sdkmetric.WithResource(res),
	), nil
}
