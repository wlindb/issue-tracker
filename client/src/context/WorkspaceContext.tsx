import { createContext, useCallback, useContext, useEffect, useState } from 'react'
import {
  listWorkspaces,
  createWorkspace as apiCreateWorkspace,
  type Workspace,
  type CreateWorkspaceRequest,
} from '@/api/generated/issueTrackerAPI'

interface WorkspaceContextValue {
  workspaces: Workspace[]
  activeWorkspace: Workspace | null
  setActiveWorkspace: (workspace: Workspace) => void
  createWorkspace: (request: CreateWorkspaceRequest) => Promise<Workspace>
  loading: boolean
}

const WorkspaceContext = createContext<WorkspaceContextValue | null>(null)

export function useWorkspace(): WorkspaceContextValue {
  const context = useContext(WorkspaceContext)
  if (context === null) {
    throw new Error('useWorkspace must be used within a WorkspaceProvider')
  }
  return context
}

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
