import { createContext } from '@lit/context'
import type { IArticleRepository } from '../domain/ArticleRepository.js'

export const articleContext =
  createContext<IArticleRepository>('article-context')
