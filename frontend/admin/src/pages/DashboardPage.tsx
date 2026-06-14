import { Link } from 'react-router-dom'

export function DashboardPage() {
  return (
    <section className="container">
      <header
        className="page-header"
        style={{ marginBottom: '24px' }}
      >
        <p
          className="eyebrow"
          lang="en"
        >
          Dashboard
        </p>
        <h1 className="title">ダッシュボード</h1>
      </header>

      <div
        style={{
          display: 'grid',
          gap: '16px',
          gridTemplateColumns: 'repeat(auto-fill, minmax(240px, 1fr))',
        }}
      >
        <Link
          to="/articles"
          className="card"
          style={{
            textDecoration: 'none',
            color: 'inherit',
            padding: '24px',
            transition: 'box-shadow 0.2s',
            width: 'auto',
            margin: 0,
          }}
        >
          <h2 style={{ fontSize: '18px', marginTop: 0 }}>記事管理</h2>
          <p
            style={{
              color: 'var(--color-text-secondary)',
              fontSize: '14px',
              marginBottom: 0,
            }}
          >
            ブログ記事やZenn等の外部記事の管理
          </p>
        </Link>

        <Link
          to="/profile"
          className="card"
          style={{
            textDecoration: 'none',
            color: 'inherit',
            padding: '24px',
            transition: 'box-shadow 0.2s',
            width: 'auto',
            margin: 0,
          }}
        >
          <h2 style={{ fontSize: '18px', marginTop: 0 }}>プロフィール編集</h2>
          <p
            style={{
              color: 'var(--color-text-secondary)',
              fontSize: '14px',
              marginBottom: 0,
            }}
          >
            ポートフォリオに表示する経歴やスキルの管理
          </p>
        </Link>

        <Link
          to="/account"
          className="card"
          style={{
            textDecoration: 'none',
            color: 'inherit',
            padding: '24px',
            transition: 'box-shadow 0.2s',
            width: 'auto',
            margin: 0,
          }}
        >
          <h2 style={{ fontSize: '18px', marginTop: 0 }}>アカウント設定</h2>
          <p
            style={{
              color: 'var(--color-text-secondary)',
              fontSize: '14px',
              marginBottom: 0,
            }}
          >
            ログインメールアドレスとパスワードの変更
          </p>
        </Link>
      </div>
    </section>
  )
}
