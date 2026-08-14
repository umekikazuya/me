package server

import (
	"net/http"

	"github.com/umekikazuya/me/pkg/middleware"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func NewHandler(router http.Handler) http.Handler {
	return middleware.RequestID(
		otelhttp.NewHandler(
			middleware.Recover(router),
			"api",
		),
	)
}
