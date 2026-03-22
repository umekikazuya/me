package me

import (
	"reflect"
	"testing"
	"time"
)

func Test_NewMe(t *testing.T) {
	type args struct {
		name string
		opts []OptFunc
	}
	tests := []struct {
		name    string
		args    args
		check   func(*testing.T, *Me)
		wantErr bool
	}{
		{
			name: "success with all valid options",
			args: args{
				name: "Taro Tanaka",
				opts: []OptFunc{
					OptDisplayNameJa("田中 太郎"),
					OptRole("Software Engineer"),
					OptLocation("Tokyo, Japan"),
					OptLikes([]string{"Go", "Rust"}),
				},
			},
			wantErr: false,
			check: func(t *testing.T, m *Me) {
				if m.DisplayName() != "Taro Tanaka" {
					t.Errorf("DisplayName() = %v, want %v", m.DisplayName(), "Taro Tanaka")
				}
				if m.DisplayNameJa() != "田中 太郎" {
					t.Errorf("DisplayNameJa() = %v, want %v", m.DisplayNameJa(), "田中 太郎")
				}
				if m.Role() != "Software Engineer" {
					t.Errorf("Role() = %v, want %v", m.Role(), "Software Engineer")
				}
				if m.Location() != "Tokyo, Japan" {
					t.Errorf("Location() = %v, want %v", m.Location(), "Tokyo, Japan")
				}
				expectedLikes := []string{"Go", "Rust"}
				if !reflect.DeepEqual(m.Likes(), expectedLikes) {
					t.Errorf("Likes() = %v, want %v", m.Likes(), expectedLikes)
				}
			},
		},
		{
			name: "error with empty mandatory name",
			args: args{
				name: "",
			},
			wantErr: true,
		},
		{
			name: "error with space only mandatory name",
			args: args{
				name: "   ",
			},
			wantErr: true,
		},
		{
			name: "error with invalid DisplayNameJa option",
			args: args{
				name: "Taro",
				opts: []OptFunc{
					OptDisplayNameJa(""),
				},
			},
			wantErr: true,
		},
		{
			name: "error with invalid Role option",
			args: args{
				name: "Taro",
				opts: []OptFunc{
					OptRole("  "),
				},
			},
			wantErr: true,
		},
		{
			name: "error with invalid Location option",
			args: args{
				name: "Taro",
				opts: []OptFunc{
					OptLocation(""),
				},
			},
			wantErr: true,
		},
		{
			name: "error with invalid Likes option (contains empty string)",
			args: args{
				name: "Taro",
				opts: []OptFunc{
					OptLikes([]string{"Go", ""}),
				},
			},
			wantErr: true,
		},
		{
			name: "error with nil option",
			args: args{
				name: "Taro",
				opts: []OptFunc{
					nil,
				},
			},
			wantErr: true,
		},
		{
			name: "success with mandatory name only",
			args: args{
				name: "Minimal Me",
			},
			wantErr: false,
			check: func(t *testing.T, m *Me) {
				if m.DisplayName() != "Minimal Me" {
					t.Errorf("DisplayName() = %v, want %v", m.DisplayName(), "Minimal Me")
				}
				if m.DisplayNameJa() != "" {
					t.Errorf("DisplayNameJa() should be empty string, got %v", m.DisplayNameJa())
				}
				if len(m.Likes()) != 0 {
					t.Errorf("Likes() should be an empty slice, got %v", m.Likes())
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMe(tt.args.name, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, got)
			}
		})
	}
}

func Test_Me_Update(t *testing.T) {
	t.Run("PUT style update: unspecified fields should be cleared", func(t *testing.T) {
		m, err := NewMe("Taro",
			OptDisplayNameJa("太郎"),
			OptRole("Engineer"),
			OptLocation("Tokyo"),
			OptLikes([]string{"Go"}),
		)
		if err != nil {
			t.Fatalf("NewMe failed: %v", err)
		}

		// 名前のみ更新。他は指定しない。
		err = m.Update("Jiro")
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		if m.DisplayName() != "Jiro" {
			t.Errorf("DisplayName() = %v, want %v", m.DisplayName(), "Jiro")
		}

		// PUT仕様: 指定しなかったフィールドは空になる
		if m.DisplayNameJa() != "" {
			t.Errorf("expected DisplayNameJa to be empty string, got %v", m.DisplayNameJa())
		}
		if m.Role() != "" {
			t.Errorf("expected Role to be empty string, got %v", m.Role())
		}
		if m.Location() != "" {
			t.Errorf("expected Location to be empty string, got %v", m.Location())
		}
		if len(m.Likes()) != 0 {
			t.Errorf("expected Likes to be empty, got %v", m.Likes())
		}
	})

	t.Run("PUT style update: partially specified fields should replace others", func(t *testing.T) {
		m, err := NewMe("Taro",
			OptRole("Engineer"),
			OptLocation("Tokyo"),
		)
		if err != nil {
			t.Fatalf("NewMe failed: %v", err)
		}

		err = m.Update("Taro", OptRole("Designer"))
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		if m.Role() != "Designer" {
			t.Errorf("Role() = %v, want %v", m.Role(), "Designer")
		}
		if m.Location() != "" {
			t.Errorf("expected Location to be empty, got %v", m.Location())
		}
	})

	t.Run("error with invalid update values: state must remain unchanged", func(t *testing.T) {
		m, err := NewMe("Taro", OptRole("Engineer"))
		if err != nil {
			t.Fatalf("NewMe failed: %v", err)
		}

		beforeName := m.DisplayName()
		beforeRole := m.Role()

		// 不正なDisplayNameでの更新試行
		if err := m.Update(""); err == nil {
			t.Error("expected error for empty display name")
		}

		// 失敗後、状態が変わっていないことをチェック
		if m.DisplayName() != beforeName {
			t.Errorf("DisplayName changed after failed update: got %v, want %v", m.DisplayName(), beforeName)
		}
		if m.Role() != beforeRole {
			t.Errorf("Role changed after failed update: got %v, want %v", m.Role(), beforeRole)
		}

		// Option内でのバリデーションエラー
		if err := m.Update("Taro", OptRole("  ")); err == nil {
			t.Error("expected error for invalid role option")
		}

		if m.DisplayName() != beforeName || m.Role() != beforeRole {
			t.Error("state must remain unchanged after failed update with invalid option")
		}
	})
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
		m, err := NewMe("Default")
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

	t.Run("Check correctly set optional fields", func(t *testing.T) {
		m, err := NewMe("Taro",
			OptDisplayNameJa("太郎"),
			OptLikes([]string{"A", "B"}),
		)
		if err != nil {
			t.Fatalf("NewMe failed: %v", err)
		}

		if m.DisplayNameJa() != "太郎" {
			t.Errorf("got %v, want %v", m.DisplayNameJa(), "太郎")
		}
		if !reflect.DeepEqual(m.Likes(), []string{"A", "B"}) {
			t.Errorf("got %v, want %v", m.Likes(), []string{"A", "B"})
		}
	})
}
