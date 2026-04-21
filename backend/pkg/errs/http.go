package errs

import (
	"encoding/json"
	"errors"
	"net/http"
)

type errorRecorder interface {
	RecordError(error)
}

// ProblemDetail は RFC 9457 (Problem Details for HTTP APIs) 準拠のエラーレスポンス。
type ProblemDetail struct {
	Type          string         `json:"type"`
	Title         string         `json:"title"`
	Status        int            `json:"status"`
	Detail        string         `json:"detail,omitempty"`
	Instance      string         `json:"instance,omitempty"`
	InvalidParams []InvalidParam `json:"invalidParams,omitempty"`
}

type InvalidParam struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

// DomainProblem は 422 Unprocessable Entity 専用のドメイン違反レスポンス。
type DomainProblem struct {
	Code    string              `json:"code"`
	Message string              `json:"message"`
	Details []DomainProblemItem `json:"details"`
}

type DomainProblemItem struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

const problemContentType = "application/problem+json"

// WriteProblem はエラーを RFC 9457 ProblemDetails (または 422 用 DomainProblem) として書き出す。
func WriteProblem(w http.ResponseWriter, r *http.Request, err error) {
	if recorder, ok := w.(errorRecorder); ok {
		recorder.RecordError(err)
	}

	// 422: ドメインエラーは独自 shape
	if errors.Is(err, ErrUnprocessable) {
		dp := toDomainProblem(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(dp) //nolint:errcheck
		return
	}

	p := toProblem(err, instanceFromRequest(r))
	w.Header().Set("Content-Type", problemContentType)
	w.WriteHeader(p.Status)
	json.NewEncoder(w).Encode(p) //nolint:errcheck
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
			Type:          "about:blank",
			Title:         "Bad Request",
			Status:        http.StatusBadRequest,
			Instance:      instance,
			InvalidParams: ve.Params,
		}
	}

	switch {
	case errors.Is(err, ErrBadRequest):
		return ProblemDetail{Type: "about:blank", Title: "Bad Request", Status: http.StatusBadRequest, Detail: msg, Instance: instance}
	case errors.Is(err, ErrNotFound):
		return ProblemDetail{Type: "about:blank", Title: "Not Found", Status: http.StatusNotFound, Detail: msg, Instance: instance}
	case errors.Is(err, ErrConflict):
		return ProblemDetail{Type: "about:blank", Title: "Conflict", Status: http.StatusConflict, Detail: msg, Instance: instance}
	case errors.Is(err, ErrUnauthenticated):
		return ProblemDetail{Type: "about:blank", Title: "Unauthorized", Status: http.StatusUnauthorized, Detail: msg, Instance: instance}
	case errors.Is(err, ErrPermissionDenied):
		return ProblemDetail{Type: "about:blank", Title: "Forbidden", Status: http.StatusForbidden, Detail: msg, Instance: instance}
	default:
		// 500: 内部エラーは detail を漏らさない
		return ProblemDetail{Type: "about:blank", Title: "Internal Server Error", Status: http.StatusInternalServerError, Instance: instance}
	}
}

func toDomainProblem(err error) DomainProblem {
	var de *DomainError
	if errors.As(err, &de) {
		details := de.Details
		if details == nil {
			details = []DomainProblemItem{}
		}
		code := de.Code
		if code == "" {
			code = "UNPROCESSABLE_ENTITY"
		}
		message := de.Message
		if message == "" {
			message = "Invariant violation"
		}
		return DomainProblem{Code: code, Message: message, Details: details}
	}
	return DomainProblem{
		Code:    "UNPROCESSABLE_ENTITY",
		Message: err.Error(),
		Details: []DomainProblemItem{},
	}
}
