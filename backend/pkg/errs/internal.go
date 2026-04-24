package errs

import "fmt"

// WrapInternal は infra 由来のエラーを 500 として扱える ErrInternal でラップして返す。
// 元のエラーも連鎖として保持するため errors.Is で元エラーにマッチする (デバッグ・テスト用)。
//
// 本関数はログを出さない — ログは handler 境界で `obs.LogIfInternal(ctx, err)` が 1 回だけ出す。
// ログ/エラーの責務を分けることで、app 層 (interactor) を logger から完全に切り離す。
func WrapInternal(op string, err error) error {
	return fmt.Errorf("%s: %w: %w", op, ErrInternal, err)
}
