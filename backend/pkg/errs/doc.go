// Package errs は、アプリケーション共通のエラー分類と
// HTTP エラーレスポンスへの変換機能を提供する。
//
// 本パッケージの公開 API は次の責務を持つ。
//   - sentinel error（ErrBadRequest など）によるカテゴリ判定（errors.Is）
//   - 構造化エラー（ValidationError ）による詳細表現（errors.As）
//   - WriteProblem による API エラーレスポンス出力
//   - WrapInternal によるインフラ起因エラーの内部エラー分類
package errs
