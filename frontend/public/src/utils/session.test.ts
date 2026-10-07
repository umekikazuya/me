import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { Session } from './session.js'

const host = () => ({
  addController: vi.fn(),
  removeController: vi.fn(),
  requestUpdate: vi.fn(),
  updateComplete: Promise.resolve(true),
})

describe('Session', () => {
  beforeEach(() => vi.useFakeTimers())
  afterEach(() => vi.useRealTimers())

  it('shows the next command only after the output and a pause', () => {
    const session = new Session(host())
    session.advance(0, 3, 100)
    expect(session.outputs).toBe(1)
    vi.advanceTimersByTime(549)
    expect(session.commands).toBe(1)
    vi.advanceTimersByTime(1)
    expect(session.commands).toBe(2)
  })

  it('holds on reveal until continueAfter is called', () => {
    const session = new Session(host())
    session.reveal(0)
    vi.advanceTimersByTime(10_000)
    expect(session.commands).toBe(1)
    session.continueAfter(0, 3, 0)
    vi.advanceTimersByTime(450)
    expect(session.commands).toBe(2)
  })

  it('stops after the last command', () => {
    const session = new Session(host())
    session.advance(0, 1, 0)
    vi.advanceTimersByTime(10_000)
    expect(session.commands).toBe(1)
  })
})
