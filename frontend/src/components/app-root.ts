import { provide } from '@lit/context'
import { Router, Routes } from '@lit-labs/router'
import type { PropertyValues } from 'lit'
import { css, html, LitElement, nothing } from 'lit'
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
import '../pages/page-admin-account.js'
import '../pages/page-admin-articles.js'
import '../pages/page-admin-dashboard.js'
import '../pages/page-admin-login.js'
import '../pages/page-admin-profile.js'
import '../pages/page-about.js'
import '../pages/page-articles.js'
import '../pages/page-not-found.js'
import '../pages/page-top.js'
import type { RouteShellElement } from './route-shell.js'

@customElement('app-root')
export class AppRoot extends LitElement {
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
  private adminReturnPath = '/admin'

  constructor() {
    super()
    new RepositoryObserver(this, this.auth)
    new RepositoryObserver(this, this.profile)
    new RepositoryObserver(this, this.article)
  }

  private onPopState = () => {
    this.currentPath = window.location.pathname
  }

  private publicRoutes = new Routes(this, [
    {
      path: '/',
      render: () => html`<page-top></page-top>`,
    },
    { path: '/articles', render: () => html`<page-articles></page-articles>` },
    {
      path: '/about',
      render: () => html`<page-about></page-about>`,
    },
    { path: '/*', render: () => html`<page-not-found></page-not-found>` },
  ])

  private adminRoutes = new Routes(this, [
    {
      path: '/admin/login',
      render: () => this.renderAdminLogin(),
    },
    {
      path: '/admin',
      render: () =>
        this.renderProtectedAdmin(
          html`<page-admin-dashboard></page-admin-dashboard>`,
        ),
    },
    {
      path: '/admin/profile',
      render: () =>
        this.renderProtectedAdmin(
          html`<page-admin-profile></page-admin-profile>`,
        ),
    },
    {
      path: '/admin/articles',
      render: () =>
        this.renderProtectedAdmin(
          html`<page-admin-articles></page-admin-articles>`,
        ),
    },
    {
      path: '/admin/account',
      render: () =>
        this.renderProtectedAdmin(
          html`<page-admin-account></page-admin-account>`,
        ),
    },
    { path: '/*', render: () => html`<page-not-found></page-not-found>` },
  ])

  render() {
    const isAdmin = this.isAdminPath(this.currentPath)
    return isAdmin
      ? html`<app-admin-shell
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

  private updateVisualEffects() {
    // Theme switching
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
    for (const cleanup of this.cleanups) {
      cleanup()
    }
    this.cleanups = []
  }

  firstUpdated() {
    void this.profile.loadPublicProfile()
    if (this.isAdminPath(this.currentPath)) {
      void this.syncAdminRouteState()
    }
  }

  protected updated(changedProperties: PropertyValues) {
    if (changedProperties.has('currentPath')) {
      this.updateVisualEffects()
      if (this.isAdminPath(this.currentPath)) {
        void this.syncAdminRouteState()
      }
    }
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
    ) {
      return
    }

    if (replace) {
      window.history.replaceState({}, '', pathname)
    } else {
      window.history.pushState({}, '', pathname)
    }

    this.currentPath = pathname
    await this.router.goto(pathname)
  }

  private isProtectedAdminPath(pathname: string) {
    return this.isAdminPath(pathname) && pathname !== '/admin/login'
  }

  private async syncAdminRouteState() {
    if (!this.isAdminPath(this.currentPath)) return

    if (this.currentPath === '/admin/login') {
      if (this.auth.status === 'unknown') {
        await this.auth.refreshSession()
      }
      if (this.auth.status === 'authenticated') {
        await this.navigateToPath(this.adminReturnPath, true, true)
      }
      return
    }

    if (!this.isProtectedAdminPath(this.currentPath)) return

    if (this.auth.status !== 'authenticated') {
      this.adminReturnPath = this.currentPath
      await this.auth.refreshSession()
    }

    if (this.auth.status !== 'authenticated') {
      await this.navigateToPath('/admin/login', true, true)
      return
    }

    if (this.currentPath === '/admin/profile' && !this.profile.adminLoaded) {
      await this.profile.loadAdminProfile()
    }
  }

  private renderAdminLogin() {
    if (this.auth.status === 'checking') {
      return this.renderAdminStatus('セッションを確認しています...')
    }

    return html`<page-admin-login></page-admin-login>`
  }

  private renderProtectedAdmin(content: unknown) {
    if (this.auth.status === 'checking' || this.auth.status === 'unknown') {
      return this.renderAdminStatus('認証状態を確認しています...')
    }

    if (this.auth.status !== 'authenticated') {
      return nothing
    }

    return content
  }

  private renderAdminStatus(message: string) {
    return html`<section class="admin-status">
      <p>${message}</p>
    </section>`
  }

  private shouldConfirmAdminNavigation(pathname: string) {
    return (
      pathname !== this.currentPath &&
      ((this.profile.adminDirty && this.currentPath === '/admin/profile') ||
        (this.article.adminDirty && this.currentPath === '/admin/articles'))
    )
  }

  static styles = css`
    :host {
      display: block;
    }

    .admin-status {
      min-height: 60dvh;
      display: grid;
      place-items: center;
      color: var(--color-text-secondary);
      font-size: 15px;
      letter-spacing: var(--tracking-wide);
    }
  `
}

declare global {
  interface HTMLElementTagNameMap {
    'app-root': AppRoot
  }
}
