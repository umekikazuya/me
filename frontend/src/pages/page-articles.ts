import { css, html, LitElement } from 'lit'
import { customElement, state } from 'lit/decorators.js'
import {
  listArticles,
  listArticleTags,
  suggestArticles,
} from '../admin/article-api.js'
import type {
  ArticleItem,
  ArticleSuggestionItem,
  ArticleTagItem,
} from '../admin/article-types.js'
import { describeApiError } from '../admin/types.js'
import { setupReveal } from '../utils/scroll.js'

interface ArticleGroup {
  key: string
  label: string
  items: ArticleItem[]
}

@customElement('page-articles')
export class PageArticles extends LitElement {
  @state()
  private articles: ArticleItem[] = []

  @state()
  private tagOptions: ArticleTagItem[] = []

  @state()
  private suggestions: ArticleSuggestionItem[] = []

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

  private cleanups: Array<() => void> = []
  private suggestTimer?: number
  private articleRequestId = 0
  private tagRequestId = 0
  private suggestionRequestId = 0

  firstUpdated() {
    const root = this.shadowRoot
    if (!root) return

    const revealEls = Array.from(
      root.querySelectorAll('.page-header, .search-area, .tag-cloud'),
    )
    this.cleanups.push(setupReveal(revealEls, true))
    void this.loadInitialData()
  }

  disconnectedCallback() {
    super.disconnectedCallback()
    if (this.suggestTimer !== undefined) {
      window.clearTimeout(this.suggestTimer)
    }
    for (const cleanup of this.cleanups) cleanup()
    this.cleanups = []
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
          ${this.loading ? html`<p class="loading">記事を読み込み中...</p>` : this.renderArticleGroups()}
        </div>

        ${this.renderLoadMore()}
      </div>
    `
  }

  private renderHeader() {
    return html`
      <header class="page-header">
        <h1 class="page-title">Articles</h1>
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
            placeholder="Search by token or title..."
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
      return html`<p class="search-status">候補を探しています...</p>`
    }
    if (this.suggestions.length === 0) return null

    return html`
      <ul class="suggestion-list">
        ${this.suggestions.map(
          (s) => html`
            <li>
              <button type="button" class="suggestion-item" @click=${() => this.handleSuggestionSelect(s)}>
                <span class="suggestion-value">${s.value}</span>
                <span class="suggestion-meta">${s.type} · ${s.count}</span>
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

  private renderTag(tag: ArticleTagItem) {
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

  private renderArticleRow(article: ArticleItem) {
    return html`
      <li class="article-row">
        <span class="article-date">${this.formatArticleDate(article.publishedAt)}</span>
        <a href=${article.url} class="article-title" target="_blank" rel="noreferrer">${article.title}</a>
        <div class="article-tags">
          ${article.tags.map(
            (tag) => html`
            <button type="button" class="article-tag" @click=${() => this.toggleTag(tag)}>${tag}</button>
          `,
          )}
        </div>
      </li>
    `
  }

  private renderLoadMore() {
    if (!this.nextCursor || this.loading) return null
    return html`
      <div class="load-more">
        <button type="button" class="ghost-button" ?disabled=${this.loadingMore} @click=${this.handleLoadMore}>
          ${this.loadingMore ? 'Loading...' : 'Load more'}
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
        this.tagOptions = tagsResult.value
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

  private handleSuggestionSelect(suggestion: ArticleSuggestionItem) {
    if (suggestion.type === 'tag') {
      this.invalidateSuggestions()
      this.query = ''
      this.appliedQuery = ''
      this.toggleTag(suggestion.value)
      return
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

      this.suggestions = suggestions.slice(0, 10)
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

  private formatArticleDate(value?: string) {
    if (!value) return '----.--'

    const date = new Date(value)
    if (Number.isNaN(date.valueOf())) return '----.--'

    return `${date.getFullYear()}.${String(date.getMonth() + 1).padStart(2, '0')}`
  }

  static styles = css`
    :host {
      display: block;
      padding-top: 80px;
    }

    *, *::before, *::after {
      box-sizing: border-box;
    }

    .container {
      max-width: 720px;
      margin: 0 auto;
      padding: var(--space-lg) var(--space-md);
    }

    .page-header {
      margin-bottom: 32px;
    }

    .page-title {
      font-family: var(--font-en);
      font-weight: 300;
      font-size: 36px;
      letter-spacing: var(--tracking-wide);
      color: var(--color-text-primary);
      margin: 0;
    }

    .page-description,
    .search-status,
    .loading,
    .message,
    .empty-state {
      color: var(--color-text-secondary);
      line-height: 1.8;
      font-size: 14px;
    }

    .search-area {
      position: relative;
      margin-bottom: 24px;
    }

    form {
      margin: 0;
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
      color: var(--color-text-mute);
    }

    .search-input:focus-visible {
      outline: none;
      border-bottom-color: var(--color-text-primary);
      box-shadow: 0 1px 0 0 var(--color-text-primary);
    }

    .suggestion-list {
      list-style: none;
      padding: 8px 0 0;
      margin: 0;
      display: grid;
      gap: 4px;
    }

    .suggestion-item,
    .tag,
    .article-tag,
    .ghost-button {
      border: 0;
      background: transparent;
      padding: 0;
      font: inherit;
      cursor: pointer;
    }

    .suggestion-item {
      width: 100%;
      display: flex;
      justify-content: space-between;
      align-items: center;
      gap: 12px;
      padding: 10px 0;
      border-bottom: 1px solid var(--color-border-subtle);
      text-align: left;
    }

    .suggestion-value {
      color: var(--color-text-primary);
      font-size: 14px;
    }

    .suggestion-meta {
      color: var(--color-text-tertiary);
      font-family: var(--font-en);
      font-size: 12px;
      letter-spacing: var(--tracking-wide);
    }

    .tag-cloud {
      display: flex;
      flex-wrap: wrap;
      column-gap: 20px;
      row-gap: 12px;
      margin-bottom: 64px;
    }

    .tag,
    .tag-toggle {
      display: inline-flex;
      align-items: center;
      gap: 4px;
      padding: 4px 0;
      color: var(--color-text-tertiary);
      font-family: var(--font-en);
      font-size: 13px;
      letter-spacing: var(--tracking-wide);
      transition: color 0.3s ease, text-shadow 0.3s ease, opacity 0.3s ease;
      background: transparent;
      border: none;
      cursor: pointer;
    }

    .tag-hash {
      font-size: 11px;
      color: var(--color-text-secondary);
    }

    .tag-name {
      color: var(--color-text-secondary);
      transition: color 0.3s ease;
    }

    .tag:hover .tag-name {
      color: var(--color-text-primary);
    }

    .tag.selected {
      color: var(--color-text-primary);
      text-shadow: 0 0 12px var(--color-glow-sharp);
    }

    .tag.selected .tag-name {
      color: var(--color-text-primary);
    }

    .tag.selected .tag-hash {
      color: var(--color-text-primary);
    }

    .tag-count {
      font-size: 11px;
      margin-left: 2px;
      font-style: italic;
      color: var(--color-text-primary);
    }

    .tag-toggle {
      color: var(--color-text-tertiary);
      font-style: italic;
      opacity: 0.6;
    }

    .tag-toggle:hover {
      opacity: 1;
      color: var(--color-text-secondary);
    }

    .message.error {
      margin-bottom: 24px;
      color: #8c5a52;
    }

    .year-group {
      margin-bottom: 48px;
    }

    .year-label {
      font-family: var(--font-en);
      font-weight: 300;
      font-size: 13px;
      letter-spacing: var(--tracking-wider);
      color: var(--color-text-primary);
      margin-bottom: 16px;
    }

    .article-list {
      list-style: none;
      padding: 0;
      margin: 0;
    }

    .article-row {
      display: grid;
      grid-template-columns: 88px 1fr auto;
      align-items: baseline;
      gap: 16px;
      padding: 16px 8px;
      border-bottom: 1px solid var(--color-border-subtle);
      transition: background 0.2s ease, transform 0.2s ease;
    }

    .article-row:hover {
      background: var(--color-bg-surface);
      transform: translateX(4px);
      border-bottom-color: var(--color-text-tertiary);
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
      text-decoration: none;
      transition: opacity 0.2s ease;
    }

    .article-title:hover {
      opacity: 0.6;
    }

    .article-tags {
      display: flex;
      gap: 4px;
      flex-wrap: wrap;
      justify-content: flex-end;
    }

    .article-tag {
      font-family: var(--font-en);
      font-weight: 300;
      font-size: 11px;
      letter-spacing: var(--tracking-wide);
      color: var(--color-text-tertiary);
      transition: opacity 0.2s ease;
    }

    .article-tag:hover,
    .ghost-button:hover,
    .tag:hover {
      opacity: 0.6;
    }

    .empty-state,
    .load-more {
      display: grid;
      justify-items: start;
      gap: 12px;
      margin-top: 16px;
    }

    .ghost-button {
      color: var(--color-text-primary);
      font-family: var(--font-en);
      font-size: 13px;
      letter-spacing: var(--tracking-wide);
    }

    @media (prefers-reduced-motion: reduce) {
      .tag,
      .article-title,
      .article-row,
      .article-tag,
      .ghost-button {
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
        justify-content: start;
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
