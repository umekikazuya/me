import { consume } from '@lit/context'
import { SignalWatcher } from '@lit-labs/signals'
import { html, LitElement } from 'lit'
import { customElement } from 'lit/decorators.js'
import { authContext } from '../contexts/auth-context.js'
import type { IAuthRepository } from '../domain/AuthRepository.js'
import '../components/admin/ui/me-text-input.js'

@customElement('page-admin-entry')
export class PageAdminEntry extends SignalWatcher(LitElement) {
  @consume({ context: authContext, subscribe: true })
  set authRepo(repo: IAuthRepository) {
    if (this._authRepo === repo) return
    this._authRepo = repo
  }
  private _authRepo!: IAuthRepository

  connectedCallback(): void {
    super.connectedCallback()
    void this.bootstrap()
  }

  private async bootstrap() {
    if (this.authRepo?.status.value === 'unknown') {
      await this.authRepo.refreshSession()
    }
  }

  render() {
    const status = this.authRepo?.status.value
    if (status === 'unknown' || 'checking') {
      return html`<p>認証状態を確認しています...</p>`
    }
    if (status === 'authenticated') {
      return html`<page-admin-dashboard></page-admin-dashboard>`
    }
    return html`<page-admin-login></page-admin-login>`
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'page-admin-entry': PageAdminEntry
  }
}
