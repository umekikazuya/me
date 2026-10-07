import { describe, expect, it } from 'vitest'
import { localTime, moodAt } from './clock.js'

describe('localTime', () => {
  it('formats the time in Tokyo regardless of the viewer zone', () => {
    expect(localTime(new Date('2026-10-06T17:41:00Z'))).toEqual({
      text: 'Wed 02:41 JST · Tokyo',
      mood: 'probably asleep',
    })
  })
})

describe('moodAt', () => {
  it('switches at 7, 10 and 19', () => {
    expect([6, 7, 10, 19].map(moodAt)).toEqual([
      'probably asleep',
      'on the train',
      'at work',
      'writing code, maybe',
    ])
  })
})
