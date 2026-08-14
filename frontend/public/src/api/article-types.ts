import type {
  ArticleItem,
  ArticleListQuery,
  ArticleSuggestionItem,
  ArticleTagItem,
  Platform,
} from '@me/types'

export type { ArticleItem, ArticleSuggestionItem, ArticleTagItem, Platform }

/** クエリパラメータは operations 由来。spec の変更が型で伝わる。 */
export type ArticleListParams = ArticleListQuery

export const articlePlatforms = ['qiita', 'zenn', 'mochiya', 'note'] as const

/**
 * 生成型は値を持てないため、ランタイム配列と spec の enum が一致しているかを
 * コンパイル時に突き合わせる。spec で platform が増減するとここが型エラーになる。
 */
type _AssertPlatformsInSync = [Platform] extends [
  (typeof articlePlatforms)[number],
]
  ? [(typeof articlePlatforms)[number]] extends [Platform]
    ? true
    : never
  : never
const _platformsInSync: _AssertPlatformsInSync = true
void _platformsInSync

export interface ArticleListResult {
  articles: ArticleItem[]
  nextCursor?: string
}

const isRecord = (value: unknown): value is Record<string, unknown> =>
  typeof value === 'object' && value !== null

const asString = (value: unknown) => (typeof value === 'string' ? value : '')

const asNumber = (value: unknown) =>
  typeof value === 'number' && Number.isFinite(value) ? value : 0

const asStringArray = (value: unknown) =>
  Array.isArray(value)
    ? value
        .filter((item): item is string => typeof item === 'string')
        .map((item) => item.trim())
        .filter(Boolean)
    : []

const isPlatform = (value: unknown): value is Platform =>
  typeof value === 'string' && articlePlatforms.includes(value as Platform)

const normalizePlatform = (value: unknown): Platform =>
  isPlatform(value) ? value : 'mochiya'

/** Go のゼロ値時刻。バックエンドが省略せず送ってくるため、ここで落とす。 */
const zeroDatePrefix = '0001-01-01T00:00:00'

const normalizeOptionalIsoDate = (value: unknown) => {
  if (typeof value !== 'string') return undefined

  const trimmed = value.trim()
  if (trimmed === '' || trimmed.startsWith(zeroDatePrefix)) return undefined

  const date = new Date(trimmed)
  return Number.isNaN(date.valueOf()) ? undefined : date.toISOString()
}

export const normalizeArticleItem = (value: unknown): ArticleItem => {
  const item = isRecord(value) ? value : {}
  const publishedAt = normalizeOptionalIsoDate(item.publishedAt)

  return {
    externalId: asString(item.externalId),
    title: asString(item.title),
    url: asString(item.url),
    platform: normalizePlatform(item.platform),
    tags: asStringArray(item.tags),
    ...(publishedAt ? { publishedAt } : {}),
  }
}

export const normalizeArticleListResponse = (
  payload: unknown,
): ArticleListResult => {
  const record = isRecord(payload) ? payload : {}
  const articles = Array.isArray(record.articles) ? record.articles : []

  return {
    articles: articles.map(normalizeArticleItem),
    nextCursor: asString(record.nextCursor) || undefined,
  }
}

export const normalizeArticleTagListResponse = (
  payload: unknown,
): ArticleTagItem[] => {
  const record = isRecord(payload) ? payload : {}
  const tags = Array.isArray(record.tags) ? record.tags : []

  return tags.map((value) => {
    const item = isRecord(value) ? value : {}
    return {
      name: asString(item.name),
      count: asNumber(item.count),
    }
  })
}

/**
 * サジェストは `type` による判別共用体。`title` バリアントは `count` を持たず
 * 代わりに `externalId` を持つため、バリアントごとに分けて組み立てる。
 * 未知の `type` は形が分からないので捏造せず捨てる。
 */
const normalizeSuggestion = (value: unknown): ArticleSuggestionItem | null => {
  if (!isRecord(value)) return null

  const suggestionValue = asString(value.value)
  if (!suggestionValue) return null

  if (value.type === 'tag' || value.type === 'token') {
    return {
      type: value.type,
      value: suggestionValue,
      count: asNumber(value.count),
    }
  }

  if (value.type === 'title') {
    const externalId = asString(value.externalId)
    return externalId
      ? { type: 'title', value: suggestionValue, externalId }
      : null
  }

  return null
}

export const normalizeArticleSuggestResponse = (
  payload: unknown,
): ArticleSuggestionItem[] => {
  const record = isRecord(payload) ? payload : {}
  const suggestions = Array.isArray(record.suggestions)
    ? record.suggestions
    : []

  return suggestions
    .map(normalizeSuggestion)
    .filter((item): item is ArticleSuggestionItem => item !== null)
}
