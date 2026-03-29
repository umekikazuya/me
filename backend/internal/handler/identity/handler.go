package identity

import (
	"encoding/json"
	"net/http"

	appidentity "github.com/umekikazuya/me/internal/app/identity"
)

type Handler struct {
	interactor appidentity.TokenService // TODO: interactor interface に差し替え
	tokenSrv   appidentity.TokenService
}

func NewHandler(interactor appidentity.TokenService, tokenSrv appidentity.TokenService) *Handler {
	return &Handler{interactor: interactor, tokenSrv: tokenSrv}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v) //nolint:errcheck
}
