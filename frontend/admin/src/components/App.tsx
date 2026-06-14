import { useMemo } from 'react'
import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { AdminShell } from './AdminShell'
import { AuthBoundary } from './AuthBoundary'
import { DashboardPage } from '../pages/DashboardPage'
import { ArticlesPage } from '../pages/ArticlesPage'
import { AuthRepository } from '../domain/AuthRepository'
import { ArticleRepository } from '../domain/ArticleRepository'
import { ProfileRepository } from '../domain/ProfileRepository'
import { ProfilePage } from '../pages/ProfilePage'
import { AccountPage } from '../pages/AccountPage'
import { useSignal } from '../hooks/useSignal'

export function App() {
  // 1回だけインスタンス化する（Singleton として扱うならグローバルでも可）
  const authRepo = useMemo(() => new AuthRepository(), [])
  const articleRepo = useMemo(() => new ArticleRepository(), [])
  const profileRepo = useMemo(() => new ProfileRepository(), [])

  const authStatus = useSignal(authRepo.status)
  const isAuthenticated = authStatus === 'authenticated'

  return (
    <BrowserRouter>
      <AdminShell authenticated={isAuthenticated}>
        <AuthBoundary authRepo={authRepo}>
          <Routes>
            <Route
              path="/"
              element={<DashboardPage />}
            />
            <Route
              path="/articles"
              element={<ArticlesPage articleRepo={articleRepo} />}
            />
            <Route
              path="/profile"
              element={<ProfilePage profileRepo={profileRepo} />}
            />
            <Route
              path="/account"
              element={<AccountPage authRepo={authRepo} />}
            />
          </Routes>
        </AuthBoundary>
      </AdminShell>
    </BrowserRouter>
  )
}
