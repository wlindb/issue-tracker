import { useEffect, useRef } from 'react'
import { type Issue } from '@/api/generated/issueTrackerAPI'
import { useNats } from '@/context/NatsContext'
import { useWorkspace } from '@/context/WorkspaceContext'

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
        try {
          onIssueCreatedRef.current(msg.json<IssueCreatedEvent>())
        } catch (err) {
          console.error('[NATS] failed to process issue.created message', err)
        }
      }
    }
    void listen()

    return () => {
      sub.unsubscribe()
    }
  }, [connection])
}
