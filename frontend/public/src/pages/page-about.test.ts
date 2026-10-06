import type { components } from '@me/types'
import { describe, expect, it } from 'vitest'
import { fetchRows } from './page-about.js'

const base = {
  displayName: 'Kazuya Umeki',
  role: 'Web Creator',
  location: 'Fukuoka, Japan',
  updatedAt: '2026-01-01T00:00:00Z',
} as components['schemas']['MeResponse']

describe('fetchRows', () => {
  it('shows role and location when the API omits empty arrays', () => {
    expect(fetchRows(base)).toEqual([
      ['Role', 'Web Creator'],
      ['Location', 'Fukuoka, Japan'],
    ])
  })

  it('adds skill groups in sort order, then certs and likes', () => {
    expect(
      fetchRows({
        ...base,
        skills: [
          { category: 'Stack', items: ['Lit', 'React'], sortOrder: 2 },
          { category: 'Lang', items: ['Go', 'TypeScript'], sortOrder: 1 },
        ],
        certifications: [{ name: 'AWS SAA', year: 2024 }],
        likes: ['Mr.Children'],
      }).map(([key]) => key),
    ).toEqual(['Role', 'Location', 'Lang', 'Stack', 'Certs', 'Likes'])
  })
})
