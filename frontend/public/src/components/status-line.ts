import type { PropertyValues } from 'lit'
import { css, html, LitElement } from 'lit'
import { customElement, property, state } from 'lit/decorators.js'
import { currentNav, NAV } from '../utils/nav.js'

const TYPE_INTERVAL = 22 // ms per character of ":cd <path>"
const COMMAND_HOLD = 200 // ms to leave the finished command visible

type Mode = 'NORMAL' | 'COMMAND'

/** URL のパスをホーム基準の表記に変える。 "/" → "~", "/articles" → "~/articles" */
export function toHomePath(pathname: string): string {
  const trimmed = pathname.replace(/\/+$/, '')
  return trimmed ? `~${trimmed}` : '~'
}

/** スクロール位置を Vim のルーラー風に表す */
export function toRuler(scrollY: number, max: number): string {
  if (max <= 0 || scrollY < 4) return 'top'
  if (scrollY >= max - 4) return 'bot'
  return `${Math.round((scrollY / max) * 100)}%`
}

@customElement('status-line')
export class StatusLine extends LitElement {
  /** 現在の location.pathname */
  @property()
  path = '/'

  @state()
  private mode: Mode = 'NORMAL'

  @state()
  private display = '~'

  @state()
  private ruler = 'top'

  private typing = 0

  private onScroll = () => {
    const max = document.documentElement.scrollHeight - window.innerHeight
    this.ruler = toRuler(window.scrollY, max)
  }

  connectedCallback() {
    super.connectedCallback()
    window.addEventListener('scroll', this.onScroll, { passive: true })
    this.onScroll()
  }

  disconnectedCallback() {
    super.disconnectedCallback()
    window.removeEventListener('scroll', this.onScroll)
    this.typing++
  }

  protected willUpdate(changed: PropertyValues<this>) {
    if (!changed.has('path')) return
    // 初回はそのまま表示し、遷移時だけ :cd をタイプする
    if (changed.get('path') === undefined) this.display = toHomePath(this.path)
    else void this.typeCommand(toHomePath(this.path))
  }

  private async typeCommand(target: string) {
    const run = ++this.typing
    const reduced = window.matchMedia(
      '(prefers-reduced-motion: reduce)',
    ).matches
    if (reduced) {
      this.display = target
      return
    }

    const wait = (ms: number) =>
      new Promise((resolve) => window.setTimeout(resolve, ms))
    const command = `:cd ${target}`
    this.mode = 'COMMAND'
    for (let i = 1; i <= command.length; i++) {
      if (run !== this.typing) return
      this.display = command.slice(0, i)
      await wait(TYPE_INTERVAL)
    }
    await wait(COMMAND_HOLD)
    if (run !== this.typing) return
    this.display = target
    this.mode = 'NORMAL'
  }

  render() {
    return html`
      <div class="inner">
        <nav aria-label="pages">
          ${NAV.map(
            (item) => html`
              <a
                href=${item.path}
                aria-current=${item === currentNav(this.path) ? 'page' : 'false'}
              >${item.label}</a>
            `,
          )}
        </nav>
        <span class="path" data-mode=${this.mode} aria-hidden="true">${this.display}</span>
        <span class="ruler" aria-hidden="true">${this.ruler}</span>
      </div>
    `
  }

  static styles = css`
    :host {
      position: fixed;
      inset-inline: 0;
      bottom: 0;
      z-index: 100;
      display: block;
      background: var(--color-bg-deep);
      padding-bottom: env(safe-area-inset-bottom, 0px);
      font-family: var(--font-mono);
      font-size: 12px;
      color: var(--color-text-tertiary);
      transition: background-color 0.4s ease;
    }

    .inner {
      max-width: var(--content-width);
      height: var(--statusline-height);
      margin: 0 auto;
      padding-inline: 24px;
      display: flex;
      align-items: center;
      gap: 14px;
      border-top: 1px solid var(--color-border);
    }

    nav {
      display: flex;
      gap: 16px;
    }

    a {
      color: inherit;
      text-decoration: none;
      transition: color 0.2s ease;
    }

    a:hover,
    a[aria-current='page'] {
      color: var(--color-text-primary);
    }

    a[aria-current='page']::after {
      content: '';
      display: inline-block;
      width: 4px;
      height: 4px;
      margin-left: 5px;
      border-radius: 50%;
      background: var(--color-accent);
      vertical-align: 2px;
    }

    a:focus-visible {
      outline: 1px solid var(--color-accent);
      outline-offset: 3px;
    }

    .path {
      margin-left: auto;
      color: var(--color-text-primary);
      white-space: nowrap;
      overflow: hidden;
      min-width: 0;
    }


    .path[data-mode='COMMAND'] {
      color: var(--color-accent);
    }

    .ruler {
      min-width: 4ch;
      text-align: right;
      font-variant-numeric: tabular-nums;
    }
  `
}

declare global {
  interface HTMLElementTagNameMap {
    'status-line': StatusLine
  }
}
