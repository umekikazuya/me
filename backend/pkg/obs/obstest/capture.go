// Package obstest はテスト用のログ/スパン capture ヘルパを提供する。
//
// `slog.SetDefault` を差し替えるパターンはテスト間の副作用を生むため避け、
// Capture を明示的に Logger に渡す DI で書く前提。
package obstest

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"sync"
	"testing"
)

// Capture は slog の JSON 出力を in-memory に貯める io.Writer 兼ヘルパ。
type Capture struct {
	mu  sync.Mutex
	buf *bytes.Buffer
}

// NewCapture は新しい Capture を返す。t は未使用だが、将来 Cleanup 等を組み込む余地を残す。
func NewCapture(t *testing.T) *Capture {
	t.Helper()
	return &Capture{buf: &bytes.Buffer{}}
}

// Handler は level 指定なし (既定 Info) の JSON Handler を返す。
func (c *Capture) Handler() slog.Handler {
	return slog.NewJSONHandler(c, nil)
}

// HandlerWithLevel は level 指定版。
func (c *Capture) HandlerWithLevel(lv slog.Leveler) slog.Handler {
	return slog.NewJSONHandler(c, &slog.HandlerOptions{Level: lv})
}

// Write は io.Writer 実装。JSON Handler の出力先として使われる。
func (c *Capture) Write(p []byte) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.buf.Write(p)
}

// Records は出力済み各行を JSON decode した結果を返す。
// decode に失敗した時点で打ち切って返す (テストヘルパ用途の割り切り)。
func (c *Capture) Records() []map[string]any {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := []map[string]any{}
	dec := json.NewDecoder(bytes.NewReader(c.buf.Bytes()))
	for dec.More() {
		var m map[string]any
		if err := dec.Decode(&m); err != nil {
			return out
		}
		out = append(out, m)
	}
	return out
}

// Reset は記録済みレコードを破棄する。
func (c *Capture) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.buf.Reset()
}
