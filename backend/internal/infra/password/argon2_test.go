package password_test

import (
	"context"
	"testing"

	"github.com/umekikazuya/me/internal/infra/password"
)

func TestArgon2PasswordManager_Verify(t *testing.T) {
	tests := []struct {
		name    string
		a       string
		b       string
		wantErr bool
	}{
		{
			name:    "ok#正常に検証できる",
			a:       "abc",
			b:       "abc",
			wantErr: false,
		},
		{
			name:    "ng#不一致",
			a:       "abc",
			b:       "ab0",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var a password.Argon2PasswordManager
			rawHashedBytes, err := a.Hash(t.Context(), tt.a)
			if err != nil {
				t.Fatalf("err = %#v", err)
			}
			gotErr := a.Verify(
				context.Background(),
				string(rawHashedBytes),
				tt.b,
			)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Verify() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Verify() succeeded unexpectedly")
			}
		})
	}
}
