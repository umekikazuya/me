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
  private hosts: Set<ReactiveControllerHost> = new Set()

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

  constructor(host: ReactiveControllerHost) {
    this.addHost(host)
  }

  addHost(host: ReactiveControllerHost) {
    this.hosts.add(host)
    host.addController(this)
  }

  hostConnected() {}
  hostDisconnected() {
    // We don't want to remove from Set here necessarily if it's a shared instance
    // but for safety in short-lived components we could.
    // However, for this app's lifecycle, keeping them is fine or we can manage it.
  }

  private requestUpdate() {
    for (const host of this.hosts) {
      host.requestUpdate()
    }
  }

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
    this.requestUpdate()
    try {
      await this._internalFetch()
    } catch {
      // Fallback handled by components
    } finally {
      this._publicLoading = false
      this.requestUpdate()
    }
  }

  async loadAdminProfile() {
    if (this._adminLoaded || this._adminLoading) return
    this._adminLoading = true
    this._adminError = ''
    this.requestUpdate()
    try {
      await this._internalFetch()
      this._adminLoaded = true
      this._adminDirty = false
    } catch (error) {
      this._adminError = describeApiError(error)
    } finally {
      this._adminLoading = false
      this.requestUpdate()
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
    this.requestUpdate()
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
      this.requestUpdate()
    }
  }

  setAdminDirty(dirty: boolean) {
    this._adminDirty = dirty
    if (dirty) {
      this._adminSuccess = ''
    }
    this.requestUpdate()
  }
}
