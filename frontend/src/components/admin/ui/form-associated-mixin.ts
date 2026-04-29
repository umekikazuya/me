import type { LitElement } from 'lit'
import { property } from 'lit/decorators.js'

type Constructor<T = {}> = new (...args: any[]) => T

export const FormAssociatedMixin = <T extends Constructor<LitElement>>(
  superClass: T,
) => {
  class FormAssociatedElement extends superClass {
    static formAssociated = true

    protected readonly _internals: ElementInternals

    @property({ type: String }) name = ''
    @property({ type: Boolean, reflect: true }) disabled = false
    @property({ type: Boolean, reflect: true }) required = false
    @property({ reflect: true }) value = ''

    constructor(...args: any[]) {
      super(...args)
      this._internals = this.attachInternals()
    }

    formResetCallback() {
      this.value = ''
      this._internals.setFormValue('')
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
  return FormAssociatedElement as any
}
