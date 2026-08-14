import { describe, expect, it } from 'vitest'
import { ApiError, describeApiError, describeProblemDetail } from './types.js'

describe('describeProblemDetail', () => {
  it('prefers detail over every other field', () => {
    expect(
      describeProblemDetail({
        detail: 'detail wins',
        title: 'title',
        invalidParams: [{ name: 'a', reason: 'b' }],
      }),
    ).toBe('detail wins')
  })

  it('joins invalidParams when detail is absent', () => {
    expect(
      describeProblemDetail({
        title: 'ignored',
        invalidParams: [
          { name: 'emailAddress', reason: 'must be a valid email address' },
          { name: 'year', reason: 'must be at least 1' },
        ],
      }),
    ).toBe(
      'emailAddress: must be a valid email address\nyear: must be at least 1',
    )
  })

  it('falls back through title, then status', () => {
    expect(describeProblemDetail({ title: 'Not Found' })).toBe('Not Found')
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
