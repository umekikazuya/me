import { consume } from '@lit/context'
import { css, html, LitElement } from 'lit'
import { customElement, state } from 'lit/decorators.js'
import { adminFormStyles } from '../admin/admin-form-styles.js'
import {
  createArticle,
  deleteArticle,
  listArticles,
  listArticleTags,
  updateArticle,
} from '../admin/article-api.js'
import {
  type ArticleDraft,
  type ArticleItem,
  type ArticlePlatform,
  type ArticleTagItem,
  articleDraftFromArticle,
  articlePlatforms,
  cloneArticleDraft,
  createEmptyArticleDraft,
} from '../admin/article-types.js'
import { describeApiError } from '../admin/types.js'
import {
  articleContext,
  type ArticleController,
} from '../contexts/article-context.js'

interface SearchFormState {
  q: string
  year: string
  platform: '' | ArticlePlatform
  tags: string[]
}

const createSearchFormState = (): SearchFormState => ({
  q: '',
  year: '',
  platform: '',
  tags: [],
})

@customElement('page-admin-articles')
export class PageAdminArticles extends LitElement {
  @consume({ context: articleContext, subscribe: true })
  articleController!: ArticleController

  @state()
  private articles: ArticleItem[] = []

  @state()
  private tagOptions: ArticleTagItem[] = []

  @state()
  private filters: SearchFormState = createSearchFormState()

  @state()
  private loading = false

  @state()
  private loadingMore = false

  @state()
  private saving = false

  @state()
  private deleting = false

  @state()
  private errorMessage = ''

  @state()
  private successMessage = ''

  @state()
  private nextCursor?: string

  @state()
  private editorMode: 'create' | 'edit' = 'create'

  @state()
  private form: ArticleDraft = createEmptyArticleDraft()

  @state()
  private baseline: ArticleDraft = createEmptyArticleDraft()

  private onBeforeUnload = (event: BeforeUnloadEvent) => {
    if (!this.articleController.adminDirty) return
    event.preventDefault()
    event.returnValue = ''
  }

  connectedCallback() {
    super.connectedCallback()
    window.addEventListener('beforeunload', this.onBeforeUnload)
  }

  disconnectedCallback() {
    super.disconnectedCallback()
    window.removeEventListener('beforeunload', this.onBeforeUnload)
  }

  firstUpdated() {
    void this.loadInitialData()
  }

  render() {
    const ac = this.articleController
    return html`
      <section class="container">
        <header class="page-header">
          <div>
            <p class="eyebrow" lang="en">Articles</p>
            <h1 class="title">記事管理</h1>
            <p class="description">
              記事の一覧確認、手動登録、更新、削除を行います。
            </p>
            <div class="meta">
              <span>Loaded: ${this.articles.length}</span>
              <span>Known tags: ${this.tagOptions.length}</span>
              <span>${this.editorMode === 'edit' ? 'Editing' : 'Creating'}</span>
            </div>
          </div>
          <button
            type="button"
            class="subtle"
            ?disabled=${this.saving || this.deleting}
            @click=${this.handleStartCreate}
          >
            新規作成
          </button>
        </header>

        ${this.errorMessage ? html`<p class="message error">${this.errorMessage}</p>` : null}
        ${
          this.successMessage
            ? html`<p class="message success">${this.successMessage}</p>`
            : null
        }

        <section class="section">
          <div class="section-header">
            <div class="section-copy">
              <h2>検索と絞り込み</h2>
              <p class="section-help">
                キーワード、年、プラットフォーム、タグで記事一覧を絞り込みます。
              </p>
            </div>
          </div>

          <form class="grid" @submit=${this.handleSearch}>
            <label class="field field-wide">
              <span>キーワード</span>
              <input
                type="search"
                .value=${this.filters.q}
                placeholder="タイトルやトークンで検索"
                @input=${(event: Event) => {
                  this.filters = {
                    ...this.filters,
                    q: (event.target as HTMLInputElement).value,
                  }
                }}
              />
            </label>

            <label class="field">
              <span>公開年</span>
              <input
                type="number"
                min="1"
                max="2100"
                .value=${this.filters.year}
                @input=${(event: Event) => {
                  this.filters = {
                    ...this.filters,
                    year: (event.target as HTMLInputElement).value,
                  }
                }}
              />
            </label>

            <label class="field">
              <span>プラットフォーム</span>
              <select
                .value=${this.filters.platform}
                @change=${(event: Event) => {
                  this.filters = {
                    ...this.filters,
                    platform: (event.target as HTMLSelectElement)
                      .value as SearchFormState['platform'],
                  }
                }}
              >
                <option value="">すべて</option>
                ${articlePlatforms.map(
                  (platform) => html`
                    <option value=${platform}>${this.platformLabel(platform)}</option>
                  `,
                )}
              </select>
            </label>

            <div class="filter-actions field-wide">
              <button type="submit" ?disabled=${this.loading || this.loadingMore}>
                ${this.loading ? '読み込み中...' : '絞り込む'}
              </button>
              <button
                type="button"
                class="subtle"
                ?disabled=${this.loading || this.loadingMore}
                @click=${this.handleClearFilters}
              >
                条件をリセット
              </button>
            </div>
          </form>

          ${
            this.tagOptions.length > 0
              ? html`
                <div class="tag-filter">
                  <p class="tag-filter-label">タグ</p>
                  <div class="tag-list">
                    ${this.tagOptions.map(
                      (tag) => html`
                        <button
                          type="button"
                          class=${
                            this.filters.tags.includes(tag.name)
                              ? 'tag-chip selected'
                              : 'tag-chip'
                          }
                          aria-pressed=${this.filters.tags.includes(tag.name)}
                          @click=${() => this.toggleFilterTag(tag.name)}
                        >
                          <span>${tag.name}</span>
                          <small>${tag.count}</small>
                        </button>
                      `,
                    )}
                  </div>
                </div>
              `
              : null
          }
        </section>

        <div class="content-grid">
          <section class="section">
            <div class="section-header">
              <div class="section-copy">
                <h2>記事一覧</h2>
                <p class="section-help">
                  一覧から記事を選ぶと右側のフォームで編集できます。
                </p>
              </div>
              <button
                type="button"
                class="subtle"
                ?disabled=${this.loading || this.loadingMore}
                @click=${this.handleRefreshArticles}
              >
                再読み込み
              </button>
            </div>

            ${
              this.loading
                ? html`<p class="loading">記事を読み込み中...</p>`
                : this.articles.length === 0
                  ? html`
                      <article class="empty-panel">
                        <p>条件に一致する記事がありません。</p>
                      </article>
                    `
                  : html`
                      <div class="stack">
                        ${this.articles.map(
                          (article) => html`
                            <article
                              class=${
                                this.form.externalId === article.externalId &&
                                this.editorMode === 'edit'
                                  ? 'article-card selected'
                                  : 'article-card'
                              }
                            >
                              <div class="article-card-header">
                                <div class="article-copy">
                                  <p class="article-platform">
                                    ${this.platformLabel(article.platform)}
                                  </p>
                                  <h3>${article.title}</h3>
                                  <a href=${article.url} target="_blank" rel="noreferrer">
                                    ${article.url}
                                  </a>
                                </div>
                                <button
                                  type="button"
                                  class="subtle"
                                  @click=${() => this.handleStartEdit(article)}
                                >
                                  編集
                                </button>
                              </div>
                              <div class="article-meta">
                                <span>ID: ${article.externalId}</span>
                                <span>
                                  公開日:
                                  ${
                                    article.publishedAt
                                      ? this.formatDateTime(article.publishedAt)
                                      : '未設定'
                                  }
                                </span>
                              </div>
                              <div class="article-tags">
                                ${
                                  article.tags.length > 0
                                    ? article.tags.map(
                                        (tag) => html`
                                        <button
                                          type="button"
                                          class="inline-tag"
                                          @click=${() => this.toggleFilterTag(tag)}
                                        >
                                          ${tag}
                                        </button>
                                      `,
                                      )
                                    : html`<span class="muted">タグなし</span>`
                                }
                              </div>
                            </article>
                          `,
                        )}
                      </div>
                    `
            }

            ${
              this.nextCursor
                ? html`
                    <button
                      type="button"
                      class="subtle"
                      ?disabled=${this.loadingMore}
                      @click=${this.handleLoadMore}
                    >
                      ${this.loadingMore ? '読み込み中...' : 'さらに読み込む'}
                    </button>
                  `
                : null
            }
          </section>

          <section class="section editor-section">
            <div class="section-header">
              <div class="section-copy">
                <h2>${this.editorMode === 'edit' ? '記事を編集' : '記事を登録'}</h2>
                <p class="section-help">
                  ${
                    this.editorMode === 'edit'
                      ? 'manual 登録した記事を更新します。externalId と platform は変更できません。'
                      : '管理画面から手動追加する記事を登録します。'
                  }
                </p>
              </div>
            </div>

            ${
              this.editorMode === 'edit'
                ? html`
                    <p class="message notice">
                      一覧 API では <code>articleUpdatedAt</code> を取得できないため、必要なら再入力してください。
                      空のまま保存するとクリアされます。
                    </p>
                  `
                : null
            }

            <form class="stack" @submit=${this.handleSubmit}>
              <div class="grid">
                <label class="field">
                  <span>externalId *</span>
                  <input
                    .value=${this.form.externalId}
                    ?readonly=${this.editorMode === 'edit'}
                    @input=${(event: Event) =>
                      this.updateForm(
                        'externalId',
                        (event.target as HTMLInputElement).value,
                      )}
                    required
                  />
                </label>

                <label class="field">
                  <span>platform *</span>
                  <select
                    .value=${this.form.platform}
                    ?disabled=${this.editorMode === 'edit'}
                    @change=${(event: Event) =>
                      this.updateForm(
                        'platform',
                        (event.target as HTMLSelectElement)
                          .value as ArticlePlatform,
                      )}
                  >
                    ${articlePlatforms.map(
                      (platform) => html`
                        <option value=${platform}>${this.platformLabel(platform)}</option>
                      `,
                    )}
                  </select>
                </label>

                <label class="field field-wide">
                  <span>title *</span>
                  <input
                    .value=${this.form.title}
                    @input=${(event: Event) =>
                      this.updateForm(
                        'title',
                        (event.target as HTMLInputElement).value,
                      )}
                    required
                  />
                </label>

                <label class="field field-wide">
                  <span>url *</span>
                  <input
                    type="url"
                    .value=${this.form.url}
                    @input=${(event: Event) =>
                      this.updateForm(
                        'url',
                        (event.target as HTMLInputElement).value,
                      )}
                    required
                  />
                </label>

                <label class="field">
                  <span>publishedAt</span>
                  <input
                    type="datetime-local"
                    .value=${this.form.publishedAt}
                    @input=${(event: Event) =>
                      this.updateForm(
                        'publishedAt',
                        (event.target as HTMLInputElement).value,
                      )}
                  />
                </label>

                <label class="field">
                  <span>articleUpdatedAt</span>
                  <input
                    type="datetime-local"
                    .value=${this.form.articleUpdatedAt}
                    @input=${(event: Event) =>
                      this.updateForm(
                        'articleUpdatedAt',
                        (event.target as HTMLInputElement).value,
                      )}
                  />
                </label>

                <label class="field field-wide">
                  <span>tags（1行につき1件）</span>
                  <textarea
                    rows="6"
                    .value=${this.form.tags.join('\n')}
                    @input=${(event: Event) =>
                      this.updateForm(
                        'tags',
                        this.splitLines(
                          (event.target as HTMLTextAreaElement).value,
                        ),
                      )}
                  ></textarea>
                </label>
              </div>

              <div class="actions">
                <div class="actions-copy">
                  <p class=${ac.adminDirty ? 'dirty-indicator dirty' : 'dirty-indicator'}>
                    ${ac.adminDirty ? '未保存の変更があります。' : '保存済みの内容です。'}
                  </p>
                </div>
                <button
                  type="button"
                  class="subtle"
                  ?disabled=${!ac.adminDirty || this.saving || this.deleting}
                  @click=${this.handleReset}
                >
                  入力を戻す
                </button>
                ${
                  this.editorMode === 'edit'
                    ? html`
                        <button
                          type="button"
                          class="danger"
                          ?disabled=${this.saving || this.deleting}
                          @click=${this.handleDelete}
                        >
                          ${this.deleting ? '削除中...' : '削除'}
                        </button>
                      `
                    : null
                }
                <button
                  type="submit"
                  ?disabled=${this.saving || this.deleting || !ac.adminDirty}
                >
                  ${
                    this.saving
                      ? this.editorMode === 'edit'
                        ? '更新中...'
                        : '登録中...'
                      : this.editorMode === 'edit'
                        ? '更新する'
                        : '登録する'
                  }
                </button>
              </div>
            </form>
          </section>
        </div>
      </section>
    `
  }

  private async loadInitialData() {
    this.loading = true
    this.errorMessage = ''

    const [articlesResult, tagsResult] = await Promise.allSettled([
      listArticles({ limit: 50 }),
      listArticleTags(),
    ])

    if (articlesResult.status === 'fulfilled') {
      this.articles = articlesResult.value.articles
      this.nextCursor = articlesResult.value.nextCursor
    } else {
      this.errorMessage = describeApiError(articlesResult.reason)
    }

    if (tagsResult.status === 'fulfilled') {
      this.tagOptions = tagsResult.value
    } else if (!this.errorMessage) {
      this.errorMessage = describeApiError(tagsResult.reason)
    }

    this.loading = false
  }

  private async reloadArticles(cursor?: string, append = false) {
    if (append) {
      this.loadingMore = true
    } else {
      this.loading = true
      this.errorMessage = ''
    }

    try {
      const result = await listArticles({
        q: this.filters.q.trim() || undefined,
        year: this.toOptionalNumber(this.filters.year),
        platform: this.filters.platform || undefined,
        tag: this.filters.tags,
        limit: 50,
        cursor,
      })

      this.articles = append
        ? [...this.articles, ...result.articles]
        : result.articles
      this.nextCursor = result.nextCursor
    } catch (error) {
      this.errorMessage = describeApiError(error)
    } finally {
      this.loading = false
      this.loadingMore = false
    }
  }

  private async refreshTags() {
    try {
      this.tagOptions = await listArticleTags()
    } catch (error) {
      if (!this.errorMessage) {
        this.errorMessage = describeApiError(error)
      }
    }
  }

  private handleSearch(event: Event) {
    event.preventDefault()
    void this.reloadArticles()
  }

  private handleRefreshArticles = () => {
    void this.reloadArticles()
  }

  private handleLoadMore = () => {
    if (!this.nextCursor) return
    void this.reloadArticles(this.nextCursor, true)
  }

  private handleClearFilters = () => {
    this.filters = createSearchFormState()
    void this.reloadArticles()
  }

  private toggleFilterTag(tagName: string) {
    const nextTags = this.filters.tags.includes(tagName)
      ? this.filters.tags.filter((tag) => tag !== tagName)
      : [...this.filters.tags, tagName]

    this.filters = {
      ...this.filters,
      tags: nextTags,
    }

    void this.reloadArticles()
  }

  private handleStartCreate = () => {
    if (!this.confirmDiscardChanges()) return
    this.startCreateMode()
  }

  private handleStartEdit(article: ArticleItem) {
    if (!this.confirmDiscardChanges()) return
    this.startEditMode(article)
  }

  private async handleSubmit(event: Event) {
    event.preventDefault()
    this.saving = true
    this.errorMessage = ''
    this.successMessage = ''

    try {
      if (this.editorMode === 'edit') {
        await updateArticle(this.form.externalId, this.form)
        this.successMessage = '記事を更新しました。'
      } else {
        await createArticle(this.form)
        this.editorMode = 'edit'
        this.successMessage = '記事を登録しました。'
      }

      this.setBaseline(cloneArticleDraft(this.form))
      await Promise.all([this.reloadArticles(), this.refreshTags()])
    } catch (error) {
      this.errorMessage = describeApiError(error)
    } finally {
      this.saving = false
    }
  }

  private async handleDelete() {
    if (this.editorMode !== 'edit') return
    if (!window.confirm('この記事を削除します。よろしいですか？')) return

    this.deleting = true
    this.errorMessage = ''
    this.successMessage = ''

    try {
      await deleteArticle(this.form.externalId)
      this.successMessage = '記事を削除しました。'
      this.startCreateMode()
      await Promise.all([this.reloadArticles(), this.refreshTags()])
    } catch (error) {
      this.errorMessage = describeApiError(error)
    } finally {
      this.deleting = false
    }
  }

  private handleReset = () => {
    if (!this.confirmDiscardChanges()) return
    this.setForm(cloneArticleDraft(this.baseline))
  }

  private startCreateMode() {
    this.editorMode = 'create'
    this.setBaseline(createEmptyArticleDraft())
    this.successMessage = ''
    this.errorMessage = ''
  }

  private startEditMode(article: ArticleItem) {
    this.editorMode = 'edit'
    this.setBaseline(articleDraftFromArticle(article))
    this.successMessage = ''
    this.errorMessage = ''
  }

  private updateForm<Key extends keyof ArticleDraft>(
    key: Key,
    value: ArticleDraft[Key],
  ) {
    this.setForm({
      ...this.form,
      [key]: value,
    })
  }

  private setBaseline(nextBaseline: ArticleDraft) {
    this.baseline = nextBaseline
    this.setForm(cloneArticleDraft(nextBaseline))
  }

  private setForm(nextForm: ArticleDraft) {
    this.form = nextForm
    this.updateDirtyState(nextForm)
  }

  private updateDirtyState(nextForm: ArticleDraft) {
    const nextDirty = JSON.stringify(nextForm) !== JSON.stringify(this.baseline)
    this.articleController.setAdminDirty(nextDirty)
  }

  private confirmDiscardChanges() {
    return (
      !this.articleController.adminDirty ||
      window.confirm('未保存の変更を破棄して切り替えてもよいですか？')
    )
  }

  private splitLines(value: string) {
    return value
      .split('\n')
      .map((item) => item.trim())
      .filter(Boolean)
  }

  private toOptionalNumber(value: string) {
    const trimmed = value.trim()
    return trimmed === '' ? undefined : Number(trimmed)
  }

  private formatDateTime(value: string) {
    return new Date(value).toLocaleString('ja-JP')
  }

  private platformLabel(platform: ArticlePlatform) {
    switch (platform) {
      case 'qiita':
        return 'Qiita'
      case 'zenn':
        return 'Zenn'
      case 'mochiya':
        return 'Mochiya'
      case 'note':
        return 'note'
    }
  }

  static styles = [
    adminFormStyles,
    css`
      :host {
        display: block;
      }

      code {
        font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
      }

      .container {
        display: grid;
        gap: 24px;
      }

      .page-header,
      .section-header,
      .article-card-header,
      .actions {
        display: flex;
        gap: 16px;
        justify-content: space-between;
        align-items: center;
        flex-wrap: wrap;
      }

      .meta,
      .article-meta,
      .article-tags,
      .tag-list,
      .filter-actions {
        display: flex;
        gap: 8px;
        flex-wrap: wrap;
      }

      .meta span,
      .article-meta span {
        display: inline-flex;
        align-items: center;
        min-height: 28px;
        padding: 0 10px;
        border: 1px solid var(--color-border);
        background: var(--color-bg-surface);
        color: var(--color-text-secondary);
        font-size: 12px;
      }

      .section,
      .article-card {
        display: grid;
        gap: 16px;
        padding: 24px;
        border: 1px solid var(--color-border);
        background: #fff;
      }

      .section-copy {
        display: grid;
        gap: 6px;
      }

      .section-help,
      .loading,
      .muted {
        color: var(--color-text-tertiary);
        font-size: 13px;
        line-height: 1.8;
      }

      .tag-filter {
        display: grid;
        gap: 12px;
      }

      .tag-filter-label,
      .article-platform {
        font-family: var(--font-en);
        font-size: 12px;
        letter-spacing: var(--tracking-wide);
        color: var(--color-text-tertiary);
      }

      .tag-chip,
      .inline-tag {
        border: 1px solid var(--color-border);
        background: var(--color-bg-surface);
        color: var(--color-text-secondary);
        padding: 6px 10px;
        font-size: 12px;
      }

      .tag-chip.selected {
        border-color: var(--admin-accent);
        color: var(--admin-accent);
        background: #e8f0fb;
      }

      .tag-chip small {
        font-size: 11px;
        color: inherit;
      }

      .content-grid {
        display: grid;
        gap: 24px;
        grid-template-columns: minmax(0, 1.2fr) minmax(320px, 0.9fr);
        align-items: start;
      }

      .stack,
      form {
        display: grid;
        gap: 16px;
      }

      .grid {
        display: grid;
        gap: 16px;
        grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
      }

      .empty-panel {
        display: grid;
        gap: 12px;
        border: 1px dashed var(--color-border);
        padding: 18px;
        color: var(--color-text-secondary);
        font-size: 14px;
      }

      .article-card.selected {
        border-color: var(--admin-accent);
        box-shadow: 0 0 0 1px color-mix(in srgb, var(--admin-accent) 30%, transparent);
      }

      .article-copy {
        display: grid;
        gap: 6px;
      }

      .article-copy h3 {
        font-size: 16px;
        font-weight: 500;
        color: var(--color-text-primary);
      }

      .article-copy a {
        font-size: 13px;
        color: var(--admin-accent);
        overflow-wrap: anywhere;
      }

      .editor-section {
        position: sticky;
        top: 24px;
      }

      .actions {
        position: sticky;
        bottom: 16px;
        padding: 14px 16px;
        border: 1px solid var(--color-border);
        background: rgba(255, 255, 255, 0.95);
        backdrop-filter: blur(12px);
      }

      .actions-copy {
        margin-right: auto;
      }

      .dirty-indicator {
        font-size: 13px;
        color: var(--color-text-tertiary);
      }

      .dirty-indicator.dirty {
        color: #9a6d2f;
      }

      h2 {
        font-size: 16px;
        font-weight: 500;
        color: var(--color-text-primary);
      }

      @media (max-width: 1080px) {
        .content-grid {
          grid-template-columns: 1fr;
        }

        .editor-section {
          position: static;
        }
      }

      @media (max-width: 720px) {
        .actions {
          flex-wrap: wrap;
        }

        .actions-copy {
          width: 100%;
          margin-right: 0;
        }
      }
    `,
  ]
}

declare global {
  interface HTMLElementTagNameMap {
    'page-admin-articles': PageAdminArticles
  }
}
