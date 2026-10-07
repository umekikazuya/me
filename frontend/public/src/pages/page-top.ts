import { consume } from '@lit/context'
import type { components } from '@me/types'
import { css, html, nothing, type TemplateResult } from 'lit'
import { customElement, state } from 'lit/decorators.js'
import { listArticles } from '../api/article-api.js'
import { BIO, NOW } from '../content/home.js'
import { profileContext } from '../contexts/profile-context.js'
import type { IProfileRepository } from '../domain/ProfileRepository.js'
import { pageStyles } from '../styles/page-styles.js'
import { localTime } from '../utils/clock.js'
import { formatDate, sanitizeUrl } from '../utils/format.js'
import { formatLogin, takeLastLogin } from '../utils/last-login.js'
import { Session } from '../utils/session.js'
import { shortHash } from './page-articles.js'
import { ShellPage } from './shell-page.js'
import '../components/load-spinner.js'

type Article = components['schemas']['ArticleItem']

const LINE_STEP_MS = 16
const LATEST_COUNT = 3

const COMMANDS = [
  'glow README.md',
  'cat now.md',
  `git log --oneline -${LATEST_COUNT}`,
  'date',
] as const

/** トップ。README・now・最新記事・時刻を1つのセッションとして順に打つ。 */
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

  @state()
  private articles: Article[] = []

  @state()
  private articlesLoading = true

  @state()
  private articlesError = ''

  private session = new Session(this)
  private lastLogin = takeLastLogin()
  private waitingForArticles = -1

  firstUpdated() {
    void this.loadArticles()
  }

  disconnectedCallback() {
    super.disconnectedCallback()
    if (this._profileRepo) {
      this._profileRepo.removeEventListener('change', this._onRepoChange)
    }
  }

  private onCommandTyped(e: Event, index: number) {
    e.stopPropagation()
    if (COMMANDS[index].startsWith('git log') && this.articlesLoading) {
      this.waitingForArticles = index
      this.session.reveal(index)
      return
    }
    this.session.advance(index, COMMANDS.length, this.outputMs(index))
  }

  private outputMs(index: number) {
    const readmeLines = 3 + (this.profileRepo.profile?.links ?? []).length
    const lines =
      [readmeLines, NOW.length, this.articles.length + 1, 1][index] ?? 1
    return lines * LINE_STEP_MS
  }

  render() {
    const done = this.session.outputs >= COMMANDS.length
    return html`
      <h1 class="sr-only">umekikazuya</h1>
      ${
        this.lastLogin
          ? html`<p class="ln m login">Last login: ${formatLogin(this.lastLogin)} on ttys001</p>`
          : nothing
      }
      ${COMMANDS.slice(0, this.session.commands).map(
        (command, index) => html`
          <section class="block">
            <shell-command
              command=${command}
              @typed=${(e: Event) => this.onCommandTyped(e, index)}
            ></shell-command>
            ${index < this.session.outputs ? this.renderOutput(index) : nothing}
          </section>
        `,
      )}
      ${done ? html`<shell-prompt></shell-prompt>` : nothing}
    `
  }

  private renderOutput(index: number): TemplateResult {
    switch (index) {
      case 0:
        return this.renderReadme()
      case 1:
        return html`${NOW.map(
          (line, i) =>
            html`<p class="ln jp ${i > 0 ? 'm' : ''}" style="--i:${i}">${line}</p>`,
        )}`
      case 2:
        return this.renderLatest()
      default: {
        const time = localTime(new Date())
        return html`<p class="ln">${time.text} <span class="m">— ${time.mood}</span></p>`
      }
    }
  }

  private renderReadme() {
    const profile = this.profileRepo.profile
    const name = profile?.displayName ?? 'umekikazuya'
    const links = profile?.links ?? []
    return html`
      <p class="ln"><span class="name">${name}</span></p>
      <p class="ln jp bio" style="--i:1">${BIO}</p>
      ${
        links.length > 0
          ? html`
            <p class="ln h2" style="--i:2">## contact</p>
            <ul class="links">
              ${links.map(
                (link, i) => html`
                  <li>
                    <a class="ln link" style="--i:${i + 3}" href=${sanitizeUrl(link.url)} target="_blank" rel="noopener noreferrer"><span class="label">${link.label ?? link.platform}</span><span class="arrow" aria-hidden="true">↗</span></a>
                  </li>
                `,
              )}
            </ul>
          `
          : nothing
      }
    `
  }

  private renderLatest() {
    if (this.articlesLoading)
      return html`<load-spinner class="ln"></load-spinner>`
    if (this.articlesError)
      return html`<p class="ln m">${this.articlesError}</p>`
    return html`
      <ul>
        ${this.articles.map(
          (article, i) => html`
            <li>
              <a class="ln row" style="--i:${i}" href=${sanitizeUrl(article.url)} target="_blank" rel="noopener noreferrer"><span class="a hash">${shortHash(article.externalId)}</span><span class="m">${formatDate(article.publishedAt).replaceAll('.', '-')}</span><span class="jp t">${article.title}</span></a>
            </li>
          `,
        )}
      </ul>
      <a class="ln more" style="--i:${this.articles.length}" href="/articles"><span class="m">→ all writing</span></a>
    `
  }

  private async loadArticles() {
    try {
      const result = await listArticles({ limit: LATEST_COUNT })
      this.articles = result.articles
    } catch {
      this.articlesError =
        '記事を読み込めませんでした。時間をおいて再読み込みしてください。'
    } finally {
      this.articlesLoading = false
      if (this.waitingForArticles >= 0) {
        const index = this.waitingForArticles
        this.waitingForArticles = -1
        this.session.continueAfter(index, COMMANDS.length, this.outputMs(index))
      }
    }
  }

  static styles = [
    pageStyles,
    css`
      p {
        margin: 0;
      }

      .login {
        margin-bottom: var(--session-gap);
      }

      .name {
        display: inline-block;
        padding: 0.1em 0.6em;
        background: var(--color-text-primary);
        color: var(--color-bg-deep);
        font-weight: 500;
        letter-spacing: 0.02em;
      }

      .bio {
        margin-top: 0.9em;
      }

      .links {
        display: flex;
        flex-wrap: wrap;
        gap: 0 2ch;
      }

      a.link {
        display: inline-flex;
        gap: 0.4ch;
        width: auto;
      }

      .label {
        text-decoration: underline;
        text-decoration-color: var(--color-text-tertiary);
        text-decoration-thickness: 1px;
        text-underline-offset: 0.3em;
        transition: text-decoration-color 0.2s ease;
      }

      .arrow {
        color: var(--color-text-tertiary);
        transition:
          color 0.2s ease,
          transform 0.2s ease;
      }

      .link:hover .label,
      .link:focus-visible .label {
        text-decoration-color: var(--color-accent);
      }

      .link:hover .arrow,
      .link:focus-visible .arrow {
        color: var(--color-accent);
        transform: translate(1px, -1px);
      }

      .h2 {
        margin-top: 0.9em;
        color: var(--color-accent);
        font-weight: 500;
      }

      .row {
        display: grid;
        grid-template-columns: 8ch 11ch minmax(0, 1fr);
      }

      .more {
        margin-top: 0.6em;
      }

      @media (max-width: 520px) {
        .row {
          grid-template-columns: 11ch minmax(0, 1fr);
        }

        .hash {
          display: none;
        }
      }
    `,
  ]
}

declare global {
  interface HTMLElementTagNameMap {
    'page-top': PageTop
  }
}
