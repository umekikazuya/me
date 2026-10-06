import { describe, expect, it } from 'vitest'
import { currentNav } from './nav.js'

describe('currentNav', () => {
  it('matches pages with or without a trailing slash', () => {
    expect(currentNav('/')?.label).toBe('home')
    expect(currentNav('/articles/')?.label).toBe('writing')
  })

  it('returns undefined for unknown paths', () => {
    expect(currentNav('/nope')).toBeUndefined()
  })
})
