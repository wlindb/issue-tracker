import { useEffect, useState } from 'react'
import { PlusCircle } from 'lucide-react'
import { ProjectCard } from '@/components/ProjectCard'
import { Button } from '@/components/ui/button'
import { CreateProjectForm } from '@/components/CreateProjectForm'
import { cn } from '@/lib/utils'
import { listProjects, type Project } from '@/api/generated/issueTrackerAPI'

export function ProjectsPage() {
  const [projects, setProjects] = useState<Project[]>([])
  const [loading, setLoading] = useState(true)
  const [creating, setCreating] = useState(false)

  useEffect(() => {
    async function load() {
      const page = await listProjects()
      setProjects(page.items)
      setLoading(false)
    }
    load()
  }, [])

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

      {loading ? (
        <p className="px-6 py-8 text-sm text-muted-foreground">Loading…</p>
      ) : (
        <div className="grid grid-cols-1 gap-4 p-6 sm:grid-cols-2 lg:grid-cols-3">
          {projects.map((project) => (
            <ProjectCard key={project.id} project={project} />
          ))}
        </div>
      )}
    </div>
  )
}
