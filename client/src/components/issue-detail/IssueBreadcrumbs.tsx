import { Link } from 'react-router-dom'
import { ChevronRightIcon } from 'lucide-react'
import type { Issue, Project } from '@/data/mock'

interface IssueBreadcrumbsProps {
  project: Project
  issue: Issue
}

export function IssueBreadcrumbs({ project, issue }: IssueBreadcrumbsProps) {
  return (
    <nav className="flex items-center gap-1 text-sm text-muted-foreground">
      <Link to="/projects" className="hover:text-foreground transition-colors">
        Projects
      </Link>
      <ChevronRightIcon className="size-3.5 shrink-0" />
      <Link
        to={`/projects/${project.identifier}`}
        className="hover:text-foreground transition-colors"
      >
        {project.name}
      </Link>
      <ChevronRightIcon className="size-3.5 shrink-0" />
      <span className="font-mono text-foreground">{issue.identifier}</span>
    </nav>
  )
}
