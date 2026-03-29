import { useEffect, useRef, useState, type KeyboardEvent } from 'react'
import { type Project } from '@/data/mock'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { createProject } from '@/api/generated/issueTrackerAPI'

interface CreateProjectFormProps {
  onSave: (project: Project) => void
  onCancel: () => void
}

export function CreateProjectForm({ onSave, onCancel }: CreateProjectFormProps) {
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const nameRef = useRef<HTMLInputElement>(null)

  useEffect(() => {
    nameRef.current?.focus()
  }, [])

  async function handleSubmit() {
    const trimmedName = name.trim()
    if (!trimmedName) return

    setSubmitting(true)
    try {
      const project = await createProject({
        name: trimmedName,
        description: description.trim() || null,
      })
      onSave(project)
    } finally {
      setSubmitting(false)
    }
  }

  function handleKeyDown(event: KeyboardEvent) {
    if (event.key === 'Escape') onCancel()
  }

  return (
    <div className="border-b border-border px-6 py-4">
      <form
        onSubmit={(e) => { e.preventDefault(); handleSubmit() }}
        onKeyDown={handleKeyDown}
        aria-label="New project"
        className="flex flex-col gap-3 rounded-lg border border-border bg-card p-4"
      >
        <h2 className="text-sm font-medium">New project</h2>
        <div className="flex flex-col gap-1.5">
          <label htmlFor="project-name" className="text-sm font-medium">
            Name <span className="text-destructive" aria-hidden>*</span>
          </label>
          <Input
            id="project-name"
            ref={nameRef}
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="e.g. Backend"
            required
          />
        </div>
        <div className="flex flex-col gap-1.5">
          <label htmlFor="project-description" className="text-sm font-medium">
            Description
          </label>
          <textarea
            id="project-description"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            placeholder="What is this project about?"
            rows={3}
            className="w-full min-w-0 resize-none rounded-lg border border-input bg-transparent px-2.5 py-1.5 text-sm transition-colors outline-none placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 dark:bg-input/30"
          />
        </div>
        <div className="flex justify-end gap-2">
          <Button type="button" variant="ghost" size="sm" onClick={onCancel} disabled={submitting}>
            Cancel
          </Button>
          <Button type="submit" size="sm" disabled={!name.trim() || submitting}>
            {submitting ? 'Creating…' : 'Create project'}
          </Button>
        </div>
      </form>
    </div>
  )
}
