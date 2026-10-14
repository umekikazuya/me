/** ステータスラインに並べるページ。 */
export const NAV = [
  { label: 'home', path: '/' },
  { label: 'writing', path: '/articles' },
  { label: 'about', path: '/about' },
] as const

/** パスに対応するページ。該当しなければ undefined。 */
export function currentNav(pathname: string) {
  const trimmed = pathname.replace(/\/+$/, '') || '/'
  return NAV.find((item) => item.path === trimmed)
}
