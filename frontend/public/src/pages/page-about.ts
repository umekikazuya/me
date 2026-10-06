import { consume } from '@lit/context'
import { css, html, LitElement, nothing } from 'lit'
import { customElement } from 'lit/decorators.js'
import { profileContext } from '../contexts/profile-context.js'
import type { IProfileRepository } from '../domain/ProfileRepository.js'
import { pageStyles } from '../styles/page-styles.js'
import { sanitizeUrl } from '../utils/format.js'
import '../components/load-spinner.js'

@customElement('page-about')
export class PageAbout extends LitElement {
  @consume({ context: profileContext, subscribe: true })
  set profileRepo(repo: IProfileRepository) {
    if (this._profileRepo) {
      this._profileRepo.removeEventListener('change', this._onRepoChange)
    }
    this._profileRepo = repo
    repo.addEventListener('change', this._onRepoChange)
    this.requestUpdate()
  }
  get profileRepo() {
    return this._profileRepo
  }
  private _profileRepo!: IProfileRepository
  private _onRepoChange = () => this.requestUpdate()

  disconnectedCallback() {
    super.disconnectedCallback()
    if (this._profileRepo) {
      this._profileRepo.removeEventListener('change', this._onRepoChange)
    }
  }

  private get sortedSkills() {
    return [...(this.profileRepo.profile?.skills ?? [])].sort(
      (a, b) => a.sortOrder - b.sortOrder,
    )
  }

  render() {
    const p = this.profileRepo.profile

    return html`
      <header class="head">
        <h1>about</h1>
      </header>
      ${
        p
          ? this.renderProfile(p)
          : this.profileRepo.error
            ? html`<p class="note">プロフィールを読み込めませんでした。時間をおいて再読み込みしてください。</p>`
            : html`<load-spinner></load-spinner>`
      }
    `
  }

  private renderProfile(p: NonNullable<IProfileRepository['profile']>) {
    return html`
      ${
        this.sortedSkills.length > 0
          ? html`
        <section>
          <h2>skills</h2>
          <ul class="rows">
            ${this.sortedSkills.map(
              (group) => html`
                <li class="row">
                  <span class="main">${group.items.join(' / ')}</span>
                  <span class="meta">${group.category}</span>
                </li>
              `,
            )}
          </ul>
        </section>`
          : nothing
      }

      ${
        (p.experiences ?? []).length > 0
          ? html`
        <section>
          <h2>experience</h2>
          <ul class="rows">
            ${(p.experiences ?? []).map((exp) => {
              const years = html`<span class="meta">${exp.startYear} — ${exp.endYear ?? 'now'}</span>`
              return html`
                <li>
                  ${
                    exp.url
                      ? html`<a class="row" href=${sanitizeUrl(exp.url)} target="_blank" rel="noopener noreferrer"><span class="main">${exp.company}</span>${years}</a>`
                      : html`<div class="row"><span class="main">${exp.company}</span>${years}</div>`
                  }
                </li>
              `
            })}
          </ul>
        </section>`
          : nothing
      }

      ${
        (p.certifications ?? []).length > 0
          ? html`
        <section>
          <h2>certifications</h2>
          <ul class="rows">
            ${(p.certifications ?? []).map(
              (cert) => html`
                <li class="row">
                  <span class="main">${cert.name}</span>
                  <span class="meta">${cert.year}</span>
                </li>
              `,
            )}
          </ul>
        </section>`
          : nothing
      }

      ${
        (p.likes ?? []).length > 0
          ? html`
        <section>
          <h2>likes</h2>
          <p class="likes">${(p.likes ?? []).join(', ')}</p>
        </section>`
          : nothing
      }
    `
  }

  static styles = [
    pageStyles,
    css`
      .head {
        padding-block: clamp(96px, 18vh, 160px) 56px;
      }

      h1 {
        margin: 0;
        font-family: var(--font-en);
        font-size: 22px;
        font-weight: 500;
        letter-spacing: var(--tracking-tight);
      }

      .likes {
        margin: 0;
      }
    `,
  ]
}

declare global {
  interface HTMLElementTagNameMap {
    'page-about': PageAbout
  }
}
