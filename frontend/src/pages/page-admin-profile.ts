import { consume } from '@lit/context'
import { css, html, LitElement } from 'lit'
import { customElement, state } from 'lit/decorators.js'
import { adminFormStyles } from '../admin/admin-form-styles.js'
import {
  cloneMeProfile,
  createEmptyMeProfile,
  type MeProfile,
} from '../admin/types.js'
import { profileContext } from '../contexts/profile-context.js'
import { RepositoryObserver } from '../controllers/RepositoryObserver.js'
import type { IProfileRepository } from '../domain/ProfileRepository.js'

// Import encapsulated components
import '../components/admin/ui/me-admin-section.js'
import '../components/admin/ui/me-text-input.js'
import '../components/admin/ui/me-textarea.js'
import '../components/admin/profile/me-profile-skills-editor.js'
import '../components/admin/profile/me-profile-certifications-editor.js'
import '../components/admin/profile/me-profile-experiences-editor.js'
import '../components/admin/profile/me-profile-links-editor.js'

@customElement('page-admin-profile')
export class PageAdminProfile extends LitElement {
  @consume({ context: profileContext, subscribe: true })
  set profileRepo(repo: IProfileRepository) {
    this._profileRepo = repo
    if (repo) new RepositoryObserver(this, repo)
  }
  get profileRepo() {
    return this._profileRepo
  }
  private _profileRepo!: IProfileRepository

  @state()
  private form: MeProfile = createEmptyMeProfile()

  private _lastSyncedData = ''

  private onBeforeUnload = (event: BeforeUnloadEvent) => {
    if (!this.profileRepo.adminDirty) return
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

  protected willUpdate() {
    const p = this.profileRepo
    if (!p) return

    // Initialize/Sync form when data is loaded and not currently being edited
    if (p.adminLoaded && !p.adminDirty) {
      const data = JSON.stringify(p.adminProfile)
      if (data !== this._lastSyncedData) {
        this._lastSyncedData = data
        this.setForm(cloneMeProfile(p.adminProfile))
      }
    }
  }

  render() {
    const p = this.profileRepo
    return html`
      <section class="container">
        <header class="page-header">
          <div>
            <p class="eyebrow" lang="en">Profile</p>
            <h1 class="title">プロフィール編集</h1>
            <p class="description">公開プロフィールの表示内容を更新します。</p>
            <div class="meta">
              <span>Skill: ${this.form.skills.length}</span>
              <span>Cert: ${this.form.certifications.length}</span>
              <span>Exp: ${this.form.experiences.length}</span>
              <span>Link: ${this.form.links.length}</span>
            </div>
          </div>
          ${this.form.updatedAt
            ? html`<p class="updated-at">
                最終更新: ${new Date(this.form.updatedAt).toLocaleString('ja-JP')}
              </p>`
            : null}
        </header>

        ${p.adminError ? html`<p class="message error">${p.adminError}</p>` : null}
        ${p.adminSuccess
          ? html`<p class="message success">${p.adminSuccess}</p>`
          : null}

        ${p.adminLoading
          ? html`<p class="loading">プロフィールを読み込み中...</p>`
          : html`
              <form @submit=${this.handleSubmit}>
                <me-admin-section
                  title="基本情報"
                  description="最低限、表示名だけあれば更新できます。未入力項目は公開画面で省略されます。"
                >
                  <div class="grid">
                    <me-text-input
                      label="表示名 *"
                      .value=${this.form.displayName}
                      required
                      @change=${(e: CustomEvent) => this.updateField('displayName', e.detail)}
                    ></me-text-input>

                    <me-text-input
                      label="表示名（日本語）"
                      .value=${this.form.displayJa}
                      @change=${(e: CustomEvent) => this.updateField('displayJa', e.detail)}
                    ></me-text-input>

                    <me-text-input
                      label="Role"
                      .value=${this.form.role}
                      @change=${(e: CustomEvent) => this.updateField('role', e.detail)}
                    ></me-text-input>

                    <me-text-input
                      label="Location"
                      .value=${this.form.location}
                      @change=${(e: CustomEvent) => this.updateField('location', e.detail)}
                    ></me-text-input>
                  </div>
                </me-admin-section>

                <me-profile-skills-editor
                  .skills=${this.form.skills}
                  @change=${(e: CustomEvent) => this.updateField('skills', e.detail)}
                ></me-profile-skills-editor>

                <me-profile-certifications-editor
                  .certifications=${this.form.certifications}
                  @change=${(e: CustomEvent) =>
                    this.updateField('certifications', e.detail)}
                ></me-profile-certifications-editor>

                <me-profile-experiences-editor
                  .experiences=${this.form.experiences}
                  @change=${(e: CustomEvent) =>
                    this.updateField('experiences', e.detail)}
                ></me-profile-experiences-editor>

                <me-profile-links-editor
                  .links=${this.form.links}
                  @change=${(e: CustomEvent) => this.updateField('links', e.detail)}
                ></me-profile-links-editor>

                <me-admin-section
                  title="Likes"
                  description="1行ごとに1件ずつ入力します。空行は保存時に除外されます。"
                >
                  <me-textarea
                    label="1行につき1件"
                    rows="6"
                    .value=${this.form.likes.join('\n')}
                    @change=${(e: CustomEvent) =>
                      this.updateField('likes', this.splitLines(e.detail))}
                  ></me-textarea>
                </me-admin-section>

                <div class="actions">
                  <div class="actions-copy">
                    <p class=${p.adminDirty ? 'dirty-indicator dirty' : 'dirty-indicator'}>
                      ${p.adminDirty
                        ? '未保存の変更があります。'
                        : '保存済みの内容です。'}
                    </p>
                  </div>
                  <button
                    type="button"
                    class="subtle"
                    ?disabled=${!p.adminDirty || p.adminSaving}
                    @click=${this.handleReset}
                  >
                    変更を元に戻す
                  </button>
                  <button type="submit" ?disabled=${p.adminSaving || !p.adminDirty}>
                    ${p.adminSaving ? '保存中...' : '保存する'}
                  </button>
                </div>
              </form>
            `}
      </section>
    `
  }

  private updateField<K extends keyof MeProfile>(key: K, value: MeProfile[K]) {
    this.setForm({ ...this.form, [key]: value })
  }

  private splitLines(value: string) {
    return value
      .split('\n')
      .map((item) => item.trim())
      .filter(Boolean)
  }

  private handleSubmit(event: Event) {
    event.preventDefault()
    void this.profileRepo.saveAdminProfile(cloneMeProfile(this.form))
  }

  private handleReset = () => {
    if (
      this.profileRepo.adminDirty &&
      !window.confirm('未保存の変更を破棄して元に戻しますか？')
    ) {
      return
    }

    this.setForm(cloneMeProfile(this.profileRepo.adminProfile))
  }

  private setForm(nextForm: MeProfile) {
    this.form = nextForm
    this.updateDirtyState(nextForm)
  }

  private updateDirtyState(nextForm: MeProfile) {
    const nextDirty = !this.profilesEqual(nextForm, this.profileRepo.adminProfile)
    this.profileRepo.setAdminDirty(nextDirty)
  }

  private profilesEqual(a: MeProfile, b: MeProfile) {
    return JSON.stringify(a) === JSON.stringify(b)
  }

  static styles = [
    adminFormStyles,
    css`
      :host {
        display: block;
      }

      .container {
        display: grid;
        gap: 24px;
      }

      .page-header {
        display: flex;
        gap: 24px;
        justify-content: space-between;
        align-items: end;
        flex-wrap: wrap;
      }

      .updated-at,
      .loading {
        color: var(--color-text-secondary);
        font-size: 14px;
        line-height: 1.8;
      }

      .meta {
        display: flex;
        gap: 8px;
        flex-wrap: wrap;
        margin-top: 16px;
      }

      .meta span {
        display: inline-flex;
        align-items: center;
        height: 28px;
        padding: 0 10px;
        border: 1px solid var(--color-border);
        background: var(--color-bg-surface);
        color: var(--color-text-secondary);
        font-size: 12px;
      }

      form {
        display: grid;
        gap: 24px;
      }

      .grid {
        display: grid;
        gap: 16px;
        grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
      }

      .actions {
        display: flex;
        align-items: center;
        gap: 12px;
        justify-content: end;
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
    'page-admin-profile': PageAdminProfile
  }
}
