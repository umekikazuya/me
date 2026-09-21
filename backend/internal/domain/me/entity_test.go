package me_test

import (
	"testing"
	"time"

	"github.com/umekikazuya/me/internal/domain/me"
)

func TestReconstruct(t *testing.T) {
	tests := []struct {
		name     string
		input    me.ReconstructInput
		assertFn func(t *testing.T, entity *me.Me)
	}{
		{
			name: "ok#正常",
			input: me.ReconstructInput{
				DisplayName:    sampleName,
				DisplayNameJa:  sampleNameJa,
				Role:           sampleRole,
				Location:       sampleLocation,
				Likes:          []string{},
				Links:          []me.Link{},
				Skills:         me.Skills{},
				Certifications: []me.Certification{},
				UpdatedAt:      baseTime,
			},
			assertFn: func(t *testing.T, entity *me.Me) {
				t.Helper()
				if entity.DisplayName() != sampleName {
					t.Errorf("entity.DisplayName = %v", entity.DisplayName())
				}
				if entity.DisplayNameJa() != sampleNameJa {
					t.Errorf("entity.DisplayNameJa = %v", entity.DisplayNameJa())
				}
				if entity.Role() != sampleRole {
					t.Errorf("entity.Role() = %v", entity.Role())
				}
				if len(entity.Likes()) != 0 {
					t.Errorf("entity.Links() = %v", entity.Links())
				}
			},
		},
		{
			name: "ok#未設定",
			input: me.ReconstructInput{
				DisplayName:    "",
				DisplayNameJa:  "",
				Role:           "",
				Location:       "",
				Likes:          []string{},
				Links:          []me.Link{},
				Skills:         me.Skills{},
				Certifications: []me.Certification{},
				UpdatedAt:      time.Time{},
			},
			assertFn: func(t *testing.T, entity *me.Me) {
				t.Helper()
				if entity.DisplayName() != "" {
					t.Errorf("entity.DisplayName = %v", entity.DisplayName())
				}
				if entity.DisplayNameJa() != "" {
					t.Errorf("entity.DisplayNameJa = %v", entity.DisplayNameJa())
				}
				if entity.Role() != "" {
					t.Errorf("entity.Role() = %v", entity.Role())
				}
				if len(entity.Likes()) != 0 {
					t.Errorf("entity.Links() = %v", entity.Links())
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := me.Reconstruct(tt.input)

			tt.assertFn(t, got)
		})
	}
}
