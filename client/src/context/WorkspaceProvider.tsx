import { useCallback, useEffect, useState } from 'react'
import {
  listWorkspaces,
  createWorkspace as apiCreateWorkspace,
  type Workspace,
  type CreateWorkspaceRequest,
} from '@/api/generated/issueTrackerAPI'
import { WorkspaceContext } from './WorkspaceContext'

const ACTIVE_WORKSPACE_KEY = 'activeWorkspaceId'

interface WorkspaceProviderProps {
  children: React.ReactNode
}

export function WorkspaceProvider({ children }: WorkspaceProviderProps) {
  const [workspaces, setWorkspaces] = useState<Workspace[]>([])
  const [activeWorkspace, setActiveWorkspaceState] = useState<Workspace | null>(null)
  const [loading, setLoading] = useState(true)

  const setActiveWorkspace = useCallback((workspace: Workspace) => {
    setActiveWorkspaceState(workspace)
    localStorage.setItem(ACTIVE_WORKSPACE_KEY, workspace.id)
  }, [])

  useEffect(() => {
    async function load() {
      const page = await listWorkspaces()
      setWorkspaces(page.items)
      if (page.items.length > 0) {
        const storedId = localStorage.getItem(ACTIVE_WORKSPACE_KEY)
        const storedWorkspace = page.items.find((ws) => ws.id === storedId)
        setActiveWorkspaceState(storedWorkspace ?? page.items[0])
      }
      setLoading(false)
    }
    load()
  }, [])

  const createWorkspace = useCallback(async (request: CreateWorkspaceRequest) => {
    const workspace = await apiCreateWorkspace(request)
    setWorkspaces((prev) => [...prev, workspace])
    setActiveWorkspace(workspace)
    return workspace
  }, [setActiveWorkspace])

  return (
    <WorkspaceContext.Provider
      value={{ workspaces, activeWorkspace, setActiveWorkspace, createWorkspace, loading }}
    >
      {children}
    </WorkspaceContext.Provider>
  )
}
