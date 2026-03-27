import { useState } from 'react'
import { PlusCircle } from 'lucide-react'
import { ProjectCard } from '@/components/ProjectCard'
import { projects as mockProjects, type Project } from '@/data/mock'
import { Button } from '@/components/ui/button'
import { CreateProjectForm } from '@/components/CreateProjectForm'
import { cn } from '@/lib/utils'

export function ProjectsPage() {
  const [projects, setProjects] = useState<Project[]>(mockProjects)
  const [creating, setCreating] = useState(false)

  function handleSave(project: Project) {
    setProjects([project, ...projects])
    setCreating(false)
  }

  return (
    <div className="flex flex-col">
      <div className="flex items-center justify-between border-b border-border px-6 py-4">
        <div>
          <h1 className="text-lg font-semibold">Projects</h1>
          <p className="text-sm text-muted-foreground">{projects.length} projects</p>
        </div>
        <Button
          variant="ghost"
          size="icon"
          aria-label="Create new project"
          onClick={() => setCreating(true)}
          className={cn(creating && 'text-primary')}
        >
          <PlusCircle className="size-5" />
        </Button>
      </div>

      {creating && (
        <CreateProjectForm
          onSave={handleSave}
          onCancel={() => setCreating(false)}
        />
      )}

      <div className="grid grid-cols-1 gap-4 p-6 sm:grid-cols-2 lg:grid-cols-3">
        {projects.map((project) => (
          <ProjectCard key={project.id} project={project} />
        ))}
      </div>
    </div>
  )
}
