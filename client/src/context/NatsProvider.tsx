import { connect } from 'nats.ws'
import { useEffect, useRef, useState } from 'react'
import type { NatsConnection } from 'nats.ws'
import { useKeycloak } from '@/auth/useKeycloak'
import { useWorkspace } from '@/context/WorkspaceContext'
import { NatsContext } from '@/context/NatsContext'

const natsWsUrl = import.meta.env.VITE_NATS_WS_URL as string | undefined

interface NatsProviderProps {
  children: React.ReactNode
}

export function NatsProvider({ children }: NatsProviderProps) {
  const { keycloak } = useKeycloak()
  const { activeWorkspace } = useWorkspace()
  const [connection, setConnection] = useState<NatsConnection | null>(null)
  const connectionRef = useRef<NatsConnection | null>(null)

  useEffect(() => {
    if (!natsWsUrl || !activeWorkspace) {
      return
    }

    const token = keycloak.token
    if (!token) {
      return
    }

    let cancelled = false

    const runConnection = async () => {
      try {
        const nc = await connect({
          servers: natsWsUrl,
          user: activeWorkspace.id,
          pass: token,
        })

        if (cancelled) {
          void nc.close()
          return
        }

        connectionRef.current = nc
        setConnection(nc)
        console.log('connected to ', nc.info)

        // When the server closes the connection, clear state
        await nc.closed()
        if (connectionRef.current === nc) {
          connectionRef.current = null
          setConnection(null)
        }
      } catch (err: unknown) {
        if (!cancelled) {
          console.error('NATS WebSocket connection failed', err)
        }
      }
    }

    void runConnection()

    return () => {
      cancelled = true
      const prev = connectionRef.current
      if (prev) {
        connectionRef.current = null
        setConnection(null)
        prev.close()
      }
    }
  }, [activeWorkspace, keycloak.token])

  return <NatsContext.Provider value={{ connection }}>{children}</NatsContext.Provider>
}
