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
    const onLeaveEnd = (event: TransitionEvent) => {
      if (event.target !== outlet) return
      outlet.removeEventListener('transitionend', onLeaveEnd)
      outlet.classList.remove('leaving')
      resolve(true)
    }

    outlet.addEventListener('transitionend', onLeaveEnd)
  })
}
