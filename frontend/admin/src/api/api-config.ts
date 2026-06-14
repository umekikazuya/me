export const DEFAULT_API_BASE_PATH = '/api'

export const resolveApiBasePath = (
  apiOrigin = import.meta.env.VITE_API_ORIGIN,
) => {
  const normalizedApiOrigin = apiOrigin?.trim()
  if (!normalizedApiOrigin) return DEFAULT_API_BASE_PATH

  return normalizedApiOrigin.replace(/\/+$/, '')
}
