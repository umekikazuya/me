import { consume } from '@lit/context'
import type { components } from '@me/types'
import { css, html, LitElement, nothing } from 'lit'
import { customElement, state } from 'lit/decorators.js'
import { listArticles } from '../api/article-api.js'
import { profileContext } from '../contexts/profile-context.js'
import type { IProfileRepository } from '../domain/ProfileRepository.js'
import { pageStyles } from '../styles/page-styles.js'
import { formatDate, sanitizeUrl } from '../utils/format.js'
import '../components/load-spinner.js'
import '../components/term-prompt.js'

type Profile = components['schemas']['MeResponse']

/** whoami の出力行。名前・肩書きと、上位のスキルカテゴリ */
export function whoamiLines(p: Profile): string[] {
  // API は空の配列を省略して返すため、必須型でもフォールバックする
  const skills = [...(p.skills ?? [])]
    .sort((a, b) => a.sortOrder - b.sortOrder)
    .slice(0, 3)
    .flatMap((group) => group.items.slice(0, 2))
  return [
    p.displayJa ? `${p.displayName} (${p.displayJa})` : p.displayName,
    [p.role, p.location].filter(Boolean).join(', '),
    ...(skills.length > 0 ? [skills.join(' / ')] : []),
  ]
}

@customElement('page-top')
export class PageTop extends LitElement {
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
  private articles: components['schemas']['ArticleItem'][] = []

  @state()
  private articlesLoading = true

  @state()
  private articlesError = ''

  firstUpdated() {
    void this.loadArticles()
  }

  disconnectedCallback() {
    super.disconnectedCallback()
    if (this._profileRepo) {
      this._profileRepo.removeEventListener('change', this._onRepoChange)
    }
  }

  render() {
    const p = this.profileRepo.profile

    return html`
      <header class="hero">
        <h1>${p?.displayName ?? html`&nbsp;`}</h1>
        <p class="role">
          ${p ? [p.role, p.location].filter(Boolean).join(', ').toLowerCase() : html`&nbsp;`}
        </p>
      </header>

      <section>
        <h2>~</h2>
        <term-prompt
          command="whoami"
          .output=${p ? whoamiLines(p) : this.profileRepo.error ? ['whoami: profile unavailable'] : []}
        ></term-prompt>
      </section>

      <section>
        <h2>writing</h2>
        ${this.renderArticles()}
        <a href="/articles" class="more">all writing →</a>
      </section>

      ${
        p && (p.links ?? []).length > 0
          ? html`
            <section>
              <h2>elsewhere</h2>
              <ul class="rows">
                ${(p.links ?? []).map(
                  (link) => html`
                    <li>
                      <a class="row" href=${sanitizeUrl(link.url)} target="_blank" rel="noopener noreferrer">
                        <span class="main">${link.label ?? link.platform}</span>
                        <span class="meta">↗</span>
                      </a>
                    </li>
                  `,
                )}
              </ul>
            </section>
          `
          : nothing
      }
    `
  }

  private renderArticles() {
    if (this.articlesLoading) {
      return html`<load-spinner class="pending"></load-spinner>`
    }
    if (this.articlesError) {
      return html`<p class="note">${this.articlesError}</p>`
    }
    return html`
      <ul class="rows">
        ${this.articles.map(
          (article) => html`
            <li>
              <a class="row" href=${sanitizeUrl(article.url)} target="_blank" rel="noopener noreferrer">
                <span class="main">${article.title}</span>
                <time class="meta" datetime=${article.publishedAt ?? nothing}>
                  ${formatDate(article.publishedAt)}
                </time>
              </a>
            </li>
          `,
        )}
      </ul>
    `
  }

  private async loadArticles() {
    this.articlesLoading = true
    this.articlesError = ''
    try {
      const result = await listArticles({ limit: 5 })
      this.articles = result.articles
    } catch {
      this.articles = []
      this.articlesError =
        '記事を読み込めませんでした。時間をおいて再読み込みしてください。'
    } finally {
      this.articlesLoading = false
    }
  }

  static styles = [
    pageStyles,
    css`
      .hero {
        padding-block: clamp(96px, 18vh, 160px) 72px;
      }

      h1 {
        margin: 0 0 6px;
        font-family: var(--font-en);
        font-size: 22px;
        font-weight: 500;
        letter-spacing: var(--tracking-tight);
        line-height: 1.4;
      }

      .role {
        margin: 0;
        font-family: var(--font-mono);
        font-size: 12px;
        color: var(--color-text-tertiary);
      }

      .pending {
        padding: 11px 0;
        border-top: 1px solid var(--color-border);
        width: 100%;
      }
    `,
  ]
}

declare global {
  interface HTMLElementTagNameMap {
    'page-top': PageTop
  }
}
