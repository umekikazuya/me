import { css, html, LitElement } from 'lit'
import { customElement, property, state } from 'lit/decorators.js'
import { adminFormStyles } from '../admin/admin-form-styles.js'
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
                  class="subtle"
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

  static styles = [
    adminFormStyles,
    css`
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
        width: min(100%, 440px);
        background: #fff;
        border: 1px solid var(--color-border);
        padding: 40px;
      }

      .description {
        margin-bottom: 28px;
      }

      form {
        display: grid;
        gap: 18px;
      }

      .password-field {
        display: grid;
        grid-template-columns: minmax(0, 1fr) auto;
        gap: 8px;
        align-items: center;
      }

      .subtle {
        min-width: 72px;
      }
    `,
  ]
}

declare global {
  interface HTMLElementTagNameMap {
    'page-admin-login': PageAdminLogin
  }
}
