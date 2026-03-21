import { LitElement, css, html } from 'lit'
import { customElement } from 'lit/decorators.js'
import { Router } from '@vaadin/router'
import { setupBackgroundShift } from '../utils/scroll.js'
import { setupCursor } from '../utils/cursor.js'
import './nav-bar.js'

@customElement('app-root')
export class AppRoot extends LitElement {
  private cleanups: Array<() => void> = []

  render() {
    return html`
      <nav-bar></nav-bar>
      <main id="outlet"></main>
    `
  }

  firstUpdated() {
    const outlet = this.shadowRoot?.querySelector('#outlet') as HTMLElement
    const router = new Router(outlet)
    router.setRoutes([
      { path: '/', component: 'page-top' },
      { path: '/articles', component: 'page-articles' },
      { path: '/about', component: 'page-about' },
      { path: '(.*)', component: 'page-not-found' },
    ])

    this.cleanups.push(setupBackgroundShift())
    this.cleanups.push(setupCursor())
    this.cleanups.push(this.setupNavigation(outlet))
  }

  private setupNavigation(outlet: HTMLElement): () => void {
    const onClick = (e: Event) => {
      if (
        e.defaultPrevented ||
        (e instanceof MouseEvent &&
          (e.button !== 0 || e.metaKey || e.ctrlKey || e.shiftKey || e.altKey))
      ) {
        return
      }
      const anchor = (e.composedPath() as Element[]).find(
        (el) => (el as HTMLElement).tagName === 'A',
      ) as HTMLAnchorElement | undefined

      if (
        !anchor?.href ||
        (anchor.target && anchor.target !== '_self') ||
        anchor.hasAttribute('download')
      )
        return
      const url = new URL(anchor.href)
      if (url.origin !== location.origin) return

      e.preventDefault()

      const reduced = window.matchMedia(
        '(prefers-reduced-motion: reduce)',
      ).matches
      if (reduced) {
        Router.go(anchor.href)
        return
      }

      if (outlet.classList.contains('leaving')) return
      outlet.classList.add('leaving')
      const onLeaveEnd = (event: TransitionEvent) => {
        if (event.target !== outlet) return
        outlet.removeEventListener('transitionend', onLeaveEnd)
        outlet.classList.remove('leaving')
        Router.go(anchor.href)
      }
      outlet.addEventListener('transitionend', onLeaveEnd)
    }

    this.shadowRoot?.addEventListener('click', onClick)
    return () => this.shadowRoot?.removeEventListener('click', onClick)
  }

  disconnectedCallback() {
    super.disconnectedCallback()
    for (const cleanup of this.cleanups) cleanup()
    this.cleanups = []
  }

  static styles = css`
    :host {
      display: block;
    }

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
}

declare global {
  interface HTMLElementTagNameMap {
    'app-root': AppRoot
  }
}
