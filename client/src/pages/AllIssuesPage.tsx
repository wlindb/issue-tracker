import { IssueRow } from '@/components/IssueRow'
import { issues } from '@/data/mock'

export function AllIssuesPage() {
  return (
    <div className="flex flex-col">
      <div className="border-b border-border px-6 py-4">
        <h1 className="text-lg font-semibold">All Issues</h1>
        <p className="text-sm text-muted-foreground">{issues.length} issues total</p>
      </div>
      <div className="py-2">
        {issues.map((issue) => (
          <IssueRow key={issue.id} issue={issue} />
        ))}
      </div>
    </div>
  )
}
