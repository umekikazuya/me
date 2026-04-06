import { css, html, LitElement } from 'lit'
import { customElement, property, state } from 'lit/decorators.js'
import {
  cloneMeProfile,
  createEmptyMeProfile,
  type MeCertification,
  type MeExperience,
  type MeLink,
  type MeProfile,
  type MeSkillGroup,
} from '../admin/types.js'

@customElement('page-admin-profile')
export class PageAdminProfile extends LitElement {
  @property({ attribute: false })
  profile: MeProfile = createEmptyMeProfile()

  @property({ type: Boolean })
  loading = false

  @property({ type: Boolean })
  saving = false

  @property()
  errorMessage = ''

  @property()
  successMessage = ''

  @state()
  private form: MeProfile = createEmptyMeProfile()

  protected willUpdate(changedProperties: Map<PropertyKey, unknown>) {
    if (changedProperties.has('profile')) {
      this.form = cloneMeProfile(this.profile)
    }
  }

  render() {
    return html`
      <section class="container">
        <header class="page-header">
          <div>
            <p class="eyebrow" lang="en">Profile</p>
            <h1 class="title">プロフィール編集</h1>
            <p class="description">
              公開プロフィールの表示内容を更新します。
            </p>
          </div>
          ${
            this.form.updatedAt
              ? html`<p class="updated-at">
                最終更新: ${new Date(this.form.updatedAt).toLocaleString('ja-JP')}
              </p>`
              : null
          }
        </header>

        ${this.errorMessage ? html`<p class="message error">${this.errorMessage}</p>` : null}
        ${
          this.successMessage
            ? html`<p class="message success">${this.successMessage}</p>`
            : null
        }

        ${
          this.loading
            ? html`<p class="loading">プロフィールを読み込み中...</p>`
            : html`
              <form @submit=${this.handleSubmit}>
                <section class="section">
                  <h2>基本情報</h2>
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
                    <h2>Skills</h2>
                    <button type="button" class="subtle" @click=${this.addSkill}>
                      カテゴリを追加
                    </button>
                  </div>
                  <div class="stack">
                    ${this.form.skills.map(
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
                                    category: (event.target as HTMLInputElement)
                                      .value,
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
                    )}
                  </div>
                </section>

                <section class="section">
                  <div class="section-header">
                    <h2>Certifications</h2>
                    <button
                      type="button"
                      class="subtle"
                      @click=${this.addCertification}
                    >
                      資格を追加
                    </button>
                  </div>
                  <div class="stack">
                    ${this.form.certifications.map(
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
                    )}
                  </div>
                </section>

                <section class="section">
                  <div class="section-header">
                    <h2>Experiences</h2>
                    <button type="button" class="subtle" @click=${this.addExperience}>
                      経歴を追加
                    </button>
                  </div>
                  <div class="stack">
                    ${this.form.experiences.map(
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
                    )}
                  </div>
                </section>

                <section class="section">
                  <div class="section-header">
                    <h2>Links</h2>
                    <button type="button" class="subtle" @click=${this.addLink}>
                      リンクを追加
                    </button>
                  </div>
                  <div class="stack">
                    ${this.form.links.map(
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
                    )}
                  </div>
                </section>

                <section class="section">
                  <h2>Likes</h2>
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
                  <button type="submit" ?disabled=${this.saving}>
                    ${this.saving ? '保存中...' : '保存する'}
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
    this.form = {
      ...this.form,
      [key]: value,
    }
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
        label: '',
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
    this.dispatchEvent(
      new CustomEvent<MeProfile>('admin-save-profile', {
        detail: cloneMeProfile(this.form),
        bubbles: true,
        composed: true,
      }),
    )
  }

  static styles = css`
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

    .eyebrow {
      font-family: var(--font-en);
      letter-spacing: var(--tracking-wider);
      color: var(--color-text-tertiary);
      margin-bottom: 12px;
    }

    .title {
      font-size: 30px;
      font-weight: 300;
      margin-bottom: 12px;
    }

    .description,
    .updated-at,
    .loading {
      color: var(--color-text-secondary);
      line-height: 1.8;
    }

    form,
    .stack {
      display: grid;
      gap: 20px;
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
      border: 1px solid var(--color-border-light);
      background: var(--color-surface);
      padding: 20px;
    }

    h2,
    h3 {
      font-weight: 300;
    }

    .grid {
      display: grid;
      gap: 16px;
      grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
    }

    .field {
      display: grid;
      gap: 8px;
      color: var(--color-text-secondary);
      font-size: 14px;
    }

    .field-wide {
      grid-column: 1 / -1;
    }

    input,
    textarea {
      width: 100%;
      border: 1px solid var(--color-border);
      background: #fff;
      padding: 12px 14px;
      color: var(--color-text-primary);
      font: inherit;
    }

    input:focus,
    textarea:focus {
      outline: none;
      border-color: var(--color-text-primary);
    }

    textarea {
      resize: vertical;
    }

    button {
      border: 0;
      background: var(--color-text-primary);
      color: #fff;
      padding: 12px 18px;
      font: inherit;
      cursor: pointer;
    }

    button:disabled {
      opacity: 0.5;
      cursor: wait;
    }

    .subtle {
      background: transparent;
      border: 1px solid var(--color-border);
      color: var(--color-text-secondary);
      padding: 10px 14px;
    }

    .danger {
      color: #9a3f3f;
    }

    .actions {
      display: flex;
      justify-content: end;
    }

    .message {
      font-size: 14px;
      line-height: 1.7;
    }

    .error {
      color: #9a3f3f;
    }

    .success {
      color: #3d7a56;
    }
  `
}

declare global {
  interface HTMLElementTagNameMap {
    'page-admin-profile': PageAdminProfile
  }
}
