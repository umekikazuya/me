import { toggleScheme } from './scheme.js'

const SCROLL_STEP = 80 // px

type Binding = () => void

const bindings: Record<string, Binding> = {
  j: () => window.scrollBy({ top: SCROLL_STEP }),
  k: () => window.scrollBy({ top: -SCROLL_STEP }),
  g: () => window.scrollTo({ top: 0, behavior: 'smooth' }),
  G: () =>
    window.scrollTo({
      top: document.documentElement.scrollHeight,
      behavior: 'smooth',
    }),
  t: () => toggleScheme(),
}

function isTyping(e: KeyboardEvent): boolean {
  const target = e.composedPath()[0] as HTMLElement | undefined
  if (!target) return false
  return (
    target.isContentEditable ||
    ['INPUT', 'TEXTAREA', 'SELECT'].includes(target.tagName)
  )
}

/**
 * Vim 風の最小限のキー操作。入力欄にフォーカスがあるときや修飾キー付きは無視する。
 * クリーンアップ関数を返す。
 */
export function setupKeys(): () => void {
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
