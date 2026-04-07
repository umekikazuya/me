import { css, html, LitElement, nothing } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import type { MeProfile } from '../admin/types.js'
import { setupAmbientLines } from '../utils/ambient.js'
import { setupFade, setupReveal } from '../utils/scroll.js'

interface Article {
  title: string
  date: string
  slug: string
}

const mockArticles: Article[] = [
  {
    title: 'TypeScriptの型システムと向き合う',
    date: '2025-03',
    slug: 'ts-type-system',
  },
  {
    title: 'Litで作るWeb Components入門',
    date: '2025-01',
    slug: 'lit-web-components',
  },
  { title: '余白という設計思想', date: '2024-11', slug: 'design-whitespace' },
]

@customElement('page-top')
export class PageTop extends LitElement {
  @property({ attribute: false }) profile: MeProfile | null = null
  @property({ type: Boolean }) loading = false

  private cleanups: Array<() => void> = []

  firstUpdated() {
    const root = this.shadowRoot
    if (!root) return
    const revealEls = Array.from(
      root.querySelectorAll('.who > *, .articles-preview > *, .contact > *'),
    )
    const fadeEls = Array.from(root.querySelectorAll('.layer-1, .layer-2'))
    const fvSection = root.querySelector('.layer-0') as HTMLElement | null

    this.cleanups.push(setupReveal(revealEls, true))
    this.cleanups.push(setupFade(fadeEls))
    if (fvSection) this.cleanups.push(setupAmbientLines(fvSection))
  }

  disconnectedCallback() {
    super.disconnectedCallback()
    for (const cleanup of this.cleanups) cleanup()
    this.cleanups = []
  }

  render() {
    return html`
      <!-- Layer 0: First View -->
      <section class="layer layer-0 js-layer-0">
        <h1 class="name ${this.loading ? 'is-loading' : ''}">
          ${this.profile?.displayName ?? ''}
        </h1>
      </section>

      <!-- Layer 1: Who I am -->
      <section class="layer layer-1">
        <div class="who">
          <p class="role ${this.loading ? 'is-loading' : ''}">${this.profile?.role ?? ''}</p>
          <p class="location ${this.loading ? 'is-loading' : ''}">${this.profile?.location ?? ''}</p>
        </div>
      </section>

      <!-- Layer 2: Articles -->
      <section class="layer layer-2">
        <div class="articles-preview">
          <h2 class="section-label">Articles</h2>
          <ul class="article-list">
            ${mockArticles.map(
              (a) => html`
              <li class="article-item">
                <span class="article-date">${a.date}</span>
                <a href="/articles" class="article-title">${a.title}</a>
              </li>
            `,
            )}
          </ul>
          <a href="/articles" class="view-all">View all</a>
        </div>
      </section>

      <!-- Layer 3: Contact -->
      <section class="layer layer-3">
        <div class="contact">
          <p class="contact-label">Say Hello</p>
          <ul class="contact-links">
            ${
              this.profile
                ? this.profile.links.map(
                    (link) => html`
                    <li>
                      <a href=${link.url} target="_blank" rel="noopener">
                        ${link.platform}
                      </a>
                    </li>
                  `,
                  )
                : nothing
            }
          </ul>
        </div>
      </section>
    `
  }

  static styles = css`
    :host {
      display: block;
    }

    .layer {
      padding: var(--space-lg, 64px) var(--space-md, 32px);
    }

    /* Layer 0 */
    .layer-0 {
      height: 100dvh;
      display: flex;
      align-items: center;
      justify-content: center;
      padding: 0;
    }

    .name {
      font-family: var(--font-en);
      font-weight: 300;
      font-size: clamp(1.5rem, 2.5vw, 2.5rem);
      letter-spacing: var(--tracking-wider);
      color: var(--color-text-primary);
      margin: 0;
      animation: breathing 7s ease-in-out infinite;
    }

    @keyframes breathing {
      0%, 100% { opacity: 0.85; }
      50% { opacity: 1; }
    }

    /* Layer 1 */
    .layer-1 {
      display: block;
    }

    .who {
      max-width: 640px;
      margin: 0 auto;
    }

    .role {
      font-family: var(--font-en);
      font-weight: 300;
      font-size: 24px;
      letter-spacing: var(--tracking-wide);
      color: var(--color-text-primary);
      margin-bottom: 8px;
    }

    .location {
      font-family: var(--font-jp);
      font-weight: 200;
      font-size: 14px;
      letter-spacing: var(--tracking-wide);
      color: var(--color-text-secondary);
      margin-bottom: 40px;
    }

    .philosophy {
      font-family: var(--font-jp);
      font-weight: 200;
      font-size: 18px;
      letter-spacing: var(--tracking-wide);
      color: var(--color-text-primary);
      line-height: 2;
    }

    /* Layer 2 */
    .layer-2 {
      display: block;
    }

    .articles-preview {
      max-width: 640px;
      margin: 0 auto;
      width: 100%;
    }

    .section-label {
      font-family: var(--font-en);
      font-weight: 300;
      font-size: 13px;
      letter-spacing: var(--tracking-wider);
      text-transform: uppercase;
      color: var(--color-text-secondary);
      margin-bottom: 32px;
    }

    .article-list {
      list-style: none;
      padding: 0;
      margin: 0 0 40px;
    }

    .article-item {
      display: flex;
      align-items: baseline;
      gap: 24px;
      padding: 16px 8px;
      border-bottom: 1px solid var(--color-border);
      transition: background 0.2s ease, transform 0.2s ease;
    }

    .article-item:hover {
      background: var(--color-surface);
      transform: translateX(4px);
    }

    .article-date {
      font-family: var(--font-en);
      font-weight: 300;
      font-size: 13px;
      letter-spacing: var(--tracking-wide);
      color: var(--color-text-secondary);
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

    .view-all {
      font-family: var(--font-en);
      font-weight: 300;
      font-size: 14px;
      letter-spacing: var(--tracking-wide);
      color: var(--color-text-secondary);
      text-decoration: none;
      transition: opacity 0.2s ease;
    }

    .view-all:hover {
      opacity: 0.6;
    }

    /* Layer 3 */
    .layer-3 {
      display: block;
    }

    .contact {
      max-width: 640px;
      margin: 0 auto;
    }

    .contact-label {
      font-family: var(--font-en);
      font-weight: 300;
      font-size: 13px;
      letter-spacing: var(--tracking-wider);
      text-transform: uppercase;
      color: var(--color-text-secondary);
      margin: 0 0 20px;
    }

    .contact-links {
      list-style: none;
      padding: 0;
      margin: 0;
    }

    .contact-links li {
      border-bottom: 1px solid var(--color-border-light);
    }

    .contact-links li:first-child {
      border-top: 1px solid var(--color-border-light);
    }

    .contact-links a {
      display: block;
      font-family: var(--font-en);
      font-weight: 300;
      font-size: 15px;
      letter-spacing: var(--tracking-wide);
      color: var(--color-text-primary);
      text-decoration: none;
      padding: 12px 0;
      transition: opacity 0.2s ease;
    }

    .contact-links a:hover {
      opacity: 0.5;
    }

    .is-loading {
      opacity: 0.3;
    }

    @media (prefers-reduced-motion: reduce) {
      .name {
        animation: none;
        opacity: 1;
      }

      .article-item {
        transition: none;
        transform: none;
      }

      .article-title,
      .view-all,
      .contact-links a {
        transition: none;
      }
    }

    @media (max-width: 640px) {
      .layer {
        padding: var(--space-md, 32px) var(--space-sm, 16px);
      }
    }
  `
}

declare global {
  interface HTMLElementTagNameMap {
    'page-top': PageTop
  }
}
