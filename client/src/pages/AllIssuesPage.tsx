import { useState } from 'react'
import { PlusCircle, Search } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { IssueRow } from '@/components/IssueRow'
import { CreateIssueForm } from '@/components/CreateIssueForm'
import { issues as mockIssues, projects, type Issue } from '@/data/mock'
import { useIssueSearch } from '@/hooks/useIssueSearch'
import { cn } from '@/lib/utils'

export function AllIssuesPage() {
  const [issues, setIssues] = useState<Issue[]>(mockIssues)
  const [creating, setCreating] = useState(false)
  const [query, setQuery] = useState('')

  const { results, isPending } = useIssueSearch(issues, query)

  function handleSave(issue: Issue) {
    setIssues([issue, ...issues])
    setCreating(false)
  }

  return (
    <div className="flex flex-col">
      <div className="flex items-center justify-between border-b border-border px-6 py-4">
        <div>
          <h1 className="text-lg font-semibold">All Issues</h1>
          <p className="text-sm text-muted-foreground">
            {query ? `${results.length} of ${issues.length} issues` : `${issues.length} issues total`}
          </p>
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

      <div className="border-b border-border px-6 py-3">
        <div className="relative">
          <Search className="absolute left-2.5 top-1/2 size-3.5 -translate-y-1/2 text-muted-foreground" />
          <Input
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Search issues…"
            aria-label="Search issues"
            className={cn('pl-8', isPending && 'opacity-60')}
          />
        </div>
      </div>

      {creating && (
        <CreateIssueForm
          projects={projects}
          onSave={handleSave}
          onCancel={() => setCreating(false)}
        />
      )}

      <div className="py-2">
        {results.length === 0 ? (
          <p className="px-6 py-8 text-sm text-muted-foreground">
            No issues match &ldquo;{query}&rdquo;.
          </p>
        ) : (
          results.map((issue) => <IssueRow key={issue.id} issue={issue} />)
        )}
      </div>
    </div>
  )
}
