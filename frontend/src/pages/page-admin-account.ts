import { consume } from '@lit/context'
import { SignalWatcher } from '@lit-labs/signals'
import { css, html, LitElement } from 'lit'
import { customElement } from 'lit/decorators.js'
import { adminFormStyles } from '../admin/admin-form-styles.js'
import { authContext } from '../contexts/auth-context.js'
import type { IAuthRepository } from '../domain/AuthRepository.js'
import '../components/admin/ui/me-text-input.js'

@customElement('page-admin-account')
export class PageAdminAccount extends SignalWatcher(LitElement) {
  @consume({ context: authContext, subscribe: true })
  set authRepo(repo: IAuthRepository) {
    this._authRepo = repo
  }
  get authRepo() {
    return this._authRepo
  }
  private _authRepo!: IAuthRepository

  render() {
    const a = this.authRepo
    const error = a.error.value
    const success = a.success.value
    const busyAction = a.accountBusyAction.value

    return html`
      <section class="container">
        <header>
          <p class="eyebrow" lang="en">Account</p>
          <h1 class="title">アカウント管理</h1>
          <p class="description">
            セッション操作とメールアドレス変更を管理します。
          </p>
        </header>

        ${error ? html`<p class="message error">${error}</p>` : null}
        ${success ? html`<p class="message success">${success}</p>` : null}

        <section class="card">
          <div>
            <h2>ログアウト</h2>
            <p>現在の端末のセッションを終了します。</p>
            <p class="note">このブラウザでの編集作業を終了するときに使います。</p>
          </div>
          <button
            type="button"
            ?disabled=${busyAction !== ''}
            @click=${this.handleLogout}
          >
            ${busyAction === 'logout' ? '実行中...' : 'ログアウトする'}
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
            ?disabled=${busyAction !== ''}
            @click=${this.handleRevokeAllSessions}
          >
            ${
              busyAction === 'revoke-sessions'
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
              name="token"
              ?disabled=${busyAction !== ''}
              required
            ></me-text-input>

            <me-text-input
              label="New email address"
              name="newEmailAddress"
              type="email"
              ?disabled=${busyAction !== ''}
              required
            ></me-text-input>

            <button type="submit" ?disabled=${busyAction !== ''}>
              ${busyAction === 'change-email' ? '送信中...' : 'メール変更を送信'}
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
    const form = event.target as HTMLFormElement
    const formData = new FormData(form)

    await this.authRepo.changeEmail({
      token: (formData.get('token') as string).trim(),
      newEmailAddress: (formData.get('newEmailAddress') as string).trim(),
    })

    if (!this.authRepo.error.value) {
      form.reset()
    }
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
