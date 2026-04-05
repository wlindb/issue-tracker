import { keycloak } from '@/keycloak'

const apiUrl = import.meta.env.VITE_API_URL ?? "https://issuetrackerapi.astraterial.com"

if (!apiUrl) {
  throw new Error('VITE_API_URL is not set')
}

export interface RequestConfig {
  url: string
  method: string
  params?: Record<string, unknown>
  data?: unknown
  headers?: Record<string, string>
  signal?: AbortSignal
}

export class ApiError extends Error {
  readonly status: number
  readonly statusText: string
  readonly url: string

  constructor(status: number, statusText: string, url: string) {
    super(`API error ${status} ${statusText}: ${url}`)
    this.name = 'ApiError'
    this.status = status
    this.statusText = statusText
    this.url = url
  }
}

export const customFetch = async <T>(config: RequestConfig): Promise<T> => {
  const token = keycloak.token

  const headers = new Headers(config.headers)
  headers.set('Content-Type', 'application/json')
  if (token !== undefined) {
    headers.set('Authorization', `Bearer ${token}`)
  }

  const url = new URL(`${apiUrl}${config.url}`)
  if (config.params) {
    for (const [key, value] of Object.entries(config.params)) {
      if (value !== undefined) {
        url.searchParams.set(key, String(value))
      }
    }
  }

  const response = await fetch(url.toString(), {
    method: config.method,
    headers,
    body: config.data !== undefined ? JSON.stringify(config.data) : undefined,
    signal: config.signal,
  })

  if (!response.ok) {
    throw new ApiError(response.status, response.statusText, url.toString())
  }

  if (response.status === 204) {
    return undefined as T
  }

  return response.json() as Promise<T>
}
