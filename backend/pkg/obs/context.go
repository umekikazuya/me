package obs

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
