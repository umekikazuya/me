import type { components } from '@me/types'
import { css, html, LitElement, nothing } from 'lit'
import { customElement, state } from 'lit/decorators.js'
import {
  listArticles,
  listArticleTags,
  suggestArticles,
} from '../api/article-api.js'
import { describeApiError } from '../api/types.js'
import { pageStyles } from '../styles/page-styles.js'
import { formatDate, sanitizeUrl } from '../utils/format.js'
import '../components/load-spinner.js'

interface ArticleGroup {
  key: string
  label: string
  items: components['schemas']['ArticleItem'][]
}

@customElement('page-articles')
export class PageArticles extends LitElement {
  @state()
  private articles: components['schemas']['ArticleItem'][] = []

  @state()
  private tagOptions: components['schemas']['ArticleTagItem'][] = []

  @state()
  private suggestions: components['schemas']['ArticleSuggestionItem'][] = []

  @state()
  private query = ''

  @state()
  private appliedQuery = ''

  @state()
  private selectedTags: string[] = []

  @state()
  private loading = false

  @state()
  private loadingMore = false

  @state()
  private showAllTags = false

  @state()
  private suggestionLoading = false

  @state()
  private errorMessage = ''

  @state()
  private nextCursor?: string

  private suggestTimer?: number
  private articleRequestId = 0
  private tagRequestId = 0
  private suggestionRequestId = 0

  firstUpdated() {
    void this.loadInitialData()
  }

  disconnectedCallback() {
    super.disconnectedCallback()
    if (this.suggestTimer !== undefined) {
      window.clearTimeout(this.suggestTimer)
    }
  }

  private get displayedTags() {
    const sorted = [...this.tagOptions].sort((a, b) => b.count - a.count)
    if (this.showAllTags) return sorted
    return sorted.slice(0, 12)
  }

  render() {
    return html`
      <div class="container">
        ${this.renderHeader()}
        ${this.renderSearchArea()}
        ${this.renderTagCloud()}

        ${this.errorMessage ? html`<p class="message error">${this.errorMessage}</p>` : null}

        <div class="timeline">
          ${this.loading ? html`<load-spinner class="pending"></load-spinner>` : this.renderArticleGroups()}
        </div>

        ${this.renderLoadMore()}
      </div>
    `
  }

  private renderHeader() {
    return html`
      <header class="page-header">
        <h1 class="page-title">writing</h1>
        ${
          this.selectedTags.length > 0 || this.appliedQuery
            ? html`<p class="page-description">${this.describeFilters()}</p>`
            : null
        }
      </header>
    `
  }

  private renderSearchArea() {
    return html`
      <div class="search-area">
        <form @submit=${this.handleSearch}>
          <input
            type="search"
            class="search-input"
            .value=${this.query}
            placeholder="/ search title or tag"
            aria-label="記事を検索"
            @input=${this.handleQueryInput}
          />
        </form>
        ${this.renderSuggestions()}
      </div>
    `
  }

  private renderSuggestions() {
    if (this.suggestionLoading) {
      return html`<load-spinner class="search-status" label="searching"></load-spinner>`
    }
    if (this.suggestions.length === 0) return null

    return html`
      <ul class="suggestion-list">
        ${this.suggestions.map(
          (s) => html`
            <li>
              <button type="button" class="suggestion-item" @click=${() => this.handleSuggestionSelect(s)}>
                <span class="suggestion-value">${s.value}</span>
                <span class="suggestion-meta">${s.type === 'title' ? '記事' : html`${s.type} · ${s.count}`}</span>
              </button>
            </li>
          `,
        )}
      </ul>
    `
  }

  private renderTagCloud() {
    return html`
      <div class="tag-cloud">
        ${this.displayedTags.map((tag) => this.renderTag(tag))}
        ${this.tagOptions.length > 12 ? this.renderTagToggle() : null}
      </div>
    `
  }

  private renderTag(tag: components['schemas']['ArticleTagItem']) {
    const isSelected = this.selectedTags.includes(tag.name)
    return html`
      <button
        type="button"
        class=${isSelected ? 'tag selected' : 'tag'}
        aria-pressed=${isSelected}
        @click=${() => this.toggleTag(tag.name)}
      >
        <span class="tag-hash">#</span>
        <span class="tag-name">${tag.name}</span>
        <small class="tag-count">${tag.count}</small>
      </button>
    `
  }

  private renderTagToggle() {
    return html`
      <button type="button" class="tag-toggle" @click=${() => (this.showAllTags = !this.showAllTags)}>
        ${this.showAllTags ? '— show less' : `+ ${this.tagOptions.length - 12} more`}
      </button>
    `
  }

  private renderArticleGroups() {
    if (this.articleGroups.length === 0) {
      return html`
        <section class="empty-state">
          <p>条件に一致する記事がありません。</p>
          <button type="button" class="ghost-button" @click=${this.clearFilters}>条件をリセット</button>
        </section>
      `
    }

    return this.articleGroups.map(
      (group) => html`
        <div class="year-group">
          <div class="year-label">${group.label}</div>
          <ul class="article-list">
            ${group.items.map((article) => this.renderArticleRow(article))}
          </ul>
        </div>
      `,
    )
  }

  private renderArticleRow(article: components['schemas']['ArticleItem']) {
    return html`
      <li class="article-row">
        <a href=${sanitizeUrl(article.url)} class="row" target="_blank" rel="noopener noreferrer">
          <span class="main">${article.title}</span>
          <time class="meta" datetime=${article.publishedAt ?? nothing}>
            ${formatDate(article.publishedAt).slice(5)}
          </time>
        </a>
        ${
          article.tags?.length
            ? html`<div class="article-tags">
              ${article.tags.map(
                (tag) =>
                  html`<button type="button" class="article-tag" @click=${() => this.toggleTag(tag)}>#${tag}</button>`,
              )}
            </div>`
            : nothing
        }
      </li>
    `
  }

  private renderLoadMore() {
    if (!this.nextCursor || this.loading) return null
    return html`
      <div class="load-more">
        <button type="button" class="ghost-button" ?disabled=${this.loadingMore} @click=${this.handleLoadMore}>
          ${this.loadingMore ? 'loading…' : 'load more'}
        </button>
      </div>
    `
  }

  private get articleGroups(): ArticleGroup[] {
    const groups = new Map<string, ArticleGroup>()

    for (const article of this.articles) {
      const date = article.publishedAt ? new Date(article.publishedAt) : null
      const key =
        date && !Number.isNaN(date.valueOf())
          ? String(date.getFullYear())
          : 'undated'
      const label = key === 'undated' ? 'Archive' : key

      const group = groups.get(key) ?? { key, label, items: [] }
      group.items.push(article)
      groups.set(key, group)
    }

    return Array.from(groups.values())
  }

  private async loadInitialData() {
    this.loading = true
    this.loadingMore = false
    this.errorMessage = ''
    const articleRequestId = ++this.articleRequestId
    const tagRequestId = ++this.tagRequestId

    const [articlesResult, tagsResult] = await Promise.allSettled([
      listArticles({ limit: 50 }),
      listArticleTags(),
    ])

    if (articleRequestId === this.articleRequestId) {
      if (articlesResult.status === 'fulfilled') {
        this.errorMessage = ''
        this.articles = articlesResult.value.articles
        this.nextCursor = articlesResult.value.nextCursor
      } else {
        this.errorMessage = describeApiError(articlesResult.reason)
      }

      this.loading = false
    }

    if (tagRequestId === this.tagRequestId) {
      if (tagsResult.status === 'fulfilled') {
        this.tagOptions = tagsResult.value.tags
      } else if (
        articleRequestId === this.articleRequestId &&
        !this.errorMessage
      ) {
        this.errorMessage = describeApiError(tagsResult.reason)
      }
    }
  }

  private async reloadArticles(cursor?: string, append = false) {
    const requestId = ++this.articleRequestId

    if (append) {
      this.loading = false
      this.loadingMore = true
    } else {
      this.loading = true
      this.loadingMore = false
      this.errorMessage = ''
    }

    try {
      const result = await listArticles({
        q: this.appliedQuery || undefined,
        tag: this.selectedTags,
        limit: 50,
        cursor,
      })

      if (requestId !== this.articleRequestId) return
      this.errorMessage = ''

      this.articles = append
        ? [...this.articles, ...result.articles]
        : result.articles
      this.nextCursor = result.nextCursor
    } catch (error) {
      if (requestId !== this.articleRequestId) return
      this.errorMessage = describeApiError(error)
    } finally {
      if (requestId === this.articleRequestId) {
        this.loading = false
        this.loadingMore = false
      }
    }
  }

  private handleQueryInput = (event: Event) => {
    this.query = (event.target as HTMLInputElement).value
    this.scheduleSuggest()
  }

  private handleSearch = (event: Event) => {
    event.preventDefault()
    this.applySearch(this.query)
  }

  private handleSuggestionSelect(
    suggestion: components['schemas']['ArticleSuggestionItem'],
  ) {
    if (suggestion.type === 'tag') {
      this.invalidateSuggestions()
      this.query = ''
      this.appliedQuery = ''
      this.toggleTag(suggestion.value)
      return
    }

    if (suggestion.type === 'title') {
      // 読み込み済みの一覧に該当記事があれば直接開く。
      // 無ければタイトルで検索してユーザーに見つけてもらう。
      const article = this.articles.find(
        (item) => item.externalId === suggestion.externalId,
      )
      if (article) {
        this.invalidateSuggestions()
        window.open(article.url, '_blank', 'noreferrer')
        return
      }
    }

    this.query = suggestion.value
    this.applySearch(suggestion.value)
  }

  private handleLoadMore = () => {
    if (!this.nextCursor) return
    void this.reloadArticles(this.nextCursor, true)
  }

  private toggleTag(tagName: string) {
    this.selectedTags = this.selectedTags.includes(tagName)
      ? this.selectedTags.filter((tag) => tag !== tagName)
      : [...this.selectedTags, tagName]

    void this.reloadArticles()
  }

  private clearFilters = () => {
    this.invalidateSuggestions()
    this.query = ''
    this.appliedQuery = ''
    this.selectedTags = []
    void this.reloadArticles()
  }

  private applySearch(query: string) {
    this.invalidateSuggestions()
    this.appliedQuery = query.trim()
    this.query = query
    void this.reloadArticles()
  }

  private scheduleSuggest() {
    this.clearSuggestTimer()

    const query = this.query.trim()
    if (query === '') {
      this.invalidateSuggestions()
      return
    }

    const requestId = ++this.suggestionRequestId
    this.suggestionLoading = true
    this.suggestTimer = window.setTimeout(() => {
      this.suggestTimer = undefined
      void this.loadSuggestions(query, requestId)
    }, 150)
  }

  private async loadSuggestions(query: string, requestId: number) {
    try {
      const suggestions = await suggestArticles(query)
      if (requestId !== this.suggestionRequestId || this.query.trim() !== query)
        return

      this.suggestions = suggestions.suggestions
    } catch {
      if (requestId !== this.suggestionRequestId || this.query.trim() !== query)
        return

      this.suggestions = []
    } finally {
      if (
        requestId === this.suggestionRequestId &&
        this.query.trim() === query
      ) {
        this.suggestionLoading = false
      }
    }
  }

  private clearSuggestTimer() {
    if (this.suggestTimer === undefined) return

    window.clearTimeout(this.suggestTimer)
    this.suggestTimer = undefined
  }

  private invalidateSuggestions() {
    this.clearSuggestTimer()
    this.suggestionRequestId += 1
    this.suggestionLoading = false
    this.suggestions = []
  }

  private describeFilters() {
    const parts: string[] = []
    if (this.appliedQuery) parts.push(`query: ${this.appliedQuery}`)
    if (this.selectedTags.length > 0)
      parts.push(`tags: ${this.selectedTags.join(', ')}`)
    return parts.join(' / ')
  }

  static styles = [
    pageStyles,
    css`
      .page-header {
        padding-block: clamp(96px, 18vh, 160px) 40px;
      }

      .page-title {
        margin: 0;
        font-family: var(--font-en);
        font-size: 22px;
        font-weight: 500;
        letter-spacing: var(--tracking-tight);
      }

      .page-description {
        margin: 6px 0 0;
        font-family: var(--font-mono);
        font-size: 12px;
        color: var(--color-text-tertiary);
      }

      .search-area {
        position: relative;
        margin-bottom: 16px;
      }

      .search-input {
        width: 100%;
        padding: 10px 0;
        border: 0;
        border-bottom: 1px solid var(--color-border);
        background: transparent;
        color: var(--color-text-primary);
        font-family: var(--font-mono);
        font-size: 13px;
        caret-color: var(--color-accent);
        transition: border-color 0.2s ease;
      }

      .search-input::placeholder {
        color: var(--color-text-tertiary);
      }

      .search-input:focus {
        outline: none;
        border-bottom-color: var(--color-text-primary);
      }

      .search-status {
        margin-top: 8px;
      }

      .suggestion-list {
        margin-top: 4px;
      }

      .suggestion-item {
        width: 100%;
        display: flex;
        justify-content: space-between;
        gap: 16px;
        padding: 8px 0;
        border: 0;
        background: none;
        color: var(--color-text-primary);
        font: inherit;
        text-align: left;
        cursor: pointer;
      }

      .suggestion-item:hover .suggestion-value {
        color: var(--color-accent);
      }

      .suggestion-meta {
        font-family: var(--font-mono);
        font-size: 12px;
        color: var(--color-text-tertiary);
      }

      .tag-cloud {
        display: flex;
        flex-wrap: wrap;
        gap: 6px 14px;
        margin-bottom: 48px;
        font-family: var(--font-mono);
        font-size: 12px;
      }

      .tag,
      .tag-toggle,
      .article-tag,
      .ghost-button {
        border: 0;
        padding: 0;
        background: none;
        font: inherit;
        color: var(--color-text-tertiary);
        cursor: pointer;
        transition: color 0.2s ease;
      }

      .tag:hover,
      .tag-toggle:hover,
      .article-tag:hover,
      .ghost-button:hover {
        color: var(--color-text-primary);
      }

      .tag.selected {
        color: var(--color-accent);
      }

      .tag-count {
        margin-left: 3px;
        font-size: 10px;
        opacity: 0.6;
      }

      .tag-hash {
        opacity: 0.5;
      }

      .year-group + .year-group {
        margin-top: 48px;
      }

      .year-label {
        margin-bottom: 12px;
        font-family: var(--font-mono);
        font-size: 12px;
        color: var(--color-text-tertiary);
      }

      .article-list > li {
        border-top: 1px solid var(--color-border);
      }

      .article-list > li:last-child {
        border-bottom: 1px solid var(--color-border);
      }

      .article-list:hover .row {
        opacity: 0.4;
      }

      .article-list .article-row:hover .row {
        opacity: 1;
      }

      .article-tags {
        display: flex;
        flex-wrap: wrap;
        gap: 4px 12px;
        margin: -4px 0 11px;
        font-family: var(--font-mono);
        font-size: 11px;
      }

      .pending {
        width: 100%;
        padding: 11px 0;
        border-top: 1px solid var(--color-border);
      }

      .message.error,
      .empty-state p {
        margin: 0 0 12px;
        font-family: var(--font-mono);
        font-size: 12px;
        color: var(--color-text-tertiary);
      }

      .ghost-button {
        font-family: var(--font-mono);
        font-size: 12px;
      }

      .ghost-button:disabled {
        cursor: default;
        opacity: 0.5;
      }

      .load-more {
        margin-top: 24px;
      }
    `,
  ]
}

declare global {
  interface HTMLElementTagNameMap {
    'page-articles': PageArticles
  }
}
