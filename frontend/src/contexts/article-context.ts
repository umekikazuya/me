import { createContext } from '@lit/context'
import type { ReactiveControllerHost } from 'lit'

export interface ArticleController {
  readonly adminDirty: boolean
  addHost(host: ReactiveControllerHost): void
  setAdminDirty(dirty: boolean): void
}

export const articleContext =
  createContext<ArticleController>('article-context')
