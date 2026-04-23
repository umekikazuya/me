import type { ReactiveController, ReactiveControllerHost } from 'lit'

/**
 * A generic Lit controller that observes any EventTarget-based repository.
 * When the repository dispatches a 'change' event, the host component is updated.
 */
export class RepositoryObserver implements ReactiveController {
  private _host: ReactiveControllerHost
  private _repository: EventTarget

  constructor(host: ReactiveControllerHost, repository: EventTarget) {
    this._host = host
    this._repository = repository
    host.addController(this)
  }

  hostConnected() {
    this._repository.addEventListener('change', this._onRepositoryChange)
  }

  hostDisconnected() {
    this._repository.removeEventListener('change', this._onRepositoryChange)
  }

  private _onRepositoryChange = () => {
    this._host.requestUpdate()
  }
}
