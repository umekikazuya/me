const EASING = 'cubic-bezier(0.22, 1, 0.36, 1)'

function isReducedMotion(): boolean {
  return window.matchMedia('(prefers-reduced-motion: reduce)').matches
}

/**
 * Intersection Observer ベースのコンテンツ出現アニメーション。
 * stagger=true のとき、各要素に 80ms ずつ遅延を付ける。
 * クリーンアップ関数を返す。
 */
export function setupReveal(els: Element[], stagger = false): () => void {
  if (isReducedMotion()) return () => {}

  els.forEach((el, i) => {
    const el_ = el as HTMLElement
    el_.style.opacity = '0'
    el_.style.transform = 'translateY(20px)'
    const delay = stagger ? `${i * 80}ms` : '0ms'
    el_.style.transition = `opacity 0.8s ${EASING} ${delay}, transform 0.8s ${EASING} ${delay}`
  })

  const observer = new IntersectionObserver(
    (entries) => {
      entries.forEach((entry) => {
        if (!entry.isIntersecting) return
        const el = entry.target as HTMLElement
        el.style.opacity = '1'
        el.style.transform = 'translateY(0)'
        observer.unobserve(entry.target)
      })
    },
    { threshold: 0.15 },
  )

  for (const el of els) observer.observe(el)
  return () => observer.disconnect()
}

/**
 * スクロールで通り過ぎた要素を opacity 0.3 にフェード。
 * スクロールバックで復帰。クリーンアップ関数を返す。
 */
export function setupFade(els: Element[]): () => void {
  if (isReducedMotion()) return () => {}

  const handler = () => {
    for (const el of els) {
      const rect = el.getBoundingClientRect()
      const faded = rect.top < -100
      ;(el as HTMLElement).style.transition = 'opacity 0.6s ease-out'
      ;(el as HTMLElement).style.opacity = faded ? '0.3' : '1'
    }
  }

  window.addEventListener('scroll', handler, { passive: true })
  handler()
  return () => window.removeEventListener('scroll', handler)
}

/**
 * スクロール量に応じて body の背景色を #f3f2ee → #ffffff へ補間。
 * クリーンアップ関数を返す。
 */
export function setupBackgroundShift(): () => void {
  if (isReducedMotion()) return () => {}

  const from = { r: 0x0d, g: 0x0d, b: 0x0c }
  const to = { r: 0x16, g: 0x15, b: 0x14 }

  const handler = () => {
    const max = document.documentElement.scrollHeight - window.innerHeight
    const progress = max > 0 ? Math.min(window.scrollY / max, 1) : 0
    const r = Math.round(from.r + (to.r - from.r) * progress)
    const g = Math.round(from.g + (to.g - from.g) * progress)
    const b = Math.round(from.b + (to.b - from.b) * progress)
    document.body.style.backgroundColor = `rgb(${r}, ${g}, ${b})`
  }

  window.addEventListener('scroll', handler, { passive: true })
  handler() // 初期値を即時適用
  return () => {
    window.removeEventListener('scroll', handler)
    document.body.style.backgroundColor = ''
  }
}
