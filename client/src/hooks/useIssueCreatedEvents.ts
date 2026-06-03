import { useEffect, useRef } from 'react'
import { type Issue } from '@/api/generated/issueTrackerAPI'
import { useNats } from '@/context/NatsContext'
import { useWorkspace } from '@/context/WorkspaceContext'

interface EventIssue {
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

export interface ApiIssueCreatedEvent {
  occurred_at: string
  Payload: EventIssue
}

function toIssue(raw: EventIssue): Issue {
  return {
    id: raw.ID,
    identifier: raw.Identifier,
    title: raw.Title,
    description: raw.Description,
    status: raw.Status as Issue['status'],
    priority: raw.Priority as Issue['priority'],
    labels: raw.Labels,
    assigneeId: raw.AssigneeID,
    projectId: raw.ProjectID,
    reporterId: raw.ReporterID,
    createdAt: raw.CreatedAt,
    updatedAt: raw.UpdatedAt,
  }
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
          const raw = msg.json<ApiIssueCreatedEvent>()
          onIssueCreatedRef.current({ occurred_at: raw.occurred_at, Payload: toIssue(raw.Payload) })
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
