package identity_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/umekikazuya/me/internal/domain/identity"
)

func TestRegisterWithGithub(t *testing.T) {
	tests := []struct {
		name     string
		inputID  string
		baseTime time.Time
		wantErr  bool
	}{
		{
			name:     "ok#正常に作成できる",
			inputID:  "123456",
			baseTime: baseTime,
			wantErr:  false,
		},
		{
			name:     "ゼロ値考慮",
			inputID:  "",
			baseTime: baseTime,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := identity.RegisterWithGithub(tt.inputID, tt.baseTime)

			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("RegisterWithGithub() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("RegisterWithGithub() succeeded unexpectedly")
			}
			if got.ID() == uuid.Nil.String() {
				t.Errorf("got.ID = %v", got.ID())
			}
			if got.GithubID() != tt.inputID {
				t.Errorf("got.GithubID = %v, want = %v", got.GithubID(), tt.inputID)
			}
			if got.CreatedAt() != baseTime {
				t.Errorf("got.CreatedAt = %v, baseTime = %v", got.CreatedAt(), baseTime)
			}
			if got.UpdatedAt() != baseTime {
				t.Errorf("got.UpdatedAt = %v, baseTime = %v", got.UpdatedAt(), baseTime)
			}
		})
	}
}
