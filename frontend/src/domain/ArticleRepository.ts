import { Repository } from './Repository.js'

export interface IArticleRepository extends Repository {
  readonly adminDirty: boolean
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
    this.notifyChange()
  }
}
