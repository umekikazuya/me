import { consume } from '@lit/context'
import { css, html, LitElement } from 'lit'
import { customElement, state } from 'lit/decorators.js'
import { adminFormStyles } from '../admin/admin-form-styles.js'
import { authContext } from '../contexts/auth-context.js'
import { RepositoryObserver } from '../controllers/RepositoryObserver.js'
import type { IAuthRepository } from '../domain/AuthRepository.js'

@customElement('page-admin-login')
export class PageAdminLogin extends LitElement {
  @consume({ context: authContext, subscribe: true })
  authRepo!: IAuthRepository

  @state()
  private emailAddress = ''

  @state()
  private password = ''

  @state()
  private passwordVisible = false

  constructor() {
    super()
    new RepositoryObserver(this, this.authRepo)
  }

  firstUpdated() {
    this.shadowRoot
      ?.querySelector<HTMLInputElement>('input[name="emailAddress"]')
      ?.focus()
  }

  render() {
    const a = this.authRepo
    return html`
      <section class="container">
        <div class="card">
          <p class="eyebrow" lang="en">Admin</p>
          <h1 class="title">ログイン</h1>
          <p class="description">
            管理画面へ入るには、メールアドレスとパスワードでログインしてください。
          </p>

          ${
            a.loginNotice
              ? html`<p class="message notice">${a.loginNotice}</p>`
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
                ?disabled=${a.loginPending}
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
                  ?disabled=${a.loginPending}
                  @input=${this.handlePasswordInput}
                  required
                />
                <button
                  type="button"
                  class="subtle"
                  ?disabled=${a.loginPending}
                  @click=${this.togglePasswordVisibility}
                >
                  ${this.passwordVisible ? '隠す' : '表示'}
                </button>
              </div>
            </label>

            ${
              a.loginError
                ? html`<p class="message error">${a.loginError}</p>`
                : null
            }

            <button type="submit" ?disabled=${a.loginPending}>
              ${a.loginPending ? 'ログイン中...' : 'ログイン'}
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

  private async handleSubmit(event: Event) {
    event.preventDefault()

    await this.authRepo.login({
      emailAddress: this.emailAddress.trim(),
      password: this.password,
    })
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
