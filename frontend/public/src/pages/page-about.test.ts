import type { components } from '@me/types'
import { describe, expect, it } from 'vitest'
import { aboutSections, outputDuration } from './page-about.js'

const base = {
  displayName: 'Kazuya Umeki',
  role: 'Web Creator',
  location: 'Fukuoka, Japan',
  updatedAt: '2026-01-01T00:00:00Z',
} as components['schemas']['MeResponse']

describe('aboutSections', () => {
  it('only runs whoami when the API omits every list', () => {
    expect(aboutSections(base)).toEqual([
      {
        kind: 'whoami',
        command: 'whoami',
        text: 'Kazuya Umeki — Web Creator, Fukuoka, Japan',
      },
    ])
  })

  it('runs the commands in a fixed order and flattens skills by sortOrder', () => {
    const sections = aboutSections({
      ...base,
      skills: [
        { category: 'Stack', items: ['Lit'], sortOrder: 2 },
        { category: 'Lang', items: ['Go', 'TypeScript'], sortOrder: 1 },
      ],
      certifications: [{ name: 'GCP ACE', year: 2025 }],
      experiences: [
        { company: 'A', startYear: 2019, endYear: 2022 },
        { company: 'B', startYear: 2022 },
      ],
      likes: ['Mr.Children'],
    })
    expect(sections.map((s) => s.command)).toEqual([
      'whoami',
      'ls skills/',
      'ls certs/',
      'history',
      'cat likes.txt',
    ])
    expect(sections[1]).toMatchObject({ items: ['Go', 'TypeScript', 'Lit'] })
    expect(sections[3]).toMatchObject({
      experiences: [{ company: 'B' }, { company: 'A' }],
    })
  })

  it('skips skills when every group is empty', () => {
    const sections = aboutSections({
      ...base,
      skills: [{ category: 'Lang', items: [], sortOrder: 1 }],
    })
    expect(sections.map((s) => s.kind)).toEqual(['whoami'])
  })
})

describe('outputDuration', () => {
  it('waits for the graph to finish drawing', () => {
    const [, graph] = aboutSections({
      ...base,
      experiences: [
        { company: 'A', startYear: 2019, endYear: 2022 },
        { company: 'B', startYear: 2022 },
      ],
    })
    expect(outputDuration(graph)).toBeGreaterThan(2 * 140)
  })
})
