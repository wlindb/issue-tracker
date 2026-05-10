import type { IssueStatus } from '@/api/generated/issueTrackerAPI'

export const STATUS_LABEL: Record<IssueStatus, string> = {
  in_progress: 'In Progress',
  todo: 'Todo',
  backlog: 'Backlog',
  done: 'Done',
  cancelled: 'Cancelled',
}
