import { Repository } from './Repository.js'

export interface ArticleEventMap {
  'article:admin-dirty-change': CustomEvent<{ dirty: boolean }>
}

/**
 * The public interface for ArticleRepository.
 */
export interface IArticleRepository extends EventTarget {
  readonly adminDirty: boolean

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

  setAdminDirty(dirty: boolean): void
}

export class ArticleRepository
  extends Repository
  implements IArticleRepository
{
  private _adminDirty = false

  get adminDirty() {
    return this._adminDirty
  }

  setAdminDirty(dirty: boolean) {
    this._adminDirty = dirty
    this.dispatchEvent(
      new CustomEvent('article:admin-dirty-change', {
        detail: { dirty },
      }),
    )
    this.notifyChange()
  }
}
