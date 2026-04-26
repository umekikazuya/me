import { consume } from '@lit/context'
import { css, html, LitElement } from 'lit'
import { customElement } from 'lit/decorators.js'
import { adminFormStyles } from '../admin/admin-form-styles.js'
import type { MeProfile } from '../admin/types.js'
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
    if (this._profileRepo === repo) return
    this._profileRepo = repo
    this._observer?.disconnect()
    if (repo) this._observer = new RepositoryObserver(this, repo)
  }
  get profileRepo() {
    return this._profileRepo
  }
  private _profileRepo!: IProfileRepository
  private _observer?: RepositoryObserver

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

  render() {
    const p = this.profileRepo
    const profile = p.adminProfile

    return html`
      <section class="container">
        <header class="page-header">
          <div>
            <p class="eyebrow" lang="en">Profile</p>
            <h1 class="title">プロフィール編集</h1>
            <p class="description">公開プロフィールの表示内容を更新します。</p>
          </div>
          ${
            profile.updatedAt
              ? html`<p class="updated-at">
                最終更新: ${new Date(profile.updatedAt).toLocaleString('ja-JP')}
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
              <form @submit=${this.handleSubmit} @input=${this.handleInput}>
                <me-admin-section
                  title="基本情報"
                  description="最低限、表示名だけあれば更新できます。未入力項目は公開画面で省略されます。"
                >
                  <div class="grid">
                    <me-text-input
                      label="表示名 *"
                      name="displayName"
                      .value=${profile.displayName}
                      required
                    ></me-text-input>

                    <me-text-input
                      label="表示名（日本語）"
                      name="displayJa"
                      .value=${profile.displayJa}
                    ></me-text-input>

                    <me-text-input
                      label="Role"
                      name="role"
                      .value=${profile.role}
                    ></me-text-input>

                    <me-text-input
                      label="Location"
                      name="location"
                      .value=${profile.location}
                    ></me-text-input>
                  </div>
                </me-admin-section>

                <me-profile-skills-editor
                  name="skills"
                  .skills=${profile.skills}
                ></me-profile-skills-editor>

                <me-profile-certifications-editor
                  name="certifications"
                  .certifications=${profile.certifications}
                ></me-profile-certifications-editor>

                <me-profile-experiences-editor
                  name="experiences"
                  .experiences=${profile.experiences}
                ></me-profile-experiences-editor>

                <me-profile-links-editor
                  name="links"
                  .links=${profile.links}
                ></me-profile-links-editor>

                <me-admin-section
                  title="Likes"
                  description="1行ごとに1件ずつ入力します。空行は保存時に除外されます。"
                >
                  <me-textarea
                    label="1行につき1件"
                    name="likes"
                    rows="6"
                    .value=${profile.likes.join('\n')}
                  ></me-textarea>
                </me-admin-section>

                <div class="actions">
                  <div class="actions-copy">
                    <p class=${p.adminDirty ? 'dirty-indicator dirty' : 'dirty-indicator'}>
                      ${
                        p.adminDirty
                          ? '未保存の変更があります。'
                          : '保存済みの内容です。'
                      }
                    </p>
                  </div>
                  <button
                    type="reset"
                    class="subtle"
                    ?disabled=${p.adminSaving}
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

  private handleInput() {
    this.profileRepo.setAdminDirty(true)
  }

  private async handleSubmit(event: Event) {
    event.preventDefault()
    const form = event.target as HTMLFormElement
    if (!form.checkValidity()) {
      form.reportValidity()
      return
    }

    const formData = new FormData(form)

    const profile: MeProfile = {
      displayName: formData.get('displayName') as string,
      displayJa: formData.get('displayJa') as string,
      role: formData.get('role') as string,
      location: formData.get('location') as string,
      skills: JSON.parse((formData.get('skills') as string) || '[]'),
      certifications: JSON.parse(
        (formData.get('certifications') as string) || '[]',
      ),
      experiences: JSON.parse((formData.get('experiences') as string) || '[]'),
      links: JSON.parse((formData.get('links') as string) || '[]'),
      likes: ((formData.get('likes') as string) || '')
        .split('\n')
        .map((s) => s.trim())
        .filter(Boolean),
      updatedAt: this.profileRepo.adminProfile.updatedAt,
    }

    await this.profileRepo.saveAdminProfile(profile)
  }

  private handleReset = (e: Event) => {
    e.preventDefault()
    if (
      this.profileRepo.adminDirty &&
      !window.confirm('未保存の変更を破棄して元に戻しますか？')
    ) {
      return
    }
    this.profileRepo.setAdminDirty(false)
    this.requestUpdate()
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
