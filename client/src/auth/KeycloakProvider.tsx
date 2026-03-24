import Keycloak from 'keycloak-js'
import { createContext, useContext, useEffect, useState } from 'react'
import { keycloak } from '../keycloak'

interface KeycloakContextValue {
  keycloak: Keycloak
  authenticated: boolean
}

const KeycloakContext = createContext<KeycloakContextValue | null>(null)

export function useKeycloak(): KeycloakContextValue {
  const context = useContext(KeycloakContext)
  if (context === null) {
    throw new Error('useKeycloak must be used within a KeycloakProvider')
  }
  return context
}

interface KeycloakProviderProps {
  children: React.ReactNode
}

export function KeycloakProvider({ children }: KeycloakProviderProps) {
  const [initialized, setInitialized] = useState(false)
  const [authenticated, setAuthenticated] = useState(false)

  useEffect(() => {
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
      })
      .catch(() => {
        setInitialized(true)
      })
  }, [])

  if (!initialized) {
    return <div>Loading...</div>
  }

  return (
    <KeycloakContext.Provider value={{ keycloak, authenticated }}>
      {children}
    </KeycloakContext.Provider>
  )
}
