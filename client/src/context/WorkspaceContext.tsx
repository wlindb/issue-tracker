import { createContext, useContext } from 'react'
import {
  type Workspace,
  type CreateWorkspaceRequest,
} from '@/api/generated/issueTrackerAPI'

export interface WorkspaceContextValue {
  workspaces: Workspace[]
  activeWorkspace: Workspace | null
  setActiveWorkspace: (workspace: Workspace) => void
  createWorkspace: (request: CreateWorkspaceRequest) => Promise<Workspace>
  loading: boolean
}

export const WorkspaceContext = createContext<WorkspaceContextValue | null>(null)

export function useWorkspace(): WorkspaceContextValue {
  const context = useContext(WorkspaceContext)
  if (context === null) {
    throw new Error('useWorkspace must be used within a WorkspaceProvider')
  }
  return context
}
