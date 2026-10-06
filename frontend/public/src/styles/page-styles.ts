import { css } from 'lit'

/** 公開ページ共通の型。セクション見出し・罫線リスト・補足テキスト */
export const pageStyles = css`
  :host {
    display: block;
  }

  section + section {
    margin-top: 64px;
  }

  h2 {
    margin: 0 0 12px;
    font-family: var(--font-mono);
    font-size: 12px;
    font-weight: 400;
    color: var(--color-text-tertiary);
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

  .rows > li {
    border-top: 1px solid var(--color-border);
  }

  .rows > li:last-child {
    border-bottom: 1px solid var(--color-border);
  }

  .row {
    display: flex;
    justify-content: space-between;
    align-items: baseline;
    gap: 24px;
    padding: 11px 0;
    transition: opacity 0.25s ease;
  }

  .rows:hover a.row {
    opacity: 0.4;
  }

  .rows:hover a.row:hover {
    opacity: 1;
  }

  .row .main {
    min-width: 0;
  }

  .meta {
    font-family: var(--font-mono);
    font-size: 12px;
    color: var(--color-text-tertiary);
    font-variant-numeric: tabular-nums;
    white-space: nowrap;
  }

  .note {
    margin: 12px 0 0;
    font-family: var(--font-mono);
    font-size: 12px;
    color: var(--color-text-tertiary);
  }

  .more {
    display: inline-block;
    margin-top: 14px;
    font-family: var(--font-mono);
    font-size: 12px;
    color: var(--color-text-tertiary);
    transition: color 0.2s ease;
  }

  .more:hover {
    color: var(--color-text-primary);
  }

  a:focus-visible,
  button:focus-visible {
    outline: 1px solid var(--color-accent);
    outline-offset: 3px;
  }
`
