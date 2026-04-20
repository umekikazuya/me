package tokenizer

import (
	"slices"
	"testing"
)

func TestKagomeTokenizer_Tokenize(t *testing.T) {
	tok, err := NewKagomeTokenizer()
	if err != nil {
		t.Fatalf("NewKagomeTokenizer: %v", err)
	}

	t.Run("extracts nouns from Japanese title", func(t *testing.T) {
		tokens := tok.Tokenize("Goの並行処理入門")
		if len(tokens) == 0 {
			t.Fatal("expected non-empty tokens")
		}
		// "Go" と "入門" は名詞として抽出されるはず
		if !slices.Contains(tokens, "入門") {
			t.Errorf("tokens = %v, want to contain 入門", tokens)
		}
	})

	t.Run("deduplicates tokens", func(t *testing.T) {
		tokens := tok.Tokenize("Go Go Go")
		count := 0
		for _, tok := range tokens {
			if tok == "Go" {
				count++
			}
		}
		if count > 1 {
			t.Errorf("Go appears %d times, want 1", count)
		}
	})

	t.Run("empty string returns empty slice", func(t *testing.T) {
		tokens := tok.Tokenize("")
		if len(tokens) != 0 {
			t.Errorf("tokens = %v, want empty", tokens)
		}
	})

	t.Run("filters out particles and symbols", func(t *testing.T) {
		tokens := tok.Tokenize("GoとRustの比較")
		// "と" "の" は助詞なので含まれないはず
		for _, tok := range tokens {
			if tok == "と" || tok == "の" {
				t.Errorf("particle %q should be filtered out, tokens = %v", tok, tokens)
			}
		}
	})
}
