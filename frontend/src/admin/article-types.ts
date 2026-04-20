const zeroDatePrefix = '0001-01-01T00:00:00'

export const articlePlatforms = ['qiita', 'zenn', 'mochiya', 'note'] as const

export type ArticlePlatform = (typeof articlePlatforms)[number]

export interface ArticleItem {
  externalId: string
  title: string
  url: string
  platform: ArticlePlatform
  publishedAt?: string
  tags: string[]
}

export interface ArticleTagItem {
  name: string
  count: number
}

export interface ArticleSuggestionItem {
  type: 'tag' | 'token' | string
  value: string
  count: number
}

export interface ArticleListParams {
  q?: string
  tag?: string[]
  year?: number
  platform?: ArticlePlatform
  limit?: number
  cursor?: string
}

export interface ArticleListResult {
  articles: ArticleItem[]
  nextCursor?: string
}

export interface ArticleDraft {
  externalId: string
  title: string
  url: string
  platform: ArticlePlatform
  publishedAt: string
  articleUpdatedAt: string
  tags: string[]
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

const isArticlePlatform = (value: unknown): value is ArticlePlatform =>
  typeof value === 'string' &&
  articlePlatforms.includes(value as ArticlePlatform)

const normalizePlatform = (value: unknown): ArticlePlatform =>
  isArticlePlatform(value) ? value : 'mochiya'

const normalizeOptionalIsoDate = (value: unknown) => {
  if (typeof value !== 'string') return undefined

  const trimmed = value.trim()
  if (trimmed === '' || trimmed.startsWith(zeroDatePrefix)) return undefined

  const date = new Date(trimmed)
  return Number.isNaN(date.valueOf()) ? undefined : date.toISOString()
}

const pad = (value: number) => String(value).padStart(2, '0')

const toDateTimeLocal = (value?: string) => {
  if (!value) return ''

  const date = new Date(value)
  if (Number.isNaN(date.valueOf())) return ''

  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(
    date.getDate(),
  )}T${pad(date.getHours())}:${pad(date.getMinutes())}`
}

const toOptionalApiDate = (value: string) => {
  const trimmed = value.trim()
  if (trimmed === '') return undefined

  const date = new Date(trimmed)
  return Number.isNaN(date.valueOf()) ? undefined : date.toISOString()
}

export const createEmptyArticleDraft = (): ArticleDraft => ({
  externalId: '',
  title: '',
  url: '',
  platform: 'mochiya',
  publishedAt: '',
  articleUpdatedAt: '',
  tags: [],
})

export const cloneArticleDraft = (draft: ArticleDraft): ArticleDraft =>
  structuredClone(draft)

export const articleDraftFromArticle = (
  article: ArticleItem,
): ArticleDraft => ({
  externalId: article.externalId,
  title: article.title,
  url: article.url,
  platform: article.platform,
  publishedAt: toDateTimeLocal(article.publishedAt),
  articleUpdatedAt: '',
  tags: [...article.tags],
})

export const normalizeArticleListResponse = (
  payload: unknown,
): ArticleListResult => {
  const record = isRecord(payload) ? payload : {}
  const articles = Array.isArray(record.articles) ? record.articles : []

  return {
    articles: articles.map((value) => {
      const item = isRecord(value) ? value : {}
      return {
        externalId: asString(item.externalId),
        title: asString(item.title),
        url: asString(item.url),
        platform: normalizePlatform(item.platform),
        publishedAt: normalizeOptionalIsoDate(item.publishedAt),
        tags: asStringArray(item.tags),
      }
    }),
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

export const normalizeArticleSuggestResponse = (
  payload: unknown,
): ArticleSuggestionItem[] => {
  const record = isRecord(payload) ? payload : {}
  const suggestions = Array.isArray(record.suggestions)
    ? record.suggestions
    : []

  return suggestions.map((value) => {
    const item = isRecord(value) ? value : {}
    return {
      type: asString(item.type),
      value: asString(item.value),
      count: asNumber(item.count),
    }
  })
}

export const toArticleCreateRequest = (draft: ArticleDraft) => ({
  externalId: draft.externalId.trim(),
  title: draft.title.trim(),
  url: draft.url.trim(),
  platform: draft.platform,
  ...(toOptionalApiDate(draft.publishedAt)
    ? { publishedAt: toOptionalApiDate(draft.publishedAt) }
    : {}),
  ...(toOptionalApiDate(draft.articleUpdatedAt)
    ? { articleUpdatedAt: toOptionalApiDate(draft.articleUpdatedAt) }
    : {}),
  tags: draft.tags.map((tag) => tag.trim()).filter(Boolean),
})

export const toArticleUpdateRequest = (draft: ArticleDraft) => ({
  title: draft.title.trim(),
  url: draft.url.trim(),
  ...(toOptionalApiDate(draft.publishedAt)
    ? { publishedAt: toOptionalApiDate(draft.publishedAt) }
    : {}),
  ...(toOptionalApiDate(draft.articleUpdatedAt)
    ? { articleUpdatedAt: toOptionalApiDate(draft.articleUpdatedAt) }
    : {}),
  tags: draft.tags.map((tag) => tag.trim()).filter(Boolean),
})
