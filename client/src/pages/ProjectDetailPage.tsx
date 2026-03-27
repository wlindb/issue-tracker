import { useState } from 'react'
import { useParams } from 'react-router-dom'
import { PlusCircle } from 'lucide-react'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { IssueGroupSection } from '@/components/IssueGroupSection'
import { CreateIssueForm } from '@/components/CreateIssueForm'
import { issues as mockIssues, projects, type Issue } from '@/data/mock'
import { groupIssuesByStatus } from '@/lib/groupIssuesByStatus'
import { cn } from '@/lib/utils'

export function ProjectDetailPage() {
  const { identifier } = useParams<{ identifier: string }>()
  const project = projects.find((p) => p.identifier === identifier)

  const [issues, setIssues] = useState<Issue[]>(mockIssues)
  const [creating, setCreating] = useState(false)

  if (!project) {
    return (
      <div className="px-6 py-8 text-sm text-muted-foreground">Project not found.</div>
    )
  }

  function handleSave(issue: Issue) {
    setIssues([issue, ...issues])
    setCreating(false)
  }

  const projectIssues = issues.filter((issue) => issue.projectId === project.id)
  const groups = groupIssuesByStatus(projectIssues)

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
          projects={projects}
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
