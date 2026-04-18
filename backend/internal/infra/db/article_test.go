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
	if ep := os.Getenv("TEST_DYNAMODB_ENDPOINT"); ep != "" {
		return ep
	}
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
	os.Exit(m.Run())
}

// fullArticleDao はGSI属性を含むテスト投入用DAO。
// productionのarticleDaoにはGSI属性がないため、テスト専用で定義する。
type fullArticleDao struct {
	PK               string   `dynamodbav:"PK"`
	SK               string   `dynamodbav:"SK"`
	GSI1PK           string   `dynamodbav:"GSI1PK,omitempty"`
	GSI1SK           string   `dynamodbav:"GSI1SK,omitempty"`
	GSI2PK           string   `dynamodbav:"GSI2PK,omitempty"`
	GSI2SK           string   `dynamodbav:"GSI2SK,omitempty"`
	GSI3PK           string   `dynamodbav:"GSI3PK,omitempty"`
	GSI3SK           string   `dynamodbav:"GSI3SK,omitempty"`
	ExternalID       string   `dynamodbav:"externalId"`
	Title            string   `dynamodbav:"title"`
	URL              string   `dynamodbav:"url"`
	Platform         string   `dynamodbav:"platform"`
	PublishedAt      string   `dynamodbav:"publishedAt,omitempty"`
	ArticleUpdatedAt string   `dynamodbav:"articleUpdatedAt,omitempty"`
	Tags             []string `dynamodbav:"tags,omitempty"`
	Tokens           []string `dynamodbav:"tokens,omitempty"`
	IsActive         bool     `dynamodbav:"isActive"`
	Year             int      `dynamodbav:"year,omitempty"`
	CreatedAt        string   `dynamodbav:"createdAt"`
	UpdatedAt        string   `dynamodbav:"updatedAt"`
}

// newTestRepo はテスト用リポジトリを生成する。
func newTestRepo() *ArticleDynamoRepo {
	return &ArticleDynamoRepo{
		client:    testClient,
		tableName: testTableName,
	}
}

// makeArticleDAO はテスト用アイテムのベースを組み立てる。
// publishedAt を指定するとGSI属性も自動でセットされる。
func makeArticleDAO(externalID, platform, publishedAt string, active bool) fullArticleDao {
	now := time.Now().UTC().Format(time.RFC3339)
	dao := fullArticleDao{
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

// putTestArticle はアイテムをDynamoDBに投入し、テスト終了時に自動削除する。
func putTestArticle(t *testing.T, dao fullArticleDao) {
	t.Helper()
	item, err := attributevalue.MarshalMap(dao)
	if err != nil {
		t.Fatalf("putTestArticle marshal: %v", err)
	}
	_, err = testClient.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String(testTableName),
		Item:      item,
	})
	if err != nil {
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
					t.Errorf("ID = %q, want %q", got.ID(), "fbeid-001")
				}
				if got.Platform() != "qiita" {
					t.Errorf("Platform = %q, want %q", got.Platform(), "qiita")
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
			name:       "not found: returns nil",
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
	t.Run("creates new article and is retrievable", func(t *testing.T) {
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

		got, err := repo.FindByExternalID(context.Background(), "save-001")
		if err != nil {
			t.Fatalf("FindByExternalID after Save: %v", err)
		}
		if got == nil {
			t.Fatal("FindByExternalID returned nil after Save")
		}
		if got.ID() != "save-001" {
			t.Errorf("ID = %q", got.ID())
		}
		if got.Title() != "保存テスト記事" {
			t.Errorf("Title = %q", got.Title())
		}
		if !got.IsActive() {
			t.Error("IsActive = false, want true")
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
			t.Errorf("Title = %q, want %q", got.Title(), "更新後タイトル")
		}
	})
}

// -----------------------------------------------------------------------
// FindAll
// -----------------------------------------------------------------------

func TestArticleDynamoRepo_FindAll(t *testing.T) {
	// 各サブテストが独立したフィクスチャを持つ
	t.Run("returns active articles in publishedAt desc order", func(t *testing.T) {
		putTestArticle(t, makeArticleDAO("fa-ord-001", "qiita", "2025-12-01T00:00:00Z", true))
		putTestArticle(t, makeArticleDAO("fa-ord-002", "zenn", "2025-06-01T00:00:00Z", true))
		putTestArticle(t, makeArticleDAO("fa-ord-003", "qiita", "2025-01-01T00:00:00Z", false)) // inactive

		repo := newTestRepo()
		result, err := repo.FindAll(context.Background(), domain.SearchCriteria{
			ActiveOnly: true,
			Limit:      10,
		})
		if err != nil {
			t.Fatalf("FindAll: %v", err)
		}

		for _, a := range result.Articles {
			if a.ID() == "fa-ord-003" {
				t.Error("inactive article must not be returned")
			}
		}

		ids := make([]string, 0, len(result.Articles))
		for _, a := range result.Articles {
			ids = append(ids, a.ID())
		}
		idx001 := slices.Index(ids, "fa-ord-001")
		idx002 := slices.Index(ids, "fa-ord-002")
		if idx001 == -1 || idx002 == -1 {
			t.Fatalf("expected fa-ord-001 and fa-ord-002 in result, got %v", ids)
		}
		if idx001 > idx002 {
			t.Errorf("expected fa-ord-001 before fa-ord-002 (publishedAt desc), got order %v", ids)
		}
	})

	t.Run("with tokens filter: returns only matching articles", func(t *testing.T) {
		dao1 := makeArticleDAO("fa-tok-001", "qiita", "2025-04-01T00:00:00Z", true)
		dao1.Tokens = []string{"設計", "パターン"}
		dao2 := makeArticleDAO("fa-tok-002", "zenn", "2025-03-01T00:00:00Z", true)
		dao2.Tokens = []string{"テスト", "go"}
		putTestArticle(t, dao1)
		putTestArticle(t, dao2)

		repo := newTestRepo()
		result, err := repo.FindAll(context.Background(), domain.SearchCriteria{
			Tokens:     []string{"設計"},
			ActiveOnly: true,
			Limit:      10,
		})
		if err != nil {
			t.Fatalf("FindAll: %v", err)
		}
		for _, a := range result.Articles {
			if a.ID() == "fa-tok-002" {
				t.Errorf("article without matching token must not be returned: %s", a.ID())
			}
		}
		ids := make([]string, 0, len(result.Articles))
		for _, a := range result.Articles {
			ids = append(ids, a.ID())
		}
		if !slices.Contains(ids, "fa-tok-001") {
			t.Errorf("expected fa-tok-001 in result, got %v", ids)
		}
	})

	t.Run("with tags filter: returns only matching articles", func(t *testing.T) {
		dao1 := makeArticleDAO("fa-tag-001", "qiita", "2025-04-01T00:00:00Z", true)
		dao1.Tags = []string{"go", "design-pattern"}
		dao2 := makeArticleDAO("fa-tag-002", "zenn", "2025-03-01T00:00:00Z", true)
		dao2.Tags = []string{"go", "testing"}
		putTestArticle(t, dao1)
		putTestArticle(t, dao2)

		repo := newTestRepo()
		result, err := repo.FindAll(context.Background(), domain.SearchCriteria{
			Tags:       []string{"design-pattern"},
			ActiveOnly: true,
			Limit:      10,
		})
		if err != nil {
			t.Fatalf("FindAll: %v", err)
		}
		for _, a := range result.Articles {
			if a.ID() == "fa-tag-002" {
				t.Errorf("article without matching tag must not be returned: %s", a.ID())
			}
		}
	})

	t.Run("with platform filter: returns only matching platform", func(t *testing.T) {
		putTestArticle(t, makeArticleDAO("fa-plt-001", "note", "2025-04-01T00:00:00Z", true))
		putTestArticle(t, makeArticleDAO("fa-plt-002", "mochiya", "2025-03-01T00:00:00Z", true))

		platform := "note"
		repo := newTestRepo()
		result, err := repo.FindAll(context.Background(), domain.SearchCriteria{
			Platform:   &platform,
			ActiveOnly: true,
			Limit:      10,
		})
		if err != nil {
			t.Fatalf("FindAll: %v", err)
		}
		for _, a := range result.Articles {
			if a.Platform() != "note" {
				t.Errorf("unexpected platform %q for article %s", a.Platform(), a.ID())
			}
		}
	})

	t.Run("with year filter: returns only matching year", func(t *testing.T) {
		putTestArticle(t, makeArticleDAO("fa-yr-001", "qiita", "2024-12-01T00:00:00Z", true))
		putTestArticle(t, makeArticleDAO("fa-yr-002", "qiita", "2025-01-01T00:00:00Z", true))

		year := 2025
		repo := newTestRepo()
		result, err := repo.FindAll(context.Background(), domain.SearchCriteria{
			Year:       &year,
			ActiveOnly: true,
			Limit:      10,
		})
		if err != nil {
			t.Fatalf("FindAll: %v", err)
		}
		for _, a := range result.Articles {
			if a.PublishedAt().Year() != 2025 {
				t.Errorf("year = %d, want 2025 for article %s", a.PublishedAt().Year(), a.ID())
			}
		}
	})

	t.Run("with limit: returns nextCursor when more items exist", func(t *testing.T) {
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
		for _, a := range got {
			if a.Platform() != "mochiya" {
				t.Errorf("unexpected platform %q for article %s", a.Platform(), a.ID())
			}
		}
		ids := make([]string, 0, len(got))
		for _, a := range got {
			ids = append(ids, a.ID())
		}
		if !slices.Contains(ids, "fbp-001") || !slices.Contains(ids, "fbp-002") {
			t.Errorf("expected fbp-001 and fbp-002, got %v", ids)
		}
	})

	t.Run("returns empty slice for unknown platform", func(t *testing.T) {
		// 存在しないプラットフォームに対しては空スライスを返す
		// (allowedPlatforms にない値はドメイン層で弾かれるが、インフラ層は値をそのまま扱う)
		repo := newTestRepo()
		// 既存データを用意しないことで、qiita に一切アイテムがない状態を確保するのは難しいため
		// ここでは返却されたアイテムがすべてそのプラットフォームであることだけを保証する
		got, err := repo.FindByPlatform(context.Background(), "note")
		if err != nil {
			t.Fatalf("FindByPlatform: %v", err)
		}
		for _, a := range got {
			if a.Platform() != "note" {
				t.Errorf("unexpected platform %q", a.Platform())
			}
		}
	})
}

// -----------------------------------------------------------------------
// AllTags
// -----------------------------------------------------------------------

func TestArticleDynamoRepo_AllTags(t *testing.T) {
	t.Run("aggregates tags from active articles with correct counts", func(t *testing.T) {
		dao1 := makeArticleDAO("at-001", "qiita", "2025-04-01T00:00:00Z", true)
		dao1.Tags = []string{"go", "design-pattern"}
		dao2 := makeArticleDAO("at-002", "zenn", "2025-03-01T00:00:00Z", true)
		dao2.Tags = []string{"go", "testing"}
		dao3 := makeArticleDAO("at-003", "qiita", "2025-02-01T00:00:00Z", false) // inactive
		dao3.Tags = []string{"go"}
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
		// inactive な at-003 の "go" は含まない → active 2件
		if counts["go"] != 2 {
			t.Errorf("go count = %d, want 2", counts["go"])
		}
		if counts["design-pattern"] != 1 {
			t.Errorf("design-pattern count = %d, want 1", counts["design-pattern"])
		}
		if counts["testing"] != 1 {
			t.Errorf("testing count = %d, want 1", counts["testing"])
		}
	})
}

// -----------------------------------------------------------------------
// AllTokens
// -----------------------------------------------------------------------

func TestArticleDynamoRepo_AllTokens(t *testing.T) {
	t.Run("aggregates tokens from active articles with correct counts", func(t *testing.T) {
		dao1 := makeArticleDAO("atok-001", "qiita", "2025-04-01T00:00:00Z", true)
		dao1.Tokens = []string{"設計", "パターン"}
		dao2 := makeArticleDAO("atok-002", "zenn", "2025-03-01T00:00:00Z", true)
		dao2.Tokens = []string{"設計", "テスト"}
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
		if counts["設計"] != 2 {
			t.Errorf("設計 count = %d, want 2", counts["設計"])
		}
		if counts["パターン"] != 1 {
			t.Errorf("パターン count = %d, want 1", counts["パターン"])
		}
		if counts["テスト"] != 1 {
			t.Errorf("テスト count = %d, want 1", counts["テスト"])
		}
	})
}
