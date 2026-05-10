import { useContext } from 'react'
import { KeycloakContext, type KeycloakContextValue } from './KeycloakContext'

export function useKeycloak(): KeycloakContextValue {
  const context = useContext(KeycloakContext)
  if (context === null) {
    throw new Error('useKeycloak must be used within a KeycloakProvider')
  }
  return context
}
