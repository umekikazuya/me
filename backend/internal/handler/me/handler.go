package me

import (
	"fmt"
	"net/http"
	"os"

	app "github.com/umekikazuya/me/internal/app/me"
	"github.com/umekikazuya/me/pkg/errs"
	"github.com/umekikazuya/me/pkg/httpx"
	"github.com/umekikazuya/me/pkg/obs"
)

type Handler struct {
	me app.Interactor
}

func NewHandler(me app.Interactor) *Handler {
	return &Handler{me: me}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	meID := os.Getenv("ME_ID")
	if meID == "" {
		errs.WriteProblem(w, r, fmt.Errorf("ME_ID environment variable is not configured: %w", errs.ErrNotFound))
		return
	}
	out, err := h.me.Get(r.Context(), meID)
	if err != nil {
		obs.LogIfInternal(r.Context(), err)
		errs.WriteProblem(w, r, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, out)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	meID := os.Getenv("ME_ID")
	if meID == "" {
		errs.WriteProblem(w, r, fmt.Errorf("ME_ID environment variable is not configured: %w", errs.ErrNotFound))
		return
	}
	var input app.InputDto
	if err := httpx.DecodeAndValidate(w, r, &input); err != nil {
		errs.WriteProblem(w, r, err)
		return
	}
	input.ID = meID
	out, err := h.me.Update(r.Context(), input)
	if err != nil {
		obs.LogIfInternal(r.Context(), err)
		errs.WriteProblem(w, r, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, out)
}
