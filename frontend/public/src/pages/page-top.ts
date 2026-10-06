import { consume } from '@lit/context'
import { css, html, nothing } from 'lit'
import { customElement } from 'lit/decorators.js'
import { profileContext } from '../contexts/profile-context.js'
import type { IProfileRepository } from '../domain/ProfileRepository.js'
import { pageStyles } from '../styles/page-styles.js'
import { sanitizeUrl } from '../utils/format.js'
import { ShellPage } from './shell-page.js'
import '../components/load-spinner.js'

/** `glow README.md` — プロフィールを Markdown として描画した体裁 */
@customElement('page-top')
export class PageTop extends ShellPage {
  @consume({ context: profileContext, subscribe: true })
  set profileRepo(repo: IProfileRepository) {
    if (this._profileRepo) {
      this._profileRepo.removeEventListener('change', this._onRepoChange)
    }
    this._profileRepo = repo
    repo.addEventListener('change', this._onRepoChange)
    this.requestUpdate()
  }
  get profileRepo() {
    return this._profileRepo
  }
  private _profileRepo!: IProfileRepository
  private _onRepoChange = () => this.requestUpdate()

  disconnectedCallback() {
    super.disconnectedCallback()
    if (this._profileRepo) {
      this._profileRepo.removeEventListener('change', this._onRepoChange)
    }
  }

  render() {
    return html`
      <shell-command command="glow README.md"></shell-command>
      ${this.typed ? this.renderOutput() : nothing}
    `
  }

  private renderOutput() {
    const p = this.profileRepo.profile
    if (!p) {
      return this.profileRepo.error
        ? html`<p class="ln m">glow: README.md: profile unavailable</p>
            <shell-nav current="home"></shell-nav>`
        : html`<load-spinner class="ln"></load-spinner>
            <shell-nav current="home"></shell-nav>`
    }

    // API は空の配列を省略して返すため、必須型でもフォールバックする
    const links = p.links ?? []
    let i = 0
    return html`
      <h1 class="ln" style="--i:${i++}"><span class="h1">${p.displayName}</span></h1>
      <p class="ln" style="--i:${i++}"></p>
      <p class="ln jp" style="--i:${i++}">${[p.displayJa, [p.role, p.location].filter(Boolean).join(', ')].filter(Boolean).join(' — ')}</p>
      ${
        links.length > 0
          ? html`
            <p class="ln" style="--i:${i++}"></p>
            <p class="ln h2" style="--i:${i++}">## links</p>
            <ul>
              ${links.map(
                (link) => html`
                  <li>
                    <a class="ln" style="--i:${i++}" href=${sanitizeUrl(link.url)} target="_blank" rel="noopener noreferrer"><span class="m">•</span> ${link.label ?? link.platform}  <span class="m">${link.url.replace(/^https?:\/\//, '')}</span></a>
                  </li>
                `,
              )}
            </ul>
          `
          : nothing
      }
      <shell-nav current="home" style="animation-delay:${i * 16}ms"></shell-nav>
    `
  }

  static styles = [
    pageStyles,
    css`
      p,
      h1 {
        margin: 0;
        font: inherit;
      }

      .h1 {
        display: inline-block;
        padding: 0 0.6em;
        background: var(--color-text-primary);
        color: var(--color-bg-deep);
        font-weight: 500;
      }

      .h2 {
        color: var(--color-accent);
        font-weight: 500;
      }
    `,
  ]
}

declare global {
  interface HTMLElementTagNameMap {
    'page-top': PageTop
  }
}
