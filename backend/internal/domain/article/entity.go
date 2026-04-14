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
	tags             []string
	tokens           []string
	publishedAt      publishedAt
	articleUpdatedAt articleUpdatedAt
	isActive         isActive
	createdAt        time.Time
	updatedAt        time.Time
}

// Opt はFOを表現
type Opt func(*Article) error

// ReconstructArticleInput はReconstructIdentityの入力型
type ReconstructArticleInput struct {
	ID               string
	Title            string
	URL              string
	Platform         string
	Tags             []string
	Tokens           []string
	PublishedAt      time.Time
	ArticleUpdatedAt time.Time
	IsActive         bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// ファクトリー関数

// newArticle はArticle集約のファクトリー関数
func newArticle(
	inputID, inputTitle, inputURL, inputPlatform string,
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
	url, err := newURL(inputURL)
	if err != nil {
		return nil, err
	}
	platform, err := newPlatform(inputPlatform)
	if err != nil {
		return nil, err
	}
	isActive, err := newIsActive(true)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	e := Article{
		id:        id,
		title:     title,
		url:       url,
		platform:  platform,
		isActive:  isActive,
		createdAt: now,
		updatedAt: now,
	}
	for _, opt := range opts {
		if opt == nil {
			return nil, errors.New("FunctionalOptionパターンの指定にミスがあります")
		}
		if err := opt(&e); err != nil {
			return nil, err
		}
	}
	return &e, nil
}

func Reconstruct(input ReconstructArticleInput) (*Article, error) {
	return &Article{
		id:               id{value: input.ID},
		title:            title{value: input.Title},
		url:              url{value: input.URL},
		platform:         platform{value: input.Platform},
		tags:             input.Tags,
		tokens:           input.Tokens,
		publishedAt:      publishedAt{value: input.PublishedAt},
		articleUpdatedAt: articleUpdatedAt{value: input.ArticleUpdatedAt},
		isActive:         isActive{value: input.IsActive},
		createdAt:        input.CreatedAt,
		updatedAt:        input.UpdatedAt,
	}, nil
}

// --- FO ---

func WithTags(inputs []string) Opt {
	return func(e *Article) error {
		e.tags = inputs
		return nil
	}
}

func WithTokens(inputs []string) Opt {
	return func(e *Article) error {
		e.tokens = inputs
		return nil
	}
}

func WithPublishedAt(input time.Time) Opt {
	return func(e *Article) error {
		value, err := newPublishedAt(input)
		if err != nil {
			return err
		}
		e.publishedAt = value
		return nil
	}
}

func WithArticleUpdatedAt(input time.Time) Opt {
	return func(e *Article) error {
		value, err := newArticleUpdatedAt(input)
		if err != nil {
			return err
		}
		e.articleUpdatedAt = value
		return nil
	}
}

// --- 振る舞い ---

// Index は記事のインデックス登録を実施

func Index(
	inputID, inputTitle, inputURL, inputPlatform string,
	opts ...Opt,
) (*Article, error) {
	return newArticle(
		inputID,
		inputTitle,
		inputURL,
		inputPlatform,
		opts...,
	)
}

// Register は記事の手動登録を実施
func Register(
	inputID, inputTitle, inputURL, inputPlatform string,
	opts ...Opt,
) (*Article, error) {
	return newArticle(
		inputID,
		inputTitle,
		inputURL,
		inputPlatform,
		opts...,
	)
}

// Reindex はインデクサーによる全上書き更新を実施
//
// id, platformの上書きは不可
func (e *Article) Reindex(
	inputTitle, inputURL string,
	opts ...Opt,
) error {
	if !e.IsActive() {
		return errors.New("コンテンツが削除されています")
	}
	title, err := newTitle(inputTitle)
	if err != nil {
		return err
	}
	url, err := newURL(inputURL)
	if err != nil {
		return err
	}
	e.title = title
	e.url = url
	e.updatedAt = time.Now()
	return nil
}

// Update は手動上書き更新
//
// id, platformの上書きは不可
func (e *Article) Update(
	inputTitle, inputURL string,
	opts ...Opt,
) error {
	if !e.IsActive() {
		return errors.New("コンテンツが削除されています")
	}
	title, err := newTitle(inputTitle)
	if err != nil {
		return err
	}
	url, err := newURL(inputURL)
	if err != nil {
		return err
	}
	e.title = title
	e.url = url
	e.updatedAt = time.Now()
	return nil
}

// Deactivate はインデクサーによる論理削除
func (e *Article) Deactivate() error {
	isActive, err := newIsActive(false)
	if err != nil {
		return err
	}
	e.isActive = isActive
	e.updatedAt = time.Now()
	return nil
}

// Remove は手動論理削除
func (e *Article) Remove() error {
	isActive, err := newIsActive(false)
	if err != nil {
		return err
	}
	e.isActive = isActive
	e.updatedAt = time.Now()
	return nil
}

// --- Getter ---

func (e *Article) ID() string {
	return e.id.value
}

func (e *Article) Title() string {
	return e.title.value
}

func (e *Article) URL() string {
	return e.url.value
}

func (e *Article) Platform() string {
	return e.platform.value
}

func (e *Article) IsActive() bool {
	return e.isActive.value
}

func (e *Article) Tags() []string {
	return e.tags
}

func (e *Article) PublishedAt() time.Time {
	return e.publishedAt.value
}

func (e *Article) ArticleUpdatedAt() time.Time {
	return e.articleUpdatedAt.value
}

func (e *Article) CreatedAt() time.Time {
	return e.createdAt
}

func (e *Article) UpdatedAt() time.Time {
	return e.updatedAt
}
