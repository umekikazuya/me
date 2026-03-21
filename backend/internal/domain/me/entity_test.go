package me

import (
	"reflect"
	"testing"
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
				if m.DisplayName().Value() != "Taro Tanaka" {
					t.Errorf("DisplayName() = %v, want %v", m.DisplayName().Value(), "Taro Tanaka")
				}
				if m.DisplayNameJa().Value() != "田中 太郎" {
					t.Errorf("DisplayNameJa() = %v, want %v", m.DisplayNameJa().Value(), "田中 太郎")
				}
				if m.Role().Value() != "Software Engineer" {
					t.Errorf("Role() = %v, want %v", m.Role().Value(), "Software Engineer")
				}
				if m.Location().Value() != "Tokyo, Japan" {
					t.Errorf("Location() = %v, want %v", m.Location().Value(), "Tokyo, Japan")
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
			name: "success with mandatory name only",
			args: args{
				name: "Minimal Me",
			},
			wantErr: false,
			check: func(t *testing.T, m *Me) {
				if m.DisplayName().Value() != "Minimal Me" {
					t.Errorf("DisplayName() = %v, want %v", m.DisplayName().Value(), "Minimal Me")
				}
				if m.DisplayNameJa() != nil {
					t.Errorf("DisplayNameJa() should be nil")
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
		// 初期状態: 全てセットされている
		m, _ := NewMe("Taro",
			OptDisplayNameJa("太郎"),
			OptRole("Engineer"),
			OptLocation("Tokyo"),
			OptLikes([]string{"Go"}),
		)

		// Updateを実行。引数は名前のみ。Optionsはなし。
		// 仕様: PUT想定なので、指定しなかったオプションフィールドはリセット(nil/空)されるべき
		err := m.Update("Jiro")
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		if m.DisplayName().Value() != "Jiro" {
			t.Errorf("DisplayName() = %v, want %v", m.DisplayName().Value(), "Jiro")
		}

		// オプションフィールドは全てリセットされていることを確認
		if m.DisplayNameJa() != nil {
			t.Error("expected DisplayNameJa to be cleared (nil)")
		}
		if m.Role() != nil {
			t.Error("expected Role to be cleared (nil)")
		}
		if m.Location() != nil {
			t.Error("expected Location to be cleared (nil)")
		}
		if len(m.Likes()) != 0 {
			t.Errorf("expected Likes to be cleared (empty), got %v", m.Likes())
		}
	})

	t.Run("PUT style update: partially specified fields should replace others", func(t *testing.T) {
		m, _ := NewMe("Taro",
			OptRole("Engineer"),
			OptLocation("Tokyo"),
		)

		// Roleだけ新しく指定。Locationは指定しない。
		err := m.Update("Taro", OptRole("Designer"))
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		if m.Role().Value() != "Designer" {
			t.Errorf("Role() = %v, want %v", m.Role().Value(), "Designer")
		}
		// 指定しなかったLocationは消えているべき
		if m.Location() != nil {
			t.Error("expected Location to be cleared")
		}
	})

	t.Run("success clearing likes explicitly", func(t *testing.T) {
		m, _ := NewMe("Taro",
			OptLikes([]string{"Go"}),
		)

		// 空のスライスを渡した場合も当然空になる
		err := m.Update("Taro", OptLikes([]string{}))
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		if len(m.Likes()) != 0 {
			t.Errorf("Likes() = %v, want empty slice", m.Likes())
		}
	})

	t.Run("error with invalid update values", func(t *testing.T) {
		m, _ := NewMe("Taro")

		// 不正なDisplayName (バリデーションエラー)
		if err := m.Update(""); err == nil {
			t.Error("expected error for empty display name")
		}

		// Option内でのバリデーションエラー
		if err := m.Update("Taro", OptRole("  ")); err == nil {
			t.Error("expected error for invalid role option")
		}
	})
}

func Test_Me_Getters(t *testing.T) {
	t.Run("Check default values for optional fields", func(t *testing.T) {
		m, err := NewMe("Default")
		if err != nil {
			t.Fatalf("NewMe failed: %v", err)
		}

		if m.DisplayNameJa() != nil {
			t.Error("expected DisplayNameJa to be nil")
		}
		if m.Role() != nil {
			t.Error("expected Role to be nil")
		}
		if m.Location() != nil {
			t.Error("expected Location to be nil")
		}
		if len(m.Likes()) != 0 {
			t.Errorf("expected Likes to be empty, got %v", m.Likes())
		}
		if !m.CreatedAt().IsZero() {
			t.Error("expected CreatedAt to be zero")
		}
		if !m.UpdatedAt().IsZero() {
			t.Error("expected UpdatedAt to be zero")
		}
	})

	t.Run("Check correctly set optional fields", func(t *testing.T) {
		m, _ := NewMe("Taro",
			OptDisplayNameJa("太郎"),
			OptLikes([]string{"A", "B"}),
		)

		if m.DisplayNameJa().Value() != "太郎" {
			t.Errorf("got %v, want %v", m.DisplayNameJa().Value(), "太郎")
		}
		if !reflect.DeepEqual(m.Likes(), []string{"A", "B"}) {
			t.Errorf("got %v, want %v", m.Likes(), []string{"A", "B"})
		}
	})
}
