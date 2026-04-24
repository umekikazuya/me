package obs

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/umekikazuya/me/pkg/errs"
)

// LogIfInternal は err が errs.ErrInternal を含むときのみ ERROR を 1 件出す。
// 通常は handler 境界 (WriteProblem の直前) で使う。
// client fault (ErrBadRequest / ErrNotFound 等) のときは何もしない。
func LogIfInternal(ctx context.Context, err error) {
	if err == nil || !errors.Is(err, errs.ErrInternal) {
		return
	}
	slog.ErrorContext(ctx, "internal error",
		AttrExceptionMessage, err.Error(),
		AttrExceptionType, fmt.Sprintf("%T", err),
	)
}
