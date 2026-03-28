package errs

import (
	"encoding/json"
	"errors"
	"net/http"
)

type ProblemDetail struct {
	Status int     `json:"status"`
	Title  string  `json:"title"`
	Detail *string `json:"detail,omitempty"`
}

func WriteProblem(w http.ResponseWriter, err error) {
	p := ToProblem(err)
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(p.Status)
	json.NewEncoder(w).Encode(p) //nolint:errcheck
}

func ToProblem(err error) ProblemDetail {
	msg := err.Error()
	detail := &msg
	switch {
	case errors.Is(err, ErrBadRequest):
		return ProblemDetail{Status: http.StatusBadRequest, Title: "Bad Request", Detail: detail}
	case errors.Is(err, ErrNotFound):
		return ProblemDetail{Status: http.StatusNotFound, Title: "Not Found", Detail: detail}
	case errors.Is(err, ErrConflict):
		return ProblemDetail{Status: http.StatusConflict, Title: "Conflict", Detail: detail}
	case errors.Is(err, ErrUnprocessable):
		return ProblemDetail{Status: http.StatusUnprocessableEntity, Title: "Unprocessable Entity", Detail: detail}
	case errors.Is(err, ErrUnauthenticated):
		return ProblemDetail{Status: http.StatusUnauthorized, Title: "Unauthorized", Detail: detail}
	case errors.Is(err, ErrPermissionDenied):
		return ProblemDetail{Status: http.StatusForbidden, Title: "Forbidden", Detail: detail}
	default:
		return ProblemDetail{Status: http.StatusInternalServerError, Title: "Internal Server Error"}
	}
}
