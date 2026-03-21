const LINE_COUNT = 3
const COLOR_PRIMARY_RGB = '44, 42, 38'
const CROSS_DURATION_MS = 12000 // 12秒で画面を横断
const REPULSION_RADIUS = 150
const REPULSION_STRENGTH = 30

interface Line {
  // 制御点の初期Y位置 (0〜1 の比率)
  y0: number
  y1: number
  y2: number
  y3: number
  // アニメーション開始時刻（オフセットでずらす）
  startTime: number
}

function rand(min: number, max: number): number {
  return min + Math.random() * (max - min)
}

function generateLine(startTime: number): Line {
  // ゆるやかなS字になるよう制御点をランダムに配置
  const base = rand(0.2, 0.8)
  return {
    y0: base + rand(-0.1, 0.1),
    y1: base + rand(-0.15, 0.15),
    y2: base + rand(-0.15, 0.15),
    y3: base + rand(-0.1, 0.1),
    startTime,
  }
}

/** ベジェ曲線上の点を返す (t: 0〜1) */
function bezierPoint(
  t: number,
  p0: number,
  p1: number,
  p2: number,
  p3: number
): number {
  const u = 1 - t
  return u * u * u * p0 + 3 * u * u * t * p1 + 3 * u * t * t * p2 + t * t * t * p3
}

export function setupAmbientLines(container: HTMLElement): () => void {
  const reduced = window.matchMedia('(prefers-reduced-motion: reduce)').matches

  const canvas = document.createElement('canvas')
  canvas.style.cssText = `
    position: absolute;
    top: 0; left: 0;
    width: 100%;
    height: 100%;
    pointer-events: none;
    z-index: 0;
  `
  container.style.position = 'relative'
  container.insertBefore(canvas, container.firstChild)

  const ctx = canvas.getContext('2d')
  if (!ctx) { canvas.remove(); return () => {} }

  const resize = () => {
    canvas.width = container.offsetWidth
    canvas.height = container.offsetHeight
  }
  resize()

  const ro = new ResizeObserver(resize)
  ro.observe(container)

  // マウス位置（container相対）
  let mouseX = -9999
  let mouseY = -9999

  const onMouseMove = (e: MouseEvent) => {
    const rect = container.getBoundingClientRect()
    mouseX = e.clientX - rect.left
    mouseY = e.clientY - rect.top
  }

  window.addEventListener('mousemove', onMouseMove, { passive: true })

  const now = performance.now()
  const lines: Line[] = Array.from({ length: LINE_COUNT }, (_, i) =>
    generateLine(now - i * (CROSS_DURATION_MS / LINE_COUNT))
  )

  let rafId: number

  const drawFrame = (time: number) => {
    const w = canvas.width
    const h = canvas.height
    ctx.clearRect(0, 0, w, h)

    for (const line of lines) {
      const elapsed = time - line.startTime
      const progress = (elapsed % CROSS_DURATION_MS) / CROSS_DURATION_MS

      // 左端から始まり右端まで進む。中心X座標
      // 線全体をオフセットとして扱う: offsetX は -w〜2w の範囲で動く
      const offsetX = (progress - 0.1) * (w * 1.2)

      // サンプル点を N 個計算し、それを Canvas パスで描画
      const STEPS = 60
      ctx.beginPath()

      for (let i = 0; i <= STEPS; i++) {
        const t = i / STEPS
        const x = offsetX + t * (w + w * 0.2) - w * 0.1

        // Y: ベジェで S字
        const rawY = bezierPoint(t, line.y0, line.y1, line.y2, line.y3) * h

        // マウス repulsion
        const dx = x - mouseX
        const dy = rawY - mouseY
        const dist = Math.sqrt(dx * dx + dy * dy)
        let repY = 0
        if (dist < REPULSION_RADIUS && dist > 0) {
          const force = (1 - dist / REPULSION_RADIUS) * REPULSION_STRENGTH
          repY = (dy / dist) * force
        }

        const y = rawY + repY

        if (i === 0) {
          ctx.moveTo(x, y)
        } else {
          ctx.lineTo(x, y)
        }
      }

      ctx.strokeStyle = `rgba(${COLOR_PRIMARY_RGB}, 0.15)`
      ctx.lineWidth = 0.5
      ctx.stroke()

      // ループ: 一周したら新しい形状に regenerate
      if (elapsed > CROSS_DURATION_MS) {
        line.y0 = rand(0.2, 0.8) + rand(-0.1, 0.1)
        line.y1 = line.y0 + rand(-0.15, 0.15)
        line.y2 = line.y0 + rand(-0.15, 0.15)
        line.y3 = line.y0 + rand(-0.1, 0.1)
        line.startTime = time
      }
    }

    if (!reduced) {
      rafId = requestAnimationFrame(drawFrame)
    }
  }

  if (!reduced) {
    rafId = requestAnimationFrame(drawFrame)
  }

  return () => {
    cancelAnimationFrame(rafId)
    window.removeEventListener('mousemove', onMouseMove)
    ro.disconnect()
    canvas.remove()
  }
}
