const LAST_VISIT_KEY = 'me.lastVisit'

/** "Mon Oct 6 01:41" 形式。 */
export function formatLogin(date: Date): string {
  const pad = (n: number) => String(n).padStart(2, '0')
  const weekday = date.toLocaleString('en-US', { weekday: 'short' })
  const month = date.toLocaleString('en-US', { month: 'short' })
  return `${weekday} ${month} ${date.getDate()} ${pad(date.getHours())}:${pad(date.getMinutes())}`
}

/** 前回の訪問日時を返し、今回の訪問を記録する。初回や保存できない環境では null。 */
export function takeLastLogin(now = new Date()): Date | null {
  try {
    const previous = localStorage.getItem(LAST_VISIT_KEY)
    localStorage.setItem(LAST_VISIT_KEY, now.toISOString())
    const date = previous ? new Date(previous) : null
    return date && !Number.isNaN(date.valueOf()) ? date : null
  } catch {
    return null
  }
}
