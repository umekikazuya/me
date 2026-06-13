import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { DashboardPage } from '../pages/DashboardPage.tsx'
import { LoginPage } from '../pages/LoginPage.tsx'
import { ArticlesPage } from '../pages/ArticlesPage.tsx'

export function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route
          path="/"
          element={<DashboardPage />}
        />
        <Route
          path="/login"
          element={<LoginPage />}
        />
        <Route
          path="/articles"
          element={<ArticlesPage />}
        />
      </Routes>
    </BrowserRouter>
  )
}
