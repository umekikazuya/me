import { computed, type ReadonlySignal, signal } from '@preact/signals-core'
import { getMe, updateMe } from '../api/me-api.js'
import {
  createEmptyMeProfile,
  describeApiError,
  type MeProfile,
} from '../api/types.js'
import { Repository } from './Repository.js'

export interface ProfileEventMap {
  'profile:public-change': CustomEvent<{ profile: MeProfile | null }>
  'profile:admin-change': CustomEvent<{ profile: MeProfile }>
}

/**
 * The public interface for ProfileRepository.
 */
export interface IProfileRepository extends EventTarget {
  readonly publicProfile: ReadonlySignal<MeProfile | null>
  readonly publicLoading: ReadonlySignal<boolean>
  readonly adminProfile: ReadonlySignal<MeProfile>
  readonly adminLoading: ReadonlySignal<boolean>
  readonly adminSaving: ReadonlySignal<boolean>
  readonly adminLoaded: ReadonlySignal<boolean>
  readonly adminError: ReadonlySignal<string>
  readonly adminSuccess: ReadonlySignal<string>
  readonly adminDirty: ReadonlySignal<boolean>

  addEventListener<K extends keyof ProfileEventMap>(
    type: K,
    listener: (e: ProfileEventMap[K]) => void,
    options?: boolean | AddEventListenerOptions,
  ): void
  addEventListener(
    type: string,
    callback: EventListenerOrEventListenerObject | null,
    options?: boolean | AddEventListenerOptions,
  ): void

  loadPublicProfile(): Promise<void>
  loadAdminProfile(): Promise<void>
  saveAdminProfile(profile: MeProfile): Promise<void>
  setAdminDirty(dirty: boolean): void
}

export class ProfileRepository
  extends Repository
  implements IProfileRepository
{
  private _publicProfile = signal<MeProfile | null>(null)
  private _publicLoading = signal(false)
  private _adminProfile = signal<MeProfile>(createEmptyMeProfile())
  private _adminLoading = signal(false)
  private _adminSaving = signal(false)
  private _adminLoaded = signal(false)
  private _adminError = signal('')
  private _adminSuccess = signal('')
  private _adminDirty = signal(false)

  private _fetchPromise: Promise<MeProfile> | null = null

  readonly publicProfile: ReadonlySignal<MeProfile | null> = computed(
    () => this._publicProfile.value,
  )
  readonly publicLoading: ReadonlySignal<boolean> = computed(
    () => this._publicLoading.value,
  )
  readonly adminProfile: ReadonlySignal<MeProfile> = computed(
    () => this._adminProfile.value,
  )
  readonly adminLoading: ReadonlySignal<boolean> = computed(
    () => this._adminLoading.value,
  )
  readonly adminSaving: ReadonlySignal<boolean> = computed(
    () => this._adminSaving.value,
  )
  readonly adminLoaded: ReadonlySignal<boolean> = computed(
    () => this._adminLoaded.value,
  )
  readonly adminError: ReadonlySignal<string> = computed(
    () => this._adminError.value,
  )
  readonly adminSuccess: ReadonlySignal<string> = computed(
    () => this._adminSuccess.value,
  )
  readonly adminDirty: ReadonlySignal<boolean> = computed(
    () => this._adminDirty.value,
  )

  private notifyPublicChange() {
    this.dispatchEvent(
      new CustomEvent('profile:public-change', {
        detail: { profile: this._publicProfile.value },
      }),
    )
  }

  private notifyAdminChange() {
    this.dispatchEvent(
      new CustomEvent('profile:admin-change', {
        detail: { profile: this._adminProfile.value },
      }),
    )
  }

  async loadPublicProfile() {
    if (this._publicProfile.value || this._publicLoading.value) return
    this._publicLoading.value = true
    try {
      await this._internalFetch()
      this.notifyPublicChange()
    } catch {
      // Fallback handled by components
    } finally {
      this._publicLoading.value = false
    }
  }

  async loadAdminProfile() {
    if (this._adminLoaded.value || this._adminLoading.value) return
    this._adminLoading.value = true
    this._adminError.value = ''
    try {
      await this._internalFetch()
      this._adminLoaded.value = true
      this._adminDirty.value = false
      this.notifyAdminChange()
    } catch (error) {
      this._adminError.value = describeApiError(error)
    } finally {
      this._adminLoading.value = false
    }
  }

  private async _internalFetch() {
    if (this._fetchPromise) return this._fetchPromise
    this._fetchPromise = getMe()
    try {
      const p = await this._fetchPromise
      this._publicProfile.value = p
      this._adminProfile.value = p
      return p
    } finally {
      this._fetchPromise = null
    }
  }

  async saveAdminProfile(profile: MeProfile) {
    this._adminSaving.value = true
    this._adminError.value = ''
    this._adminSuccess.value = ''
    try {
      const saved = await updateMe(profile)
      this._adminProfile.value = saved
      this._publicProfile.value = saved
      this._adminLoaded.value = true
      this._adminDirty.value = false
      this._adminSuccess.value = 'プロフィールを更新しました。'
      this.notifyAdminChange()
      this.notifyPublicChange()
    } catch (error) {
      this._adminError.value = describeApiError(error)
    } finally {
      this._adminSaving.value = false
    }
  }

  setAdminDirty(dirty: boolean) {
    this._adminDirty.value = dirty
    if (dirty) {
      this._adminSuccess.value = ''
    }
  }
}
