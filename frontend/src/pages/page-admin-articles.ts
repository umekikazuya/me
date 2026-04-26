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
import { articleContext } from '../contexts/article-context.js'
import { RepositoryObserver } from '../controllers/RepositoryObserver.js'
import type { IArticleRepository } from '../domain/ArticleRepository.js'

// Import encapsulated components
import '../components/admin/ui/me-admin-section.js'
import '../components/admin/ui/me-text-input.js'
import '../components/admin/ui/me-textarea.js'
import '../components/admin/ui/me-select.js'

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
  set articleRepo(repo: IArticleRepository) {
    if (this._articleRepo === repo) return
    this._articleRepo = repo
    this._observer?.disconnect()
    if (repo) this._observer = new RepositoryObserver(this, repo)
  }
  get articleRepo() {
    return this._articleRepo
  }
  private _articleRepo!: IArticleRepository
  private _observer?: RepositoryObserver

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
  private baseline: ArticleDraft = createEmptyArticleDraft()

  private onBeforeUnload = (event: BeforeUnloadEvent) => {
    if (!this.articleRepo.adminDirty) return
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
    const ac = this.articleRepo
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

        <me-admin-section
          title="検索と絞り込み"
          description="キーワード、年、プラットフォーム、タグで記事一覧を絞り込みます。"
        >
          <form class="grid" @submit=${this.handleSearch}>
            <me-text-input
              label="キーワード"
              name="q"
              type="search"
              placeholder="タイトルやトークンで検索"
              class="field-wide"
              .value=${this.filters.q}
            ></me-text-input>

            <me-text-input
              label="公開年"
              name="year"
              type="number"
              .value=${this.filters.year}
            ></me-text-input>

            <me-select
              label="プラットフォーム"
              name="platform"
              .value=${this.filters.platform}
            >
              <option value="">すべて</option>
              ${articlePlatforms.map(
                (p) =>
                  html`<option value=${p}>${this.platformLabel(p)}</option>`,
              )}
            </me-select>

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
        </me-admin-section>

        <div class="content-grid">
          <me-admin-section
            title="記事一覧"
            description="一覧から記事を選ぶと右側のフォームで編集できます。"
          >
            <button
              slot="header-actions"
              type="button"
              class="subtle"
              ?disabled=${this.loading || this.loadingMore}
              @click=${this.handleRefreshArticles}
            >
              再読み込み
            </button>

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
                                this.baseline.externalId ===
                                  article.externalId &&
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
          </me-admin-section>

          <me-admin-section
            class="editor-section"
            title=${this.editorMode === 'edit' ? '記事を編集' : '記事を登録'}
            description=${
              this.editorMode === 'edit'
                ? 'manual 登録した記事を更新します。externalId と platform は変更できません。'
                : '管理画面から手動追加する記事を登録します。'
            }
          >
            ${
              this.editorMode === 'edit'
                ? html`
                    <p class="message notice">
                      一覧 API では <code>articleUpdatedAt</code> を取得できないため、必要なら再入力してください。
                    </p>
                  `
                : null
            }

            <form class="stack" @submit=${this.handleSubmit} @input=${this.handleInput}>
              <div class="grid">
                <me-text-input
                  label="externalId *"
                  name="externalId"
                  .value=${this.baseline.externalId}
                  ?readonly=${this.editorMode === 'edit'}
                  required
                ></me-text-input>

                <me-select
                  label="platform *"
                  name="platform"
                  .value=${this.baseline.platform}
                  ?disabled=${this.editorMode === 'edit'}
                  required
                >
                  ${articlePlatforms.map(
                    (p) =>
                      html`<option value=${p}>${this.platformLabel(p)}</option>`,
                  )}
                </me-select>

                <me-text-input
                  label="title *"
                  name="title"
                  class="field-wide"
                  .value=${this.baseline.title}
                  required
                ></me-text-input>

                <me-text-input
                  label="url *"
                  name="url"
                  type="url"
                  class="field-wide"
                  .value=${this.baseline.url}
                  required
                ></me-text-input>

                <me-text-input
                  label="publishedAt"
                  name="publishedAt"
                  type="datetime-local"
                  .value=${this.baseline.publishedAt}
                ></me-text-input>

                <me-text-input
                  label="articleUpdatedAt"
                  name="articleUpdatedAt"
                  type="datetime-local"
                  .value=${this.baseline.articleUpdatedAt}
                ></me-text-input>

                <me-textarea
                  label="tags（1行につき1件）"
                  name="tags"
                  class="field-wide"
                  rows="6"
                  .value=${this.baseline.tags.join('\n')}
                ></me-textarea>
              </div>

              <div class="actions">
                <div class="actions-copy">
                  <p class=${ac.adminDirty ? 'dirty-indicator dirty' : 'dirty-indicator'}>
                    ${ac.adminDirty ? '未保存の変更があります。' : '保存済みの内容です。'}
                  </p>
                </div>
                <button
                  type="reset"
                  class="subtle"
                  ?disabled=${this.saving || this.deleting}
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
          </me-admin-section>
        </div>
      </section>
    `
  }

  private handleInput() {
    this.articleRepo.setAdminDirty(true)
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
    const form = event.target as HTMLFormElement
    const formData = new FormData(form)

    this.filters = {
      ...this.filters,
      q: (formData.get('q') as string) || '',
      year: (formData.get('year') as string) || '',
      platform: (formData.get('platform') as SearchFormState['platform']) || '',
    }

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
    const form = event.target as HTMLFormElement
    if (!form.checkValidity()) {
      form.reportValidity()
      return
    }

    const formData = new FormData(form)
    const draft: ArticleDraft = {
      externalId: formData.get('externalId') as string,
      platform: formData.get('platform') as ArticlePlatform,
      title: formData.get('title') as string,
      url: formData.get('url') as string,
      publishedAt: (formData.get('publishedAt') as string) || '',
      articleUpdatedAt: (formData.get('articleUpdatedAt') as string) || '',
      tags: ((formData.get('tags') as string) || '')
        .split('\n')
        .map((s) => s.trim())
        .filter(Boolean),
    }

    this.saving = true
    this.errorMessage = ''
    this.successMessage = ''

    try {
      if (this.editorMode === 'edit') {
        await updateArticle(draft.externalId, draft)
        this.successMessage = '記事を更新しました。'
      } else {
        await createArticle(draft)
        this.editorMode = 'edit'
        this.successMessage = '記事を登録しました。'
      }

      this.setBaseline(cloneArticleDraft(draft))
      await Promise.all([this.reloadArticles(), this.refreshTags()])
      this.articleRepo.setAdminDirty(false)
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
      await deleteArticle(this.baseline.externalId)
      this.successMessage = '記事を削除しました。'
      this.startCreateMode()
      await Promise.all([this.reloadArticles(), this.refreshTags()])
    } catch (error) {
      this.errorMessage = describeApiError(error)
    } finally {
      this.deleting = false
    }
  }

  private handleReset = (e: Event) => {
    e.preventDefault()
    if (!this.confirmDiscardChanges()) return
    this.articleRepo.setAdminDirty(false)
    this.requestUpdate()
  }

  private startCreateMode() {
    this.editorMode = 'create'
    this.setBaseline(createEmptyArticleDraft())
    this.successMessage = ''
    this.errorMessage = ''
    this.articleRepo.setAdminDirty(false)
  }

  private startEditMode(article: ArticleItem) {
    this.editorMode = 'edit'
    this.setBaseline(articleDraftFromArticle(article))
    this.successMessage = ''
    this.errorMessage = ''
    this.articleRepo.setAdminDirty(false)
  }

  private setBaseline(nextBaseline: ArticleDraft) {
    this.baseline = nextBaseline
  }

  private confirmDiscardChanges() {
    return (
      !this.articleRepo.adminDirty ||
      window.confirm('未保存の変更を破棄して切り替えてもよいですか？')
    )
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

      .article-card {
        display: grid;
        gap: 16px;
        padding: 24px;
        border: 1px solid var(--color-border);
        background: #fff;
      }

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
        cursor: pointer;
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

      .field-wide {
        grid-column: 1 / -1;
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
        z-index: 10;
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
