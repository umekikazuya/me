import {
  createArticle,
  deleteArticle,
  listArticles,
  listArticleTags,
  updateArticle,
} from '../api/article-api.js'
import type {
  ArticleDraft,
  ArticleItem,
  ArticlePlatform,
  ArticleTagItem,
} from '../api/article-types.js'
import { describeApiError } from '../api/types.js'
import {
  createInitialState,
  type IState,
  Repository,
  type StateStatus,
} from './Repository.js'

export interface ArticleEventMap {
  'article:created': CustomEvent<ArticleItem>
  'article:updated': CustomEvent<ArticleItem>
  'article:deleted': CustomEvent<{ externalId: string }>
}

export interface AdminArticleData {
  articles: ArticleItem[]
  tagOptions: ArticleTagItem[]
  nextCursor?: string
  adminDirty: boolean
}

/**
 * The public interface for ArticleRepository.
 */
export interface IArticleRepository extends EventTarget {
  readonly articles: ArticleItem[]
  readonly tagOptions: ArticleTagItem[]
  readonly adminDirty: boolean
  readonly isLoading: boolean
  readonly error: string

  addEventListener<K extends keyof ArticleEventMap>(
    type: K,
    listener: (e: ArticleEventMap[K]) => void,
    options?: boolean | AddEventListenerOptions,
  ): void
  addEventListener(
    type: string,
    callback: EventListenerOrEventListenerObject | null,
    options?: boolean | AddEventListenerOptions,
  ): void

  loadInitialData(): Promise<void>
  reloadArticles(params?: {
    q?: string
    year?: number
    platform?: string
    tag?: string[]
    cursor?: string
    append?: boolean
  }): Promise<void>
  createArticle(draft: ArticleDraft): Promise<void>
  updateArticle(externalId: string, draft: ArticleDraft): Promise<void>
  deleteArticle(externalId: string): Promise<void>
  setAdminDirty(dirty: boolean): void
}

const DEFAULT_ARTICLE_DATA: AdminArticleData = {
  articles: [],
  tagOptions: [],
  nextCursor: undefined,
  adminDirty: false,
}

export class ArticleRepository
  extends Repository
  implements IArticleRepository
{
  private _state: IState<AdminArticleData> = createInitialState(DEFAULT_ARTICLE_DATA)

  get articles(): ArticleItem[] {
    return this._state.data?.articles ?? []
  }

  get tagOptions(): ArticleTagItem[] {
    return this._state.data?.tagOptions ?? []
  }

  get adminDirty(): boolean {
    return this._state.data?.adminDirty ?? false
  }

  get isLoading(): boolean {
    return this._state.status === 'loading'
  }

  get error(): string {
    return this._state.error?.message ?? ''
  }

  private setState(patch: Partial<IState<AdminArticleData>>): void {
    this._state = { ...this._state, ...patch }
    this.emitChange()
  }

  async loadInitialData() {
    const gen = this.nextGeneration()
    this.setState({ status: 'loading', error: null })

    try {
      const [articlesResult, tagsResult] = await Promise.all([
        listArticles({ limit: 50 }),
        listArticleTags(),
      ])

      if (!this.isCurrent(gen)) return

      this.setState({
        status: 'success',
        data: {
          ...this.ensureData(),
          articles: articlesResult.articles,
          nextCursor: articlesResult.nextCursor,
          tagOptions: tagsResult,
        },
      })
    } catch (error) {
      if (!this.isCurrent(gen)) return
      this.setState({
        status: 'error',
        error: { code: 'LOAD_FAILED', message: describeApiError(error) },
      })
    }
  }

  async reloadArticles(
    params: {
      q?: string
      year?: number
      platform?: string
      tag?: string[]
      cursor?: string
      append?: boolean
    } = {},
  ) {
    const gen = this.nextGeneration()
    const isAppend = params.append ?? false
    if (!isAppend) {
      this.setState({ status: 'loading', error: null })
    }

    try {
      const result = await listArticles({
        q: params.q,
        year: params.year,
        platform: params.platform as ArticlePlatform,
        tag: params.tag,
        limit: 50,
        cursor: params.cursor,
      })

      if (!this.isCurrent(gen)) return

      const nextArticles = isAppend
        ? [...this.articles, ...result.articles]
        : result.articles

      this.setState({
        status: 'success',
        data: {
          ...this.ensureData(),
          articles: nextArticles,
          nextCursor: result.nextCursor,
        },
      })
    } catch (error) {
      if (!this.isCurrent(gen)) return
      this.setState({
        status: 'error',
        error: { code: 'RELOAD_FAILED', message: describeApiError(error) },
      })
    }
  }

  async createArticle(draft: ArticleDraft) {
    this.setState({ error: null })
    try {
      await createArticle(draft)
      this.patchData({ adminDirty: false }, 'success')
      await this.loadInitialData()
    } catch (error) {
      this.setState({
        error: { code: 'CREATE_FAILED', message: describeApiError(error) },
      })
    }
  }

  async updateArticle(externalId: string, draft: ArticleDraft) {
    this.setState({ error: null })
    try {
      await updateArticle(externalId, draft)
      this.patchData({ adminDirty: false }, 'success')
      await this.loadInitialData()
    } catch (error) {
      this.setState({
        error: { code: 'UPDATE_FAILED', message: describeApiError(error) },
      })
    }
  }

  async deleteArticle(externalId: string) {
    this.setState({ error: null })
    try {
      await deleteArticle(externalId)
      await this.loadInitialData()
    } catch (error) {
      this.setState({
        error: { code: 'DELETE_FAILED', message: describeApiError(error) },
      })
    }
  }

  setAdminDirty(dirty: boolean) {
    this.patchData({ adminDirty: dirty })
  }

  private ensureData(): AdminArticleData {
    return this._state.data ?? DEFAULT_ARTICLE_DATA
  }

  private patchData(patch: Partial<AdminArticleData>, status?: StateStatus) {
    this.setState({
      status: status ?? this._state.status,
      data: { ...this.ensureData(), ...patch },
    })
  }
}
