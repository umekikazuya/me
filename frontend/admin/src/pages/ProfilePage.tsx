import { useEffect, useSyncExternalStore } from 'react'
import type { MeProfile } from '../api/types'
import type { ProfileRepository } from '../domain/ProfileRepository'

// モジュールレベルに切り出すことで ProfilePage の複雑度を下げる
function parseJSON<T>(str: string, fallback: T): T {
  try {
    return str ? (JSON.parse(str) as T) : fallback
  } catch {
    return fallback
  }
}

function useProfileStore(profileRepo: ProfileRepository) {
  const subscribe = (cb: () => void) => {
    profileRepo.addEventListener('profile:admin-change', cb)
    profileRepo.addEventListener('change', cb)
    return () => {
      profileRepo.removeEventListener('profile:admin-change', cb)
      profileRepo.removeEventListener('change', cb)
    }
  }
  return {
    profile: useSyncExternalStore(subscribe, () => profileRepo.adminProfile),
    error: useSyncExternalStore(subscribe, () => profileRepo.adminError),
    success: useSyncExternalStore(subscribe, () => profileRepo.adminSuccess),
    loading: useSyncExternalStore(subscribe, () => profileRepo.adminLoading),
    saving: useSyncExternalStore(subscribe, () => profileRepo.adminSaving),
    dirty: useSyncExternalStore(subscribe, () => profileRepo.adminDirty),
    loaded: useSyncExternalStore(subscribe, () => profileRepo.adminLoaded),
  }
}

function BasicInfoFields({ profile }: { profile: MeProfile }) {
  return (
    <div
      className="card"
      style={{ width: '100%', margin: 0 }}
    >
      <h2 style={{ fontSize: '18px', marginTop: 0 }}>基本情報</h2>
      <p
        style={{
          fontSize: '13px',
          color: 'var(--color-text-secondary)',
          marginBottom: '16px',
        }}
      >
        最低限、表示名だけあれば更新できます。未入力項目は公開画面で省略されます。
      </p>
      <div
        style={{
          display: 'grid',
          gap: '16px',
          gridTemplateColumns: 'repeat(auto-fit, minmax(220px, 1fr))',
        }}
      >
        <div className="form-field">
          <label htmlFor="profile-displayName">表示名 *</label>
          <input
            id="profile-displayName"
            type="text"
            name="displayName"
            defaultValue={profile.displayName}
            required
          />
        </div>
        <div className="form-field">
          <label htmlFor="profile-displayJa">表示名（日本語）</label>
          <input
            id="profile-displayJa"
            type="text"
            name="displayJa"
            defaultValue={profile.displayJa}
          />
        </div>
        <div className="form-field">
          <label htmlFor="profile-role">Role</label>
          <input
            id="profile-role"
            type="text"
            name="role"
            defaultValue={profile.role}
          />
        </div>
        <div className="form-field">
          <label htmlFor="profile-location">Location</label>
          <input
            id="profile-location"
            type="text"
            name="location"
            defaultValue={profile.location}
          />
        </div>
      </div>
    </div>
  )
}

function DetailInfoFields({ profile }: { profile: MeProfile }) {
  return (
    <div
      className="card"
      style={{ width: '100%', margin: 0 }}
    >
      <h2 style={{ fontSize: '18px', marginTop: 0 }}>詳細情報 (JSON)</h2>
      <p
        style={{
          fontSize: '13px',
          color: 'var(--color-text-secondary)',
          marginBottom: '16px',
        }}
      >
        各配列をJSON形式で入力してください。
      </p>
      <div style={{ display: 'grid', gap: '16px' }}>
        <div className="form-field">
          <label htmlFor="profile-skills">Skills</label>
          <textarea
            id="profile-skills"
            name="skills"
            rows={4}
            defaultValue={JSON.stringify(profile.skills ?? [], null, 2)}
          />
        </div>
        <div className="form-field">
          <label htmlFor="profile-certifications">Certifications</label>
          <textarea
            id="profile-certifications"
            name="certifications"
            rows={4}
            defaultValue={JSON.stringify(profile.certifications ?? [], null, 2)}
          />
        </div>
        <div className="form-field">
          <label htmlFor="profile-experiences">Experiences</label>
          <textarea
            id="profile-experiences"
            name="experiences"
            rows={6}
            defaultValue={JSON.stringify(profile.experiences ?? [], null, 2)}
          />
        </div>
        <div className="form-field">
          <label htmlFor="profile-links">Links</label>
          <textarea
            id="profile-links"
            name="links"
            rows={4}
            defaultValue={JSON.stringify(profile.links ?? [], null, 2)}
          />
        </div>
        <div className="form-field">
          <label htmlFor="profile-likes">1行につき1件</label>
          <textarea
            id="profile-likes"
            name="likes"
            rows={6}
            defaultValue={(profile.likes ?? []).join('\n')}
          />
        </div>
      </div>
    </div>
  )
}

export function ProfilePage({
  profileRepo,
}: {
  profileRepo: ProfileRepository
}) {
  const { profile, error, success, loading, saving, dirty, loaded } =
    useProfileStore(profileRepo)

  useEffect(() => {
    if (!loaded) {
      void profileRepo.loadAdminProfile()
    }
  }, [profileRepo, loaded])

  useEffect(() => {
    const handleBeforeUnload = (e: BeforeUnloadEvent) => {
      if (dirty) {
        e.preventDefault()
        e.returnValue = ''
      }
    }
    window.addEventListener('beforeunload', handleBeforeUnload)
    return () => window.removeEventListener('beforeunload', handleBeforeUnload)
  }, [dirty])

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    const form = e.currentTarget
    if (!form.checkValidity()) {
      form.reportValidity()
      return
    }
    const formData = new FormData(form)
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
    }
    await profileRepo.saveAdminProfile(newProfile)
  }

  const handleReset = () => {
    if (dirty && !window.confirm('未保存の変更を破棄して元に戻しますか？')) {
      return
    }
    profileRepo.setAdminDirty(false)
    window.location.reload()
  }

  return (
    <section className="container">
      <header
        className="page-header"
        style={{
          marginBottom: '24px',
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'end',
          flexWrap: 'wrap',
          gap: '24px',
        }}
      >
        <div>
          <p
            className="eyebrow"
            lang="en"
          >
            Profile
          </p>
          <h1 className="title">プロフィール編集</h1>
          <p
            className="description"
            style={{ marginBottom: 0 }}
          >
            公開プロフィールの表示内容を更新します。
          </p>
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
        <p style={{ color: 'var(--color-text-secondary)' }}>
          プロフィールを読み込み中...
        </p>
      ) : (
        <form
          onSubmit={handleSubmit}
          onInput={() => profileRepo.setAdminDirty(true)}
          style={{ display: 'grid', gap: '24px' }}
        >
          <BasicInfoFields profile={profile} />
          <DetailInfoFields profile={profile} />

          <div
            style={{
              display: 'flex',
              alignItems: 'center',
              gap: '12px',
              justifyContent: 'flex-end',
              position: 'sticky',
              bottom: '16px',
              padding: '14px 16px',
              border: '1px solid var(--color-border)',
              background: 'rgba(255, 255, 255, 0.95)',
              backdropFilter: 'blur(12px)',
              zIndex: 10,
              borderRadius: '4px',
            }}
          >
            <div style={{ marginRight: 'auto' }}>
              <p
                style={{
                  fontSize: '13px',
                  margin: 0,
                  color: dirty ? '#9a6d2f' : 'var(--color-text-tertiary)',
                }}
              >
                {dirty ? '未保存の変更があります。' : '保存済みの内容です。'}
              </p>
            </div>
            <button
              type="button"
              className="subtle"
              disabled={saving}
              onClick={handleReset}
            >
              変更を元に戻す
            </button>
            <button
              type="submit"
              disabled={saving || !dirty}
            >
              {saving ? '保存中...' : '保存する'}
            </button>
          </div>
        </form>
      )}
    </section>
  )
}
