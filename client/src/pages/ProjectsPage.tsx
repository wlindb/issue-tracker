import { useEffect, useState } from 'react'
import { PlusCircle } from 'lucide-react'
import { ProjectCard } from '@/components/ProjectCard'
import { Button } from '@/components/ui/button'
import { CreateProjectForm } from '@/components/CreateProjectForm'
import { cn } from '@/lib/utils'
import { listProjects, type Project } from '@/api/generated/issueTrackerAPI'
import { useWorkspace } from '@/context/WorkspaceContext'
import { useProjectCreatedEvents } from '@/hooks/useProjectCreatedEvents'

export function ProjectsPage() {
  const [projects, setProjects] = useState<Project[]>([])
  const [loading, setLoading] = useState(true)
  const [creating, setCreating] = useState(false)
  const { activeWorkspace } = useWorkspace()

  const upsertProject = (nextProject: Project) =>
    setProjects((previous) => {
      const existingIndex = previous.findIndex((existing) => existing.id === nextProject.id)
      if (existingIndex === -1) {
        return [nextProject, ...previous]
      }

      const updated = [...previous]
      updated[existingIndex] = nextProject
      return updated
    })

  useProjectCreatedEvents((event) => upsertProject(event.payload))

  useEffect(() => {
    if (!activeWorkspace) return
    const workspaceId = activeWorkspace.id
    async function load() {
      const page = await listProjects(workspaceId)
      setProjects(page.items)
      setLoading(false)
    }
    load()
  }, [activeWorkspace])

  function handleSave(project: Project) {
    upsertProject(project)
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
