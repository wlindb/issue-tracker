import { useEffect, useRef } from 'react'
import { type Issue, type Label, type IssueStatus, type IssuePriority } from '@/api/generated/issueTrackerAPI'
import { useNats } from '@/context/NatsContext'
import { useWorkspace } from '@/context/WorkspaceContext'

interface EventIssue {
  id: string
  identifier: string
  projectId: string
  title: string
  description?: string | null
  status: IssueStatus
  priority: IssuePriority
  labels: Label[]
  assigneeId?: string | null
  reporterId: string
  createdAt: string
  updatedAt: string
}

export interface IssuePriorityUpdatedEvent {
  occurred_at: string
  payload: Issue
}

export interface ApiIssuePriorityUpdatedEvent {
  occurred_at: string
  payload: EventIssue
}

function toIssue(raw: EventIssue): Issue {
  return raw
}

export function useIssuePriorityUpdatedEvents(onIssuePriorityUpdated: (event: IssuePriorityUpdatedEvent) => void) {
  const { connection } = useNats()
  const { activeWorkspace } = useWorkspace()
  const activeWorkspaceRef = useRef(activeWorkspace)
  activeWorkspaceRef.current = activeWorkspace
  const onIssuePriorityUpdatedRef = useRef(onIssuePriorityUpdated)
  onIssuePriorityUpdatedRef.current = onIssuePriorityUpdated

  useEffect(() => {
    const workspace = activeWorkspaceRef.current
    if (!connection || !workspace) return

    const sub = connection.subscribe(`workspaces.${workspace.id}.issues.priority_updated`)

    const listen = async () => {
      for await (const msg of sub) {
        try {
          const raw = msg.json<ApiIssuePriorityUpdatedEvent>()
          onIssuePriorityUpdatedRef.current({ occurred_at: raw.occurred_at, payload: toIssue(raw.payload) })
        } catch (err) {
          console.error('[NATS] failed to process issue.priority_updated message', err)
        }
      }
    }
    void listen()

    return () => {
      sub.unsubscribe()
    }
  }, [connection])
}
