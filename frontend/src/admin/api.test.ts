import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { apiRequest } from './api.js'
import { resolveApiBasePath } from './api-config.js'

const setApiOrigin = (value?: string) => {
  if (value === undefined) {
    Reflect.deleteProperty(import.meta.env, 'VITE_API_ORIGIN')
    return
  }

  Reflect.set(import.meta.env, 'VITE_API_ORIGIN', value)
}

const jsonResponse = (body: unknown) =>
  new Response(JSON.stringify(body), {
    headers: { 'Content-Type': 'application/json' },
  })

describe('resolveApiBasePath', () => {
  it('falls back to /api when API origin is not configured', () => {
    expect(resolveApiBasePath(undefined)).toBe('/api')
    expect(resolveApiBasePath('   ')).toBe('/api')
  })

  it('removes trailing slashes from a configured API origin', () => {
    expect(resolveApiBasePath('https://api.example.com/')).toBe(
      'https://api.example.com',
    )
    expect(resolveApiBasePath('https://api.example.com///')).toBe(
      'https://api.example.com',
    )
  })
})

describe('apiRequest', () => {
  const originalFetch = globalThis.fetch

  beforeEach(() => {
    globalThis.fetch = vi.fn() as typeof fetch
    setApiOrigin(undefined)
  })

  afterEach(() => {
    globalThis.fetch = originalFetch
    setApiOrigin(undefined)
  })

  it('uses the existing /api base path by default', async () => {
    vi.mocked(globalThis.fetch).mockResolvedValue(jsonResponse({ ok: true }))

    await apiRequest<{ ok: boolean }>('/me')

    expect(globalThis.fetch).toHaveBeenCalledWith(
      '/api/me',
      expect.objectContaining({
        credentials: 'include',
      }),
    )
  })

  it('uses a configured absolute API origin', async () => {
    setApiOrigin('https://api.example.com/')
    vi.mocked(globalThis.fetch).mockResolvedValue(jsonResponse({ ok: true }))

    await apiRequest<{ ok: boolean }>('/me')

    expect(globalThis.fetch).toHaveBeenCalledWith(
      'https://api.example.com/me',
      expect.any(Object),
    )
  })

  it('preserves JSON request behavior while using the resolved API base path', async () => {
    vi.mocked(globalThis.fetch).mockResolvedValue(jsonResponse({ ok: true }))

    await apiRequest<{ ok: boolean }>('/auth/login', {
      method: 'POST',
      body: {
        emailAddress: 'hello@example.com',
        password: 'passw0rd',
      },
    })

    const [, request] = vi.mocked(globalThis.fetch).mock.calls[0]
    const headers = new Headers(request?.headers)

    expect(headers.get('Content-Type')).toBe('application/json')
    expect(headers.get('X-Requested-With')).toBe('XMLHttpRequest')
    expect(request?.body).toBe(
      JSON.stringify({
        emailAddress: 'hello@example.com',
        password: 'passw0rd',
      }),
    )
  })
})
