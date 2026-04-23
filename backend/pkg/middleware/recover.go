package middleware

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/umekikazuya/me/pkg/errs"
	"github.com/umekikazuya/me/pkg/obs"
)

// RecoverOption は Recover の振る舞いを差し替える関数型オプション。
type RecoverOption func(*recoverConfig)

type recoverConfig struct {
	logger *slog.Logger
}

// WithLogger はテスト等で slog.Default() 以外の Logger を注入する。
func WithLogger(l *slog.Logger) RecoverOption {
	return func(c *recoverConfig) { c.logger = l }
}

// Recover は panic を捕捉し、500 ProblemDetails を返しつつ ERROR ログを残す。
// アクセスログではなく「プロセス保護 + 未処理エラーの可視化」のための middleware。
//
// 属性名は OpenTelemetry Semantic Conventions に準拠する (pkg/obs/attr.go)。
// Logger は opts で明示注入できる。未指定時は slog.Default()。
func Recover(next http.Handler, opts ...RecoverOption) http.Handler {
	cfg := recoverConfig{}
	for _, o := range opts {
		o(&cfg)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			rec := recover()
			if rec == nil {
				return
			}
			err, ok := rec.(error)
			if !ok {
				err = fmt.Errorf("panic: %v", rec)
			}
			logger := cfg.logger
			if logger == nil {
				logger = slog.Default()
			}
			logger.ErrorContext(r.Context(), "unhandled panic",
				obs.AttrExceptionMessage, err.Error(),
				obs.AttrExceptionType, fmt.Sprintf("%T", err),
				obs.AttrExceptionStack, string(debug.Stack()),
				obs.AttrHTTPMethod, r.Method,
				obs.AttrURLPath, r.URL.Path,
			)
			errs.WriteProblem(w, r, errors.Join(errs.ErrInternal, err))
		}()
		next.ServeHTTP(w, r)
	})
}
