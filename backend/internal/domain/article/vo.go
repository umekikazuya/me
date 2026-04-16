package article

import (
	"errors"
	"slices"
	"strings"
	"time"
)

type (
	id               struct{ value string }
	title            struct{ value string }
	url              struct{ value string }
	platform         struct{ value string }
	publishedAt      struct{ value time.Time }
	articleUpdatedAt struct{ value time.Time }
	// tags
	// tokens
	isActive struct{ value bool }

	// SearchCriteria は検索パラメータを表現
	SearchCriteria struct {
		Keyword    *string
		Tags       []string
		Year       *int
		Platform   *string
		ActiveOnly bool
		Limit      int
		Cursor     *string
	}
	TagCount struct {
		Name  string
		Count int
	}
	TokenCount struct {
		Value string
		Count int
	}
)

// allowedPlatforms はプラットフォーム許可リスト
var allowedPlatforms = []string{"qiita", "zenn", "mochiya", "note"}

// newID はidオブジェクトを生成
func newID(
	input string,
) (id, error) {
	err := validateNonEmpty(input)
	if err != nil {
		return id{}, err
	}
	return id{
		value: input,
	}, nil
}

// title はtitleオブジェクトを生成
func newTitle(
	input string,
) (title, error) {
	err := validateNonEmpty(input)
	if err != nil {
		return title{}, err
	}
	return title{
		value: input,
	}, nil
}

// newURL はurlオブジェクトを生成
func newURL(
	input string,
) (url, error) {
	err := validateNonEmpty(input)
	if err != nil {
		return url{}, err
	}
	return url{
		value: input,
	}, nil
}

// newPlatform はplatformオブジェクトを生成
func newPlatform(
	input string,
) (platform, error) {
	if !slices.Contains(allowedPlatforms, input) {
		return platform{}, errors.New("許可されたプラットフォームではありません")
	}
	err := validateNonEmpty(input)
	if err != nil {
		return platform{}, err
	}
	return platform{
		value: input,
	}, nil
}

// newPublishedAt はpublishedAtオブジェクトを生成
func newPublishedAt(
	input time.Time,
) (publishedAt, error) {
	if input.After(time.Now()) {
		return publishedAt{}, errors.New("公開日時の指定に不備があります")
	}
	return publishedAt{
		value: input,
	}, nil
}

// newArticleUpdatedAt はarticleUpdatedAtオブジェクトを生成
func newArticleUpdatedAt(
	input time.Time,
) (articleUpdatedAt, error) {
	if input.After(time.Now()) {
		return articleUpdatedAt{}, errors.New("更新日時の指定に不備があります")
	}
	return articleUpdatedAt{
		value: input,
	}, nil
}

// newIsActive はisActiveオブジェクトを生成
func newIsActive(
	input bool,
) (isActive, error) {
	return isActive{
		value: input,
	}, nil
}

// validateNonEmpty
func validateNonEmpty(value string) error {
	if strings.TrimSpace(value) == "" {
		return errors.New("must not be empty")
	}
	return nil
}
