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
    box-shadow: 0 0 0 3px color-mix(in srgb, var(--admin-accent) 12%, transparent);
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
    cursor: not-allowed;
  }

  button.is-loading {
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
    border: 1px solid var(--color-danger);
    color: var(--color-danger);
  }

  .danger:hover:not(:disabled) {
    background: var(--color-danger-bg);
  }

  /* ── Messages ────────────────────────────────── */
  .message {
    font-size: 14px;
    line-height: 1.7;
    padding: 10px 14px;
    border-left: 2px solid;
  }

  .error {
    color: var(--color-danger);
    background: var(--color-danger-bg);
    border-color: color-mix(in srgb, var(--color-danger) 40%, transparent);
  }

  .success {
    color: var(--color-success);
    background: var(--color-success-bg);
    border-color: color-mix(in srgb, var(--color-success) 40%, transparent);
  }

  .notice {
    color: var(--color-notice);
    background: var(--color-notice-bg);
    border-color: color-mix(in srgb, var(--color-notice) 35%, transparent);
  }
`
