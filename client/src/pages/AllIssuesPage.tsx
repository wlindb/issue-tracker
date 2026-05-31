import { useEffect, useState } from 'react'
import { PlusCircle, Search } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { IssueRow } from '@/components/IssueRow'
import { CreateIssueForm } from '@/components/CreateIssueForm'
import { useIssueSearch } from '@/hooks/useIssueSearch'
import { cn } from '@/lib/utils'
import { listProjects, listIssues, type Issue, type Project } from '@/api/generated/issueTrackerAPI'
import { useWorkspace } from '@/context/WorkspaceContext'
import { useIssueCreatedEvents } from '@/hooks/useIssueCreatedEvents'

export function AllIssuesPage() {
  const [issues, setIssues] = useState<Issue[]>([])
  const [projects, setProjects] = useState<Project[]>([])
  const [loading, setLoading] = useState(true)
  const [creating, setCreating] = useState(false)
  const [query, setQuery] = useState('')
  const { activeWorkspace } = useWorkspace()

  const { results, isPending } = useIssueSearch(issues, query)
  const prependIfMissing = (nextIssue: Issue) =>
    setIssues((previous) =>
      previous.some((existing) => existing.id === nextIssue.id) ? previous : [nextIssue, ...previous]
    )

  useIssueCreatedEvents((event) => {
    const issue: Issue = {
      id: event.Payload.ID,
      identifier: event.Payload.Identifier,
      projectId: event.Payload.ProjectID,
      title: event.Payload.Title,
      description: event.Payload.Description,
      status: event.Payload.Status as Issue['status'],
      priority: event.Payload.Priority as Issue['priority'],
      labels: event.Payload.Labels,
      assigneeId: event.Payload.AssigneeID,
      reporterId: event.Payload.ReporterID,
      createdAt: event.Payload.CreatedAt,
      updatedAt: event.Payload.UpdatedAt,
    }
    prependIfMissing(issue)
  })

  useEffect(() => {
    if (!activeWorkspace) return
    const workspaceId = activeWorkspace.id
    async function load() {
      const projectPage = await listProjects(workspaceId)
      setProjects(projectPage.items)
      const issuePages = await Promise.all(
        projectPage.items.map((p) => listIssues(workspaceId, { project_id: p.id }))
      )
      setIssues(issuePages.flatMap((page) => page.items))
      setLoading(false)
    }
    load()
  }, [activeWorkspace])

  function handleSave(issue: Issue) {
    prependIfMissing(issue)
    setCreating(false)
  }

  return (
    <div className="flex flex-col">
      <div className="flex items-center justify-between border-b border-border px-6 py-4">
        <div>
          <h1 className="text-lg font-semibold">All Issues</h1>
          <p className="text-sm text-muted-foreground">
            {query ? `${results.length} of ${issues.length} issues` : `${issues.length} issues total`}
          </p>
        </div>
        <Button
          variant="ghost"
          size="icon"
          aria-label="Create new issue"
          onClick={() => setCreating(true)}
          className={cn(creating && 'text-primary')}
        >
          <PlusCircle className="size-5" />
        </Button>
      </div>

      <div className="border-b border-border px-6 py-3">
        <div className="relative">
          <Search className="absolute left-2.5 top-1/2 size-3.5 -translate-y-1/2 text-muted-foreground" />
          <Input
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Search issues…"
            aria-label="Search issues"
            className={cn('pl-8', isPending && 'opacity-60')}
          />
        </div>
      </div>

      {creating && (
        <CreateIssueForm
          projects={projects}
          onSave={handleSave}
          onCancel={() => setCreating(false)}
        />
      )}

      <div className="py-2">
        {loading ? (
          <p className="px-6 py-8 text-sm text-muted-foreground">Loading…</p>
        ) : results.length === 0 ? (
          query ? (
            <p className="px-6 py-8 text-sm text-muted-foreground">
              No issues match &ldquo;{query}&rdquo;.
            </p>
          ) : (
            <p className="px-6 py-8 text-sm text-muted-foreground">No issues yet.</p>
          )
        ) : (
          results.map((issue) => <IssueRow key={issue.id} issue={issue} />)
        )}
      </div>
    </div>
  )
}
