package identity

import (
	"context"
	"net/http"
)

type contextKey string

const identityIDKey contextKey = "identityID"

// CSRFMiddleware は X-Requested-With ヘッダーで簡易 CSRF 対策を行う
func CSRFMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: X-Requested-With: XMLHttpRequest の検証
		// if r.Header.Get("X-Requested-With") != "XMLHttpRequest" {
		//     http.Error(w, "forbidden", http.StatusForbidden)
		//     return
		// }
		next.ServeHTTP(w, r)
	})
}

// AuthMiddleware は meAccessToken Cookie を検証し、IdentityID をコンテキストに注入する
// POST /identity/logout、GET /identity/refresh 以外の保護エンドポイントで使用する
func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO:
		// 1. r.Cookie("meAccessToken") でトークン取得
		// 2. h.tokenSrv.Validate(ctx, token) で検証
		// 3. gojwt.ParseUnverified で sub (IdentityID) を取得
		// 4. context.WithValue(r.Context(), identityIDKey, sub) で注入
		// 5. next.ServeHTTP(w, r.WithContext(ctx))
		next.ServeHTTP(w, r)
	})
}

// RefreshMiddleware は meAccessToken を ParseUnverified して IdentityID をコンテキストに注入する
// POST /identity/refresh 専用。AT の署名検証は行わない（期限切れ AT からの更新を許容するため）
// TODO: SessionRepo に GSI(tokenHash) を追加し FindByTokenHash を実装したら本 middleware を削除して
//
//	RefreshTokens の input から IdentityID を除去する
func (h *Handler) RefreshMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO:
		// 1. r.Cookie("meAccessToken") でトークン取得（期限切れでも可）
		// 2. gojwt.NewParser().ParseUnverified(token, &jwt.RegisteredClaims{}) で sub 取得
		// 3. context.WithValue で identityIDKey に注入
		// 4. next.ServeHTTP(w, r.WithContext(ctx))
		next.ServeHTTP(w, r)
	})
}

// identityIDFromContext はコンテキストから IdentityID を取り出す
func identityIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(identityIDKey).(string)
	return id, ok
}
