import { useState } from 'react'
import { useParams } from 'react-router-dom'
import { Separator } from '@/components/ui/separator'
import {
  comments as allComments,
  issues,
  projects,
  users,
  type Comment,
  type Priority,
  type Status,
} from '@/data/mock'
import { EditableText } from '@/components/issue-detail/EditableText'
import { IssueBreadcrumbs } from '@/components/issue-detail/IssueBreadcrumbs'
import { IssueMetaSidebar } from '@/components/issue-detail/IssueMetaSidebar'
import { CommentSection } from '@/components/issue-detail/CommentSection'

export function IssueDetailPage() {
  const { identifier } = useParams<{ identifier: string }>()

  const baseIssue = issues.find((i) => i.identifier === identifier)
  const project = baseIssue ? projects.find((p) => p.id === baseIssue.projectId) : undefined

  const [title, setTitle] = useState(baseIssue?.title ?? '')
  const [description, setDescription] = useState(baseIssue?.description ?? '')
  const [status, setStatus] = useState<Status>(baseIssue?.status ?? 'backlog')
  const [priority, setPriority] = useState<Priority>(baseIssue?.priority ?? 'none')
  const [assigneeId, setAssigneeId] = useState(baseIssue?.assigneeId ?? null)
  const [labels, setLabels] = useState(baseIssue?.labels ?? [])
  const [issueComments, setIssueComments] = useState<Comment[]>(
    allComments.filter((c) => c.issueId === baseIssue?.id),
  )

  if (!baseIssue || !project) {
    return (
      <div className="flex h-full items-center justify-center text-muted-foreground">
        Issue not found.
      </div>
    )
  }

  function handleAddComment(body: string) {
    setIssueComments((prev) => [
      ...prev,
      {
        id: `comment-${Date.now()}`,
        issueId: baseIssue!.id,
        authorId: 'user-1',
        authorName: 'Alice',
        body,
        createdAt: new Date().toISOString(),
      },
    ])
  }

  function handleAssigneeChange(newId: string | null) {
    setAssigneeId(newId)
  }

  const issue = { ...baseIssue, title, description, status, priority, assigneeId, labels }

  return (
    <div className="flex flex-col">
      <div className="mx-auto w-full max-w-5xl px-6 py-6">
        <IssueBreadcrumbs project={project} issue={issue} />

        <div className="mt-6 flex gap-8">
          {/* Main content */}
          <div className="min-w-0 flex-1 flex flex-col gap-6">
            <EditableText
              value={title}
              onSave={setTitle}
              placeholder="Issue title..."
              className="text-xl font-semibold"
              inputClassName="text-xl font-semibold"
            />

            <div className="flex flex-col gap-1">
              <span className="px-2.5 text-xs font-medium uppercase tracking-wide text-muted-foreground">
                Description
              </span>
              <EditableText
                value={description ?? ''}
                onSave={setDescription}
                placeholder="Add a description..."
                multiline
                className="min-h-[80px] leading-relaxed"
              />
            </div>

            <Separator />

            <CommentSection comments={issueComments} onAddComment={handleAddComment} />
          </div>

          {/* Meta sidebar */}
          <aside className="w-52 shrink-0">
            <IssueMetaSidebar
              issue={issue}
              users={users}
              onStatusChange={setStatus}
              onPriorityChange={setPriority}
              onAssigneeChange={handleAssigneeChange}
              onLabelsChange={setLabels}
            />
          </aside>
        </div>
      </div>
    </div>
  )
}
