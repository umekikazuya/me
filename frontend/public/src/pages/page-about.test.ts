import type { components } from '@me/types'
import { describe, expect, it } from 'vitest'
import { profileRows } from './page-about.js'

const base = {
  displayName: 'Kazuya Umeki',
  role: 'Web Creator',
  location: 'Fukuoka, Japan',
  updatedAt: '2026-01-01T00:00:00Z',
} as components['schemas']['MeResponse']

describe('profileRows', () => {
  it('is empty when the API omits every list', () => {
    expect(profileRows(base)).toEqual([])
  })

  it('labels only the first skill line and keeps sort order', () => {
    expect(
      profileRows({
        ...base,
        skills: [
          { category: 'Stack', items: ['Lit', 'React'], sortOrder: 2 },
          { category: 'Lang', items: ['Go', 'TypeScript'], sortOrder: 1 },
        ],
        certifications: [{ name: 'AWS SAA', year: 2024 }],
        likes: ['Mr.Children'],
      }),
    ).toEqual([
      ['skills', 'Go, TypeScript'],
      ['', 'Lit, React'],
      ['certs', 'AWS SAA'],
      ['likes', 'Mr.Children'],
    ])
  })

  it('drops skill groups without items', () => {
    expect(
      profileRows({
        ...base,
        skills: [{ category: 'Lang', items: [], sortOrder: 1 }],
      }),
    ).toEqual([])
  })
})
