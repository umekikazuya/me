import { describe, expect, it } from 'vitest'
import { toHomePath, toRuler } from './status-line.js'

describe('toHomePath', () => {
  it('maps the root to ~', () => {
    expect(toHomePath('/')).toBe('~')
  })

  it('prefixes nested paths with ~ and drops trailing slashes', () => {
    expect(toHomePath('/articles/')).toBe('~/articles')
  })
})

describe('toRuler', () => {
  it('reports top, bot and a percentage in between', () => {
    expect(toRuler(0, 1000)).toBe('top')
    expect(toRuler(1000, 1000)).toBe('bot')
    expect(toRuler(250, 1000)).toBe('25%')
  })

  it('treats a page without scroll as top', () => {
    expect(toRuler(0, 0)).toBe('top')
  })
})
