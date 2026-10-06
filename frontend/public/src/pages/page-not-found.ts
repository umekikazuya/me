import { css, html, LitElement } from 'lit'
import { customElement } from 'lit/decorators.js'
import { pageStyles } from '../styles/page-styles.js'
import '../components/term-prompt.js'

@customElement('page-not-found')
export class PageNotFound extends LitElement {
  render() {
    const path = window.location.pathname
    return html`
      <header class="head">
        <h1>404</h1>
      </header>
      <term-prompt
        command=${`cd ${path}`}
        .output=${[`cd: no such file or directory: ${path}`]}
      ></term-prompt>
      <a href="/" class="more">cd ~ →</a>
    `
  }

  static styles = [
    pageStyles,
    css`
      .head {
        padding-block: clamp(96px, 18vh, 160px) 56px;
      }

      h1 {
        margin: 0;
        font-family: var(--font-mono);
        font-size: 22px;
        font-weight: 400;
      }
    `,
  ]
}

declare global {
  interface HTMLElementTagNameMap {
    'page-not-found': PageNotFound
  }
}
