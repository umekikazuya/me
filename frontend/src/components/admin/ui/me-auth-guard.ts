import { consume } from '@lit/context'
import { SignalWatcher } from '@lit-labs/signals'
import { css, html, LitElement, nothing } from 'lit'
import { customElement } from 'lit/decorators.js'
import { authContext } from '../../../contexts/auth-context.js'
import type { IAuthRepository } from '../../../domain/AuthRepository.js'

/**
 * A declarative component that protects its content based on authentication status.
 * Uses SignalWatcher for fine-grained reactivity.
 */
@customElement('me-auth-guard')
export class MeAuthGuard extends SignalWatcher(LitElement) {
  @consume({ context: authContext, subscribe: true })
  set authRepo(repo: IAuthRepository) {
    this._authRepo = repo
  }
  get authRepo() {
    return this._authRepo
  }
  private _authRepo!: IAuthRepository

  connectedCallback() {
    super.connectedCallback()
    this.checkSession()
  }

  private async checkSession() {
    if (this.authRepo?.status.value === 'unknown') {
      await this.authRepo.refreshSession()
    }
  }

  render() {
    const status = this.authRepo?.status.value

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

    *, *::before, *::after {
      box-sizing: border-box;
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
