import { useMemo } from 'react'
import { type Issue } from '@/api/generated/issueTrackerAPI'
import { useDebounce } from './useDebounce'

const DEBOUNCE_MS = 300

interface UseIssueSearchResult {
  results: Issue[]
  isPending: boolean
}

/**
 * Searches issues by query string.
 *
 * Currently performs client-side filtering. To switch to a backend search,
 * replace the `useMemo` block with a data-fetching call (e.g. useQuery) using
 * `debouncedQuery` as the search parameter. The returned interface is unchanged.
 */
export function useIssueSearch(issues: Issue[], query: string): UseIssueSearchResult {
  const debouncedQuery = useDebounce(query, DEBOUNCE_MS)
  const isPending = query !== debouncedQuery

  const results = useMemo(() => {
    const normalized = debouncedQuery.trim().toLowerCase()
    if (!normalized) return issues

    return issues.filter(
      (issue) =>
        issue.title.toLowerCase().includes(normalized) ||
        issue.identifier.toLowerCase().includes(normalized),
    )
  }, [issues, debouncedQuery])

  return { results, isPending }
}
