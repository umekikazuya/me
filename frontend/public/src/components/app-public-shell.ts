import { css, html, LitElement } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import type { RouteShellElement } from './route-shell.js'
import { playLeaveTransition, routeShellStyles } from './route-shell.js'
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
        opacity: 1;
      }

      #outlet {
        display: block;
        max-width: var(--content-width);
        margin: 0 auto;
        padding-inline: 24px;
        padding-top: clamp(64px, 14vh, 128px);
      }

    `,
  ]
}

declare global {
  interface HTMLElementTagNameMap {
    'app-public-shell': AppPublicShell
  }
}
