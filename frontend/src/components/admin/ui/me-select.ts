import { css, html, LitElement } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import { classMap } from 'lit/directives/class-map.js'

@customElement('me-select')
export class MeSelect extends LitElement {
  @property() label = ''
  @property() value = ''
  @property({ type: Boolean }) disabled = false
  @property({ type: Boolean }) required = false

  private _onChange(e: Event) {
    const select = e.target as HTMLSelectElement
    this.value = select.value
    this.dispatchEvent(
      new CustomEvent('change', {
        detail: select.value,
        bubbles: true,
        composed: true,
      }),
    )
  }

  render() {
    return html`
      <div class=${classMap({ field: true, disabled: this.disabled })}>
        ${this.label ? html`<label class="label">${this.label}</label>` : null}
        <div class="select-wrapper">
          <select
            .value=${this.value}
            ?disabled=${this.disabled}
            ?required=${this.required}
            @change=${this._onChange}
          >
            <slot></slot>
          </select>
        </div>
      </div>
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

    .label {
      font-size: 13px;
      font-weight: 400;
      color: var(--color-text-secondary);
    }

    .select-wrapper {
      position: relative;
      display: grid;
    }

    select {
      width: 100%;
      height: 40px;
      padding: 0 32px 0 12px;
      border: 1px solid var(--color-border);
      background: #ffffff;
      color: var(--color-text-primary);
      font-family: inherit;
      font-size: 14px;
      border-radius: 4px;
      transition: border-color 0.2s ease, box-shadow 0.2s ease;
      outline: none;
      appearance: none;
      cursor: pointer;
    }

    .select-wrapper::after {
      content: "";
      position: absolute;
      right: 12px;
      top: 50%;
      width: 0;
      height: 0;
      border-left: 5px solid transparent;
      border-right: 5px solid transparent;
      border-top: 5px solid var(--color-text-tertiary);
      transform: translateY(-50%);
      pointer-events: none;
    }

    select:focus {
      border-color: var(--admin-accent);
      box-shadow: 0 0 0 1px var(--admin-accent);
    }

    select:disabled {
      background: var(--color-bg-deep);
      color: var(--color-text-tertiary);
      cursor: not-allowed;
    }
  `
}

declare global {
  interface HTMLElementTagNameMap {
    'me-select': MeSelect
  }
}
