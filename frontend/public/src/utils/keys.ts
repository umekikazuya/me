import { moveCursor } from './cursorline.js'
import { toggleScheme } from './scheme.js'

export const NAV = [
  { key: '1', label: 'home', path: '/' },
  { key: '2', label: 'writing', path: '/articles' },
  { key: '3', label: 'about', path: '/about' },
] as const

function isTyping(e: KeyboardEvent): boolean {
  const target = e.composedPath()[0] as HTMLElement | undefined
  if (!target) return false
  return (
    target.isContentEditable ||
    ['INPUT', 'TEXTAREA', 'SELECT'].includes(target.tagName)
  )
}

/**
 * Vim 風のキー操作。j/k/g/G で行カーソル、1/2/3 でページ移動、t でテーマ。
 * 入力欄にフォーカスがあるときや修飾キー付きは無視する。クリーンアップ関数を返す。
 */
export function setupKeys(navigate: (path: string) => void): () => void {
  const bindings: Record<string, () => void> = {
    j: () => moveCursor('next'),
    k: () => moveCursor('prev'),
    g: () => moveCursor('first'),
    G: () => moveCursor('last'),
    t: () => toggleScheme(),
  }
  for (const item of NAV) bindings[item.key] = () => navigate(item.path)

  const onKeyDown = (e: KeyboardEvent) => {
    if (e.metaKey || e.ctrlKey || e.altKey || isTyping(e)) return
    const binding = bindings[e.key]
    if (!binding) return
    e.preventDefault()
    binding()
  }
  window.addEventListener('keydown', onKeyDown)
  return () => window.removeEventListener('keydown', onKeyDown)
}
