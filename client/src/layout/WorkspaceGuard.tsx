import { Navigate, Outlet } from 'react-router-dom'
import { useWorkspace } from '@/context/WorkspaceContext'

export function WorkspaceGuard() {
  const { activeWorkspace, loading } = useWorkspace()

  if (loading) {
    return <div className="flex min-h-screen items-center justify-center">Loading…</div>
  }

  if (!activeWorkspace) {
    return <Navigate to="/create-workspace" replace />
  }

  return <Outlet />
}
