import { useState } from 'react'
import { ChevronDownIcon, ChevronRightIcon } from 'lucide-react'
import type { IssueGroup } from '@/lib/groupIssuesByStatus'
import { StatusIcon, STATUS_LABEL } from './StatusIcon'
import { IssueRow } from './IssueRow'

interface IssueGroupSectionProps {
  group: IssueGroup
}

export function IssueGroupSection({ group }: IssueGroupSectionProps) {
  const [expanded, setExpanded] = useState(true)

  return (
    <div>
      <button
        onClick={() => setExpanded((e) => !e)}
        className="flex w-full items-center gap-2 px-4 py-1.5 hover:bg-muted/40 transition-colors"
      >
        {expanded ? (
          <ChevronDownIcon className="size-3.5 text-muted-foreground" />
        ) : (
          <ChevronRightIcon className="size-3.5 text-muted-foreground" />
        )}
        <StatusIcon status={group.status} />
        <span className="text-xs font-medium uppercase tracking-wide text-muted-foreground">
          {STATUS_LABEL[group.status]}
        </span>
        <span className="ml-auto text-xs text-muted-foreground">{group.issues.length}</span>
      </button>
      {expanded && (
        <div>
          {group.issues.map((issue) => (
            <IssueRow key={issue.id} issue={issue} />
          ))}
        </div>
      )}
    </div>
  )
}
