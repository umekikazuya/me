package article

import (
	"context"
)

type SyncFetcher interface {
	Fetch(ctx context.Context, platform string) error // TODO: FetchArticles を返却する必要がある
}
