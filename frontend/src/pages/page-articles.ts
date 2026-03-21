import { LitElement, css, html } from 'lit'
import { customElement } from 'lit/decorators.js'
import { setupReveal } from '../utils/scroll.js'

interface ArticleItem {
  date: string
  title: string
  tags: string[]
}

interface ArticleGroup {
  year: number
  items: ArticleItem[]
}

const mockTags: string[] = [
  'TypeScript',
  'Web Components',
  'Design',
  'Go',
  'Architecture',
  'CSS',
]

const mockArticles: ArticleGroup[] = [
  {
    year: 2025,
    items: [
      {
        date: '03',
        title: 'TypeScriptの型システムと向き合う',
        tags: ['TypeScript'],
      },
      {
        date: '01',
        title: 'Litで作るWeb Components入門',
        tags: ['Web Components', 'TypeScript'],
      },
    ],
  },
  {
    year: 2024,
    items: [
      { date: '11', title: '余白という設計思想', tags: ['Design', 'CSS'] },
      {
        date: '08',
        title: 'GoでつくるClean Architecture',
        tags: ['Go', 'Architecture'],
      },
      {
        date: '03',
        title: 'CSSカスタムプロパティ設計の実践',
        tags: ['CSS', 'Design'],
      },
    ],
  },
]

@customElement('page-articles')
export class PageArticles extends LitElement {
  private cleanups: Array<() => void> = []

  firstUpdated() {
    const root = this.shadowRoot
    if (!root) return
    const revealEls = Array.from(
      root.querySelectorAll(
        '.page-header, .search-area, .tag-cloud, .year-group',
      ),
    )
    this.cleanups.push(setupReveal(revealEls, true))
  }

  disconnectedCallback() {
    super.disconnectedCallback()
    for (const cleanup of this.cleanups) cleanup()
    this.cleanups = []
  }

  render() {
    return html`
      <div class="container">
        <header class="page-header">
          <h1 class="page-title">Articles</h1>
        </header>

        <div class="search-area">
          <input
            type="search"
            class="search-input"
            placeholder="Search..."
            aria-label="記事を検索"
          />
        </div>

        <div class="tag-cloud">
          ${mockTags.map((tag) => html`<span class="tag">${tag}</span>`)}
        </div>

        <div class="timeline">
          ${mockArticles.map(
            (group) => html`
            <div class="year-group">
              <div class="year-label">${group.year}</div>
              <ul class="article-list">
                ${group.items.map(
                  (article) => html`
                  <li class="article-row">
                    <span class="article-date">${group.year}.${article.date}</span>
                    <span class="article-title">${article.title}</span>
                    <div class="article-tags">
                      ${article.tags.map((t) => html`<span class="article-tag">${t}</span>`)}
                    </div>
                  </li>
                `,
                )}
              </ul>
            </div>
          `,
          )}
        </div>
      </div>
    `
  }

  static styles = css`
    :host {
      display: block;
      padding-top: 80px;
    }

    .container {
      max-width: 720px;
      margin: 0 auto;
      padding: var(--space-lg) var(--space-md);
    }

    .page-header {
      margin-bottom: 48px;
    }

    .page-title {
      font-family: var(--font-en);
      font-weight: 300;
      font-size: 36px;
      letter-spacing: var(--tracking-wide);
      color: var(--color-text-primary);
      margin: 0;
    }

    .search-area {
      margin-bottom: 32px;
    }

    .search-input {
      width: 100%;
      background: transparent;
      border: none;
      border-bottom: 0.5px solid var(--color-border);
      outline: none;
      font-family: var(--font-jp);
      font-weight: 200;
      font-size: 15px;
      color: var(--color-text-primary);
      padding: 8px 0;
      letter-spacing: 0.04em;
      border-radius: 0;
    }

    .search-input::placeholder {
      color: var(--color-text-tertiary);
    }

    .tag-cloud {
      display: flex;
      flex-wrap: wrap;
      gap: 8px;
      margin-bottom: 64px;
    }

    .tag {
      font-family: var(--font-en);
      font-weight: 300;
      font-size: 13px;
      letter-spacing: var(--tracking-wide);
      color: var(--color-text-secondary);
      cursor: pointer;
      transition: opacity 0.2s ease;
    }

    .tag:hover {
      opacity: 0.6;
    }

    .year-group {
      margin-bottom: 48px;
    }

    .year-label {
      font-family: var(--font-en);
      font-weight: 300;
      font-size: 13px;
      letter-spacing: var(--tracking-wider);
      color: var(--color-text-tertiary);
      margin-bottom: 16px;
    }

    .article-list {
      list-style: none;
      padding: 0;
      margin: 0;
    }

    .article-row {
      display: grid;
      grid-template-columns: 80px 1fr auto;
      align-items: baseline;
      gap: 16px;
      padding: 16px 8px;
      border-bottom: 1px solid var(--color-border-light);
      transition: background 0.2s ease, transform 0.2s ease;
    }

    .article-row:hover {
      background: var(--color-surface);
      transform: translateX(4px);
    }

    .article-date {
      font-family: var(--font-en);
      font-weight: 300;
      font-size: 13px;
      letter-spacing: var(--tracking-wide);
      color: var(--color-text-tertiary);
      white-space: nowrap;
    }

    .article-title {
      font-family: var(--font-jp);
      font-weight: 200;
      font-size: 15px;
      letter-spacing: 0.04em;
      color: var(--color-text-primary);
      cursor: pointer;
      transition: opacity 0.2s ease;
    }

    .article-title:hover {
      opacity: 0.6;
    }

    .article-tags {
      display: flex;
      gap: 4px;
    }

    .article-tag {
      font-family: var(--font-en);
      font-weight: 300;
      font-size: 11px;
      letter-spacing: var(--tracking-wide);
      color: var(--color-text-tertiary);
    }

    @media (prefers-reduced-motion: reduce) {
      .tag,
      .article-title,
      .article-row {
        transition: none;
        transform: none;
      }
    }

    @media (max-width: 640px) {
      .container {
        padding: 48px 24px;
      }

      .article-row {
        grid-template-columns: 1fr;
        gap: 4px;
      }

      .article-tags {
        margin-top: 4px;
      }
    }
  `
}

declare global {
  interface HTMLElementTagNameMap {
    'page-articles': PageArticles
  }
}
