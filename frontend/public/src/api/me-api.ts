import { apiRequest } from './api.js'
import { normalizeMeResponse } from './types.js'

export const getMe = async () =>
  normalizeMeResponse(await apiRequest<unknown>('/me', { method: 'GET' }))
