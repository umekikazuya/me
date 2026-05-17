package health

import (
	"net/http"

	"github.com/umekikazuya/me/pkg/httpx"
)

type Handler struct{}

type Res struct {
	Status string `json:"status"`
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Exec(
	w http.ResponseWriter,
	r *http.Request,
) {
	httpx.WriteJSON(
		w,
		http.StatusOK,
		Res{Status: "ok"},
	)
}
