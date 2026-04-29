import { css, html, LitElement, type PropertyValues } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import { FormAssociatedMixin } from './form-associated-mixin.js'

/**
 * A standard-compliant, form-associated text input component.
 * Fully leverages ElementInternals and delegatesFocus for a native feel.
 * Targets WCAG 2.1 AA compliance.
 */
@customElement('me-text-input')
export class MeTextInput extends FormAssociatedMixin(LitElement) {
  static shadowRootOptions: ShadowRootInit = {
    ...LitElement.shadowRootOptions,
    delegatesFocus: true,
  }

  @property() label = ''
  @property() autocomplete = ''
  @property({ reflect: true }) type:
    | 'text'
    | 'number'
    | 'email'
    | 'password'
    | 'url'
    | 'datetime-local'
    | 'search' = 'text'
  @property({ type: Boolean, reflect: true }) readonly = false
  @property({ reflect: true }) placeholder = ''

  private _inputId = `me-input-${Math.random().toString(36).slice(2, 9)}`

  private _onInput(e: Event) {
    const input = e.target as HTMLInputElement
    this.value = input.value
    this.syncValidity(input)
  }

  protected update(changedProperties: PropertyValues): void {
    const input = this.renderRoot.querySelector('input')
    if (input) {
      this.syncValidity()
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

    *, *::before, *::after {
      box-sizing: border-box;
    }

    .field {
      width: 100%;
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
