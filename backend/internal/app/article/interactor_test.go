package article

import (
	"context"
	"errors"
	"testing"
	"time"

	domain "github.com/umekikazuya/me/internal/domain/article"
)

// -----------------------------------------------------------------------
// スタブ
// -----------------------------------------------------------------------

type stubRepo struct {
	findAll          func(context.Context, domain.SearchCriteria) (domain.FindAllResult, error)
	findByExternalID func(context.Context, string) (*domain.Article, error)
	findByPlatform   func(context.Context, string) ([]domain.Article, error)
	save             func(context.Context, *domain.Article) error
	allTags          func(context.Context) ([]domain.TagCount, error)
	allTokens        func(context.Context) ([]domain.TokenCount, error)
}

func (s *stubRepo) FindAll(ctx context.Context, c domain.SearchCriteria) (domain.FindAllResult, error) {
	if s.findAll != nil {
		return s.findAll(ctx, c)
	}
	return domain.FindAllResult{}, nil
}

func (s *stubRepo) FindByExternalID(ctx context.Context, id string) (*domain.Article, error) {
	if s.findByExternalID != nil {
		return s.findByExternalID(ctx, id)
	}
	return nil, nil
}

func (s *stubRepo) FindByPlatform(ctx context.Context, platform string) ([]domain.Article, error) {
	if s.findByPlatform != nil {
		return s.findByPlatform(ctx, platform)
	}
	return nil, nil
}

func (s *stubRepo) Save(ctx context.Context, a *domain.Article) error {
	if s.save != nil {
		return s.save(ctx, a)
	}
	return nil
}

func (s *stubRepo) AllTags(ctx context.Context) ([]domain.TagCount, error) {
	if s.allTags != nil {
		return s.allTags(ctx)
	}
	return nil, nil
}

func (s *stubRepo) AllTokens(ctx context.Context) ([]domain.TokenCount, error) {
	if s.allTokens != nil {
		return s.allTokens(ctx)
	}
	return nil, nil
}

type stubFetcher struct {
	fetch func(context.Context, string) ([]FetchedArticle, error)
}

func (s *stubFetcher) Fetch(ctx context.Context, platform string) ([]FetchedArticle, error) {
	if s.fetch != nil {
		return s.fetch(ctx, platform)
	}
	return nil, nil
}

type stubTokenizer struct{}

func (s *stubTokenizer) Tokenize(text string) []string {
	return []string{text}
}

func newInteractor(repo domain.Repo) Interactor {
	return NewInteractor(repo, &stubFetcher{}, &stubTokenizer{})
}

func makeArticle(id, title, platform string, active bool) *domain.Article {
	a, _ := domain.Reconstruct(domain.ReconstructArticleInput{
		ID:        id,
		Title:     title,
		URL:       "https://example.com/" + id,
		Platform:  platform,
		IsActive:  active,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	return a
}

// -----------------------------------------------------------------------
// GetTagsAll
// -----------------------------------------------------------------------

func TestInteractor_GetTagsAll(t *testing.T) {
	t.Run("returns tags mapped to DTO", func(t *testing.T) {
		repo := &stubRepo{
			allTags: func(_ context.Context) ([]domain.TagCount, error) {
				return []domain.TagCount{
					{Name: "go", Count: 5},
					{Name: "design", Count: 2},
				}, nil
			},
		}
		got, err := newInteractor(repo).GetTagsAll(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got.Tags) != 2 {
			t.Fatalf("len(Tags) = %d, want 2", len(got.Tags))
		}
		if got.Tags[0].Name != "go" || got.Tags[0].Count != 5 {
			t.Errorf("Tags[0] = %+v, want {go 5}", got.Tags[0])
		}
		if got.Tags[1].Name != "design" || got.Tags[1].Count != 2 {
			t.Errorf("Tags[1] = %+v, want {design 2}", got.Tags[1])
		}
	})

	t.Run("propagates repo error", func(t *testing.T) {
		repoErr := errors.New("db error")
		repo := &stubRepo{
			allTags: func(_ context.Context) ([]domain.TagCount, error) {
				return nil, repoErr
			},
		}
		_, err := newInteractor(repo).GetTagsAll(context.Background())
		if !errors.Is(err, repoErr) {
			t.Errorf("err = %v, want %v", err, repoErr)
		}
	})
}

// -----------------------------------------------------------------------
// GetSuggests
// -----------------------------------------------------------------------

func TestInteractor_GetSuggests(t *testing.T) {
	t.Run("returns tags and tokens matching prefix, sorted by count desc", func(t *testing.T) {
		repo := &stubRepo{
			allTags: func(_ context.Context) ([]domain.TagCount, error) {
				return []domain.TagCount{
					{Name: "golang", Count: 3},
					{Name: "go", Count: 7},
					{Name: "design", Count: 5}, // prefix 不一致
				}, nil
			},
			allTokens: func(_ context.Context) ([]domain.TokenCount, error) {
				return []domain.TokenCount{
					{Value: "goroutine", Count: 2},
					{Value: "gorm", Count: 4},
					{Value: "設計", Count: 1}, // prefix 不一致
				}, nil
			},
		}

		got, err := newInteractor(repo).GetSuggests(context.Background(), InputGetSuggestDto{Q: "go"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// "design"・"設計" は含まない
		for _, s := range got.Suggestions {
			if s.Value == "design" || s.Value == "設計" {
				t.Errorf("unexpected suggest: %s", s.Value)
			}
		}

		// "go"(7), "gorm"(4), "golang"(3), "goroutine"(2) の順
		wantOrder := []string{"go", "gorm", "golang", "goroutine"}
		if len(got.Suggestions) != len(wantOrder) {
			t.Fatalf("len(Suggests) = %d, want %d", len(got.Suggestions), len(wantOrder))
		}
		for i, want := range wantOrder {
			if got.Suggestions[i].Value != want {
				t.Errorf("Suggests[%d].Value = %q, want %q", i, got.Suggestions[i].Value, want)
			}
		}
	})

	t.Run("returns empty when nothing matches prefix", func(t *testing.T) {
		repo := &stubRepo{
			allTags: func(_ context.Context) ([]domain.TagCount, error) {
				return []domain.TagCount{{Name: "rust", Count: 1}}, nil
			},
			allTokens: func(_ context.Context) ([]domain.TokenCount, error) {
				return []domain.TokenCount{{Value: "ownership", Count: 1}}, nil
			},
		}
		got, err := newInteractor(repo).GetSuggests(context.Background(), InputGetSuggestDto{Q: "go"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got.Suggestions) != 0 {
			t.Errorf("expected empty, got %v", got.Suggestions)
		}
	})

	t.Run("propagates allTags error", func(t *testing.T) {
		repoErr := errors.New("tags error")
		repo := &stubRepo{
			allTags: func(_ context.Context) ([]domain.TagCount, error) {
				return nil, repoErr
			},
		}
		_, err := newInteractor(repo).GetSuggests(context.Background(), InputGetSuggestDto{Q: "go"})
		if !errors.Is(err, repoErr) {
			t.Errorf("err = %v, want %v", err, repoErr)
		}
	})

	t.Run("propagates allTokens error", func(t *testing.T) {
		repoErr := errors.New("tokens error")
		repo := &stubRepo{
			allTags: func(_ context.Context) ([]domain.TagCount, error) {
				return nil, nil
			},
			allTokens: func(_ context.Context) ([]domain.TokenCount, error) {
				return nil, repoErr
			},
		}
		_, err := newInteractor(repo).GetSuggests(context.Background(), InputGetSuggestDto{Q: "go"})
		if !errors.Is(err, repoErr) {
			t.Errorf("err = %v, want %v", err, repoErr)
		}
	})
}

// -----------------------------------------------------------------------
// Search
// -----------------------------------------------------------------------

func TestInteractor_Search(t *testing.T) {
	t.Run("returns articles mapped to DTO", func(t *testing.T) {
		publishedAt := time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC)
		article, _ := domain.Reconstruct(domain.ReconstructArticleInput{
			ID:          "art-001",
			Title:       "Go入門",
			URL:         "https://example.com/art-001",
			Platform:    "qiita",
			Tags:        []string{"go"},
			PublishedAt: publishedAt,
			IsActive:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		})
		cursor := "next-cursor"
		repo := &stubRepo{
			findAll: func(_ context.Context, _ domain.SearchCriteria) (domain.FindAllResult, error) {
				return domain.FindAllResult{
					Articles:   []domain.Article{*article},
					NextCursor: &cursor,
				}, nil
			},
		}

		got, err := newInteractor(repo).Search(context.Background(), InputSearchDto{Limit: 10})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got.Articles) != 1 {
			t.Fatalf("len(Articles) = %d, want 1", len(got.Articles))
		}
		a := got.Articles[0]
		if a.ExternalID != "art-001" {
			t.Errorf("ExternalID = %q, want art-001", a.ExternalID)
		}
		if a.Title != "Go入門" {
			t.Errorf("Title = %q, want Go入門", a.Title)
		}
		if a.Platform != "qiita" {
			t.Errorf("Platform = %q, want qiita", a.Platform)
		}
		if !a.PublishedAt.Equal(publishedAt) {
			t.Errorf("PublishedAt = %v, want %v", a.PublishedAt, publishedAt)
		}
		if got.NextCursor != cursor {
			t.Errorf("NextCursor = %q, want %q", got.NextCursor, cursor)
		}
	})

	t.Run("always searches active articles only", func(t *testing.T) {
		var gotCriteria domain.SearchCriteria
		repo := &stubRepo{
			findAll: func(_ context.Context, c domain.SearchCriteria) (domain.FindAllResult, error) {
				gotCriteria = c
				return domain.FindAllResult{}, nil
			},
		}
		_, err := newInteractor(repo).Search(context.Background(), InputSearchDto{Limit: 10})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !gotCriteria.ActiveOnly {
			t.Error("ActiveOnly must be true")
		}
	})

	t.Run("maps filters to criteria", func(t *testing.T) {
		year := 2025
		platform := "qiita"
		cursor := "cur"
		var gotCriteria domain.SearchCriteria
		repo := &stubRepo{
			findAll: func(_ context.Context, c domain.SearchCriteria) (domain.FindAllResult, error) {
				gotCriteria = c
				return domain.FindAllResult{}, nil
			},
		}
		_, err := newInteractor(repo).Search(context.Background(), InputSearchDto{
			Tag:        []string{"go", "design"},
			Year:       &year,
			Platform:   &platform,
			Limit:      20,
			NextCursor: &cursor,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(gotCriteria.Tags) != 2 || gotCriteria.Tags[0] != "go" {
			t.Errorf("Tags = %v", gotCriteria.Tags)
		}
		if gotCriteria.Year == nil || *gotCriteria.Year != 2025 {
			t.Errorf("Year = %v", gotCriteria.Year)
		}
		if gotCriteria.Platform == nil || *gotCriteria.Platform != "qiita" {
			t.Errorf("Platform = %v", gotCriteria.Platform)
		}
		if gotCriteria.Limit != 20 {
			t.Errorf("Limit = %d, want 20", gotCriteria.Limit)
		}
		if gotCriteria.Cursor == nil || *gotCriteria.Cursor != cursor {
			t.Errorf("Cursor = %v", gotCriteria.Cursor)
		}
	})

	t.Run("Q is tokenized and set to criteria.Tokens", func(t *testing.T) {
		q := "Go入門"
		var gotCriteria domain.SearchCriteria
		repo := &stubRepo{
			findAll: func(_ context.Context, c domain.SearchCriteria) (domain.FindAllResult, error) {
				gotCriteria = c
				return domain.FindAllResult{}, nil
			},
		}
		_, err := newInteractor(repo).Search(context.Background(), InputSearchDto{Q: &q, Limit: 10})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// stubTokenizer は text をそのまま1トークンとして返す
		if len(gotCriteria.Tokens) != 1 || gotCriteria.Tokens[0] != q {
			t.Errorf("Tokens = %v, want [%q]", gotCriteria.Tokens, q)
		}
	})

	t.Run("Q nil sets no tokens", func(t *testing.T) {
		var gotCriteria domain.SearchCriteria
		repo := &stubRepo{
			findAll: func(_ context.Context, c domain.SearchCriteria) (domain.FindAllResult, error) {
				gotCriteria = c
				return domain.FindAllResult{}, nil
			},
		}
		_, err := newInteractor(repo).Search(context.Background(), InputSearchDto{Limit: 10})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(gotCriteria.Tokens) != 0 {
			t.Errorf("Tokens = %v, want empty", gotCriteria.Tokens)
		}
	})

	t.Run("propagates repo error", func(t *testing.T) {
		repoErr := errors.New("db error")
		repo := &stubRepo{
			findAll: func(_ context.Context, _ domain.SearchCriteria) (domain.FindAllResult, error) {
				return domain.FindAllResult{}, repoErr
			},
		}
		_, err := newInteractor(repo).Search(context.Background(), InputSearchDto{Limit: 10})
		if !errors.Is(err, repoErr) {
			t.Errorf("err = %v, want %v", err, repoErr)
		}
	})
}

// -----------------------------------------------------------------------
// Update
// -----------------------------------------------------------------------

func TestInteractor_Update(t *testing.T) {
	t.Run("updates title and saves", func(t *testing.T) {
		original := makeArticle("upd-001", "旧タイトル", "qiita", true)
		var saved *domain.Article
		repo := &stubRepo{
			findByExternalID: func(_ context.Context, _ string) (*domain.Article, error) {
				return original, nil
			},
			save: func(_ context.Context, a *domain.Article) error {
				saved = a
				return nil
			},
		}
		err := newInteractor(repo).Update(context.Background(), InputUpdateDto{
			ExternalID: "upd-001",
			Title:      "新タイトル",
			URL:        "https://example.com/upd-001",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if saved == nil {
			t.Fatal("Save not called")
		}
		if saved.Title() != "新タイトル" {
			t.Errorf("Title = %q, want 新タイトル", saved.Title())
		}
	})

	t.Run("returns error when article not found", func(t *testing.T) {
		repo := &stubRepo{
			findByExternalID: func(_ context.Context, _ string) (*domain.Article, error) {
				return nil, nil
			},
		}
		err := newInteractor(repo).Update(context.Background(), InputUpdateDto{ExternalID: "not-exist"})
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("propagates domain error when article is inactive", func(t *testing.T) {
		inactive := makeArticle("upd-002", "削除済み", "qiita", false)
		repo := &stubRepo{
			findByExternalID: func(_ context.Context, _ string) (*domain.Article, error) {
				return inactive, nil
			},
		}
		err := newInteractor(repo).Update(context.Background(), InputUpdateDto{
			ExternalID: "upd-002",
			Title:      "新タイトル",
			URL:        "https://example.com/upd-002",
		})
		if err == nil {
			t.Error("expected error for inactive article, got nil")
		}
	})

	t.Run("propagates repo error on FindByExternalID", func(t *testing.T) {
		repoErr := errors.New("db error")
		repo := &stubRepo{
			findByExternalID: func(_ context.Context, _ string) (*domain.Article, error) {
				return nil, repoErr
			},
		}
		err := newInteractor(repo).Update(context.Background(), InputUpdateDto{ExternalID: "upd-001"})
		if !errors.Is(err, repoErr) {
			t.Errorf("err = %v, want %v", err, repoErr)
		}
	})

	t.Run("propagates repo error on Save", func(t *testing.T) {
		repoErr := errors.New("save error")
		repo := &stubRepo{
			findByExternalID: func(_ context.Context, _ string) (*domain.Article, error) {
				return makeArticle("upd-003", "タイトル", "qiita", true), nil
			},
			save: func(_ context.Context, _ *domain.Article) error {
				return repoErr
			},
		}
		err := newInteractor(repo).Update(context.Background(), InputUpdateDto{
			ExternalID: "upd-003",
			Title:      "新タイトル",
			URL:        "https://example.com/upd-003",
		})
		if !errors.Is(err, repoErr) {
			t.Errorf("err = %v, want %v", err, repoErr)
		}
	})
}

// -----------------------------------------------------------------------
// Remove
// -----------------------------------------------------------------------

func TestInteractor_Remove(t *testing.T) {
	t.Run("deactivates article and saves", func(t *testing.T) {
		original := makeArticle("rm-001", "記事", "qiita", true)
		var saved *domain.Article
		repo := &stubRepo{
			findByExternalID: func(_ context.Context, _ string) (*domain.Article, error) {
				return original, nil
			},
			save: func(_ context.Context, a *domain.Article) error {
				saved = a
				return nil
			},
		}
		err := newInteractor(repo).Remove(context.Background(), InputRemoveDto{ExternalID: "rm-001"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if saved == nil {
			t.Fatal("Save not called")
		}
		if saved.IsActive() {
			t.Error("IsActive must be false after Remove")
		}
	})

	t.Run("returns error when article not found", func(t *testing.T) {
		repo := &stubRepo{
			findByExternalID: func(_ context.Context, _ string) (*domain.Article, error) {
				return nil, nil
			},
		}
		err := newInteractor(repo).Remove(context.Background(), InputRemoveDto{ExternalID: "not-exist"})
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("propagates repo error on FindByExternalID", func(t *testing.T) {
		repoErr := errors.New("db error")
		repo := &stubRepo{
			findByExternalID: func(_ context.Context, _ string) (*domain.Article, error) {
				return nil, repoErr
			},
		}
		err := newInteractor(repo).Remove(context.Background(), InputRemoveDto{ExternalID: "rm-001"})
		if !errors.Is(err, repoErr) {
			t.Errorf("err = %v, want %v", err, repoErr)
		}
	})

	t.Run("propagates repo error on Save", func(t *testing.T) {
		repoErr := errors.New("save error")
		repo := &stubRepo{
			findByExternalID: func(_ context.Context, _ string) (*domain.Article, error) {
				return makeArticle("rm-002", "記事", "qiita", true), nil
			},
			save: func(_ context.Context, _ *domain.Article) error {
				return repoErr
			},
		}
		err := newInteractor(repo).Remove(context.Background(), InputRemoveDto{ExternalID: "rm-002"})
		if !errors.Is(err, repoErr) {
			t.Errorf("err = %v, want %v", err, repoErr)
		}
	})
}

// -----------------------------------------------------------------------
// Register
// -----------------------------------------------------------------------

func TestInteractor_Register(t *testing.T) {
	baseInput := InputRegisterDto{
		ExternalID:  "reg-001",
		Title:       "新規記事",
		URL:         "https://qiita.com/reg-001",
		Platform:    "qiita",
		PublishedAt: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		Tags:        []string{"go"},
	}

	t.Run("saves new article", func(t *testing.T) {
		var saved *domain.Article
		repo := &stubRepo{
			findByExternalID: func(_ context.Context, _ string) (*domain.Article, error) {
				return nil, nil // 未登録
			},
			save: func(_ context.Context, a *domain.Article) error {
				saved = a
				return nil
			},
		}
		if err := newInteractor(repo).Register(context.Background(), baseInput); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if saved == nil {
			t.Fatal("Save not called")
		}
		if saved.ID() != "reg-001" {
			t.Errorf("ID = %q, want reg-001", saved.ID())
		}
		if saved.Title() != "新規記事" {
			t.Errorf("Title = %q, want 新規記事", saved.Title())
		}
		if saved.Platform() != "qiita" {
			t.Errorf("Platform = %q, want qiita", saved.Platform())
		}
		if !saved.IsActive() {
			t.Error("IsActive must be true for new article")
		}
	})

	t.Run("returns error when article already exists", func(t *testing.T) {
		repo := &stubRepo{
			findByExternalID: func(_ context.Context, _ string) (*domain.Article, error) {
				return makeArticle("reg-001", "既存記事", "qiita", true), nil
			},
		}
		if err := newInteractor(repo).Register(context.Background(), baseInput); err == nil {
			t.Error("expected error for duplicate article, got nil")
		}
	})

	t.Run("propagates repo error on FindByExternalID", func(t *testing.T) {
		repoErr := errors.New("db error")
		repo := &stubRepo{
			findByExternalID: func(_ context.Context, _ string) (*domain.Article, error) {
				return nil, repoErr
			},
		}
		err := newInteractor(repo).Register(context.Background(), baseInput)
		if !errors.Is(err, repoErr) {
			t.Errorf("err = %v, want %v", err, repoErr)
		}
	})

	t.Run("propagates repo error on Save", func(t *testing.T) {
		repoErr := errors.New("save error")
		repo := &stubRepo{
			findByExternalID: func(_ context.Context, _ string) (*domain.Article, error) {
				return nil, nil
			},
			save: func(_ context.Context, _ *domain.Article) error {
				return repoErr
			},
		}
		err := newInteractor(repo).Register(context.Background(), baseInput)
		if !errors.Is(err, repoErr) {
			t.Errorf("err = %v, want %v", err, repoErr)
		}
	})
}

// -----------------------------------------------------------------------
// Sync
// -----------------------------------------------------------------------

func TestInteractor_Sync(t *testing.T) {
	t.Run("indexes new articles not yet in repo", func(t *testing.T) {
		var saveCount int
		repo := &stubRepo{
			findByPlatform: func(_ context.Context, _ string) ([]domain.Article, error) {
				return nil, nil // 既存なし
			},
			save: func(_ context.Context, _ *domain.Article) error {
				saveCount++
				return nil
			},
		}
		fetcher := &stubFetcher{
			fetch: func(_ context.Context, _ string) ([]FetchedArticle, error) {
				return []FetchedArticle{
					{ExternalID: "sync-001", Title: "記事A", URL: "https://qiita.com/sync-001", Platform: "qiita"},
					{ExternalID: "sync-002", Title: "記事B", URL: "https://qiita.com/sync-002", Platform: "qiita"},
				}, nil
			},
		}

		result := NewInteractor(repo, fetcher, &stubTokenizer{}).Sync(context.Background(), "qiita")

		if result.Indexed != 2 {
			t.Errorf("Indexed = %d, want 2", result.Indexed)
		}
		if saveCount != 2 {
			t.Errorf("Save called %d times, want 2", saveCount)
		}
	})

	t.Run("reindexes existing articles with changed content", func(t *testing.T) {
		existing := makeArticle("sync-003", "旧タイトル", "zenn", true)
		var saved *domain.Article
		repo := &stubRepo{
			findByPlatform: func(_ context.Context, _ string) ([]domain.Article, error) {
				return []domain.Article{*existing}, nil
			},
			save: func(_ context.Context, a *domain.Article) error {
				saved = a
				return nil
			},
		}
		fetcher := &stubFetcher{
			fetch: func(_ context.Context, _ string) ([]FetchedArticle, error) {
				return []FetchedArticle{
					{ExternalID: "sync-003", Title: "新タイトル", URL: "https://zenn.dev/sync-003", Platform: "zenn"},
				}, nil
			},
		}

		result := NewInteractor(repo, fetcher, &stubTokenizer{}).Sync(context.Background(), "zenn")

		if result.Reindexed != 1 {
			t.Errorf("Reindexed = %d, want 1", result.Reindexed)
		}
		if saved == nil {
			t.Fatal("Save not called")
		}
		if saved.Title() != "新タイトル" {
			t.Errorf("Title = %q, want 新タイトル", saved.Title())
		}
	})

	t.Run("deactivates articles removed from platform", func(t *testing.T) {
		existing := makeArticle("sync-004", "消えた記事", "note", true)
		var saved *domain.Article
		repo := &stubRepo{
			findByPlatform: func(_ context.Context, _ string) ([]domain.Article, error) {
				return []domain.Article{*existing}, nil
			},
			save: func(_ context.Context, a *domain.Article) error {
				saved = a
				return nil
			},
		}
		fetcher := &stubFetcher{
			fetch: func(_ context.Context, _ string) ([]FetchedArticle, error) {
				return nil, nil // フィードから消えた
			},
		}

		result := NewInteractor(repo, fetcher, &stubTokenizer{}).Sync(context.Background(), "note")

		if result.Deactivated != 1 {
			t.Errorf("Deactivated = %d, want 1", result.Deactivated)
		}
		if saved == nil {
			t.Fatal("Save not called")
		}
		if saved.IsActive() {
			t.Error("IsActive must be false after deactivation")
		}
	})

	t.Run("propagates fetcher error", func(t *testing.T) {
		repo := &stubRepo{}
		fetcher := &stubFetcher{
			fetch: func(_ context.Context, _ string) ([]FetchedArticle, error) {
				return nil, errors.New("fetch error")
			},
		}
		result := NewInteractor(repo, fetcher, &stubTokenizer{}).Sync(context.Background(), "qiita")
		if len(result.Errors) == 0 {
			t.Error("expected Errors to be non-empty")
		}
	})
}
