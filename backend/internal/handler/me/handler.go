package me

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	app "github.com/umekikazuya/me/internal/app/me"
	"github.com/umekikazuya/me/pkg/errs"
)

var validate = validator.New()

type meInteractor interface {
	Create(ctx context.Context, input app.InputDto) (*app.OutputDto, error)
	Update(ctx context.Context, input app.InputDto) (*app.OutputDto, error)
	Get(ctx context.Context) (*app.OutputDto, error)
}

type Handler struct {
	me meInteractor
}

func NewHandler(me meInteractor) *Handler {
	return &Handler{me: me}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	out, err := h.me.Get(r.Context())
	if err != nil {
		errs.WriteProblem(w, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var input app.InputDto
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errs.WriteProblem(w, fmt.Errorf("decode request body: %w", errs.ErrBadRequest))
		return
	}
	if err := validate.Struct(input); err != nil {
		errs.WriteProblem(w, fmt.Errorf("%s: %w", err.Error(), errs.ErrBadRequest))
		return
	}
	out, err := h.me.Create(r.Context(), input)
	if err != nil {
		errs.WriteProblem(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, out)
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
