import {
  changeEmail,
  login,
  logout,
  refreshSession,
  revokeAllSessions,
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

export interface IAuthRepository extends Repository {
  readonly status: AdminSessionStatus
  readonly loginPending: boolean
  readonly loginError: string
  readonly loginNotice: string
  readonly accountBusyAction: string
  readonly accountError: string
  readonly accountSuccess: string

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

  async login(input: AdminLoginInput) {
    this._loginPending = true
    this._loginError = ''
    this.notifyChange()

    try {
      await login(input)
      this._status = 'authenticated'
      this._loginNotice = ''
      this._loginError = ''
    } catch (error) {
      this._loginError = describeApiError(error)
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
      await logout()
      this._status = 'guest'
      this._loginNotice = 'ログアウトしました。'
      this._accountSuccess = 'ログアウトしました。'
    } catch (error) {
      this._accountError = describeApiError(error)
    } finally {
      this._accountBusyAction = ''
      this.notifyChange()
    }
  }

  async refreshSession() {
    if (this.sessionBootstrap) return this.sessionBootstrap

    this._status = 'checking'
    this.notifyChange()

    this.sessionBootstrap = (async () => {
      try {
        await refreshSession()
        this._status = 'authenticated'
        this._loginError = ''
        this._loginNotice = ''
      } catch (error) {
        const isUnauthorized = error instanceof ApiError && error.status === 401
        this._status = 'guest'
        if (isUnauthorized) {
          this._loginNotice =
            'セッションの有効期限が切れました。再度ログインしてください。'
        }
        if (!isUnauthorized) {
          this._loginError = describeApiError(error)
        }
      } finally {
        this.sessionBootstrap = undefined
        this.notifyChange()
      }
    })()

    return this.sessionBootstrap
  }

  async revokeAllSessions() {
    this._accountBusyAction = 'revoke-sessions'
    this._accountError = ''
    this._accountSuccess = ''
    this.notifyChange()

    try {
      await revokeAllSessions()
      this._status = 'guest'
      this._loginNotice =
        '全セッションを終了しました。必要に応じて再度ログインしてください。'
      this._accountSuccess = '全セッションを失効させました。'
    } catch (error) {
      this._accountError = describeApiError(error)
    } finally {
      this._accountBusyAction = ''
      this.notifyChange()
    }
  }

  async changeEmail(input: ChangeEmailInput) {
    this._accountBusyAction = 'change-email'
    this._accountError = ''
    this._accountSuccess = ''
    this.notifyChange()

    try {
      await changeEmail(input)
      this._accountSuccess = 'メールアドレス変更を送信しました。'
    } catch (error) {
      this._accountError = describeApiError(error)
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
