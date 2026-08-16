package me

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	app "github.com/umekikazuya/me/internal/app/me"
	"github.com/umekikazuya/me/pkg/errs"
)

type mockInteractor struct {
	createFn        func(ctx context.Context, in app.InputDto) (*app.OutputDto, error)
	updateProfileFn func(ctx context.Context, in app.InputUpdateProfile) (*app.OutputDto, error)
	updateLinksFn   func(ctx context.Context, in app.InputUpdateLinks) (*app.OutputDto, error)
	updateLikesFn   func(ctx context.Context, in app.InputUpdateLikes) (*app.OutputDto, error)
	getFn           func(ctx context.Context, id string) (*app.OutputDto, error)
}

func (m *mockInteractor) Create(ctx context.Context, in app.InputDto) (*app.OutputDto, error) {
	if m.createFn != nil {
		return m.createFn(ctx, in)
	}
	return nil, nil
}

func (m *mockInteractor) Get(ctx context.Context, id string) (*app.OutputDto, error) {
	return m.getFn(ctx, id)
}

func (m *mockInteractor) UpdateProfile(ctx context.Context, in app.InputUpdateProfile) (*app.OutputDto, error) {
	return m.updateProfileFn(ctx, in)
}

func (m *mockInteractor) UpdateLinks(ctx context.Context, in app.InputUpdateLinks) (*app.OutputDto, error) {
	return m.updateLinksFn(ctx, in)
}

func (m *mockInteractor) UpdateLikes(ctx context.Context, in app.InputUpdateLikes) (*app.OutputDto, error) {
	return m.updateLikesFn(ctx, in)
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
			h, err := NewHandler(
				&mockInteractor{
					getFn: tt.getFn,
				}, "me-id")
			if err != nil {
				t.Fatalf("NewHandler() error = %v", err)
			}
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/me", nil)
			h.Get(w, r)
			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestNewHandler(t *testing.T) {
	t.Run("missing me id", func(t *testing.T) {
		_, err := NewHandler(&mockInteractor{}, "")
		if err == nil {
			t.Fatal("NewHandler() error = nil, want non-nil")
		}
	})
}
