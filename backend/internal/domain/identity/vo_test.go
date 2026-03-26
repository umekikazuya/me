package identity

import (
	"strings"
	"testing"
)

func TestNewEmail(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "valid", input: "user@example.com", want: "user@example.com", wantErr: false},
		{name: "empty", input: "", wantErr: true},
		{name: "no @ symbol", input: "userexample.com", wantErr: true},
		{name: "missing local part", input: "@example.com", wantErr: true},
		{name: "missing domain", input: "user@", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEmail(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Value() != tt.want {
				t.Errorf("NewEmail().Value() = %v, want %v", got.Value(), tt.want)
			}
		})
	}
}

func TestNewPassword(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{name: "valid", input: "Password1", wantErr: false},
		{name: "exactly 8 chars", input: "Passwor1", wantErr: false},
		{name: "exactly 72 chars", input: strings.Repeat("Aa", 36), wantErr: false},
		{name: "too short (7 chars)", input: "Pass1Ab", wantErr: true},
		{name: "too long (73 chars)", input: strings.Repeat("Aa", 36) + "B", wantErr: true},
		{name: "no uppercase", input: "password1", wantErr: true},
		{name: "no lowercase", input: "PASSWORD1", wantErr: true},
		{name: "empty", input: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewPassword(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewPasswordHash(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		wantErr bool
	}{
		{name: "valid hash", input: []byte("$2a$12$somebcrypthashvalue"), wantErr: false},
		{name: "empty hash", input: []byte{}, wantErr: true},
		{name: "nil hash", input: nil, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPasswordHash(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && string(got.Value()) != string(tt.input) {
				t.Errorf("NewPasswordHash().Value() = %v, want %v", got.Value(), tt.input)
			}
		})
	}
}

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
