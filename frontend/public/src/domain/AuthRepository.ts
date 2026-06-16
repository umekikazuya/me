import {
  changeEmail as apiChangeEmail,
  login as apiLogin,
  logout as apiLogout,
  refreshSession as apiRefreshSession,
  revokeAllSessions as apiRevokeAllSessions,
} from '../api/auth-api.js'
import {
  type AdminLoginInput,
  ApiError,
  type ChangeEmailInput,
  describeApiError,
} from '../api/types.js'
import {
  createInitialState,
  type IState,
  Repository,
  type StateStatus,
} from './Repository.js'

export type AdminSessionStatus =
  | 'unknown'
  | 'checking'
  | 'authenticated'
  | 'guest'

export interface AuthData {
  status: AdminSessionStatus
  loginNotice: string
  accountBusyAction: string
  accountSuccess: string
}

export interface AuthEventMap {
  'auth:login-success': CustomEvent<void>
  'auth:logout': CustomEvent<void>
}

/**
 * The public interface for AuthRepository.
 */
export interface IAuthRepository extends EventTarget {
  readonly status: AdminSessionStatus
  readonly loginPending: boolean
  readonly accountBusyAction: string
  readonly error: string
  readonly success: string
  readonly notice: string

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

const DEFAULT_AUTH_DATA: AuthData = {
  status: 'unknown',
  loginNotice: '',
  accountBusyAction: '',
  accountSuccess: '',
}

export class AuthRepository extends Repository implements IAuthRepository {
  private _state: IState<AuthData> = createInitialState(DEFAULT_AUTH_DATA)
  private sessionBootstrap?: Promise<void>

  get status(): AdminSessionStatus {
    return this._state.data?.status ?? 'unknown'
  }

  get loginPending(): boolean {
    return (
      this._state.status === 'loading' &&
      this._state.data?.accountBusyAction === 'login'
    )
  }

  get accountBusyAction(): string {
    return this._state.data?.accountBusyAction ?? ''
  }

  get error(): string {
    return this._state.error?.message ?? ''
  }

  get success(): string {
    return this._state.data?.accountSuccess ?? ''
  }

  get notice(): string {
    return this._state.data?.loginNotice ?? ''
  }

  private setState(patch: Partial<IState<AuthData>>): void {
    this._state = { ...this._state, ...patch }
    this.emitChange()
  }

  async login(input: AdminLoginInput) {
    const gen = this.nextGeneration()
    this.patchData({ accountBusyAction: 'login' }, 'loading')

    try {
      await apiLogin(input)
      if (!this.isCurrent(gen)) return
      this.setState({
        status: 'success',
        data: { ...DEFAULT_AUTH_DATA, status: 'authenticated' },
      })
      this.dispatchEvent(new CustomEvent('auth:login-success'))
    } catch (error) {
      if (!this.isCurrent(gen)) return
      this.setState({
        status: 'error',
        error: { code: 'LOGIN_FAILED', message: describeApiError(error) },
        data: { ...this.ensureData(), accountBusyAction: '' },
      })
    }
  }

  async logout() {
    this.patchData({ accountBusyAction: 'logout', accountSuccess: '' })
    try {
      await apiLogout()
      this.setState({
        status: 'success',
        data: {
          ...DEFAULT_AUTH_DATA,
          status: 'guest',
          loginNotice: 'ログアウトしました。',
          accountSuccess: 'ログアウトしました。',
        },
      })
      this.dispatchEvent(new CustomEvent('auth:logout'))
    } catch (error) {
      this.setState({
        status: 'error',
        error: { code: 'LOGOUT_FAILED', message: describeApiError(error) },
        data: { ...this.ensureData(), accountBusyAction: '' },
      })
    }
  }

  async refreshSession() {
    if (this.sessionBootstrap) return this.sessionBootstrap
    const gen = this.nextGeneration()
    this.patchData({ status: 'checking' }, 'loading')

    this.sessionBootstrap = (async () => {
      try {
        await apiRefreshSession()
        if (this.isCurrent(gen)) {
          this.patchData({ status: 'authenticated' }, 'success')
        }
      } catch (error) {
        if (this.isCurrent(gen)) {
          this.handleRefreshError(error)
        }
      } finally {
        this.sessionBootstrap = undefined
      }
    })()
    return this.sessionBootstrap
  }

  private handleRefreshError(error: unknown) {
    const isUnauthorized = error instanceof ApiError && error.status === 401
    this.setState({
      status: isUnauthorized ? 'success' : 'error',
      error: isUnauthorized
        ? null
        : { code: 'REFRESH_FAILED', message: describeApiError(error) },
      data: {
        ...this.ensureData(),
        status: 'guest',
        loginNotice: isUnauthorized ? 'セッションが切れました。' : '',
      },
    })
  }

  async revokeAllSessions() {
    this.patchData({ accountBusyAction: 'revoke-sessions' })
    try {
      await apiRevokeAllSessions()
      this.setState({
        status: 'success',
        data: {
          ...DEFAULT_AUTH_DATA,
          status: 'guest',
          loginNotice: '全セッションを失効させました。',
          accountSuccess: '全セッションを失効させました。',
        },
      })
    } catch (error) {
      this.setState({
        status: 'error',
        error: { code: 'REVOKE_FAILED', message: describeApiError(error) },
        data: { ...this.ensureData(), accountBusyAction: '' },
      })
    }
  }

  async changeEmail(input: ChangeEmailInput) {
    this.patchData({ accountBusyAction: 'change-email', accountSuccess: '' })
    try {
      await apiChangeEmail(input)
      this.patchData({
        accountBusyAction: '',
        accountSuccess: 'メールアドレス変更を送信しました。',
      })
    } catch (error) {
      this.setState({
        status: 'error',
        error: {
          code: 'CHANGE_EMAIL_FAILED',
          message: describeApiError(error),
        },
        data: { ...this.ensureData(), accountBusyAction: '' },
      })
    }
  }

  clearLoginNotice() {
    this.patchData({ loginNotice: '' })
  }

  private ensureData(): AuthData {
    return this._state.data ?? DEFAULT_AUTH_DATA
  }

  private patchData(patch: Partial<AuthData>, status?: StateStatus) {
    this.setState({
      status: status ?? this._state.status,
      data: { ...this.ensureData(), ...patch },
    })
  }
}
