import { createContext } from '@lit/context'
import type { ReactiveControllerHost } from 'lit'
import type { MeProfile } from '../admin/types.js'

export interface ProfileController {
  readonly publicProfile: MeProfile | null
  readonly publicLoading: boolean
  readonly adminProfile: MeProfile
  readonly adminLoading: boolean
  readonly adminSaving: boolean
  readonly adminLoaded: boolean
  readonly adminError: string
  readonly adminSuccess: string
  readonly adminDirty: boolean

  addHost(host: ReactiveControllerHost): void
  loadPublicProfile(): Promise<void>
  loadAdminProfile(): Promise<void>
  saveAdminProfile(profile: MeProfile): Promise<void>
  setAdminDirty(dirty: boolean): void
}

export const profileContext =
  createContext<ProfileController>('profile-context')
