import { ApiError, type ProblemDetail } from './types.js'

const API_BASE_PATH = '/api'
const REQUESTED_WITH_HEADER = 'XMLHttpRequest'

interface ApiRequestOptions extends Omit<RequestInit, 'body' | 'headers'> {
  body?: BodyInit | object
  headers?: HeadersInit
}

const isPlainObject = (value: unknown): value is Record<string, unknown> =>
  typeof value === 'object' &&
  value !== null &&
  !(value instanceof FormData) &&
  !(value instanceof URLSearchParams) &&
  !(value instanceof Blob) &&
  !(value instanceof ArrayBuffer)

const parseJson = async (response: Response) => {
  const contentType = response.headers.get('content-type') ?? ''
  if (!contentType.includes('application/json')) return undefined

  return response.json()
}

export async function apiRequest<T>(
  path: string,
  options: ApiRequestOptions = {},
): Promise<T> {
  const headers = new Headers(options.headers)
  headers.set('Accept', 'application/json')
  headers.set('X-Requested-With', REQUESTED_WITH_HEADER)

  let body = options.body as BodyInit | undefined
  if (isPlainObject(options.body)) {
    headers.set('Content-Type', 'application/json')
    body = JSON.stringify(options.body)
  }

  const response = await fetch(`${API_BASE_PATH}${path}`, {
    ...options,
    body,
    headers,
    credentials: 'include',
  })

  if (!response.ok) {
    const problem = (await parseJson(response)) as ProblemDetail | undefined
    throw new ApiError(
      problem?.message ||
        problem?.detail ||
        problem?.title ||
        `API request failed with status ${response.status}`,
      response.status,
      problem,
    )
  }

  if (response.status === 204) return undefined as T
  return (await parseJson(response)) as T
}
