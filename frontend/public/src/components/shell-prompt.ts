import { css, html, LitElement } from 'lit'
import { customElement } from 'lit/decorators.js'

/** 出力末尾の入力待ちプロンプト。 */
@customElement('shell-prompt')
export class ShellPrompt extends LitElement {
  render() {
    return html`<span class="p">~ $ </span><span class="cursor"></span>`
  }

  static styles = css`
    :host {
      display: block;
      margin-top: var(--session-gap);
      animation: in 0.14s ease both;
    }

    .p {
      color: var(--color-text-tertiary);
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
    'shell-prompt': ShellPrompt
  }
}
