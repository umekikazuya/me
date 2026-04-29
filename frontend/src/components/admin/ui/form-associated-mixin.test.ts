import { afterEach, beforeEach, describe, expect, it } from 'vitest'
import { customElement } from 'lit/decorators.js'
import { FormAssociatedMixin } from './form-associated-mixin.ts'
import { html, LitElement } from 'lit'

@customElement('test-form-associated')
class TestFormAssociated extends FormAssociatedMixin(LitElement) {
  get internalsForTest() {
    return this._internals
  }
  get input() {
    return this.renderRoot.querySelector('input')!
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
