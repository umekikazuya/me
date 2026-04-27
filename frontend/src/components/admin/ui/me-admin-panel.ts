import { css, html, LitElement } from 'lit'
import { customElement, property } from 'lit/decorators.js'

@customElement('me-admin-panel')
export class MeAdminPanel extends LitElement {
  @property() title = ''

  render() {
    return html`
      <article class="panel">
        <div class="header">
          <h3 class="title">${this.title}</h3>
          <div class="actions">
            <slot name="header-actions"></slot>
          </div>
        </div>
        <div class="content">
          <slot></slot>
        </div>
      </article>
    `
  }

  static styles = css`
    :host {
      display: block;
    }

    .panel {
      display: grid;
      gap: 16px;
      border: 1px solid var(--color-border-subtle);
      background: var(--color-bg-surface);
      padding: 20px;
    }

    .header {
      display: flex;
      justify-content: space-between;
      gap: 16px;
      align-items: center;
      flex-wrap: wrap;
    }

    .title {
      font-size: 14px;
      font-weight: 500;
      color: var(--color-text-secondary);
      margin: 0;
    }

    .content {
      display: grid;
      gap: 16px;
    }
  `
}

declare global {
  interface HTMLElementTagNameMap {
    'me-admin-panel': MeAdminPanel
  }
}
