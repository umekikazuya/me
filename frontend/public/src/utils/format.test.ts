import { describe, expect, it } from 'vitest'
import { formatDate, sanitizeUrl } from './format.js'

describe('formatDate', () => {
  it('formats a date as yyyy.mm.dd', () => {
    expect(formatDate('2026-09-08T12:00:00')).toBe('2026.09.08')
  })

  it('falls back for missing or invalid values', () => {
    expect(formatDate()).toBe('----.--.--')
    expect(formatDate('nope')).toBe('----.--.--')
  })
})

describe('sanitizeUrl', () => {
  it('keeps http(s) and mailto links', () => {
    expect(sanitizeUrl(' https://github.com/x ')).toBe('https://github.com/x')
    expect(sanitizeUrl('mailto:a@example.com')).toBe('mailto:a@example.com')
  })

  it('neutralizes other schemes', () => {
    expect(sanitizeUrl('javascript:alert(1)')).toBe('#')
  })
})
