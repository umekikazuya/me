import { css, html, LitElement } from 'lit'
import { customElement } from 'lit/decorators.js'
import { playLeaveTransition, routeShellStyles } from './route-shell.js'

@customElement('app-admin-shell')
export class AppAdminShell extends LitElement {
  render() {
    return html`
      <main id="outlet">
        <slot></slot>
      </main>
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
      }

      #outlet {
        display: block;
        min-height: 100dvh;
      }
    `,
  ]
}

declare global {
  interface HTMLElementTagNameMap {
    'app-admin-shell': AppAdminShell
  }
}
