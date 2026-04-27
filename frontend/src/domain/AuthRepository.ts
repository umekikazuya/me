import { computed, type ReadonlySignal, signal } from '@preact/signals-core'
import {
  changeEmail as apiChangeEmail,
  login as apiLogin,
  logout as apiLogout,
  refreshSession as apiRefreshSession,
  revokeAllSessions as apiRevokeAllSessions,
} from '../admin/auth-api.js'
import {
  type AdminLoginInput,
  ApiError,
  type ChangeEmailInput,
  describeApiError,
} from '../admin/types.js'
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
  readonly state: ReadonlySignal<IState<AuthData>>

  readonly status: ReadonlySignal<AdminSessionStatus>
  readonly loginPending: ReadonlySignal<boolean>
  readonly accountBusyAction: ReadonlySignal<string>
  readonly error: ReadonlySignal<string>
  readonly success: ReadonlySignal<string>
  readonly notice: ReadonlySignal<string>

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

  public status = computed(() => this._state.value.data?.status ?? 'unknown')

  public loginPending = computed(
    () =>
      this._state.value.status === 'loading' &&
      this._state.value.data?.accountBusyAction === 'login',
  )

  public accountBusyAction = computed(
    () => this._state.value.data?.accountBusyAction ?? '',
  )
  public error = computed(() => this._state.value.error?.message ?? '')
  public success = computed(() => this._state.value.data?.accountSuccess ?? '')
  public notice = computed(() => this._state.value.data?.loginNotice ?? '')

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
    } catch (error) {
      if (!this.isCurrent(gen)) return
      this.updateState(this._state, {
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
    } catch (error) {
      this.updateState(this._state, {
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
    } catch (error) {
      this.updateState(this._state, {
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
      this.updateState(this._state, {
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
    return this._state.value.data ?? DEFAULT_AUTH_DATA
  }

  private patchData(patch: Partial<AuthData>, status?: StateStatus) {
    this.updateState(this._state, {
      status: status ?? this._state.value.status,
      data: { ...this.ensureData(), ...patch },
    })
  }
}
