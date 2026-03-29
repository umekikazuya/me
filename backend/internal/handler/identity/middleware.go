package identity

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type contextKey string

const identityIDKey contextKey = "identityID"

const XRequestedWith = "X-Requested-With"

// CSRFMiddleware は X-Requested-With ヘッダーで簡易 CSRF 対策を行う
func CSRFMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			if r.Header.Get(XRequestedWith) != "XMLHttpRequest" {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		},
	)
}

// AuthMiddleware は meAccessToken Cookie を検証し、IdentityID をコンテキストに注入する
// POST /identity/logout、GET /identity/refresh 以外の保護エンドポイントで使用する
func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("meAccessToken")
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		identityID, err := h.tokenSrv.ValidateAT(r.Context(), cookie.Value)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), identityIDKey, identityID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RefreshMiddleware は meAccessToken を ParseUnverified して IdentityID をコンテキストに注入する
// POST /identity/refresh 専用。AT の署名検証は行わない（期限切れ AT からの更新を許容するため）
// TODO: SessionRepo に GSI(tokenHash) を追加し FindByTokenHash を実装したら本 middleware を削除して
//
//	RefreshTokens の input から IdentityID を除去する
func (h *Handler) RefreshMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("meAccessToken")
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		claims, _, err := jwt.NewParser().ParseUnverified(cookie.Value, &jwt.RegisteredClaims{})
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		id, err := claims.Claims.GetSubject()
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		_, err = uuid.Parse(id)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), identityIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// identityIDFromContext はコンテキストから IdentityID を取り出す
func identityIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(identityIDKey).(string)
	return id, ok
}
