import { useEffect, useState } from 'react'
import { listLabels, type Label } from '@/api/generated/issueTrackerAPI'
import { useDebounce } from './useDebounce'

const DEBOUNCE_MS = 300

interface UseLabelSearchResult {
  results: Label[]
  isPending: boolean
}

/**
 * Searches workspace labels by name via the backend, debounced.
 */
export function useLabelSearch(workspaceId: string, query: string): UseLabelSearchResult {
  const debouncedQuery = useDebounce(query, DEBOUNCE_MS)
  const isPending = query !== debouncedQuery
  const [results, setResults] = useState<Label[]>([])

  useEffect(() => {
    let cancelled = false

    listLabels(workspaceId, { search: debouncedQuery || undefined }).then((page) => {
      if (!cancelled) setResults(page.items)
    })

    return () => {
      cancelled = true
    }
  }, [workspaceId, debouncedQuery])

  return { results, isPending }
}
