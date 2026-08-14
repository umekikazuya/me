import type { components, operations } from '@me/types'
import { apiRequest } from './api.js'

type ArticleListQuery = NonNullable<
  operations['listArticles']['parameters']['query']
>

const buildArticleQuery = (params: ArticleListQuery) => {
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

export const listArticles = (params: ArticleListQuery = {}) =>
  apiRequest<components['schemas']['ArticleListResponse']>(
    `/articles${buildArticleQuery(params)}`,
    { method: 'GET' },
  )

export const listArticleTags = () =>
  apiRequest<components['schemas']['ArticleTagListResponse']>(
    '/articles/meta/tags',
    { method: 'GET' },
  )

export const suggestArticles = (query: string) =>
  apiRequest<components['schemas']['ArticleSuggestResponse']>(
    `/articles/meta/suggest?q=${encodeURIComponent(query)}`,
    { method: 'GET' },
  )
