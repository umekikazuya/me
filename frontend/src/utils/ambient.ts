const LINE_COUNT = 3
const COLOR_PRIMARY_RGB = '209, 205, 199'
const CROSS_DURATION_MS = 12000 // 12秒で画面を横断
const REPULSION_RADIUS = 150
const REPULSION_STRENGTH = 30

interface Line {
  // 制御点の初期Y位置 (0〜1 の比率)
  y0: number
  y1: number
  y2: number
  // 横断の開始時刻
  startTime: number
  // 遅延
  delay: number
}

function createLine(now: number, delay = 0): Line {
  return {
    y0: Math.random(),
    y1: Math.random(),
    y2: Math.random(),
    startTime: now + delay,
    delay,
  }
}

export function setupAmbientLines(container: HTMLElement): () => void {
  const canvas = document.createElement('canvas')
  canvas.style.cssText = `
    position: absolute;
    top: 0; left: 0;
    width: 100%; height: 100%;
    pointer-events: none;
    z-index: 0;
  `
  container.appendChild(canvas)

  const ctx = canvas.getContext('2d')
  if (!ctx) return () => {}

  const ro = new ResizeObserver(() => {
    const dpr = window.devicePixelRatio || 1
    canvas.width = container.clientWidth * dpr
    canvas.height = container.clientHeight * dpr
    ctx.scale(dpr, dpr)
  })
  ro.observe(container)

  let mouseX = -1000
  let mouseY = -1000
  const onMouseMove = (e: MouseEvent) => {
    const rect = container.getBoundingClientRect()
    mouseX = e.clientX - rect.left
    mouseY = e.clientY - rect.top
  }
  window.addEventListener('mousemove', onMouseMove, { passive: true })

  let rafId: number
  const lines: Line[] = []
  const now = performance.now()
  for (let i = 0; i < LINE_COUNT; i++) {
    lines.push(createLine(now, i * 2000))
  }

  const drawFrame = (time: number) => {
    const w = container.clientWidth
    const h = container.clientHeight
    ctx.clearRect(0, 0, w, h)

    for (let i = 0; i < lines.length; i++) {
      updateAndDrawLine(
        ctx,
        lines[i],
        i,
        time,
        w,
        h,
        mouseX,
        mouseY,
        (newLine) => {
          lines[i] = newLine
        },
      )
    }

    rafId = requestAnimationFrame(drawFrame)
  }

  rafId = requestAnimationFrame(drawFrame)

  return () => {
    cancelAnimationFrame(rafId)
    window.removeEventListener('mousemove', onMouseMove)
    ro.disconnect()
    canvas.remove()
  }
}

function updateAndDrawLine(
  ctx: CanvasRenderingContext2D,
  line: Line,
  index: number,
  time: number,
  w: number,
  h: number,
  mx: number,
  my: number,
  onReset: (l: Line) => void,
) {
  const elapsed = time - line.startTime
  if (elapsed < 0) return

  const progress = elapsed / CROSS_DURATION_MS
  if (progress > 1) {
    onReset(createLine(time))
    return
  }

  const x = progress * w
  const baseY = h * (line.y0 + (line.y2 - line.y0) * progress)
  const drift = Math.sin(progress * Math.PI * 2 + index) * (h * 0.05)
  let y = baseY + drift

  // Repulsion
  const dx = x - mx
  const dy = y - my
  const dist = Math.sqrt(dx * dx + dy * dy)
  if (dist < REPULSION_RADIUS) {
    const force = (1 - dist / REPULSION_RADIUS) * REPULSION_STRENGTH
    y += (dy / dist) * force
  }

  const opacity = Math.sin(progress * Math.PI) * 0.15
  ctx.beginPath()
  ctx.moveTo(x - w * 0.1, y)
  ctx.lineTo(x, y)
  ctx.strokeStyle = `rgba(${COLOR_PRIMARY_RGB}, ${opacity})`
  ctx.lineWidth = 1
  ctx.stroke()
}
