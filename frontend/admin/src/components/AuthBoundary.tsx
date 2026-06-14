import React, { useEffect } from 'react'
import { useSignal } from '../hooks/useSignal'
import type { IAuthRepository } from '../domain/AuthRepository'
import { LoginPage } from '../pages/LoginPage'

interface AuthBoundaryProps {
  authRepo: IAuthRepository
  children: React.ReactNode
}

export function AuthBoundary({ authRepo, children }: AuthBoundaryProps) {
  const status = useSignal(authRepo.status)

  useEffect(() => {
    if (status === 'unknown') {
      void authRepo.refreshSession()
    }
  }, [authRepo, status])

  if (status === 'unknown' || status === 'checking') {
    return (
      <div
        className="status"
        style={{ minHeight: '60dvh', display: 'grid', placeItems: 'center' }}
      >
        <p>認証状態を確認しています...</p>
      </div>
    )
  }

  if (status === 'authenticated') {
    return <>{children}</>
  }

  return <LoginPage authRepo={authRepo} />
}
