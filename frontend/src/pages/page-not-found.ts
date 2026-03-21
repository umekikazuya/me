import { LitElement, css, html } from 'lit'
import { customElement } from 'lit/decorators.js'

@customElement('page-not-found')
export class PageNotFound extends LitElement {
  render() {
    return html`
      <div class="container">
        <p class="code">404</p>
        <p class="message" lang="en">Page not found.</p>
        <a href="/" class="back" lang="en">Return home</a>
      </div>
    `
  }

  static styles = css`
    :host {
      display: block;
      padding-top: 80px;
    }

    .container {
      max-width: 640px;
      margin: 0 auto;
      padding: 120px 32px;
    }

    .code {
      font-family: var(--font-en);
      font-weight: 300;
      font-size: 72px;
      letter-spacing: var(--tracking-wide);
      color: var(--color-text-tertiary);
      margin: 0 0 16px;
    }

    .message {
      font-family: var(--font-jp);
      font-weight: 200;
      font-size: 15px;
      letter-spacing: var(--tracking-wide);
      color: var(--color-text-secondary);
      margin: 0 0 48px;
    }

    .back {
      font-family: var(--font-en);
      font-weight: 300;
      font-size: 14px;
      letter-spacing: var(--tracking-wide);
      color: var(--color-text-secondary);
      text-decoration: none;
      transition: opacity 0.2s ease;
    }

    .back:hover {
      opacity: 0.6;
    }

    @media (prefers-reduced-motion: reduce) {
      .back {
        transition: none;
      }
    }
  `
}

declare global {
  interface HTMLElementTagNameMap {
    'page-not-found': PageNotFound
  }
}
