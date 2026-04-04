package me

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/umekikazuya/me/internal/app/me"
	app "github.com/umekikazuya/me/internal/app/me"
	"github.com/umekikazuya/me/pkg/errs"
)

var validate = validator.New()

type Handler struct {
	me me.Interactor
}

func NewHandler(me me.Interactor) *Handler {
	return &Handler{me: me}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	out, err := h.me.Get(r.Context(), os.Getenv("ME_ID"))
	if err != nil {
		errs.WriteProblem(w, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var input app.InputDto
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errs.WriteProblem(w, fmt.Errorf("decode request body: %w", errs.ErrBadRequest))
		return
	}
	if err := validate.Struct(input); err != nil {
		errs.WriteProblem(w, fmt.Errorf("%s: %w", err.Error(), errs.ErrBadRequest))
		return
	}
	input.ID = os.Getenv("ME_ID")
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
