import { useCallback, useEffect, useState } from 'react'
import {
  listWorkspaces,
  createWorkspace as apiCreateWorkspace,
  type Workspace,
  type CreateWorkspaceRequest,
} from '@/api/generated/issueTrackerAPI'
import { WorkspaceContext } from './WorkspaceContext'

interface WorkspaceProviderProps {
  children: React.ReactNode
}

export function WorkspaceProvider({ children }: WorkspaceProviderProps) {
  const [workspaces, setWorkspaces] = useState<Workspace[]>([])
  const [activeWorkspace, setActiveWorkspace] = useState<Workspace | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    async function load() {
      const page = await listWorkspaces()
      setWorkspaces(page.items)
      if (page.items.length > 0) {
        setActiveWorkspace(page.items[0])
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
  }, [])

  return (
    <WorkspaceContext.Provider
      value={{ workspaces, activeWorkspace, setActiveWorkspace, createWorkspace, loading }}
    >
      {children}
    </WorkspaceContext.Provider>
  )
}
