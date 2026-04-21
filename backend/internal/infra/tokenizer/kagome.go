package tokenizer

import (
	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

type KagomeTokenizer struct {
	t *tokenizer.Tokenizer
}

func NewKagomeTokenizer() (*KagomeTokenizer, error) {
	t, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	if err != nil {
		return nil, err
	}
	return &KagomeTokenizer{t: t}, nil
}

// Tokenize はテキストを形態素解析し、名詞・動詞・形容詞の表層形をユニークなトークン列として返す。
func (k *KagomeTokenizer) Tokenize(text string) []string {
	tokens := k.t.Tokenize(text)
	seen := make(map[string]struct{}, len(tokens))
	result := make([]string, 0, len(tokens))
	for _, tok := range tokens {
		features := tok.Features()
		if len(features) == 0 {
			continue
		}
		pos := features[0]
		if pos != "名詞" && pos != "動詞" && pos != "形容詞" {
			continue
		}
		surface := tok.Surface
		if surface == "" {
			continue
		}
		if _, ok := seen[surface]; ok {
			continue
		}
		seen[surface] = struct{}{}
		result = append(result, surface)
	}
	return result
}
