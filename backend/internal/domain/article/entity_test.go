package article

import (
	"slices"
	"testing"
	"time"
)

var (
	validID       = "b0adacee33b2774d7089"
	validTitle    = "Factory Methodパターン入門"
	validURL      = "https://qiita.com/umekikazuya/items/b0adacee33b2774d7089"
	validPlatform = "qiita"
	validTokens   = []string{"SOLID", "原則"}
	validTags     = []string{"design-pattern", "go"}
	pastTime      = time.Now().AddDate(0, 0, -1)
	futureTime    = time.Now().AddDate(0, 0, 1)
)

// --- FO ---

func TestWithTags(t *testing.T) {
	tags := []string{"go", "aws"}
	entity, err := Index(
		validID,
		validTitle,
		validURL,
		validPlatform,
		WithTags(tags),
	)
	if err != nil {
		t.Fatal(err)
	}
	// Act
	got := entity.Tags()

	// Assert
	if len(got) != len(tags) {
		t.Errorf(
			"タグ数の想定は %d です。実際の取得数 %d 。",
			len(got),
			len(tags),
		)
	}

	// Act take2
	// WithTags に渡したスライスを後から書き換えた際に Article 内部の tags が変わらないことを担保
	tags[0] = "CHANGED"
	got2 := entity.Tags()
	// Assert
	if got2[0] == tags[0] {
		t.Errorf(
			"WithTags() は防御的コピーを保存すべきです: 入力スライスの外部操作によって内部状態が変更しました (want %q,  got %q)",
			"go", got2[0],
		)
	}

	// Act take3
	// getterの返り値を上書きした後に、エンティティに影響がないことを担保
	got[0] = "CHANGED"
	if entity.Tags()[0] == got[0] {
		t.Errorf(
			"Tags() の戻り値を変更した際、内部状態も変更されています。防御的コピーを返してください (want %q, got %q)",
			"go", entity.Tags()[0],
		)
	}
}

func TestWithTokens(t *testing.T) {
	tokens := []string{"go", "aws"}
	entity, err := Index(
		validID,
		validTitle,
		validURL,
		validPlatform,
		WithTokens(tokens),
	)
	if err != nil {
		t.Fatal(err)
	}
	// Act
	got := entity.Tokens()

	// Assert
	if len(got) != len(tokens) {
		t.Errorf(
			"トークン数の想定は %d です。実際の取得数 %d 。",
			len(got), len(tokens),
		)
	}

	// Act take2
	// WithTokens() に渡したスライスを後から書き換えた際に Article 内部の tokens が変わらないことを担保
	tokens[0] = "CHANGED"
	got2 := entity.Tokens()

	// Assert
	if got2[0] == tokens[0] {
		t.Errorf(
			"WithTokens() は防御的コピーを保存すべきです: 外部でのスライス操作が内部状態に影響を与えています (want %q,  got %q)",
			"go", got2[0],
		)
	}

	// Act take3
	// getterの返り値を上書きした後に、エンティティに影響がないことを担保
	got[0] = "CHANGED"
	if entity.Tokens()[0] == got[0] {
		t.Errorf(
			"Tokens() の戻り値を変更した際、内部状態も変更されています。防御的コピーを返してください (want %q, got %q)",
			"go", entity.Tokens()[0],
		)
	}
}

// --- Index ---

func TestIndex(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		title    string
		url      string
		platform string
		opts     []Opt
		wantErr  bool
	}{
		{
			name:     "valid – all fields",
			id:       validID,
			title:    validTitle,
			url:      validURL,
			platform: validPlatform,
			opts: []Opt{
				WithTags(validTags),
				WithTokens(validTokens),
				WithPublishedAt(pastTime),
				WithArticleUpdatedAt(pastTime),
			},
			wantErr: false,
		},
		{
			name:     "valid – opts omitted",
			id:       validID,
			title:    validTitle,
			url:      validURL,
			platform: validPlatform,
			wantErr:  false,
		},
		{
			name:     "invalid – empty externalId",
			id:       "",
			title:    validTitle,
			url:      validURL,
			platform: validPlatform,
			wantErr:  true,
		},
		{
			name:     "invalid – empty title",
			id:       validID,
			title:    "",
			url:      validURL,
			platform: validPlatform,
			wantErr:  true,
		},
		{
			name:     "invalid – empty url",
			id:       validID,
			title:    validTitle,
			url:      "",
			platform: validPlatform,
			wantErr:  true,
		},
		{
			name:     "invalid – unknown platform",
			id:       validID,
			title:    validTitle,
			url:      validURL,
			platform: "twitter",
			wantErr:  true,
		},
		{
			name:     "invalid – future publishedAt",
			id:       validID,
			title:    validTitle,
			url:      validURL,
			platform: validPlatform,
			opts:     []Opt{WithPublishedAt(futureTime)},
			wantErr:  true,
		},
		{
			name:     "invalid – future articleUpdatedAt",
			id:       validID,
			title:    validTitle,
			url:      validURL,
			platform: validPlatform,
			opts:     []Opt{WithArticleUpdatedAt(futureTime)},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Index(
				tt.id,
				tt.title,
				tt.url,
				tt.platform,
				tt.opts...,
			)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("Index() failed: %v", err)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Index() succeeded unexpectedly")
			}
			if got == nil {
				t.Fatal("Index() returned nil")
			}
			if !got.isActive.value {
				t.Errorf("Index() isActive = false, want true")
			}
			if got.createdAt.IsZero() {
				t.Errorf("Index() createdAt is zero")
			}
			if got.updatedAt.IsZero() {
				t.Errorf("Index() updatedAt is zero")
			}
		})
	}
}

// --- Register ---

func TestRegister(t *testing.T) {
	// Index と Register は振る舞いが同じ（入り口の意味だけが違う）
	t.Run("valid – same invariants as Index", func(t *testing.T) {
		got, err := Register(
			validID, validTitle, validURL, validPlatform,
			WithTags(validTags),
			WithTokens(validTokens),
			WithPublishedAt(pastTime),
		)
		if err != nil {
			t.Fatalf("Register() failed: %v", err)
		}
		if got == nil {
			t.Fatal("Register() returned nil")
		}
		if !got.isActive.value {
			t.Errorf("Register() isActive = false, want true")
		}
	})

	t.Run("invalid – unknown platform", func(t *testing.T) {
		_, err := Register(
			validID,
			validTitle,
			validURL,
			"twitter",
		)
		if err == nil {
			t.Fatal("Register() succeeded unexpectedly")
		}
	})
}

// --- Reindex ---

func TestArticle_Reindex(t *testing.T) {
	t.Run("success – updates fields and updatedAt", func(
		t *testing.T,
	) {
		a := activeArticle(t)
		prev := a.updatedAt
		reindexedPublishedAt := time.Now().AddDate(0, 0, -2)
		reindexedArticleUpdatedAt := time.Now().AddDate(0, 0, -3)

		err := a.Reindex(
			"新しいタイトル", validURL,
			WithTags([]string{"new-tag"}),
			WithTokens([]string{"新しい"}),
			WithPublishedAt(reindexedPublishedAt),
			WithArticleUpdatedAt(reindexedArticleUpdatedAt),
		)
		if err != nil {
			t.Fatalf("Reindex() failed: %v", err)
		}
		if a.title.value != "新しいタイトル" {
			t.Errorf("Reindex() title = %q, want %q", a.title.value, "新しいタイトル")
		}
		if !slices.Equal(a.tags, []string{"new-tag"}) {
			t.Errorf("Reindex() tags = %v, want %v", a.tags, []string{"new-tag"})
		}
		if !slices.Equal(a.tokens, []string{"新しい"}) {
			t.Errorf("Reindex() tokens = %v, want %v", a.tokens, []string{"新しい"})
		}
		if !a.publishedAt.value.Equal(reindexedPublishedAt) {
			t.Errorf("Reindex() publishedAt = %v, want %v", a.publishedAt.value, reindexedPublishedAt)
		}
		if !a.articleUpdatedAt.value.Equal(reindexedArticleUpdatedAt) {
			t.Errorf("Reindex() articleUpdatedAt = %v, want %v", a.articleUpdatedAt.value, reindexedArticleUpdatedAt)
		}
		if !a.updatedAt.After(prev) {
			t.Errorf("Reindex() did not advance updatedAt")
		}
	})

	t.Run("error – inactive article", func(t *testing.T) {
		a := inactiveArticle(t)
		err := a.Reindex("新しいタイトル", validURL)
		if err == nil {
			t.Fatal("Reindex() should return error for inactive article")
		}
	})

	t.Run("error – empty title", func(t *testing.T) {
		a := activeArticle(t)
		err := a.Reindex("", validURL)
		if err == nil {
			t.Fatal("Reindex() should return error for empty title")
		}
	})

	t.Run("error – future publishedAt", func(t *testing.T) {
		a := activeArticle(t)
		prevTitle := a.title.value
		prevPublishedAt := a.publishedAt.value
		err := a.Reindex(
			validTitle,
			validURL,
			WithPublishedAt(futureTime),
		)
		if err == nil {
			t.Fatal("Reindex() should return error for future publishedAt")
		}
		if a.title.value != prevTitle {
			t.Errorf("Reindex() changed title on failed update: got %q, want %q", a.title.value, prevTitle)
		}
		if !a.publishedAt.value.Equal(prevPublishedAt) {
			t.Errorf("Reindex() changed publishedAt on failed update: got %v, want %v", a.publishedAt.value, prevPublishedAt)
		}
	})
}

// --- Update ---

func TestArticle_Update(t *testing.T) {
	t.Run("success – updates fields and updatedAt", func(t *testing.T) {
		a := activeArticle(t)
		prev := a.updatedAt
		updatedPublishedAt := time.Now().AddDate(0, 0, -4)

		err := a.Update(
			"更新タイトル", validURL,
			WithTags([]string{"updated-tag"}),
			WithPublishedAt(updatedPublishedAt),
		)
		if err != nil {
			t.Fatalf("Update() failed: %v", err)
		}
		if a.title.value != "更新タイトル" {
			t.Errorf("Update() title = %q, want %q", a.title.value, "更新タイトル")
		}
		if !slices.Equal(a.tags, []string{"updated-tag"}) {
			t.Errorf("Update() tags = %v, want %v", a.tags, []string{"updated-tag"})
		}
		if !a.publishedAt.value.Equal(updatedPublishedAt) {
			t.Errorf("Update() publishedAt = %v, want %v", a.publishedAt.value, updatedPublishedAt)
		}
		if !a.updatedAt.After(prev) {
			t.Errorf("Update() did not advance updatedAt")
		}
	})

	t.Run("error – inactive article", func(t *testing.T) {
		a := inactiveArticle(t)
		err := a.Update("更新タイトル", validURL)
		if err == nil {
			t.Fatal("Update() should return error for inactive article")
		}
	})

	t.Run("error – future publishedAt", func(t *testing.T) {
		a := activeArticle(t)
		prevTitle := a.title.value
		prevPublishedAt := a.publishedAt.value
		err := a.Update(
			validTitle,
			validURL,
			WithPublishedAt(futureTime),
		)
		if err == nil {
			t.Fatal("Update() should return error for future publishedAt")
		}
		if a.title.value != prevTitle {
			t.Errorf("Update() changed title on failed update: got %q, want %q", a.title.value, prevTitle)
		}
		if !a.publishedAt.value.Equal(prevPublishedAt) {
			t.Errorf("Update() changed publishedAt on failed update: got %v, want %v", a.publishedAt.value, prevPublishedAt)
		}
	})
}

// --- Deactivate ---

func TestArticle_Deactivate(t *testing.T) {
	t.Run("sets isActive=false and advances updatedAt", func(t *testing.T) {
		a := activeArticle(t)
		prev := a.updatedAt

		if err := a.Deactivate(); err != nil {
			t.Fatalf("Deactivate() failed: %v", err)
		}
		if a.isActive.value {
			t.Errorf("Deactivate() isActive = true, want false")
		}
		if !a.updatedAt.After(prev) {
			t.Errorf("Deactivate() did not advance updatedAt")
		}
	})
}

// --- Remove ---

func TestArticle_Remove(t *testing.T) {
	t.Run("sets isActive=false and advances updatedAt", func(t *testing.T) {
		a := activeArticle(t)
		prev := a.updatedAt

		if err := a.Remove(); err != nil {
			t.Fatalf("Remove() failed: %v", err)
		}
		if a.isActive.value {
			t.Errorf("Remove() isActive = true, want false")
		}
		if !a.updatedAt.After(prev) {
			t.Errorf("Remove() did not advance updatedAt")
		}
	})
}

// --- helpers ---

func activeArticle(t *testing.T) *Article {
	t.Helper()
	a, err := Index(
		validID, validTitle, validURL, validPlatform,
		WithTags(validTags),
		WithTokens(validTokens),
		WithPublishedAt(pastTime),
		WithArticleUpdatedAt(pastTime),
	)
	if err != nil {
		t.Fatalf("activeArticle: %v", err)
	}
	return a
}

func inactiveArticle(t *testing.T) *Article {
	t.Helper()
	a := activeArticle(t)
	if err := a.Deactivate(); err != nil {
		t.Fatalf("inactiveArticle: %v", err)
	}
	return a
}
