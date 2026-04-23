package middleware

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/umekikazuya/me/pkg/errs"
)

// Recover は panic を捕捉し、500 ProblemDetails を返しつつ ERROR ログを残す。
// アクセスログではなく「プロセス保護 + 未処理エラーの可視化」のための middleware。
func Recover(next http.Handler) http.Handler {
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
			slog.ErrorContext(r.Context(), "unhandled panic",
				"error", err,
				"stack", string(debug.Stack()),
				"method", r.Method,
				"path", r.URL.Path,
			)
			errs.WriteProblem(w, r, errors.Join(errs.ErrInternal, err))
		}()
		next.ServeHTTP(w, r)
	})
}
