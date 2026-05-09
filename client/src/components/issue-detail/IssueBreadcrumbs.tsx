import { Link } from 'react-router-dom'
import { ChevronRightIcon } from 'lucide-react'
import type { Project } from '@/api/generated/issueTrackerAPI'

interface IssueBreadcrumbsProps {
  project: Project
}

export function IssueBreadcrumbs({ project }: IssueBreadcrumbsProps) {
  return (
    <nav className="flex items-center gap-1 text-sm text-muted-foreground">
      <Link to="/projects" className="hover:text-foreground transition-colors">
        Projects
      </Link>
      <ChevronRightIcon className="size-3.5 shrink-0" />
      <Link
        to={`/projects/${project.id}`}
        className="hover:text-foreground transition-colors"
      >
        {project.name}
      </Link>
    </nav>
  )
}
