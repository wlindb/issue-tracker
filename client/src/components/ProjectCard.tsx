import { Link } from 'react-router-dom'
import { Badge } from '@/components/ui/badge'
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import type { Project } from '@/api/generated/issueTrackerAPI'

interface ProjectCardProps {
  project: Project
}

export function ProjectCard({ project }: ProjectCardProps) {
  return (
    <Link to={`/projects/${project.identifier}`} className="block">
      <Card className="flex h-full flex-col transition-colors hover:bg-muted/40">
        <CardHeader className="pb-2">
          <div className="flex items-center gap-2">
            <CardTitle className="text-sm font-semibold">{project.name}</CardTitle>
            <Badge variant="outline" className="ml-auto font-mono text-xs">
              {project.identifier}
            </Badge>
          </div>
        </CardHeader>
        <CardContent className="pb-3">
          <p className="line-clamp-2 text-sm text-muted-foreground">
            {project.description}
          </p>
        </CardContent>
        <CardFooter className="mt-auto">
          <Badge variant="secondary" className="text-xs">
            {project.issueCount} issues
          </Badge>
        </CardFooter>
      </Card>
    </Link>
  )
}
