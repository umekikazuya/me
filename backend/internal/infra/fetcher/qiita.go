package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	app "github.com/umekikazuya/me/internal/app/article"
)

const (
	qiitaBaseURL  = "https://qiita.com/api/v2"
	qiitaPerPage  = 100
	qiitaPlatform = "qiita"
)

type QiitaFetcher struct {
	token      string
	httpClient *http.Client
}

func NewQiitaFetcher(token string) *QiitaFetcher {
	return &QiitaFetcher{
		token:      token,
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

func (f *QiitaFetcher) fetch(ctx context.Context) ([]app.FetchedArticle, error) {
	var all []app.FetchedArticle
	for page := 1; ; page++ {
		items, err := f.fetchPage(ctx, page)
		if err != nil {
			return nil, err
		}
		all = append(all, items...)
		if len(items) < qiitaPerPage {
			break
		}
	}
	return all, nil
}

type qiitaItem struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	URL       string `json:"url"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Private   bool   `json:"private"`
	Tags      []struct {
		Name string `json:"name"`
	} `json:"tags"`
}

func (f *QiitaFetcher) fetchPage(ctx context.Context, page int) ([]app.FetchedArticle, error) {
	url := fmt.Sprintf("%s/authenticated_user/items?page=%s&per_page=%s",
		qiitaBaseURL, strconv.Itoa(page), strconv.Itoa(qiitaPerPage))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("qiita: build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+f.token)

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("qiita: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("qiita: unexpected status: %d", resp.StatusCode)
	}

	var items []qiitaItem
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, fmt.Errorf("qiita: decode response: %w", err)
	}

	result := make([]app.FetchedArticle, 0, len(items))
	for _, item := range items {
		if item.Private {
			continue
		}
		tags := make([]string, 0, len(item.Tags))
		for _, t := range item.Tags {
			tags = append(tags, t.Name)
		}
		publishedAt, _ := time.Parse(time.RFC3339, item.CreatedAt)
		updatedAt, _ := time.Parse(time.RFC3339, item.UpdatedAt)
		result = append(result, app.FetchedArticle{
			ExternalID:       item.ID,
			Title:            item.Title,
			URL:              item.URL,
			Platform:         qiitaPlatform,
			PublishedAt:      publishedAt,
			ArticleUpdatedAt: updatedAt,
			Tags:             tags,
			Body:             item.Body,
		})
	}
	return result, nil
}
