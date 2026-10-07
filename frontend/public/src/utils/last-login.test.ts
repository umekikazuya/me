import { beforeEach, describe, expect, it } from 'vitest'
import { formatLogin, takeLastLogin } from './last-login.js'

describe('takeLastLogin', () => {
  beforeEach(() => localStorage.clear())

  it('is null on the first visit and returns the previous visit after', () => {
    const first = new Date('2026-10-05T10:00:00Z')
    expect(takeLastLogin(first)).toBeNull()
    expect(takeLastLogin(new Date('2026-10-06T10:00:00Z'))).toEqual(first)
  })
})

describe('formatLogin', () => {
  it('formats like a terminal login banner', () => {
    expect(formatLogin(new Date(2026, 9, 6, 1, 41))).toBe('Tue Oct 6 01:41')
  })
})
