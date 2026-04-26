import { css, html, LitElement } from 'lit'
import { customElement, property } from 'lit/decorators.js'

/**
 * A standard-compliant, form-associated text input component.
 * Fully leverages ElementInternals and delegatesFocus for a native feel.
 */
@customElement('me-text-input')
export class MeTextInput extends LitElement {
  static formAssociated = true

  @property({ reflect: true }) label = ''
  @property({ reflect: true }) name = ''
  @property({ reflect: true }) autocomplete = ''
  @property({ reflect: true }) value: string | number = ''
  @property({ reflect: true }) type:
    | 'text'
    | 'number'
    | 'email'
    | 'password'
    | 'url'
    | 'datetime-local'
    | 'search' = 'text'
  @property({ type: Boolean, reflect: true }) disabled = false
  @property({ type: Boolean, reflect: true }) required = false
  @property({ type: Boolean, reflect: true }) readonly = false
  @property({ reflect: true }) placeholder = ''

  private _internals: ElementInternals
  private _inputId = `me-input-${Math.random().toString(36).slice(2, 9)}`

  constructor() {
    super()
    this._internals = this.attachInternals()
  }

  /**
   * Overrides createRenderRoot to enable focus delegation.
   * This means focusing the <me-text-input> element will automatically
   * focus the inner <input> tag.
   */
  protected createRenderRoot() {
    return this.attachShadow({ mode: 'open', delegatesFocus: true })
  }

  // --- Form Association Callbacks ---

  formResetCallback() {
    this.value = ''
    this._internals.setFormValue('')
  }

  formDisabledCallback(disabled: boolean) {
    this.disabled = disabled
  }

  // --- Validation ---

  checkValidity() {
    return this._internals.checkValidity()
  }

  reportValidity() {
    return this._internals.reportValidity()
  }

  private _onInput(e: Event) {
    const input = e.target as HTMLInputElement
    this.value = input.value
    this._syncInternals()

    // Bubbles the event as a standard 'change' event for parent listeners
    this.dispatchEvent(
      new CustomEvent('change', {
        detail: input.value,
        bubbles: true,
        composed: true,
      }),
    )
  }

  protected updated(changedProperties: Map<PropertyKey, unknown>) {
    if (changedProperties.has('value')) {
      this._syncInternals()
    }
  }

  private _syncInternals() {
    const val = String(this.value ?? '')
    this._internals.setFormValue(val)

    // Simple native validation sync
    const input = this.shadowRoot?.querySelector('input')
    if (input) {
      this._internals.setValidity(
        input.validity,
        input.validationMessage,
        input,
      )
    }
  }

  render() {
    return html`
      <div class="field" ?disabled=${this.disabled}>
        ${
          this.label
            ? html`<label class="label" for=${this._inputId} part="label">${this.label}</label>`
            : null
        }
        <input
          id=${this._inputId}
          part="input"
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

    /* Standard states via CSS classes or attributes */
    :host([disabled]) input {
      background: var(--color-bg-deep);
      color: var(--color-text-tertiary);
      cursor: not-allowed;
    }

    input[readonly] {
      background: var(--color-bg-dim);
      border-color: var(--color-border-subtle);
    }

    /* Validation styles */
    input:invalid:not(:placeholder-shown) {
      border-color: var(--color-danger);
    }
  `
}

declare global {
  interface HTMLElementTagNameMap {
    'me-text-input': MeTextInput
  }
}
