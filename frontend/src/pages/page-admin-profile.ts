import { consume } from '@lit/context'
import { css, html, LitElement } from 'lit'
import { customElement, state } from 'lit/decorators.js'
import { adminFormStyles } from '../admin/admin-form-styles.js'
import {
  cloneMeProfile,
  createEmptyMeProfile,
  type MeCertification,
  type MeExperience,
  type MeLink,
  type MeProfile,
  type MeSkillGroup,
} from '../admin/types.js'
import { profileContext } from '../contexts/profile-context.js'
import { RepositoryObserver } from '../controllers/RepositoryObserver.js'
import type { IProfileRepository } from '../domain/ProfileRepository.js'

@customElement('page-admin-profile')
export class PageAdminProfile extends LitElement {
  @consume({ context: profileContext, subscribe: true })
  profileRepo!: IProfileRepository

  @state()
  private form: MeProfile = createEmptyMeProfile()

  private _lastSyncedData = ''

  constructor() {
    super()
    new RepositoryObserver(this, this.profileRepo)
  }

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
            <p class="description">
              公開プロフィールの表示内容を更新します。
            </p>
            <div class="meta">
              <span>Skill: ${this.form.skills.length}</span>
              <span>Certification: ${this.form.certifications.length}</span>
              <span>Experience: ${this.form.experiences.length}</span>
              <span>Link: ${this.form.links.length}</span>
            </div>
          </div>
          ${
            this.form.updatedAt
              ? html`<p class="updated-at">
                最終更新: ${new Date(this.form.updatedAt).toLocaleString('ja-JP')}
              </p>`
              : null
          }
        </header>

        ${p.adminError ? html`<p class="message error">${p.adminError}</p>` : null}
        ${
          p.adminSuccess
            ? html`<p class="message success">${p.adminSuccess}</p>`
            : null
        }

        ${
          p.adminLoading
            ? html`<p class="loading">プロフィールを読み込み中...</p>`
            : html`
              <form @submit=${this.handleSubmit}>
                <section class="section">
                  <div class="section-copy">
                    <h2>基本情報</h2>
                    <p class="section-help">
                      最低限、表示名だけあれば更新できます。未入力項目は公開画面で省略されます。
                    </p>
                  </div>
                  <div class="grid">
                    <label class="field">
                      <span>表示名 *</span>
                      <input
                        .value=${this.form.displayName}
                        @input=${(event: Event) =>
                          this.updateField(
                            'displayName',
                            (event.target as HTMLInputElement).value,
                          )}
                        required
                      />
                    </label>
                    <label class="field">
                      <span>表示名（日本語）</span>
                      <input
                        .value=${this.form.displayJa}
                        @input=${(event: Event) =>
                          this.updateField(
                            'displayJa',
                            (event.target as HTMLInputElement).value,
                          )}
                      />
                    </label>
                    <label class="field">
                      <span>Role</span>
                      <input
                        .value=${this.form.role}
                        @input=${(event: Event) =>
                          this.updateField(
                            'role',
                            (event.target as HTMLInputElement).value,
                          )}
                      />
                    </label>
                    <label class="field">
                      <span>Location</span>
                      <input
                        .value=${this.form.location}
                        @input=${(event: Event) =>
                          this.updateField(
                            'location',
                            (event.target as HTMLInputElement).value,
                          )}
                      />
                    </label>
                  </div>
                </section>

                <section class="section">
                  <div class="section-header">
                    <div class="section-copy">
                      <h2>Skills</h2>
                      <p class="section-help">
                        カテゴリごとに整理し、Items は1行ずつ入力すると編集しやすいです。
                      </p>
                    </div>
                    <button type="button" class="subtle" @click=${this.addSkill}>
                      カテゴリを追加
                    </button>
                  </div>
                  <div class="stack">
                    ${
                      this.form.skills.length === 0
                        ? this.renderEmptyPanel(
                            'まだ skill カテゴリがありません。',
                            'カテゴリを追加',
                            this.addSkill,
                          )
                        : this.form.skills.map(
                            (skill, index) => html`
                            <article class="panel">
                              <div class="panel-header">
                                <h3>カテゴリ ${index + 1}</h3>
                                <button
                                  type="button"
                                  class="subtle danger"
                                  @click=${() => this.removeSkill(index)}
                                >
                                  削除
                                </button>
                              </div>
                              <div class="grid">
                                <label class="field">
                                  <span>Category</span>
                                  <input
                                    .value=${skill.category}
                                    @input=${(event: Event) =>
                                      this.updateSkill(index, {
                                        category: (
                                          event.target as HTMLInputElement
                                        ).value,
                                      })}
                                  />
                                </label>
                                <label class="field">
                                  <span>Sort order</span>
                                  <input
                                    type="number"
                                    .value=${String(skill.sortOrder)}
                                    @input=${(event: Event) =>
                                      this.updateSkill(index, {
                                        sortOrder: Number(
                                          (event.target as HTMLInputElement)
                                            .value || '0',
                                        ),
                                      })}
                                  />
                                </label>
                                <label class="field field-wide">
                                  <span>Items（改行区切り）</span>
                                  <textarea
                                    rows="4"
                                    .value=${skill.items.join('\n')}
                                    @input=${(event: Event) =>
                                      this.updateSkill(index, {
                                        items: this.splitLines(
                                          (event.target as HTMLTextAreaElement)
                                            .value,
                                        ),
                                      })}
                                  ></textarea>
                                </label>
                              </div>
                            </article>
                          `,
                          )
                    }
                  </div>
                </section>

                <section class="section">
                  <div class="section-header">
                    <div class="section-copy">
                      <h2>Certifications</h2>
                      <p class="section-help">
                        month は任意です。年だけでも掲載できます。
                      </p>
                    </div>
                    <button
                      type="button"
                      class="subtle"
                      @click=${this.addCertification}
                    >
                      資格を追加
                    </button>
                  </div>
                  <div class="stack">
                    ${
                      this.form.certifications.length === 0
                        ? this.renderEmptyPanel(
                            '資格がまだありません。',
                            '資格を追加',
                            this.addCertification,
                          )
                        : this.form.certifications.map(
                            (certification, index) => html`
                        <article class="panel">
                          <div class="panel-header">
                            <h3>資格 ${index + 1}</h3>
                            <button
                              type="button"
                              class="subtle danger"
                              @click=${() => this.removeCertification(index)}
                            >
                              削除
                            </button>
                          </div>
                          <div class="grid">
                            <label class="field">
                              <span>資格名</span>
                              <input
                                .value=${certification.name}
                                @input=${(event: Event) =>
                                  this.updateCertification(index, {
                                    name: (event.target as HTMLInputElement)
                                      .value,
                                  })}
                              />
                            </label>
                            <label class="field">
                              <span>Issuer</span>
                              <input
                                .value=${certification.issuer}
                                @input=${(event: Event) =>
                                  this.updateCertification(index, {
                                    issuer: (event.target as HTMLInputElement)
                                      .value,
                                  })}
                              />
                            </label>
                            <label class="field">
                              <span>Year</span>
                              <input
                                type="number"
                                .value=${String(certification.year)}
                                @input=${(event: Event) =>
                                  this.updateCertification(index, {
                                    year: Number(
                                      (event.target as HTMLInputElement)
                                        .value || '0',
                                    ),
                                  })}
                              />
                            </label>
                            <label class="field">
                              <span>Month</span>
                              <input
                                type="number"
                                min="1"
                                max="12"
                                .value=${
                                  certification.month
                                    ? String(certification.month)
                                    : ''
                                }
                                @input=${(event: Event) =>
                                  this.updateCertification(index, {
                                    month: this.toOptionalNumber(
                                      (event.target as HTMLInputElement).value,
                                    ),
                                  })}
                              />
                            </label>
                          </div>
                        </article>
                          `,
                          )
                    }
                  </div>
                </section>

                <section class="section">
                  <div class="section-header">
                    <div class="section-copy">
                      <h2>Experiences</h2>
                      <p class="section-help">
                        endYear を空にすると、継続中の経歴として扱えます。
                      </p>
                    </div>
                    <button type="button" class="subtle" @click=${this.addExperience}>
                      経歴を追加
                    </button>
                  </div>
                  <div class="stack">
                    ${
                      this.form.experiences.length === 0
                        ? this.renderEmptyPanel(
                            '経歴がまだありません。',
                            '経歴を追加',
                            this.addExperience,
                          )
                        : this.form.experiences.map(
                            (experience, index) => html`
                        <article class="panel">
                          <div class="panel-header">
                            <h3>経歴 ${index + 1}</h3>
                            <button
                              type="button"
                              class="subtle danger"
                              @click=${() => this.removeExperience(index)}
                            >
                              削除
                            </button>
                          </div>
                          <div class="grid">
                            <label class="field">
                              <span>Company</span>
                              <input
                                .value=${experience.company}
                                @input=${(event: Event) =>
                                  this.updateExperience(index, {
                                    company: (event.target as HTMLInputElement)
                                      .value,
                                  })}
                              />
                            </label>
                            <label class="field">
                              <span>URL</span>
                              <input
                                .value=${experience.url}
                                @input=${(event: Event) =>
                                  this.updateExperience(index, {
                                    url: (event.target as HTMLInputElement)
                                      .value,
                                  })}
                              />
                            </label>
                            <label class="field">
                              <span>Start year</span>
                              <input
                                type="number"
                                .value=${String(experience.startYear)}
                                @input=${(event: Event) =>
                                  this.updateExperience(index, {
                                    startYear: Number(
                                      (event.target as HTMLInputElement)
                                        .value || '0',
                                    ),
                                  })}
                              />
                            </label>
                            <label class="field">
                              <span>End year</span>
                              <input
                                type="number"
                                .value=${
                                  experience.endYear
                                    ? String(experience.endYear)
                                    : ''
                                }
                                @input=${(event: Event) =>
                                  this.updateExperience(index, {
                                    endYear: this.toOptionalNumber(
                                      (event.target as HTMLInputElement).value,
                                    ),
                                  })}
                              />
                            </label>
                          </div>
                        </article>
                          `,
                          )
                    }
                  </div>
                </section>

                <section class="section">
                  <div class="section-header">
                    <div class="section-copy">
                      <h2>Links</h2>
                      <p class="section-help">
                        platform と URL は必須です。label は公開側で見せたい名前を指定します。
                      </p>
                    </div>
                    <button type="button" class="subtle" @click=${this.addLink}>
                      リンクを追加
                    </button>
                  </div>
                  <div class="stack">
                    ${
                      this.form.links.length === 0
                        ? this.renderEmptyPanel(
                            'リンクがまだありません。',
                            'リンクを追加',
                            this.addLink,
                          )
                        : this.form.links.map(
                            (link, index) => html`
                        <article class="panel">
                          <div class="panel-header">
                            <h3>リンク ${index + 1}</h3>
                            <button
                              type="button"
                              class="subtle danger"
                              @click=${() => this.removeLink(index)}
                            >
                              削除
                            </button>
                          </div>
                          <div class="grid">
                            <label class="field">
                              <span>Platform</span>
                              <input
                                .value=${link.platform}
                                @input=${(event: Event) =>
                                  this.updateLink(index, {
                                    platform: (event.target as HTMLInputElement)
                                      .value,
                                  })}
                              />
                            </label>
                            <label class="field">
                              <span>URL</span>
                              <input
                                .value=${link.url}
                                @input=${(event: Event) =>
                                  this.updateLink(index, {
                                    url: (event.target as HTMLInputElement)
                                      .value,
                                  })}
                              />
                            </label>
                          </div>
                        </article>
                          `,
                          )
                    }
                  </div>
                </section>

                <section class="section">
                  <div class="section-copy">
                    <h2>Likes</h2>
                    <p class="section-help">
                      1行ごとに1件ずつ入力します。空行は保存時に除外されます。
                    </p>
                  </div>
                  <label class="field">
                    <span>1行につき1件</span>
                    <textarea
                      rows="6"
                      .value=${this.form.likes.join('\n')}
                      @input=${(event: Event) =>
                        this.updateField(
                          'likes',
                          this.splitLines(
                            (event.target as HTMLTextAreaElement).value,
                          ),
                        )}
                    ></textarea>
                  </label>
                </section>

                <div class="actions">
                  <div class="actions-copy">
                    <p class=${p.adminDirty ? 'dirty-indicator dirty' : 'dirty-indicator'}>
                      ${p.adminDirty ? '未保存の変更があります。' : '保存済みの内容です。'}
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
            `
        }
      </section>
    `
  }

  private updateField<Key extends keyof MeProfile>(
    key: Key,
    value: MeProfile[Key],
  ) {
    this.setForm({
      ...this.form,
      [key]: value,
    })
  }

  private updateSkill(index: number, patch: Partial<MeSkillGroup>) {
    const skills = [...this.form.skills]
    skills[index] = { ...skills[index], ...patch }
    this.updateField('skills', skills)
  }

  private updateCertification(index: number, patch: Partial<MeCertification>) {
    const certifications = [...this.form.certifications]
    certifications[index] = { ...certifications[index], ...patch }
    this.updateField('certifications', certifications)
  }

  private updateExperience(index: number, patch: Partial<MeExperience>) {
    const experiences = [...this.form.experiences]
    experiences[index] = { ...experiences[index], ...patch }
    this.updateField('experiences', experiences)
  }

  private updateLink(index: number, patch: Partial<MeLink>) {
    const links = [...this.form.links]
    links[index] = { ...links[index], ...patch }
    this.updateField('links', links)
  }

  private addSkill = () => {
    this.updateField('skills', [
      ...this.form.skills,
      { category: '', items: [], sortOrder: this.form.skills.length },
    ])
  }

  private removeSkill(index: number) {
    this.updateField(
      'skills',
      this.form.skills.filter((_, itemIndex) => itemIndex !== index),
    )
  }

  private addCertification = () => {
    this.updateField('certifications', [
      ...this.form.certifications,
      {
        name: '',
        issuer: '',
        year: new Date().getFullYear(),
      },
    ])
  }

  private removeCertification(index: number) {
    this.updateField(
      'certifications',
      this.form.certifications.filter((_, itemIndex) => itemIndex !== index),
    )
  }

  private addExperience = () => {
    this.updateField('experiences', [
      ...this.form.experiences,
      {
        company: '',
        url: '',
        startYear: new Date().getFullYear(),
      },
    ])
  }

  private removeExperience(index: number) {
    this.updateField(
      'experiences',
      this.form.experiences.filter((_, itemIndex) => itemIndex !== index),
    )
  }

  private addLink = () => {
    this.updateField('links', [
      ...this.form.links,
      {
        platform: '',
        url: '',
      },
    ])
  }

  private removeLink(index: number) {
    this.updateField(
      'links',
      this.form.links.filter((_, itemIndex) => itemIndex !== index),
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
    const nextDirty = !this.profilesEqual(
      nextForm,
      this.profileRepo.adminProfile,
    )
    this.profileRepo.setAdminDirty(nextDirty)
  }

  private profilesEqual(a: MeProfile, b: MeProfile) {
    return JSON.stringify(a) === JSON.stringify(b)
  }

  private renderEmptyPanel(
    message: string,
    actionLabel: string,
    onClick: () => void,
  ) {
    return html`<article class="empty-panel">
      <p>${message}</p>
      <button type="button" class="subtle" @click=${onClick}>${actionLabel}</button>
    </article>`
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
        min-height: 28px;
        padding: 0 10px;
        border: 1px solid var(--color-border);
        background: var(--color-bg-surface);
        color: var(--color-text-secondary);
        font-size: 12px;
      }

      form,
      .stack {
        display: grid;
        gap: 20px;
      }

      .section-copy {
        display: grid;
        gap: 6px;
      }

      .section-help {
        color: var(--color-text-tertiary);
        font-size: 13px;
        line-height: 1.8;
      }

      .section {
        display: grid;
        gap: 16px;
        padding: 24px;
        border: 1px solid var(--color-border);
        background: #fff;
      }

      .section-header,
      .panel-header {
        display: flex;
        justify-content: space-between;
        gap: 16px;
        align-items: center;
        flex-wrap: wrap;
      }

      .panel {
        display: grid;
        gap: 16px;
        border: 1px solid var(--color-border-subtle);
        background: var(--color-bg-surface);
        padding: 20px;
      }

      .empty-panel {
        display: grid;
        justify-items: start;
        gap: 12px;
        border: 1px dashed var(--color-border);
        padding: 18px;
        color: var(--color-text-secondary);
        font-size: 14px;
      }

      h2 {
        font-size: 16px;
        font-weight: 500;
        color: var(--color-text-primary);
      }

      h3 {
        font-size: 14px;
        font-weight: 500;
        color: var(--color-text-secondary);
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
