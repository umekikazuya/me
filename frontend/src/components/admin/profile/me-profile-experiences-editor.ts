import { css, html, LitElement } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import { adminFormStyles } from '../../../admin/admin-form-styles.js'
import type { MeExperience } from '../../../admin/types.js'
import '../ui/me-admin-field.js'
import '../ui/me-admin-panel.js'
import '../ui/me-admin-section.js'

@customElement('me-profile-experiences-editor')
export class MeProfileExperiencesEditor extends LitElement {
  @property({ type: Array }) experiences: MeExperience[] = []

  private dispatchChange(next: MeExperience[]) {
    this.dispatchEvent(
      new CustomEvent<MeExperience[]>('change', {
        detail: next,
        bubbles: true,
        composed: true,
      }),
    )
  }

  private addItem = () => {
    const next = [
      ...this.experiences,
      { company: '', url: '', startYear: new Date().getFullYear() },
    ]
    this.dispatchChange(next)
  }

  private removeItem(index: number) {
    const next = this.experiences.filter((_, i) => i !== index)
    this.dispatchChange(next)
  }

  private updateItem(index: number, patch: Partial<MeExperience>) {
    const next = [...this.experiences]
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
        title="Experiences"
        description="endYear を空にすると、継続中の経歴として扱えます。"
      >
        <button
          slot="header-actions"
          type="button"
          class="subtle"
          @click=${this.addItem}
        >
          経歴を追加
        </button>

        <div class="stack">
          ${
            this.experiences.length === 0
              ? html`<p class="empty-text">経歴がまだありません。</p>`
              : this.experiences.map(
                  (exp, index) => html`
                  <me-admin-panel title="経歴 ${index + 1}">
                    <button
                      slot="header-actions"
                      type="button"
                      class="subtle danger"
                      @click=${() => this.removeItem(index)}
                    >
                      削除
                    </button>

                    <div class="grid">
                      <me-admin-field label="Company">
                        <input
                          .value=${exp.company}
                          @input=${(e: Event) =>
                            this.updateItem(index, {
                              company: (e.target as HTMLInputElement).value,
                            })}
                        />
                      </me-admin-field>
                      <me-admin-field label="URL">
                        <input
                          .value=${exp.url}
                          @input=${(e: Event) =>
                            this.updateItem(index, {
                              url: (e.target as HTMLInputElement).value,
                            })}
                        />
                      </me-admin-field>
                      <me-admin-field label="Start year">
                        <input
                          type="number"
                          .value=${String(exp.startYear)}
                          @input=${(e: Event) =>
                            this.updateItem(index, {
                              startYear: Number(
                                (e.target as HTMLInputElement).value || '0',
                              ),
                            })}
                        />
                      </me-admin-field>
                      <me-admin-field label="End year">
                        <input
                          type="number"
                          .value=${exp.endYear ? String(exp.endYear) : ''}
                          @input=${(e: Event) =>
                            this.updateItem(index, {
                              endYear: this.toOptionalNumber(
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
    'me-profile-experiences-editor': MeProfileExperiencesEditor
  }
}
