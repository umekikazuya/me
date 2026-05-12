package me

import (
	"errors"
	"net/http"
	"strings"

	app "github.com/umekikazuya/me/internal/app/me"
	"github.com/umekikazuya/me/pkg/errs"
	"github.com/umekikazuya/me/pkg/httpx"
	"github.com/umekikazuya/me/pkg/obs"
)

type Handler struct {
	me   app.Interactor
	meID string
}

func NewHandler(me app.Interactor, meID string) (*Handler, error) {
	meID = strings.TrimSpace(meID)
	if meID == "" {
		return nil, errors.New("ME_ID is not set")
	}
	return &Handler{me: me, meID: meID}, nil
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	out, err := h.me.Get(r.Context(), h.meID)
	if err != nil {
		obs.LogIfInternal(r.Context(), err)
		errs.WriteProblem(w, r, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, out)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var input app.InputDto
	if err := httpx.DecodeAndValidate(w, r, &input); err != nil {
		errs.WriteProblem(w, r, err)
		return
	}
	input.ID = h.meID
	out, err := h.me.Update(r.Context(), input)
	if err != nil {
		obs.LogIfInternal(r.Context(), err)
		errs.WriteProblem(w, r, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, out)
}
