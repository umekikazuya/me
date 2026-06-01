package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/umekikazuya/me/cmd/api/di"
	"github.com/umekikazuya/me/pkg/middleware"
	"github.com/umekikazuya/me/pkg/obs"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const shutdownTimeout = 30 * time.Second

func main() {
	ctx := context.Background()

	prov, shutdown, err := obs.Bootstrap(ctx, obs.Config{
		ServiceName:   "api",
		Level:         obs.ParseLevel(os.Getenv("LOG_LEVEL")),
		SensitiveKeys: []string{"password", "password_hash", "authorization", "cookie", "set-cookie", "token", "refresh_token"},
		AddSource:     true,
		EnableTraces:  true,
		EnableMetrics: true,
	})
	if err != nil {
		slog.Error("観測性基盤の初期化に失敗しました", "error", err)
		os.Exit(1)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if shutdownErr := shutdown(shutdownCtx); shutdownErr != nil {
			slog.Error("観測性基盤のシャットダウンに失敗しました", "error", shutdownErr)
		}
	}()
	slog.SetDefault(prov.Logger)

	handlers, err := di.NewHandlers(ctx)
	if err != nil {
		slog.Error("ハンドラーの初期化に失敗しました", "error", err)
		os.Exit(1)
	}

	r := di.NewRouter(*handlers)
	handler := middleware.RequestID(
		otelhttp.NewHandler(
			middleware.Recover(r),
			"api",
		),
	)
	adapter := httpadapter.NewV2(handler)

	lambda.Start(func(lambdaCtx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		res, proxyErr := adapter.ProxyWithContext(lambdaCtx, req)
		if proxyErr != nil {
			slog.ErrorContext(lambdaCtx, "Lambda ハンドラーでエラーが発生しました", "error", proxyErr)
			return events.APIGatewayV2HTTPResponse{
				StatusCode: 500,
				Body:       `{"title":"Internal Server Error","status":500}`,
				Headers: map[string]string{
					"Content-Type": "application/problem+json",
				},
			}, nil
		}

		if res.StatusCode == 0 {
			return events.APIGatewayV2HTTPResponse{}, errors.New("empty response from lambda http adapter")
		}

		return res, nil
	})
}
