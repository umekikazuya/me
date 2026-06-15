import { computed, type ReadonlySignal, signal } from '@preact/signals-core'
import { getMe } from '../api/me-api.js'
import { describeApiError, type MeProfile } from '../api/types.js'
import { createInitialState, type IState, Repository } from './Repository.js'

/**
 * The public interface for ProfileRepository.
 * Scoped to public-facing profile data only.
 */
export interface IProfileRepository {
  readonly state: ReadonlySignal<IState<MeProfile>>
  readonly profile: ReadonlySignal<MeProfile | null>
  readonly isLoading: ReadonlySignal<boolean>
  readonly error: ReadonlySignal<string>

  loadProfile(): Promise<void>
}

export class ProfileRepository
  extends Repository
  implements IProfileRepository
{
  private _state = signal<IState<MeProfile>>(createInitialState<MeProfile>())

  readonly state: ReadonlySignal<IState<MeProfile>> = computed(
    () => this._state.value,
  )
  readonly profile: ReadonlySignal<MeProfile | null> = computed(
    () => this._state.value.data,
  )
  readonly isLoading: ReadonlySignal<boolean> = computed(
    () => this._state.value.status === 'loading',
  )
  readonly error: ReadonlySignal<string> = computed(
    () => this._state.value.error?.message ?? '',
  )

  private _fetchPromise: Promise<MeProfile> | null = null

  async loadProfile() {
    if (this._state.value.data || this._state.value.status === 'loading') return

    const gen = this.nextGeneration()
    this.updateState(this._state, { status: 'loading', error: null })

    try {
      if (!this._fetchPromise) {
        this._fetchPromise = getMe()
      }
      const profile = await this._fetchPromise
      if (!this.isCurrent(gen)) return
      this.updateState(this._state, { status: 'success', data: profile })
    } catch (error) {
      if (!this.isCurrent(gen)) return
      this.updateState(this._state, {
        status: 'error',
        error: { code: 'LOAD_FAILED', message: describeApiError(error) },
      })
    } finally {
      this._fetchPromise = null
    }
  }
}
