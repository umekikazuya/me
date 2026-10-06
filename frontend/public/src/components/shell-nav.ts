import { css, html, LitElement } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import { NAV } from '../utils/keys.js'

/** 出力の末尾に置く「1 home 2 writing 3 about」と、次の入力待ちプロンプト */
@customElement('shell-nav')
export class ShellNav extends LitElement {
  @property()
  current = ''

  render() {
    return html`
      <nav aria-label="pages">
        ${NAV.map(
          (item) => html`
            <a
              href=${item.path}
              aria-current=${item.label === this.current ? 'page' : 'false'}
            ><span class="k">${item.key}</span> ${item.label}</a>
          `,
        )}
      </nav>
      <p><span class="k">~ $ </span><span class="cursor"></span></p>
    `
  }

  static styles = css`
    :host {
      display: block;
      margin-top: 1.9em;
      animation: in 0.14s ease both;
    }

    nav {
      display: flex;
      flex-wrap: wrap;
      gap: 4px 20px;
    }

    a {
      color: var(--color-text-tertiary);
      text-decoration: none;
      transition: color 0.2s ease;
    }

    a:hover,
    a:focus-visible,
    a[aria-current='page'] {
      color: var(--color-text-primary);
      outline: none;
    }

    a[aria-current='page']::after {
      content: '';
      display: inline-block;
      width: 4px;
      height: 4px;
      margin-left: 6px;
      border-radius: 50%;
      background: var(--color-accent);
      vertical-align: 2px;
    }

    .k {
      color: var(--color-text-tertiary);
    }

    p {
      margin: 0;
    }

    .cursor {
      display: inline-block;
      width: 0.55em;
      height: 1.05em;
      vertical-align: -0.16em;
      background: var(--color-accent);
      animation: blink 1.05s steps(1) infinite;
    }

    @keyframes blink {
      50% {
        opacity: 0;
      }
    }

    @keyframes in {
      from {
        opacity: 0;
      }
    }

    @media (prefers-reduced-motion: reduce) {
      :host,
      .cursor {
        animation: none;
      }
    }
  `
}

declare global {
  interface HTMLElementTagNameMap {
    'shell-nav': ShellNav
  }
}
