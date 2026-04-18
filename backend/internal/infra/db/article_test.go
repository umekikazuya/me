//go:build integration

package db

import (
	"context"
	"fmt"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	domain "github.com/umekikazuya/me/internal/domain/article"
)

const testTableName = "me.test"

func testEndpoint() string {
	port := os.Getenv("FLOCI_PORT")
	if port == "" {
		port = "4566"
	}
	return "http://localhost:" + port
}

var testClient *dynamodb.Client

func TestMain(m *testing.M) {
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion("ap-northeast-1"),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider("test", "test", ""),
		),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "config: %v\n", err)
		os.Exit(1)
	}
	testClient = dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(testEndpoint())
	})

	if _, err := testClient.DescribeTable(context.Background(), &dynamodb.DescribeTableInput{
		TableName: aws.String(testTableName),
	}); err != nil {
		fmt.Fprintf(os.Stderr, "table %q not found: %v\nrun: docker compose up --build -d floci\n", testTableName, err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

// -----------------------------------------------------------------------
// ヘルパー
// -----------------------------------------------------------------------

func newTestRepo() *ArticleDynamoRepo {
	return &ArticleDynamoRepo{
		client:    testClient,
		tableName: testTableName,
	}
}

// makeArticleDAO はテスト投入用アイテムを組み立てる。
// publishedAt を渡すと GSI1/GSI2/GSI3 属性も自動でセットする。
func makeArticleDAO(externalID, platform, publishedAt string, active bool) articleDao {
	now := time.Now().UTC().Format(time.RFC3339)
	dao := articleDao{
		PK:         articlePKPrefix + externalID,
		SK:         articleSK,
		ExternalID: externalID,
		Title:      "Test Article " + externalID,
		URL:        "https://example.com/" + externalID,
		Platform:   platform,
		IsActive:   active,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if publishedAt != "" {
		dao.PublishedAt = publishedAt
		dao.GSI1PK = "ARTICLES"
		dao.GSI1SK = publishedAt
		dao.GSI3PK = "PLATFORM#" + platform
		dao.GSI3SK = publishedAt
		if t, err := time.Parse(time.RFC3339, publishedAt); err == nil {
			dao.Year = t.Year()
			dao.GSI2PK = fmt.Sprintf("YEAR#%d", dao.Year)
			dao.GSI2SK = publishedAt
		}
	}
	return dao
}

// putTestArticle はアイテムを投入し、テスト終了後に自動削除する。
func putTestArticle(t *testing.T, dao articleDao) {
	t.Helper()
	item, err := attributevalue.MarshalMap(dao)
	if err != nil {
		t.Fatalf("putTestArticle marshal: %v", err)
	}
	if _, err = testClient.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String(testTableName),
		Item:      item,
	}); err != nil {
		t.Fatalf("putTestArticle PutItem: %v", err)
	}
	t.Cleanup(func() {
		_, _ = testClient.DeleteItem(context.Background(), &dynamodb.DeleteItemInput{
			TableName: aws.String(testTableName),
			Key: map[string]types.AttributeValue{
				"PK": &types.AttributeValueMemberS{Value: dao.PK},
				"SK": &types.AttributeValueMemberS{Value: dao.SK},
			},
		})
	})
}

// articleIDs は []domain.Article から ID 一覧を返す。
func articleIDs(articles []domain.Article) []string {
	ids := make([]string, 0, len(articles))
	for _, a := range articles {
		ids = append(ids, a.ID())
	}
	return ids
}

// -----------------------------------------------------------------------
// FindByExternalID
// -----------------------------------------------------------------------

func TestArticleDynamoRepo_FindByExternalID(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(t *testing.T)
		externalID string
		wantNil    bool
		wantErr    bool
		check      func(*testing.T, *domain.Article)
	}{
		{
			name: "found: returns correct article",
			setup: func(t *testing.T) {
				dao := makeArticleDAO("fbeid-001", "qiita", "2025-04-01T00:00:00Z", true)
				dao.Tags = []string{"go", "design-pattern"}
				dao.Tokens = []string{"設計", "パターン"}
				putTestArticle(t, dao)
			},
			externalID: "fbeid-001",
			check: func(t *testing.T, got *domain.Article) {
				if got.ID() != "fbeid-001" {
					t.Errorf("ID = %q, want fbeid-001", got.ID())
				}
				if got.Platform() != "qiita" {
					t.Errorf("Platform = %q, want qiita", got.Platform())
				}
				if got.Title() != "Test Article fbeid-001" {
					t.Errorf("Title = %q", got.Title())
				}
				if !got.IsActive() {
					t.Error("IsActive = false, want true")
				}
				if !slices.Equal(got.Tags(), []string{"go", "design-pattern"}) {
					t.Errorf("Tags = %v", got.Tags())
				}
				if got.PublishedAt().IsZero() {
					t.Error("PublishedAt is zero")
				}
			},
		},
		{
			name:       "not found: returns nil, no error",
			setup:      func(t *testing.T) {},
			externalID: "does-not-exist",
			wantNil:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)
			repo := newTestRepo()
			got, err := repo.FindByExternalID(context.Background(), tt.externalID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr = %v", err, tt.wantErr)
			}
			if tt.wantNil {
				if got != nil {
					t.Errorf("expected nil, got %+v", got)
				}
				return
			}
			if got == nil {
				t.Fatal("got nil, want article")
			}
			if tt.check != nil {
				tt.check(t, got)
			}
		})
	}
}

// -----------------------------------------------------------------------
// Save
// -----------------------------------------------------------------------

func TestArticleDynamoRepo_Save(t *testing.T) {
	t.Run("creates new article: retrievable via GetItem and GSI", func(t *testing.T) {
		article, err := domain.Index(
			"save-001",
			"保存テスト記事",
			"https://qiita.com/save-001",
			"qiita",
			domain.WithTags([]string{"go"}),
			domain.WithTokens([]string{"保存", "テスト"}),
			domain.WithPublishedAt(time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC)),
		)
		if err != nil {
			t.Fatalf("Index: %v", err)
		}
		t.Cleanup(func() {
			_, _ = testClient.DeleteItem(context.Background(), &dynamodb.DeleteItemInput{
				TableName: aws.String(testTableName),
				Key: map[string]types.AttributeValue{
					"PK": &types.AttributeValueMemberS{Value: articlePKPrefix + "save-001"},
					"SK": &types.AttributeValueMemberS{Value: articleSK},
				},
			})
		})

		repo := newTestRepo()
		if err := repo.Save(context.Background(), article); err != nil {
			t.Fatalf("Save: %v", err)
		}

		// GetItem 経由で取得できること
		got, err := repo.FindByExternalID(context.Background(), "save-001")
		if err != nil {
			t.Fatalf("FindByExternalID after Save: %v", err)
		}
		if got == nil {
			t.Fatal("FindByExternalID returned nil after Save")
		}
		if got.Title() != "保存テスト記事" {
			t.Errorf("Title = %q", got.Title())
		}

		// GSI1 経由（FindAll）でも取得できること。GSI に正しく書かれているかを検証する。
		result, err := repo.FindAll(context.Background(), domain.SearchCriteria{
			ActiveOnly: true,
			Limit:      100,
		})
		if err != nil {
			t.Fatalf("FindAll after Save: %v", err)
		}
		if !slices.Contains(articleIDs(result.Articles), "save-001") {
			t.Errorf("save-001 not found in FindAll result — GSI1 may not be set correctly")
		}
	})

	t.Run("overwrites existing article on re-save", func(t *testing.T) {
		putTestArticle(t, makeArticleDAO("save-002", "zenn", "2025-05-01T00:00:00Z", true))

		updated, err := domain.Reconstruct(domain.ReconstructArticleInput{
			ID:        "save-002",
			Title:     "更新後タイトル",
			URL:       "https://zenn.dev/save-002",
			Platform:  "zenn",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
		if err != nil {
			t.Fatalf("Reconstruct: %v", err)
		}

		repo := newTestRepo()
		if err := repo.Save(context.Background(), updated); err != nil {
			t.Fatalf("Save: %v", err)
		}

		got, err := repo.FindByExternalID(context.Background(), "save-002")
		if err != nil {
			t.Fatalf("FindByExternalID: %v", err)
		}
		if got.Title() != "更新後タイトル" {
			t.Errorf("Title = %q, want 更新後タイトル", got.Title())
		}
	})
}

// -----------------------------------------------------------------------
// FindAll
// -----------------------------------------------------------------------

func TestArticleDynamoRepo_FindAll(t *testing.T) {
	t.Run("returns active articles in publishedAt desc order, excludes inactive", func(t *testing.T) {
		putTestArticle(t, makeArticleDAO("fa-ord-001", "qiita", "2025-12-01T00:00:00Z", true))
		putTestArticle(t, makeArticleDAO("fa-ord-002", "zenn", "2025-06-01T00:00:00Z", true))
		putTestArticle(t, makeArticleDAO("fa-ord-003", "qiita", "2025-01-01T00:00:00Z", false)) // inactive

		repo := newTestRepo()
		result, err := repo.FindAll(context.Background(), domain.SearchCriteria{
			ActiveOnly: true,
			Limit:      100,
		})
		if err != nil {
			t.Fatalf("FindAll: %v", err)
		}

		ids := articleIDs(result.Articles)

		// inactive は返らない
		if slices.Contains(ids, "fa-ord-003") {
			t.Error("inactive fa-ord-003 must not be returned")
		}
		// active 2件は返る
		if !slices.Contains(ids, "fa-ord-001") {
			t.Errorf("fa-ord-001 not in result: %v", ids)
		}
		if !slices.Contains(ids, "fa-ord-002") {
			t.Errorf("fa-ord-002 not in result: %v", ids)
		}
		// publishedAt 降順: 001(12月) が 002(6月) より前
		idx001 := slices.Index(ids, "fa-ord-001")
		idx002 := slices.Index(ids, "fa-ord-002")
		if idx001 > idx002 {
			t.Errorf("expected fa-ord-001 before fa-ord-002, got %v", ids)
		}
	})

	t.Run("without active filter: includes inactive articles", func(t *testing.T) {
		putTestArticle(t, makeArticleDAO("fa-inc-001", "qiita", "2025-04-01T00:00:00Z", true))
		putTestArticle(t, makeArticleDAO("fa-inc-002", "qiita", "2025-03-01T00:00:00Z", false)) // inactive

		repo := newTestRepo()
		result, err := repo.FindAll(context.Background(), domain.SearchCriteria{
			ActiveOnly: false,
			Limit:      100,
		})
		if err != nil {
			t.Fatalf("FindAll: %v", err)
		}
		ids := articleIDs(result.Articles)
		if !slices.Contains(ids, "fa-inc-002") {
			t.Errorf("inactive fa-inc-002 should be included when ActiveOnly=false, got %v", ids)
		}
	})

	t.Run("tokens filter: AND condition", func(t *testing.T) {
		dao1 := makeArticleDAO("fa-tok-001", "qiita", "2025-04-01T00:00:00Z", true)
		dao1.Tokens = []string{"設計", "パターン", "go"} // 全部持つ
		dao2 := makeArticleDAO("fa-tok-002", "zenn", "2025-03-01T00:00:00Z", true)
		dao2.Tokens = []string{"設計"} // パターンなし
		dao3 := makeArticleDAO("fa-tok-003", "qiita", "2025-02-01T00:00:00Z", true)
		dao3.Tokens = []string{"パターン"} // 設計なし
		putTestArticle(t, dao1)
		putTestArticle(t, dao2)
		putTestArticle(t, dao3)

		repo := newTestRepo()
		result, err := repo.FindAll(context.Background(), domain.SearchCriteria{
			Tokens:     []string{"設計", "パターン"}, // AND: 両方必須
			ActiveOnly: true,
			Limit:      100,
		})
		if err != nil {
			t.Fatalf("FindAll: %v", err)
		}
		ids := articleIDs(result.Articles)

		if !slices.Contains(ids, "fa-tok-001") {
			t.Errorf("fa-tok-001 (has both tokens) must be in result: %v", ids)
		}
		if slices.Contains(ids, "fa-tok-002") {
			t.Errorf("fa-tok-002 (missing パターン) must not be in result: %v", ids)
		}
		if slices.Contains(ids, "fa-tok-003") {
			t.Errorf("fa-tok-003 (missing 設計) must not be in result: %v", ids)
		}
	})

	t.Run("tags filter: AND condition", func(t *testing.T) {
		dao1 := makeArticleDAO("fa-tag-001", "qiita", "2025-04-01T00:00:00Z", true)
		dao1.Tags = []string{"go", "design-pattern"}
		dao2 := makeArticleDAO("fa-tag-002", "zenn", "2025-03-01T00:00:00Z", true)
		dao2.Tags = []string{"go"} // design-pattern なし
		dao3 := makeArticleDAO("fa-tag-003", "qiita", "2025-02-01T00:00:00Z", true)
		dao3.Tags = []string{"design-pattern"} // go なし
		putTestArticle(t, dao1)
		putTestArticle(t, dao2)
		putTestArticle(t, dao3)

		repo := newTestRepo()
		result, err := repo.FindAll(context.Background(), domain.SearchCriteria{
			Tags:       []string{"go", "design-pattern"}, // AND
			ActiveOnly: true,
			Limit:      100,
		})
		if err != nil {
			t.Fatalf("FindAll: %v", err)
		}
		ids := articleIDs(result.Articles)

		if !slices.Contains(ids, "fa-tag-001") {
			t.Errorf("fa-tag-001 (has both tags) must be in result: %v", ids)
		}
		if slices.Contains(ids, "fa-tag-002") {
			t.Errorf("fa-tag-002 (missing design-pattern) must not be in result: %v", ids)
		}
		if slices.Contains(ids, "fa-tag-003") {
			t.Errorf("fa-tag-003 (missing go) must not be in result: %v", ids)
		}
	})

	t.Run("platform filter: returns only matching platform", func(t *testing.T) {
		putTestArticle(t, makeArticleDAO("fa-plt-001", "note", "2025-04-01T00:00:00Z", true))
		putTestArticle(t, makeArticleDAO("fa-plt-002", "mochiya", "2025-03-01T00:00:00Z", true))

		platform := "note"
		repo := newTestRepo()
		result, err := repo.FindAll(context.Background(), domain.SearchCriteria{
			Platform:   &platform,
			ActiveOnly: true,
			Limit:      100,
		})
		if err != nil {
			t.Fatalf("FindAll: %v", err)
		}
		ids := articleIDs(result.Articles)

		if !slices.Contains(ids, "fa-plt-001") {
			t.Errorf("fa-plt-001 must be in result: %v", ids)
		}
		if slices.Contains(ids, "fa-plt-002") {
			t.Errorf("fa-plt-002 (mochiya) must not be in result: %v", ids)
		}
		for _, a := range result.Articles {
			if a.Platform() != "note" {
				t.Errorf("unexpected platform %q for article %s", a.Platform(), a.ID())
			}
		}
	})

	t.Run("year filter: returns only matching year", func(t *testing.T) {
		putTestArticle(t, makeArticleDAO("fa-yr-001", "qiita", "2024-12-01T00:00:00Z", true))
		putTestArticle(t, makeArticleDAO("fa-yr-002", "qiita", "2025-01-01T00:00:00Z", true))

		year := 2025
		repo := newTestRepo()
		result, err := repo.FindAll(context.Background(), domain.SearchCriteria{
			Year:       &year,
			ActiveOnly: true,
			Limit:      100,
		})
		if err != nil {
			t.Fatalf("FindAll: %v", err)
		}
		ids := articleIDs(result.Articles)

		if !slices.Contains(ids, "fa-yr-002") {
			t.Errorf("fa-yr-002 (2025) must be in result: %v", ids)
		}
		if slices.Contains(ids, "fa-yr-001") {
			t.Errorf("fa-yr-001 (2024) must not be in result: %v", ids)
		}
	})

	t.Run("limit: returns nextCursor when more items exist", func(t *testing.T) {
		putTestArticle(t, makeArticleDAO("fa-cur-001", "qiita", "2025-04-01T00:00:00Z", true))
		putTestArticle(t, makeArticleDAO("fa-cur-002", "qiita", "2025-03-01T00:00:00Z", true))

		repo := newTestRepo()
		result, err := repo.FindAll(context.Background(), domain.SearchCriteria{
			ActiveOnly: true,
			Limit:      1,
		})
		if err != nil {
			t.Fatalf("FindAll: %v", err)
		}
		if len(result.Articles) != 1 {
			t.Errorf("len(Articles) = %d, want 1", len(result.Articles))
		}
		if result.NextCursor == nil {
			t.Error("NextCursor is nil, want non-nil")
		}
	})

	t.Run("cursor: second page returns different articles", func(t *testing.T) {
		putTestArticle(t, makeArticleDAO("fa-pg-001", "qiita", "2025-04-01T00:00:00Z", true))
		putTestArticle(t, makeArticleDAO("fa-pg-002", "qiita", "2025-03-01T00:00:00Z", true))

		repo := newTestRepo()
		first, err := repo.FindAll(context.Background(), domain.SearchCriteria{
			ActiveOnly: true,
			Limit:      1,
		})
		if err != nil || first.NextCursor == nil {
			t.Fatalf("first page: err=%v, cursor=%v", err, first.NextCursor)
		}

		second, err := repo.FindAll(context.Background(), domain.SearchCriteria{
			ActiveOnly: true,
			Limit:      1,
			Cursor:     first.NextCursor,
		})
		if err != nil {
			t.Fatalf("second page: %v", err)
		}
		if len(second.Articles) == 0 {
			t.Fatal("second page returned no articles")
		}
		if first.Articles[0].ID() == second.Articles[0].ID() {
			t.Errorf("second page returned same article as first: %s", first.Articles[0].ID())
		}
	})

	t.Run("last page has nil NextCursor", func(t *testing.T) {
		putTestArticle(t, makeArticleDAO("fa-lastpg-001", "note", "2025-04-01T00:00:00Z", true))
		putTestArticle(t, makeArticleDAO("fa-lastpg-002", "note", "2025-03-01T00:00:00Z", true))

		// note プラットフォームに絞ることで、この 2 件だけがヒットする状態を作る
		platform := "note"
		repo := newTestRepo()
		result, err := repo.FindAll(context.Background(), domain.SearchCriteria{
			Platform:   &platform,
			ActiveOnly: true,
			Limit:      100,
		})
		if err != nil {
			t.Fatalf("FindAll: %v", err)
		}
		if result.NextCursor != nil {
			t.Error("NextCursor should be nil on last page")
		}
	})
}

// -----------------------------------------------------------------------
// FindByPlatform
// -----------------------------------------------------------------------

func TestArticleDynamoRepo_FindByPlatform(t *testing.T) {
	t.Run("returns all articles for given platform", func(t *testing.T) {
		putTestArticle(t, makeArticleDAO("fbp-001", "mochiya", "2025-04-01T00:00:00Z", true))
		putTestArticle(t, makeArticleDAO("fbp-002", "mochiya", "2025-03-01T00:00:00Z", true))
		putTestArticle(t, makeArticleDAO("fbp-003", "note", "2025-02-01T00:00:00Z", true))

		repo := newTestRepo()
		got, err := repo.FindByPlatform(context.Background(), "mochiya")
		if err != nil {
			t.Fatalf("FindByPlatform: %v", err)
		}
		ids := articleIDs(got)

		if !slices.Contains(ids, "fbp-001") {
			t.Errorf("fbp-001 must be in result: %v", ids)
		}
		if !slices.Contains(ids, "fbp-002") {
			t.Errorf("fbp-002 must be in result: %v", ids)
		}
		if slices.Contains(ids, "fbp-003") {
			t.Errorf("fbp-003 (note) must not be in result: %v", ids)
		}
		for _, a := range got {
			if a.Platform() != "mochiya" {
				t.Errorf("unexpected platform %q for article %s", a.Platform(), a.ID())
			}
		}
	})

	t.Run("returns empty slice when no articles exist for platform", func(t *testing.T) {
		// mochiya のみ挿入し、note を検索 → 0件
		putTestArticle(t, makeArticleDAO("fbp-empty-001", "mochiya", "2025-04-01T00:00:00Z", true))

		repo := newTestRepo()
		got, err := repo.FindByPlatform(context.Background(), "note")
		if err != nil {
			t.Fatalf("FindByPlatform: %v", err)
		}
		// note の記事は挿入していないため 0件
		for _, a := range got {
			if a.ID() == "fbp-empty-001" {
				t.Error("mochiya article must not appear in note query")
			}
		}
	})
}

// -----------------------------------------------------------------------
// AllTags
// -----------------------------------------------------------------------

func TestArticleDynamoRepo_AllTags(t *testing.T) {
	t.Run("aggregates tags from active articles, excludes inactive", func(t *testing.T) {
		// 他テストと衝突しないよう "alltags-" プレフィックスで固有のタグを使う
		dao1 := makeArticleDAO("at-001", "qiita", "2025-04-01T00:00:00Z", true)
		dao1.Tags = []string{"alltags-go", "alltags-design"}
		dao2 := makeArticleDAO("at-002", "zenn", "2025-03-01T00:00:00Z", true)
		dao2.Tags = []string{"alltags-go", "alltags-test"}
		dao3 := makeArticleDAO("at-003", "qiita", "2025-02-01T00:00:00Z", false) // inactive
		dao3.Tags = []string{"alltags-go"}
		putTestArticle(t, dao1)
		putTestArticle(t, dao2)
		putTestArticle(t, dao3)

		repo := newTestRepo()
		tags, err := repo.AllTags(context.Background())
		if err != nil {
			t.Fatalf("AllTags: %v", err)
		}

		counts := make(map[string]int, len(tags))
		for _, tc := range tags {
			counts[tc.Name] = tc.Count
		}

		// inactive な at-003 の "alltags-go" は含まない → active 2件
		if counts["alltags-go"] != 2 {
			t.Errorf("alltags-go count = %d, want 2", counts["alltags-go"])
		}
		if counts["alltags-design"] != 1 {
			t.Errorf("alltags-design count = %d, want 1", counts["alltags-design"])
		}
		if counts["alltags-test"] != 1 {
			t.Errorf("alltags-test count = %d, want 1", counts["alltags-test"])
		}
	})
}

// -----------------------------------------------------------------------
// AllTokens
// -----------------------------------------------------------------------

func TestArticleDynamoRepo_AllTokens(t *testing.T) {
	t.Run("aggregates tokens from active articles", func(t *testing.T) {
		// 他テストと衝突しないよう "alltokens-" プレフィックスで固有のトークンを使う
		dao1 := makeArticleDAO("atok-001", "qiita", "2025-04-01T00:00:00Z", true)
		dao1.Tokens = []string{"alltokens-設計", "alltokens-パターン"}
		dao2 := makeArticleDAO("atok-002", "zenn", "2025-03-01T00:00:00Z", true)
		dao2.Tokens = []string{"alltokens-設計", "alltokens-テスト"}
		putTestArticle(t, dao1)
		putTestArticle(t, dao2)

		repo := newTestRepo()
		tokens, err := repo.AllTokens(context.Background())
		if err != nil {
			t.Fatalf("AllTokens: %v", err)
		}

		counts := make(map[string]int, len(tokens))
		for _, tc := range tokens {
			counts[tc.Value] = tc.Count
		}

		if counts["alltokens-設計"] != 2 {
			t.Errorf("alltokens-設計 count = %d, want 2", counts["alltokens-設計"])
		}
		if counts["alltokens-パターン"] != 1 {
			t.Errorf("alltokens-パターン count = %d, want 1", counts["alltokens-パターン"])
		}
		if counts["alltokens-テスト"] != 1 {
			t.Errorf("alltokens-テスト count = %d, want 1", counts["alltokens-テスト"])
		}
	})
}
