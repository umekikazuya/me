import { computed, type Signal, signal } from '@lit-labs/signals'
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
  readonly state: Signal.State<IState<AuthData>>

  readonly status: Signal.Computed<AdminSessionStatus>
  readonly loginPending: Signal.Computed<boolean>
  readonly accountBusyAction: Signal.Computed<string>
  readonly error: Signal.Computed<string>
  readonly success: Signal.Computed<string>
  readonly notice: Signal.Computed<string>

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
  private _state = signal<IState<AuthData>>(
    createInitialState(DEFAULT_AUTH_DATA),
  )
  private sessionBootstrap?: Promise<void>

  get state() {
    return this._state
  }

  public status = computed(() => this._state.get().data?.status ?? 'unknown')

  public loginPending = computed(
    () =>
      this._state.get().status === 'loading' &&
      this._state.get().data?.accountBusyAction === 'login',
  )

  public accountBusyAction = computed(
    () => this._state.get().data?.accountBusyAction ?? '',
  )
  public error = computed(() => this._state.get().error?.message ?? '')
  public success = computed(() => this._state.get().data?.accountSuccess ?? '')
  public notice = computed(() => this._state.get().data?.loginNotice ?? '')

  async login(input: AdminLoginInput) {
    const gen = this.nextGeneration()
    this.patchData({ accountBusyAction: 'login' }, 'loading')

    try {
      await apiLogin(input)
      if (!this.isCurrent(gen)) return
      this.updateState(this._state, {
        status: 'success',
        data: { ...DEFAULT_AUTH_DATA, status: 'authenticated' },
      })
      this.dispatchEvent(new CustomEvent('auth:login-success'))
      this.notifyChange()
    } catch (error) {
      if (!this.isCurrent(gen)) return
      this.updateState(this._state, {
        status: 'error',
        error: { code: 'LOGIN_FAILED', message: describeApiError(error) },
        data: { ...this.ensureData(), accountBusyAction: '' },
      })
      this.notifyChange()
    }
  }

  async logout() {
    this.patchData({ accountBusyAction: 'logout', accountSuccess: '' })
    try {
      await apiLogout()
      this.updateState(this._state, {
        status: 'success',
        data: {
          ...DEFAULT_AUTH_DATA,
          status: 'guest',
          loginNotice: 'ログアウトしました。',
          accountSuccess: 'ログアウトしました。',
        },
      })
      this.dispatchEvent(new CustomEvent('auth:logout'))
      this.notifyChange()
    } catch (error) {
      this.updateState(this._state, {
        status: 'error',
        error: { code: 'LOGOUT_FAILED', message: describeApiError(error) },
        data: { ...this.ensureData(), accountBusyAction: '' },
      })
      this.notifyChange()
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
    this.updateState(this._state, {
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
    this.notifyChange()
  }

  async revokeAllSessions() {
    this.patchData({ accountBusyAction: 'revoke-sessions' })
    try {
      await apiRevokeAllSessions()
      this.updateState(this._state, {
        status: 'success',
        data: {
          ...DEFAULT_AUTH_DATA,
          status: 'guest',
          loginNotice: '全セッションを失効させました。',
          accountSuccess: '全セッションを失効させました。',
        },
      })
      this.notifyChange()
    } catch (error) {
      this.updateState(this._state, {
        status: 'error',
        error: { code: 'REVOKE_FAILED', message: describeApiError(error) },
        data: { ...this.ensureData(), accountBusyAction: '' },
      })
      this.notifyChange()
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
      this.updateState(this._state, {
        status: 'error',
        error: {
          code: 'CHANGE_EMAIL_FAILED',
          message: describeApiError(error),
        },
        data: { ...this.ensureData(), accountBusyAction: '' },
      })
      this.notifyChange()
    }
  }

  clearLoginNotice() {
    this.patchData({ loginNotice: '' })
  }

  private ensureData(): AuthData {
    return this._state.get().data ?? DEFAULT_AUTH_DATA
  }

  private patchData(patch: Partial<AuthData>, status?: StateStatus) {
    this.updateState(this._state, {
      status: status ?? this._state.get().status,
      data: { ...this.ensureData(), ...patch },
    })
    this.notifyChange()
  }
}
