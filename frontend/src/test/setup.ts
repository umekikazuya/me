import { vi } from 'vitest'

if (
  typeof HTMLElement !== 'undefined' &&
  !HTMLElement.prototype.attachInternals
) {
  HTMLElement.prototype.attachInternals = function () {
    return {
      setFormValue: vi.fn(),
      setValidity: vi.fn(),
      checkValidity: vi.fn(() => true),
      reportValidity: vi.fn(() => true),
      validationMessage: '',
      willValidate: true,
      validity: {},
      states: new Set(),
      form: null,
      labels: [],
      shadowRoot: this.shadowRoot,
    } as unknown as ElementInternals
  }
}
