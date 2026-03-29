package identity

import (
	"net/http"
)

const (
	accessTokenCookieName  = "meAccessToken"
	refreshTokenCookieName = "meRefreshToken"

	atMaxAge = 15 * 60           // 15 minutes
	rtMaxAge = 30 * 24 * 60 * 60 // 30 days
)

// setTokenCookies は AT と RT を HttpOnly Cookie としてレスポンスにセットする
func setTokenCookies(w http.ResponseWriter, at, rt string) {
	http.SetCookie(w, &http.Cookie{
		Name:     accessTokenCookieName,
		Value:    at,
		MaxAge:   atMaxAge,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     refreshTokenCookieName,
		Value:    rt,
		MaxAge:   rtMaxAge,
		Path:     "/identity/refresh",
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		HttpOnly: true,
	})
}

// clearTokenCookies は AT と RT Cookie を削除する（ログアウト・全セッション解除用）
func clearTokenCookies(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   accessTokenCookieName,
		MaxAge: -1,
	})
	http.SetCookie(w, &http.Cookie{
		Name:   refreshTokenCookieName,
		MaxAge: -1,
	})
}
