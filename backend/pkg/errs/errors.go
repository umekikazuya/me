package errs

import "errors"

var (
	// ErrBadRequest は 400 Bad Request 相当の入力不備を表す。
	ErrBadRequest = errors.New("bad request")
	// ErrNotFound は 404 Not Found 相当の未存在を表す。
	ErrNotFound = errors.New("not found")
	// ErrConflict は 409 Conflict 相当の競合を表す。
	ErrConflict = errors.New("conflict")
	// ErrPermissionDenied は 403 Forbidden 相当の認可拒否を表す。
	ErrPermissionDenied = errors.New("permission denied") // 権限がない
	// ErrUnauthenticated は 401 Unauthorized 相当の未認証を表す。
	ErrUnauthenticated = errors.New("unauthenticated") // 認証されていない
	// ErrInternal は 500 Internal Server Error 相当の内部失敗を表す。
	ErrInternal = errors.New("internal")
)

// ValidationError は HTTP レベルの入力バリデーションエラーを表す。
// 400 Bad Request として ProblemDetails.invalidParams に展開される。
type ValidationError struct {
	// Params は不正パラメータの一覧。各要素は invalidParams に展開される。
	Params []InvalidParam
}

// Error は error インターフェース実装。
func (e *ValidationError) Error() string { return "validation failed" }

// Is は ValidationError を ErrBadRequest と同一カテゴリとして扱う。
func (e *ValidationError) Is(target error) bool { return target == ErrBadRequest }
