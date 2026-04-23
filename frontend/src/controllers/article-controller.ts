import type { ReactiveController, ReactiveControllerHost } from 'lit'
import type { ArticleController as IArticleController } from '../contexts/article-context.js'

export class ArticleController
  implements ReactiveController, IArticleController
{
  private host: ReactiveControllerHost
  private _adminDirty = false

  constructor(host: ReactiveControllerHost) {
    this.host = host
    host.addController(this)
  }

  hostConnected() {}

  get adminDirty() {
    return this._adminDirty
  }

  setAdminDirty(dirty: boolean) {
    this._adminDirty = dirty
    this.host.requestUpdate()
  }
}
