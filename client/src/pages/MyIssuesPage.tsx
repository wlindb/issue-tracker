import { useEffect, useState } from 'react'
import { IssueGroupSection } from '@/components/IssueGroupSection'
import { groupIssuesByStatus } from '@/lib/groupIssuesByStatus'
import { getMe, listProjects, listIssues, type Issue } from '@/api/generated/issueTrackerAPI'

export function MyIssuesPage() {
  const [issues, setIssues] = useState<Issue[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    async function load() {
      const [me, projectPage] = await Promise.all([getMe(), listProjects()])
      const issuePages = await Promise.all(
        projectPage.items.map((p) => listIssues({ project_id: p.id, assigneeId: me.id }))
      )
      setIssues(issuePages.flatMap((page) => page.items))
      setLoading(false)
    }
    load()
  }, [])

  const groups = groupIssuesByStatus(issues)

  return (
    <div className="flex flex-col">
      <div className="border-b border-border px-6 py-4">
        <h1 className="text-lg font-semibold">My Issues</h1>
        <p className="text-sm text-muted-foreground">{issues.length} issues assigned to you</p>
      </div>
      <div className="py-2">
        {loading ? (
          <p className="px-6 py-8 text-sm text-muted-foreground">Loading…</p>
        ) : groups.length === 0 ? (
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
