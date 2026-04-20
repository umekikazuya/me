import { css, html, LitElement } from 'lit'
import { customElement } from 'lit/decorators.js'

@customElement('page-admin-dashboard')
export class PageAdminDashboard extends LitElement {
  render() {
    return html`
      <section class="container">
        <div class="hero">
          <p class="eyebrow" lang="en">Dashboard</p>
          <h1 class="title">管理画面</h1>
          <p class="description">
            プロフィール、記事、アカウントの管理をここから行えます。
          </p>
        </div>

        <div class="cards">
          <a class="card" href="/admin/articles">
            <h2>記事</h2>
            <p>記事の一覧確認、手動登録、更新、削除を行います。</p>
          </a>

          <a class="card" href="/admin/profile">
            <h2>プロフィール</h2>
            <p>公開プロフィールの内容を更新します。</p>
          </a>

          <a class="card" href="/admin/account">
            <h2>アカウント</h2>
            <p>メールアドレス変更やセッション操作を行います。</p>
          </a>
        </div>
      </section>
    `
  }

  static styles = css`
    :host {
      display: block;
    }

    .container {
      display: grid;
      gap: 32px;
    }

    .eyebrow {
      font-family: var(--font-en);
      letter-spacing: var(--tracking-wider);
      color: var(--color-text-tertiary);
      margin-bottom: 12px;
    }

    .title {
      font-weight: 300;
      font-size: 30px;
      margin-bottom: 12px;
    }

    .description {
      color: var(--color-text-secondary);
      line-height: 1.8;
    }

    .cards {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
      gap: 20px;
    }

    .card {
      display: grid;
      gap: 10px;
      padding: 24px;
      border: 1px solid var(--color-border);
      background: #fff;
      transition:
        transform 0.3s var(--easing-smooth),
        border-color 0.3s var(--easing-smooth);
    }

    .card:hover {
      transform: translateY(-4px);
      border-color: var(--color-text-primary);
    }

    h2 {
      font-size: 20px;
      font-weight: 300;
    }

    p {
      color: var(--color-text-secondary);
      line-height: 1.8;
    }
  `
}

declare global {
  interface HTMLElementTagNameMap {
    'page-admin-dashboard': PageAdminDashboard
  }
}
