package article

import (
	"fmt"
	"net/http"
	"strconv"

	app "github.com/umekikazuya/me/internal/app/article"
	"github.com/umekikazuya/me/pkg/errs"
	"github.com/umekikazuya/me/pkg/httpx"
)

type Handler struct {
	interactor app.Interactor
}

func NewHandler(interactor app.Interactor) *Handler {
	return &Handler{interactor: interactor}
}

// Search handles GET /articles
func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	input := app.InputSearchDto{Limit: 50}

	if v := q.Get("q"); v != "" {
		input.Q = &v
	}
	input.Tag = q["tag"]
	if v := q.Get("platform"); v != "" {
		input.Platform = &v
	}
	if v := q.Get("year"); v != "" {
		year, err := strconv.Atoi(v)
		if err != nil || year <= 0 {
			errs.WriteProblem(w, r, fmt.Errorf("year must be a positive integer: %w", errs.ErrBadRequest))
			return
		}
		input.Year = &year
	}
	if v := q.Get("limit"); v != "" {
		limit, err := strconv.Atoi(v)
		if err != nil || limit < 1 || limit > 100 {
			errs.WriteProblem(w, r, fmt.Errorf("limit must be between 1 and 100: %w", errs.ErrBadRequest))
			return
		}
		input.Limit = limit
	}
	if v := q.Get("cursor"); v != "" {
		input.NextCursor = &v
	}

	out, err := h.interactor.Search(r.Context(), input)
	if err != nil {
		errs.WriteProblem(w, r, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, out)
}

// GetTagsAll handles GET /articles/meta/tags
func (h *Handler) GetTagsAll(w http.ResponseWriter, r *http.Request) {
	out, err := h.interactor.GetTagsAll(r.Context())
	if err != nil {
		errs.WriteProblem(w, r, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, out)
}

// GetSuggests handles GET /articles/meta/suggest
func (h *Handler) GetSuggests(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		errs.WriteProblem(w, r, fmt.Errorf("q is required: %w", errs.ErrBadRequest))
		return
	}
	out, err := h.interactor.GetSuggests(r.Context(), app.InputGetSuggestDto{Q: q})
	if err != nil {
		errs.WriteProblem(w, r, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, out)
}

// Register handles POST /articles
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var input app.InputRegisterDto
	if err := httpx.DecodeAndValidate(w, r, &input); err != nil {
		errs.WriteProblem(w, r, err)
		return
	}
	if err := h.interactor.Register(r.Context(), input); err != nil {
		errs.WriteProblem(w, r, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// Update handles PUT /articles/{externalId}
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var input app.InputUpdateDto
	if err := httpx.DecodeAndValidate(w, r, &input); err != nil {
		errs.WriteProblem(w, r, err)
		return
	}
	input.ExternalID = r.PathValue("externalId")
	if input.ExternalID == "" {
		errs.WriteProblem(w, r, fmt.Errorf("externalId is required %w", errs.ErrBadRequest))
		return
	}
	if err := h.interactor.Update(r.Context(), input); err != nil {
		errs.WriteProblem(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Remove handles DELETE /articles/{externalId}
func (h *Handler) Remove(w http.ResponseWriter, r *http.Request) {
	externalID := r.PathValue("externalId")
	if externalID == "" {
		errs.WriteProblem(w, r, fmt.Errorf("externalId is required %w", errs.ErrBadRequest))
		return
	}
	if err := h.interactor.Remove(r.Context(), app.InputRemoveDto{ExternalID: externalID}); err != nil {
		errs.WriteProblem(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
