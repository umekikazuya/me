import { css, html, LitElement } from 'lit'
import { customElement, property } from 'lit/decorators.js'

@customElement('me-select')
export class MeSelect extends LitElement {
  static formAssociated = true

  @property({ reflect: true }) label = ''
  @property({ reflect: true }) name = ''
  @property({ reflect: true }) value = ''
  @property({ type: Boolean, reflect: true }) disabled = false
  @property({ type: Boolean, reflect: true }) required = false

  private _internals: ElementInternals
  private _inputId = `me-input-${Math.random().toString(36).slice(2, 9)}`

  constructor() {
    super()
    this._internals = this.attachInternals()
  }

  protected createRenderRoot() {
    return this.attachShadow({ mode: 'open', delegatesFocus: true })
  }

  formResetCallback() {
    this.value = ''
    this._syncInternals()
  }

  formDisabledCallback(disabled: boolean) {
    this.disabled = disabled
  }

  private _onChange(e: Event) {
    const select = e.target as HTMLSelectElement
    this.value = select.value
    this._syncInternals()

    this.dispatchEvent(
      new CustomEvent('change', {
        detail: select.value,
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
    const select = this.shadowRoot?.querySelector('select')
    if (select) {
      this._internals.setValidity(
        select.validity,
        select.validationMessage,
        select,
      )
      this._internals.ariaInvalid = select.checkValidity() ? 'false' : 'true'
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
        <div class="select-wrapper">
          <select
            id=${this._inputId}
            part="select"
            .name=${this.name}
            .value=${this.value ?? ''}
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
      cursor: pointer;
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

    :host([disabled]) select {
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
