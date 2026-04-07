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
        /* デザイントークンを admin 用に上書き */
        --font-en: system-ui, -apple-system, sans-serif;
        --font-jp: "Noto Sans JP", system-ui, sans-serif;
        --color-bg-top: #f5f5f5;
        --color-bg-bottom: #ffffff;
        --color-text-primary: #1a1a1a;
        --color-text-secondary: #6b6b6b;
        --color-text-tertiary: #9a9a9a;
        --color-border: #d9d9d9;
        --color-border-light: #e8e8e8;
        --color-surface: #f0f0f0;
        --tracking-wide: 0.02em;
        --tracking-wider: 0.04em;
        /* admin 専用トークン */
        --admin-accent: #0057b8;
        --admin-accent-hover: #004494;
        --admin-sidebar-width: 220px;

        display: block;
        background: var(--color-bg-top);
        font-family: var(--font-jp);
      }

      .layout {
        min-height: 100dvh;
        display: grid;
        grid-template-columns: var(--admin-sidebar-width) minmax(0, 1fr);
      }

      .sidebar {
        display: grid;
        align-content: start;
        gap: 4px;
        padding: 24px 12px;
        border-right: 1px solid var(--color-border);
        background: #ffffff;
      }

      .sidebar a {
        display: block;
        padding: 10px 12px;
        font-family: var(--font-jp);
        font-size: 14px;
        font-weight: 400;
        letter-spacing: var(--tracking-wide);
        color: var(--color-text-secondary);
        border-radius: 4px;
        transition: background 0.15s ease, color 0.15s ease;
      }

      .sidebar a:hover {
        background: var(--color-surface);
        color: var(--color-text-primary);
      }

      .sidebar a.active {
        background: #e8f0fb;
        color: var(--admin-accent);
        font-weight: 500;
      }

      #outlet {
        display: block;
        min-height: 100dvh;
        padding: 40px 48px;
        background: #ffffff;
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
          border-bottom: 1px solid var(--color-border);
          padding: 12px 16px;
        }

        #outlet {
          padding: 28px 24px 48px;
          background: var(--color-bg-top);
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
