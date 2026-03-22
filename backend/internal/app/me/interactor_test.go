package me

import (
	"context"
	"errors"
	"reflect"
	"testing"

	domain "github.com/umekikazuya/me/internal/domain/me"
)

// MockRepo is a mock implementation of domain.Repo
type MockRepo struct {
	findFn   func(ctx context.Context) (*domain.Me, error)
	saveFn   func(ctx context.Context, e *domain.Me) error
	existsFn func(ctx context.Context) (bool, error)
}

func (m *MockRepo) Find(ctx context.Context) (*domain.Me, error) {
	return m.findFn(ctx)
}

func (m *MockRepo) Save(ctx context.Context, e *domain.Me) error {
	return m.saveFn(ctx, e)
}

func (m *MockRepo) Exists(ctx context.Context) (bool, error) {
	return m.existsFn(ctx)
}

func TestInteractor_Create(t *testing.T) {
	displayJa := "田中 太郎"
	role := "Engineer"
	location := "Tokyo"
	likes := []string{"Go", "Rust"}

	notExists := func(_ context.Context) (bool, error) { return false, nil }

	tests := []struct {
		name     string
		input    InputDto
		existsFn func(ctx context.Context) (bool, error)
		saveFn   func(ctx context.Context, e *domain.Me) error
		wantErr  bool
		check    func(*testing.T, *OutputDto)
	}{
		{
			name: "success: full fields provided",
			input: InputDto{
				DisplayName: "Taro",
				DisplayJa:   &displayJa,
				Role:        &role,
				Location:    &location,
				Likes:       likes,
			},
			existsFn: notExists,
			saveFn:   func(ctx context.Context, e *domain.Me) error { return nil },
			check: func(t *testing.T, got *OutputDto) {
				if got.DisplayName != "Taro" || got.DisplayJa != displayJa || got.Role != role || got.Location != location {
					t.Errorf("unexpected output fields: %+v", got)
				}
				if !reflect.DeepEqual(got.Likes, likes) {
					t.Errorf("got likes %v, want %v", got.Likes, likes)
				}
			},
		},
		{
			name: "success: minimal fields (nil pointers provided)",
			input: InputDto{
				DisplayName: "Minimal",
				DisplayJa:   nil,
				Role:        nil,
				Location:    nil,
				Likes:       nil,
			},
			existsFn: notExists,
			saveFn:   func(ctx context.Context, e *domain.Me) error { return nil },
			check: func(t *testing.T, got *OutputDto) {
				// マッパーが nil を安全に（空文字などで）扱えているか検証
				if got.DisplayJa != "" || got.Role != "" || got.Location != "" || len(got.Likes) != 0 {
					t.Errorf("expected empty values for nil inputs, got: %+v", got)
				}
			},
		},
		{
			name:     "error: domain validation (empty name)",
			input:    InputDto{DisplayName: ""},
			existsFn: notExists,
			wantErr:  true,
		},
		{
			name:     "error: repository failure",
			input:    InputDto{DisplayName: "Taro"},
			existsFn: notExists,
			saveFn: func(ctx context.Context, e *domain.Me) error {
				return errors.New("database error")
			},
			wantErr: true,
		},
		{
			name:     "error: conflict",
			input:    InputDto{DisplayName: "Taro"},
			existsFn: func(_ context.Context) (bool, error) { return true, nil },
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Interactor{
				repo: &MockRepo{existsFn: tt.existsFn, saveFn: tt.saveFn},
			}
			got, err := i.Create(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Interactor.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, got)
			}
		})
	}
}

func TestInteractor_Update_PUTBehavior(t *testing.T) {
	// PUT（置換）スタイルの更新: 指定しなかったフィールドが消えることを検証
	role := "Designer"

	tests := []struct {
		name  string
		input InputDto
		check func(*testing.T, *OutputDto)
	}{
		{
			name: "update with partial fields (others should be cleared)",
			input: InputDto{
				DisplayName: "NewName",
				Role:        &role,
				// DisplayJa, Location, Likes は nil/空
			},
			check: func(t *testing.T, got *OutputDto) {
				if got.DisplayName != "NewName" {
					t.Errorf("got name %v, want NewName", got.DisplayName)
				}
				if got.Role != role {
					t.Errorf("got role %v, want %v", got.Role, role)
				}
				// 指定しなかったフィールドは空になっているべき
				if got.DisplayJa != "" {
					t.Errorf("expected DisplayJa to be cleared, got %v", got.DisplayJa)
				}
				if got.Location != "" {
					t.Errorf("expected Location to be cleared, got %v", got.Location)
				}
				if len(got.Likes) != 0 {
					t.Errorf("expected Likes to be cleared, got %v", got.Likes)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 初期状態の Entity を準備
			initialMe, _ := domain.NewMe("OldName", domain.OptRole("OldRole"))

			i := &Interactor{
				repo: &MockRepo{
					findFn: func(ctx context.Context) (*domain.Me, error) {
						return initialMe, nil
					},
					saveFn: func(ctx context.Context, e *domain.Me) error { return nil },
				},
			}
			got, err := i.Update(context.Background(), tt.input)
			if err != nil {
				t.Fatalf("Update failed: %v", err)
			}
			tt.check(t, got)
		})
	}
}

func TestInteractor_Get(t *testing.T) {
	displayJa := "田中 太郎"

	tests := []struct {
		name    string
		findFn  func(ctx context.Context) (*domain.Me, error)
		wantErr bool
		check   func(*testing.T, *OutputDto)
	}{
		{
			name: "success get",
			findFn: func(ctx context.Context) (*domain.Me, error) {
				e, _ := domain.NewMe("Taro", domain.OptDisplayNameJa(displayJa))
				return e, nil
			},
			wantErr: false,
			check: func(t *testing.T, got *OutputDto) {
				if got.DisplayName != "Taro" || got.DisplayJa != displayJa {
					t.Errorf("unexpected output: %+v", got)
				}
			},
		},
		{
			name: "error repo find",
			findFn: func(ctx context.Context) (*domain.Me, error) {
				return nil, errors.New("not found")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Interactor{
				repo: &MockRepo{findFn: tt.findFn},
			}
			got, err := i.Get(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("Interactor.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, got)
			}
		})
	}
}
