import { css, html, LitElement } from 'lit'
import { customElement, property } from 'lit/decorators.js'

@customElement('me-admin-section')
export class MeAdminSection extends LitElement {
  @property() title = ''
  @property() description = ''

  render() {
    return html`
      <section class="section">
        <div class="header">
          <div class="copy">
            <h2 class="title">${this.title}</h2>
            ${
              this.description
                ? html`<p class="description">${this.description}</p>`
                : null
            }
          </div>
          <div class="actions">
            <slot name="header-actions"></slot>
          </div>
        </div>
        <div class="content">
          <slot></slot>
        </div>
      </section>
    `
  }

  static styles = css`
    :host {
      display: block;
    }

    *, *::before, *::after {
      box-sizing: border-box;
    }

    .section {
      display: grid;
      gap: 16px;
      padding: 24px;
      border: 1px solid var(--color-border);
      background: var(--color-bg-surface);
    }

    .header {
      display: flex;
      justify-content: space-between;
      gap: 16px;
      align-items: center;
      flex-wrap: wrap;
    }

    .copy {
      display: grid;
      gap: 6px;
    }

    .title {
      font-size: 16px;
      font-weight: 500;
      color: var(--color-text-primary);
      margin: 0;
    }

    .description {
      color: var(--color-text-tertiary);
      font-size: 13px;
      line-height: 1.8;
      margin: 0;
    }

    .content {
      display: grid;
      gap: 20px;
    }
  `
}

declare global {
  interface HTMLElementTagNameMap {
    'me-admin-section': MeAdminSection
  }
}
