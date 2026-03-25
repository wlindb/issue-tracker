import { IssueGroupSection } from '@/components/IssueGroupSection'
import { CURRENT_USER_ID, issues } from '@/data/mock'
import { groupIssuesByStatus } from '@/lib/groupIssuesByStatus'

export function MyIssuesPage() {
  const myIssues = issues.filter((issue) => issue.assigneeId === CURRENT_USER_ID)
  const groups = groupIssuesByStatus(myIssues)

  return (
    <div className="flex flex-col">
      <div className="border-b border-border px-6 py-4">
        <h1 className="text-lg font-semibold">My Issues</h1>
        <p className="text-sm text-muted-foreground">{myIssues.length} issues assigned to you</p>
      </div>
      <div className="py-2">
        {groups.length === 0 ? (
          <p className="px-6 py-8 text-sm text-muted-foreground">No issues assigned to you.</p>
        ) : (
          groups.map((group) => (
            <IssueGroupSection key={group.status} group={group} />
          ))
        )}
      </div>
    </div>
  )
}
