import { Route, Routes } from 'react-router-dom'
import { useKeycloak } from './auth/KeycloakProvider'

function App() {
  const { keycloak } = useKeycloak()
  const username = keycloak.tokenParsed?.preferred_username as string | undefined

  return (
    <Routes>
      <Route
        path="/"
        element={
          <div>
            <p>Logged in as <strong>{username}</strong></p>
            <button onClick={() => keycloak.logout()}>Logout</button>
          </div>
        }
      />
    </Routes>
  )
}

export default App
