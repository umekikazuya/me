package me

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/google/uuid"
	domain "github.com/umekikazuya/me/internal/domain/me"
)

// MockRepo is a mock implementation of domain.Repo
type MockRepo struct {
	findByIDFn func(ctx context.Context, id string) (*domain.Me, error)
	saveFn     func(ctx context.Context, e *domain.Me) error
	existsFn   func(ctx context.Context, id string) (bool, error)
}

func (m *MockRepo) FindByID(ctx context.Context, id string) (*domain.Me, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *MockRepo) Save(ctx context.Context, e *domain.Me) error {
	return m.saveFn(ctx, e)
}

func (m *MockRepo) Exists(ctx context.Context, id string) (bool, error) {
	return m.existsFn(ctx, id)
}

func TestInteractor_Create(t *testing.T) {
	testID := uuid.New().String()

	tests := []struct {
		name     string
		input    InputDto
		wantErr  bool
		assertFn func(*testing.T, *OutputDto, *memoryMeRepo)
	}{
		{
			name: "ok#full fields provided",
			input: InputDto{
				ID: testID,
			},
			assertFn: func(t *testing.T, got *OutputDto, repo *memoryMeRepo) {
				e, err := repo.FindByID(t.Context(), testID)
				if err != nil {
					t.Fatalf("err = %v", err)
				}
				if e.ID() != testID {
					t.Errorf("e.ID() = %v, want = %v", e.ID(), testID)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMeRepo()
			i := &interactor{repo: repo, id: testID}
			got, err := i.Create(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Interactor.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.assertFn != nil {
				tt.assertFn(t, got, repo)
			}
		})
	}
}

func TestInteractor_Get(t *testing.T) {
	testID := uuid.New()

	tests := []struct {
		name     string
		seedFn   func(t *testing.T, repo *memoryMeRepo)
		wantErr  bool
		assertFn func(t *testing.T, repo *memoryMeRepo)
	}{
		{
			name: "ok#正常に取得できる",
			seedFn: func(t *testing.T, repo *memoryMeRepo) {
				t.Helper()
				repo.seedData(t, domain.ReconstructInput{
					ID:             testID,
					Name:           "abcde",
					DisplayJa:      new(string),
					Role:           new(string),
					Location:       new(string),
					Likes:          []string{},
					Links:          []domain.Link{},
					Certifications: []domain.Certification{},
					CreatedAt:      time.Time{},
					UpdatedAt:      time.Time{},
				})
			},
			wantErr: false,
			assertFn: func(t *testing.T, repo *memoryMeRepo) {
				t.Helper()
				e, err := repo.FindByID(t.Context(), testID.String())
				if err != nil {
					t.Fatal(err)
				}
				if e.DisplayName() != "abcde" {
					t.Errorf("e.DisplayName = %v, want = abcde", e.DisplayName())
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMeRepo()
			tt.seedFn(t, repo)
			i := &interactor{repo: repo, id: testID.String()}
			_, err := i.Get(t.Context(), testID.String())
			if (err != nil) != tt.wantErr {
				t.Errorf("Interactor.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.assertFn != nil {
				tt.assertFn(t, repo)
			}
		})
	}
}

func Test_interactor_UpdateLikes(t *testing.T) {
	testID := uuid.New()
	tests := []struct {
		name     string
		in       InputUpdateLikes
		seedFn   func(*testing.T, *memoryMeRepo)
		wantErr  bool
		assertFn func(*testing.T, *memoryMeRepo)
	}{
		{
			name: "ok#正常に更新できる",
			in:   InputUpdateLikes{"abc", "bcd"},
			seedFn: func(t *testing.T, repo *memoryMeRepo) {
				t.Helper()
				repo.seedData(t, domain.ReconstructInput{
					ID:    testID,
					Likes: []string{},
				})
			},
			wantErr: false,
			assertFn: func(t *testing.T, repo *memoryMeRepo) {
				t.Helper()
				e, err := repo.FindByID(t.Context(), testID.String())
				if err != nil {
					t.Fatalf("err = %#v", err)
				}
				if len(e.Likes()) != 2 {
					t.Errorf("len(e.Likes()) = %v, want = 2", len(e.Likes()))
				}
			},
		},
		{
			name: "ok#既存のデータを上書きできる",
			in:   InputUpdateLikes{"abc", "bcd"},
			seedFn: func(t *testing.T, repo *memoryMeRepo) {
				t.Helper()
				repo.seedData(t, domain.ReconstructInput{
					ID:    testID,
					Likes: []string{"xyz"},
				})
			},
			wantErr: false,
			assertFn: func(t *testing.T, repo *memoryMeRepo) {
				t.Helper()
				e, err := repo.FindByID(t.Context(), testID.String())
				if err != nil {
					t.Fatalf("err = %#v", err)
				}
				if len(e.Likes()) != 2 {
					t.Errorf("len(e.Likes()) = %v, want = 2", len(e.Likes()))
				}
				if slices.Contains(e.Likes(), "xyz") {
					t.Errorf("e.Likes() = %v", e.Likes())
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			repo := newMeRepo()
			tt.seedFn(t, repo)
			i := interactor{repo: repo, id: testID.String()}

			// Act
			_, gotErr := i.UpdateLikes(t.Context(), tt.in)

			// Assert
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("UpdateLikes() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("UpdateLikes() succeeded unexpectedly")
			}
			tt.assertFn(t, repo)
		})
	}
}

func Test_interactor_UpdateLinks(t *testing.T) {
	testID := uuid.New()
	tests := []struct {
		name     string
		in       InputUpdateLinks
		seedFn   func(*testing.T, *memoryMeRepo)
		wantErr  bool
		assertFn func(*testing.T, *memoryMeRepo)
	}{
		{
			name: "",
			in: InputUpdateLinks{
				InputLink{
					Label:    "abc",
					Platform: "abc",
					URL:      "https://example.com",
				},
				InputLink{
					Platform: "def",
					URL:      "https://example.com",
				},
			},
			seedFn: func(t *testing.T, repo *memoryMeRepo) {
				t.Helper()
				repo.seedData(t, domain.ReconstructInput{
					ID:        testID,
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
				})
			},
			wantErr: false,
			assertFn: func(t *testing.T, repo *memoryMeRepo) {
				t.Helper()
				e, err := repo.FindByID(t.Context(), testID.String())
				if err != nil {
					t.Fatalf("err = %#v", err)
				}
				if len(e.Links()) != 2 {
					t.Errorf("len(e.Links()) = %v, want = 2", len(e.Likes()))
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			repo := newMeRepo()
			tt.seedFn(t, repo)
			i := interactor{repo: repo, id: testID.String()}

			// Act
			_, gotErr := i.UpdateLinks(t.Context(), tt.in)

			// Assert
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("UpdateLinks() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("UpdateLinks() succeeded unexpectedly")
			}
			tt.assertFn(t, repo)
		})
	}
}

func Test_interactor_UpdateProfile(t *testing.T) {
	testID := uuid.New()
	tests := []struct {
		name     string
		in       InputUpdateProfile
		seedFn   func(*testing.T, *memoryMeRepo)
		wantErr  bool
		assertFn func(*testing.T, *memoryMeRepo)
	}{
		{
			name: "",
			in: InputUpdateProfile{
				Location:    "tokyo",
				DisplayName: "abc",
				DisplayJa:   "あいう",
				Role:        "role",
			},
			seedFn: func(t *testing.T, repo *memoryMeRepo) {
				t.Helper()
				repo.seedData(t, domain.ReconstructInput{ID: testID})
			},
			wantErr: false,
			assertFn: func(t *testing.T, repo *memoryMeRepo) {
				t.Helper()
				e, err := repo.FindByID(t.Context(), testID.String())
				if err != nil {
					t.Fatalf("err = %#v", err)
				}
				if e.Role() != "role" {
					t.Errorf("e.Role() = %v, want = role", e.Role())
				}
				if e.DisplayName() != "abc" {
					t.Errorf("e.DisplayName = %v, want = abc", e.DisplayName())
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			repo := newMeRepo()
			tt.seedFn(t, repo)
			i := interactor{repo: repo, id: testID.String()}

			// Act
			_, gotErr := i.UpdateProfile(t.Context(), tt.in)

			// Assert
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("UpdateProfile() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("UpdateProfile() succeeded unexpectedly")
			}
			tt.assertFn(t, repo)
		})
	}
}
