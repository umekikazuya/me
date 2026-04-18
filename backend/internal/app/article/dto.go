package article

import (
	"time"
)

type OutputArticleItemDto struct {
	ExternalID  string    `json:"externalId"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	Platform    string    `json:"platform"`
	PublishedAt time.Time `json:"publishedAt"`
	Tags        []string  `json:"tags"`
}

// InputSearchDto は検索パラメータ(クエリパラメータ想定)
type InputSearchDto struct {
	Q          *string  `json:"q"` // ","でAND検索
	Tag        []string `json:"tag"`
	Year       *int     `json:"year"`
	Platform   *string  `json:"platform"`
	Limit      int      `json:"limit" validate:"required,min=1,max=100"`
	NextCursor *string  `json:"cursor"`
}

// OutputSearchDto はArticle検索結果を表現
type OutputSearchDto struct {
	Articles   []OutputArticleItemDto `json:"articles"`
	NextCursor string                 `json:"nextCursor,omitempty"`
}

type OutputTagItemDto struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type OutputTagAllDto struct {
	Tags []OutputTagItemDto `json:"tags"`
}

type InputGetSuggestDto struct {
	Q string `json:"q" validate:"q,required,min=1"`
}

type OutputGetSuggestItemDto struct {
	Value string `json:"value"`
	Count int    `json:"count"`
}

type OutputGetSuggestAllDto struct {
	Suggests []OutputGetSuggestItemDto `json:"suggests"`
}

type InputRegisterDto struct {
	ExternalID       string    `json:"externalId" validate:"externalId,required,min=1,max=256"`
	Title            string    `json:"title" validate:"title,required,min=1,max=500"`
	URL              string    `json:"url" validate:"url,required,url"`
	Platform         string    `json:"platform" validate:"required"`
	PublishedAt      time.Time `json:"publishedAt" validate:""`
	ArticleUpdatedAt time.Time `json:"articleUpdatedAt" validate:""`
	Tags             []string  `json:"tags"`
}

type InputUpdateDto struct {
	ExternalID       string    `json:"-"`
	Title            string    `json:"title" validate:"title,required,min=1,max=500"`
	URL              string    `json:"url" validate:"url,required,url"`
	PublishedAt      time.Time `json:"publishedAt" validate:""`
	ArticleUpdatedAt time.Time `json:"articleUpdatedAt" validate:""`
	Tags             []string  `json:"tags" validate:"min=1,max=1"`
}

type InputRemoveDto struct {
	ExternalID string `json:"-"`
}
