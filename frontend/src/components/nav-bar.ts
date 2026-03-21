import { LitElement, css, html } from 'lit'
import { customElement } from 'lit/decorators.js'

@customElement('nav-bar')
export class NavBar extends LitElement {
  render() {
    return html`
      <nav>
        <a href="/" class="brand">umekikazuya</a>
        <div class="links">
          <a href="/articles">Articles</a>
          <a href="/about">About</a>
        </div>
      </nav>
    `
  }

  static styles = css`
    :host {
      display: block;
    }

    nav {
      position: fixed;
      top: 0;
      left: 0;
      right: 0;
      z-index: 100;
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: var(--space-sm) var(--space-md);
    }

    a {
      font-family: var(--font-en);
      font-weight: 300;
      font-size: 15px;
      letter-spacing: var(--tracking-wide);
      color: var(--color-text-primary);
      text-decoration: none;
      opacity: 1;
      transition: opacity 0.2s ease;
      position: relative;
    }

    a::after {
      content: '';
      position: absolute;
      left: 0;
      bottom: -2px;
      width: 100%;
      height: 0.5px;
      background: var(--color-text-primary);
      transform: scaleX(0);
      transform-origin: left;
      transition: transform 0.3s var(--easing-smooth);
    }

    a:hover::after {
      transform: scaleX(1);
    }

    .links {
      display: flex;
      gap: 40px;
    }

    nav:hover a {
      opacity: 0.5;
    }

    nav:hover a:hover {
      opacity: 1;
    }

    @media (prefers-reduced-motion: reduce) {
      a,
      a::after {
        transition: none;
      }
    }

    @media (max-width: 640px) {
      nav {
        padding: 20px 24px;
      }

      .links {
        gap: 24px;
      }
    }
  `
}

declare global {
  interface HTMLElementTagNameMap {
    'nav-bar': NavBar
  }
}
