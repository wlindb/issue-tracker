import { ProjectCard } from '@/components/ProjectCard'
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
          <ProjectCard key={project.id} project={project} />
        ))}
      </div>
    </div>
  )
}
