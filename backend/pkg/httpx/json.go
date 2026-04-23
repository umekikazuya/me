package httpx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/umekikazuya/me/pkg/errs"
)

const maxJSONBodyBytes int64 = 1 << 20

var validate = validator.New()

// WriteJSON は成功レスポンスを application/json で書き出す。
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v == nil {
		return
	}
	json.NewEncoder(w).Encode(v) //nolint:errcheck
}

// DecodeAndValidate はリクエストボディを dst にデコードし、validator タグで検証する。
// サイズ上限 1MB。未知フィールドは拒否。trailing content も拒否。
// バリデーション失敗時は errs.ValidationError を返す (400 として扱われる)。
func DecodeAndValidate(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxJSONBodyBytes)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return fmt.Errorf("decode request body: %w", errs.ErrBadRequest)
	}
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return fmt.Errorf("decode request body: %w", errs.ErrBadRequest)
	}
	if err := validate.Struct(dst); err != nil {
		return &errs.ValidationError{Params: toInvalidParams(err)}
	}
	return nil
}

func toInvalidParams(err error) []errs.InvalidParam {
	var ves validator.ValidationErrors
	if !errors.As(err, &ves) {
		return nil
	}
	out := make([]errs.InvalidParam, 0, len(ves))
	for _, fe := range ves {
		out = append(out, errs.InvalidParam{
			Name:   fe.Field(),
			Reason: fe.Tag(),
		})
	}
	return out
}
