import { css } from 'lit'

/** 1行 = .ln。シェルの出力として並べる公開ページ共通の型 */
export const pageStyles = css`
  :host {
    display: block;
    font-family: var(--font-mono);
    font-size: 13px;
    line-height: 1.9;
  }

  a {
    color: inherit;
    text-decoration: none;
  }

  ul {
    list-style: none;
    margin: 0;
    padding: 0;
  }

  .ln {
    display: block;
    min-height: 1.9em;
    margin-inline: -10px;
    padding-inline: 10px;
    border-radius: 3px;
    white-space: pre-wrap;
    overflow-wrap: anywhere;
    /* keep j/k targets clear of the fixed statusline */
    scroll-margin-block: 48px;
    animation: line-in 0.14s ease both;
    animation-delay: calc(var(--i, 0) * 16ms);
  }

  a.ln,
  button.ln {
    width: calc(100% + 20px);
    border: 0;
    background: none;
    color: inherit;
    font: inherit;
    text-align: left;
    cursor: pointer;
  }

  a.ln:hover,
  button.ln:hover,
  a.ln:focus-visible,
  button.ln:focus-visible {
    background: var(--color-bg-hl);
    outline: none;
  }

  .m {
    color: var(--color-text-tertiary);
  }

  .a {
    color: var(--color-accent);
  }

  .jp {
    font-family: var(--font-jp);
  }

  @keyframes line-in {
    from {
      opacity: 0;
    }
  }

  @media (prefers-reduced-motion: reduce) {
    .ln {
      animation: none;
    }
  }
`
