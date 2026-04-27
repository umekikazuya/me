package obs

// Attribute key constants.
//
// 値は OpenTelemetry Semantic Conventions (v1.26.0) に準拠する。
// 新規コードでは定数を経由して書き、リテラルの揺れを避ける。
//
// `request.id` のみ sem-conv 範囲外の独自拡張である。
const (
	AttrServiceName      = "service.name"
	AttrServiceVersion   = "service.version"
	AttrExceptionMessage = "exception.message"
	AttrExceptionType    = "exception.type"
	AttrExceptionStack   = "exception.stacktrace"
	AttrHTTPMethod       = "http.request.method"
	AttrURLPath          = "url.path"
	AttrTraceID          = "trace_id"
	AttrSpanID           = "span_id"
	AttrRequestID        = "request.id"
	AttrOp               = "code.function"
)
