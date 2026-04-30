import { css, html, LitElement } from 'lit'
import { customElement, property } from 'lit/decorators.js'

/**
 * A standard-compliant, form-associated textarea component.
 */
@customElement('me-textarea')
export class MeTextarea extends LitElement {
  static formAssociated = true

  @property({ reflect: true }) label = ''
  @property({ reflect: true }) name = ''
  @property({ type: Number, reflect: true }) rows = 4
  @property({ reflect: true }) value = ''
  @property({ type: Boolean, reflect: true }) disabled = false
  @property({ type: Boolean, reflect: true }) required = false
  @property({ reflect: true }) placeholder = ''

  private _internals: ElementInternals
  private _inputId = `me-input-${Math.random().toString(36).slice(2, 9)}`

  constructor() {
    super()
    this._internals = this.attachInternals()
  }

  static shadowRootOptions: ShadowRootInit = {
    ...LitElement.shadowRootOptions,
    delegatesFocus: true,
  }

  formResetCallback() {
    this.value = ''
    this._syncInternals()
  }

  formDisabledCallback(disabled: boolean) {
    this.disabled = disabled
  }

  private _onInput(e: Event) {
    const input = e.target as HTMLTextAreaElement
    this.value = input.value
    this._syncInternals()

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
    this._internals.setFormValue(this.value ?? '')
    const input = this.shadowRoot?.querySelector('textarea')
    if (input) {
      this._internals.setValidity(
        input.validity,
        input.validationMessage,
        input,
      )
      this._internals.ariaInvalid = input.checkValidity() ? 'false' : 'true'
      this._internals.ariaRequired = this.required ? 'true' : 'false'
    }
  }

  render() {
    return html`
      <div class="field">
        ${
          this.label
            ? html`<label class="label" for=${this._inputId} part="label">${this.label}</label>`
            : null
        }
        <textarea
          id=${this._inputId}
          part="textarea"
          .name=${this.name}
          .rows=${this.rows}
          .value=${this.value ?? ''}
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

    *, *::before, *::after {
      box-sizing: border-box;
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

    :host([disabled]) textarea {
      background: var(--color-bg-deep);
      color: var(--color-text-tertiary);
      cursor: not-allowed;
    }

    textarea:invalid:not(:placeholder-shown) {
      border-color: var(--color-danger);
    }
  `
}

declare global {
  interface HTMLElementTagNameMap {
    'me-textarea': MeTextarea
  }
}
