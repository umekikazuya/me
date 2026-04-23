import type { ReactiveController, ReactiveControllerHost } from 'lit'
import type { ArticleController as IArticleController } from '../contexts/article-context.js'

export class ArticleController
  implements ReactiveController, IArticleController
{
  private hosts: Set<ReactiveControllerHost> = new Set()
  private _adminDirty = false

  constructor(host: ReactiveControllerHost) {
    this.addHost(host)
  }

  addHost(host: ReactiveControllerHost) {
    this.hosts.add(host)
    host.addController(this)
  }

  hostConnected() {}

  private requestUpdate() {
    for (const host of this.hosts) {
      host.requestUpdate()
    }
  }

  get adminDirty() {
    return this._adminDirty
  }

  setAdminDirty(dirty: boolean) {
    this._adminDirty = dirty
    this.requestUpdate()
  }
}
