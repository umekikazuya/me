import { consume } from '@lit/context'
import { css, html, LitElement, nothing } from 'lit'
import { customElement } from 'lit/decorators.js'
import { authContext } from '../../../contexts/auth-context.js'
import { RepositoryObserver } from '../../../controllers/RepositoryObserver.js'
import type { IAuthRepository } from '../../../domain/AuthRepository.js'

/**
 * A declarative component that protects its content based on authentication status.
 */
@customElement('me-auth-guard')
export class MeAuthGuard extends LitElement {
  @consume({ context: authContext, subscribe: true })
  set authRepo(repo: IAuthRepository) {
    if (this._authRepo === repo) return
    this._authRepo = repo
    this._observer?.disconnect()
    if (repo) this._observer = new RepositoryObserver(this, repo)
  }
  get authRepo() {
    return this._authRepo
  }
  private _authRepo!: IAuthRepository
  private _observer?: RepositoryObserver

  connectedCallback() {
    super.connectedCallback()
    this.checkSession()
  }

  private async checkSession() {
    if (this.authRepo?.status === 'unknown') {
      await this.authRepo.refreshSession()
    }
  }

  render() {
    const status = this.authRepo?.status

    if (status === 'checking' || status === 'unknown') {
      return html`
        <div class="guard-status">
          <p>認証状態を確認しています...</p>
        </div>
      `
    }

    if (status === 'authenticated') {
      return html`<slot></slot>`
    }

    // Guest or failed - render nothing, parent orchestrator will handle redirect
    return nothing
  }

  static styles = css`
    :host {
      display: contents;
    }

    .guard-status {
      min-height: 60dvh;
      display: grid;
      place-items: center;
      color: var(--color-text-secondary);
      font-size: 15px;
      letter-spacing: var(--tracking-wide);
    }
  `
}

declare global {
  interface HTMLElementTagNameMap {
    'me-auth-guard': MeAuthGuard
  }
}
