import { describe, expect, it } from 'vitest'
import { gitLogCommand, shortHash } from './page-articles.js'

describe('gitLogCommand', () => {
  it('is plain git log without filters', () => {
    expect(gitLogCommand('', [])).toBe('git log --oneline')
  })

  it('adds grep and tag pathspecs', () => {
    expect(gitLogCommand('neovim', ['Go', 'Git'])).toBe(
      'git log --oneline --grep="neovim" -- tag:Go tag:Git',
    )
  })
})

describe('shortHash', () => {
  it('takes the first 7 characters and pads short ids', () => {
    expect(shortHash('5fa8a38049328e301cb4')).toBe('5fa8a38')
    expect(shortHash('a1')).toBe('a1     ')
  })
})
