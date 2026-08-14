import { apiRequest } from './api'
import type { ChangeEmailRequest, LoginRequest } from './types'

export const login = (input: LoginRequest) =>
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

export const changeEmail = (input: ChangeEmailRequest) =>
  apiRequest<void>('/auth/email', {
    method: 'PUT',
    body: input,
  })
