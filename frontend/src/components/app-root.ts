import { provide } from '@lit/context'
import { Router, Routes } from '@lit-labs/router'
import { SignalWatcher } from '@lit-labs/signals'
import type { PropertyValues } from 'lit'
import { css, html, LitElement } from 'lit'
import { customElement, state } from 'lit/decorators.js'
import { articleContext } from '../contexts/article-context.js'
import { authContext } from '../contexts/auth-context.js'
import { profileContext } from '../contexts/profile-context.js'
import { RepositoryObserver } from '../controllers/RepositoryObserver.js'
import { ArticleRepository } from '../domain/ArticleRepository.js'
import { AuthRepository } from '../domain/AuthRepository.js'
import { ProfileRepository } from '../domain/ProfileRepository.js'
import { setupCursor } from '../utils/cursor.js'
import { setupBackgroundShift } from '../utils/scroll.js'
import './app-admin-shell.js'
import './app-public-shell.js'
import '../components/admin/ui/me-admin-auth-boundary.js'
import '../pages/page-admin-account.js'
import '../pages/page-admin-articles.js'
import '../pages/page-admin-profile.js'
import '../pages/page-about.js'
import '../pages/page-articles.js'
import '../pages/page-not-found.js'
import '../pages/page-admin-entry.js'
import '../pages/page-top.js'
import type { RouteShellElement } from './route-shell.js'

@customElement('app-root')
export class AppRoot extends SignalWatcher(LitElement) {
  @provide({ context: authContext })
  auth = new AuthRepository()

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

  private adminRoutes = new Routes(this, [
    {
      path: '/admin',
      render: () => html`
        <page-admin-entry></page-admin-entry>
      `,
    },
    {
      path: '/admin/profile',
      render: () => html`
        <me-admin-auth-boundary>
          <page-admin-profile></page-admin-profile>
        </me-admin-auth-boundary>
      `,
    },
    {
      path: '/admin/articles',
      render: () => html`
        <me-admin-auth-boundary>
          <page-admin-articles></page-admin-articles>
        </me-admin-auth-boundary>
      `,
    },
    {
      path: '/admin/account',
      render: () => html`
        <me-admin-auth-boundary>
          <page-admin-account></page-admin-account>
        </me-admin-auth-boundary>
      `,
    },
    { path: '/*', render: () => html`<page-not-found></page-not-found>` },
  ])

  constructor() {
    super()
    new RepositoryObserver(this, this.auth)
  }

  render() {
    const isAdmin = this.isAdminPath(this.currentPath)
    const status = this.auth.status.value // Pure Signal consumption

    return isAdmin
      ? html`<app-admin-shell
          .authenticated=${status === 'authenticated'}
          .currentPath=${this.currentPath}
          >${this.adminRoutes.outlet()}</app-admin-shell
        >`
      : html`<app-public-shell>${this.publicRoutes.outlet()}</app-public-shell>`
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
    }
  }

  private updateVisualEffects() {
    if (typeof window === 'undefined') return

    const theme = this.isAdminPath(this.currentPath) ? 'admin' : 'public'
    document.documentElement.setAttribute('data-theme', theme)

    this.teardownVisualEffects()

    if (!this.isAdminPath(this.currentPath)) {
      this.cleanups.push(setupBackgroundShift())
      this.cleanups.push(setupCursor())
    }
    this.cleanups.push(this.setupNavigation())
  }

  private teardownVisualEffects() {
    for (const cleanup of this.cleanups) cleanup()
    this.cleanups = []
  }

  firstUpdated() {
    void this.profile.loadPublicProfile()
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
      'app-public-shell, app-admin-shell',
    ) as RouteShellElement | null
    return shell ? await shell.playLeaveTransition() : true
  }

  private isAdminPath(pathname: string) {
    return pathname === '/admin' || pathname.startsWith('/admin/')
  }

  private async navigate(anchor: HTMLAnchorElement) {
    if (anchor.href === location.href) return
    await this.navigateToPath(new URL(anchor.href).pathname)
  }

  private async navigateToPath(
    pathname: string,
    replace = false,
    force = false,
  ) {
    if (pathname === this.currentPath) return
    if (
      !force &&
      this.shouldConfirmAdminNavigation(pathname) &&
      !window.confirm('未保存の変更があります。ページを移動してもよいですか？')
    )
      return

    if (replace) window.history.replaceState({}, '', pathname)
    else window.history.pushState({}, '', pathname)

    this.currentPath = pathname
    await this.router.goto(pathname)
  }

  private shouldConfirmAdminNavigation(pathname: string) {
    const isProfile = this.currentPath === '/admin/profile'
    const isArticles = this.currentPath === '/admin/articles'
    return (
      pathname !== this.currentPath &&
      ((this.profile.adminDirty && isProfile) ||
        (this.article.adminDirty && isArticles))
    )
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
