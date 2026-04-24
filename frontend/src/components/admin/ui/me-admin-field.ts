import { css, html, LitElement } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import { classMap } from 'lit/directives/class-map.js'

@customElement('me-admin-field')
export class MeAdminField extends LitElement {
  @property() label = ''
  @property({ type: Boolean }) wide = false

  render() {
    return html`
      <label class=${classMap({ field: true, 'field-wide': this.wide })}>
        <span class="label-text">${this.label}</span>
        <div class="input-container">
          <slot></slot>
        </div>
      </label>
    `
  }

  static styles = css`
    :host {
      display: block;
    }

    .field {
      display: grid;
      gap: 6px;
    }

    .label-text {
      font-size: 13px;
      font-weight: 400;
      color: var(--color-text-secondary);
    }

    .input-container {
      display: grid;
    }

    /* Support for grid layout in parents */
    :host([wide]) {
      grid-column: 1 / -1;
    }
  `
}

declare global {
  interface HTMLElementTagNameMap {
    'me-admin-field': MeAdminField
  }
}
