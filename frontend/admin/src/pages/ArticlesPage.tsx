import { useEffect } from 'react';
import { useSignal } from '../hooks/useSignal';
import type { ArticleRepository } from '../domain/ArticleRepository';
import type { ArticleItem } from '../api/article-types';

export function ArticlesPage({ articleRepo }: { articleRepo: ArticleRepository }) {
  const articles = useSignal(articleRepo.articles);
  const loading = useSignal(articleRepo.isLoading);
  const error = useSignal(articleRepo.error);

  useEffect(() => {
    // 初回マウント時に記事一覧を読み込む
    void articleRepo.loadInitialData();
  }, [articleRepo]);

  return (
    <section className="container">
      <header className="page-header" style={{ marginBottom: '24px', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <div>
          <p className="eyebrow" lang="en">Articles</p>
          <h1 className="title">記事管理</h1>
          <p className="description" style={{ marginBottom: 0 }}>
            外部プラットフォームや自前ブログの記事を管理します。
          </p>
        </div>
        <button type="button" onClick={() => alert('TODO: 新規作成ダイアログを開く')}>
          新規作成
        </button>
      </header>

      {error && <p className="message error">{error}</p>}

      {loading ? (
        <p>読み込み中...</p>
      ) : (
        <div className="card" style={{ width: '100%', padding: '0' }}>
          <table style={{ width: '100%', borderCollapse: 'collapse', textAlign: 'left' }}>
            <thead>
              <tr style={{ borderBottom: '1px solid var(--color-border)', backgroundColor: 'var(--color-bg-surface)' }}>
                <th style={{ padding: '12px 16px', fontWeight: 500 }}>タイトル</th>
                <th style={{ padding: '12px 16px', fontWeight: 500 }}>プラットフォーム</th>
                <th style={{ padding: '12px 16px', fontWeight: 500 }}>公開日</th>
                <th style={{ padding: '12px 16px', width: '80px' }}></th>
              </tr>
            </thead>
            <tbody>
              {articles.length === 0 ? (
                <tr>
                  <td colSpan={4} style={{ padding: '24px', textAlign: 'center', color: 'var(--color-text-secondary)' }}>
                    記事がありません。
                  </td>
                </tr>
              ) : (
                articles.map((article: ArticleItem) => (
                  <tr key={article.externalId} style={{ borderBottom: '1px solid var(--color-border)' }}>
                    <td style={{ padding: '12px 16px' }}>
                      <div style={{ fontWeight: 500 }}>{article.title}</div>
                      <a href={article.url} target="_blank" rel="noreferrer" style={{ fontSize: '12px', color: 'var(--admin-accent)' }}>
                        {article.url}
                      </a>
                    </td>
                    <td style={{ padding: '12px 16px' }}>{article.platform}</td>
                    <td style={{ padding: '12px 16px', fontSize: '14px' }}>
                      {article.publishedAt ? new Date(article.publishedAt).toLocaleDateString('ja-JP') : '-'}
                    </td>
                    <td style={{ padding: '12px 16px' }}>
                      <button type="button" className="subtle" style={{ padding: '6px 12px' }} onClick={() => alert('TODO: 編集機能')}>
                        編集
                      </button>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      )}
    </section>
  );
}
