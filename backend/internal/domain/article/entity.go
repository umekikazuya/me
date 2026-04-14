package article

import (
	"errors"
	"time"
)

// 型定義

// Article はArticle集約を表現
type Article struct {
	id               id
	title            title
	url              url
	platform         platform
	publishedAt      publishedAt
	articleUpdatedAt articleUpdatedAt
	isActive         isActive
	createdAt        time.Time
	updatedAt        time.Time
}

// Opt はFOを表現
type Opt func(*Article) error

// ファクトリー関数

// newArticle はArticle集約のファクトリー関数
func newArticle(
	inputID, inputTitle, inputURL, inputPlatform string,
	inputIsActive bool,
	opts ...Opt,
) (*Article, error) {
	id, err := newID(inputID)
	if err != nil {
		return nil, err
	}
	title, err := newTitle(inputTitle)
	if err != nil {
		return nil, err
	}
	platform, err := newPlatform(inputPlatform)
	if err != nil {
		return nil, err
	}
	isActive, err := newIsActive(inputIsActive)
	if err != nil {
		return nil, err
	}
	e := Article{
		id:       id,
		title:    title,
		platform: platform,
		isActive: isActive,
	}
	for _, opt := range opts {
		if opt == nil {
			return nil, errors.New("nil option is not allowed")
		}
		if err := opt(&e); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

// --- FO ---

func WithTags(inputs []string) Opt {
	return nil
}

func WithTokens(inputs []string) Opt {
	return nil
}

func WithPublishedAt(input time.Time) Opt {
	return nil
}

func WithArticleUpdatedAt(input time.Time) Opt {
	return nil
}

// --- 振る舞い ---

// Index は記事のインデックス登録を実施

func Index(
	inputID, inputTitle, inputURL, inputPlatform string,
	opts ...Opt,
) (*Article, error) {
	return nil, nil
}

// Register は記事の手動登録を実施
func Register(
	inputID, inputTitle, inputURL, inputPlatform string,
	opts ...Opt,
) (*Article, error) {
	return nil, nil
}

// Reindex はインデクサーによる全上書き更新を実施
//
// id, platformの上書きは不可
func (e *Article) Reindex(
	inputTitle, inputURL string,
	opts ...Opt,
) error {
	return nil
}

// Update は手動上書き更新
//
// id, platformの上書きは不可
func (e *Article) Update(
	inputTitle, inputURL string,
	opts ...Opt,
) error {
	return nil
}

// Deactivate はインデクサーによる論理削除
func (e *Article) Deactivate() error {
	return nil
}

// Remove は手動論理削除
func (e *Article) Remove() error {
	return nil
}
