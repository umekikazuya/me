package me

import (
	domain "github.com/umekikazuya/me/internal/domain/me"
)

// toOutputDto はドメインをレスポンス形式に変換する関数
func toOutputDto(e domain.Me) *OutputDto {
	links := make([]struct {
		Platform string `json:"platform"`
		URL      string `json:"url"`
	}, 0, len(e.Links()))
	for _, l := range e.Links() {
		links = append(links, struct {
			Platform string `json:"platform"`
			URL      string `json:"url"`
		}{
			Platform: l.Platform(),
			URL:      l.URL(),
		})
	}

	return &OutputDto{
		Likes:       e.Likes(),
		Links:       links,
		Location:    e.Location(),
		DisplayName: e.DisplayName(),
		DisplayJa:   e.DisplayNameJa(),
		Role:        e.Role(),
		CreatedAt:   e.CreatedAt(),
		UpdatedAt:   e.UpdatedAt(),
	}
}
