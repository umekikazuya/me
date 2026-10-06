import { css, html, LitElement, nothing } from 'lit'
import { customElement, property, state } from 'lit/decorators.js'
import { whenRevealed } from '../utils/boot.js'

const START_DELAY = 260 // ms after the page is shown
const KEY_MIN = 38 // ms
const KEY_JITTER = 46 // ms: uneven gaps read as a person typing
const RETURN_DELAY = 220 // ms between the last key and stdout

/**
 * "~ $ <command>" をタイプする1行。打ち終えたら `typed` イベントを出す。
 * ページはこれを合図に出力行を描画する。
 */
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
    const wait = (ms: number) =>
      new Promise((resolve) => window.setTimeout(resolve, ms))

    if (!window.matchMedia('(prefers-reduced-motion: reduce)').matches) {
      await whenRevealed()
      await wait(START_DELAY)
      for (const ch of this.command) {
        if (run !== this.run) return
        this.shown += ch
        await wait(KEY_MIN + Math.random() * KEY_JITTER)
      }
      await wait(RETURN_DELAY)
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
      min-height: 1.9em;
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
