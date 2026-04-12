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
