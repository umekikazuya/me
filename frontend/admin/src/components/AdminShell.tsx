import type React from 'react'
import { NavLink } from 'react-router-dom'

interface AdminShellProps {
  authenticated: boolean
  children: React.ReactNode
}

export function AdminShell({ authenticated, children }: AdminShellProps) {
  return (
    <div className={`layout ${authenticated ? 'with-sidebar' : ''}`}>
      {authenticated && (
        <aside className="sidebar">
          <NavLink
            to="/"
            end
            className={({ isActive }) => (isActive ? 'active' : '')}
          >
            Dashboard
          </NavLink>
          <NavLink
            to="/articles"
            className={({ isActive }) => (isActive ? 'active' : '')}
          >
            Articles
          </NavLink>
          <NavLink
            to="/profile"
            className={({ isActive }) => (isActive ? 'active' : '')}
          >
            Profile
          </NavLink>
          <NavLink
            to="/account"
            className={({ isActive }) => (isActive ? 'active' : '')}
          >
            Account
          </NavLink>
        </aside>
      )}
      <main id="outlet">{children}</main>
    </div>
  )
}
