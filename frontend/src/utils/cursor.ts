const TRAIL_LIFETIME = 2200 // ms
const MAX_POINTS = 120
const COLOR_PRIMARY = '#d1cdc7'
const COLOR_PRIMARY_RGB = '209, 205, 199'

type Point = { x: number; y: number; t: number }
type DotState = 'click' | 'hover' | 'idle'

const DOT_STYLES: Record<
  DotState,
  { size: string; margin: string; opacity: string; duration: string }
> = {
  click: { size: '4px', margin: '1px', opacity: '0.8', duration: '0.1s' },
  hover: { size: '32px', margin: '-13px', opacity: '0.2', duration: '0.4s' },
  idle: { size: '6px', margin: '0px', opacity: '0.6', duration: '0.3s' },
}

function isTouchDevice(): boolean {
  return window.matchMedia('(pointer: coarse)').matches
}

/** Catmull-Rom spline: draw the segment between p1 and p2 */
function catmullRomSegment(
  ctx: CanvasRenderingContext2D,
  p0: Point,
  p1: Point,
  p2: Point,
  p3: Point,
) {
  // Convert catmull-rom to cubic bezier control points
  const cp1x = p1.x + (p2.x - p0.x) / 6
  const cp1y = p1.y + (p2.y - p0.y) / 6
  const cp2x = p2.x - (p3.x - p1.x) / 6
  const cp2y = p2.y - (p3.y - p1.y) / 6
  ctx.bezierCurveTo(cp1x, cp1y, cp2x, cp2y, p2.x, p2.y)
}

function drawTrail(
  ctx: CanvasRenderingContext2D,
  canvas: HTMLCanvasElement,
  points: Point[],
  now: number,
) {
  ctx.clearRect(0, 0, canvas.width, canvas.height)
  if (points.length < 2) return

  // Draw segments with opacity based on age
  for (let i = 1; i < points.length; i++) {
    const age = now - points[i].t
    const opacity = Math.max(0, 1 - age / TRAIL_LIFETIME) * 0.08
    if (opacity <= 0) continue

    const p0 = points[Math.max(0, i - 2)]
    const p1 = points[i - 1]
    const p2 = points[i]
    const p3 = points[Math.min(points.length - 1, i + 1)]

    ctx.beginPath()
    ctx.moveTo(p1.x, p1.y)
    catmullRomSegment(ctx, p0, p1, p2, p3)
    ctx.strokeStyle = `rgba(${COLOR_PRIMARY_RGB}, ${opacity})`
    ctx.lineWidth = 1
    ctx.lineCap = 'round'
    ctx.stroke()
  }
}

function createTouchRipple(x: number, y: number) {
  const ripple = document.createElement('div')
  ripple.style.cssText = `
    position: fixed;
    left: ${x}px;
    top: ${y}px;
    width: 0;
    height: 0;
    border-radius: 50%;
    background: rgba(${COLOR_PRIMARY_RGB}, 0.05);
    transform: translate(-50%, -50%);
    pointer-events: none;
    z-index: 9999;
    animation: ripple-expand 0.6s ease-out forwards;
  `
  document.body.appendChild(ripple)
  ripple.addEventListener('animationend', () => ripple.remove())
}

export function setupCursor(): () => void {
  const reduced = window.matchMedia('(prefers-reduced-motion: reduce)').matches
  if (reduced) return () => {}
  const isTouch = isTouchDevice()

  // Touch: only ripple, no cursor/trail
  if (isTouch) {
    const onTouch = (e: TouchEvent) => {
      for (const touch of Array.from(e.changedTouches)) {
        createTouchRipple(touch.clientX, touch.clientY)
      }
    }
    window.addEventListener('touchstart', onTouch, { passive: true })
    return () => window.removeEventListener('touchstart', onTouch)
  }

  // Spotlight overlay
  const spotlight = document.createElement('div')
  spotlight.style.cssText = `
    position: fixed;
    top: 0; left: 0;
    width: 600px; height: 600px;
    border-radius: 50%;
    background: radial-gradient(circle, rgba(${COLOR_PRIMARY_RGB}, 0.03) 0%, rgba(${COLOR_PRIMARY_RGB}, 0) 70%);
    pointer-events: none;
    z-index: 1;
    will-change: transform;
    transform: translate(-50%, -50%);
  `

  // Cursor dot
  const dot = document.createElement('div')
  dot.style.cssText = `
    position: fixed;
    top: 0; left: 0;
    width: 6px; height: 6px;
    border-radius: 50%;
    background: ${COLOR_PRIMARY};
    opacity: 0.6;
    pointer-events: none;
    z-index: 9999;
    mix-blend-mode: screen;
    will-change: transform;
    transition: width 0.3s ease-out, height 0.3s ease-out, opacity 0.3s ease-out, margin 0.3s ease-out;
  `

  // Trail canvas
  const canvas = document.createElement('canvas')
  canvas.style.cssText = `
    position: fixed;
    top: 0; left: 0;
    width: 100vw; height: 100vh;
    pointer-events: none;
    z-index: 9998;
    will-change: transform;
  `
  const ctx = canvas.getContext('2d')
  if (!ctx) return () => {}

  document.body.style.cursor = 'none'
  document.body.appendChild(canvas)
  document.body.appendChild(spotlight)
  document.body.appendChild(dot)

  const resizeCanvas = () => {
    canvas.width = window.innerWidth
    canvas.height = window.innerHeight
  }
  resizeCanvas()
  window.addEventListener('resize', resizeCanvas, { passive: true })

  let targetX = -100
  let targetY = -100
  let currentX = -100
  let currentY = -100
  let spotlightX = -100
  let spotlightY = -100
  let hovering = false
  let clicking = false
  let inViewport = false
  const points: Point[] = []
  let rafId: number
  let lastX = -100
  let lastY = -100

  const setDotSize = () => {
    const state: DotState = clicking ? 'click' : hovering ? 'hover' : 'idle'
    const s = DOT_STYLES[state]
    dot.style.width = s.size
    dot.style.height = s.size
    dot.style.margin = s.margin
    dot.style.opacity = state === 'idle' && !inViewport ? '0' : s.opacity
    dot.style.transitionDuration =
      state === 'idle' && !inViewport ? '0.5s' : s.duration
  }

  const tick = () => {
    currentX = targetX
    currentY = targetY

    // Smooth spotlight follow
    spotlightX += (targetX - spotlightX) * 0.1
    spotlightY += (targetY - spotlightY) * 0.1

    dot.style.transform = `translate(${currentX - 3}px, ${currentY - 3}px)`
    if (inViewport) {
      spotlight.style.opacity = '1'
      spotlight.style.transform = `translate(${spotlightX - 300}px, ${spotlightY - 300}px)`
    } else {
      spotlight.style.opacity = '0'
    }

    // Add point if moved enough
    const dx = currentX - lastX
    const dy = currentY - lastY
    if (dx * dx + dy * dy > 4) {
      points.push({ x: currentX, y: currentY, t: performance.now() })
      lastX = currentX
      lastY = currentY
    }

    // Expire old points
    const now = performance.now()
    while (points.length > 0 && now - points[0].t > TRAIL_LIFETIME) {
      points.shift()
    }
    if (points.length > MAX_POINTS) points.splice(0, points.length - MAX_POINTS)

    drawTrail(ctx, canvas, points, now)

    rafId = requestAnimationFrame(tick)
  }

  const onMouseMove = (e: MouseEvent) => {
    targetX = e.clientX
    targetY = e.clientY
    inViewport = true
    setDotSize()
  }

  const onMouseLeave = () => {
    inViewport = false
    setDotSize()
  }

  const onMouseEnter = () => {
    inViewport = true
    setDotSize()
  }

  const onMouseOver = (e: MouseEvent) => {
    const target = e.target as HTMLElement
    const isLink =
      target.tagName === 'A' ||
      target.tagName === 'BUTTON' ||
      !!target.closest('a') ||
      !!target.closest('button')
    if (hovering !== isLink) {
      hovering = isLink
      setDotSize()
    }
  }

  const onMouseDown = () => {
    clicking = true
    setDotSize()
  }

  const onMouseUp = () => {
    clicking = false
    setDotSize()
  }

  window.addEventListener('mousemove', onMouseMove, { passive: true })
  document.addEventListener('mouseleave', onMouseLeave)
  document.addEventListener('mouseenter', onMouseEnter)
  window.addEventListener('mouseover', onMouseOver, { passive: true })
  window.addEventListener('mousedown', onMouseDown, { passive: true })
  window.addEventListener('mouseup', onMouseUp, { passive: true })

  tick()

  return () => {
    cancelAnimationFrame(rafId)
    window.removeEventListener('mousemove', onMouseMove)
    document.removeEventListener('mouseleave', onMouseLeave)
    document.removeEventListener('mouseenter', onMouseEnter)
    window.removeEventListener('mouseover', onMouseOver)
    window.removeEventListener('mousedown', onMouseDown)
    window.removeEventListener('mouseup', onMouseUp)
    window.removeEventListener('resize', resizeCanvas)
    dot.remove()
    canvas.remove()
    spotlight.remove()
    document.body.style.cursor = ''
  }
}
