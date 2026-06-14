import React, { useState } from 'react';
import { useSignal } from '../hooks/useSignal';
import type { AuthRepository } from '../domain/AuthRepository';

export function AccountPage({ authRepo }: { authRepo: AuthRepository }) {
  const error = useSignal(authRepo.error);
  const notice = useSignal(authRepo.notice);
  
  const [emailPending, setEmailPending] = useState(false);

  const handleEmailChange = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const form = e.currentTarget;
    const formData = new FormData(form);
    
    setEmailPending(true);
    await authRepo.changeEmail({
      token: '', // TODO: implement token flow if needed
      newEmailAddress: formData.get('newEmailAddress') as string,
    });
    setEmailPending(false);
    form.reset();
  };

  const handleLogout = async () => {
    if (window.confirm('ログアウトしますか？')) {
      await authRepo.logout();
    }
  };

  return (
    <section className="container">
      <header className="page-header" style={{ marginBottom: '24px', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <div>
          <p className="eyebrow" lang="en">Account</p>
          <h1 className="title">アカウント設定</h1>
          <p className="description" style={{ marginBottom: 0 }}>
            ログインメールアドレスを変更できます。
          </p>
        </div>
        <button type="button" className="subtle" onClick={handleLogout}>
          ログアウト
        </button>
      </header>

      {error && <p className="message error">{error}</p>}
      {notice && <p className="message notice">{notice}</p>}

      <div style={{ display: 'grid', gap: '24px' }}>
        <div className="card" style={{ width: '100%', margin: 0 }}>
          <h2 style={{ fontSize: '18px', marginTop: 0 }}>メールアドレス変更</h2>
          <form onSubmit={handleEmailChange}>
            <div className="form-field">
              <label>新しいメールアドレス</label>
              <input type="email" name="newEmailAddress" required disabled={emailPending} />
            </div>
            <button type="submit" disabled={emailPending} style={{ justifySelf: 'start', marginTop: '16px' }}>
              {emailPending ? '変更中...' : 'メールアドレスを変更'}
            </button>
          </form>
        </div>
      </div>
    </section>
  );
}
