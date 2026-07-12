import { useEffect, useRef, useState } from 'react'
import { upsertCurrentUser } from '@/api/generated/issueTrackerAPI'
import { keycloak } from '../keycloak'
import { KeycloakContext } from './KeycloakContext'

interface KeycloakProviderProps {
  children: React.ReactNode
}

export function KeycloakProvider({ children }: KeycloakProviderProps) {
  const [initialized, setInitialized] = useState(false)
  const [authenticated, setAuthenticated] = useState(false)
  const [error, setError] = useState<Error | null>(null)
  const initCalled = useRef(false)

  useEffect(() => {
    if (initCalled.current) return
    initCalled.current = true

    keycloak.onTokenExpired = () => {
      keycloak.updateToken(30).catch(() => {
        keycloak.login()
      })
    }

    keycloak
      .init({ onLoad: 'login-required' })
      .then((result) => {
        setAuthenticated(result)
        setInitialized(true)
        if (result) {
          upsertCurrentUser().catch((err: unknown) => {
            console.error('Failed to sync current user', err)
          })
        }
      })
      .catch((err: unknown) => {
        setError(err instanceof Error ? err : new Error('Keycloak initialization failed'))
        // No setInitialized — there is no recovery from a failed init
      })
  }, [])

  if (error !== null) {
    throw error
  }

  if (!initialized || !authenticated) {
    return <div>Loading...</div>
  }

  return (
    <KeycloakContext.Provider value={{ keycloak, authenticated }}>
      {children}
    </KeycloakContext.Provider>
  )
}
