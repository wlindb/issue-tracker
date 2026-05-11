import { connect } from 'nats.ws'
import { useEffect, useRef, useState } from 'react'
import type { NatsConnection } from 'nats.ws'
import { useKeycloak } from '@/auth/useKeycloak'
import { useWorkspace } from '@/context/WorkspaceContext'
import { NatsContext } from '@/context/NatsContext'

const natsWsUrl = import.meta.env.VITE_NATS_WS_URL ?? "wss://issuetrackerapi.astraterial.com:4233"

if (!natsWsUrl) {
  console.warn('[NATS] VITE_NATS_WS_URL is not set — WebSocket connection disabled')
}

interface NatsProviderProps {
  children: React.ReactNode
}

export function NatsProvider({ children }: NatsProviderProps) {
  const { keycloak } = useKeycloak()
  const { activeWorkspace } = useWorkspace()
  const [connection, setConnection] = useState<NatsConnection | null>(null)
  const [error, setError] = useState<Error | null>(null)
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
        setError(null)
        console.log('[NATS] connecting', { url: natsWsUrl, workspace: activeWorkspace.id })

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
        console.log('[NATS] connected', nc.info)

        const closeErr = await nc.closed()
        if (connectionRef.current === nc) {
          connectionRef.current = null
          setConnection(null)
        }
        if (closeErr) {
          console.error('[NATS] connection closed with error', closeErr)
          if (!cancelled) {
            setError(closeErr instanceof Error ? closeErr : new Error(String(closeErr)))
          }
        } else {
          console.log('[NATS] connection closed cleanly')
        }
      } catch (err: unknown) {
        if (!cancelled) {
          const message = err instanceof Error ? err.message : String(err)
          console.error('[NATS] connection failed', { url: natsWsUrl, workspace: activeWorkspace.id, error: message })
          setError(err instanceof Error ? err : new Error(message))
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

  return <NatsContext.Provider value={{ connection, error }}>{children}</NatsContext.Provider>
}
