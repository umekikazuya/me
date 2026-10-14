import { CITY, TIME_ZONE } from '../content/home.js'

/** 時刻帯ごとの一言。 */
export function moodAt(hour: number): string {
  if (hour < 7) return 'probably asleep'
  if (hour < 10) return 'on the train'
  if (hour < 19) return 'at work'
  return 'writing code, maybe'
}

/** `date` の出力。"Tue 02:41 JST · Tokyo" と一言。 */
export function localTime(now: Date): { text: string; mood: string } {
  const parts = new Intl.DateTimeFormat('en-US', {
    timeZone: TIME_ZONE,
    weekday: 'short',
    hour: '2-digit',
    minute: '2-digit',
    hourCycle: 'h23',
  }).formatToParts(now)
  const get = (type: Intl.DateTimeFormatPartTypes) =>
    parts.find((p) => p.type === type)?.value ?? ''
  return {
    text: `${get('weekday')} ${get('hour')}:${get('minute')} JST · ${CITY}`,
    mood: moodAt(Number(get('hour'))),
  }
}
