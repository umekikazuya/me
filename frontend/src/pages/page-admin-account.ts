import { css, html, LitElement } from 'lit'
import { customElement, property, state } from 'lit/decorators.js'
import type { ChangeEmailInput } from '../admin/types.js'

@customElement('page-admin-account')
export class PageAdminAccount extends LitElement {
  @property()
  busyAction = ''

  @property()
  errorMessage = ''

  @property()
  successMessage = ''

  @state()
  private token = ''

  @state()
  private newEmailAddress = ''

  render() {
    return html`
      <section class="container">
        <header>
          <p class="eyebrow" lang="en">Account</p>
          <h1 class="title">アカウント管理</h1>
          <p class="description">
            セッション操作とメールアドレス変更を管理します。
          </p>
        </header>

        ${this.errorMessage ? html`<p class="message error">${this.errorMessage}</p>` : null}
        ${
          this.successMessage
            ? html`<p class="message success">${this.successMessage}</p>`
            : null
        }

        <section class="card">
          <div>
            <h2>ログアウト</h2>
            <p>現在の端末のセッションを終了します。</p>
          </div>
          <button
            type="button"
            ?disabled=${this.busyAction !== ''}
            @click=${this.handleLogout}
          >
            ${this.busyAction === 'logout' ? '実行中...' : 'ログアウトする'}
          </button>
        </section>

        <section class="card">
          <div>
            <h2>全セッション失効</h2>
            <p>他の端末を含む全セッションを失効させます。</p>
          </div>
          <button
            type="button"
            class="danger"
            ?disabled=${this.busyAction !== ''}
            @click=${this.handleRevokeAllSessions}
          >
            ${
              this.busyAction === 'revoke-sessions'
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
          </div>

          <form @submit=${this.handleChangeEmail}>
            <label class="field">
              <span>Token</span>
              <input
                .value=${this.token}
                @input=${(event: Event) => {
                  this.token = (event.target as HTMLInputElement).value
                }}
                required
              />
            </label>
            <label class="field">
              <span>New email address</span>
              <input
                type="email"
                .value=${this.newEmailAddress}
                @input=${(event: Event) => {
                  this.newEmailAddress = (
                    event.target as HTMLInputElement
                  ).value
                }}
                required
              />
            </label>
            <button type="submit" ?disabled=${this.busyAction !== ''}>
              ${this.busyAction === 'change-email' ? '送信中...' : 'メール変更を送信'}
            </button>
          </form>
        </section>
      </section>
    `
  }

  private handleLogout = () => {
    this.dispatchEvent(
      new CustomEvent('admin-logout', {
        bubbles: true,
        composed: true,
      }),
    )
  }

  private handleRevokeAllSessions = () => {
    this.dispatchEvent(
      new CustomEvent('admin-revoke-sessions', {
        bubbles: true,
        composed: true,
      }),
    )
  }

  private handleChangeEmail(event: Event) {
    event.preventDefault()

    const detail: ChangeEmailInput = {
      token: this.token.trim(),
      newEmailAddress: this.newEmailAddress.trim(),
    }

    this.dispatchEvent(
      new CustomEvent<ChangeEmailInput>('admin-change-email', {
        detail,
        bubbles: true,
        composed: true,
      }),
    )
  }

  static styles = css`
    :host {
      display: block;
    }

    .container {
      display: grid;
      gap: 24px;
    }

    .eyebrow {
      font-family: var(--font-en);
      letter-spacing: var(--tracking-wider);
      color: var(--color-text-tertiary);
      margin-bottom: 12px;
    }

    .title {
      font-size: 30px;
      font-weight: 300;
      margin-bottom: 12px;
    }

    .description,
    p {
      color: var(--color-text-secondary);
      line-height: 1.8;
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
      font-size: 22px;
      font-weight: 300;
      margin-bottom: 10px;
    }

    form {
      display: grid;
      gap: 16px;
    }

    .field {
      display: grid;
      gap: 8px;
      color: var(--color-text-secondary);
      font-size: 14px;
    }

    input {
      border: 1px solid var(--color-border);
      background: #fff;
      padding: 12px 14px;
      font: inherit;
      color: var(--color-text-primary);
    }

    input:focus {
      outline: none;
      border-color: var(--color-text-primary);
    }

    button {
      justify-self: start;
      border: 0;
      background: var(--color-text-primary);
      color: #fff;
      padding: 12px 18px;
      font: inherit;
      cursor: pointer;
    }

    button:disabled {
      opacity: 0.5;
      cursor: wait;
    }

    .danger {
      background: #8c4a4a;
    }

    .message {
      font-size: 14px;
      line-height: 1.7;
    }

    .error {
      color: #9a3f3f;
    }

    .success {
      color: #3d7a56;
    }
  `
}

declare global {
  interface HTMLElementTagNameMap {
    'page-admin-account': PageAdminAccount
  }
}
