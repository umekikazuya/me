package errs

import "fmt"

func New(errType error, detail string) error {
	return fmt.Errorf("%w: %s", errType, detail)
}

// WrapInternal は infra 由来のエラーを 500 として扱える ErrInternal でラップして返す。
// 元のエラーも連鎖として保持するため errors.Is で元エラーにマッチする (デバッグ・テスト用)。
//
// ログ/エラーの責務を分けることで、app 層 (interactor) を logger から完全に切り離す。
// err は非 nil を想定。
func WrapInternal(op string, err error) error {
	return fmt.Errorf("%s: %w: %w", op, ErrInternal, err)
}
