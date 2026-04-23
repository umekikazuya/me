import type { ReactiveController, ReactiveControllerHost } from 'lit'
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
import type {
  AdminSessionStatus,
  AuthController as IAuthController,
} from '../contexts/auth-context.js'

export class AuthController implements ReactiveController, IAuthController {
  private hosts: Set<ReactiveControllerHost> = new Set()

  private _status: AdminSessionStatus = 'unknown'
  private _loginPending = false
  private _loginError = ''
  private _loginNotice = ''
  private _accountBusyAction = ''
  private _accountError = ''
  private _accountSuccess = ''

  private sessionBootstrap?: Promise<void>

  constructor(host: ReactiveControllerHost) {
    this.addHost(host)
  }

  addHost(host: ReactiveControllerHost) {
    this.hosts.add(host)
    host.addController(this)
  }

  hostConnected() {}

  private requestUpdate() {
    for (const host of this.hosts) {
      host.requestUpdate()
    }
  }

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
    this.requestUpdate()

    try {
      await login(input)
      this._status = 'authenticated'
      this._loginNotice = ''
      this._loginError = ''
    } catch (error) {
      this._loginError = describeApiError(error)
    } finally {
      this._loginPending = false
      this.requestUpdate()
    }
  }

  async logout() {
    this._accountBusyAction = 'logout'
    this._accountError = ''
    this._accountSuccess = ''
    this.requestUpdate()

    try {
      await logout()
      this._status = 'guest'
      this._loginNotice = 'ログアウトしました。'
      this._accountSuccess = 'ログアウトしました。'
    } catch (error) {
      this._accountError = describeApiError(error)
    } finally {
      this._accountBusyAction = ''
      this.requestUpdate()
    }
  }

  async refreshSession() {
    if (this.sessionBootstrap) return this.sessionBootstrap

    this._status = 'checking'
    this.requestUpdate()

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
        this.requestUpdate()
      }
    })()

    return this.sessionBootstrap
  }

  async revokeAllSessions() {
    this._accountBusyAction = 'revoke-sessions'
    this._accountError = ''
    this._accountSuccess = ''
    this.requestUpdate()

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
      this.requestUpdate()
    }
  }

  async changeEmail(input: ChangeEmailInput) {
    this._accountBusyAction = 'change-email'
    this._accountError = ''
    this._accountSuccess = ''
    this.requestUpdate()

    try {
      await changeEmail(input)
      this._accountSuccess = 'メールアドレス変更を送信しました。'
    } catch (error) {
      this._accountError = describeApiError(error)
    } finally {
      this._accountBusyAction = ''
      this.requestUpdate()
    }
  }

  clearLoginNotice() {
    this._loginNotice = ''
    this.requestUpdate()
  }
}
