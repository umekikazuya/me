package me

import (
	domain "github.com/umekikazuya/me/internal/domain/me"
)

// toOutputDto はドメインをレスポンス形式に変換する関数
func toOutputDto(e domain.Me) *OutputDto {
	return &OutputDto{
		Likes:       e.Likes(),
		Location:    e.Location(),
		DisplayName: e.DisplayName(),
		DisplayJa:   e.DisplayNameJa(),
		Role:        e.Role(),
		CreatedAt:   e.CreatedAt(),
		UpdatedAt:   e.UpdatedAt(),
	}
}
