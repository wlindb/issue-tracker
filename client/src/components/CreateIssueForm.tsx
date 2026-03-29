import { useEffect, useRef, useState, type KeyboardEvent } from 'react'
import { type Issue, type Project } from '@/data/mock'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { createIssue } from '@/api/generated/issueTrackerAPI'

interface CreateIssueFormProps {
  projects: Project[]
  defaultProjectId?: string
  onSave: (issue: Issue) => void
  onCancel: () => void
}

export function CreateIssueForm({ projects, defaultProjectId, onSave, onCancel }: CreateIssueFormProps) {
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [projectId, setProjectId] = useState(defaultProjectId ?? '')
  const [submitting, setSubmitting] = useState(false)
  const titleRef = useRef<HTMLInputElement>(null)

  useEffect(() => {
    titleRef.current?.focus()
  }, [])

  async function handleSubmit() {
    const trimmedTitle = title.trim()
    if (!trimmedTitle || !projectId) return

    setSubmitting(true)
    try {
      const issue = await createIssue({
        projectId,
        title: trimmedTitle,
        description: description.trim() || null,
        status: 'backlog',
        priority: 'none',
      })
      onSave(issue)
    } finally {
      setSubmitting(false)
    }
  }

  function handleKeyDown(event: KeyboardEvent) {
    if (event.key === 'Escape') onCancel()
  }

  const canSubmit = title.trim().length > 0 && projectId !== ''

  return (
    <div className="border-b border-border px-6 py-4">
      <form
        onSubmit={(e) => { e.preventDefault(); handleSubmit() }}
        onKeyDown={handleKeyDown}
        aria-label="New issue"
        className="flex flex-col gap-3 rounded-lg border border-border bg-card p-4"
      >
        <h2 className="text-sm font-medium">New issue</h2>
        <div className="flex flex-col gap-1.5">
          <label htmlFor="issue-project" className="text-sm font-medium">
            Project <span className="text-destructive" aria-hidden>*</span>
          </label>
          <select
            id="issue-project"
            value={projectId}
            onChange={(e) => setProjectId(e.target.value)}
            disabled={!!defaultProjectId}
            required
            className="h-8 w-full min-w-0 rounded-lg border border-input bg-transparent px-2.5 py-1 text-sm transition-colors outline-none focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 disabled:pointer-events-none disabled:cursor-not-allowed disabled:opacity-50 dark:bg-input/30"
          >
            {!defaultProjectId && <option value="">Select a project…</option>}
            {projects.map((p) => (
              <option key={p.id} value={p.id}>{p.name}</option>
            ))}
          </select>
        </div>
        <div className="flex flex-col gap-1.5">
          <label htmlFor="issue-title" className="text-sm font-medium">
            Title <span className="text-destructive" aria-hidden>*</span>
          </label>
          <Input
            id="issue-title"
            ref={titleRef}
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="Issue title"
            required
          />
        </div>
        <div className="flex flex-col gap-1.5">
          <label htmlFor="issue-description" className="text-sm font-medium">
            Description
          </label>
          <textarea
            id="issue-description"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            placeholder="Add more details…"
            rows={3}
            className="w-full min-w-0 resize-none rounded-lg border border-input bg-transparent px-2.5 py-1.5 text-sm transition-colors outline-none placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 dark:bg-input/30"
          />
        </div>
        <div className="flex justify-end gap-2">
          <Button type="button" variant="ghost" size="sm" onClick={onCancel} disabled={submitting}>
            Cancel
          </Button>
          <Button type="submit" size="sm" disabled={!canSubmit || submitting}>
            {submitting ? 'Creating…' : 'Create issue'}
          </Button>
        </div>
      </form>
    </div>
  )
}
