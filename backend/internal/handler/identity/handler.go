package identity

import (
	"net/http"

	app "github.com/umekikazuya/me/internal/app/identity"
	"github.com/umekikazuya/me/pkg/errs"
	"github.com/umekikazuya/me/pkg/httpx"
	"github.com/umekikazuya/me/pkg/obs"
)

type Handler struct {
	interactor app.Interactor
	tokenSrv   app.TokenService
}

func NewHandler(interactor app.Interactor, tokenSrv app.TokenService) *Handler {
	return &Handler{interactor: interactor, tokenSrv: tokenSrv}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var input app.InputLoginDto
	err := httpx.DecodeAndValidate(w, r, &input)
	if err != nil {
		errs.WriteProblem(w, r, err)
		return
	}
	out, err := h.interactor.Login(
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

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var input app.InputRegisterDto
	err := httpx.DecodeAndValidate(w, r, &input)
	if err != nil {
		errs.WriteProblem(w, r, err)
		return
	}
	err = h.interactor.Register(
		r.Context(),
		input,
	)
	if err != nil {
		obs.LogIfInternal(r.Context(), err)
		errs.WriteProblem(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	identityID, ok := identityIDFromContext(r.Context())
	if !ok {
		errs.WriteProblem(w, r, errs.ErrUnauthenticated)
		return
	}
	var input app.InputResetPasswordDto
	err := httpx.DecodeAndValidate(w, r, &input)
	if err != nil {
		errs.WriteProblem(w, r, err)
		return
	}
	input.ID = identityID
	err = h.interactor.ResetPassword(
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

func (h *Handler) ChangeEmailAddress(w http.ResponseWriter, r *http.Request) {
	identityID, ok := identityIDFromContext(r.Context())
	if !ok {
		errs.WriteProblem(w, r, errs.ErrUnauthenticated)
		return
	}
	var input app.InputChangeEmailDto
	err := httpx.DecodeAndValidate(w, r, &input)
	if err != nil {
		errs.WriteProblem(w, r, err)
		return
	}
	input.ID = identityID
	err = h.interactor.ChangeEmail(
		r.Context(),
		input,
	)
	if err != nil {
		obs.LogIfInternal(r.Context(), err)
		errs.WriteProblem(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
