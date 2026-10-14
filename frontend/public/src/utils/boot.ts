const SEEN_KEY = 'me.booted'
const MIN_DURATION_MS = 900
const MAX_DURATION_MS = 3000
const READY_HOLD_MS = 450
const LEAVE_DURATION_MS = 600
const REVEAL_OFFSET_MS = 200

export type BootOptions = {
  tasks: Array<Promise<unknown>>
  force?: boolean
  onReveal?: () => void
}

function isReducedMotion(): boolean {
  return window.matchMedia('(prefers-reduced-motion: reduce)').matches
}

function hasBooted(): boolean {
  try {
    return sessionStorage.getItem(SEEN_KEY) === '1'
  } catch {
    return false
  }
}

function markBooted() {
  try {
    sessionStorage.setItem(SEEN_KEY, '1')
  } catch {}
}

let resolveRevealed = () => {}
const revealed = new Promise<void>((resolve) => {
  resolveRevealed = resolve
})

/** ブート画面が抜け始めたら解決する。 */
export function whenRevealed(): Promise<void> {
  return revealed
}

const sleep = (ms: number) =>
  new Promise<void>((resolve) => window.setTimeout(resolve, ms))

/** 表示値を実際の進捗へ一定割合ずつ近づける。差が僅かなら吸着する。 */
export function easeProgress(shown: number, real: number, rate = 0.1): number {
  const next = shown + (real - shown) * rate
  return real - next < 0.4 ? real : next
}

/** index.html の #boot をタスクの完了に合わせて進め、閉じる。 */
export async function runBoot({
  tasks,
  force = false,
  onReveal: notify = () => {},
}: BootOptions) {
  const onReveal = () => {
    resolveRevealed()
    notify()
  }
  const el = document.getElementById('boot')
  const num = el?.querySelector<HTMLElement>('.boot-num')
  const bar = el?.querySelector<HTMLElement>('.boot-bar')
  if (!el || !num || !bar) {
    onReveal()
    return
  }

  if ((!force && hasBooted()) || isReducedMotion()) {
    el.dataset.state = 'done'
    onReveal()
    return
  }

  document.documentElement.removeAttribute('data-boot')
  el.dataset.state = 'running'
  let skip = () => {}
  const skipped = new Promise<void>((resolve) => {
    skip = resolve
  })
  window.addEventListener('keydown', skip, { capture: true, once: true })
  window.addEventListener('pointerdown', skip, { capture: true, once: true })

  let real = 0
  let shown = 0
  let rafId = 0
  const paint = (value: number) => {
    num.textContent = String(Math.round(value))
    bar.style.transform = `scaleX(${value / 100})`
  }
  const tick = () => {
    shown = easeProgress(shown, real)
    paint(shown)
    rafId = requestAnimationFrame(tick)
  }
  tick()

  let done = 0
  const settled = Promise.all(
    tasks.map((task) =>
      task
        .catch(() => {})
        .finally(() => {
          done++
          real = (done / tasks.length) * 100
        }),
    ),
  )

  await Promise.race([
    Promise.all([settled, sleep(MIN_DURATION_MS)]),
    sleep(MAX_DURATION_MS),
    skipped,
  ])
  real = 100
  await Promise.race([sleep(READY_HOLD_MS), skipped])

  cancelAnimationFrame(rafId)
  paint(100)
  window.removeEventListener('keydown', skip, { capture: true })
  window.removeEventListener('pointerdown', skip, { capture: true })
  markBooted()

  el.dataset.state = 'leaving'
  await sleep(REVEAL_OFFSET_MS)
  onReveal()
  await sleep(LEAVE_DURATION_MS - REVEAL_OFFSET_MS)
  el.dataset.state = 'done'
}
