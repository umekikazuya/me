import { css, html, LitElement } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import { adminFormStyles } from '../../../admin/admin-form-styles.js'
import type { MeExperience } from '../../../admin/types.js'
import '../ui/me-admin-panel.js'
import '../ui/me-admin-section.js'
import '../ui/me-text-input.js'

@customElement('me-profile-experiences-editor')
export class MeProfileExperiencesEditor extends LitElement {
  static formAssociated = true

  @property({ type: Array }) experiences: MeExperience[] = []
  @property() name = ''

  private _internals: ElementInternals

  constructor() {
    super()
    this._internals = this.attachInternals()
  }

  private dispatchChange(next: MeExperience[]) {
    this.experiences = next
    this._internals.setFormValue(JSON.stringify(next))

    this.dispatchEvent(
      new CustomEvent<MeExperience[]>('change', {
        detail: next,
        bubbles: true,
        composed: true,
      }),
    )
  }

  updated(changedProperties: Map<PropertyKey, unknown>) {
    if (changedProperties.has('experiences')) {
      this._internals.setFormValue(JSON.stringify(this.experiences ?? []))
    }
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
                      <me-text-input
                        label="Company"
                        .value=${exp.company}
                        @change=${(e: CustomEvent) =>
                          this.updateItem(index, { company: e.detail })}
                      ></me-text-input>

                      <me-text-input
                        label="URL"
                        type="url"
                        .value=${exp.url}
                        @change=${(e: CustomEvent) =>
                          this.updateItem(index, { url: e.detail })}
                      ></me-text-input>

                      <me-text-input
                        label="Start year"
                        type="number"
                        .value=${String(exp.startYear)}
                        @change=${(e: CustomEvent) =>
                          this.updateItem(index, {
                            startYear: Number(e.detail || '0'),
                          })}
                      ></me-text-input>

                      <me-text-input
                        label="End year"
                        type="number"
                        .value=${exp.endYear ? String(exp.endYear) : ''}
                        @change=${(e: CustomEvent) =>
                          this.updateItem(index, {
                            endYear: e.detail ? Number(e.detail) : undefined,
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
    'me-profile-experiences-editor': MeProfileExperiencesEditor
  }
}
