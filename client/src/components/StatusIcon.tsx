import {
  CheckCircle2Icon,
  CircleDashedIcon,
  CircleDotIcon,
  CircleIcon,
  XCircleIcon,
} from 'lucide-react'
import { cn } from '@/lib/utils'
import type { IssueStatus } from '@/api/generated/issueTrackerAPI'

const config: Record<IssueStatus, { icon: React.ElementType; className: string }> = {
  in_progress: { icon: CircleDotIcon, className: 'text-blue-500' },
  todo: { icon: CircleIcon, className: 'text-muted-foreground' },
  backlog: { icon: CircleDashedIcon, className: 'text-muted-foreground' },
  done: { icon: CheckCircle2Icon, className: 'text-emerald-500' },
  cancelled: { icon: XCircleIcon, className: 'text-muted-foreground' },
}

export const STATUS_LABEL: Record<IssueStatus, string> = {
  in_progress: 'In Progress',
  todo: 'Todo',
  backlog: 'Backlog',
  done: 'Done',
  cancelled: 'Cancelled',
}

interface StatusIconProps {
  status: IssueStatus
  className?: string
}

export function StatusIcon({ status, className }: StatusIconProps) {
  const { icon: Icon, className: colorClass } = config[status]
  return <Icon className={cn('size-3.5 shrink-0', colorClass, className)} />
}
