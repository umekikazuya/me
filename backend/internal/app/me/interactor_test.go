package me

import (
	"errors"
	"reflect"
	"slices"
	"testing"
	"time"

	"github.com/google/uuid"
	domain "github.com/umekikazuya/me/internal/domain/me"
)

var (
	targetID = uuid.New()
	now      = time.Now()
)

const (
	sampleLocation      = "tokyo"
	sampleName          = "abc name"
	sampleNameJa        = "あいう"
	sampleRole          = "role"
	sampleTagName       = "tagA"
	sampleTagNameParent = "parentTagA"
)

func TestInteractor_Get(t *testing.T) {
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
					Name:           sampleName,
					DisplayJa:      new(string),
					Role:           new(string),
					Location:       new(string),
					Likes:          []string{},
					Links:          []domain.Link{},
					Certifications: []domain.Certification{},
					UpdatedAt:      time.Time{},
				})
			},
			wantErr: false,
			assertFn: func(t *testing.T, repo *memoryMeRepo) {
				t.Helper()
				e, err := repo.Find(t.Context())
				if err != nil {
					t.Fatal(err)
				}
				if e.DisplayName() != sampleName {
					t.Errorf("e.DisplayName = %v, want = abcde", e.DisplayName())
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMeRepo()
			tt.seedFn(t, repo)
			i := &interactor{repo: repo}
			_, err := i.Get(t.Context())
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
					Likes: []string{},
				})
			},
			wantErr: false,
			assertFn: func(t *testing.T, repo *memoryMeRepo) {
				t.Helper()
				e, err := repo.Find(t.Context())
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
					Likes: []string{"xyz"},
				})
			},
			wantErr: false,
			assertFn: func(t *testing.T, repo *memoryMeRepo) {
				t.Helper()
				e, err := repo.Find(t.Context())
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
			i := interactor{repo: repo}

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
	tests := []struct {
		name     string
		in       InputUpdateLinks
		seedFn   func(*testing.T, *memoryMeRepo)
		wantErr  bool
		assertFn func(*testing.T, *memoryMeRepo)
	}{
		{
			name: "ok#正常に更新できるか",
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
					UpdatedAt: time.Time{},
				})
			},
			wantErr: false,
			assertFn: func(t *testing.T, repo *memoryMeRepo) {
				t.Helper()
				e, err := repo.Find(t.Context())
				if err != nil {
					t.Fatalf("err = %#v", err)
				}
				if len(e.Links()) != 2 {
					t.Errorf("len(e.Links()) = %v, want = 2", len(e.Links()))
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			repo := newMeRepo()
			tt.seedFn(t, repo)
			i := interactor{repo: repo}

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

func Test_interactor_AddSkill(t *testing.T) {
	tests := []struct {
		name     string
		in       InputAddSkill
		seedFn   func(t *testing.T, repo *memoryMeRepo)
		wantErr  error
		assertFn func(t *testing.T, repo *memoryMeRepo)
	}{
		{
			name: "ok",
			in: InputAddSkill{
				Name:   sampleTagName,
				Parent: sampleTagNameParent,
			},
			seedFn: func(t *testing.T, repo *memoryMeRepo) {
				t.Helper()
				repo.seedData(t, domain.ReconstructInput{
					Name:      sampleName,
					Skills:    domain.Skills{},
					UpdatedAt: now,
				})
			},
			wantErr: nil,
			assertFn: func(t *testing.T, repo *memoryMeRepo) {
				t.Helper()
				e, err := repo.Find(t.Context())
				if err != nil {
					t.Fatalf("err = %v", err)
				}
				if e == nil {
					t.Fatal("e = nil")
				}
				if len(e.Skills()) != 1 {
					t.Errorf("len(e.Skills()) = %d", len(e.Skills()))
				}
				if len(e.Skills()[sampleTagNameParent].Items) != 1 {
					t.Errorf("len(e.Skills) = %d", len(e.Skills()))
				}
				if !slices.Contains(e.Skills()[sampleTagNameParent].Items, sampleTagName) {
					t.Errorf("e.Skills() = %v", e.Skills()[sampleTagNameParent].Items)
				}
			},
		},
		{
			name: "ok_original",
			in: InputAddSkill{
				Name:   sampleTagName,
				Parent: sampleTagNameParent,
			},
			seedFn: func(t *testing.T, repo *memoryMeRepo) {
				t.Helper()
				repo.seedData(t, domain.ReconstructInput{
					Name: sampleName,
					Skills: domain.Skills{
						sampleTagNameParent: {Items: []string{"def"}},
					},
					UpdatedAt: now,
				})
			},
			wantErr: nil,
			assertFn: func(t *testing.T, repo *memoryMeRepo) {
				t.Helper()
				e, err := repo.Find(t.Context())
				if err != nil {
					t.Fatalf("err = %v", err)
				}
				if e == nil {
					t.Fatal("e = nil")
				}
				if len(e.Skills()) != 1 {
					t.Errorf("len(e.Skills) = %d", len(e.Skills()))
				}
				if len(e.Skills()[sampleTagNameParent].Items) != 2 {
					t.Errorf(
						"len(e.Skills()[sampleTagNameParent].Items) = %d, e.Skills() = %v",
						len(e.Skills()[sampleTagNameParent].Items),
						e.Skills(),
					)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMeRepo()
			tt.seedFn(t, repo)
			sut := NewInteractor(repo)

			_, err := sut.AddSkill(t.Context(), tt.in)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want = %v", err, tt.wantErr)
			}
			tt.assertFn(t, repo)
		})
	}
}

func Test_interactor_RemoveSkill(t *testing.T) {
	tests := []struct {
		name    string
		in      InputRemoveSkill
		seedFn  func(t *testing.T, repo *memoryMeRepo)
		want    *OutputDto
		wantErr error
	}{
		{
			name: "ok#正常にスキルを削除出来る",
			in: InputRemoveSkill{
				Name:   sampleTagName,
				Parent: sampleTagNameParent,
			},
			seedFn: func(t *testing.T, repo *memoryMeRepo) {
				t.Helper()
				repo.seedData(t, domain.ReconstructInput{
					Name: sampleName,
					Skills: domain.Skills{
						sampleTagNameParent: {
							Items: []string{sampleTagName},
						},
					},
					UpdatedAt: now,
				})
			},
			want: &OutputDto{
				Skills: []struct {
					Category  string   "json:\"category\""
					Items     []string "json:\"items\""
					SortOrder int      "json:\"sortOrder\""
				}{},
			},
			wantErr: nil,
		},
		{
			name: "ok#正常にスキルを削除出来る",
			in: InputRemoveSkill{
				Name:   sampleTagName,
				Parent: sampleTagNameParent,
			},
			seedFn: func(t *testing.T, repo *memoryMeRepo) {
				t.Helper()
				repo.seedData(t, domain.ReconstructInput{
					Name: sampleName,
					Skills: domain.Skills{
						sampleTagNameParent: {
							Items: []string{"def", sampleTagName},
						},
					},
					UpdatedAt: now,
				})
			},
			want: &OutputDto{
				Skills: []struct {
					Category  string   "json:\"category\""
					Items     []string "json:\"items\""
					SortOrder int      "json:\"sortOrder\""
				}{
					{
						Category:  sampleTagNameParent,
						Items:     []string{"def"},
						SortOrder: 0,
					},
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMeRepo()
			tt.seedFn(t, repo)
			sut := NewInteractor(repo)

			got, err := sut.RemoveSkill(t.Context(), tt.in)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v", err)
			}
			if !reflect.DeepEqual(got.Skills, tt.want.Skills) {
				t.Errorf("got.Skills = %v, want = %v", got.Skills, tt.want.Skills)
			}
		})
	}
}

func Test_interactor_UpdateProfile(t *testing.T) {
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
				repo.seedData(t, domain.ReconstructInput{})
			},
			wantErr: false,
			assertFn: func(t *testing.T, repo *memoryMeRepo) {
				t.Helper()
				e, err := repo.Find(t.Context())
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
			i := interactor{repo: repo}

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
