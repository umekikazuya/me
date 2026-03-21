package me

// Repo はMe集約のリポジトリ設計
type Repo interface {
	Find() (*Me, error)
	Save(e *Me) error
	Exists() bool
}
