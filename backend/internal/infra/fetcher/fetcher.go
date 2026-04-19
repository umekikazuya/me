package fetcher

import (
	"context"
	"fmt"

	app "github.com/umekikazuya/me/internal/app/article"
)

// platformFetcher は platform ごとの取得戦略
type platformFetcher interface {
	fetch(ctx context.Context) ([]app.FetchedArticle, error)
}

// Dispatcher は platform 名で platformFetcher を選択する Registry
type Dispatcher struct {
	fetchers map[string]platformFetcher
}

var _ app.PlatformArticleFetcher = (*Dispatcher)(nil)

func NewDispatcher(fetchers map[string]platformFetcher) *Dispatcher {
	return &Dispatcher{fetchers: fetchers}
}

// NewDefaultDispatcher は qiita / zenn を登録した Dispatcher を返す。
func NewDefaultDispatcher(qiitaToken, zennUsername string) *Dispatcher {
	return NewDispatcher(map[string]platformFetcher{
		"qiita": NewQiitaFetcher(qiitaToken),
		"zenn":  NewZennFetcher(zennUsername),
	})
}

func (d *Dispatcher) Fetch(ctx context.Context, platform string) ([]app.FetchedArticle, error) {
	f, ok := d.fetchers[platform]
	if !ok {
		return nil, fmt.Errorf("unsupported platform: %s", platform)
	}
	return f.fetch(ctx)
}
