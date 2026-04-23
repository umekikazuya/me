package errs

import "errors"

var (
	ErrBadRequest       = errors.New("bad request")
	ErrNotFound         = errors.New("not found")
	ErrConflict         = errors.New("conflict")
	ErrUnprocessable    = errors.New("unprocessable entity")
	ErrPermissionDenied = errors.New("permission denied") // 権限がない
	ErrUnauthenticated  = errors.New("unauthenticated")   // 認証されていない
	ErrInternal         = errors.New("internal")
)

// ValidationError は HTTP レベルの入力バリデーションエラーを表す。
// 400 Bad Request として ProblemDetails.invalidParams に展開される。
type ValidationError struct {
	Params []InvalidParam
}

func (e *ValidationError) Error() string { return "validation failed" }

func (e *ValidationError) Is(target error) bool { return target == ErrBadRequest }

// DomainError はドメイン不変条件違反を表す。
// 422 Unprocessable Entity として DomainProblem.details に展開される。
type DomainError struct {
	Code    string
	Message string
	Details []DomainProblemItem
}

func (e *DomainError) Error() string { return e.Message }

func (e *DomainError) Is(target error) bool { return target == ErrUnprocessable }
