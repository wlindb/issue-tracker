import type { Issue, Status } from '@/data/mock'

const STATUS_ORDER: Status[] = [
  'in_progress',
  'in_review',
  'todo',
  'backlog',
  'done',
  'cancelled',
]

export interface IssueGroup {
  status: Status
  issues: Issue[]
}

export function groupIssuesByStatus(issues: Issue[]): IssueGroup[] {
  const map = new Map<Status, Issue[]>()

  for (const issue of issues) {
    const group = map.get(issue.status) ?? []
    group.push(issue)
    map.set(issue.status, group)
  }

  return STATUS_ORDER
    .filter((status) => (map.get(status)?.length ?? 0) > 0)
    .map((status) => ({ status, issues: map.get(status)! }))
}
