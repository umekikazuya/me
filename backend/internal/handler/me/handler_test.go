package me

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	app "github.com/umekikazuya/me/internal/app/me"
	"github.com/umekikazuya/me/pkg/errs"
)

type mockInteractor struct {
	createFn func(ctx context.Context, input app.InputDto) (*app.OutputDto, error)
	updateFn func(ctx context.Context, input app.InputDto) (*app.OutputDto, error)
	getFn    func(ctx context.Context) (*app.OutputDto, error)
}

func (m *mockInteractor) Create(ctx context.Context, input app.InputDto) (*app.OutputDto, error) {
	return m.createFn(ctx, input)
}

func (m *mockInteractor) Update(ctx context.Context, input app.InputDto) (*app.OutputDto, error) {
	return m.updateFn(ctx, input)
}

func (m *mockInteractor) Get(ctx context.Context) (*app.OutputDto, error) {
	return m.getFn(ctx)
}

func TestHandler_Get(t *testing.T) {
	tests := []struct {
		name       string
		getFn      func(ctx context.Context) (*app.OutputDto, error)
		wantStatus int
	}{
		{
			name: "success",
			getFn: func(_ context.Context) (*app.OutputDto, error) {
				return &app.OutputDto{DisplayName: "Taro"}, nil
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "not found",
			getFn: func(_ context.Context) (*app.OutputDto, error) {
				return nil, fmt.Errorf("get me: %w", errs.ErrNotFound)
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "internal error",
			getFn: func(_ context.Context) (*app.OutputDto, error) {
				return nil, errors.New("unexpected")
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(&mockInteractor{getFn: tt.getFn})
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/me", nil)
			h.Get(w, r)
			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestHandler_Create(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		createFn   func(ctx context.Context, input app.InputDto) (*app.OutputDto, error)
		wantStatus int
	}{
		{
			name: "success",
			body: `{"displayName":"Taro"}`,
			createFn: func(_ context.Context, _ app.InputDto) (*app.OutputDto, error) {
				return &app.OutputDto{DisplayName: "Taro"}, nil
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "invalid json",
			body:       `{invalid}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "conflict",
			body: `{"displayName":"Taro"}`,
			createFn: func(_ context.Context, _ app.InputDto) (*app.OutputDto, error) {
				return nil, fmt.Errorf("create me: %w", errs.ErrConflict)
			},
			wantStatus: http.StatusConflict,
		},
		{
			name: "internal error",
			body: `{"displayName":"Taro"}`,
			createFn: func(_ context.Context, _ app.InputDto) (*app.OutputDto, error) {
				return nil, errors.New("unexpected")
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(&mockInteractor{createFn: tt.createFn})
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/me", strings.NewReader(tt.body))
			r.Header.Set("Content-Type", "application/json")
			h.Create(w, r)
			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestHandler_Update(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		updateFn   func(ctx context.Context, input app.InputDto) (*app.OutputDto, error)
		wantStatus int
	}{
		{
			name: "success",
			body: `{"displayName":"NewName"}`,
			updateFn: func(_ context.Context, _ app.InputDto) (*app.OutputDto, error) {
				return &app.OutputDto{DisplayName: "NewName"}, nil
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid json",
			body:       `{invalid}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "not found",
			body: `{"displayName":"NewName"}`,
			updateFn: func(_ context.Context, _ app.InputDto) (*app.OutputDto, error) {
				return nil, fmt.Errorf("update me: %w", errs.ErrNotFound)
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "internal error",
			body: `{"displayName":"NewName"}`,
			updateFn: func(_ context.Context, _ app.InputDto) (*app.OutputDto, error) {
				return nil, errors.New("unexpected")
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(&mockInteractor{updateFn: tt.updateFn})
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPut, "/me", strings.NewReader(tt.body))
			r.Header.Set("Content-Type", "application/json")
			h.Update(w, r)
			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestHandler_Create_ResponseBody(t *testing.T) {
	h := NewHandler(&mockInteractor{
		createFn: func(_ context.Context, _ app.InputDto) (*app.OutputDto, error) {
			return &app.OutputDto{DisplayName: "Taro"}, nil
		},
	})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/me", strings.NewReader(`{"displayName":"Taro"}`))
	r.Header.Set("Content-Type", "application/json")
	h.Create(w, r)

	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}
	var out app.OutputDto
	if err := json.NewDecoder(w.Body).Decode(&out); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if out.DisplayName != "Taro" {
		t.Errorf("DisplayName = %q, want Taro", out.DisplayName)
	}
}
