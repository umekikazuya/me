package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/umekikazuya/me/pkg/obs"
)

const requestIDHeader = "X-Request-ID"

// RequestID は X-Request-ID をクライアントから受け取るか、無ければ UUIDv4 を生成し、
// context とレスポンスヘッダに載せる。
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(requestIDHeader)
		if _, err := uuid.Parse(id); err != nil {
			id = uuid.NewString()
		}
		w.Header().Set(requestIDHeader, id)
		next.ServeHTTP(w, r.WithContext(obs.WithRequestID(r.Context(), id)))
	})
}
