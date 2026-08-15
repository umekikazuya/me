package me

import (
	"reflect"
	"slices"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewMe(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		inputID string
		wantErr bool
	}{
		{
			name:    "ok#正常に生成出来る",
			inputID: uuid.New().String(),
			wantErr: false,
		},
		{
			name:    "ng#uuid形式じゃない場合false",
			inputID: "test-id",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := NewMe(tt.inputID)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("NewMe() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("NewMe() succeeded unexpectedly")
			}
			if got.id.String() != tt.inputID {
				t.Errorf("got.id.String() = %v, want = %v", got.id.String(), tt.inputID)
			}
		})
	}
}

func Test_Reconstruct(t *testing.T) {
	displayJa := "田中 太郎"
	role := "Engineer"
	location := "Tokyo"
	fixedTime := func(s string) time.Time {
		t, _ := time.Parse(time.RFC3339, s)
		return t
	}
	createdAt := fixedTime("2024-01-01T00:00:00Z")
	updatedAt := fixedTime("2024-06-01T00:00:00Z")

	tests := []struct {
		name  string
		input ReconstructInput
		check func(*testing.T, *Me)
	}{
		{
			name: "full fields",
			input: ReconstructInput{
				Name:      "Taro",
				DisplayJa: &displayJa,
				Role:      &role,
				Location:  &location,
				Likes:     []string{"Go", "Rust"},
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			check: func(t *testing.T, m *Me) {
				if m.DisplayName() != "Taro" {
					t.Errorf("DisplayName() = %v, want Taro", m.DisplayName())
				}
				if m.DisplayNameJa() != displayJa {
					t.Errorf("DisplayNameJa() = %v, want %v", m.DisplayNameJa(), displayJa)
				}
				if m.Role() != role {
					t.Errorf("Role() = %v, want %v", m.Role(), role)
				}
				if m.Location() != location {
					t.Errorf("Location() = %v, want %v", m.Location(), location)
				}
				if !reflect.DeepEqual(m.Likes(), []string{"Go", "Rust"}) {
					t.Errorf("Likes() = %v, want [Go Rust]", m.Likes())
				}
				if !m.CreatedAt().Equal(createdAt) {
					t.Errorf("CreatedAt() = %v, want %v", m.CreatedAt(), createdAt)
				}
				if !m.UpdatedAt().Equal(updatedAt) {
					t.Errorf("UpdatedAt() = %v, want %v", m.UpdatedAt(), updatedAt)
				}
			},
		},
		{
			name: "minimal fields (nil optional)",
			input: ReconstructInput{
				Name:      "Minimal",
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			check: func(t *testing.T, m *Me) {
				if m.DisplayName() != "Minimal" {
					t.Errorf("DisplayName() = %v, want Minimal", m.DisplayName())
				}
				if m.DisplayNameJa() != "" {
					t.Errorf("DisplayNameJa() = %v, want empty", m.DisplayNameJa())
				}
				if m.Role() != "" {
					t.Errorf("Role() = %v, want empty", m.Role())
				}
				if m.Location() != "" {
					t.Errorf("Location() = %v, want empty", m.Location())
				}
				if len(m.Likes()) != 0 {
					t.Errorf("Likes() = %v, want empty", m.Likes())
				}
			},
		},
		{
			name: "createdAt and updatedAt are preserved",
			input: ReconstructInput{
				Name:      "Taro",
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			check: func(t *testing.T, m *Me) {
				if !m.CreatedAt().Equal(createdAt) {
					t.Errorf("CreatedAt() = %v, want %v", m.CreatedAt(), createdAt)
				}
				if !m.UpdatedAt().Equal(updatedAt) {
					t.Errorf("UpdatedAt() = %v, want %v", m.UpdatedAt(), updatedAt)
				}
			},
		},
		{
			name: "links are restored",
			input: ReconstructInput{
				Name: "Taro",
				Links: []Link{
					{platform: "github", url: "https://github.com/example"},
				},
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			check: func(t *testing.T, m *Me) {
				links := m.Links()
				if len(links) != 1 {
					t.Fatalf("Links() len = %d, want 1", len(links))
				}
				if links[0].Platform() != "github" {
					t.Errorf("Platform() = %v, want github", links[0].Platform())
				}
				if links[0].URL() != "https://github.com/example" {
					t.Errorf("URL() = %v, want https://github.com/example", links[0].URL())
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Reconstruct(tt.input)
			tt.check(t, got)
		})
	}
}

func Test_Me_Getters(t *testing.T) {
	t.Run("Check default values and timestamps", func(t *testing.T) {
		m, err := NewMe(uuid.New().String())
		if err != nil {
			t.Fatalf("NewMe failed: %v", err)
		}

		if m.DisplayNameJa() != "" {
			t.Errorf("expected empty string, got %v", m.DisplayNameJa())
		}
		if m.Role() != "" {
			t.Errorf("expected empty string, got %v", m.Role())
		}

		// CreatedAt, UpdatedAt の検証を復活
		if m.CreatedAt().IsZero() {
			t.Error("expected CreatedAt to be set (not zero)")
		}
		if m.UpdatedAt().IsZero() {
			t.Error("expected UpdatedAt to be set (not zero)")
		}
	})
}

func TestMe_updateProfile(t *testing.T) {
	testID := uuid.New().String()
	baseTime := time.Now().Add(24 * time.Hour)
	tests := []struct {
		name     string
		baseTime time.Time
		in       []OptProfileFunc
		wantErr  bool
		assertFn func(t *testing.T, e *Me, baseTime time.Time)
	}{
		{
			name: "ok#正常に更新できる",
			in: []OptProfileFunc{
				OptDisplayNameJa("abc"),
				OptLocation("abc"),
				OptDisplayName("abc"),
				OptRole("abc"),
			},
			baseTime: baseTime,
			wantErr:  false,
			assertFn: func(t *testing.T, e *Me, baseTime time.Time) {
				t.Helper()
				if e.profile.displayName != "abc" {
					t.Errorf("e.profile.displayName = %v", e.profile.displayName)
				}
				if e.profile.displayNameJa != "abc" {
					t.Errorf("e.profile.displayNameJa = %v", e.profile.displayNameJa)
				}
				if e.profile.role != "abc" {
					t.Errorf("e.profile.role = %v", e.profile.role)
				}
				if e.profile.location != "abc" {
					t.Errorf("e.profile.location = %v", e.profile.location)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := NewMe(testID)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			gotErr := e.UpdateProfile(tt.baseTime, tt.in...)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("updateProfile() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("updateProfile() succeeded unexpectedly")
			}
			if e.profile.displayName == "" {
				t.Errorf("必須フィールド e.profile.displayName = %v", e.profile.displayName)
			}
			if !e.updatedAt.Equal(baseTime) {
				t.Errorf("e.updatedAt = %v, baseTime = %v", e.updatedAt, baseTime)
			}
			tt.assertFn(t, e, tt.baseTime)
		})
	}
}

func TestMe_UpdateLikes(t *testing.T) {
	testID := uuid.New().String()
	baseTime := time.Now().Add(24 * time.Hour)
	tests := []struct {
		name     string
		in       []string
		baseTime time.Time
		wantErr  bool
		assertFn func(t *testing.T, e *Me)
	}{
		{
			name:     "ok#正常に更新できる",
			in:       []string{"go", "rust"},
			baseTime: baseTime,
			wantErr:  false,
			assertFn: func(t *testing.T, e *Me) {
				t.Helper()
				if !slices.Contains(e.likes, like{"go"}) {
					t.Fatalf("e.likes = %v", e.likes)
				}
				if !e.updatedAt.Equal(baseTime) {
					t.Errorf("e.updatedAt = %v, want = %v", e.updatedAt, baseTime)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := NewMe(testID)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			gotErr := e.UpdateLikes(tt.in, tt.baseTime)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("UpdateLikes() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("UpdateLikes() succeeded unexpectedly")
			}
			tt.assertFn(t, e)
		})
	}
}

func TestMe_UpdateLinks(t *testing.T) {
	testID := uuid.New().String()
	baseTime := time.Now().Add(24 * time.Hour)

	tests := []struct {
		name     string
		in       []Link
		wantErr  bool
		assertFn func(t *testing.T, e *Me)
	}{
		{
			name: "ok#正常に更新できる",
			in: []Link{
				{
					platform: "a",
					url:      "example.com",
				},
			},
			wantErr: false,
			assertFn: func(t *testing.T, e *Me) {
				t.Helper()
				links := e.links
				l := links[0]
				if l.platform != "a" {
					t.Errorf("l.platform = %v, want = %v", l.platform, "a")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := NewMe(testID)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			gotErr := e.UpdateLinks(tt.in, baseTime)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("UpdateLinks() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("UpdateLinks() succeeded unexpectedly")
			}
			tt.assertFn(t, e)
		})
	}
}
