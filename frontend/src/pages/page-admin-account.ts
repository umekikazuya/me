import { consume } from '@lit/context'
import { css, html, LitElement } from 'lit'
import { customElement, state } from 'lit/decorators.js'
import { adminFormStyles } from '../admin/admin-form-styles.js'
import { authContext } from '../contexts/auth-context.js'
import { RepositoryObserver } from '../controllers/RepositoryObserver.js'
import type { IAuthRepository } from '../domain/AuthRepository.js'
import '../components/admin/ui/me-text-input.js'

@customElement('page-admin-account')
export class PageAdminAccount extends LitElement {
  @consume({ context: authContext, subscribe: true })
  authRepo!: IAuthRepository

  @state()
  private token = ''

  @state()
  private newEmailAddress = ''

  @state()
  private lastSubmittedAction = ''

  constructor() {
    super()
    new RepositoryObserver(this, this.authRepo)
  }

  protected updated(changedProperties: Map<PropertyKey, unknown>) {
    if (
      changedProperties.has('authRepo') &&
      this.authRepo.accountSuccess &&
      this.lastSubmittedAction === 'change-email'
    ) {
      this.token = ''
      this.newEmailAddress = ''
      this.lastSubmittedAction = ''
    }
  }

  render() {
    const a = this.authRepo
    return html`
      <section class="container">
        <header>
          <p class="eyebrow" lang="en">Account</p>
          <h1 class="title">アカウント管理</h1>
          <p class="description">
            セッション操作とメールアドレス変更を管理します。
          </p>
        </header>

        ${a.accountError ? html`<p class="message error">${a.accountError}</p>` : null}
        ${
          a.accountSuccess
            ? html`<p class="message success">${a.accountSuccess}</p>`
            : null
        }

        <section class="card">
          <div>
            <h2>ログアウト</h2>
            <p>現在の端末のセッションを終了します。</p>
            <p class="note">このブラウザでの編集作業を終了するときに使います。</p>
          </div>
          <button
            type="button"
            ?disabled=${a.accountBusyAction !== ''}
            @click=${this.handleLogout}
          >
            ${a.accountBusyAction === 'logout' ? '実行中...' : 'ログアウトする'}
          </button>
        </section>

        <section class="card">
          <div>
            <h2>全セッション失効</h2>
            <p>他の端末を含む全セッションを失効させます。</p>
            <p class="note">
              共有端末や漏洩が心配な場合に使います。現在の端末も再ログインが必要になります。
            </p>
          </div>
          <button
            type="button"
            class="danger"
            ?disabled=${a.accountBusyAction !== ''}
            @click=${this.handleRevokeAllSessions}
          >
            ${
              a.accountBusyAction === 'revoke-sessions'
                ? '実行中...'
                : 'すべてのセッションを終了'
            }
          </button>
        </section>

        <section class="card card-form">
          <div>
            <h2>メールアドレス変更</h2>
            <p>
              API 仕様に合わせて、変更トークンと新しいメールアドレスを送信します。
            </p>
            <p class="note">
              トークンが未発行なら、バックエンド側のメール変更フロー準備後に利用してください。
            </p>
          </div>

          <form @submit=${this.handleChangeEmail}>
            <me-text-input
              label="Token"
              .value=${this.token}
              ?disabled=${a.accountBusyAction !== ''}
              required
              @change=${(e: CustomEvent) => (this.token = e.detail)}
            ></me-text-input>

            <me-text-input
              label="New email address"
              type="email"
              .value=${this.newEmailAddress}
              ?disabled=${a.accountBusyAction !== ''}
              required
              @change=${(e: CustomEvent) => (this.newEmailAddress = e.detail)}
            ></me-text-input>

            <button type="submit" ?disabled=${a.accountBusyAction !== ''}>
              ${a.accountBusyAction === 'change-email' ? '送信中...' : 'メール変更を送信'}
            </button>
          </form>
        </section>
      </section>
    `
  }

  private handleLogout = async () => {
    if (!window.confirm('現在の端末からログアウトします。よろしいですか？'))
      return

    await this.authRepo.logout()
  }

  private handleRevokeAllSessions = async () => {
    if (
      !window.confirm(
        'すべてのセッションを終了します。現在の端末も再ログインが必要になります。実行しますか？',
      )
    ) {
      return
    }

    await this.authRepo.revokeAllSessions()
  }

  private async handleChangeEmail(event: Event) {
    event.preventDefault()
    this.lastSubmittedAction = 'change-email'

    await this.authRepo.changeEmail({
      token: this.token.trim(),
      newEmailAddress: this.newEmailAddress.trim(),
    })
  }

  static styles = [
    adminFormStyles,
    css`
      :host {
        display: block;
      }

      .container {
        display: grid;
        gap: 24px;
      }

      p {
        color: var(--color-text-secondary);
        line-height: 1.8;
        font-size: 14px;
      }

      .note {
        font-size: 13px;
        color: var(--color-text-tertiary);
        margin-top: 8px;
      }

      .card {
        display: grid;
        gap: 20px;
        padding: 24px;
        border: 1px solid var(--color-border);
        background: #fff;
      }

      .card-form {
        gap: 24px;
      }

      h2 {
        font-size: 18px;
        font-weight: 500;
        margin-bottom: 8px;
        color: var(--color-text-primary);
      }

      form {
        display: grid;
        gap: 16px;
      }

      button {
        justify-self: start;
      }
    `,
  ]
}

declare global {
  interface HTMLElementTagNameMap {
    'page-admin-account': PageAdminAccount
  }
}
