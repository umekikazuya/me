import {
  login as apiLogin,
  logout as apiLogout,
  refreshSession as apiRefreshSession,
  revokeAllSessions as apiRevokeAllSessions,
  changeEmail as apiChangeEmail,
} from '../admin/auth-api.js'
import {
  ApiError,
  type AdminLoginInput,
  type ChangeEmailInput,
  describeApiError,
} from '../admin/types.js'
import { Repository } from './Repository.js'

export type AdminSessionStatus =
  | 'unknown'
  | 'checking'
  | 'authenticated'
  | 'guest'

export interface AuthEventMap {
  'auth:status-change': CustomEvent<{ status: AdminSessionStatus }>
}

/**
 * The public interface for AuthRepository.
 */
export interface IAuthRepository extends EventTarget {
  readonly status: AdminSessionStatus
  readonly loginPending: boolean
  readonly loginError: string
  readonly loginNotice: string
  readonly accountBusyAction: string
  readonly accountError: string
  readonly accountSuccess: string

  addEventListener<K extends keyof AuthEventMap>(
    type: K,
    listener: (e: AuthEventMap[K]) => void,
    options?: boolean | AddEventListenerOptions,
  ): void
  addEventListener(
    type: string,
    callback: EventListenerOrEventListenerObject | null,
    options?: boolean | AddEventListenerOptions,
  ): void

  login(input: AdminLoginInput): Promise<void>
  logout(): Promise<void>
  refreshSession(): Promise<void>
  revokeAllSessions(): Promise<void>
  changeEmail(input: ChangeEmailInput): Promise<void>
  clearLoginNotice(): void
}

export class AuthRepository extends Repository implements IAuthRepository {
  private _status: AdminSessionStatus = 'unknown'
  private _loginPending = false
  private _loginError = ''
  private _loginNotice = ''
  private _accountBusyAction = ''
  private _accountError = ''
  private _accountSuccess = ''

  private sessionBootstrap?: Promise<void>

  get status() {
    return this._status
  }
  get loginPending() {
    return this._loginPending
  }
  get loginError() {
    return this._loginError
  }
  get loginNotice() {
    return this._loginNotice
  }
  get accountBusyAction() {
    return this._accountBusyAction
  }
  get accountError() {
    return this._accountError
  }
  get accountSuccess() {
    return this._accountSuccess
  }

  private dispatchStatusChange() {
    this.dispatchEvent(
      new CustomEvent('auth:status-change', {
        detail: { status: this._status },
      }),
    )
    this.notifyChange()
  }

  async login(input: AdminLoginInput) {
    this._loginPending = true
    this._loginError = ''
    this.notifyChange()
    try {
      await apiLogin(input)
      this._status = 'authenticated'
      this._loginNotice = ''
      this._loginError = ''
      this.dispatchStatusChange()
    } catch (error) {
      this._loginError = describeApiError(error)
      this.notifyChange()
    } finally {
      this._loginPending = false
      this.notifyChange()
    }
  }

  async logout() {
    this._accountBusyAction = 'logout'
    this._accountError = ''
    this._accountSuccess = ''
    this.notifyChange()
    try {
      await apiLogout()
      this._status = 'guest'
      this._loginNotice = 'ログアウトしました。'
      this._accountSuccess = 'ログアウトしました。'
      this.dispatchStatusChange()
    } catch (error) {
      this._accountError = describeApiError(error)
      this.notifyChange()
    } finally {
      this._accountBusyAction = ''
      this.notifyChange()
    }
  }

  async refreshSession() {
    if (this.sessionBootstrap) return this.sessionBootstrap
    this._status = 'checking'
    this.dispatchStatusChange()
    this.sessionBootstrap = (async () => {
      try {
        await apiRefreshSession()
        this._status = 'authenticated'
        this.dispatchStatusChange()
      } catch (error) {
        const isUnauthorized = error instanceof ApiError && error.status === 401
        this._status = 'guest'
        if (isUnauthorized) this._loginNotice = 'セッションが切れました。'
        this.dispatchStatusChange()
      } finally {
        this.sessionBootstrap = undefined
      }
    })()
    return this.sessionBootstrap
  }

  async revokeAllSessions() {
    this._accountBusyAction = 'revoke-sessions'
    this.notifyChange()
    try {
      await apiRevokeAllSessions()
      this._status = 'guest'
      this.dispatchStatusChange()
    } catch (error) {
      this._accountError = describeApiError(error)
      this.notifyChange()
    } finally {
      this._accountBusyAction = ''
      this.notifyChange()
    }
  }

  async changeEmail(input: ChangeEmailInput) {
    this._accountBusyAction = 'change-email'
    this.notifyChange()
    try {
      await apiChangeEmail(input)
      this._accountSuccess = 'メールアドレス変更を送信しました。'
      this.notifyChange()
    } catch (error) {
      this._accountError = describeApiError(error)
      this.notifyChange()
    } finally {
      this._accountBusyAction = ''
      this.notifyChange()
    }
  }

  clearLoginNotice() {
    this._loginNotice = ''
    this.notifyChange()
  }
}
