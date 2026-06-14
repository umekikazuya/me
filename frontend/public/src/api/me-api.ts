import { apiRequest } from './api.js'
import { type MeProfile, normalizeMeResponse, toMeRequest } from './types.js'

export const getMe = async () =>
  normalizeMeResponse(await apiRequest<unknown>('/me', { method: 'GET' }))

export const updateMe = async (profile: MeProfile) =>
  normalizeMeResponse(
    await apiRequest<unknown>('/me', {
      method: 'PUT',
      body: toMeRequest(profile),
    }),
  )
