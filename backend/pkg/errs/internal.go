package errs

import (
	"context"
	"fmt"
	"log/slog"
)

// WrapInternal は infra 由来のエラーをログに出力し、500 として扱える ErrInternal でラップして返す。
// 元のエラーも連鎖として保持するため errors.Is で元エラーにマッチする (デバッグ・テスト用)。
// クライアントへ漏らしたくない内部詳細はログにのみ残り、レスポンスは空の ProblemDetail になる。
func WrapInternal(ctx context.Context, op string, err error) error {
	slog.ErrorContext(ctx, op, "error", err)
	return fmt.Errorf("%s: %w: %w", op, ErrInternal, err)
}
