import type { components } from '@me/types'
import { apiRequest } from './api.js'

export const getMe = () =>
  apiRequest<components['schemas']['MeResponse']>('/me', { method: 'GET' })
