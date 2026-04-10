import { useEffect, useState } from 'react'
import { IssueGroupSection } from '@/components/IssueGroupSection'
import { groupIssuesByStatus } from '@/lib/groupIssuesByStatus'
import { listProjects, listIssues, type Issue } from '@/api/generated/issueTrackerAPI'
import { useWorkspace } from '@/context/WorkspaceContext'
import { useKeycloak } from '@/auth/KeycloakProvider'

export function MyIssuesPage() {
  const [issues, setIssues] = useState<Issue[]>([])
  const [loading, setLoading] = useState(true)
  const { activeWorkspace } = useWorkspace()
  const { keycloak } = useKeycloak()

  useEffect(() => {
    if (!activeWorkspace) return
    const workspaceId = activeWorkspace.id
    const userId = keycloak.tokenParsed?.sub
    if (!userId) return
    async function load() {
      const projectPage = await listProjects(workspaceId)
      const issuePages = await Promise.all(
        projectPage.items.map((p) => listIssues(workspaceId, { project_id: p.id, assigneeId: userId }))
      )
      setIssues(issuePages.flatMap((page) => page.items))
      setLoading(false)
    }
    load()
  }, [activeWorkspace, keycloak.tokenParsed?.sub])

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
