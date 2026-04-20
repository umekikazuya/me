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

// NewDefaultDispatcher は設定済みの platform のみ登録した Dispatcher を返す。
// トークン・ユーザー名が空の platform は登録しない。
func NewDefaultDispatcher(qiitaToken, zennUsername string) *Dispatcher {
	fetchers := make(map[string]platformFetcher)
	if qiitaToken != "" {
		fetchers["qiita"] = NewQiitaFetcher(qiitaToken)
	}
	if zennUsername != "" {
		fetchers["zenn"] = NewZennFetcher(zennUsername)
	}
	return NewDispatcher(fetchers)
}

func (d *Dispatcher) Fetch(ctx context.Context, platform string) ([]app.FetchedArticle, error) {
	f, ok := d.fetchers[platform]
	if !ok {
		return nil, fmt.Errorf("unsupported platform: %s", platform)
	}
	return f.fetch(ctx)
}
