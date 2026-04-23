import { getMe, updateMe } from '../admin/me-api.js'
import {
  createEmptyMeProfile,
  describeApiError,
  type MeProfile,
} from '../admin/types.js'
import { Repository } from './Repository.js'

export interface IProfileRepository extends Repository {
  readonly publicProfile: MeProfile | null
  readonly publicLoading: boolean
  readonly adminProfile: MeProfile
  readonly adminLoading: boolean
  readonly adminSaving: boolean
  readonly adminLoaded: boolean
  readonly adminError: string
  readonly adminSuccess: string
  readonly adminDirty: boolean

  loadPublicProfile(): Promise<void>
  loadAdminProfile(): Promise<void>
  saveAdminProfile(profile: MeProfile): Promise<void>
  setAdminDirty(dirty: boolean): void
}

export class ProfileRepository
  extends Repository
  implements IProfileRepository
{
  private _publicProfile: MeProfile | null = null
  private _publicLoading = false
  private _adminProfile = createEmptyMeProfile()
  private _adminLoading = false
  private _adminSaving = false
  private _adminLoaded = false
  private _adminError = ''
  private _adminSuccess = ''
  private _adminDirty = false

  private _fetchPromise: Promise<MeProfile> | null = null

  get publicProfile() {
    return this._publicProfile
  }
  get publicLoading() {
    return this._publicLoading
  }
  get adminProfile() {
    return this._adminProfile
  }
  get adminLoading() {
    return this._adminLoading
  }
  get adminSaving() {
    return this._adminSaving
  }
  get adminLoaded() {
    return this._adminLoaded
  }
  get adminError() {
    return this._adminError
  }
  get adminSuccess() {
    return this._adminSuccess
  }
  get adminDirty() {
    return this._adminDirty
  }

  async loadPublicProfile() {
    if (this._publicProfile || this._publicLoading) return
    this._publicLoading = true
    this.notifyChange()
    try {
      await this._internalFetch()
    } catch {
      // Fallback handled by components
    } finally {
      this._publicLoading = false
      this.notifyChange()
    }
  }

  async loadAdminProfile() {
    if (this._adminLoaded || this._adminLoading) return
    this._adminLoading = true
    this._adminError = ''
    this.notifyChange()
    try {
      await this._internalFetch()
      this._adminLoaded = true
      this._adminDirty = false
    } catch (error) {
      this._adminError = describeApiError(error)
    } finally {
      this._adminLoading = false
      this.notifyChange()
    }
  }

  private async _internalFetch() {
    if (this._fetchPromise) return this._fetchPromise
    this._fetchPromise = getMe()
    try {
      const p = await this._fetchPromise
      this._publicProfile = p
      this._adminProfile = p
      return p
    } finally {
      this._fetchPromise = null
    }
  }

  async saveAdminProfile(profile: MeProfile) {
    this._adminSaving = true
    this._adminError = ''
    this._adminSuccess = ''
    this.notifyChange()
    try {
      this._adminProfile = await updateMe(profile)
      this._publicProfile = this._adminProfile
      this._adminLoaded = true
      this._adminDirty = false
      this._adminSuccess = 'プロフィールを更新しました。'
    } catch (error) {
      this._adminError = describeApiError(error)
    } finally {
      this._adminSaving = false
      this.notifyChange()
    }
  }

  setAdminDirty(dirty: boolean) {
    this._adminDirty = dirty
    if (dirty) {
      this._adminSuccess = ''
    }
    this.notifyChange()
  }
}
