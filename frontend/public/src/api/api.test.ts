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

  it('returns undefined for a 204 response', async () => {
    vi.mocked(globalThis.fetch).mockResolvedValue(
      new Response(null, { status: 204 }),
    )

    await expect(
      apiRequest('/auth/logout', { method: 'POST' }),
    ).resolves.toBeUndefined()
  })

  it('throws an ApiError carrying the problem details', async () => {
    vi.mocked(globalThis.fetch).mockResolvedValue(
      new Response(JSON.stringify({ detail: 'bad input', status: 400 }), {
        status: 400,
        headers: { 'Content-Type': 'application/problem+json' },
      }),
    )

    await expect(apiRequest('/me')).rejects.toMatchObject({
      name: 'ApiError',
      status: 400,
      message: 'bad input',
    })
  })

  // プロキシやロードバランサは JSON でないエラーボディを返すことがある。
  it('throws an ApiError with a status fallback for a non-JSON error body', async () => {
    vi.mocked(globalThis.fetch).mockResolvedValue(
      new Response('<html>Bad Gateway</html>', {
        status: 502,
        headers: { 'Content-Type': 'text/html' },
      }),
    )

    await expect(apiRequest('/me')).rejects.toMatchObject({
      name: 'ApiError',
      status: 502,
      message: 'API request failed with status 502',
    })
  })
})
