import { css, html, LitElement } from 'lit'
import { customElement } from 'lit/decorators.js'
import { playLeaveTransition, routeShellStyles } from './route-shell.js'
import './nav-bar.js'

@customElement('app-public-shell')
export class AppPublicShell extends LitElement {
  render() {
    return html`
      <nav-bar></nav-bar>
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
      }
    `,
  ]
}

declare global {
  interface HTMLElementTagNameMap {
    'app-public-shell': AppPublicShell
  }
}
