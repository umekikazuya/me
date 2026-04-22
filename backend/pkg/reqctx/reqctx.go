// Package reqctx はリクエストスコープで context に載せる値の key / accessor を提供する。
// middleware 層と logger 層の双方から依存されるため、依存方向を片方向に保つ目的でここに集約する。
package reqctx

import "context"

type requestIDCtxKey struct{}

// WithRequestID は context に RequestID を載せて返す。
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDCtxKey{}, id)
}

// RequestIDFromContext は context から RequestID を取り出す。未設定の場合は空文字。
func RequestIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	id, _ := ctx.Value(requestIDCtxKey{}).(string)
	return id
}
