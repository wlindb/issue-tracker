import { XIcon } from 'lucide-react'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import type { Issue, IssuePriority, IssueStatus, Label, WorkspaceMember } from '@/api/generated/issueTrackerAPI'
import { PriorityIcon } from '@/components/PriorityIcon'
import { StatusIcon } from '@/components/StatusIcon'
import { STATUS_LABEL } from '@/components/statusLabel'
import { LabelDot } from '@/components/LabelDot'
import { LabelPicker } from '@/components/issue-detail/LabelPicker'
import { cn } from '@/lib/utils'

const PRIORITY_LABEL: Record<IssuePriority, string> = {
  urgent: 'Urgent',
  high: 'High',
  medium: 'Medium',
  low: 'Low',
  none: 'No priority',
}

const ALL_STATUSES: IssueStatus[] = ['backlog', 'todo', 'in_progress', 'done', 'cancelled']
const ALL_PRIORITIES: IssuePriority[] = ['none', 'urgent', 'high', 'medium', 'low']

function getInitials(email: string): string {
  return email.split('@')[0].slice(0, 2).toUpperCase()
}

const selectClass = cn(
  'w-full rounded-md border border-input bg-transparent px-2.5 py-1 text-sm',
  'focus:outline-none focus:ring-2 focus:ring-ring cursor-pointer',
)

interface IssueMetaSidebarProps {
  issue: Issue
  members: WorkspaceMember[]
  workspaceId: string
  onStatusChange: (status: IssueStatus) => void
  onPriorityChange: (priority: IssuePriority) => void
  onAssigneeChange: (assigneeId: string | null) => void
  onLabelsChange: (labels: Label[]) => void
}

export function IssueMetaSidebar({
  issue,
  members,
  workspaceId,
  onStatusChange,
  onPriorityChange,
  onAssigneeChange,
  onLabelsChange,
}: IssueMetaSidebarProps) {
  function handleRemoveLabel(labelId: string) {
    onLabelsChange(issue.labels.filter((l) => l.id !== labelId))
  }

  return (
    <div className="flex flex-col gap-5">
      <MetaRow label="Status">
        <div className="relative">
          <select
            value={issue.status}
            onChange={(e) => onStatusChange(e.target.value as IssueStatus)}
            className={selectClass}
          >
            {ALL_STATUSES.map((s) => (
              <option key={s} value={s}>
                {STATUS_LABEL[s]}
              </option>
            ))}
          </select>
        </div>
        {/* Visual preview below the select */}
        <div className="flex items-center gap-1.5 px-0.5">
          <StatusIcon status={issue.status} />
          <span className="text-xs text-muted-foreground">{STATUS_LABEL[issue.status]}</span>
        </div>
      </MetaRow>

      <MetaRow label="Priority">
        <select
          value={issue.priority}
          onChange={(e) => onPriorityChange(e.target.value as IssuePriority)}
          className={selectClass}
        >
          {ALL_PRIORITIES.map((p) => (
            <option key={p} value={p}>
              {PRIORITY_LABEL[p]}
            </option>
          ))}
        </select>
        <div className="flex items-center gap-1.5 px-0.5">
          <PriorityIcon priority={issue.priority} className="size-3.5" />
          <span className="text-xs text-muted-foreground">{PRIORITY_LABEL[issue.priority]}</span>
        </div>
      </MetaRow>

      <MetaRow label="Assignee">
        <select
          value={issue.assigneeId ?? ''}
          onChange={(e) => {
            const userId = e.target.value
            onAssigneeChange(userId || null)
          }}
          className={selectClass}
        >
          <option value="">Unassigned</option>
          {members.map((m) => (
            <option key={m.id} value={m.id}>
              {m.email}
            </option>
          ))}
        </select>
        {issue.assigneeId && (() => {
          const assignee = members.find((m) => m.id === issue.assigneeId)
          return assignee ? (
            <div className="flex items-center gap-2 px-0.5">
              <Avatar size="sm">
                <AvatarFallback>{getInitials(assignee.email)}</AvatarFallback>
              </Avatar>
              <span className="text-xs text-muted-foreground">{assignee.email}</span>
            </div>
          ) : null
        })()}
      </MetaRow>

      <MetaRow label="Labels">
        <div className="flex flex-wrap items-center gap-1">
          {issue.labels.map((label) => (
            <Badge key={label.id} variant="outline" className="gap-1 pr-1 text-xs">
              <LabelDot labelId={label.id} />
              {label.name}
              <button
                type="button"
                onClick={() => handleRemoveLabel(label.id)}
                className="rounded hover:text-destructive transition-colors"
              >
                <XIcon className="size-3" />
              </button>
            </Badge>
          ))}
          <LabelPicker
            workspaceId={workspaceId}
            attachedLabels={issue.labels}
            onLabelsChange={onLabelsChange}
          />
        </div>
      </MetaRow>
    </div>
  )
}

function MetaRow({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <div className="flex flex-col gap-1.5">
      <span className="text-xs font-medium uppercase tracking-wide text-muted-foreground">
        {label}
      </span>
      {children}
    </div>
  )
}
