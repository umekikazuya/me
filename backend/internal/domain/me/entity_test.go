package me_test

import (
	"slices"
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

func TestMe_updateProfile(t *testing.T) {
	tests := []struct {
		name     string
		baseTime time.Time
		in       []me.OptProfileFunc
		wantErr  bool
		assertFn func(t *testing.T, e *me.Me, baseTime time.Time)
	}{
		{
			name: "ok#正常に更新できる",
			in: []me.OptProfileFunc{
				me.OptDisplayNameJa("abc"),
				me.OptLocation("abc"),
				me.OptDisplayName("abc"),
				me.OptRole("abc"),
			},
			baseTime: baseTime,
			wantErr:  false,
			assertFn: func(t *testing.T, e *me.Me, baseTime time.Time) {
				t.Helper()
				if e.DisplayName() != "abc" {
					t.Errorf("e.DisplayName() = %v", e.DisplayName())
				}
				if e.DisplayNameJa() != "abc" {
					t.Errorf("e.profile.displayNameJa = %v", e.DisplayNameJa())
				}
				if e.Role() != "abc" {
					t.Errorf("e.profile.role = %v", e.Role())
				}
				if e.Location() != "abc" {
					t.Errorf("e.profile.location = %v", e.Location())
				}
			},
		},
		{
			name:     "ng#Optが空",
			baseTime: baseTime,
			in:       []me.OptProfileFunc{},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &me.Me{}
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
			if !e.UpdatedAt().Equal(baseTime) {
				t.Errorf("e.UpdatedAt = %v, baseTime = %v", e.UpdatedAt(), tt.baseTime)
			}
			tt.assertFn(t, e, tt.baseTime)
		})
	}
}

func TestMe_UpdateLikes(t *testing.T) {
	tests := []struct {
		name     string
		in       []string
		baseTime time.Time
		wantErr  bool
		assertFn func(t *testing.T, e *me.Me)
	}{
		{
			name:     "ok#正常に更新できる",
			in:       []string{"go", "rust"},
			baseTime: baseTime,
			wantErr:  false,
			assertFn: func(t *testing.T, e *me.Me) {
				t.Helper()
				if !slices.Equal(e.Likes(), []string{"go", "rust"}) {
					t.Errorf("e.Likes() = %v", e.Likes())
				}
				if !e.UpdatedAt().Equal(baseTime) {
					t.Errorf("e.updatedAt = %v, want = %v", e.UpdatedAt(), baseTime)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &me.Me{}
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

// func TestMe_UpdateLinks(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		in       []me.Link
// 		wantErr  bool
// 		assertFn func(t *testing.T, e *me.Me)
// 	}{
// 		{
// 			name: "ok#正常に更新できる",
// 			in: []me.Link{
// 				{
// 					platform: "a",
// 					url:      "example.com",
// 				},
// 			},
// 			wantErr: false,
// 			assertFn: func(t *testing.T, e *Me) {
// 				t.Helper()
// 				links := e.links
// 				l := links[0]
// 				if l.platform != "a" {
// 					t.Errorf("l.platform = %v, want = %v", l.platform, "a")
// 				}
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			e := &me.Me{}
// 			gotErr := e.UpdateLinks(tt.in, baseTime)
// 			if gotErr != nil {
// 				if !tt.wantErr {
// 					t.Errorf("UpdateLinks() failed: %v", gotErr)
// 				}
// 				return
// 			}
// 			if tt.wantErr {
// 				t.Fatal("UpdateLinks() succeeded unexpectedly")
// 			}
// 			tt.assertFn(t, e)
// 		})
// 	}
// }
//
// func TestMe_AddSkill(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		item     string
// 		parent   string
// 		seedFn   func(t *testing.T, e *me.Me)
// 		wantErr  bool
// 		assertFn func(t *testing.T, e *me.Me)
// 	}{
// 		{
// 			name:   "ok#正常に追加出来る",
// 			item:   skillItemNameA,
// 			parent: skillCategoryNameA,
// 			seedFn: func(t *testing.T, e *me.Me) {
// 				t.Helper()
// 				e.Skills() = Skills{
// 					skillCategoryNameA: {
// 						Items: []string{"def"},
// 					},
// 				}
// 			},
// 			wantErr: false,
// 			assertFn: func(t *testing.T, e *Me) {
// 				t.Helper()
// 				if len(e.skills[skillCategoryNameA].Items) != 2 {
// 					t.Errorf("len(e.skills) = %d", 2)
// 				}
// 				if !slices.Contains(e.skills[skillCategoryNameA].Items, skillItemNameA) {
// 					t.Fatalf("e.skills = %v", e.skills)
// 				}
// 			},
// 		},
// 		{
// 			name:   "ok#正常に追加できる(カテゴリ含め)",
// 			item:   skillItemNameA,
// 			parent: skillCategoryNameA,
// 			seedFn: func(t *testing.T, e *Me) {
// 				t.Helper()
// 			},
// 			wantErr: false,
// 			assertFn: func(t *testing.T, e *Me) {
// 				t.Helper()
// 				if !slices.Contains(e.skills[skillCategoryNameA].Items, skillItemNameA) {
// 					t.Fatalf("e.skills = %v", e.skills)
// 				}
// 			},
// 		},
// 		{
// 			name:   "ng#既に登録済み",
// 			item:   skillItemNameA,
// 			parent: skillCategoryNameA,
// 			seedFn: func(t *testing.T, e *Me) {
// 				t.Helper()
// 				e.skills = Skills{
// 					skillCategoryNameA: {
// 						Items: []string{skillItemNameA},
// 					},
// 				}
// 			},
// 			wantErr:  true,
// 			assertFn: func(t *testing.T, e *Me) {},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Arrange
// 			e, err := NewMe(targetID.String())
// 			if err != nil {
// 				t.Fatalf("could not construct receiver type: %v", err)
// 			}
// 			tt.seedFn(t, e)
// 			// Act
// 			gotErr := e.AddSkill(tt.item, tt.parent, baseTime)
// 			if gotErr != nil {
// 				if !tt.wantErr {
// 					t.Errorf("AddSkill() failed: %v", gotErr)
// 				}
// 				return
// 			}
// 			if tt.wantErr {
// 				t.Fatal("AddSkill() succeeded unexpectedly")
// 			}
// 			// Assert
// 			tt.assertFn(t, e)
// 		})
// 	}
// }
//
// func TestMe_RemoveSkill(t *testing.T) {
// 	tests := []struct {
// 		name         string
// 		itemName     string
// 		categoryName string
// 		baseTime     time.Time
// 		seedFn       func(t *testing.T, e *Me)
// 		wantErr      bool
// 		assertFn     func(t *testing.T, e *Me)
// 	}{
// 		{
// 			name:         "ok#正常に削除出来る(カテゴリごと)",
// 			itemName:     skillItemNameA,
// 			categoryName: skillCategoryNameA,
// 			baseTime:     baseTime.Add(2 * time.Hour),
// 			seedFn: func(t *testing.T, e *Me) {
// 				t.Helper()
// 				e.skills = Skills{
// 					skillCategoryNameA: {
// 						Items: []string{skillItemNameA},
// 					},
// 				}
// 			},
// 			wantErr: false,
// 			assertFn: func(t *testing.T, e *Me) {
// 				t.Helper()
// 				if len(e.skills) != 0 {
// 					t.Errorf("e.skills = %v", e.skills)
// 				}
// 			},
// 		},
// 		{
// 			name:         "ng#存在しない",
// 			itemName:     skillItemNameA,
// 			categoryName: skillCategoryNameA,
// 			baseTime:     baseTime.Add(2 * time.Hour),
// 			seedFn: func(t *testing.T, e *Me) {
// 				t.Helper()
// 			},
// 			wantErr: true,
// 			assertFn: func(t *testing.T, e *Me) {
// 				t.Helper()
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Arrange
// 			e, err := NewMe(targetID.String())
// 			if err != nil {
// 				t.Fatalf("could not construct receiver type: %v", err)
// 			}
// 			tt.seedFn(t, e)
// 			// Act
// 			gotErr := e.RemoveSkill(tt.itemName, tt.categoryName, tt.baseTime)
// 			if gotErr != nil {
// 				if !tt.wantErr {
// 					t.Errorf("RemoveSkill() failed: %v", gotErr)
// 				}
// 				return
// 			}
// 			if tt.wantErr {
// 				t.Fatal("RemoveSkill() succeeded unexpectedly")
// 			}
// 			// Assert
// 			tt.assertFn(t, e)
// 		})
// 	}
// }
