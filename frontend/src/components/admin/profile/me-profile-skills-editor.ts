import { css, html, LitElement } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import { adminFormStyles } from '../../../admin/admin-form-styles.js'
import type { MeSkillGroup } from '../../../admin/types.js'
import '../ui/me-admin-panel.js'
import '../ui/me-admin-section.js'
import '../ui/me-text-input.js'
import '../ui/me-textarea.js'

@customElement('me-profile-skills-editor')
export class MeProfileSkillsEditor extends LitElement {
  static formAssociated = true

  @property({ type: Array }) skills: MeSkillGroup[] = []
  @property() name = ''

  private _internals: ElementInternals

  constructor() {
    super()
    this._internals = this.attachInternals()
  }

  private dispatchChange(nextSkills: MeSkillGroup[]) {
    this.skills = nextSkills
    this._internals.setFormValue(JSON.stringify(nextSkills))

    this.dispatchEvent(
      new CustomEvent<MeSkillGroup[]>('change', {
        detail: nextSkills,
        bubbles: true,
        composed: true,
      }),
    )
  }

  updated(changedProperties: Map<PropertyKey, unknown>) {
    if (changedProperties.has('skills')) {
      this._internals.setFormValue(JSON.stringify(this.skills ?? []))
    }
  }

  private addSkill = () => {
    const nextSkills = [
      ...this.skills,
      { category: '', items: [], sortOrder: this.skills.length },
    ]
    this.dispatchChange(nextSkills)
  }

  private removeSkill(index: number) {
    const nextSkills = this.skills.filter((_, i) => i !== index)
    this.dispatchChange(nextSkills)
  }

  private updateSkill(index: number, patch: Partial<MeSkillGroup>) {
    const nextSkills = [...this.skills]
    nextSkills[index] = { ...nextSkills[index], ...patch }
    this.dispatchChange(nextSkills)
  }

  private splitLines(value: string) {
    return value
      .split('\n')
      .map((item) => item.trim())
      .filter(Boolean)
  }

  render() {
    return html`
      <me-admin-section
        title="Skills"
        description="カテゴリごとに整理し、Items は1行ずつ入力すると編集しやすいです。"
      >
        <button
          slot="header-actions"
          type="button"
          class="subtle"
          @click=${this.addSkill}
        >
          カテゴリを追加
        </button>

        <div class="stack">
          ${
            this.skills.length === 0
              ? html`<p class="empty-text">まだ skill カテゴリがありません。</p>`
              : this.skills.map(
                  (skill, index) => html`
                  <me-admin-panel title="カテゴリ ${index + 1}">
                    <button
                      slot="header-actions"
                      type="button"
                      class="subtle danger"
                      @click=${() => this.removeSkill(index)}
                    >
                      削除
                    </button>

                    <div class="grid">
                      <me-text-input
                        label="Category"
                        .value=${skill.category}
                        @change=${(e: CustomEvent) =>
                          this.updateSkill(index, { category: e.detail })}
                      ></me-text-input>
                      
                      <me-text-input
                        label="Sort order"
                        type="number"
                        .value=${String(skill.sortOrder)}
                        @change=${(e: CustomEvent) =>
                          this.updateSkill(index, {
                            sortOrder: Number(e.detail || '0'),
                          })}
                      ></me-text-input>

                      <me-textarea
                        label="Items（改行区切り）"
                        .value=${skill.items.join('\n')}
                        rows="4"
                        class="field-wide"
                        @change=${(e: CustomEvent) =>
                          this.updateSkill(index, {
                            items: this.splitLines(e.detail),
                          })}
                      ></me-textarea>
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

      .field-wide {
        grid-column: 1 / -1;
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
    'me-profile-skills-editor': MeProfileSkillsEditor
  }
}
