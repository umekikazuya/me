package me

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	app "github.com/umekikazuya/me/internal/app/me"
	"github.com/umekikazuya/me/pkg/errs"
)

var validate = validator.New()

type Handler struct {
	me app.Interactor
}

func NewHandler(me app.Interactor) *Handler {
	return &Handler{me: me}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	meID := os.Getenv("ME_ID")
	if meID == "" {
		errs.WriteProblem(w, fmt.Errorf("ME_ID environment variable is not configured: %w", errs.ErrNotFound))
		return
	}
	out, err := h.me.Get(r.Context(), meID)
	if err != nil {
		errs.WriteProblem(w, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	meID := os.Getenv("ME_ID")
	if meID == "" {
		errs.WriteProblem(w, fmt.Errorf("ME_ID environment variable is not configured: %w", errs.ErrNotFound))
		return
	}
	var input app.InputDto
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errs.WriteProblem(w, fmt.Errorf("decode request body: %w", errs.ErrBadRequest))
		return
	}
	if err := validate.Struct(input); err != nil {
		errs.WriteProblem(w, fmt.Errorf("%s: %w", err.Error(), errs.ErrBadRequest))
		return
	}
	input.ID = meID
	out, err := h.me.Update(r.Context(), input)
	if err != nil {
		errs.WriteProblem(w, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v) //nolint:errcheck
}
