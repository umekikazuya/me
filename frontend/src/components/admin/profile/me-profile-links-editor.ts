import { css, html, LitElement } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import { adminFormStyles } from '../../../admin/admin-form-styles.js'
import type { MeLink } from '../../../admin/types.js'
import '../ui/me-admin-panel.js'
import '../ui/me-admin-section.js'
import '../ui/me-text-input.js'

@customElement('me-profile-links-editor')
export class MeProfileLinksEditor extends LitElement {
  static formAssociated = true

  @property({ type: Array }) links: MeLink[] = []
  @property() name = ''

  private _internals: ElementInternals

  constructor() {
    super()
    this._internals = this.attachInternals()
  }

  private dispatchChange(next: MeLink[]) {
    this.links = next
    this._internals.setFormValue(JSON.stringify(next))

    this.dispatchEvent(
      new CustomEvent<MeLink[]>('change', {
        detail: next,
        bubbles: true,
        composed: true,
      }),
    )
  }

  updated(changedProperties: Map<PropertyKey, unknown>) {
    if (changedProperties.has('links')) {
      this._internals.setFormValue(JSON.stringify(this.links ?? []))
    }
  }

  private addItem = () => {
    const next = [...this.links, { platform: '', url: '' }]
    this.dispatchChange(next)
  }

  private removeItem(index: number) {
    const next = this.links.filter((_, i) => i !== index)
    this.dispatchChange(next)
  }

  private updateItem(index: number, patch: Partial<MeLink>) {
    const next = [...this.links]
    next[index] = { ...next[index], ...patch }
    this.dispatchChange(next)
  }

  render() {
    return html`
      <me-admin-section
        title="Links"
        description="platform と URL は必須です。label は公開側で見せたい名前を指定します。"
      >
        <button
          slot="header-actions"
          type="button"
          class="subtle"
          @click=${this.addItem}
        >
          リンクを追加
        </button>

        <div class="stack">
          ${
            this.links.length === 0
              ? html`<p class="empty-text">リンクがまだありません。</p>`
              : this.links.map(
                  (link, index) => html`
                  <me-admin-panel title="リンク ${index + 1}">
                    <button
                      slot="header-actions"
                      type="button"
                      class="subtle danger"
                      @click=${() => this.removeItem(index)}
                    >
                      削除
                    </button>

                    <div class="grid">
                      <me-text-input
                        label="Platform"
                        .value=${link.platform}
                        @change=${(e: CustomEvent) =>
                          this.updateItem(index, { platform: e.detail })}
                      ></me-text-input>

                      <me-text-input
                        label="URL"
                        type="url"
                        .value=${link.url}
                        @change=${(e: CustomEvent) =>
                          this.updateItem(index, { url: e.detail })}
                      ></me-text-input>
                    </div>
                  </me-admin-panel>
                `,
                )
          }
        </div>
      </me-admin-section>
    `
  }

  static styles = [
    adminFormStyles,
    css`
      :host {
        display: block;
      }

      .stack {
        display: grid;
        gap: 20px;
      }

      .grid {
        display: grid;
        gap: 16px;
        grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
      }

      .empty-text {
        color: var(--color-text-tertiary);
        font-size: 14px;
        font-style: italic;
      }
    `,
  ]
}

declare global {
  interface HTMLElementTagNameMap {
    'me-profile-links-editor': MeProfileLinksEditor
  }
}
