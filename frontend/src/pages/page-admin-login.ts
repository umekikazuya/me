import { consume } from '@lit/context'
import { css, html, LitElement } from 'lit'
import { customElement, state } from 'lit/decorators.js'
import { adminFormStyles } from '../admin/admin-form-styles.js'
import { authContext } from '../contexts/auth-context.js'
import { RepositoryObserver } from '../controllers/RepositoryObserver.js'
import type { IAuthRepository } from '../domain/AuthRepository.js'
import '../components/admin/ui/me-text-input.js'

@customElement('page-admin-login')
export class PageAdminLogin extends LitElement {
  @consume({ context: authContext, subscribe: true })
  set authRepo(repo: IAuthRepository) {
    if (this._authRepo === repo) return
    this._authRepo = repo
    if (this._observer) this._observer.disconnect()
    if (repo) this._observer = new RepositoryObserver(this, repo)
  }
  get authRepo() {
    return this._authRepo
  }
  private _authRepo!: IAuthRepository
  private _observer?: RepositoryObserver

  @state()
  private passwordVisible = false

  firstUpdated() {
    this.shadowRoot
      ?.querySelector<HTMLElement>('me-text-input[name="emailAddress"]')
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
            <me-text-input
              label="メールアドレス"
              type="email"
              name="emailAddress"
              autocomplete="email"
              ?disabled=${a.loginPending}
              required
            ></me-text-input>

            <div class="password-field-container">
              <me-text-input
                label="パスワード"
                .type=${this.passwordVisible ? 'text' : 'password'}
                name="password"
                autocomplete="current-password"
                ?disabled=${a.loginPending}
                required
              ></me-text-input>
              <button
                type="button"
                class="subtle password-toggle"
                ?disabled=${a.loginPending}
                @click=${this.togglePasswordVisibility}
              >
                ${this.passwordVisible ? '隠す' : '表示'}
              </button>
            </div>

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

  private togglePasswordVisibility = () => {
    this.passwordVisible = !this.passwordVisible
  }

  private async handleSubmit(event: Event) {
    event.preventDefault()
    const form = event.target as HTMLFormElement
    const formData = new FormData(form)

    await this.authRepo.login({
      emailAddress: (formData.get('emailAddress') as string).trim(),
      password: formData.get('password') as string,
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

      .password-field-container {
        position: relative;
        display: grid;
      }

      .password-toggle {
        position: absolute;
        right: 0;
        top: 0;
        height: 20px;
        font-size: 12px;
      }

      button[type="submit"] {
        margin-top: 8px;
      }
    `,
  ]
}

declare global {
  interface HTMLElementTagNameMap {
    'page-admin-login': PageAdminLogin
  }
}
