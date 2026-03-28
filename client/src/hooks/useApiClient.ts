import { useMemo } from 'react'
import { customFetch } from '@/api/mutator'

export function useApiClient() {
  return useMemo(
    () => ({
      get: <T>(url: string, params?: Record<string, unknown>) =>
        customFetch<T>({ url, method: 'GET', params }),
      post: <T>(url: string, data: unknown) =>
        customFetch<T>({ url, method: 'POST', data }),
      patch: <T>(url: string, data: unknown) =>
        customFetch<T>({ url, method: 'PATCH', data }),
      delete: <T>(url: string) =>
        customFetch<T>({ url, method: 'DELETE' }),
    }),
    [], // customFetch is a stable module-level function
  )
}
