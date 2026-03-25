import { useParams } from 'react-router-dom'
import { Badge } from '@/components/ui/badge'
import { IssueGroupSection } from '@/components/IssueGroupSection'
import { issues, projects } from '@/data/mock'
import { groupIssuesByStatus } from '@/lib/groupIssuesByStatus'

export function ProjectDetailPage() {
  const { id } = useParams<{ id: string }>()
  const project = projects.find((p) => p.id === id)

  if (!project) {
    return (
      <div className="px-6 py-8 text-sm text-muted-foreground">Project not found.</div>
    )
  }

  const projectIssues = issues.filter((issue) => issue.projectId === project.id)
  const groups = groupIssuesByStatus(projectIssues)

  return (
    <div className="flex flex-col">
      <div className="border-b border-border px-6 py-4">
        <div className="flex items-center gap-2">
          <span
            className="size-2.5 shrink-0 rounded-full"
            style={{ backgroundColor: project.color }}
          />
          <h1 className="text-lg font-semibold">{project.name}</h1>
          <Badge variant="outline" className="font-mono text-xs">
            {project.identifier}
          </Badge>
        </div>
        <p className="mt-1 text-sm text-muted-foreground">{project.description}</p>
      </div>
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
