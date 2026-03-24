import Keycloak from 'keycloak-js'

const url = import.meta.env.VITE_KEYCLOAK_URL
const realm = import.meta.env.VITE_KEYCLOAK_REALM
const clientId = import.meta.env.VITE_KEYCLOAK_CLIENT_ID

if (!url || !realm || !clientId) {
  throw new Error(
    'Missing required Keycloak configuration. ' +
      'Ensure VITE_KEYCLOAK_URL, VITE_KEYCLOAK_REALM, and VITE_KEYCLOAK_CLIENT_ID are set.',
  )
}

export const keycloak = new Keycloak({ url, realm, clientId })
