import type { ReactiveController, ReactiveControllerHost } from 'lit'

const NEXT_COMMAND_PAUSE_MS = 450

/** 複数のコマンドを順に打つページの進行。出力が流れ終わってから次のコマンドを出す。 */
export class Session implements ReactiveController {
  commands = 1
  outputs = 0
  private host: ReactiveControllerHost
  private timer?: number

  constructor(host: ReactiveControllerHost) {
    this.host = host
    host.addController(this)
  }

  hostDisconnected() {
    window.clearTimeout(this.timer)
  }

  /** index 番目の出力を出し、outputMs 後に次のコマンドを打ち始める。 */
  advance(index: number, total: number, outputMs: number) {
    this.reveal(index)
    this.continueAfter(index, total, outputMs)
  }

  /** index 番目までの出力を出す。次のコマンドにはまだ進まない。 */
  reveal(index: number) {
    if (index < this.outputs) return
    this.outputs = index + 1
    this.host.requestUpdate()
  }

  /** outputMs 後に index の次のコマンドを打ち始める。 */
  continueAfter(index: number, total: number, outputMs: number) {
    if (index + 1 >= total) return
    window.clearTimeout(this.timer)
    this.timer = window.setTimeout(() => {
      this.commands = Math.max(this.commands, index + 2)
      this.host.requestUpdate()
    }, outputMs + NEXT_COMMAND_PAUSE_MS)
  }
}
