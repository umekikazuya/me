/** ISO 日時を "2026.09.28" 形式に。不正・未設定は "----.--.--" */
export function formatDate(value?: string): string {
  if (!value) return '----.--.--'
  const date = new Date(value)
  if (Number.isNaN(date.valueOf())) return '----.--.--'
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${date.getFullYear()}.${pad(date.getMonth() + 1)}.${pad(date.getDate())}`
}

/** http(s) と mailto 以外のリンクは潰す */
export function sanitizeUrl(url: string): string {
  const trimmed = url.trim()
  return /^(https?|mailto):/i.test(trimmed) ? trimmed : '#'
}
