import { css, html, LitElement, nothing } from 'lit'
import { customElement, property, state } from 'lit/decorators.js'
import { listArticles } from '../admin/article-api.js'
import type { ArticleItem } from '../admin/article-types.js'
import type { MeProfile } from '../admin/types.js'
import { setupAmbientLines } from '../utils/ambient.js'
import { setupFade, setupReveal } from '../utils/scroll.js'

@customElement('page-top')
export class PageTop extends LitElement {
  @property({ attribute: false }) profile: MeProfile | null = null
  @property({ type: Boolean }) loading = false

  @state()
  private articles: ArticleItem[] = []

  @state()
  private articlesLoading = true

  @state()
  private articlesError = ''

  private cleanups: Array<() => void> = []

  firstUpdated() {
    const root = this.shadowRoot
    if (!root) return
    const revealEls = Array.from(
      root.querySelectorAll('.who > *, .articles-preview > *, .contact > *'),
    )
    const fadeEls = Array.from(root.querySelectorAll('.layer-1, .layer-2'))
    const fvSection = root.querySelector('.layer-0') as HTMLElement | null

    this.cleanups.push(setupReveal(revealEls, true))
    this.cleanups.push(setupFade(fadeEls))
    if (fvSection) this.cleanups.push(setupAmbientLines(fvSection))
    void this.loadArticles()
  }

  disconnectedCallback() {
    super.disconnectedCallback()
    for (const cleanup of this.cleanups) cleanup()
    this.cleanups = []
  }

  render() {
    return html`
      <!-- Layer 0: First View -->
      <section class="layer layer-0 js-layer-0">
        <h1 class="name ${this.loading ? 'is-loading' : ''}">
          ${this.profile?.displayName ?? ''}
        </h1>
      </section>

      <!-- Layer 1: Who I am -->
      <section class="layer layer-1">
        <div class="who">
          <p class="role ${this.loading ? 'is-loading' : ''}">${this.profile?.role ?? ''}</p>
          <p class="location ${this.loading ? 'is-loading' : ''}">${this.profile?.location ?? ''}</p>
        </div>
      </section>

      <!-- Layer 2: Articles -->
      <section class="layer layer-2">
        <div class="articles-preview">
          <h2 class="section-label">Articles</h2>
          <ul class="article-list">
            ${this.renderArticlePreview()}
          </ul>
          ${
            this.articlesError
              ? html`<p class="articles-error">${this.articlesError}</p>`
              : null
          }
          <a href="/articles" class="view-all">View all</a>
        </div>
      </section>

      <!-- Layer 3: Contact -->
      <section class="layer layer-3">
        <div class="contact">
          <p class="contact-label">Say Hello</p>
          <ul class="contact-links">
            ${
              this.profile
                ? this.profile.links.map(
                    (link) => html`
                    <li>
                      <a href=${link.url} target="_blank" rel="noopener">
                        ${link.platform}
                      </a>
                    </li>
                  `,
                  )
                : nothing
            }
          </ul>
        </div>
      </section>
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
      this.articlesError = '記事を読み込めませんでした。'
    } finally {
      this.articlesLoading = false
    }
  }

  private renderArticlePreview() {
    if (this.articlesLoading) {
      return Array.from(
        { length: 3 },
        (_, index) => html`
        <li class="article-item">
          <span class="article-date is-loading">----</span>
          <span class="article-title is-loading">Loading article ${index + 1}</span>
        </li>
      `,
      )
    }

    if (this.articles.length === 0) {
      return html`
        <li class="article-item">
          <span class="article-date">----</span>
          <a href="/articles" class="article-title">記事一覧を見る</a>
        </li>
      `
    }

    return this.articles.map(
      (article) => html`
        <li class="article-item">
          <span class="article-date">${this.formatArticleDate(article.publishedAt)}</span>
          <a
            href=${article.url}
            class="article-title"
            target="_blank"
            rel="noopener noreferrer"
          >
            ${article.title}
          </a>
        </li>
      `,
    )
  }

  private formatArticleDate(value?: string) {
    if (!value) return '----.--'

    const date = new Date(value)
    if (Number.isNaN(date.valueOf())) return '----.--'

    return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`
  }

  static styles = css`
    :host {
      display: block;
    }

    .layer {
      padding: var(--space-lg, 64px) var(--space-md, 32px);
    }

    /* Layer 0 */
    .layer-0 {
      height: 100dvh;
      display: flex;
      align-items: center;
      justify-content: center;
      padding: 0;
    }

    .name {
      font-family: var(--font-en);
      font-weight: 300;
      font-size: clamp(1.5rem, 2.5vw, 2.5rem);
      letter-spacing: var(--tracking-wider);
      color: var(--color-text-primary);
      margin: 0;
      animation: breathing 7s ease-in-out infinite;
    }

    @keyframes breathing {
      0%, 100% { opacity: 0.85; }
      50% { opacity: 1; }
    }

    /* Layer 1 */
    .layer-1 {
      display: block;
    }

    .who {
      max-width: 640px;
      margin: 0 auto;
    }

    .role {
      font-family: var(--font-en);
      font-weight: 300;
      font-size: 24px;
      letter-spacing: var(--tracking-wide);
      color: var(--color-text-primary);
      margin-bottom: 8px;
    }

    .location {
      font-family: var(--font-jp);
      font-weight: 200;
      font-size: 14px;
      letter-spacing: var(--tracking-wide);
      color: var(--color-text-secondary);
      margin-bottom: 40px;
    }

    .philosophy {
      font-family: var(--font-jp);
      font-weight: 200;
      font-size: 18px;
      letter-spacing: var(--tracking-wide);
      color: var(--color-text-primary);
      line-height: 2;
    }

    /* Layer 2 */
    .layer-2 {
      display: block;
    }

    .articles-preview {
      max-width: 640px;
      margin: 0 auto;
      width: 100%;
    }

    .section-label {
      font-family: var(--font-en);
      font-weight: 300;
      font-size: 13px;
      letter-spacing: var(--tracking-wider);
      text-transform: uppercase;
      color: var(--color-text-secondary);
      margin-bottom: 32px;
    }

    .article-list {
      list-style: none;
      padding: 0;
      margin: 0 0 40px;
    }

    .article-item {
      display: flex;
      align-items: baseline;
      gap: 24px;
      padding: 16px 8px;
      border-bottom: 1px solid var(--color-border);
      transition: background 0.2s ease, transform 0.2s ease;
    }

    .article-item:hover {
      background: var(--color-surface);
      transform: translateX(4px);
    }

    .article-date {
      font-family: var(--font-en);
      font-weight: 300;
      font-size: 13px;
      letter-spacing: var(--tracking-wide);
      color: var(--color-text-secondary);
      white-space: nowrap;
    }

    .article-title {
      font-family: var(--font-jp);
      font-weight: 200;
      font-size: 15px;
      letter-spacing: 0.04em;
      color: var(--color-text-primary);
      text-decoration: none;
      transition: opacity 0.2s ease;
    }

    .article-title:hover {
      opacity: 0.6;
    }

    .view-all {
      font-family: var(--font-en);
      font-weight: 300;
      font-size: 14px;
      letter-spacing: var(--tracking-wide);
      color: var(--color-text-secondary);
      text-decoration: none;
      transition: opacity 0.2s ease;
    }

    .view-all:hover {
      opacity: 0.6;
    }

    .articles-error {
      margin: -24px 0 24px;
      color: var(--color-text-tertiary);
      font-size: 13px;
      line-height: 1.8;
    }

    /* Layer 3 */
    .layer-3 {
      display: block;
    }

    .contact {
      max-width: 640px;
      margin: 0 auto;
    }

    .contact-label {
      font-family: var(--font-en);
      font-weight: 300;
      font-size: 13px;
      letter-spacing: var(--tracking-wider);
      text-transform: uppercase;
      color: var(--color-text-secondary);
      margin: 0 0 20px;
    }

    .contact-links {
      list-style: none;
      padding: 0;
      margin: 0;
    }

    .contact-links li {
      border-bottom: 1px solid var(--color-border-light);
    }

    .contact-links li:first-child {
      border-top: 1px solid var(--color-border-light);
    }

    .contact-links a {
      display: block;
      font-family: var(--font-en);
      font-weight: 300;
      font-size: 15px;
      letter-spacing: var(--tracking-wide);
      color: var(--color-text-primary);
      text-decoration: none;
      padding: 12px 0;
      transition: opacity 0.2s ease;
    }

    .contact-links a:hover {
      opacity: 0.5;
    }

    .is-loading {
      opacity: 0.3;
    }

    @media (prefers-reduced-motion: reduce) {
      .name {
        animation: none;
        opacity: 1;
      }

      .article-item {
        transition: none;
        transform: none;
      }

      .article-title,
      .view-all,
      .contact-links a {
        transition: none;
      }
    }

    @media (max-width: 640px) {
      .layer {
        padding: var(--space-md, 32px) var(--space-sm, 16px);
      }
    }
  `
}

declare global {
  interface HTMLElementTagNameMap {
    'page-top': PageTop
  }
}
