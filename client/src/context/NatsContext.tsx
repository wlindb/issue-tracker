import type { NatsConnection } from 'nats.ws'
import { createContext, useContext } from 'react'

export interface NatsContextValue {
  connection: NatsConnection | null
  error: Error | null
}

export const NatsContext = createContext<NatsContextValue | null>(null)

export function useNats(): NatsContextValue {
  const context = useContext(NatsContext)
  if (context === null) {
    throw new Error('useNats must be used within a NatsProvider')
  }
  return context
}
