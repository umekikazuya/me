import type { ReactiveController, ReactiveControllerHost } from 'lit'

/**
 * A standard adapter that observes a domain repository.
 * Bridges the Domain's EventTarget to the UI's reactive lifecycle.
 * Manages memory safety internally via AbortController.
 */
export class RepositoryObserver implements ReactiveController {
  private _host: ReactiveControllerHost
  private _repository: EventTarget
  private _onChange?: () => void
  private _abortController?: AbortController

  constructor(
    host: ReactiveControllerHost,
    repository: EventTarget,
    onChange?: () => void,
  ) {
    this._host = host
    this._repository = repository
    this._onChange = onChange
    host.addController(this)
  }

  hostConnected() {
    this.connect()
  }

  hostDisconnected() {
    this.disconnect()
  }

  /**
   * Starts observing the repository with a fresh AbortController.
   */
  connect() {
    this.disconnect() // Ensure previous is cleared
    this._abortController = new AbortController()

    this._repository.addEventListener('change', this._onRepositoryChange, {
      signal: this._abortController.signal,
    })
  }

  /**
   * Stops all observations.
   */
  disconnect() {
    this._abortController?.abort()
    this._abortController = undefined
  }

  private _onRepositoryChange = () => {
    this._host.requestUpdate()
    if (this._onChange) {
      this._onChange()
    }
  }
}
