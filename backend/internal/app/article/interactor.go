package article

import (
	"context"
	"fmt"
	"sort"
	"strings"

	domain "github.com/umekikazuya/me/internal/domain/article"
	"github.com/umekikazuya/me/pkg/errs"
)

var _ Interactor = (*interactor)(nil)

type Interactor interface {
	Search(ctx context.Context, input InputSearchDto) (*OutputSearchDto, error)
	GetTagsAll(ctx context.Context) (*OutputTagAllDto, error)
	GetSuggests(ctx context.Context, input InputGetSuggestDto) (*OutputGetSuggestAllDto, error)
	Register(ctx context.Context, input InputRegisterDto) error
	Update(ctx context.Context, input InputUpdateDto) error
	Remove(ctx context.Context, input InputRemoveDto) error
	Sync(ctx context.Context, platform string) domain.IndexingResult
}

type interactor struct {
	repo      domain.Repo
	fetcher   PlatformArticleFetcher
	tokenizer Tokenizer
}

// NewInteractor はユースケースの初期化クラス
func NewInteractor(
	repo domain.Repo,
	fetcher PlatformArticleFetcher,
	tokenizer Tokenizer,
) Interactor {
	return &interactor{
		repo:      repo,
		fetcher:   fetcher,
		tokenizer: tokenizer,
	}
}

func (i *interactor) Search(ctx context.Context, input InputSearchDto) (*OutputSearchDto, error) {
	criteria := domain.SearchCriteria{
		Tags:       input.Tag,
		Year:       input.Year,
		Platform:   input.Platform,
		ActiveOnly: true,
		Limit:      input.Limit,
		Cursor:     input.NextCursor,
	}
	tokens := i.tokenizer.Tokenize(*input.Q)
	if strings.TrimSpace(*input.Q) != "" && len(tokens) == 0 {
		return &OutputSearchDto{Articles: []OutputArticleItemDto{}}, nil
	}
	criteria.Tokens = tokens
	result, err := i.repo.FindAll(ctx, criteria)
	if err != nil {
		return nil, err
	}

	items := make([]OutputArticleItemDto, 0, len(result.Articles))
	for _, a := range result.Articles {
		items = append(items, OutputArticleItemDto{
			ExternalID:  a.ID(),
			Title:       a.Title(),
			URL:         a.URL(),
			Platform:    a.Platform(),
			PublishedAt: a.PublishedAt(),
			Tags:        a.Tags(),
		})
	}

	out := &OutputSearchDto{Articles: items}
	if result.NextCursor != nil {
		out.NextCursor = *result.NextCursor
	}
	return out, nil
}

func (i *interactor) GetTagsAll(ctx context.Context) (*OutputTagAllDto, error) {
	tags, err := i.repo.AllTags(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]OutputTagItemDto, 0, len(tags))
	for _, t := range tags {
		items = append(items, OutputTagItemDto{Name: t.Name, Count: t.Count})
	}
	return &OutputTagAllDto{Tags: items}, nil
}

func (i *interactor) GetSuggests(ctx context.Context, input InputGetSuggestDto) (*OutputGetSuggestAllDto, error) {
	tags, err := i.repo.AllTags(ctx)
	if err != nil {
		return nil, err
	}
	tokens, err := i.repo.AllTokens(ctx)
	if err != nil {
		return nil, err
	}

	var suggestions []OutputGetSuggestItemDto
	for _, t := range tags {
		if strings.HasPrefix(t.Name, input.Q) {
			suggestions = append(suggestions, OutputGetSuggestItemDto{Type: "tag", Value: t.Name, Count: t.Count})
		}
	}
	for _, t := range tokens {
		if strings.HasPrefix(t.Value, input.Q) {
			suggestions = append(suggestions, OutputGetSuggestItemDto{Type: "token", Value: t.Value, Count: t.Count})
		}
	}

	sort.Slice(suggestions, func(a, b int) bool {
		return suggestions[a].Count > suggestions[b].Count
	})

	// TODO: タイトルサジェストを追加する
	// input.Q をトークンとして完全一致する記事を返す（トークン逆引き）
	// 実装時は以下が必要:
	//   - domain.Repo に FindByToken(ctx, token string) メソッドを追加（GSI exact match）
	//   - OutputGetSuggestItemDto に ExternalID string フィールドを追加
	//   - 結果は publishedAt 降順でソートし type="title" として追加
	if suggestions == nil {
		suggestions = []OutputGetSuggestItemDto{}
	}
	return &OutputGetSuggestAllDto{Suggestions: suggestions}, nil
}

func (i *interactor) Register(ctx context.Context, input InputRegisterDto) error {
	existing, err := i.repo.FindByExternalID(ctx, input.ExternalID)
	if err != nil {
		return err
	}
	if existing != nil {
		return fmt.Errorf("article already exists: %s: %w", input.ExternalID, errs.ErrConflict)
	}

	opts := []domain.Opt{
		domain.WithTags(input.Tags),
		domain.WithTokens(i.tokenizer.Tokenize(input.Title)),
	}
	if !input.PublishedAt.IsZero() {
		opts = append(opts, domain.WithPublishedAt(input.PublishedAt))
	}
	if !input.ArticleUpdatedAt.IsZero() {
		opts = append(opts, domain.WithArticleUpdatedAt(input.ArticleUpdatedAt))
	}

	article, err := domain.Register(input.ExternalID, input.Title, input.URL, input.Platform, opts...)
	if err != nil {
		return fmt.Errorf("%s: %w", err.Error(), errs.ErrUnprocessable)
	}
	return i.repo.Save(ctx, article)
}

func (i *interactor) Update(ctx context.Context, input InputUpdateDto) error {
	article, err := i.repo.FindByExternalID(ctx, input.ExternalID)
	if err != nil {
		return err
	}
	if article == nil {
		return fmt.Errorf("article not found: %s: %w", input.ExternalID, errs.ErrNotFound)
	}

	// TODO: タイトル変更時にトークンを再生成する（現状 Sync/Index ルートのみ対応）
	opts := []domain.Opt{
		domain.WithTags(input.Tags),
	}
	if !input.PublishedAt.IsZero() {
		opts = append(opts, domain.WithPublishedAt(input.PublishedAt))
	}
	if !input.ArticleUpdatedAt.IsZero() {
		opts = append(opts, domain.WithArticleUpdatedAt(input.ArticleUpdatedAt))
	}

	if err := article.Update(input.Title, input.URL, opts...); err != nil {
		return fmt.Errorf("%s: %w", err.Error(), errs.ErrUnprocessable)
	}
	return i.repo.Save(ctx, article)
}

func (i *interactor) Remove(ctx context.Context, input InputRemoveDto) error {
	article, err := i.repo.FindByExternalID(ctx, input.ExternalID)
	if err != nil {
		return err
	}
	if article == nil {
		return fmt.Errorf("article not found: %s: %w", input.ExternalID, errs.ErrNotFound)
	}

	if err := article.Remove(); err != nil {
		return fmt.Errorf("%s: %w", err.Error(), errs.ErrUnprocessable)
	}
	return i.repo.Save(ctx, article)
}

func (i *interactor) Sync(ctx context.Context, platform string) domain.IndexingResult {
	var result domain.IndexingResult

	fetched, err := i.fetcher.Fetch(ctx, platform)
	if err != nil {
		result.Errors = append(result.Errors, err)
		return result
	}

	existing, err := i.repo.FindByPlatform(ctx, platform)
	if err != nil {
		result.Errors = append(result.Errors, err)
		return result
	}

	existingByID := make(map[string]*domain.Article, len(existing))
	for idx := range existing {
		a := existing[idx]
		existingByID[a.ID()] = &a
	}

	fetchedIDs := make(map[string]struct{}, len(fetched))
	for _, f := range fetched {
		fetchedIDs[f.ExternalID] = struct{}{}

		if article, ok := existingByID[f.ExternalID]; ok {
			opts := []domain.Opt{
				domain.WithTags(f.Tags),
				domain.WithTokens(i.tokenizer.Tokenize(strings.Join([]string{f.Title, f.Body}, " "))),
			}
			if !f.PublishedAt.IsZero() {
				opts = append(opts, domain.WithPublishedAt(f.PublishedAt))
			}
			if !f.ArticleUpdatedAt.IsZero() {
				opts = append(opts, domain.WithArticleUpdatedAt(f.ArticleUpdatedAt))
			}
			if err := article.Reindex(f.Title, f.URL, opts...); err != nil {
				result.Errors = append(result.Errors, err)
				continue
			}
			if err := i.repo.Save(ctx, article); err != nil {
				result.Errors = append(result.Errors, err)
				continue
			}
			result.Reindexed++
		} else {
			opts := []domain.Opt{
				domain.WithTags(f.Tags),
				domain.WithTokens(i.tokenizer.Tokenize(strings.Join([]string{f.Title, f.Body}, " "))),
			}
			if !f.PublishedAt.IsZero() {
				opts = append(opts, domain.WithPublishedAt(f.PublishedAt))
			}
			if !f.ArticleUpdatedAt.IsZero() {
				opts = append(opts, domain.WithArticleUpdatedAt(f.ArticleUpdatedAt))
			}
			article, err := domain.Index(f.ExternalID, f.Title, f.URL, f.Platform, opts...)
			if err != nil {
				result.Errors = append(result.Errors, err)
				continue
			}
			if err := i.repo.Save(ctx, article); err != nil {
				result.Errors = append(result.Errors, err)
				continue
			}
			result.Indexed++
		}
	}

	for _, article := range existingByID {
		if _, ok := fetchedIDs[article.ID()]; ok {
			continue
		}
		if err := article.Deactivate(); err != nil {
			result.Errors = append(result.Errors, err)
			continue
		}
		if err := i.repo.Save(ctx, article); err != nil {
			result.Errors = append(result.Errors, err)
			continue
		}
		result.Deactivated++
	}

	return result
}
