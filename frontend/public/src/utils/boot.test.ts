import { describe, expect, it } from 'vitest'
import { easeProgress } from './boot.js'

describe('easeProgress', () => {
  it('moves part of the way toward the target', () => {
    expect(easeProgress(0, 100)).toBeCloseTo(10)
  })

  it('snaps to the target once close enough', () => {
    expect(easeProgress(99.7, 100)).toBe(100)
  })

  it('never overshoots', () => {
    let shown = 0
    for (let i = 0; i < 200; i++) shown = easeProgress(shown, 50)
    expect(shown).toBe(50)
  })
})
