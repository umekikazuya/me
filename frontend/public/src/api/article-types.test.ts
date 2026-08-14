import type { ArticleListResponse, ArticleSuggestResponse } from '@me/types'
import { describe, expect, it } from 'vitest'
import {
  normalizeArticleListResponse,
  normalizeArticleSuggestResponse,
  normalizeArticleTagListResponse,
} from './article-types.js'

/** spec が変わるとコンパイルエラーになる、ワイヤ契約に縛られたフィクスチャ。 */
const wireList: ArticleListResponse = {
  articles: [
    {
      externalId: 'b0adacee33b2774d7089',
      title: 'サトシ:「行け!!ピカチュウ!!」',
      url: 'https://qiita.com/umekikazuya/items/b0adacee33b2774d7089',
      platform: 'qiita',
      publishedAt: '2025-04-14T00:00:00Z',
      tags: ['design-pattern'],
    },
  ],
  nextCursor: 'eyJsYXN0X2tl',
}

describe('normalizeArticleListResponse', () => {
  it('accepts a spec-shaped payload, normalizing timestamps to a canonical ISO form', () => {
    expect(normalizeArticleListResponse(wireList)).toEqual({
      articles: [
        {
          externalId: 'b0adacee33b2774d7089',
          title: 'サトシ:「行け!!ピカチュウ!!」',
          url: 'https://qiita.com/umekikazuya/items/b0adacee33b2774d7089',
          platform: 'qiita',
          publishedAt: '2025-04-14T00:00:00.000Z',
          tags: ['design-pattern'],
        },
      ],
      nextCursor: 'eyJsYXN0X2tl',
    })
  })

  it.each([
    null,
    undefined,
    'x',
    42,
    {},
  ])('returns no articles for the malformed payload %p', (payload) => {
    expect(normalizeArticleListResponse(payload).articles).toEqual([])
  })

  // backend は publishedAt に omitempty を付けておらず、未設定時に
  // Go のゼロ値時刻を文字列で送ってくる。仕様上は「省略」なのでここで落とす。
  it.each([
    '0001-01-01T00:00:00Z',
    '0001-01-01T00:00:00+09:00',
    '0001-01-01T00:00:00.000Z',
  ])('drops the Go zero timestamp %s', (publishedAt) => {
    const [article] = normalizeArticleListResponse({
      articles: [{ publishedAt }],
    }).articles

    expect(article).not.toHaveProperty('publishedAt')
  })

  it('drops empty and unparsable timestamps', () => {
    const { articles } = normalizeArticleListResponse({
      articles: [{ publishedAt: '' }, { publishedAt: 'not-a-date' }, {}],
    })

    for (const article of articles) {
      expect(article).not.toHaveProperty('publishedAt')
    }
  })

  // backend は日時を .Local() で返しており、オフセット無しの JST 文字列が
  // 届きうる。バックエンドが UTC に修正されたらこのテストが赤くなる。
  it('normalizes timestamps to ISO 8601 UTC', () => {
    const [article] = normalizeArticleListResponse({
      articles: [{ publishedAt: '2025-06-01T09:00:00+09:00' }],
    }).articles

    expect(article.publishedAt).toBe('2025-06-01T00:00:00.000Z')
  })

  it.each([
    'hatena',
    undefined,
    42,
    '',
  ])('falls back to mochiya for the unknown platform %p', (platform) => {
    const [article] = normalizeArticleListResponse({
      articles: [{ platform }],
    }).articles

    expect(article.platform).toBe('mochiya')
  })

  it('trims tags and drops empty entries', () => {
    const [article] = normalizeArticleListResponse({
      articles: [{ tags: [' a ', '', null, 'b', '   '] }],
    }).articles

    expect(article.tags).toEqual(['a', 'b'])
  })

  it('normalizes an empty nextCursor to undefined', () => {
    expect(
      normalizeArticleListResponse({ articles: [], nextCursor: '' }).nextCursor,
    ).toBeUndefined()
  })
})

describe('normalizeArticleTagListResponse', () => {
  it('normalizes tags and defaults a missing count to zero', () => {
    expect(
      normalizeArticleTagListResponse({
        tags: [{ name: 'SOLID', count: 5 }, { name: 'no-count' }, null],
      }),
    ).toEqual([
      { name: 'SOLID', count: 5 },
      { name: 'no-count', count: 0 },
      { name: '', count: 0 },
    ])
  })

  it.each([
    null,
    undefined,
    'x',
    {},
  ])('returns an empty list for the malformed payload %p', (payload) => {
    expect(normalizeArticleTagListResponse(payload)).toEqual([])
  })
})

describe('normalizeArticleSuggestResponse', () => {
  const wireSuggest: ArticleSuggestResponse = {
    suggestions: [
      { type: 'tag', value: 'SOLID', count: 5 },
      { type: 'token', value: '設計', count: 8 },
      { type: 'title', value: '【SOLID原則】', externalId: 'a8751e422bf1' },
    ],
  }

  it('preserves each variant of the discriminated union', () => {
    expect(normalizeArticleSuggestResponse(wireSuggest)).toEqual([
      { type: 'tag', value: 'SOLID', count: 5 },
      { type: 'token', value: '設計', count: 8 },
      { type: 'title', value: '【SOLID原則】', externalId: 'a8751e422bf1' },
    ])
  })

  it('defaults a missing count to zero for tag and token variants', () => {
    expect(
      normalizeArticleSuggestResponse({
        suggestions: [{ type: 'tag', value: 'x' }],
      }),
    ).toEqual([{ type: 'tag', value: 'x', count: 0 }])
  })

  // 未知の type は形が分からないため、count や externalId を捏造せず捨てる。
  it('drops entries that cannot be placed in the union', () => {
    expect(
      normalizeArticleSuggestResponse({
        suggestions: [
          { type: 'unknown', value: 'x', count: 1 },
          { type: 'tag', value: '' },
          { type: 'title', value: 'no external id' },
          null,
          'x',
        ],
      }),
    ).toEqual([])
  })

  it.each([
    null,
    undefined,
    'x',
    {},
  ])('returns an empty list for the malformed payload %p', (payload) => {
    expect(normalizeArticleSuggestResponse(payload)).toEqual([])
  })
})
