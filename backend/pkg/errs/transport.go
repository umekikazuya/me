package errs

// ProblemDetail は RFC 9457 (Problem Details for HTTP APIs) 準拠のエラーレスポンス。
type ProblemDetail struct {
	// Type は問題種別 URI。既定では about:blank を使う。
	Type string `json:"type"`
	// Title は HTTP ステータスに対応する短い要約。
	Title string `json:"title"`
	// Status は HTTP ステータスコード。
	Status int `json:"status"`
	// Detail は人間可読の詳細説明。
	Detail string `json:"detail,omitempty"`
	// Instance は問題が発生したリソース（通常はリクエストパス）。
	Instance string `json:"instance,omitempty"`
	// InvalidParams は入力パラメータ違反の詳細。
	InvalidParams []InvalidParam `json:"invalidParams,omitempty"`
}

// InvalidParam は単一パラメータの入力違反を表す。
type InvalidParam struct {
	// Name は不正だったパラメータ名。
	Name string `json:"name"`
	// Reason は不正理由。
	Reason string `json:"reason"`
}
