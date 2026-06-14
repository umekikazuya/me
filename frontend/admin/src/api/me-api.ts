import { apiRequest } from './api'
import { type MeProfile, normalizeMeResponse, toMeRequest } from './types'

export const getMe = async () =>
  normalizeMeResponse(await apiRequest<unknown>('/me', { method: 'GET' }))

export const updateMe = async (profile: MeProfile) =>
  normalizeMeResponse(
    await apiRequest<unknown>('/me', {
      method: 'PUT',
      body: toMeRequest(profile),
    }),
  )
