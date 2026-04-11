import { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import { Separator } from '@/components/ui/separator'
import {
  getIssue,
  getProject,
  listComments,
  createComment,
  updateIssueTitle,
  updateIssueDescription,
  updateIssueStatus,
  updateIssuePriority,
  updateIssueAssignee,
  type Issue,
  type Project,
  type Comment,
  type IssueStatus,
  type IssuePriority,
} from '@/api/generated/issueTrackerAPI'
import { useWorkspace } from '@/context/WorkspaceContext'
import { EditableText } from '@/components/issue-detail/EditableText'
import { IssueBreadcrumbs } from '@/components/issue-detail/IssueBreadcrumbs'
import { IssueMetaSidebar } from '@/components/issue-detail/IssueMetaSidebar'
import { CommentSection } from '@/components/issue-detail/CommentSection'

export function IssueDetailPage() {
  const { issueId } = useParams<{ issueId: string }>()
  const { activeWorkspace } = useWorkspace()

  const [issue, setIssue] = useState<Issue | null>(null)
  const [project, setProject] = useState<Project | null>(null)
  const [comments, setComments] = useState<Comment[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(false)

  useEffect(() => {
    if (!activeWorkspace || !issueId) return
    const workspaceId = activeWorkspace.id
    const resolvedIssueId = issueId
    setLoading(true)
    setError(false)
    async function load() {
      try {
        const fetchedIssue = await getIssue(workspaceId, resolvedIssueId)
        const [fetchedProject, commentsPage] = await Promise.all([
          getProject(workspaceId, fetchedIssue.projectId),
          listComments(workspaceId, resolvedIssueId),
        ])
        setIssue(fetchedIssue)
        setProject(fetchedProject)
        setComments(commentsPage.items)
      } catch {
        setError(true)
      } finally {
        setLoading(false)
      }
    }
    load()
  }, [activeWorkspace, issueId])

  async function handleTitleSave(title: string) {
    if (!activeWorkspace || !issue) return
    const updated = await updateIssueTitle(activeWorkspace.id, issue.id, { title })
    setIssue(updated)
  }

  async function handleDescriptionSave(description: string) {
    if (!activeWorkspace || !issue) return
    const updated = await updateIssueDescription(activeWorkspace.id, issue.id, {
      description: description || null,
    })
    setIssue(updated)
  }

  async function handleStatusChange(status: IssueStatus) {
    if (!activeWorkspace || !issue) return
    const updated = await updateIssueStatus(activeWorkspace.id, issue.id, { status })
    setIssue(updated)
  }

  async function handlePriorityChange(priority: IssuePriority) {
    if (!activeWorkspace || !issue) return
    const updated = await updateIssuePriority(activeWorkspace.id, issue.id, { priority })
    setIssue(updated)
  }

  async function handleAssigneeChange(assigneeId: string | null) {
    if (!activeWorkspace || !issue) return
    const updated = await updateIssueAssignee(activeWorkspace.id, issue.id, { assigneeId })
    setIssue(updated)
  }

  function handleLabelsChange(labels: string[]) {
    if (!issue) return
    setIssue({ ...issue, labels })
  }

  async function handleAddComment(body: string) {
    if (!activeWorkspace || !issue) return
    const comment = await createComment(activeWorkspace.id, issue.id, { body })
    setComments((prev) => [...prev, comment])
  }

  if (loading) {
    return (
      <div className="flex h-full items-center justify-center text-muted-foreground">
        Loading…
      </div>
    )
  }

  if (error || !issue || !project) {
    return (
      <div className="flex h-full items-center justify-center text-muted-foreground">
        Issue not found.
      </div>
    )
  }

  return (
    <div className="flex flex-col">
      <div className="mx-auto w-full max-w-5xl px-6 py-6">
        <IssueBreadcrumbs project={project} issue={issue} />

        <div className="mt-6 flex gap-8">
          {/* Main content */}
          <div className="min-w-0 flex-1 flex flex-col gap-6">
            <EditableText
              value={issue.title}
              onSave={handleTitleSave}
              placeholder="Issue title..."
              className="text-xl font-semibold"
              inputClassName="text-xl font-semibold"
            />

            <div className="flex flex-col gap-1">
              <span className="px-2.5 text-xs font-medium uppercase tracking-wide text-muted-foreground">
                Description
              </span>
              <EditableText
                value={issue.description ?? ''}
                onSave={handleDescriptionSave}
                placeholder="Add a description..."
                multiline
                className="min-h-[80px] leading-relaxed"
              />
            </div>

            <Separator />

            <CommentSection comments={comments} onAddComment={handleAddComment} />
          </div>

          {/* Meta sidebar */}
          <aside className="w-52 shrink-0">
            <IssueMetaSidebar
              issue={issue}
              users={[]}
              onStatusChange={handleStatusChange}
              onPriorityChange={handlePriorityChange}
              onAssigneeChange={handleAssigneeChange}
              onLabelsChange={handleLabelsChange}
            />
          </aside>
        </div>
      </div>
    </div>
  )
}
