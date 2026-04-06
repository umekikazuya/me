import { Router, Routes } from '@lit-labs/router'
import type { PropertyValues } from 'lit'
import { css, html, LitElement, nothing } from 'lit'
import { customElement, state } from 'lit/decorators.js'
import {
  changeEmail,
  login,
  logout,
  refreshSession,
  revokeAllSessions,
} from '../admin/auth-api.js'
import { getMe, updateMe } from '../admin/me-api.js'
import {
  type AdminLoginInput,
  ApiError,
  type ChangeEmailInput,
  createEmptyMeProfile,
  describeApiError,
  type MeProfile,
} from '../admin/types.js'
import type { RouteShellElement } from './route-shell.js'
import './app-admin-shell.js'
import './app-public-shell.js'
import '../pages/page-admin-account.js'
import '../pages/page-admin-dashboard.js'
import '../pages/page-admin-login.js'
import '../pages/page-admin-profile.js'
import '../pages/page-about.js'
import '../pages/page-articles.js'
import '../pages/page-not-found.js'
import '../pages/page-top.js'
import { setupCursor } from '../utils/cursor.js'
import { setupBackgroundShift } from '../utils/scroll.js'

@customElement('app-root')
export class AppRoot extends LitElement {
  @state()
  private currentPath = window.location.pathname

  @state()
  private adminSessionStatus:
    | 'unknown'
    | 'checking'
    | 'authenticated'
    | 'guest' = 'unknown'

  @state()
  private adminLoginPending = false

  @state()
  private adminLoginError = ''

  @state()
  private adminProfile = createEmptyMeProfile()

  @state()
  private adminProfileLoading = false

  @state()
  private adminProfileSaving = false

  @state()
  private adminProfileLoaded = false

  @state()
  private adminProfileError = ''

  @state()
  private adminProfileSuccess = ''

  @state()
  private adminAccountBusyAction = ''

  @state()
  private adminAccountError = ''

  @state()
  private adminAccountSuccess = ''

  private cleanups: Array<() => void> = []
  private router = new Router(this, [])
  private adminReturnPath = '/admin'
  private adminSessionBootstrap?: Promise<void>
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
          html`<page-admin-profile
            .profile=${this.adminProfile}
            .loading=${this.adminProfileLoading}
            .saving=${this.adminProfileSaving}
            .errorMessage=${this.adminProfileError}
            .successMessage=${this.adminProfileSuccess}
            @admin-save-profile=${this.handleAdminProfileSave}
          ></page-admin-profile>`,
        ),
    },
    {
      path: '/admin/account',
      render: () =>
        this.renderProtectedAdmin(
          html`<page-admin-account
            .busyAction=${this.adminAccountBusyAction}
            .errorMessage=${this.adminAccountError}
            .successMessage=${this.adminAccountSuccess}
            @admin-logout=${this.handleAdminLogout}
            @admin-revoke-sessions=${this.handleAdminRevokeSessions}
            @admin-change-email=${this.handleAdminChangeEmail}
          ></page-admin-account>`,
        ),
    },
    { path: '/*', render: () => html`<page-not-found></page-not-found>` },
  ])

  render() {
    return this.isAdminPath(this.currentPath)
      ? html`<app-admin-shell
          .authenticated=${this.adminSessionStatus === 'authenticated'}
          .busy=${this.adminSessionStatus === 'checking'}
          .currentPath=${this.currentPath}
          >${this.adminRoutes.outlet()}</app-admin-shell
        >`
      : html`<app-public-shell>${this.publicRoutes.outlet()}</app-public-shell>`
  }

  connectedCallback() {
    super.connectedCallback()
    window.addEventListener('popstate', this.onPopState)
  }

  firstUpdated() {
    this.cleanups.push(setupBackgroundShift())
    this.cleanups.push(setupCursor())
    this.cleanups.push(this.setupNavigation())
    if (this.isAdminPath(this.currentPath)) {
      void this.syncAdminRouteState()
    }
  }

  protected updated(changedProperties: PropertyValues) {
    if (
      changedProperties.has('currentPath') &&
      this.isAdminPath(this.currentPath)
    ) {
      void this.syncAdminRouteState()
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

  private async navigateToPath(pathname: string, replace = false) {
    if (pathname === this.currentPath) return

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
      if (this.adminSessionStatus === 'unknown') {
        await this.bootstrapAdminSession()
      }
      if (this.adminSessionStatus === 'authenticated') {
        await this.navigateToPath(this.adminReturnPath, true)
      }
      return
    }

    if (!this.isProtectedAdminPath(this.currentPath)) return

    if (this.adminSessionStatus !== 'authenticated') {
      this.adminReturnPath = this.currentPath
      await this.bootstrapAdminSession()
    }

    if (this.adminSessionStatus !== 'authenticated') {
      await this.navigateToPath('/admin/login', true)
      return
    }

    if (this.currentPath === '/admin/profile' && !this.adminProfileLoaded) {
      await this.loadAdminProfile()
    }
  }

  private async bootstrapAdminSession() {
    if (this.adminSessionBootstrap) {
      await this.adminSessionBootstrap
      return
    }

    this.adminSessionStatus = 'checking'
    this.adminSessionBootstrap = (async () => {
      try {
        await refreshSession()
        this.adminSessionStatus = 'authenticated'
        this.adminLoginError = ''
      } catch (error) {
        const isUnauthorized = error instanceof ApiError && error.status === 401
        this.adminSessionStatus = 'guest'
        if (this.currentPath === '/admin/login' || isUnauthorized) {
          this.adminLoginError = ''
        } else {
          this.adminLoginError = describeApiError(error)
        }
      } finally {
        this.adminSessionBootstrap = undefined
      }
    })()

    await this.adminSessionBootstrap
  }

  private renderAdminLogin() {
    if (this.adminSessionStatus === 'checking') {
      return this.renderAdminStatus('セッションを確認しています...')
    }

    return html`<page-admin-login
      .submitting=${this.adminLoginPending}
      .errorMessage=${this.adminLoginError}
      @admin-login-submit=${this.handleAdminLogin}
    ></page-admin-login>`
  }

  private renderProtectedAdmin(content: unknown) {
    if (
      this.adminSessionStatus === 'checking' ||
      this.adminSessionStatus === 'unknown'
    ) {
      return this.renderAdminStatus('認証状態を確認しています...')
    }

    if (this.adminSessionStatus !== 'authenticated') {
      return nothing
    }

    return content
  }

  private renderAdminStatus(message: string) {
    return html`<section class="admin-status">
      <p>${message}</p>
    </section>`
  }

  private async handleAdminLogin(event: CustomEvent<AdminLoginInput>) {
    this.adminLoginPending = true
    this.adminLoginError = ''
    try {
      await login(event.detail)
      this.adminSessionStatus = 'authenticated'
      this.adminProfile = createEmptyMeProfile()
      this.adminProfileLoaded = false
      this.adminProfileError = ''
      this.adminProfileSuccess = ''
      await this.navigateToPath(this.adminReturnPath)
    } catch (error) {
      this.adminLoginError = describeApiError(error)
    } finally {
      this.adminLoginPending = false
    }
  }

  private async loadAdminProfile() {
    this.adminProfileLoading = true
    this.adminProfileError = ''
    try {
      this.adminProfile = await getMe()
      this.adminProfileLoaded = true
    } catch (error) {
      this.adminProfileError = describeApiError(error)
    } finally {
      this.adminProfileLoading = false
    }
  }

  private async handleAdminProfileSave(event: CustomEvent<MeProfile>) {
    this.adminProfileSaving = true
    this.adminProfileError = ''
    this.adminProfileSuccess = ''

    try {
      this.adminProfile = await updateMe(event.detail)
      this.adminProfileLoaded = true
      this.adminProfileSuccess = 'プロフィールを更新しました。'
    } catch (error) {
      this.adminProfileError = describeApiError(error)
    } finally {
      this.adminProfileSaving = false
    }
  }

  private async handleAdminLogout() {
    this.adminAccountBusyAction = 'logout'
    this.adminAccountError = ''
    this.adminAccountSuccess = ''
    try {
      await logout()
      this.adminSessionStatus = 'guest'
      this.adminProfileLoaded = false
      this.adminAccountSuccess = 'ログアウトしました。'
      await this.navigateToPath('/admin/login')
    } catch (error) {
      this.adminAccountError = describeApiError(error)
    } finally {
      this.adminAccountBusyAction = ''
    }
  }

  private async handleAdminRevokeSessions() {
    this.adminAccountBusyAction = 'revoke-sessions'
    this.adminAccountError = ''
    this.adminAccountSuccess = ''
    try {
      await revokeAllSessions()
      this.adminSessionStatus = 'guest'
      this.adminProfileLoaded = false
      this.adminAccountSuccess = '全セッションを失効させました。'
      await this.navigateToPath('/admin/login')
    } catch (error) {
      this.adminAccountError = describeApiError(error)
    } finally {
      this.adminAccountBusyAction = ''
    }
  }

  private async handleAdminChangeEmail(event: CustomEvent<ChangeEmailInput>) {
    this.adminAccountBusyAction = 'change-email'
    this.adminAccountError = ''
    this.adminAccountSuccess = ''
    try {
      await changeEmail(event.detail)
      this.adminAccountSuccess = 'メールアドレス変更を送信しました。'
    } catch (error) {
      this.adminAccountError = describeApiError(error)
    } finally {
      this.adminAccountBusyAction = ''
    }
  }

  disconnectedCallback() {
    super.disconnectedCallback()
    window.removeEventListener('popstate', this.onPopState)
    for (const cleanup of this.cleanups) cleanup()
    this.cleanups = []
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
