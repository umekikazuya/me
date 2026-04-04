import { Router, Routes } from '@lit-labs/router'
import { css, html, LitElement } from 'lit'
import { customElement } from 'lit/decorators.js'
import type { RouteShellElement } from './route-shell.js'
import './app-admin-shell.js'
import './app-public-shell.js'
import '../pages/page-about.js'
import '../pages/page-articles.js'
import '../pages/page-not-found.js'
import '../pages/page-top.js'
import { setupCursor } from '../utils/cursor.js'
import { setupBackgroundShift } from '../utils/scroll.js'

@customElement('app-root')
export class AppRoot extends LitElement {
  private cleanups: Array<() => void> = []
  private router = new Router(this, [])
  private publicRoutes = new Routes(this, [
    { path: '/', render: () => html`<page-top></page-top>` },
    { path: '/articles', render: () => html`<page-articles></page-articles>` },
    { path: '/about', render: () => html`<page-about></page-about>` },
    { path: '/*', render: () => html`<page-not-found></page-not-found>` },
  ])
  private adminRoutes = new Routes(this, [
    { path: '/admin', render: () => html`<page-not-found></page-not-found>` },
    { path: '/admin/*', render: () => html`<page-not-found></page-not-found>` },
    { path: '/*', render: () => html`<page-not-found></page-not-found>` },
  ])

  render() {
    return this.isAdminPath(location.pathname)
      ? html`<app-admin-shell>${this.adminRoutes.outlet()}</app-admin-shell>`
      : html`<app-public-shell>${this.publicRoutes.outlet()}</app-public-shell>`
  }

  firstUpdated() {
    this.cleanups.push(setupBackgroundShift())
    this.cleanups.push(setupCursor())
    this.cleanups.push(this.setupNavigation())
  }

  private setupNavigation(): () => void {
    const onClick = async (e: Event) => {
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
        await this.navigate(anchor)
        return
      }

      const shell = this.shadowRoot?.querySelector(
        'app-public-shell, app-admin-shell',
      ) as RouteShellElement | null
      if (shell) {
        const ready = await shell.playLeaveTransition()
        if (!ready) return
      }

      await this.navigate(anchor)
    }

    this.shadowRoot?.addEventListener('click', onClick)
    return () => this.shadowRoot?.removeEventListener('click', onClick)
  }

  private isAdminPath(pathname: string) {
    return pathname === '/admin' || pathname.startsWith('/admin/')
  }

  private async navigate(anchor: HTMLAnchorElement) {
    if (anchor.href === location.href) return

    window.history.pushState({}, '', anchor.href)
    await this.router.goto(new URL(anchor.href).pathname)
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
  `
}

declare global {
  interface HTMLElementTagNameMap {
    'app-root': AppRoot
  }
}
