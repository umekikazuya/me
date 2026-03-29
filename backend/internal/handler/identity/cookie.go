package identity

import "net/http"

const (
	accessTokenCookieName  = "meAccessToken"
	refreshTokenCookieName = "meRefreshToken"

	atMaxAge = 15 * 60           // 15 minutes
	rtMaxAge = 30 * 24 * 60 * 60 // 30 days
)

// setTokenCookies は AT と RT を HttpOnly Cookie としてレスポンスにセットする
func setTokenCookies(w http.ResponseWriter, at, rt string) {
	// TODO: Set-Cookie meAccessToken
	//   HttpOnly=true, Secure=true, SameSite=Strict, Path="/", MaxAge=atMaxAge

	// TODO: Set-Cookie meRefreshToken
	//   HttpOnly=true, Secure=true, SameSite=Strict, Path="/identity/refresh", MaxAge=rtMaxAge
	//   Path をリフレッシュエンドポイントに限定することで RT の露出範囲を最小化
}

// clearTokenCookies は AT と RT Cookie を削除する（ログアウト・全セッション解除用）
func clearTokenCookies(w http.ResponseWriter) {
	// TODO: meAccessToken を MaxAge=-1 で上書き（削除）
	// TODO: meRefreshToken を MaxAge=-1 で上書き（削除）
}
