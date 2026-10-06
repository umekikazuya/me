import { css, html, LitElement } from 'lit'
import { customElement, property, state } from 'lit/decorators.js'

const FRAMES = '⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏'
const FRAME_INTERVAL = 80 // ms

@customElement('load-spinner')
export class LoadSpinner extends LitElement {
  @property()
  label = 'loading'

  @state()
  private frame = 0

  private timer?: number

  connectedCallback() {
    super.connectedCallback()
    if (window.matchMedia('(prefers-reduced-motion: reduce)').matches) return
    this.timer = window.setInterval(() => {
      this.frame = (this.frame + 1) % FRAMES.length
    }, FRAME_INTERVAL)
  }

  disconnectedCallback() {
    super.disconnectedCallback()
    window.clearInterval(this.timer)
  }

  render() {
    return html`<span class="sp" aria-hidden="true">${FRAMES[this.frame]}</span>
      <span role="status">${this.label}</span>`
  }

  static styles = css`
    :host {
      display: inline-flex;
      gap: 8px;
      font-family: var(--font-mono);
      font-size: 12px;
      color: var(--color-text-tertiary);
    }

    .sp {
      color: var(--color-accent);
    }
  `
}

declare global {
  interface HTMLElementTagNameMap {
    'load-spinner': LoadSpinner
  }
}
