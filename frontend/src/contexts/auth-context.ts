import { createContext } from '@lit/context'
import type { ReactiveControllerHost } from 'lit'
import type { AdminLoginInput, ChangeEmailInput } from '../admin/types.js'

export type AdminSessionStatus =
  | 'unknown'
  | 'checking'
  | 'authenticated'
  | 'guest'

export interface AuthController {
  readonly status: AdminSessionStatus
  readonly loginPending: boolean
  readonly loginError: string
  readonly loginNotice: string
  readonly accountBusyAction: string
  readonly accountError: string
  readonly accountSuccess: string

  addHost(host: ReactiveControllerHost): void
  login(input: AdminLoginInput): Promise<void>
  logout(): Promise<void>
  refreshSession(): Promise<void>
  revokeAllSessions(): Promise<void>
  changeEmail(input: ChangeEmailInput): Promise<void>
  clearLoginNotice(): void
}

export const authContext = createContext<AuthController>('auth-context')
