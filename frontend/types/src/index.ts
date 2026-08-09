/**
 * OpenAPI 生成型の公開エントリポイント。
 *
 * `docs/developments/api/openapi.yaml` が唯一の SoT。このファイルは
 * `components['schemas'][...]` を扱いやすい名前に付け替えるだけの薄い層で、
 * 手書きの型定義は一切置かない。
 *
 * 明示的に列挙しているのは意図的。spec でスキーマ名が変わると
 * このファイル自体がコンパイルエラーになり、ドリフトの検知点として機能する。
 *
 * ランタイム値をここに置かないこと。型のみに保つことで、CI の再生成 diff 検知の
 * 対象が `src/generated.ts` 1 ファイルに閉じる。
 */
import type { components, operations } from './generated.js'

export type { components, operations, paths, webhooks } from './generated.js'

type Schemas = components['schemas']

/* ── Articles ─────────────────────────────────────────── */

export type Platform = Schemas['Platform']
export type ArticleItem = Schemas['ArticleItem']
export type ArticleListResponse = Schemas['ArticleListResponse']
export type ArticleTagItem = Schemas['ArticleTagItem']
export type ArticleTagListResponse = Schemas['ArticleTagListResponse']

/**
 * `type` による判別共用体。全体と各バリアントの両方を公開しているのは、
 * 消費側が narrow 後の型を関数シグネチャに書けるようにするため。
 */
export type ArticleSuggestionItem = Schemas['ArticleSuggestionItem']
export type ArticleTagSuggestion = Schemas['ArticleTagSuggestion']
export type ArticleTokenSuggestion = Schemas['ArticleTokenSuggestion']
export type ArticleTitleSuggestion = Schemas['ArticleTitleSuggestion']
export type ArticleSuggestResponse = Schemas['ArticleSuggestResponse']

export type ArticleCreateRequest = Schemas['ArticleCreateRequest']
export type ArticleUpdateRequest = Schemas['ArticleUpdateRequest']

/* ── Me ───────────────────────────────────────────────── */

export type MeRequest = Schemas['MeRequest']
export type MeResponse = Schemas['MeResponse']
export type MeSkillGroup = Schemas['MeSkillGroup']
export type MeCertification = Schemas['MeCertification']
export type MeExperience = Schemas['MeExperience']
export type MeLink = Schemas['MeLink']

/* ── Identity ─────────────────────────────────────────── */

export type LoginRequest = Schemas['LoginRequest']
export type RegisterRequest = Schemas['RegisterRequest']
export type PasswordResetRequest = Schemas['PasswordResetRequest']
export type ChangeEmailRequest = Schemas['ChangeEmailRequest']

/* ── Errors ───────────────────────────────────────────── */

/** RFC 9457 準拠の HTTP エラー (400/401/403/404/409/500)。 */
export type ProblemDetails = Schemas['ProblemDetails']
export type InvalidParam = Schemas['InvalidParam']

/** ドメイン不変条件違反 (422) の専用形。ProblemDetails とは shape で区別する。 */
export type DomainErrorResponse = Schemas['DomainErrorResponse']
export type DomainErrorFieldDetail = Schemas['DomainErrorFieldDetail']

/* ── クエリパラメータ (operations から導出) ─────────────── */

/**
 * 手書きで持つとクエリパラメータの追加・変更が FE に伝わらないため、
 * operations から導出する。`platform` は Platform を経由するので enum も自動追従する。
 */
export type ArticleListQuery = NonNullable<
  operations['listArticles']['parameters']['query']
>
export type ArticleSuggestQuery = NonNullable<
  operations['suggestArticles']['parameters']['query']
>
