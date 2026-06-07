import { html, LitElement } from 'lit'
import { customElement } from 'lit/decorators.js'
import { afterEach, beforeEach, describe, expect, it } from 'vitest'
import { FormAssociatedMixin } from './form-associated-mixin.ts'

@customElement('test-form-associated')
class TestFormAssociated extends FormAssociatedMixin(LitElement) {
  get internalsForTest(): ElementInternals {
    return (this as unknown as { _internals: ElementInternals })._internals
  }
  get input() {
    const input = this.renderRoot.querySelector('input')
    if (!input) throw new Error('Input element not found')
    return input
  }
  runSyncValidity() {
    this.syncValidity(this.input)
  }

  render() {
    return html`
    <input
      .value=${this.value}
      ?disabled=${this.disabled}
      ?required=${this.required}
    >
    `
  }
}

describe('FormAssociatedMixin', () => {
  let el: TestFormAssociated

  beforeEach(async () => {
    el = document.createElement('test-form-associated') as TestFormAssociated
    document.body.appendChild(el)
    await el.updateComplete
  })

  afterEach(() => {
    el.remove()
  })

  it('ElementInternals が初期化されていること', () => {
    expect(el.internalsForTest).toBeDefined()
    expect(el.internalsForTest.role).toBeUndefined()
  })

  it('disabled プロパティが内部の input に反映されること', async () => {
    el.disabled = true
    await el.updateComplete
    expect(el.input.disabled).toBe(true)
  })

  it('formResetCallback で値がリセットされること', async () => {
    el.value = 'dirty value'
    await el.updateComplete

    el.formResetCallback()
    expect(el.value).toBe('')
  })

  it('syncValidity でバリデーション状態が同期されること', async () => {
    el.required = true
    el.value = ''
    await el.updateComplete

    el.runSyncValidity()
    expect(el.internalsForTest.validity.valid).toBe(false)

    el.value = 'filled'
    await el.updateComplete
    el.runSyncValidity()
    expect(el.internalsForTest.validity.valid).toBe(true)
  })
})
