import { type ReadonlySignal, signal } from '@preact/signals-core'
import { describe, expect, it, vi } from 'vitest'
import type { AdminLoginInput, ChangeEmailInput } from '../../../admin/types.js'
import type {
  AdminSessionStatus,
  AuthData,
  IAuthRepository,
} from '../../../domain/AuthRepository.js'
import type { IState } from '../../../domain/Repository.js'
import type { MeAdminAuthBoundary } from './me-admin-auth-boundary.ts'
import './me-admin-auth-boundary.ts'

class TestAuthRepository extends EventTarget implements IAuthRepository {
  readonly state = signal<IState<AuthData>>({
    status: 'success',
    data: {
      status: 'authenticated',
      loginNotice: '',
      accountBusyAction: '',
      accountSuccess: '',
    },
    error: null,
  })
  readonly status: ReadonlySignal<AdminSessionStatus>
  readonly loginPending = signal(false)
  readonly accountBusyAction = signal('')
  readonly error = signal('')
  readonly success = signal('')
  readonly notice = signal('')
  refreshSession = vi.fn(async () => {})

  constructor(status: AdminSessionStatus = 'authenticated') {
    super()
    this.status = signal(status)
  }

  async login(_input: AdminLoginInput) {}

  async logout() {}

  async revokeAllSessions() {}

  async changeEmail(_input: ChangeEmailInput) {}

  clearLoginNotice() {}
}

describe('MeAdminAuthBoundary', () => {
  it('replaces the previous RepositoryObserver cleanly when authRepo changes', async () => {
    const el = document.createElement(
      'me-admin-auth-boundary',
    ) as MeAdminAuthBoundary
    const removeControllerSpy = vi.spyOn(el, 'removeController')
    const firstRepo = new TestAuthRepository()
    const secondRepo = new TestAuthRepository()

    el.authRepo = firstRepo
    document.body.appendChild(el)
    await el.updateComplete
    const firstObserver = Reflect.get(el, '_observer')
    expect(firstObserver).toBeDefined()

    const disconnectSpy = vi.spyOn(firstObserver, 'disconnect')

    el.authRepo = secondRepo
    await el.updateComplete

    expect(removeControllerSpy).toHaveBeenCalledWith(firstObserver)
    expect(disconnectSpy).toHaveBeenCalledOnce()
    expect(Reflect.get(el, '_observer')).not.toBe(firstObserver)

    el.remove()
  })
})
