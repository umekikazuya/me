import { apiRequest } from './api.js'
import type { AdminLoginInput, ChangeEmailInput } from './types.js'

export const login = (input: AdminLoginInput) =>
  apiRequest<void>('/auth/login', {
    method: 'POST',
    body: input,
  })

export const refreshSession = () =>
  apiRequest<void>('/auth/refresh', {
    method: 'POST',
  })

export const logout = () =>
  apiRequest<void>('/auth/logout', {
    method: 'POST',
  })

export const revokeAllSessions = () =>
  apiRequest<void>('/auth/sessions', {
    method: 'DELETE',
  })

export const changeEmail = (input: ChangeEmailInput) =>
  apiRequest<void>('/auth/email', {
    method: 'PUT',
    body: input,
  })
