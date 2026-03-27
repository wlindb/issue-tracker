import Keycloak from 'keycloak-js'

// TODO: Fix this with proper env loading and no hardcoded fallback values
const url = import.meta.env.VITE_KEYCLOAK_URL ?? "https://keycloak.astraterial.com"
const realm = import.meta.env.VITE_KEYCLOAK_REALM ?? "issue-tracker"
const clientId = import.meta.env.VITE_KEYCLOAK_CLIENT_ID ?? "issue-tracker-frontend"

if (!url || !realm || !clientId) {
  throw new Error(
    'Missing required Keycloak configuration. ' +
    'Ensure VITE_KEYCLOAK_URL, VITE_KEYCLOAK_REALM, and VITE_KEYCLOAK_CLIENT_ID are set.',
  )
}

export const keycloak = new Keycloak({ url, realm, clientId })
