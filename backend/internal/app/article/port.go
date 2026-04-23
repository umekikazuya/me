package article

import (
	"context"
)

type PlatformArticleFetcher interface {
	Fetch(
		ctx context.Context,
		platform string,
	) ([]FetchedArticle, error)
}

type Tokenizer interface {
	Tokenize(text string) []string
}
