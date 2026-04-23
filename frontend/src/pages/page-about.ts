import { css, html, LitElement } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import type { MeProfile } from '../admin/types.js'
import { setupReveal } from '../utils/scroll.js'

@customElement('page-about')
export class PageAbout extends LitElement {
  @property({ attribute: false }) profile: MeProfile | null = null
  @property({ type: Boolean }) loading = false

  private cleanups: Array<() => void> = []

  firstUpdated() {
    const root = this.shadowRoot
    if (!root) return
    const revealEls = Array.from(
      root.querySelectorAll('.page-header, .section'),
    )
    this.cleanups.push(setupReveal(revealEls, true))
  }

  disconnectedCallback() {
    super.disconnectedCallback()
    for (const cleanup of this.cleanups) cleanup()
    this.cleanups = []
  }

  private get sortedSkills() {
    return [...(this.profile?.skills ?? [])].sort(
      (a, b) => a.sortOrder - b.sortOrder,
    )
  }

  render() {
    const p = this.profile
    const cls = this.loading ? 'is-loading' : ''

    return html`
      <div class="container ${cls}">
        <header class="page-header">
          <h1 class="page-title">About</h1>
        </header>

        <section class="section">
          <h2 class="section-title">Skills</h2>
          <ul class="list">
            ${this.sortedSkills.map(
              (group) => html`
                <li>
                  <span class="skill-category">${group.category}</span>
                  <span class="skill-items">${group.items.join(' / ')}</span>
                </li>
              `,
            )}
          </ul>
        </section>

        <section class="section">
          <h2 class="section-title">Certifications</h2>
          <ul class="list">
            ${(p?.certifications ?? []).map(
              (cert) => html`
                <li>
                  <span class="cert-name">${cert.name}</span>
                  <span class="cert-meta">
                    ${cert.issuer ? html`${cert.issuer} &middot; ` : ''}${cert.year}
                  </span>
                </li>
              `,
            )}
          </ul>
        </section>

        <section class="section">
          <h2 class="section-title">Experience</h2>
          <ul class="list">
            ${(p?.experiences ?? []).map(
              (exp) => html`
                <li>
                  <span class="exp-years">
                    ${exp.startYear} — ${exp.endYear ?? '現在'}
                  </span>
                  <span class="exp-company">${exp.company}</span>
                </li>
              `,
            )}
          </ul>
        </section>

        <section class="section">
          <h2 class="section-title">Likes</h2>
          <ul class="list">
            ${(p?.likes ?? []).map((like) => html`<li>${like}</li>`)}
          </ul>
        </section>
      </div>
    `
  }

  static styles = css`
    :host {
      display: block;
      padding-top: 80px;
    }

    .container {
      max-width: 640px;
      margin: 0 auto;
      padding: var(--space-lg) var(--space-md);
    }

    .page-header {
      margin-bottom: 64px;
    }

    .page-title {
      font-family: var(--font-en);
      font-weight: 300;
      font-size: 36px;
      letter-spacing: var(--tracking-wide);
      color: var(--color-text-primary);
      margin: 0;
    }

    .section {
      margin-bottom: 56px;
    }

    .section-title {
      font-family: var(--font-en);
      font-weight: 300;
      font-size: 13px;
      letter-spacing: var(--tracking-wider);
      text-transform: uppercase;
      color: var(--color-text-secondary);
      margin: 0 0 20px;
    }

    .section-text {
      font-family: var(--font-jp);
      font-weight: 200;
      font-size: 15px;
      letter-spacing: 0.04em;
      line-height: 2.2;
      color: var(--color-text-primary);
    }

    .list {
      list-style: none;
      padding: 0;
      margin: 0;
    }

    .list li {
      font-family: var(--font-jp);
      font-weight: 200;
      font-size: 15px;
      letter-spacing: 0.04em;
      color: var(--color-text-primary);
      padding: 12px 0;
      border-bottom: 1px solid var(--color-border-subtle);
      line-height: 1.6;
      display: flex;
      flex-direction: column;
      gap: 2px;
    }

    .list li:first-child {
      border-top: 1px solid var(--color-border-subtle);
    }

    .skill-category,
    .exp-years,
    .cert-meta {
      font-family: var(--font-en);
      font-size: 12px;
      letter-spacing: var(--tracking-wide);
      color: var(--color-text-secondary);
    }

    .is-loading {
      opacity: 0.3;
    }

    @media (max-width: 640px) {
      .container {
        padding: 48px 24px;
      }
    }
  `
}

declare global {
  interface HTMLElementTagNameMap {
    'page-about': PageAbout
  }
}
