package obs

import (
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"
)

// RecoverProcess は HTTP 外 (batch 等) の panic を ERROR ログに落とし、re-panic しない。
// usage: defer obs.RecoverProcess(ctx, "batch.main")
func RecoverProcess(ctx context.Context, op string) {
	rec := recover()
	if rec == nil {
		return
	}
	err, ok := rec.(error)
	if !ok {
		err = fmt.Errorf("panic: %v", rec)
	}
	slog.ErrorContext(ctx, "unhandled panic",
		AttrOp, op,
		AttrExceptionMessage, err.Error(),
		AttrExceptionStack, string(debug.Stack()),
	)
}
