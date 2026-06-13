import { apiRequest } from './api.js'
import {
  type ArticleDraft,
  type ArticleListParams,
  type ArticleSuggestionItem,
  normalizeArticleListResponse,
  normalizeArticleSuggestResponse,
  normalizeArticleTagListResponse,
  toArticleCreateRequest,
  toArticleUpdateRequest,
} from './article-types.js'

const buildArticleQuery = (params: ArticleListParams) => {
  const query = new URLSearchParams()

  if (params.q?.trim()) query.set('q', params.q.trim())
  if (params.platform) query.set('platform', params.platform)
  if (params.year) query.set('year', String(params.year))
  if (params.limit) query.set('limit', String(params.limit))
  if (params.cursor) query.set('cursor', params.cursor)
  for (const tag of params.tag ?? []) {
    const trimmed = tag.trim()
    if (trimmed) query.append('tag', trimmed)
  }

  const serialized = query.toString()
  return serialized ? `?${serialized}` : ''
}

export const listArticles = async (params: ArticleListParams = {}) =>
  normalizeArticleListResponse(
    await apiRequest<unknown>(`/articles${buildArticleQuery(params)}`, {
      method: 'GET',
    }),
  )

export const listArticleTags = async () =>
  normalizeArticleTagListResponse(
    await apiRequest<unknown>('/articles/meta/tags', {
      method: 'GET',
    }),
  )

export const suggestArticles = async (
  query: string,
): Promise<ArticleSuggestionItem[]> =>
  normalizeArticleSuggestResponse(
    await apiRequest<unknown>(
      `/articles/meta/suggest?q=${encodeURIComponent(query)}`,
      {
        method: 'GET',
      },
    ),
  )

export const createArticle = async (draft: ArticleDraft) =>
  apiRequest<void>('/articles', {
    method: 'POST',
    body: toArticleCreateRequest(draft),
  })

export const updateArticle = async (externalId: string, draft: ArticleDraft) =>
  apiRequest<void>(`/articles/${encodeURIComponent(externalId)}`, {
    method: 'PUT',
    body: toArticleUpdateRequest(draft),
  })

export const deleteArticle = async (externalId: string) =>
  apiRequest<void>(`/articles/${encodeURIComponent(externalId)}`, {
    method: 'DELETE',
  })
