import { Link } from 'react-router-dom'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import type { Issue } from '@/data/mock'
import { PriorityIcon } from './PriorityIcon'
import { StatusBadge } from './StatusBadge'

interface IssueRowProps {
  issue: Issue
}

function getInitials(name: string): string {
  return name
    .split(' ')
    .map((part) => part[0])
    .slice(0, 2)
    .join('')
    .toUpperCase()
}

export function IssueRow({ issue }: IssueRowProps) {
  return (
    <Link
      to={`/issues/${issue.identifier}`}
      className="flex h-9 items-center gap-3 border-b border-border/50 px-4 hover:bg-muted/40 transition-colors"
    >
      <PriorityIcon priority={issue.priority} />
      <span className="w-[72px] shrink-0 font-mono text-xs text-muted-foreground">
        {issue.identifier}
      </span>
      <span className="flex-1 truncate text-sm">{issue.title}</span>
      <StatusBadge status={issue.status} />
      {issue.assigneeName ? (
        <Avatar size="sm" className="shrink-0">
          <AvatarFallback>{getInitials(issue.assigneeName)}</AvatarFallback>
        </Avatar>
      ) : (
        <div className="size-6 shrink-0" />
      )}
    </Link>
  )
}
