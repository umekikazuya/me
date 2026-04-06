import { css, html, LitElement } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import type { RouteShellElement } from './route-shell.js'
import { playLeaveTransition, routeShellStyles } from './route-shell.js'

@customElement('app-admin-shell')
export class AppAdminShell extends LitElement implements RouteShellElement {
  @property({ type: Boolean })
  authenticated = false

  @property()
  currentPath = '/admin'

  @property({ type: Boolean })
  busy = false

  render() {
    return html`
      <div class="layout">
        ${
          this.authenticated
            ? html`
              <aside class="sidebar">
                <a href="/admin" class=${this.navClass('/admin')}>Dashboard</a>
                <a href="/admin/profile" class=${this.navClass('/admin/profile')}
                  >Profile</a
                >
                <a href="/admin/account" class=${this.navClass('/admin/account')}
                  >Account</a
                >
              </aside>
            `
            : null
        }
        <main id="outlet">
          ${
            this.busy
              ? html`<p class="status">セッションを確認しています...</p>`
              : null
          }
          <slot></slot>
        </main>
      </div>
    `
  }

  playLeaveTransition() {
    return playLeaveTransition(this.outlet)
  }

  private get outlet() {
    return this.shadowRoot?.querySelector('#outlet') as HTMLElement | null
  }

  private navClass(path: string) {
    return this.currentPath === path ? 'active' : ''
  }

  static styles = [
    routeShellStyles,
    css`
      :host {
        display: block;
        background: linear-gradient(
          180deg,
          rgba(243, 242, 238, 0.9) 0%,
          rgba(255, 255, 255, 1) 100%
        );
      }

      .layout {
        min-height: 100dvh;
        display: grid;
        grid-template-columns: 240px minmax(0, 1fr);
      }

      .sidebar {
        display: grid;
        align-content: start;
        gap: 10px;
        padding: 32px 20px;
        border-right: 1px solid var(--color-border-light);
        background: rgba(255, 255, 255, 0.72);
        backdrop-filter: blur(18px);
      }

      .sidebar a {
        display: block;
        padding: 12px 14px;
        font-family: var(--font-en);
        letter-spacing: var(--tracking-wide);
        color: var(--color-text-secondary);
        border: 1px solid transparent;
      }

      .sidebar a.active {
        color: var(--color-text-primary);
        background: #fff;
        border-color: var(--color-border);
      }

      #outlet {
        display: block;
        min-height: 100dvh;
        padding: 48px;
      }

      .status {
        margin-bottom: 20px;
        color: var(--color-text-secondary);
        font-size: 14px;
      }

      @media (max-width: 960px) {
        .layout {
          grid-template-columns: 1fr;
        }

        .sidebar {
          grid-auto-flow: column;
          grid-auto-columns: max-content;
          overflow-x: auto;
          border-right: 0;
          border-bottom: 1px solid var(--color-border-light);
        }

        #outlet {
          padding: 32px 24px 48px;
        }
      }
    `,
  ]
}

declare global {
  interface HTMLElementTagNameMap {
    'app-admin-shell': AppAdminShell
  }
}
