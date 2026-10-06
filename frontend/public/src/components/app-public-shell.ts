import { css, html, LitElement } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import type { RouteShellElement } from './route-shell.js'
import { playLeaveTransition, routeShellStyles } from './route-shell.js'
import './nav-bar.js'
import './status-line.js'

@customElement('app-public-shell')
export class AppPublicShell extends LitElement implements RouteShellElement {
  /** ブート画面が抜けたら true。それまで本体は隠しておく */
  @property({ type: Boolean, reflect: true })
  revealed = false

  @property()
  path = '/'

  render() {
    return html`
      <div class="frame">
        <nav-bar></nav-bar>
        <main id="outlet">
          <slot></slot>
        </main>
      </div>
      <status-line .path=${this.path}></status-line>
    `
  }

  playLeaveTransition() {
    return playLeaveTransition(this.outlet)
  }

  private get outlet() {
    return this.shadowRoot?.querySelector('#outlet') as HTMLElement | null
  }

  static styles = [
    routeShellStyles,
    css`
      :host {
        display: block;
        padding-bottom: calc(var(--statusline-height) + 64px);
      }

      .frame {
        opacity: 0;
      }

      :host([revealed]) .frame {
        animation: rise 0.9s var(--easing-rise) both;
      }

      #outlet {
        display: block;
        max-width: var(--content-width);
        margin: 0 auto;
        padding-inline: 24px;
      }

      @keyframes rise {
        from {
          opacity: 0;
          transform: translateY(10px);
          filter: blur(4px);
        }
        to {
          opacity: 1;
          transform: none;
          filter: none;
        }
      }

      @media (prefers-reduced-motion: reduce) {
        :host([revealed]) .frame {
          animation: none;
          opacity: 1;
        }
      }
    `,
  ]
}

declare global {
  interface HTMLElementTagNameMap {
    'app-public-shell': AppPublicShell
  }
}
