import { css, html, LitElement } from 'lit'
import { customElement, property, state } from 'lit/decorators.js'
import type { AdminLoginInput } from '../admin/types.js'

@customElement('page-admin-login')
export class PageAdminLogin extends LitElement {
  @property({ type: Boolean })
  submitting = false

  @property()
  errorMessage = ''

  @property()
  noticeMessage = ''

  @state()
  private emailAddress = ''

  @state()
  private password = ''

  @state()
  private passwordVisible = false

  firstUpdated() {
    this.shadowRoot
      ?.querySelector<HTMLInputElement>('input[name="emailAddress"]')
      ?.focus()
  }

  render() {
    return html`
      <section class="container">
        <div class="card">
          <p class="eyebrow" lang="en">Admin</p>
          <h1 class="title">ログイン</h1>
          <p class="description">
            管理画面へ入るには、メールアドレスとパスワードでログインしてください。
          </p>

          ${
            this.noticeMessage
              ? html`<p class="message notice">${this.noticeMessage}</p>`
              : null
          }

          <form @submit=${this.handleSubmit}>
            <label class="field">
              <span>メールアドレス</span>
              <input
                type="email"
                name="emailAddress"
                autocomplete="email"
                .value=${this.emailAddress}
                ?disabled=${this.submitting}
                @input=${this.handleEmailInput}
                required
              />
            </label>

            <label class="field">
              <span>パスワード</span>
              <div class="password-field">
                <input
                  type=${this.passwordVisible ? 'text' : 'password'}
                  name="password"
                  autocomplete="current-password"
                  .value=${this.password}
                  ?disabled=${this.submitting}
                  @input=${this.handlePasswordInput}
                  required
                />
                <button
                  type="button"
                  class="ghost"
                  ?disabled=${this.submitting}
                  @click=${this.togglePasswordVisibility}
                >
                  ${this.passwordVisible ? '隠す' : '表示'}
                </button>
              </div>
            </label>

            ${
              this.errorMessage
                ? html`<p class="message error">${this.errorMessage}</p>`
                : null
            }

            <button type="submit" ?disabled=${this.submitting}>
              ${this.submitting ? 'ログイン中...' : 'ログイン'}
            </button>
          </form>
        </div>
      </section>
    `
  }

  private handleEmailInput(event: Event) {
    this.emailAddress = (event.target as HTMLInputElement).value
  }

  private handlePasswordInput(event: Event) {
    this.password = (event.target as HTMLInputElement).value
  }

  private togglePasswordVisibility = () => {
    this.passwordVisible = !this.passwordVisible
  }

  private handleSubmit(event: Event) {
    event.preventDefault()

    const detail: AdminLoginInput = {
      emailAddress: this.emailAddress.trim(),
      password: this.password,
    }

    this.dispatchEvent(
      new CustomEvent<AdminLoginInput>('admin-login-submit', {
        detail,
        bubbles: true,
        composed: true,
      }),
    )
  }

  static styles = css`
    :host {
      display: block;
      min-height: 100dvh;
    }

    .container {
      min-height: 100dvh;
      display: grid;
      place-items: center;
      padding: var(--space-lg) var(--space-md);
    }

    .card {
      width: min(100%, 480px);
      background: rgba(255, 255, 255, 0.9);
      border: 1px solid var(--color-border-light);
      padding: 40px;
      box-shadow: 0 24px 80px rgba(44, 42, 38, 0.08);
      backdrop-filter: blur(20px);
    }

    .eyebrow {
      font-family: var(--font-en);
      font-size: 14px;
      letter-spacing: var(--tracking-wider);
      color: var(--color-text-tertiary);
      margin-bottom: 12px;
    }

    .title {
      font-family: var(--font-jp);
      font-weight: 300;
      font-size: 28px;
      margin-bottom: 12px;
    }

    .description {
      color: var(--color-text-secondary);
      line-height: 1.8;
      margin-bottom: 32px;
    }

    form {
      display: grid;
      gap: 20px;
    }

    .password-field {
      display: grid;
      grid-template-columns: minmax(0, 1fr) auto;
      gap: 10px;
      align-items: center;
    }

    .field {
      display: grid;
      gap: 8px;
      font-size: 14px;
      color: var(--color-text-secondary);
    }

    input {
      width: 100%;
      border: 1px solid var(--color-border);
      background: #fff;
      padding: 14px 16px;
      font: inherit;
      color: var(--color-text-primary);
    }

    input:focus {
      outline: none;
      border-color: var(--color-text-primary);
    }

    button {
      border: 0;
      background: var(--color-text-primary);
      color: #fff;
      padding: 14px 18px;
      font: inherit;
      cursor: pointer;
      transition: opacity 0.2s ease;
    }

    .ghost {
      background: transparent;
      border: 1px solid var(--color-border);
      color: var(--color-text-secondary);
      min-width: 72px;
    }

    button:hover {
      opacity: 0.85;
    }

    button:disabled {
      opacity: 0.5;
      cursor: wait;
    }

    .message {
      font-size: 14px;
      line-height: 1.7;
    }

    .error {
      color: #9a3f3f;
    }

    .notice {
      color: #5a6b85;
      background: rgba(90, 107, 133, 0.08);
      padding: 12px 14px;
      border-left: 2px solid rgba(90, 107, 133, 0.35);
    }
  `
}

declare global {
  interface HTMLElementTagNameMap {
    'page-admin-login': PageAdminLogin
  }
}
