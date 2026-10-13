package identity

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	app "github.com/umekikazuya/me/internal/app/identity"
	"github.com/umekikazuya/me/pkg/errs"
	"github.com/umekikazuya/me/pkg/obs"
)

type Handler struct {
	interactor app.Interactor
	tokenSrv   app.TokenService
}

const cookieName = "authState"

func NewHandler(interactor app.Interactor, tokenSrv app.TokenService) *Handler {
	return &Handler{interactor: interactor, tokenSrv: tokenSrv}
}

func getState() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("システムエラー: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (h *Handler) LoginFromGithub(
	w http.ResponseWriter,
	r *http.Request,
) {
	state, err := getState()
	if err != nil {
		errs.WriteProblem(w, r, err)
		return
	}
	http.SetCookie(
		w,
		&http.Cookie{Name: cookieName, Value: state, Quoted: false, Expires: time.Now().Add(20 * time.Minute), HttpOnly: true, Secure: true, SameSite: http.SameSiteLaxMode},
	)
	http.Redirect(
		w,
		r,
		h.interactor.LoginWithGithub(r.Context(), app.InputLoginWithGithubDto{State: state}),
		http.StatusFound,
	)
}

func (h *Handler) CallbackGithub(
	w http.ResponseWriter,
	r *http.Request,
) {
	c, err := r.Cookie(cookieName)
	if err != nil {
		errs.WriteProblem(w, r, errs.New(errs.ErrBadRequest, err.Error()))
		return
	}
	if c.Value != r.URL.Query().Get("state") {
		errs.WriteProblem(w, r, errs.New(errs.ErrBadRequest, "state が不正"))
		return
	}
	rawCode := r.URL.Query().Get("code")
	if rawCode == "" {
		errs.WriteProblem(w, r, errs.New(errs.ErrBadRequest, "code が不正"))
		return
	}
	out, err := h.interactor.CallbackFromGithub(
		r.Context(),
		app.InputCallbackFromGithubDto{
			Code:     rawCode,
			BaseTime: time.Now(),
		},
	)
	if err != nil {
		errs.WriteProblem(w, r, err)
		return
	}
	setTokenCookies(w, out.AT, out.RT)
	http.SetCookie(
		w,
		&http.Cookie{Name: cookieName, Value: "", Quoted: false, Expires: time.Now().Add(20 * time.Minute), HttpOnly: true, MaxAge: -1, Secure: true, SameSite: http.SameSiteLaxMode},
	)
	http.Redirect(w, r, "https://www.umekikazuya.me", http.StatusFound) // TODO: 環境変数使う
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	identityID, ok := identityIDFromContext(r.Context())
	if !ok {
		errs.WriteProblem(w, r, errs.ErrUnauthenticated)
		return
	}
	rtCookie, err := r.Cookie(refreshTokenCookieName)
	if err != nil {
		errs.WriteProblem(w, r, errs.ErrUnauthenticated)
		return
	}
	err = h.interactor.Logout(
		r.Context(),
		app.InputLogoutDto{
			IdentityID: identityID,
			RT:         rtCookie.Value,
		},
	)
	if err != nil {
		obs.LogIfInternal(r.Context(), err)
		errs.WriteProblem(w, r, err)
		return
	}
	clearTokenCookies(w)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) RevokeSessions(w http.ResponseWriter, r *http.Request) {
	identityID, ok := identityIDFromContext(r.Context())
	if !ok {
		errs.WriteProblem(w, r, errs.ErrUnauthenticated)
		return
	}
	var input app.InputRevokeAllSessionsDto
	input.IdentityID = identityID
	err := h.interactor.RevokeAllSessions(
		r.Context(),
		input,
	)
	if err != nil {
		obs.LogIfInternal(r.Context(), err)
		errs.WriteProblem(w, r, err)
		return
	}
	clearTokenCookies(w)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	identityID, ok := identityIDFromContext(r.Context())
	if !ok {
		errs.WriteProblem(w, r, errs.ErrUnauthenticated)
		return
	}
	rtCookie, err := r.Cookie(refreshTokenCookieName)
	if err != nil {
		errs.WriteProblem(w, r, errs.ErrUnauthenticated)
		return
	}
	input := app.InputRefreshTokensDto{
		IdentityID: identityID,
		RT:         rtCookie.Value,
	}
	out, err := h.interactor.RefreshTokens(
		r.Context(),
		input,
	)
	if err != nil {
		obs.LogIfInternal(r.Context(), err)
		errs.WriteProblem(w, r, err)
		return
	}
	setTokenCookies(w, out.AT, out.RT)
	w.WriteHeader(http.StatusNoContent)
}
