import { createContext } from '@lit/context'

export interface ArticleController {
  readonly adminDirty: boolean
  setAdminDirty(dirty: boolean): void
}

export const articleContext = createContext<ArticleController>('article-context')
