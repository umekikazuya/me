import type { LitElement } from 'lit'
import { property } from 'lit/decorators.js'

interface FormAssociatedInterface {
  name: string
  disabled: boolean
  required: boolean
  value: string
  formResetCallback(): void
  formDisabledCallback(disabled: boolean): void
}

type Constructor<T = {}> = new (...args: any[]) => T

export const FormAssociatedMixin = <T extends Constructor<LitElement>>(
  superClass: T,
): T & Constructor<FormAssociatedInterface> => {
  class FormAssociatedElement
    extends superClass
    implements FormAssociatedInterface
  {
    static formAssociated = true

    protected readonly _internals: ElementInternals

    @property({ type: String }) name = ''
    @property({ type: Boolean, reflect: true }) disabled = false
    @property({ type: Boolean, reflect: true }) required = false
    @property() value = ''

    constructor(...args: any[]) {
      super(...args)
      this._internals = this.attachInternals()
    }

    formResetCallback() {
      this.value = ''
    }

    formDisabledCallback(disabled: boolean) {
      this.disabled = disabled
    }

    protected syncValidity(
      input: HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement,
    ) {
      this._internals.setValidity(
        input.validity,
        input.validationMessage,
        input,
      )
    }

    protected updated(changedProperties: Map<PropertyKey, unknown>) {
      super.updated?.(changedProperties)
      if (changedProperties.has('value')) {
        this._internals.setFormValue(this.value)
      }
    }
  }
  return FormAssociatedElement as T & Constructor<FormAssociatedInterface>
}
