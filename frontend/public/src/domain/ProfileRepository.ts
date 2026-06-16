import { computed, Signal, signal } from '@lit-labs/signals'
import { getMe } from '../api/me-api.js'
import { describeApiError, type MeProfile } from '../api/types.js'
import { createInitialState, type IState, Repository } from './Repository.js'

/**
 * The public interface for ProfileRepository.
 * Scoped to public-facing profile data only.
 */
export interface IProfileRepository {
  readonly state: Signal.Computed<IState<MeProfile>>
  readonly profile: Signal.Computed<MeProfile | null>
  readonly isLoading: Signal.Computed<boolean>
  readonly error: Signal.Computed<string>

  loadProfile(): Promise<void>
}

export class ProfileRepository
  extends Repository
  implements IProfileRepository
{
  private _state = signal<IState<MeProfile>>(createInitialState<MeProfile>())

  readonly state: Signal.Computed<IState<MeProfile>> = computed(() =>
    this._state.get(),
  )
  readonly profile: Signal.Computed<MeProfile | null> = computed(
    () => this._state.get().data,
  )
  readonly isLoading: Signal.Computed<boolean> = computed(
    () => this._state.get().status === 'loading',
  )
  readonly error: Signal.Computed<string> = computed(
    () => this._state.get().error?.message ?? '',
  )

  private _fetchPromise: Promise<MeProfile> | null = null

  async loadProfile() {
    if (this._state.get().data || this._state.get().status === 'loading') return

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
