import type { components } from '@me/types'
import { describe, expect, it } from 'vitest'
import { whoamiLines } from './page-top.js'

const base: components['schemas']['MeResponse'] = {
  displayName: 'Kazuya Umeki',
  role: 'Software Engineer',
  location: 'Fukuoka, Japan',
  skills: [
    { category: 'infra', items: ['Docker'], sortOrder: 2 },
    { category: 'lang', items: ['Go', 'TypeScript', 'PHP'], sortOrder: 1 },
  ],
  certifications: [],
  experiences: [],
  links: [],
  likes: [],
  updatedAt: '2026-01-01T00:00:00Z',
}

describe('whoamiLines', () => {
  it('lists name, role and the top skills in sort order', () => {
    expect(whoamiLines(base)).toEqual([
      'Kazuya Umeki',
      'Software Engineer, Fukuoka, Japan',
      'Go / TypeScript / Docker',
    ])
  })

  it('tolerates skills omitted by the API', () => {
    const { skills: _, ...rest } = base
    expect(whoamiLines(rest as typeof base)).toHaveLength(2)
  })

  it('adds the Japanese name and omits empty skills', () => {
    expect(
      whoamiLines({ ...base, displayJa: '梅木 和弥', skills: [] }),
    ).toEqual(['Kazuya Umeki (梅木 和弥)', 'Software Engineer, Fukuoka, Japan'])
  })
})
