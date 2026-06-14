import React, { useState } from 'react';
import { useSignal } from '../hooks/useSignal';
import type { IAuthRepository } from '../domain/AuthRepository';

interface LoginPageProps {
  authRepo: IAuthRepository;
}

export function LoginPage({ authRepo }: LoginPageProps) {
  const [passwordVisible, setPasswordVisible] = useState(false);
  const error = useSignal(authRepo.error);
  const notice = useSignal(authRepo.notice);
  const isPending = useSignal(authRepo.loginPending);

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const form = e.currentTarget;
    const formData = new FormData(form);

    await authRepo.login({
      emailAddress: (formData.get('emailAddress') as string).trim(),
      password: formData.get('password') as string,
    });
  };

  return (
    <section className="container" style={{ minHeight: '100dvh', display: 'grid', placeItems: 'center', padding: '24px 16px' }}>
      <div className="card">
        <p className="eyebrow" lang="en">Admin</p>
        <h1 className="title">ログイン</h1>
        <p className="description">
          管理画面へ入るには、メールアドレスとパスワードでログインしてください。
        </p>

        {notice && <p className="message notice">{notice}</p>}

        <form onSubmit={handleSubmit}>
          <div className="form-field">
            <label htmlFor="emailAddress">メールアドレス</label>
            <input
              id="emailAddress"
              type="email"
              name="emailAddress"
              autoComplete="email"
              disabled={isPending}
              required
              autoFocus
            />
          </div>

          <div className="form-field password-field-container">
            <label htmlFor="password">パスワード</label>
            <div style={{ display: 'flex', gap: '8px' }}>
              <input
                id="password"
                type={passwordVisible ? 'text' : 'password'}
                name="password"
                autoComplete="current-password"
                disabled={isPending}
                required
                style={{ flex: 1 }}
              />
              <button
                type="button"
                className="subtle"
                onClick={() => setPasswordVisible(!passwordVisible)}
                disabled={isPending}
              >
                {passwordVisible ? '隠す' : '表示'}
              </button>
            </div>
          </div>

          {error && <p className="message error">{error}</p>}

          <button type="submit" disabled={isPending} style={{ marginTop: '8px' }}>
            {isPending ? 'ログイン中...' : 'ログイン'}
          </button>
        </form>
      </div>
    </section>
  );
}
