import { Navigate, Route, Routes } from 'react-router-dom'
import { AppLayout } from './layout/AppLayout'
import { AllIssuesPage } from './pages/AllIssuesPage'
import { IssueDetailPage } from './pages/IssueDetailPage'
import { MyIssuesPage } from './pages/MyIssuesPage'
import { ProjectDetailPage } from './pages/ProjectDetailPage'
import { ProjectsPage } from './pages/ProjectsPage'

function App() {
  return (
    <Routes>
      <Route element={<AppLayout />}>
        <Route index element={<Navigate to="/my-issues" replace />} />
        <Route path="/my-issues" element={<MyIssuesPage />} />
        <Route path="/all-issues" element={<AllIssuesPage />} />
        <Route path="/projects" element={<ProjectsPage />} />
        <Route path="/projects/:identifier" element={<ProjectDetailPage />} />
        <Route path="/issues/:identifier" element={<IssueDetailPage />} />
      </Route>
    </Routes>
  )
}

export default App
