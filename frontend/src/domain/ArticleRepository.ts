import { computed, type ReadonlySignal, signal } from '@preact/signals-core'
import {
  createArticle,
  deleteArticle,
  listArticles,
  listArticleTags,
  updateArticle,
} from '../admin/article-api.js'
import type {
  ArticleDraft,
  ArticleItem,
  ArticlePlatform,
  ArticleTagItem,
} from '../admin/article-types.js'
import { describeApiError } from '../admin/types.js'
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
  readonly state: ReadonlySignal<IState<AdminArticleData>>

  readonly articles: ReadonlySignal<ArticleItem[]>
  readonly tagOptions: ReadonlySignal<ArticleTagItem[]>
  readonly adminDirty: ReadonlySignal<boolean>
  readonly isLoading: ReadonlySignal<boolean>
  readonly error: ReadonlySignal<string>

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
  private _state = signal<IState<AdminArticleData>>(
    createInitialState(DEFAULT_ARTICLE_DATA),
  )

  get state() {
    return this._state
  }
  public articles = computed(() => this._state.value.data?.articles ?? [])
  public tagOptions = computed(() => this._state.value.data?.tagOptions ?? [])
  public adminDirty = computed(
    () => this._state.value.data?.adminDirty ?? false,
  )
  public isLoading = computed(() => this._state.value.status === 'loading')
  public error = computed(() => this._state.value.error?.message ?? '')

  async loadInitialData() {
    const gen = this.nextGeneration()
    this.updateState(this._state, { status: 'loading', error: null })

    try {
      const [articlesResult, tagsResult] = await Promise.all([
        listArticles({ limit: 50 }),
        listArticleTags(),
      ])

      if (!this.isCurrent(gen)) return

      this.updateState(this._state, {
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
      this.updateState(this._state, {
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
      this.updateState(this._state, { status: 'loading', error: null })
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
        ? [...this.articles.value, ...result.articles]
        : result.articles

      this.updateState(this._state, {
        status: 'success',
        data: {
          ...this.ensureData(),
          articles: nextArticles,
          nextCursor: result.nextCursor,
        },
      })
    } catch (error) {
      if (!this.isCurrent(gen)) return
      this.updateState(this._state, {
        status: 'error',
        error: { code: 'RELOAD_FAILED', message: describeApiError(error) },
      })
    }
  }

  async createArticle(draft: ArticleDraft) {
    this.updateState(this._state, { error: null })
    try {
      await createArticle(draft)
      this.patchData({ adminDirty: false }, 'success')
      // No full result from API, rely on reload in the interactor logic or full refresh
      await this.loadInitialData()
    } catch (error) {
      this.updateState(this._state, {
        error: { code: 'CREATE_FAILED', message: describeApiError(error) },
      })
    }
  }

  async updateArticle(externalId: string, draft: ArticleDraft) {
    this.updateState(this._state, { error: null })
    try {
      await updateArticle(externalId, draft)
      this.patchData({ adminDirty: false }, 'success')
      await this.loadInitialData()
    } catch (error) {
      this.updateState(this._state, {
        error: { code: 'UPDATE_FAILED', message: describeApiError(error) },
      })
    }
  }

  async deleteArticle(externalId: string) {
    this.updateState(this._state, { error: null })
    try {
      await deleteArticle(externalId)
      await this.loadInitialData()
    } catch (error) {
      this.updateState(this._state, {
        error: { code: 'DELETE_FAILED', message: describeApiError(error) },
      })
    }
  }

  setAdminDirty(dirty: boolean) {
    this.patchData({ adminDirty: dirty })
  }

  private ensureData(): AdminArticleData {
    return this._state.value.data ?? DEFAULT_ARTICLE_DATA
  }

  private patchData(patch: Partial<AdminArticleData>, status?: StateStatus) {
    this.updateState(this._state, {
      status: status ?? this._state.value.status,
      data: { ...this.ensureData(), ...patch },
    })
  }
}
