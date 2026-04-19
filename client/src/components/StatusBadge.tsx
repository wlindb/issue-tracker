import { Badge } from '@/components/ui/badge'
import { cn } from '@/lib/utils'
import type { IssueStatus } from '@/api/generated/issueTrackerAPI'
import { STATUS_LABEL } from './StatusIcon'

const statusColors: Record<IssueStatus, string> = {
  in_progress: 'bg-blue-500/10 text-blue-600 dark:text-blue-400',
  todo: 'bg-muted text-muted-foreground',
  backlog: 'bg-muted text-muted-foreground',
  done: 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-400',
  cancelled: 'bg-muted text-muted-foreground line-through',
}

interface StatusBadgeProps {
  status: IssueStatus
  className?: string
}

export function StatusBadge({ status, className }: StatusBadgeProps) {
  return (
    <Badge
      variant="secondary"
      className={cn('shrink-0 border-0', statusColors[status], className)}
    >
      {STATUS_LABEL[status]}
    </Badge>
  )
}
