import type { components } from '@me/types'
import { getMe } from '../api/me-api.js'
import { describeApiError } from '../api/types.js'
import { createInitialState, type IState, Repository } from './Repository.js'

/**
 * The public interface for ProfileRepository.
 * Scoped to public-facing profile data only.
 */
export interface IProfileRepository extends EventTarget {
  readonly profile: components['schemas']['MeResponse'] | null
  readonly isLoading: boolean
  readonly error: string

  addEventListener(
    type: 'change',
    listener: (e: Event) => void,
    options?: boolean | AddEventListenerOptions,
  ): void
  addEventListener(
    type: string,
    callback: EventListenerOrEventListenerObject | null,
    options?: boolean | AddEventListenerOptions,
  ): void

  loadProfile(): Promise<void>
}

export class ProfileRepository
  extends Repository
  implements IProfileRepository
{
  private _state: IState<components['schemas']['MeResponse']> =
    createInitialState<components['schemas']['MeResponse']>()
  private _fetchPromise: Promise<components['schemas']['MeResponse']> | null =
    null

  get profile(): components['schemas']['MeResponse'] | null {
    return this._state.data
  }

  get isLoading(): boolean {
    return this._state.status === 'loading'
  }

  get error(): string {
    return this._state.error?.message ?? ''
  }

  private setState(
    patch: Partial<IState<components['schemas']['MeResponse']>>,
  ): void {
    this._state = { ...this._state, ...patch }
    this.emitChange()
  }

  async loadProfile() {
    if (this._state.data || this._state.status === 'loading') return

    const gen = this.nextGeneration()
    this.setState({ status: 'loading', error: null })

    try {
      if (!this._fetchPromise) {
        this._fetchPromise = getMe()
      }
      const profile = await this._fetchPromise
      if (!this.isCurrent(gen)) return
      this.setState({ status: 'success', data: profile })
    } catch (error) {
      if (!this.isCurrent(gen)) return
      this.setState({
        status: 'error',
        error: { code: 'LOAD_FAILED', message: describeApiError(error) },
      })
    } finally {
      this._fetchPromise = null
    }
  }
}
