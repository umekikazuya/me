import { describe, expect, it } from 'vitest'
import { timeAgo } from './build-info.js'

describe('timeAgo', () => {
  const now = new Date('2026-10-07T12:00:00Z')

  it('uses minutes, hours, then days', () => {
    expect(timeAgo(new Date('2026-10-07T11:55:00Z'), now)).toBe('5m ago')
    expect(timeAgo(new Date('2026-10-07T09:00:00Z'), now)).toBe('3h ago')
    expect(timeAgo(new Date('2026-10-05T12:00:00Z'), now)).toBe('2d ago')
  })
})
