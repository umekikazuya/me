package me

import "context"

// Repo はMe集約のリポジトリ設計
type Repo interface {
	Find(ctx context.Context) (*Me, error)
	Save(ctx context.Context, me *Me) error
	Exists(ctx context.Context) (bool, error)
}
