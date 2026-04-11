import { useEffect, useRef, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { PlusCircle } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { useWorkspace } from '@/context/WorkspaceContext'

export function CreateWorkspacePage() {
  const [name, setName] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const nameRef = useRef<HTMLInputElement>(null)
  const navigate = useNavigate()
  const { createWorkspace } = useWorkspace()

  useEffect(() => {
    nameRef.current?.focus()
  }, [])

  async function handleSubmit(e: React.SubmitEvent<HTMLFormElement>) {
    e.preventDefault();

    const trimmedName = name.trim()
    if (!trimmedName) return

    setSubmitting(true)
    try {
      await createWorkspace({ name: trimmedName })
      navigate('/', { replace: true })
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="flex min-h-screen w-full items-center justify-center bg-background p-4">
      <div className="w-full max-w-xl">
        <form
          onSubmit={handleSubmit}
          className="flex flex-col gap-6 rounded-lg border border-border bg-card p-8"
        >
          <div className="flex flex-col gap-2">
            <h1 className="text-base font-semibold">Create new workspace</h1>
            <p className="text-sm text-muted-foreground">
              Workspaces are shared environments where teams can work on projects, cycles and issues.
            </p>
          </div>

          <div className="flex flex-col gap-1.5">
            <label htmlFor="workspace-name" className="text-sm font-medium">
              Name <span className="text-destructive" aria-hidden>*</span>
            </label>
            <Input
              id="workspace-name"
              ref={nameRef}
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="e.g. My Team"
              required
            />
          </div>

          <Button type="submit" disabled={!name.trim() || submitting}>
            <PlusCircle className="size-4" />
            {submitting ? 'Creating…' : 'Create workspace'}
          </Button>
        </form>
      </div>
    </div>
  )
}
