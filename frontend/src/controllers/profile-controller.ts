import type { ReactiveController, ReactiveControllerHost } from 'lit'
import { getMe, updateMe } from '../admin/me-api.js'
import {
  createEmptyMeProfile,
  describeApiError,
  type MeProfile,
} from '../admin/types.js'
import type { ProfileController as IProfileController } from '../contexts/profile-context.js'

export class ProfileController
  implements ReactiveController, IProfileController
{
  private host: ReactiveControllerHost

  private _publicProfile: MeProfile | null = null
  private _publicLoading = false
  private _adminProfile = createEmptyMeProfile()
  private _adminLoading = false
  private _adminSaving = false
  private _adminLoaded = false
  private _adminError = ''
  private _adminSuccess = ''
  private _adminDirty = false

  constructor(host: ReactiveControllerHost) {
    this.host = host
    host.addController(this)
  }

  hostConnected() {}

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
    this._publicLoading = true
    this.host.requestUpdate()
    try {
      this._publicProfile = await getMe()
    } catch {
      // Fallback handled by components
    } finally {
      this._publicLoading = false
      this.host.requestUpdate()
    }
  }

  async loadAdminProfile() {
    this._adminLoading = true
    this._adminError = ''
    this.host.requestUpdate()
    try {
      this._adminProfile = await getMe()
      this._adminLoaded = true
      this._adminDirty = false
    } catch (error) {
      this._adminError = describeApiError(error)
    } finally {
      this._adminLoading = false
      this.host.requestUpdate()
    }
  }

  async saveAdminProfile(profile: MeProfile) {
    this._adminSaving = true
    this._adminError = ''
    this._adminSuccess = ''
    this.host.requestUpdate()
    try {
      this._adminProfile = await updateMe(profile)
      this._adminLoaded = true
      this._adminDirty = false
      this._adminSuccess = 'プロフィールを更新しました。'
    } catch (error) {
      this._adminError = describeApiError(error)
    } finally {
      this._adminSaving = false
      this.host.requestUpdate()
    }
  }

  setAdminDirty(dirty: boolean) {
    this._adminDirty = dirty
    if (dirty) {
      this._adminSuccess = ''
    }
    this.host.requestUpdate()
  }
}
