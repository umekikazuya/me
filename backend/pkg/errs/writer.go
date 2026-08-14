package errs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	problemContentType = "application/problem+json"
	problemTypeBlank   = "about:blank"
)

type (
	ErrorSink    interface{ Set(error) }
	errorSinkKey struct{}
)

func WithErrorSink(
	ctx context.Context,
	sink ErrorSink,
) context.Context {
	return context.WithValue(ctx, errorSinkKey{}, sink)
}

// WriteProblem はエラーを HTTP レスポンスへ書き出す。
//
// 返却形式は以下の契約に従う。
//   - ProblemDetail（application/problem+json）
//   - err が nil: 呼び出し側バグとして 500 を返す
func WriteProblem(w http.ResponseWriter, r *http.Request, err error) {
	// nil error は呼び出し側のバグを示すため、500 を返して安全側に倒す
	if err == nil {
		p := ProblemDetail{
			Type:     problemTypeBlank,
			Title:    "Internal Server Error",
			Status:   http.StatusInternalServerError,
			Instance: instanceFromRequest(r),
		}
		w.Header().Set("Content-Type", problemContentType)
		w.WriteHeader(p.Status)
		json.NewEncoder(w).Encode(p) //nolint:errcheck,gosec
		return
	}

	// Context にエラー情報を伝播
	if r != nil {
		if s, ok := r.Context().Value(errorSinkKey{}).(ErrorSink); ok {
			s.Set(err)
		}
	}

	p := toProblem(err, instanceFromRequest(r))
	w.Header().Set("Content-Type", problemContentType)
	w.WriteHeader(p.Status)
	json.NewEncoder(w).Encode(p) //nolint:errcheck,gosec
}

func instanceFromRequest(r *http.Request) string {
	if r == nil || r.URL == nil {
		return ""
	}
	return r.URL.Path
}

func toProblem(err error, instance string) ProblemDetail {
	msg := err.Error()

	// 400: ValidationError は invalidParams に展開
	var ve *ValidationError
	if errors.As(err, &ve) {
		return ProblemDetail{
			Type:          problemTypeBlank,
			Title:         "Bad Request",
			Status:        http.StatusBadRequest,
			Instance:      instance,
			InvalidParams: ve.Params,
		}
	}

	switch {
	case errors.Is(err, ErrBadRequest):
		return ProblemDetail{
			Type:     problemTypeBlank,
			Title:    "Bad Request",
			Status:   http.StatusBadRequest,
			Detail:   badRequestDetail(err),
			Instance: instance,
		}
	case errors.Is(err, ErrNotFound):
		return ProblemDetail{Type: problemTypeBlank, Title: "Not Found", Status: http.StatusNotFound, Detail: msg, Instance: instance}
	case errors.Is(err, ErrConflict):
		return ProblemDetail{Type: problemTypeBlank, Title: "Conflict", Status: http.StatusConflict, Detail: msg, Instance: instance}
	case errors.Is(err, ErrPermissionDenied):
		return ProblemDetail{Type: problemTypeBlank, Title: "Forbidden", Status: http.StatusForbidden, Detail: msg, Instance: instance}
	default:
		// 500: 内部エラーは detail を漏らさない
		return ProblemDetail{Type: problemTypeBlank, Title: "Internal Server Error", Status: http.StatusInternalServerError, Instance: instance}
	}
}

func badRequestDetail(err error) string {
	var syntaxErr *json.SyntaxError
	if errors.As(err, &syntaxErr) {
		return fmt.Sprintf("malformed JSON at position %d", syntaxErr.Offset)
	}
	if errors.Is(err, io.ErrUnexpectedEOF) {
		return "malformed JSON"
	}
	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &typeErr) {
		if typeErr.Field != "" {
			return fmt.Sprintf("invalid value type for field %q", typeErr.Field)
		}
		if typeErr.Offset > 0 {
			return fmt.Sprintf("invalid value type at position %d", typeErr.Offset)
		}
		return "invalid value type"
	}
	if errors.Is(err, io.EOF) {
		return "request body must not be empty"
	}

	msg := err.Error()
	if strings.Contains(msg, "http: request body too large") {
		return "request body must not exceed 1MB"
	}
	if i := strings.Index(msg, "json: unknown field "); i >= 0 {
		field := strings.Trim(strings.TrimPrefix(msg[i:], "json: unknown field "), `"`)
		return fmt.Sprintf("unknown field %q", field)
	}
	const badRequestPrefix = "bad request: "
	if strings.HasPrefix(msg, badRequestPrefix) {
		return strings.TrimPrefix(msg, badRequestPrefix)
	}
	return msg
}
