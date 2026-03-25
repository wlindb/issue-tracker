import { Link } from 'react-router-dom'
import { Badge } from '@/components/ui/badge'
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { projects } from '@/data/mock'

export function ProjectsPage() {
  return (
    <div className="flex flex-col">
      <div className="border-b border-border px-6 py-4">
        <h1 className="text-lg font-semibold">Projects</h1>
        <p className="text-sm text-muted-foreground">{projects.length} projects</p>
      </div>
      <div className="grid grid-cols-1 gap-4 p-6 sm:grid-cols-2 lg:grid-cols-3">
        {projects.map((project) => (
          <Link key={project.id} to={`/projects/${project.id}`} className="block">
            <Card className="h-full transition-colors hover:bg-muted/40">
              <CardHeader className="pb-2">
                <div className="flex items-center gap-2">
                  <span
                    className="size-2.5 shrink-0 rounded-full"
                    style={{ backgroundColor: project.color }}
                  />
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
              <CardFooter>
                <Badge variant="secondary" className="text-xs">
                  {project.issueCount} issues
                </Badge>
              </CardFooter>
            </Card>
          </Link>
        ))}
      </div>
    </div>
  )
}
