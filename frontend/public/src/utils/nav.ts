/** サイト内のページ。ステータスバーにタブとして並ぶ */
export const NAV = [
  { label: 'home', path: '/' },
  { label: 'writing', path: '/articles' },
  { label: 'about', path: '/about' },
] as const

/** パスに対応するタブ。該当しなければ undefined（404 など） */
export function currentNav(pathname: string) {
  const trimmed = pathname.replace(/\/+$/, '') || '/'
  return NAV.find((item) => item.path === trimmed)
}
