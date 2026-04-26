import { createContext } from '@lit/context'
import type { IAuthRepository } from '../domain/AuthRepository.js'

export const authContext = createContext<IAuthRepository>('auth-context')
