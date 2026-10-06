import { moveCursor } from './cursorline.js'
import { toggleScheme } from './scheme.js'

function isTyping(e: KeyboardEvent): boolean {
  const target = e.composedPath()[0] as HTMLElement | undefined
  if (!target) return false
  return (
    target.isContentEditable ||
    ['INPUT', 'TEXTAREA', 'SELECT'].includes(target.tagName)
  )
}

/**
 * Vim 風のキー操作。j/k/g/G で行カーソル、t でテーマ。
 * 入力欄にフォーカスがあるときや修飾キー付きは無視する。クリーンアップ関数を返す。
 */
export function setupKeys(): () => void {
  const bindings: Record<string, () => void> = {
    j: () => moveCursor('next'),
    k: () => moveCursor('prev'),
    g: () => moveCursor('first'),
    G: () => moveCursor('last'),
    t: () => toggleScheme(),
  }
  // 押しっぱなしで連続させたいのは行移動だけ。テーマ切替は1回に限る
  const repeatable = new Set(['j', 'k'])

  const onKeyDown = (e: KeyboardEvent) => {
    if (e.metaKey || e.ctrlKey || e.altKey || isTyping(e)) return
    const binding = bindings[e.key]
    if (!binding || (e.repeat && !repeatable.has(e.key))) return
    e.preventDefault()
    binding()
  }
  window.addEventListener('keydown', onKeyDown)
  return () => window.removeEventListener('keydown', onKeyDown)
}
