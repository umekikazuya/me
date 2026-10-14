import { css, html, LitElement, nothing } from 'lit'
import { customElement, property, state } from 'lit/decorators.js'
import { whenRevealed } from '../utils/boot.js'

const START_DELAY_MS = 260
const KEY_MIN_MS = 38
const KEY_JITTER_MS = 46
const RETURN_DELAY_MS = 220

/** コマンドを1文字ずつタイプし、打ち終えたら `typed` を発火する。 */
@customElement('shell-command')
export class ShellCommand extends LitElement {
  @property()
  command = ''

  @state()
  private shown = ''

  @state()
  private done = false

  private run = 0

  connectedCallback() {
    super.connectedCallback()
    void this.play()
  }

  disconnectedCallback() {
    super.disconnectedCallback()
    this.run++
  }

  private async play() {
    const run = ++this.run
    this.shown = ''
    this.done = false
    const wait = (ms: number) =>
      new Promise((resolve) => window.setTimeout(resolve, ms))

    if (!window.matchMedia('(prefers-reduced-motion: reduce)').matches) {
      await whenRevealed()
      await wait(START_DELAY_MS)
      for (const ch of this.command) {
        if (run !== this.run) return
        this.shown += ch
        await wait(KEY_MIN_MS + Math.random() * KEY_JITTER_MS)
      }
      await wait(RETURN_DELAY_MS)
    }
    if (run !== this.run) return
    this.shown = this.command
    this.done = true
    this.dispatchEvent(new Event('typed', { bubbles: true, composed: true }))
  }

  render() {
    return html`<span class="p">~ $ </span>${this.shown}${
      this.done ? nothing : html`<span class="cursor"></span>`
    }`
  }

  static styles = css`
    :host {
      display: block;
      min-height: calc(var(--output-line-height) * 1em);
      margin-bottom: var(--output-gap);
      white-space: pre-wrap;
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

    @media (prefers-reduced-motion: reduce) {
      .cursor {
        animation: none;
      }
    }
  `
}

declare global {
  interface HTMLElementTagNameMap {
    'shell-command': ShellCommand
  }
}
