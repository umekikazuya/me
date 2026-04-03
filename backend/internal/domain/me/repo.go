package me

import "context"

// Repo はMe集約のリポジトリ設計
type Repo interface {
	FindByID(ctx context.Context, id string) (*Me, error)
	Save(ctx context.Context, me *Me) error
	Exists(ctx context.Context, id string) (bool, error)
}
