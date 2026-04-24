import { css, html, LitElement } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import { adminFormStyles } from '../../../admin/admin-form-styles.js'
import type { MeCertification } from '../../../admin/types.js'
import '../ui/me-admin-field.js'
import '../ui/me-admin-panel.js'
import '../ui/me-admin-section.js'

@customElement('me-profile-certifications-editor')
export class MeProfileCertificationsEditor extends LitElement {
  @property({ type: Array }) certifications: MeCertification[] = []

  private dispatchChange(next: MeCertification[]) {
    this.dispatchEvent(
      new CustomEvent<MeCertification[]>('change', {
        detail: next,
        bubbles: true,
        composed: true,
      }),
    )
  }

  private addItem = () => {
    const next = [
      ...this.certifications,
      { name: '', issuer: '', year: new Date().getFullYear() },
    ]
    this.dispatchChange(next)
  }

  private removeItem(index: number) {
    const next = this.certifications.filter((_, i) => i !== index)
    this.dispatchChange(next)
  }

  private updateItem(index: number, patch: Partial<MeCertification>) {
    const next = [...this.certifications]
    next[index] = { ...next[index], ...patch }
    this.dispatchChange(next)
  }

  private toOptionalNumber(value: string) {
    const trimmed = value.trim()
    return trimmed === '' ? undefined : Number(trimmed)
  }

  render() {
    return html`
      <me-admin-section
        title="Certifications"
        description="month は任意です。年だけでも掲載できます。"
      >
        <button
          slot="header-actions"
          type="button"
          class="subtle"
          @click=${this.addItem}
        >
          資格を追加
        </button>

        <div class="stack">
          ${
            this.certifications.length === 0
              ? html`<p class="empty-text">資格がまだありません。</p>`
              : this.certifications.map(
                  (cert, index) => html`
                  <me-admin-panel title="資格 ${index + 1}">
                    <button
                      slot="header-actions"
                      type="button"
                      class="subtle danger"
                      @click=${() => this.removeItem(index)}
                    >
                      削除
                    </button>

                    <div class="grid">
                      <me-admin-field label="資格名">
                        <input
                          .value=${cert.name}
                          @input=${(e: Event) =>
                            this.updateItem(index, {
                              name: (e.target as HTMLInputElement).value,
                            })}
                        />
                      </me-admin-field>
                      <me-admin-field label="Issuer">
                        <input
                          .value=${cert.issuer}
                          @input=${(e: Event) =>
                            this.updateItem(index, {
                              issuer: (e.target as HTMLInputElement).value,
                            })}
                        />
                      </me-admin-field>
                      <me-admin-field label="Year">
                        <input
                          type="number"
                          .value=${String(cert.year)}
                          @input=${(e: Event) =>
                            this.updateItem(index, {
                              year: Number(
                                (e.target as HTMLInputElement).value || '0',
                              ),
                            })}
                        />
                      </me-admin-field>
                      <me-admin-field label="Month">
                        <input
                          type="number"
                          min="1"
                          max="12"
                          .value=${cert.month ? String(cert.month) : ''}
                          @input=${(e: Event) =>
                            this.updateItem(index, {
                              month: this.toOptionalNumber(
                                (e.target as HTMLInputElement).value,
                              ),
                            })}
                        />
                      </me-admin-field>
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
    'me-profile-certifications-editor': MeProfileCertificationsEditor
  }
}
