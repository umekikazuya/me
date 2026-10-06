const SEEN_KEY = 'me.booted'
const MIN_DURATION = 900 // ms: shorter than this reads as a flicker
const MAX_DURATION = 3000 // ms: past this, hand off to in-page loaders
const READY_HOLD = 450 // ms: let "100" sit for a beat
const LEAVE_DURATION = 600 // ms: matches #boot transition in boot.css
const REVEAL_OFFSET = 200 // ms: start the page entrance while #boot fades

export type BootOptions = {
  tasks: Array<Promise<unknown>>
  force?: boolean
  /** 画面が抜け始めたタイミングで呼ばれる。本体の登場演出をここで重ねる */
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
  } catch {
    // storage unavailable (private mode etc.): boot again next time
  }
}

let resolveRevealed = () => {}
const revealed = new Promise<void>((resolve) => {
  resolveRevealed = resolve
})

/** ブート画面が抜けて本体が見え始めたら解決する。ページ内の演出の開始合図に使う */
export function whenRevealed(): Promise<void> {
  return revealed
}

const sleep = (ms: number) =>
  new Promise<void>((resolve) => window.setTimeout(resolve, ms))

/**
 * 表示中の値を実際の進捗へ近づける。数字が飛ばないよう毎フレーム差分の一部だけ進め、
 * 残りがごく僅かになったら吸着させる。
 */
export function easeProgress(shown: number, real: number, rate = 0.1): number {
  const next = shown + (real - shown) * rate
  return real - next < 0.4 ? real : next
}

/**
 * index.html の #boot を、渡されたタスクの完了状況に合わせて進める。
 * 2回目以降のセッションや reduced-motion では即座に閉じる。
 * キー入力・クリックでスキップできる。
 */
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
        .catch(() => {
          // a failed task still counts; pages render their own error state
        })
        .finally(() => {
          done++
          real = (done / tasks.length) * 100
        }),
    ),
  )

  await Promise.race([
    Promise.all([settled, sleep(MIN_DURATION)]),
    sleep(MAX_DURATION),
    skipped,
  ])
  real = 100
  await Promise.race([sleep(READY_HOLD), skipped])

  cancelAnimationFrame(rafId)
  paint(100)
  window.removeEventListener('keydown', skip, { capture: true })
  window.removeEventListener('pointerdown', skip, { capture: true })
  markBooted()

  el.dataset.state = 'leaving'
  await sleep(REVEAL_OFFSET)
  onReveal()
  await sleep(LEAVE_DURATION - REVEAL_OFFSET)
  el.dataset.state = 'done'
}
