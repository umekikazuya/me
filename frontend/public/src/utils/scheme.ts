const SCHEME_KEY = 'me.scheme'

export type Scheme = 'light' | 'dark'

export function currentScheme(): Scheme {
  const explicit = document.documentElement.getAttribute('data-scheme')
  if (explicit === 'light' || explicit === 'dark') return explicit
  return window.matchMedia('(prefers-color-scheme: dark)').matches
    ? 'dark'
    : 'light'
}

export function toggleScheme(): Scheme {
  const next: Scheme = currentScheme() === 'dark' ? 'light' : 'dark'
  document.documentElement.setAttribute('data-scheme', next)
  try {
    localStorage.setItem(SCHEME_KEY, next)
  } catch {
    // storage unavailable: the choice lasts for this page view only
  }
  return next
}
