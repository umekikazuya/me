package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	app "github.com/umekikazuya/me/internal/app/article"
)

const (
	zennBaseURL  = "https://zenn.dev/api"
	zennCount    = 96
	zennPlatform = "zenn"
)

type ZennFetcher struct {
	username   string
	httpClient *http.Client
}

func NewZennFetcher(username string) *ZennFetcher {
	return &ZennFetcher{
		username:   username,
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

func (f *ZennFetcher) fetch(ctx context.Context) ([]app.FetchedArticle, error) {
	var all []app.FetchedArticle
	for page := 1; ; page++ {
		items, hasNext, err := f.fetchPage(ctx, page)
		if err != nil {
			return nil, err
		}
		all = append(all, items...)
		if !hasNext {
			break
		}
	}
	return all, nil
}

type zennArticlesResponse struct {
	Articles []zennArticle `json:"articles"`
	NextPage *int          `json:"next_page"`
}

type zennArticle struct {
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	Path        string `json:"path"`
	PublishedAt string `json:"published_at"`
	BodyUpdated string `json:"body_updated_at"`
}

func (f *ZennFetcher) fetchPage(ctx context.Context, page int) ([]app.FetchedArticle, bool, error) {
	url := fmt.Sprintf("%s/articles?username=%s&order=latest&count=%d&page=%d",
		zennBaseURL, f.username, zennCount, page)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, false, fmt.Errorf("zenn: build request: %w", err)
	}

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return nil, false, fmt.Errorf("zenn: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("zenn: unexpected status: %d", resp.StatusCode)
	}

	var body zennArticlesResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, false, fmt.Errorf("zenn: decode response: %w", err)
	}

	result := make([]app.FetchedArticle, 0, len(body.Articles))
	for _, a := range body.Articles {
		articleURL := "https://zenn.dev" + a.Path
		publishedAt, _ := time.Parse(time.RFC3339, a.PublishedAt)
		updatedAt, _ := time.Parse(time.RFC3339, a.BodyUpdated)
		result = append(result, app.FetchedArticle{
			ExternalID:       a.Slug,
			Title:            a.Title,
			URL:              articleURL,
			Platform:         zennPlatform,
			PublishedAt:      publishedAt,
			ArticleUpdatedAt: updatedAt,
		})
	}
	return result, body.NextPage != nil, nil
}
