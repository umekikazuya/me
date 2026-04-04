import { css, html, LitElement } from 'lit'
import { customElement } from 'lit/decorators.js'
import { setupReveal } from '../utils/scroll.js'

@customElement('page-about')
export class PageAbout extends LitElement {
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

  render() {
    return html`
      <div class="container">
        <header class="page-header">
          <h1 class="page-title">About</h1>
        </header>

        <section class="section">
          <h2 class="section-title">Philosophy</h2>
          <p class="section-text">
            サンプルテキスト。<br />
            サンプルテキスト。<br />
            サンプルテキスト。
          </p>
        </section>

        <section class="section">
          <h2 class="section-title">Skills</h2>
          <ul class="list">
            <li>TypeScript / JavaScript</li>
            <li>Go</li>
            <li>Web Components / Lit</li>
            <li>React / Vue</li>
            <li>Node.js</li>
            <li>PostgreSQL / MySQL</li>
            <li>Docker / Kubernetes</li>
          </ul>
        </section>

        <section class="section">
          <h2 class="section-title">Certifications</h2>
          <ul class="list">
            <li>サンプルテキスト。</li>
            <li>サンプルテキスト。</li>
          </ul>
        </section>

        <section class="section">
          <h2 class="section-title">Experience</h2>
          <ul class="list">
            <li>2023 — 現在　某スタートアップ　Software Engineer</li>
            <li>2021 — 2023　某SIer　Backend Engineer</li>
          </ul>
        </section>

        <section class="section">
          <h2 class="section-title">Likes</h2>
          <ul class="list">
            <li>サンプルテキスト。</li>
            <li>サンプルテキスト。</li>
            <li>サンプルテキスト。</li>
            <li>サンプルテキスト。</li>
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
      border-bottom: 1px solid var(--color-border-light);
      line-height: 1.6;
    }

    .list li:first-child {
      border-top: 1px solid var(--color-border-light);
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
