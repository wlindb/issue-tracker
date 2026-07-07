import { useEffect, useRef } from 'react'
import { type Comment } from '@/api/generated/issueTrackerAPI'
import { useNats } from '@/context/NatsContext'
import { useWorkspace } from '@/context/WorkspaceContext'

interface RawComment {
  ID: string
  IssueID: string
  AuthorID: string
  Body: string
  CreatedAt: string
  UpdatedAt: string
}

export interface CommentCreatedEvent {
  occurred_at: string
  payload: Comment
}

export interface ApiCommentCreatedEvent {
  occurred_at: string
  payload: RawComment
}

function toComment(raw: RawComment): Comment {
  return {
    id: raw.ID,
    issueId: raw.IssueID,
    authorId: raw.AuthorID,
    body: raw.Body,
    createdAt: raw.CreatedAt,
    updatedAt: raw.UpdatedAt,
  }
}

export function useCommentCreatedEvents(issueId: string | undefined, onCommentCreated: (event: CommentCreatedEvent) => void) {
  const { connection } = useNats()
  const { activeWorkspace } = useWorkspace()
  const activeWorkspaceRef = useRef(activeWorkspace)
  activeWorkspaceRef.current = activeWorkspace
  const onCommentCreatedRef = useRef(onCommentCreated)
  onCommentCreatedRef.current = onCommentCreated

  useEffect(() => {
    const workspace = activeWorkspaceRef.current
    if (!connection || !workspace || !issueId) return

    const sub = connection.subscribe(`workspaces.${workspace.id}.issues.${issueId}.comments.created`)

    const listen = async () => {
      for await (const msg of sub) {
        try {
          const raw = msg.json<ApiCommentCreatedEvent>()
          onCommentCreatedRef.current({ occurred_at: raw.occurred_at, payload: toComment(raw.payload) })
        } catch (err) {
          console.error('[NATS] failed to process comments.created message', err)
        }
      }
    }
    void listen()

    return () => {
      sub.unsubscribe()
    }
  }, [connection, issueId])
}
