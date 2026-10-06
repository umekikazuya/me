import { LitElement } from 'lit'
import { state } from 'lit/decorators.js'
import { CursorLine } from '../utils/cursorline.js'
import '../components/shell-command.js'
import '../components/shell-prompt.js'

/**
 * 「コマンドを打つ → 出力が流れる」ページの土台。
 * <shell-command> の typed を受けて `typed` が true になったら出力を描画する。
 */
export class ShellPage extends LitElement {
  @state()
  protected typed = false

  protected cursor = new CursorLine(this)

  connectedCallback() {
    super.connectedCallback()
    this.addEventListener('typed', this.onTyped)
  }

  disconnectedCallback() {
    super.disconnectedCallback()
    this.removeEventListener('typed', this.onTyped)
  }

  private onTyped = () => {
    this.typed = true
  }
}
