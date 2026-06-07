import { createContext } from '@lit/context'
import type { IProfileRepository } from '../domain/ProfileRepository.js'

export const profileContext =
  createContext<IProfileRepository>('profile-context')
