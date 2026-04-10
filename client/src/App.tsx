import { Navigate, Route, Routes } from 'react-router-dom'
import { AppLayout } from './layout/AppLayout'
import { AllIssuesPage } from './pages/AllIssuesPage'
import { CreateWorkspacePage } from './pages/CreateWorkspacePage'
import { IssueDetailPage } from './pages/IssueDetailPage'
import { MyIssuesPage } from './pages/MyIssuesPage'
import { ProjectDetailPage } from './pages/ProjectDetailPage'
import { ProjectsPage } from './pages/ProjectsPage'
import { WorkspaceProvider } from './context/WorkspaceContext'
import { WorkspaceGuard } from './layout/WorkspaceGuard'

function App() {
  return (
    <WorkspaceProvider>
      <Routes>
        <Route path="/create-workspace" element={<CreateWorkspacePage />} />
        <Route element={<WorkspaceGuard />}>
          <Route element={<AppLayout />}>
            <Route index element={<Navigate to="/my-issues" replace />} />
            <Route path="/my-issues" element={<MyIssuesPage />} />
            <Route path="/all-issues" element={<AllIssuesPage />} />
            <Route path="/projects" element={<ProjectsPage />} />
            <Route path="/projects/:identifier" element={<ProjectDetailPage />} />
            <Route path="/issues/:identifier" element={<IssueDetailPage />} />
          </Route>
        </Route>
      </Routes>
    </WorkspaceProvider>
  )
}

export default App
