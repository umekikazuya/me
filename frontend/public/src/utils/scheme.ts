const SCHEME_KEY = 'me.scheme'

export type Scheme = 'light' | 'dark'

/** 明示指定がなければ OS 設定に従う。 */
export function currentScheme(): Scheme {
  const explicit = document.documentElement.getAttribute('data-scheme')
  if (explicit === 'light' || explicit === 'dark') return explicit
  return window.matchMedia('(prefers-color-scheme: dark)').matches
    ? 'dark'
    : 'light'
}

/** ライト/ダークを切り替えて保存する。 */
export function toggleScheme(): Scheme {
  const next: Scheme = currentScheme() === 'dark' ? 'light' : 'dark'
  document.documentElement.setAttribute('data-scheme', next)
  try {
    localStorage.setItem(SCHEME_KEY, next)
  } catch {}
  return next
}
