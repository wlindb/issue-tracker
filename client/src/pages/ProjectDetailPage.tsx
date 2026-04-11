import { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import { PlusCircle } from 'lucide-react'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { IssueGroupSection } from '@/components/IssueGroupSection'
import { CreateIssueForm } from '@/components/CreateIssueForm'
import { getProject, listIssues, type Issue, type Project } from '@/api/generated/issueTrackerAPI'
import { useWorkspace } from '@/context/WorkspaceContext'
import { groupIssuesByStatus } from '@/lib/groupIssuesByStatus'
import { cn } from '@/lib/utils'

export function ProjectDetailPage() {
  const { projectId } = useParams<{ projectId: string }>()
  const { activeWorkspace } = useWorkspace()

  const [project, setProject] = useState<Project | null>(null)
  const [issues, setIssues] = useState<Issue[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(false)
  const [creating, setCreating] = useState(false)

  useEffect(() => {
    if (!activeWorkspace || !projectId) return
    const workspaceId = activeWorkspace.id
    setLoading(true)
    setError(false)
    async function load() {
      try {
        const [fetchedProject, issuePage] = await Promise.all([
          getProject(workspaceId, projectId!),
          listIssues(workspaceId, { project_id: projectId! }),
        ])
        setProject(fetchedProject)
        setIssues(issuePage.items)
      } catch {
        setError(true)
      } finally {
        setLoading(false)
      }
    }
    load()
  }, [activeWorkspace, projectId])

  function handleSave(issue: Issue) {
    setIssues((prev) => [issue, ...prev])
    setCreating(false)
  }

  if (loading) {
    return (
      <p className="px-6 py-8 text-sm text-muted-foreground">Loading…</p>
    )
  }

  if (error || !project) {
    return (
      <div className="px-6 py-8 text-sm text-muted-foreground">Project not found.</div>
    )
  }

  const groups = groupIssuesByStatus(issues)

  return (
    <div className="flex flex-col">
      <div className="flex items-start justify-between border-b border-border px-6 py-4">
        <div>
          <div className="flex items-center gap-2">
            <h1 className="text-lg font-semibold">{project.name}</h1>
            <Badge variant="outline" className="font-mono text-xs">
              {project.identifier}
            </Badge>
          </div>
          <p className="mt-1 text-sm text-muted-foreground">{project.description}</p>
        </div>
        <Button
          variant="ghost"
          size="icon"
          aria-label="Create new issue"
          onClick={() => setCreating(true)}
          className={cn('shrink-0', creating && 'text-primary')}
        >
          <PlusCircle className="size-5" />
        </Button>
      </div>

      {creating && (
        <CreateIssueForm
          projects={[project]}
          defaultProjectId={project.id}
          onSave={handleSave}
          onCancel={() => setCreating(false)}
        />
      )}

      <div className="py-2">
        {groups.length === 0 ? (
          <p className="px-6 py-8 text-sm text-muted-foreground">No issues in this project.</p>
        ) : (
          groups.map((group) => (
            <IssueGroupSection key={group.status} group={group} />
          ))
        )}
      </div>
    </div>
  )
}
