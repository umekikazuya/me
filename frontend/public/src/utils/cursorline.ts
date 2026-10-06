import type { ReactiveController, ReactiveControllerHost } from 'lit'

type Host = ReactiveControllerHost & HTMLElement

/** 現在表示中のページ。j/k はこのページの行だけを動く */
let active: CursorLine | null = null

/**
 * Neovim の cursorline 相当。ページ内の選択可能な行（a.ln / button.ln）を
 * j/k でフォーカス移動させ、Enter はブラウザ標準のリンク・ボタン動作に任せる。
 */
export class CursorLine implements ReactiveController {
  private host: Host

  constructor(host: Host) {
    this.host = host
    host.addController(this)
  }

  hostConnected() {
    active = this
  }

  hostDisconnected() {
    if (active === this) active = null
  }

  private get lines(): HTMLElement[] {
    const root = this.host.shadowRoot ?? this.host
    return Array.from(root.querySelectorAll<HTMLElement>('a.ln, button.ln'))
  }

  private get index(): number {
    const focused = (this.host.shadowRoot ?? document).activeElement
    return this.lines.indexOf(focused as HTMLElement)
  }

  move(to: 'next' | 'prev' | 'first' | 'last') {
    const lines = this.lines
    if (lines.length === 0) return
    const current = this.index
    const target = {
      next: current < 0 ? 0 : Math.min(current + 1, lines.length - 1),
      prev: current < 0 ? 0 : Math.max(current - 1, 0),
      first: 0,
      last: lines.length - 1,
    }[to]
    lines[target].focus({ preventScroll: true })
    lines[target].scrollIntoView({ block: 'nearest' })
  }
}

export function moveCursor(to: 'next' | 'prev' | 'first' | 'last') {
  active?.move(to)
}
