import { consume } from '@lit/context'
import { SignalWatcher } from '@lit-labs/signals'
import { css, html, LitElement } from 'lit'
import { customElement } from 'lit/decorators.js'
import { authContext } from '../../../contexts/auth-context.js'
import { RepositoryObserver } from '../../../controllers/RepositoryObserver.js'
import type { IAuthRepository } from '../../../domain/AuthRepository.js'
import '../../../pages/page-admin-login.js'

@customElement('me-admin-auth-boundary')
export class MeAdminAuthBoundary extends SignalWatcher(LitElement) {
  @consume({ context: authContext, subscribe: true })
  set authRepo(repo: IAuthRepository) {
    if (this._authRepo === repo) return
    this._authRepo = repo
    this._observer?.disconnect()
    if (repo) {
      this._observer = new RepositoryObserver(this, repo)
      void this.bootstrap()
    }
  }
  get authRepo() {
    return this._authRepo
  }
  private _authRepo!: IAuthRepository
  private _observer?: RepositoryObserver

  private async bootstrap() {
    if (this.authRepo?.status.value === 'unknown') {
      await this.authRepo.refreshSession()
    }
  }

  render() {
    const status = this.authRepo?.status.value

    if (status === 'unknown' || status === 'checking') {
      return html`
        <div class="status">
          <p>認証状態を確認しています...</p>
        </div>
      `
    }

    if (status === 'authenticated') {
      return html`<slot></slot>`
    }

    return html`<page-admin-login></page-admin-login>`
  }

  static styles = css`
    :host {
      display: contents;
    }

    *, *::before, *::after {
      box-sizing: border-box;
    }

    .status {
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
    'me-admin-auth-boundary': MeAdminAuthBoundary
  }
}
