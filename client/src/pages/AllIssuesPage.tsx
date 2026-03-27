import { useState } from 'react'
import { PlusCircle } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { IssueRow } from '@/components/IssueRow'
import { CreateIssueForm } from '@/components/CreateIssueForm'
import { issues as mockIssues, projects, type Issue } from '@/data/mock'
import { cn } from '@/lib/utils'

export function AllIssuesPage() {
  const [issues, setIssues] = useState<Issue[]>(mockIssues)
  const [creating, setCreating] = useState(false)

  function handleSave(issue: Issue) {
    setIssues([issue, ...issues])
    setCreating(false)
  }

  return (
    <div className="flex flex-col">
      <div className="flex items-center justify-between border-b border-border px-6 py-4">
        <div>
          <h1 className="text-lg font-semibold">All Issues</h1>
          <p className="text-sm text-muted-foreground">{issues.length} issues total</p>
        </div>
        <Button
          variant="ghost"
          size="icon"
          aria-label="Create new issue"
          onClick={() => setCreating(true)}
          className={cn(creating && 'text-primary')}
        >
          <PlusCircle className="size-5" />
        </Button>
      </div>

      {creating && (
        <CreateIssueForm
          projects={projects}
          onSave={handleSave}
          onCancel={() => setCreating(false)}
        />
      )}

      <div className="py-2">
        {issues.map((issue) => (
          <IssueRow key={issue.id} issue={issue} />
        ))}
      </div>
    </div>
  )
}
