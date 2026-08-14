import type { MeResponse } from '@me/types'
import { describe, expect, it } from 'vitest'
import {
  ApiError,
  describeApiError,
  describeProblemDetail,
  normalizeMeResponse,
} from './types.js'

/**
 * ワイヤ契約そのままのフィクスチャ。`MeResponse` 型注釈が付いているので、
 * spec が変わるとこのフィクスチャがコンパイルエラーになり、
 * normalize の追従漏れが CI で検知される。
 */
const wireProfile: MeResponse = {
  displayName: 'Kazuya Umeki',
  displayJa: '梅木 和弥',
  role: 'Web Creator',
  location: 'Fukuoka, Japan',
  skills: [
    { category: 'Frontend', items: ['Lit', 'TypeScript'], sortOrder: 0 },
  ],
  certifications: [
    { name: 'GCP ACE', issuer: 'Google Cloud', year: 2025, month: 6 },
  ],
  experiences: [
    {
      company: 'モチヤ株式会社',
      url: 'https://www.mochiya.ad.jp/',
      startYear: 2022,
    },
  ],
  links: [{ platform: 'github', url: 'https://github.com/umekikazuya' }],
  likes: ['Mr.Children'],
  createdAt: '2026-01-01T00:00:00Z',
  updatedAt: '2026-03-22T10:00:00Z',
}

describe('normalizeMeResponse', () => {
  it('passes a spec-shaped payload through unchanged', () => {
    expect(normalizeMeResponse(wireProfile)).toEqual({
      displayName: 'Kazuya Umeki',
      displayJa: '梅木 和弥',
      role: 'Web Creator',
      location: 'Fukuoka, Japan',
      skills: [
        { category: 'Frontend', items: ['Lit', 'TypeScript'], sortOrder: 0 },
      ],
      certifications: [
        { name: 'GCP ACE', issuer: 'Google Cloud', year: 2025, month: 6 },
      ],
      experiences: [
        {
          company: 'モチヤ株式会社',
          url: 'https://www.mochiya.ad.jp/',
          startYear: 2022,
        },
      ],
      links: [{ platform: 'github', url: 'https://github.com/umekikazuya' }],
      likes: ['Mr.Children'],
      updatedAt: '2026-03-22T10:00:00Z',
    })
  })

  it.each([
    null,
    undefined,
    'x',
    42,
    [],
  ])('returns an empty profile for the non-object payload %p', (payload) => {
    const profile = normalizeMeResponse(payload)

    expect(profile.displayName).toBe('')
    expect(profile.skills).toEqual([])
    expect(profile.certifications).toEqual([])
    expect(profile.experiences).toEqual([])
    expect(profile.links).toEqual([])
    expect(profile.likes).toEqual([])
    expect(profile.updatedAt).toBeUndefined()
  })

  // backend の OutputDto は全配列に omitempty が付いており、
  // かつ mapper が skills / experiences をセットしていないため実際に欠落する。
  it('fills omitted collections with empty arrays', () => {
    const profile = normalizeMeResponse({ displayName: 'x' })

    expect(profile.skills).toEqual([])
    expect(profile.experiences).toEqual([])
    expect(profile.certifications).toEqual([])
    expect(profile.links).toEqual([])
  })

  it('survives malformed array elements', () => {
    const profile = normalizeMeResponse({
      skills: [null, 1, {}, { category: 'A' }],
      links: [null, { platform: 'github' }],
    })

    expect(profile.skills).toEqual([
      { category: '', items: [], sortOrder: 0 },
      { category: '', items: [], sortOrder: 0 },
      { category: '', items: [], sortOrder: 0 },
      { category: 'A', items: [], sortOrder: 0 },
    ])
    expect(profile.links).toHaveLength(2)
  })

  // year / startYear は spec 上 required。欠落は契約違反なので
  // 実行時の今年などを捏造せず、エントリごと除外する。
  it('drops certifications and experiences that violate the required-year contract', () => {
    const profile = normalizeMeResponse({
      certifications: [{ name: 'no year' }, { name: 'ok', year: 2020 }],
      experiences: [
        { company: 'no start' },
        { company: 'ok', startYear: 2019 },
      ],
    })

    expect(profile.certifications).toEqual([{ name: 'ok', year: 2020 }])
    expect(profile.experiences).toEqual([{ company: 'ok', startYear: 2019 }])
  })

  it('omits optional fields instead of emitting empty strings', () => {
    const profile = normalizeMeResponse({
      certifications: [{ name: 'x', year: 2020, issuer: '' }],
      experiences: [{ company: 'y', startYear: 2019, url: '' }],
    })

    expect(profile.certifications[0]).not.toHaveProperty('issuer')
    expect(profile.certifications[0]).not.toHaveProperty('month')
    expect(profile.experiences[0]).not.toHaveProperty('url')
    expect(profile.experiences[0]).not.toHaveProperty('endYear')
  })

  it('normalizes an empty updatedAt to undefined', () => {
    expect(normalizeMeResponse({ updatedAt: '' }).updatedAt).toBeUndefined()
    expect(
      normalizeMeResponse({ updatedAt: '2026-01-01T00:00:00Z' }).updatedAt,
    ).toBe('2026-01-01T00:00:00Z')
  })

  it('drops non-string entries from likes', () => {
    expect(normalizeMeResponse({ likes: ['a', 1, null, 'b'] }).likes).toEqual([
      'a',
      'b',
    ])
  })
})

describe('describeProblemDetail', () => {
  it('prefers detail over every other field', () => {
    expect(
      describeProblemDetail({
        detail: 'detail wins',
        title: 'title',
        message: 'message',
        invalidParams: [{ name: 'a', reason: 'b' }],
      }),
    ).toBe('detail wins')
  })

  it('joins invalidParams and domain error details when detail is absent', () => {
    expect(
      describeProblemDetail({
        title: 'ignored',
        invalidParams: [
          { name: 'emailAddress', reason: 'must be a valid email address' },
        ],
        details: [{ field: 'experiences[0].endYear', message: 'too small' }],
      }),
    ).toBe(
      'emailAddress: must be a valid email address\nexperiences[0].endYear: too small',
    )
  })

  it('falls back through title, message, then status', () => {
    expect(describeProblemDetail({ title: 'Not Found' })).toBe('Not Found')
    expect(describeProblemDetail({ message: 'Invariant violation' })).toBe(
      'Invariant violation',
    )
    expect(describeProblemDetail(undefined, 502)).toBe(
      'API request failed with status 502',
    )
    expect(describeProblemDetail(undefined)).toBe('API request failed')
  })

  it('handles partial invalidParams entries', () => {
    expect(
      describeProblemDetail({ invalidParams: [{ reason: 'required' }] }),
    ).toBe('required')
    expect(describeProblemDetail({ invalidParams: [{ name: 'email' }] })).toBe(
      'email',
    )
  })
})

describe('describeApiError', () => {
  it('describes the problem carried by an ApiError', () => {
    const error = new ApiError('ignored', 400, { detail: 'bad input' })
    expect(describeApiError(error)).toBe('bad input')
  })

  it('falls back to the message of a plain Error', () => {
    expect(describeApiError(new Error('boom'))).toBe('boom')
  })

  it('returns a generic message for non-Error values', () => {
    expect(describeApiError('boom')).toBe('予期しないエラーが発生しました。')
  })
})
