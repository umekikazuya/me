import type { components } from '@me/types'
import { css, html, nothing } from 'lit'
import { customElement, query, state } from 'lit/decorators.js'
import { listArticles, suggestArticles } from '../api/article-api.js'
import { describeApiError } from '../api/types.js'
import { pageStyles } from '../styles/page-styles.js'
import { formatDate, sanitizeUrl } from '../utils/format.js'
import { ShellPage } from './shell-page.js'
import '../components/load-spinner.js'

type Article = components['schemas']['ArticleItem']

/** 絞り込み条件を git log のオプションとして表す */
export function gitLogCommand(query: string, tags: string[]): string {
  const parts = ['git log --oneline']
  if (query) parts.push(`--grep="${query}"`)
  if (tags.length > 0) parts.push(`-- ${tags.map((t) => `tag:${t}`).join(' ')}`)
  return parts.join(' ')
}

/** 記事 ID の先頭7文字を短縮ハッシュに見立てる */
export function shortHash(externalId: string): string {
  return externalId.slice(0, 7).padEnd(7, ' ')
}

/** `git log --oneline` — 記事を commit として並べる */
@customElement('page-articles')
export class PageArticles extends ShellPage {
  @state()
  private articles: Article[] = []

  @state()
  private suggestions: components['schemas']['ArticleSuggestionItem'][] = []

  @state()
  private query = ''

  @state()
  private appliedQuery = ''

  @state()
  private selectedTags: string[] = []

  @state()
  private loading = true

  @state()
  private loadingMore = false

  @state()
  private suggestionLoading = false

  @state()
  private errorMessage = ''

  @state()
  private nextCursor?: string

  @query('input')
  private searchInput?: HTMLInputElement

  private suggestTimer?: number
  private articleRequestId = 0
  private suggestionRequestId = 0

  // vim の検索と同じく "/" で検索欄へ
  private onSlash = (e: KeyboardEvent) => {
    if (e.key !== '/' || e.metaKey || e.ctrlKey || e.altKey) return
    const target = e.composedPath()[0] as HTMLElement | undefined
    if (target && ['INPUT', 'TEXTAREA'].includes(target.tagName)) return
    e.preventDefault()
    this.searchInput?.focus()
  }

  connectedCallback() {
    super.connectedCallback()
    window.addEventListener('keydown', this.onSlash)
  }

  firstUpdated() {
    void this.reloadArticles()
  }

  disconnectedCallback() {
    super.disconnectedCallback()
    window.removeEventListener('keydown', this.onSlash)
    this.clearSuggestTimer()
  }

  render() {
    return html`
      <h1 class="sr-only">writing</h1>
      <shell-command command="git log --oneline"></shell-command>
      ${this.typed ? this.renderOutput() : nothing}
    `
  }

  private renderOutput() {
    let i = 0
    return html`
      ${this.renderFilters()}
      ${
        this.loading
          ? html`<load-spinner class="ln"></load-spinner>`
          : this.articles.length > 0
            ? html`
              <ul>
                ${this.articles.map((article) => html`<li>${this.renderRow(article, i++)}</li>`)}
              </ul>
              ${this.errorMessage ? html`<p class="ln m" style="--i:${i++}">${this.errorMessage}</p>` : nothing}
              ${this.renderMore(i++)}
            `
            : this.errorMessage
              ? html`<p class="ln m">${this.errorMessage}</p>`
              : html`<button type="button" class="ln m" @click=${this.clearFilters}>no matching commits — reset filters</button>`
      }
      <shell-prompt style="animation-delay:${i * 16}ms"></shell-prompt>
    `
  }

  private renderFilters() {
    const filtered = this.appliedQuery || this.selectedTags.length > 0
    return html`
      <form class="search" @submit=${this.handleSearch}>
        <label class="m" for="q">/</label>
        <input
          id="q"
          type="search"
          autocomplete="off"
          spellcheck="false"
          .value=${this.query}
          placeholder="search title or tag"
          aria-label="記事を検索"
          @input=${this.handleQueryInput}
        />
      </form>
      ${
        this.suggestionLoading
          ? html`<load-spinner class="ln" label="searching"></load-spinner>`
          : this.suggestions.map(
              (s) => html`
                <button type="button" class="ln" @click=${() => this.handleSuggestionSelect(s)}>
                  <span class="m">  ${s.type === 'title' ? 'title' : s.type}</span>  <span class="jp">${s.value}</span>
                </button>
              `,
            )
      }
      ${
        filtered
          ? html`
            <p class="ln m">${gitLogCommand(this.appliedQuery, this.selectedTags)}</p>
            ${this.selectedTags.map(
              (tag) =>
                html`<button type="button" class="ln m" @click=${() => this.toggleTag(tag)}>  tag:${tag}  ×</button>`,
            )}
            <button type="button" class="ln m" @click=${this.clearFilters}>  reset</button>
            <p class="ln"></p>
          `
          : nothing
      }
    `
  }

  private renderRow(article: Article, i: number) {
    const tags = article.tags ?? []
    return html`<a class="ln row" style="--i:${i}" href=${sanitizeUrl(article.url)} target="_blank" rel="noopener noreferrer"><span class="a">${shortHash(article.externalId)}</span> ${
      tags.length > 0
        ? html`<span class="m">(${tags.map((t) => `tag: ${t}`).join(', ')})</span> `
        : nothing
    }<span class="jp t">${article.title}</span> <span class="m">${formatDate(article.publishedAt).replaceAll('.', '-')}</span></a>`
  }

  private renderMore(i: number) {
    if (!this.nextCursor) {
      return html`<p class="ln m" style="--i:${i}">(END)</p>`
    }
    return html`<button type="button" class="ln m" style="--i:${i}" ?disabled=${this.loadingMore} @click=${this.handleLoadMore}>${this.loadingMore ? 'loading…' : ':more'}</button>`
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

  static styles = [
    pageStyles,
    css`
      p {
        margin: 0;
      }

      .search {
        display: flex;
        gap: 1ch;
        min-height: 1.9em;
      }

      input {
        flex: 1;
        min-width: 0;
        padding: 0;
        border: 0;
        background: none;
        color: var(--color-text-primary);
        font: inherit;
        caret-color: var(--color-accent);
      }

      input::placeholder {
        color: var(--color-text-tertiary);
      }

      input:focus {
        outline: none;
      }

      input::-webkit-search-cancel-button {
        display: none;
      }

      .row:hover .t,
      .row:focus-visible .t {
        color: var(--color-text-primary);
      }

      button.ln:disabled {
        cursor: default;
      }
    `,
  ]
}

declare global {
  interface HTMLElementTagNameMap {
    'page-articles': PageArticles
  }
}
