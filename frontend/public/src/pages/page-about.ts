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

/** fastfetch の key: value 行。API が省略した空の配列は行ごと出さない */
export function fetchRows(p: Profile): Array<[string, string]> {
  const rows: Array<[string, string]> = [
    ['Role', p.role],
    ['Location', p.location],
    ...[...(p.skills ?? [])]
      .sort((a, b) => a.sortOrder - b.sortOrder)
      .map((group): [string, string] => [
        group.category,
        group.items.join(', '),
      ]),
  ]
  const certs = (p.certifications ?? []).map((c) => c.name)
  if (certs.length > 0) rows.push(['Certs', certs.join(', ')])
  if ((p.likes ?? []).length > 0) rows.push(['Likes', p.likes.join(', ')])
  return rows.filter(([, value]) => value)
}

/** `fastfetch` と、経歴を `git log --graph` で */
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
      <shell-command command="fastfetch"></shell-command>
      ${this.typed ? this.renderOutput() : nothing}
    `
  }

  private renderOutput() {
    const p = this.profileRepo.profile
    if (!p) {
      return this.profileRepo.error
        ? html`<p class="ln m">fastfetch: profile unavailable</p>
            <shell-nav current="about"></shell-nav>`
        : html`<load-spinner class="ln"></load-spinner>
            <shell-nav current="about"></shell-nav>`
    }

    const rows = fetchRows(p)
    const width = Math.max(...rows.map(([key]) => key.length)) + 2
    const experiences = [...(p.experiences ?? [])].sort(
      (a, b) => b.startYear - a.startYear,
    )
    let i = 0
    return html`
      <p class="ln" style="--i:${i++}"><span class="a">umekikazuya</span><span class="m">@</span><span class="a">me</span></p>
      <p class="ln m" style="--i:${i++}">${'-'.repeat(14)}</p>
      ${rows.map(
        ([key, value]) =>
          html`<p class="ln" style="--i:${i++}"><span class="a">${`${key}:`.padEnd(width)}</span><span class="jp">${value}</span></p>`,
      )}
      <p class="ln" style="--i:${i++}"></p>
      <p class="ln" style="--i:${i++}" aria-hidden="true"><span class="sw fg"></span><span class="sw muted"></span><span class="sw border"></span><span class="sw accent"></span></p>
      ${
        experiences.length > 0
          ? html`
            <p class="ln" style="--i:${i++}"></p>
            <p class="ln" style="--i:${i++}"><span class="m">~ $</span> git log --graph --format="%s"  <span class="m"># experience</span></p>
            <ul>
              ${experiences.map((exp, n) => {
                const line = html`<span class="a">*</span> <span class="jp">${exp.company}</span>  <span class="m">${exp.startYear} — ${exp.endYear ?? 'now'}</span>`
                const edge =
                  n < experiences.length - 1
                    ? html`<p class="ln a" style="--i:${i++}" aria-hidden="true">|</p>`
                    : nothing
                return html`
                  <li>
                    ${
                      exp.url
                        ? html`<a class="ln" style="--i:${i++}" href=${sanitizeUrl(exp.url)} target="_blank" rel="noopener noreferrer">${line}</a>`
                        : html`<p class="ln" style="--i:${i++}">${line}</p>`
                    }
                    ${edge}
                  </li>
                `
              })}
            </ul>
          `
          : nothing
      }
      <shell-nav current="about" style="animation-delay:${i * 16}ms"></shell-nav>
    `
  }

  static styles = [
    pageStyles,
    css`
      p {
        margin: 0;
      }

      .sw {
        display: inline-block;
        width: 2.2em;
        height: 1.1em;
        margin-right: 2px;
        vertical-align: -0.2em;
      }

      .sw.fg {
        background: var(--color-text-primary);
      }

      .sw.muted {
        background: var(--color-text-tertiary);
      }

      .sw.border {
        background: var(--color-border);
      }

      .sw.accent {
        background: var(--color-accent);
      }
    `,
  ]
}

declare global {
  interface HTMLElementTagNameMap {
    'page-about': PageAbout
  }
}
