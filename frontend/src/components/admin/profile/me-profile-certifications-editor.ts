import { css, html, LitElement } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import { adminFormStyles } from '../../../admin/admin-form-styles.js'
import type { MeCertification } from '../../../admin/types.js'
import '../ui/me-admin-panel.js'
import '../ui/me-admin-section.js'
import '../ui/me-text-input.js'

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
                      <me-text-input
                        label="資格名"
                        .value=${cert.name}
                        @change=${(e: CustomEvent) =>
                          this.updateItem(index, { name: e.detail })}
                      ></me-text-input>

                      <me-text-input
                        label="Issuer"
                        .value=${cert.issuer}
                        @change=${(e: CustomEvent) =>
                          this.updateItem(index, { issuer: e.detail })}
                      ></me-text-input>

                      <me-text-input
                        label="Year"
                        type="number"
                        .value=${String(cert.year)}
                        @change=${(e: CustomEvent) =>
                          this.updateItem(index, {
                            year: Number(e.detail || '0'),
                          })}
                      ></me-text-input>

                      <me-text-input
                        label="Month"
                        type="number"
                        .value=${cert.month ? String(cert.month) : ''}
                        @change=${(e: CustomEvent) =>
                          this.updateItem(index, {
                            month: e.detail ? Number(e.detail) : undefined,
                          })}
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
    'me-profile-certifications-editor': MeProfileCertificationsEditor
  }
}
