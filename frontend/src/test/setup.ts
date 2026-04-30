import { vi } from 'vitest'

if (
  typeof HTMLElement !== 'undefined' &&
  !HTMLElement.prototype.attachInternals
) {
  HTMLElement.prototype.attachInternals = function () {
    const internals = {
      setFormValue: vi.fn(),
      setValidity: vi.fn(),
      checkValidity: vi.fn(() => true),
      reportValidity: vi.fn(() => true),
      validationMessage: '',
      willValidate: true,
      validity: {} as ValidityState,
      states: new Set(),
      form: null,
      labels: [],
      shadowRoot: this.shadowRoot,
    }
    internals.setValidity = vi.fn(
      (validity: ValidityState) => (internals.validity = validity),
    )
    return internals as unknown as ElementInternals
  }
}
