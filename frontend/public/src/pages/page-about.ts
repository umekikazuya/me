import { consume } from '@lit/context'
import type { components } from '@me/types'
import { css, html, nothing, type TemplateResult } from 'lit'
import { customElement } from 'lit/decorators.js'
import { profileContext } from '../contexts/profile-context.js'
import type { IProfileRepository } from '../domain/ProfileRepository.js'
import { pageStyles } from '../styles/page-styles.js'
import { sanitizeUrl } from '../utils/format.js'
import { Session } from '../utils/session.js'
import { ShellPage } from './shell-page.js'
import '../components/load-spinner.js'

type Profile = components['schemas']['MeResponse']
type Experience = components['schemas']['MeExperience']

const LINE_STEP_MS = 16
const NODE_STEP_MS = 140

export type AboutSection =
  | { kind: 'whoami'; command: string; text: string }
  | { kind: 'skills'; command: string; items: string[] }
  | { kind: 'certs'; command: string; certs: Profile['certifications'] }
  | { kind: 'experience'; command: string; experiences: Experience[] }
  | { kind: 'likes'; command: string; likes: string[] }

/** about で順に打つコマンド。空の項目はコマンドごと除く。 */
export function aboutSections(p: Profile): AboutSection[] {
  const whoami = [
    p.displayName,
    [p.role, p.location].filter(Boolean).join(', '),
  ]
    .filter(Boolean)
    .join(' — ')
  const skills = [...(p.skills ?? [])]
    .sort((a, b) => a.sortOrder - b.sortOrder)
    .flatMap((group) => group.items)
  const certs = p.certifications ?? []
  const experiences = [...(p.experiences ?? [])].sort(
    (a, b) => b.startYear - a.startYear,
  )
  const likes = p.likes ?? []

  const sections: AboutSection[] = [
    { kind: 'whoami', command: 'whoami', text: whoami },
  ]
  if (skills.length > 0)
    sections.push({ kind: 'skills', command: 'ls skills/', items: skills })
  if (certs.length > 0)
    sections.push({ kind: 'certs', command: 'ls certs/', certs })
  if (experiences.length > 0)
    sections.push({
      kind: 'experience',
      command: 'history',
      experiences,
    })
  if (likes.length > 0)
    sections.push({ kind: 'likes', command: 'cat likes.txt', likes })
  return sections
}

/** "2025.06"、月がなければ "2025"。 */
export function certDate(cert: Profile['certifications'][number]): string {
  return cert.month
    ? `${cert.year}.${String(cert.month).padStart(2, '0')}`
    : String(cert.year)
}

/** 出力が流れ終わるまでの時間 (ms)。 */
export function outputDuration(section: AboutSection): number {
  switch (section.kind) {
    case 'whoami':
      return LINE_STEP_MS
    case 'skills':
      return LINE_STEP_MS
    case 'certs':
      return section.certs.length * LINE_STEP_MS
    case 'experience':
      return LINE_STEP_MS + section.experiences.length * NODE_STEP_MS + 300
    case 'likes':
      return section.likes.length * LINE_STEP_MS
  }
}

/** コマンドを順に打ち、出力が流れ終わるたびに次を打つセッション。 */
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

  private session = new Session(this)

  disconnectedCallback() {
    super.disconnectedCallback()
    if (this._profileRepo) {
      this._profileRepo.removeEventListener('change', this._onRepoChange)
    }
  }

  private onCommandTyped(index: number, sections: AboutSection[]) {
    this.session.advance(
      index,
      sections.length,
      outputDuration(sections[index]),
    )
  }

  protected willUpdate() {
    const p = this.profileRepo.profile
    if (p && this.typed && this.session.outputs === 0) {
      this.onCommandTyped(0, aboutSections(p))
    }
  }

  render() {
    const p = this.profileRepo.profile
    const sections: AboutSection[] = p
      ? aboutSections(p)
      : [{ kind: 'whoami', command: 'whoami', text: '' }]
    const done = p !== null && this.session.outputs >= sections.length
    return html`
      <h1 class="sr-only">about</h1>
      ${sections.slice(0, this.session.commands).map(
        (section, index) => html`
          <section class="block">
            <shell-command
              command=${section.command}
              @typed=${(e: Event) => {
                if (index === 0) return
                e.stopPropagation()
                this.onCommandTyped(index, sections)
              }}
            ></shell-command>
            ${
              index < this.session.outputs
                ? this.renderOutput(section)
                : index === 0 && this.typed
                  ? this.renderPending()
                  : nothing
            }
          </section>
        `,
      )}
      ${
        done
          ? html`<shell-prompt style="animation-delay:${outputDuration(sections[sections.length - 1])}ms"></shell-prompt>`
          : nothing
      }
    `
  }

  private renderPending() {
    return this.profileRepo.error
      ? html`<p class="ln m">whoami: profile unavailable</p>`
      : html`<load-spinner class="ln"></load-spinner>`
  }

  private renderOutput(section: AboutSection): TemplateResult {
    switch (section.kind) {
      case 'whoami':
        return html`<p class="ln jp">${section.text}</p>`
      case 'skills':
        return html`
          <ul class="ln ls">
            ${section.items.map((item) => html`<li>${item}</li>`)}
          </ul>
        `
      case 'certs':
        return html`
          <ul>
            ${section.certs.map(
              (cert, i) =>
                html`<li class="ln" style="--i:${i}">${cert.name}<span class="m">  ${certDate(cert)}</span></li>`,
            )}
          </ul>
        `
      case 'experience':
        return this.renderGraph(section.experiences)
      case 'likes':
        return html`
          <ul>
            ${section.likes.map((like, i) => html`<li class="ln jp" style="--i:${i}">${like}</li>`)}
          </ul>
        `
    }
  }

  private renderGraph(experiences: Experience[]) {
    return html`
      <ol
        class="graph"
        style="--start:${LINE_STEP_MS}ms; --step:${NODE_STEP_MS}ms; --n:${experiences.length}"
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
      p {
        margin: 0;
      }


      .ls {
        display: flex;
        flex-wrap: wrap;
        gap: 0 4ch;
      }

      .graph {
        --rail-x: 4px;
        --draw: calc(var(--n) * var(--step));
        position: relative;
        margin: 0;
        padding: 0;
        list-style: none;
      }

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
