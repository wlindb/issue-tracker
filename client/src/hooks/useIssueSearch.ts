import { useEffect, useState } from 'react'
import { searchIssues, type Issue } from '@/api/generated/issueTrackerAPI'
import { useDebounce } from './useDebounce'

const DEBOUNCE_MS = 300

interface UseIssueSearchResult {
  results: Issue[]
  isPending: boolean
  error: string | null
}

/**
 * Searches workspace issues by title via the backend, debounced.
 */
export function useIssueSearch(workspaceId: string, query: string): UseIssueSearchResult {
  const debouncedQuery = useDebounce(query, DEBOUNCE_MS)
  const isPending = query !== debouncedQuery
  const [results, setResults] = useState<Issue[]>([])
  const [error, setError] = useState<string | null>(null)
  const hasActiveQuery = Boolean(workspaceId && debouncedQuery)

  useEffect(() => {
    if (!workspaceId || !debouncedQuery) {
      return
    }

    let cancelled = false

    const search = async () => {
      try {
        const page = await searchIssues(workspaceId, { query: debouncedQuery })
        if (cancelled) {
          return
        }
        setResults(page.items)
        setError(null)
      } catch {
        if (cancelled) {
          return
        }
        setError('Failed to search issues.')
      }
    }

    search()

    return () => {
      cancelled = true
    }
  }, [workspaceId, debouncedQuery])

  return {
    results: hasActiveQuery ? results : [],
    isPending,
    error: hasActiveQuery ? error : null,
  }
}
