import { css, html, LitElement } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import { classMap } from 'lit/directives/class-map.js'

@customElement('me-text-input')
export class MeTextInput extends LitElement {
  static formAssociated = true

  @property() label = ''
  @property() name = ''
  @property() autocomplete = ''
  @property() value: string | number = ''
  @property() type:
    | 'text'
    | 'number'
    | 'email'
    | 'password'
    | 'url'
    | 'datetime-local'
    | 'search' = 'text'
  @property({ type: Boolean }) disabled = false
  @property({ type: Boolean }) required = false
  @property({ type: Boolean }) readonly = false
  @property() placeholder = ''

  private _internals: ElementInternals
  private _inputId = `me-input-${Math.random().toString(36).slice(2, 9)}`

  constructor() {
    super()
    this._internals = this.attachInternals()
  }

  focus(options?: FocusOptions) {
    this.shadowRoot?.querySelector('input')?.focus(options)
  }

  // Native form callbacks
  formResetCallback() {
    this.value = ''
    this._internals.setFormValue('')
  }

  formDisabledCallback(disabled: boolean) {
    this.disabled = disabled
  }

  private _onInput(e: Event) {
    const input = e.target as HTMLInputElement
    this.value = input.value
    this._internals.setFormValue(input.value)

    // We still dispatch a change event for convenience,
    // but the form now sees the value automatically.
    this.dispatchEvent(
      new CustomEvent('change', {
        detail: input.value,
        bubbles: true,
        composed: true,
      }),
    )
  }

  updated(changedProperties: Map<PropertyKey, unknown>) {
    if (changedProperties.has('value')) {
      this._internals.setFormValue(String(this.value ?? ''))
    }
  }

  render() {
    return html`
      <div class=${classMap({ field: true, disabled: this.disabled })}>
        ${
          this.label
            ? html`<label class="label" for=${this._inputId}>${this.label}</label>`
            : null
        }
        <input
          id=${this._inputId}
          .type=${this.type}
          .name=${this.name}
          .autocomplete=${this.autocomplete}
          .value=${String(this.value ?? '')}
          .placeholder=${this.placeholder}
          ?disabled=${this.disabled}
          ?required=${this.required}
          ?readonly=${this.readonly}
          @input=${this._onInput}
        />
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
      cursor: pointer;
    }

    input {
      width: 100%;
      height: 40px;
      padding: 0 12px;
      border: 1px solid var(--color-border);
      background: #ffffff;
      color: var(--color-text-primary);
      font-family: inherit;
      font-size: 14px;
      border-radius: 4px;
      transition: border-color 0.2s ease, box-shadow 0.2s ease;
      outline: none;
    }

    input:focus {
      border-color: var(--admin-accent);
      box-shadow: 0 0 0 1px var(--admin-accent);
    }

    input:disabled {
      background: var(--color-bg-deep);
      color: var(--color-text-tertiary);
      cursor: not-allowed;
    }

    input[readonly] {
      background: var(--color-bg-dim);
      border-color: var(--color-border-subtle);
    }
  `
}

declare global {
  interface HTMLElementTagNameMap {
    'me-text-input': MeTextInput
  }
}
