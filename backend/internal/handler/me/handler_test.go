package me

import (
	"context"
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
	getFn    func(ctx context.Context, id string) (*app.OutputDto, error)
}

func (m *mockInteractor) Create(ctx context.Context, input app.InputDto) (*app.OutputDto, error) {
	if m.createFn != nil {
		return m.createFn(ctx, input)
	}
	return nil, nil
}

func (m *mockInteractor) Update(ctx context.Context, input app.InputDto) (*app.OutputDto, error) {
	return m.updateFn(ctx, input)
}

func (m *mockInteractor) Get(ctx context.Context, id string) (*app.OutputDto, error) {
	return m.getFn(ctx, id)
}

func TestHandler_Get(t *testing.T) {
	tests := []struct {
		name       string
		getFn      func(ctx context.Context, id string) (*app.OutputDto, error)
		wantStatus int
	}{
		{
			name: "success",
			getFn: func(_ context.Context, _ string) (*app.OutputDto, error) {
				return &app.OutputDto{DisplayName: "Taro"}, nil
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "not found",
			getFn: func(_ context.Context, _ string) (*app.OutputDto, error) {
				return nil, fmt.Errorf("get me: %w", errs.ErrNotFound)
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "internal error",
			getFn: func(_ context.Context, _ string) (*app.OutputDto, error) {
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
