package identity

import (
	"encoding/json"
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

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v) //nolint:errcheck
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var input app.InputLoginDto
	// TODO: バリデーション
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
	var input app.InputLogoutDto
	// TODO: バリデーション
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
	// TODO: バリデーション
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
	var input app.InputRefreshTokensDto
	// TODO: バリデーション
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
	// TODO: バリデーション
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
	var input app.InputResetPasswordDto
	// TODO: バリデーション
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
	var input app.InputChangeEmailDto
	// TODO: バリデーション
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
