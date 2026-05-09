import Keycloak from 'keycloak-js'
import { createContext } from 'react'

export interface KeycloakContextValue {
  keycloak: Keycloak
  authenticated: boolean
}

export const KeycloakContext = createContext<KeycloakContextValue | null>(null)
