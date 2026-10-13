package identity

import (
	"testing"
)

func TestNewTokenHash(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "valid sha256 hex", input: "a3f9b2c1d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1", want: "a3f9b2c1d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1", wantErr: false},
		{name: "empty", input: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTokenHash(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTokenHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Value() != tt.want {
				t.Errorf("NewTokenHash().Value() = %v, want %v", got.Value(), tt.want)
			}
		})
	}
}
