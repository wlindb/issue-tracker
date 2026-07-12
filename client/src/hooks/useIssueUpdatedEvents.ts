import { type Issue } from '@/api/generated/issueTrackerAPI'
import { useIssueStatusUpdatedEvents } from '@/hooks/useIssueStatusUpdatedEvents'
import { useIssueTitleUpdatedEvents } from '@/hooks/useIssueTitleUpdatedEvents'
import { useIssuePriorityUpdatedEvents } from '@/hooks/useIssuePriorityUpdatedEvents'
import { useIssueAssigneeUpdatedEvents } from '@/hooks/useIssueAssigneeUpdatedEvents'
import { useIssueDescriptionUpdatedEvents } from '@/hooks/useIssueDescriptionUpdatedEvents'

interface IssueUpdatedEvent {
  payload: Issue
}

export function useIssueUpdatedEvents(issueId: string | undefined, onIssueUpdated: (issue: Issue) => void) {
  const handleIssueUpdated = (event: IssueUpdatedEvent) => {
    if (!issueId || event.payload.id !== issueId) return
    onIssueUpdated(event.payload)
  }

  useIssueStatusUpdatedEvents(handleIssueUpdated)
  useIssueTitleUpdatedEvents(handleIssueUpdated)
  useIssuePriorityUpdatedEvents(handleIssueUpdated)
  useIssueAssigneeUpdatedEvents(handleIssueUpdated)
  useIssueDescriptionUpdatedEvents(handleIssueUpdated)
}
