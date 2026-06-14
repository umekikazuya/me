import { useEffect, useSyncExternalStore } from 'react';
import type { ProfileRepository } from '../domain/ProfileRepository';
import type { MeProfile } from '../api/types';

export function ProfilePage({ profileRepo }: { profileRepo: ProfileRepository }) {
  const subscribe = (cb: () => void) => {
    profileRepo.addEventListener('profile:admin-change', cb);
    profileRepo.addEventListener('change', cb);
    return () => {
      profileRepo.removeEventListener('profile:admin-change', cb);
      profileRepo.removeEventListener('change', cb);
    };
  };

  const profile = useSyncExternalStore(subscribe, () => profileRepo.adminProfile);
  const error = useSyncExternalStore(subscribe, () => profileRepo.adminError);
  const success = useSyncExternalStore(subscribe, () => profileRepo.adminSuccess);
  const loading = useSyncExternalStore(subscribe, () => profileRepo.adminLoading);
  const saving = useSyncExternalStore(subscribe, () => profileRepo.adminSaving);
  const dirty = useSyncExternalStore(subscribe, () => profileRepo.adminDirty);
  const loaded = useSyncExternalStore(subscribe, () => profileRepo.adminLoaded);

  useEffect(() => {
    if (!loaded) {
      void profileRepo.loadAdminProfile();
    }
  }, [profileRepo, loaded]);

  useEffect(() => {
    const handleBeforeUnload = (e: BeforeUnloadEvent) => {
      if (dirty) {
        e.preventDefault();
        e.returnValue = '';
      }
    };
    window.addEventListener('beforeunload', handleBeforeUnload);
    return () => window.removeEventListener('beforeunload', handleBeforeUnload);
  }, [dirty]);

  const handleInput = () => {
    profileRepo.setAdminDirty(true);
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const form = e.currentTarget;
    if (!form.checkValidity()) {
      form.reportValidity();
      return;
    }

    const formData = new FormData(form);
    
    // JSON parse safe helper
    const parseJSON = (str: string, fallback: any) => {
      try { return str ? JSON.parse(str) : fallback; }
      catch { return fallback; }
    };

    const newProfile: MeProfile = {
      displayName: formData.get('displayName') as string,
      displayJa: formData.get('displayJa') as string,
      role: formData.get('role') as string,
      location: formData.get('location') as string,
      skills: parseJSON(formData.get('skills') as string, []),
      certifications: parseJSON(formData.get('certifications') as string, []),
      experiences: parseJSON(formData.get('experiences') as string, []),
      links: parseJSON(formData.get('links') as string, []),
      likes: ((formData.get('likes') as string) || '')
        .split('\n')
        .map((s) => s.trim())
        .filter(Boolean),
      updatedAt: profile.updatedAt,
    };

    await profileRepo.saveAdminProfile(newProfile);
  };

  const handleReset = () => {
    if (dirty && !window.confirm('未保存の変更を破棄して元に戻しますか？')) {
      return;
    }
    profileRepo.setAdminDirty(false);
    // フォームのリセットはキーを変更するなどして対応する（簡略化のためリロード）
    window.location.reload();
  };

  return (
    <section className="container">
      <header className="page-header" style={{ marginBottom: '24px', display: 'flex', justifyContent: 'space-between', alignItems: 'end', flexWrap: 'wrap', gap: '24px' }}>
        <div>
          <p className="eyebrow" lang="en">Profile</p>
          <h1 className="title">プロフィール編集</h1>
          <p className="description" style={{ marginBottom: 0 }}>公開プロフィールの表示内容を更新します。</p>
        </div>
        {profile.updatedAt && (
          <p style={{ color: 'var(--color-text-secondary)', fontSize: '14px' }}>
            最終更新: {new Date(profile.updatedAt).toLocaleString('ja-JP')}
          </p>
        )}
      </header>

      {error && <p className="message error">{error}</p>}
      {success && <p className="message notice">{success}</p>}

      {loading ? (
        <p style={{ color: 'var(--color-text-secondary)' }}>プロフィールを読み込み中...</p>
      ) : (
        <form onSubmit={handleSubmit} onInput={handleInput} style={{ display: 'grid', gap: '24px' }}>
          
          <div className="card" style={{ width: '100%', margin: 0 }}>
            <h2 style={{ fontSize: '18px', marginTop: 0 }}>基本情報</h2>
            <p style={{ fontSize: '13px', color: 'var(--color-text-secondary)', marginBottom: '16px' }}>最低限、表示名だけあれば更新できます。未入力項目は公開画面で省略されます。</p>
            
            <div style={{ display: 'grid', gap: '16px', gridTemplateColumns: 'repeat(auto-fit, minmax(220px, 1fr))' }}>
              <div className="form-field">
                <label>表示名 *</label>
                <input type="text" name="displayName" defaultValue={profile.displayName} required />
              </div>
              <div className="form-field">
                <label>表示名（日本語）</label>
                <input type="text" name="displayJa" defaultValue={profile.displayJa} />
              </div>
              <div className="form-field">
                <label>Role</label>
                <input type="text" name="role" defaultValue={profile.role} />
              </div>
              <div className="form-field">
                <label>Location</label>
                <input type="text" name="location" defaultValue={profile.location} />
              </div>
            </div>
          </div>

          <div className="card" style={{ width: '100%', margin: 0 }}>
            <h2 style={{ fontSize: '18px', marginTop: 0 }}>詳細情報 (JSON)</h2>
            <p style={{ fontSize: '13px', color: 'var(--color-text-secondary)', marginBottom: '16px' }}>各配列をJSON形式で入力してください。</p>
            
            <div style={{ display: 'grid', gap: '16px' }}>
              <div className="form-field">
                <label>Skills</label>
                <textarea name="skills" rows={4} defaultValue={JSON.stringify(profile.skills || [], null, 2)} />
              </div>
              <div className="form-field">
                <label>Certifications</label>
                <textarea name="certifications" rows={4} defaultValue={JSON.stringify(profile.certifications || [], null, 2)} />
              </div>
              <div className="form-field">
                <label>Experiences</label>
                <textarea name="experiences" rows={6} defaultValue={JSON.stringify(profile.experiences || [], null, 2)} />
              </div>
              <div className="form-field">
                <label>Links</label>
                <textarea name="links" rows={4} defaultValue={JSON.stringify(profile.links || [], null, 2)} />
              </div>
            </div>
          </div>

          <div className="card" style={{ width: '100%', margin: 0 }}>
            <h2 style={{ fontSize: '18px', marginTop: 0 }}>Likes</h2>
            <p style={{ fontSize: '13px', color: 'var(--color-text-secondary)', marginBottom: '16px' }}>1行ごとに1件ずつ入力します。空行は保存時に除外されます。</p>
            <div className="form-field">
              <label>1行につき1件</label>
              <textarea name="likes" rows={6} defaultValue={(profile.likes || []).join('\n')} />
            </div>
          </div>

          <div style={{ 
            display: 'flex', alignItems: 'center', gap: '12px', justifyContent: 'flex-end',
            position: 'sticky', bottom: '16px', padding: '14px 16px',
            border: '1px solid var(--color-border)', background: 'rgba(255, 255, 255, 0.95)',
            backdropFilter: 'blur(12px)', zIndex: 10, borderRadius: '4px'
          }}>
            <div style={{ marginRight: 'auto' }}>
              <p style={{ fontSize: '13px', margin: 0, color: dirty ? '#9a6d2f' : 'var(--color-text-tertiary)' }}>
                {dirty ? '未保存の変更があります。' : '保存済みの内容です。'}
              </p>
            </div>
            <button type="button" className="subtle" disabled={saving} onClick={handleReset}>
              変更を元に戻す
            </button>
            <button type="submit" disabled={saving || !dirty}>
              {saving ? '保存中...' : '保存する'}
            </button>
          </div>
        </form>
      )}
    </section>
  );
}
