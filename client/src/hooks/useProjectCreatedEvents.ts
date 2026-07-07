import { useEffect, useRef } from 'react'
import { type Project } from '@/api/generated/issueTrackerAPI'
import { useNats } from '@/context/NatsContext'
import { useWorkspace } from '@/context/WorkspaceContext'

interface RawProject {
  ID: string
  Identifier: string
  Name: string
  Description?: string | null
  OwnerID: string
  CreatedAt: string
  UpdatedAt: string
}

export interface ProjectCreatedEvent {
  occurred_at: string
  payload: Project
}

export interface ApiProjectCreatedEvent {
  occurred_at: string
  payload: RawProject
}

function toProject(raw: RawProject): Project {
  return {
    id: raw.ID,
    identifier: raw.Identifier,
    name: raw.Name,
    description: raw.Description,
    ownerId: raw.OwnerID,
    issueCount: 0,
    createdAt: raw.CreatedAt,
    updatedAt: raw.UpdatedAt,
  }
}

export function useProjectCreatedEvents(onProjectCreated: (event: ProjectCreatedEvent) => void) {
  const { connection } = useNats()
  const { activeWorkspace } = useWorkspace()
  const activeWorkspaceRef = useRef(activeWorkspace)
  activeWorkspaceRef.current = activeWorkspace
  const onProjectCreatedRef = useRef(onProjectCreated)
  onProjectCreatedRef.current = onProjectCreated

  useEffect(() => {
    const workspace = activeWorkspaceRef.current
    if (!connection || !workspace) return

    const sub = connection.subscribe(`workspaces.${workspace.id}.projects.created`)

    const listen = async () => {
      for await (const msg of sub) {
        try {
          const raw = msg.json<ApiProjectCreatedEvent>()
          onProjectCreatedRef.current({ occurred_at: raw.occurred_at, payload: toProject(raw.payload) })
        } catch (err) {
          console.error('[NATS] failed to process projects.created message', err)
        }
      }
    }
    void listen()

    return () => {
      sub.unsubscribe()
    }
  }, [connection])
}
