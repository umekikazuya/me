import { css, html, LitElement } from 'lit'
import { customElement } from 'lit/decorators.js'
import { toggleScheme } from '../utils/scheme.js'

@customElement('nav-bar')
export class NavBar extends LitElement {
  private onTheme(e: Event) {
    e.preventDefault()
    toggleScheme()
  }

  render() {
    return html`
      <nav>
        <a href="/" class="brand">umekikazuya</a>
        <div class="links">
          <a href="/articles">writing</a>
          <a href="/about">about</a>
          <button type="button" @click=${this.onTheme}>theme</button>
        </div>
      </nav>
    `
  }

  static styles = css`
    :host {
      display: block;
    }

    nav {
      max-width: var(--content-width);
      margin: 0 auto;
      padding: 32px 24px 0;
      display: flex;
      justify-content: space-between;
      font-family: var(--font-mono);
      font-size: 12px;
      color: var(--color-text-tertiary);
    }

    .links {
      display: flex;
      gap: 20px;
    }

    a,
    button {
      color: inherit;
      text-decoration: none;
      transition: color 0.2s ease;
    }

    button {
      font: inherit;
      background: none;
      border: 0;
      padding: 0;
      cursor: pointer;
    }

    a:hover,
    button:hover {
      color: var(--color-text-primary);
    }

    :focus-visible {
      outline: 1px solid var(--color-accent);
      outline-offset: 3px;
    }
  `
}

declare global {
  interface HTMLElementTagNameMap {
    'nav-bar': NavBar
  }
}
