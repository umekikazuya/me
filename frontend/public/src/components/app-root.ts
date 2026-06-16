import { provide } from '@lit/context'
import { Router, Routes } from '@lit-labs/router'
import type { PropertyValues } from 'lit'
import { css, html, LitElement } from 'lit'
import { customElement, state } from 'lit/decorators.js'
import { articleContext } from '../contexts/article-context.js'
import { profileContext } from '../contexts/profile-context.js'
import { ArticleRepository } from '../domain/ArticleRepository.js'
import { ProfileRepository } from '../domain/ProfileRepository.js'
import { setupCursor } from '../utils/cursor.js'
import { setupBackgroundShift } from '../utils/scroll.js'
import { initGA, trackPageView } from '../utils/analytics.js'
import '../pages/page-about.js'
import '../pages/page-articles.js'
import '../pages/page-not-found.js'
import '../pages/page-top.js'
import './app-public-shell.js'
import type { RouteShellElement } from './route-shell.js'

@customElement('app-root')
export class AppRoot extends LitElement {
  @provide({ context: profileContext })
  profile = new ProfileRepository()

  @provide({ context: articleContext })
  article = new ArticleRepository()

  @state()
  private currentPath = window.location.pathname

  private cleanups: Array<() => void> = []
  private router = new Router(this, [])

  private onPopState = () => {
    this.currentPath = window.location.pathname
  }

  private publicRoutes = new Routes(this, [
    { path: '/', render: () => html`<page-top></page-top>` },
    { path: '/articles', render: () => html`<page-articles></page-articles>` },
    { path: '/about', render: () => html`<page-about></page-about>` },
    { path: '/*', render: () => html`<page-not-found></page-not-found>` },
  ])

  render() {
    return html`<app-public-shell>${this.publicRoutes.outlet()}</app-public-shell>`
  }

  connectedCallback() {
    super.connectedCallback()
    window.addEventListener('popstate', this.onPopState)
    this.updateVisualEffects()
  }

  disconnectedCallback() {
    super.disconnectedCallback()
    window.removeEventListener('popstate', this.onPopState)
    this.teardownVisualEffects()
  }

  protected updated(changedProperties: PropertyValues) {
    if (changedProperties.has('currentPath')) {
      this.updateVisualEffects()
      trackPageView(this.currentPath)
    }
  }

  private updateVisualEffects() {
    if (typeof window === 'undefined') return

    document.documentElement.setAttribute('data-theme', 'public')

    this.teardownVisualEffects()

    this.cleanups.push(setupBackgroundShift())
    this.cleanups.push(setupCursor())
    this.cleanups.push(this.setupNavigation())
  }

  private teardownVisualEffects() {
    for (const cleanup of this.cleanups) cleanup()
    this.cleanups = []
  }

  firstUpdated() {
    void this.profile.loadProfile()
    initGA()
    trackPageView(this.currentPath)
  }

  private setupNavigation(): () => void {
    const onClick = async (e: Event) => {
      if (this.shouldPreventNavigation(e)) return
      const anchor = this.findAnchor(e)
      if (!anchor) return

      e.preventDefault()
      if (!this.isReducedMotion()) {
        const ready = await this.playTransition()
        if (!ready) return
      }
      await this.navigate(anchor)
    }
    this.shadowRoot?.addEventListener('click', onClick)
    return () => this.shadowRoot?.removeEventListener('click', onClick)
  }

  private shouldPreventNavigation(e: Event) {
    return (
      e.defaultPrevented ||
      (e instanceof MouseEvent &&
        (e.button !== 0 || e.metaKey || e.ctrlKey || e.shiftKey || e.altKey))
    )
  }

  private findAnchor(e: Event) {
    const anchor = (e.composedPath() as Element[]).find(
      (el) => (el as HTMLElement).tagName === 'A',
    ) as HTMLAnchorElement | undefined
    if (
      !anchor?.href ||
      (anchor.target && anchor.target !== '_self') ||
      anchor.hasAttribute('download')
    )
      return null
    const url = new URL(anchor.href)
    return url.origin === location.origin ? anchor : null
  }

  private isReducedMotion() {
    return window.matchMedia('(prefers-reduced-motion: reduce)').matches
  }

  private async playTransition() {
    const shell = this.shadowRoot?.querySelector(
      'app-public-shell',
    ) as RouteShellElement | null
    return shell ? await shell.playLeaveTransition() : true
  }

  private async navigate(anchor: HTMLAnchorElement) {
    if (anchor.href === location.href) return
    await this.navigateToPath(new URL(anchor.href).pathname)
  }

  private async navigateToPath(pathname: string, replace = false) {
    if (pathname === this.currentPath) return

    if (replace) window.history.replaceState({}, '', pathname)
    else window.history.pushState({}, '', pathname)

    this.currentPath = pathname
    await this.router.goto(pathname)
  }

  static styles = css`
    :host { display: block; }
  `
}

declare global {
  interface HTMLElementTagNameMap {
    'app-root': AppRoot
  }
}
