package identity

import (
	"net/http"

	app "github.com/umekikazuya/me/internal/app/identity"
	"github.com/umekikazuya/me/pkg/errs"
	"github.com/umekikazuya/me/pkg/obs"
)

type Handler struct {
	interactor app.Interactor
	tokenSrv   app.TokenService
}

func NewHandler(interactor app.Interactor, tokenSrv app.TokenService) *Handler {
	return &Handler{interactor: interactor, tokenSrv: tokenSrv}
}

func (h *Handler) LoginFromGithub(
	w http.ResponseWriter,
	r *http.Request,
) {
}

func (h *Handler) CallbackGithub(
	w http.ResponseWriter,
	r *http.Request,
) {
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
