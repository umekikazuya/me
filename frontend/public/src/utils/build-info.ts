/** "2d ago" のような経過時間。 */
export function timeAgo(from: Date, now = new Date()): string {
  const minutes = Math.max(
    0,
    Math.round((now.valueOf() - from.valueOf()) / 60000),
  )
  if (minutes < 60) return `${minutes}m ago`
  const hours = Math.round(minutes / 60)
  if (hours < 24) return `${hours}h ago`
  return `${Math.round(hours / 24)}d ago`
}

/** ビルド時に埋め込んだデプロイ元の commit。取れなかった場合は null。 */
export function buildInfo(now = new Date()) {
  if (!__BUILD_COMMIT__) return null
  const at = new Date(__BUILD_COMMITTED_AT__)
  return {
    commit: __BUILD_COMMIT__,
    ago: Number.isNaN(at.valueOf()) ? '' : timeAgo(at, now),
  }
}
