package identity

import (
	"encoding/json"
	"fmt"
	"net/http"

	app "github.com/umekikazuya/me/internal/app/identity"
	"github.com/umekikazuya/me/pkg/errs"
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
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errs.WriteProblem(w, fmt.Errorf(
			"decode request body: %w",
			errs.ErrBadRequest,
		))
		return
	}
	out, err := h.interactor.Login(
		r.Context(),
		input,
	)
	if err != nil {
		errs.WriteProblem(w, err)
		return
	}
	setTokenCookies(w, out.AT, out.RT)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	identityID, ok := identityIDFromContext(r.Context())
	if !ok {
		errs.WriteProblem(w, errs.ErrUnauthenticated)
	}
	var input app.InputLogoutDto
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errs.WriteProblem(w, fmt.Errorf(
			"decode request body: %w",
			errs.ErrBadRequest,
		))
		return
	}
	input.IdentityID = identityID
	err := h.interactor.Logout(
		r.Context(),
		input,
	)
	if err != nil {
		errs.WriteProblem(w, err)
		return
	}
	clearTokenCookies(w)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) RevokeSessions(w http.ResponseWriter, r *http.Request) {
	var input app.InputRevokeAllSessionsDto
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errs.WriteProblem(w, fmt.Errorf(
			"decode request body: %w",
			errs.ErrBadRequest,
		))
		return
	}
	err := h.interactor.RevokeAllSessions(
		r.Context(),
		input,
	)
	if err != nil {
		errs.WriteProblem(w, err)
		return
	}
	clearTokenCookies(w)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	identityID, ok := identityIDFromContext(r.Context())
	if !ok {
		errs.WriteProblem(w, errs.ErrUnauthenticated)
		return
	}
	rtCookie, err := r.Cookie(refreshTokenCookieName)
	if err != nil {
		errs.WriteProblem(w, errs.ErrUnauthenticated)
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
		errs.WriteProblem(w, err)
		return
	}
	setTokenCookies(w, out.AT, out.RT)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var input app.InputRegisterDto
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errs.WriteProblem(w, fmt.Errorf(
			"decode request body: %w",
			errs.ErrBadRequest,
		))
		return
	}
	err := h.interactor.Register(
		r.Context(),
		input,
	)
	if err != nil {
		errs.WriteProblem(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	identityID, ok := identityIDFromContext(r.Context())
	if !ok {
		errs.WriteProblem(w, errs.ErrUnauthenticated)
		return
	}
	var input app.InputResetPasswordDto
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errs.WriteProblem(w, fmt.Errorf(
			"decode request body: %w",
			errs.ErrBadRequest,
		))
		return
	}
	input.ID = identityID
	err := h.interactor.ResetPassword(
		r.Context(),
		input,
	)
	if err != nil {
		errs.WriteProblem(w, err)
		return
	}
	clearTokenCookies(w)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ChangeEmailAddress(w http.ResponseWriter, r *http.Request) {
	identityID, ok := identityIDFromContext(r.Context())
	if !ok {
		errs.WriteProblem(w, errs.ErrUnauthenticated)
		return
	}
	var input app.InputChangeEmailDto
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errs.WriteProblem(w, fmt.Errorf(
			"decode request body: %w",
			errs.ErrBadRequest,
		))
		return
	}
	input.ID = identityID
	err := h.interactor.ChangeEmail(
		r.Context(),
		input,
	)
	if err != nil {
		errs.WriteProblem(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
