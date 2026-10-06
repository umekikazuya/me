import { consume } from '@lit/context'
import type { components } from '@me/types'
import { css, html, nothing } from 'lit'
import { customElement } from 'lit/decorators.js'
import { profileContext } from '../contexts/profile-context.js'
import type { IProfileRepository } from '../domain/ProfileRepository.js'
import { pageStyles } from '../styles/page-styles.js'
import { sanitizeUrl } from '../utils/format.js'
import { ShellPage } from './shell-page.js'
import '../components/load-spinner.js'

type Profile = components['schemas']['MeResponse']
type Experience = components['schemas']['MeExperience']

const LINE_STEP = 16 // ms between streamed lines; matches pageStyles
const NODE_STEP = 140 // ms for the graph rail to travel one commit

/**
 * whoami の「ラベル / 値」の行。スキルはグループごとに1行で、ラベルは先頭行だけに付ける。
 * API が省略した空の配列や空の値は行ごと出さない。
 */
export function profileRows(p: Profile): Array<[string, string]> {
  const skills = [...(p.skills ?? [])]
    .sort((a, b) => a.sortOrder - b.sortOrder)
    .map((group) => group.items.join(', '))
    .filter(Boolean)
  const certs = (p.certifications ?? []).map((c) => c.name).join(', ')
  const likes = (p.likes ?? []).join(', ')
  return [
    ...skills.map((value, n): [string, string] => [
      n === 0 ? 'skills' : '',
      value,
    ]),
    ['certs', certs],
    ['likes', likes],
  ].filter(([, value]) => value) as Array<[string, string]>
}

/** `whoami` — 肩書き、スキルなどの一覧、経歴 */
@customElement('page-about')
export class PageAbout extends ShellPage {
  @consume({ context: profileContext, subscribe: true })
  set profileRepo(repo: IProfileRepository) {
    if (this._profileRepo) {
      this._profileRepo.removeEventListener('change', this._onRepoChange)
    }
    this._profileRepo = repo
    repo.addEventListener('change', this._onRepoChange)
    this.requestUpdate()
  }
  get profileRepo() {
    return this._profileRepo
  }
  private _profileRepo!: IProfileRepository
  private _onRepoChange = () => this.requestUpdate()

  disconnectedCallback() {
    super.disconnectedCallback()
    if (this._profileRepo) {
      this._profileRepo.removeEventListener('change', this._onRepoChange)
    }
  }

  render() {
    return html`
      <h1 class="sr-only">about</h1>
      <shell-command command="whoami"></shell-command>
      ${this.typed ? this.renderOutput() : nothing}
    `
  }

  private renderOutput() {
    const p = this.profileRepo.profile
    if (!p) {
      return this.profileRepo.error
        ? html`<p class="ln m">whoami: profile unavailable</p>
            <shell-prompt></shell-prompt>`
        : html`<load-spinner class="ln"></load-spinner>
            <shell-prompt></shell-prompt>`
    }

    const rows = profileRows(p)
    const headline = [p.role, p.location].filter(Boolean).join(' — ')
    const experiences = [...(p.experiences ?? [])].sort(
      (a, b) => b.startYear - a.startYear,
    )
    let i = 0
    return html`
      ${headline ? html`<p class="ln jp" style="--i:${i++}">${headline}</p>` : nothing}
      ${
        rows.length > 0
          ? html`
            <p class="ln" style="--i:${i++}"></p>
            <dl>
              ${rows.map(
                ([label, value]) => html`
                  <div class="ln kv" style="--i:${i++}">
                    <dt class="m">${label}</dt>
                    <dd class="jp">${value}</dd>
                  </div>
                `,
              )}
            </dl>
          `
          : nothing
      }
      ${experiences.length > 0 ? this.renderGraph(experiences, i++) : nothing}
      <shell-prompt
        style="animation-delay:${i * LINE_STEP + experiences.length * NODE_STEP}ms"
      ></shell-prompt>
    `
  }

  /**
   * 経歴を git log --graph に見立てる。縦線が上から伸び、届いた順にコミット（会社）が現れる。
   * 在籍中は HEAD として塗りつぶす。
   */
  private renderGraph(experiences: Experience[], start: number) {
    return html`
      <p class="ln"></p>
      <h2 class="ln m" style="--i:${start}">experience</h2>
      <ol
        class="graph"
        style="--start:${(start + 1) * LINE_STEP}ms; --step:${NODE_STEP}ms; --n:${experiences.length}"
      >
        ${experiences.map((exp, n) => {
          const head = exp.endYear === undefined || exp.endYear === null
          const body = html`<span class="node" aria-hidden="true"></span><span class="m">${exp.startYear} — ${exp.endYear ?? 'now'}</span><span class="jp">${exp.company}${
            head ? html` <span class="ref">(HEAD)</span>` : nothing
          }</span>`
          return html`
            <li class=${head ? 'commit head' : 'commit'} style="--n:${n}">
              ${
                exp.url
                  ? html`<a class="ln row" href=${sanitizeUrl(exp.url)} target="_blank" rel="noopener noreferrer">${body}</a>`
                  : html`<p class="ln row">${body}</p>`
              }
            </li>
          `
        })}
      </ol>
    `
  }

  static styles = [
    pageStyles,
    css`
      p,
      h2,
      dl,
      dd {
        margin: 0;
        font: inherit;
      }

      /* label / value の2列。値が折り返しても値の列に揃う */
      .kv {
        display: grid;
        grid-template-columns: 13ch minmax(0, 1fr);
      }

      /* ---- experience graph ---- */
      .graph {
        --rail-x: 4px;
        --draw: calc(var(--n) * var(--step));
        position: relative;
        margin: 0;
        padding: 0;
        list-style: none;
      }

      /* the rail: grows downward from the first commit to the last */
      .graph::before {
        content: '';
        position: absolute;
        left: var(--rail-x);
        top: 0.95em;
        bottom: 0.95em;
        width: 1px;
        background: var(--color-text-tertiary);
        transform-origin: top;
        animation: rail var(--draw) linear both;
        animation-delay: var(--start);
      }

      .row {
        display: grid;
        grid-template-columns: 22px 13ch minmax(0, 1fr);
        align-items: center;
        animation: line-in 0.3s ease both;
        animation-delay: calc(var(--start) + var(--n) * var(--step));
      }

      .node {
        position: relative;
        z-index: 1;
        width: 9px;
        height: 9px;
        margin-left: calc(var(--rail-x) - 4px);
        border: 1px solid var(--color-text-tertiary);
        border-radius: 50%;
        background: var(--color-bg-deep);
        animation: pop 0.35s var(--easing-rise) both;
        animation-delay: calc(var(--start) + var(--n) * var(--step));
        transition: transform 0.2s ease;
      }

      .row:hover .node,
      .row:focus-visible .node {
        transform: scale(1.35);
      }

      .head .node {
        border-color: var(--color-accent);
        background: var(--color-accent);
      }

      /* HEAD: one ripple once it lands */
      .head .node::after {
        content: '';
        position: absolute;
        inset: -1px;
        border-radius: 50%;
        border: 1px solid var(--color-accent);
        opacity: 0;
        animation: ripple 1.2s ease-out both;
        animation-delay: calc(var(--start) + var(--n) * var(--step) + 200ms);
      }

      .ref {
        color: var(--color-accent);
      }

      @keyframes rail {
        from {
          transform: scaleY(0);
        }
      }

      @keyframes pop {
        from {
          transform: scale(0);
        }
      }

      @keyframes ripple {
        0% {
          opacity: 0.8;
          transform: scale(1);
        }
        100% {
          opacity: 0;
          transform: scale(3.2);
        }
      }

      @media (prefers-reduced-motion: reduce) {
        .graph::before,
        .row,
        .node,
        .head .node::after {
          animation: none;
        }
      }
    `,
  ]
}

declare global {
  interface HTMLElementTagNameMap {
    'page-about': PageAbout
  }
}
