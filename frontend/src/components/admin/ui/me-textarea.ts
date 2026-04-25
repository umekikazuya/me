import { css, html, LitElement } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import { classMap } from 'lit/directives/class-map.js'

@customElement('me-textarea')
export class MeTextarea extends LitElement {
  @property() label = ''
  @property() value = ''
  @property({ type: Number }) rows = 4
  @property({ type: Boolean }) disabled = false
  @property({ type: Boolean }) required = false
  @property() placeholder = ''

  private _onInput(e: Event) {
    const input = e.target as HTMLTextAreaElement
    this.value = input.value
    this.dispatchEvent(
      new CustomEvent('change', {
        detail: input.value,
        bubbles: true,
        composed: true,
      }),
    )
  }

  render() {
    return html`
      <div class=${classMap({ field: true, disabled: this.disabled })}>
        ${this.label ? html`<label class="label">${this.label}</label>` : null}
        <textarea
          .rows=${this.rows}
          .value=${this.value}
          .placeholder=${this.placeholder}
          ?disabled=${this.disabled}
          ?required=${this.required}
          @input=${this._onInput}
        ></textarea>
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

    textarea {
      width: 100%;
      padding: 10px 12px;
      border: 1px solid var(--color-border);
      background: #ffffff;
      color: var(--color-text-primary);
      font-family: inherit;
      font-size: 14px;
      line-height: 1.6;
      border-radius: 4px;
      transition: border-color 0.2s ease, box-shadow 0.2s ease;
      outline: none;
      resize: vertical;
    }

    textarea:focus {
      border-color: var(--admin-accent);
      box-shadow: 0 0 0 1px var(--admin-accent);
    }

    textarea:disabled {
      background: var(--color-bg-deep);
      color: var(--color-text-tertiary);
      cursor: not-allowed;
    }
  `
}

declare global {
  interface HTMLElementTagNameMap {
    'me-textarea': MeTextarea
  }
}
