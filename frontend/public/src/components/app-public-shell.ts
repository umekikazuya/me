import { css, html, LitElement } from 'lit'
import { customElement } from 'lit/decorators.js'
import type { RouteShellElement } from './route-shell.js'
import { playLeaveTransition, routeShellStyles } from './route-shell.js'
import './nav-bar.js'

@customElement('app-public-shell')
export class AppPublicShell extends LitElement implements RouteShellElement {
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
        animation: entrance 2s cubic-bezier(0.22, 1, 0.36, 1) both;
      }

      #outlet {
        display: block;
      }

      @keyframes entrance {
        from {
          opacity: 0;
          transform: translateY(10px);
        }
        to {
          opacity: 1;
          transform: translateY(0);
        }
      }

      @media (prefers-reduced-motion: reduce) {
        :host {
          animation: none;
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
