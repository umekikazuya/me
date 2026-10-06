import { css, html, LitElement, nothing } from 'lit'
import { customElement, property, state } from 'lit/decorators.js'
import { whenRevealed } from '../utils/boot.js'

const START_DELAY = 500 // ms after the page is revealed
const KEY_MIN = 50 // ms
const KEY_JITTER = 60 // ms: uneven gaps read as a person typing
const OUTPUT_DELAY = 300 // ms

/**
 * "$ <command>" をタイプしてから出力を表示する1行ターミナル。
 * output が空の間は打ち終えてもカーソルで待つので、データ取得と並走できる。
 */
@customElement('term-prompt')
export class TermPrompt extends LitElement {
  @property()
  command = ''

  @property({ type: Array })
  output: string[] = []

  @state()
  private typed = ''

  @state()
  private typingDone = false

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

    if (window.matchMedia('(prefers-reduced-motion: reduce)').matches) {
      this.typed = this.command
      this.typingDone = true
      return
    }

    await whenRevealed()
    await wait(START_DELAY)
    for (const ch of this.command) {
      if (run !== this.run) return
      this.typed += ch
      await wait(KEY_MIN + Math.random() * KEY_JITTER)
    }
    await wait(OUTPUT_DELAY)
    if (run === this.run) this.typingDone = true
  }

  render() {
    const showOutput = this.typingDone && this.output.length > 0
    // white-space: pre-wrap なので <p> の中に余計な空白を入れない
    const cursor = html`<span class="cursor"></span>`
    return html`<p><span class="pr">$ </span>${this.typed}${showOutput ? nothing : cursor}</p>${
      showOutput
        ? html`<p class="out">${this.output.join('\n')}</p><p><span class="pr">$ </span>${cursor}</p>`
        : nothing
    }`
  }

  static styles = css`
    :host {
      display: block;
      min-height: 7.6em;
      font-family: var(--font-mono);
      font-size: 12.5px;
      line-height: 1.9;
    }

    p {
      margin: 0;
      white-space: pre-wrap;
    }

    .pr,
    .out {
      color: var(--color-text-tertiary);
    }

    .out {
      margin-bottom: 1.9em;
    }

    .cursor {
      display: inline-block;
      width: 0.5em;
      height: 1em;
      vertical-align: -0.12em;
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
    'term-prompt': TermPrompt
  }
}
