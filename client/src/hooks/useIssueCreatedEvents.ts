import { useEffect, useRef } from 'react'
import { useNats } from '@/context/NatsContext'
import { useWorkspace } from '@/context/WorkspaceContext'

interface Issue {
  ID: string
  Identifier: string
  Title: string
  Description: string | null
  Status: string
  Priority: string
  Labels: string[]
  AssigneeID: string | null
  ProjectID: string
  ReporterID: string
  CreatedAt: string
  UpdatedAt: string
}

export interface IssueCreatedEvent {
  occurred_at: string
  Payload: Issue
}

export function useIssueCreatedEvents(onIssueCreated: (event: IssueCreatedEvent) => void) {
  const { connection } = useNats()
  const { activeWorkspace } = useWorkspace()
  const activeWorkspaceRef = useRef(activeWorkspace)
  activeWorkspaceRef.current = activeWorkspace
  const onIssueCreatedRef = useRef(onIssueCreated)
  onIssueCreatedRef.current = onIssueCreated

  useEffect(() => {
    const workspace = activeWorkspaceRef.current
    if (!connection || !workspace) return

    const sub = connection.subscribe(`workspaces.${workspace.id}.issues.created`)

    const listen = async () => {
      for await (const msg of sub) {
        onIssueCreatedRef.current(msg.json<IssueCreatedEvent>())
      }
    }
    void listen()

    return () => {
      sub.unsubscribe()
    }
  }, [connection])
}
