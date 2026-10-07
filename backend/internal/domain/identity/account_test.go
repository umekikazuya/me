package identity_test

import (
	"testing"

	"github.com/umekikazuya/me/internal/domain/identity"
)

func TestRegisterWithGithub(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		want    *identity.Account
		wantErr bool
	}{
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := identity.RegisterWithGithub(tt.id)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("RegisterWithGithub() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("RegisterWithGithub() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("RegisterWithGithub() = %v, want %v", got, tt.want)
			}
		})
	}
}
