import { css } from 'lit'

export interface RouteShellElement extends HTMLElement {
  playLeaveTransition(): Promise<boolean>
}

export const routeShellStyles = css`
  #outlet {
    opacity: 1;
    transform: translateY(0);
  }

  #outlet.leaving {
    opacity: 0;
    transform: translateY(-10px);
    transition:
      opacity 0.3s var(--easing-smooth),
      transform 0.3s var(--easing-smooth);
  }
`

export function playLeaveTransition(outlet: HTMLElement | null) {
  if (!outlet) return Promise.resolve(true)
  if (outlet.classList.contains('leaving')) return Promise.resolve(false)

  outlet.classList.add('leaving')

  return new Promise<boolean>((resolve) => {
    let timeoutId: number | undefined
    let settled = false

    const cleanup = () => {
      if (timeoutId !== undefined) window.clearTimeout(timeoutId)
      outlet.removeEventListener('transitionend', onLeaveEnd)
      outlet.removeEventListener('transitioncancel', onLeaveCancel)
      outlet.classList.remove('leaving')
    }

    const finish = () => {
      if (settled) return
      settled = true
      cleanup()
      resolve(true)
    }

    const onLeaveEnd = (event: TransitionEvent) => {
      if (event.target !== outlet) return
      finish()
    }

    const onLeaveCancel = (event: TransitionEvent) => {
      if (event.target !== outlet) return
      finish()
    }

    outlet.addEventListener('transitionend', onLeaveEnd)
    outlet.addEventListener('transitioncancel', onLeaveCancel)
    timeoutId = window.setTimeout(finish, 500)
  })
}
