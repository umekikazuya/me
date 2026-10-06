import { html, nothing } from 'lit'
import { customElement } from 'lit/decorators.js'
import { pageStyles } from '../styles/page-styles.js'
import { ShellPage } from './shell-page.js'

@customElement('page-not-found')
export class PageNotFound extends ShellPage {
  render() {
    const path = window.location.pathname
    return html`
      <shell-command command=${`cd ${path}`}></shell-command>
      ${
        this.typed
          ? html`
            <p class="ln m">cd: no such file or directory: ${path}</p>
            <shell-nav></shell-nav>
          `
          : nothing
      }
    `
  }

  static styles = pageStyles
}

declare global {
  interface HTMLElementTagNameMap {
    'page-not-found': PageNotFound
  }
}
