package errs

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

func TestToProblem_BadRequestDecodeErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		err    error
		detail string
	}{
		{
			name:   "malformed json with offset",
			err:    fmt.Errorf("%w: %w", ErrBadRequest, &json.SyntaxError{Offset: 18}),
			detail: "malformed JSON at position 18",
		},
		{
			name:   "unknown field",
			err:    fmt.Errorf("%w: %w", ErrBadRequest, errors.New(`json: unknown field "nickname"`)),
			detail: `unknown field "nickname"`,
		},
		{
			name:   "explicit bad request detail",
			err:    New(ErrBadRequest, "Content-Type must be application/json"),
			detail: "Content-Type must be application/json",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := toProblem(tt.err, "/recipient")
			if got.Status != 400 {
				t.Fatalf("status = %d, want 400", got.Status)
			}
			if got.Title != "Bad Request" {
				t.Fatalf("title = %q, want %q", got.Title, "Bad Request")
			}
			if got.Detail != tt.detail {
				t.Fatalf("detail = %q, want %q", got.Detail, tt.detail)
			}
		})
	}
}
