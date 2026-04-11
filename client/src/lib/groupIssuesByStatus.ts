import type { Issue, IssueStatus } from '@/api/generated/issueTrackerAPI'

const STATUS_ORDER: IssueStatus[] = [
  'in_progress',
  'todo',
  'backlog',
  'done',
  'cancelled',
]

export interface IssueGroup {
  status: IssueStatus
  issues: Issue[]
}

export function groupIssuesByStatus(issues: Issue[]): IssueGroup[] {
  const map = new Map<IssueStatus, Issue[]>()

  for (const issue of issues) {
    const group = map.get(issue.status) ?? []
    group.push(issue)
    map.set(issue.status, group)
  }

  return STATUS_ORDER
    .filter((status) => (map.get(status)?.length ?? 0) > 0)
    .map((status) => ({ status, issues: map.get(status)! }))
}
