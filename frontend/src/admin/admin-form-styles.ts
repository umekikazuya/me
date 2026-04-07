import { css } from 'lit'

export const adminFormStyles = css`
  /* ── Page header ─────────────────────────────── */
  .eyebrow {
    font-family: var(--font-en);
    font-size: 13px;
    letter-spacing: var(--tracking-wider);
    color: var(--admin-accent);
    margin-bottom: 10px;
  }

  .title {
    font-size: 28px;
    font-weight: 500;
    color: var(--color-text-primary);
    margin-bottom: 10px;
  }

  .description {
    font-size: 14px;
    color: var(--color-text-secondary);
    line-height: 1.8;
  }

  /* ── Field ───────────────────────────────────── */
  .field {
    display: grid;
    gap: 8px;
    color: var(--color-text-secondary);
    font-size: 14px;
  }

  .field-wide {
    grid-column: 1 / -1;
  }

  /* ── Form controls ───────────────────────────── */
  input,
  textarea {
    width: 100%;
    border: 1px solid var(--color-border);
    background: #fff;
    padding: 10px 12px;
    color: var(--color-text-primary);
    font: inherit;
    font-size: 14px;
    transition: border-color 0.15s ease, box-shadow 0.15s ease;
  }

  input:focus,
  textarea:focus {
    outline: none;
    border-color: var(--admin-accent);
    box-shadow: 0 0 0 3px rgba(0, 87, 184, 0.12);
  }

  textarea {
    resize: vertical;
  }

  /* ── Buttons ─────────────────────────────────── */
  button {
    border: 0;
    background: var(--admin-accent);
    color: #fff;
    padding: 10px 18px;
    font: inherit;
    font-size: 14px;
    cursor: pointer;
    transition: background 0.15s ease;
  }

  button:hover:not(:disabled) {
    background: var(--admin-accent-hover);
  }

  button:disabled {
    opacity: 0.5;
    cursor: wait;
  }

  .subtle {
    background: transparent;
    border: 1px solid var(--color-border);
    color: var(--color-text-secondary);
  }

  .subtle:hover:not(:disabled) {
    background: var(--color-surface);
    color: var(--color-text-primary);
  }

  .danger {
    background: transparent;
    border: 1px solid #c0392b;
    color: #c0392b;
  }

  .danger:hover:not(:disabled) {
    background: rgba(192, 57, 43, 0.06);
  }

  /* ── Messages ────────────────────────────────── */
  .message {
    font-size: 14px;
    line-height: 1.7;
    padding: 10px 14px;
    border-left: 2px solid;
  }

  .error {
    color: #9a3f3f;
    background: rgba(154, 63, 63, 0.06);
    border-color: rgba(154, 63, 63, 0.4);
  }

  .success {
    color: #3d7a56;
    background: rgba(61, 122, 86, 0.06);
    border-color: rgba(61, 122, 86, 0.4);
  }

  .notice {
    color: #5a6b85;
    background: rgba(90, 107, 133, 0.08);
    border-color: rgba(90, 107, 133, 0.35);
  }
`
