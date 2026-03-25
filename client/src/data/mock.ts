export type Status = 'backlog' | 'todo' | 'in_progress' | 'in_review' | 'done' | 'cancelled'
export type Priority = 'no_priority' | 'urgent' | 'high' | 'medium' | 'low'

export interface Project {
  id: string
  name: string
  identifier: string
  description: string
  issueCount: number
  color: string
}

export interface Issue {
  id: string
  identifier: string
  title: string
  status: Status
  priority: Priority
  projectId: string
  assigneeId: string | null
  assigneeName: string | null
  createdAt: string
}

export const CURRENT_USER_ID = 'user-1'

export const projects: Project[] = [
  {
    id: 'proj-1',
    name: 'Platform',
    identifier: 'PLAT',
    description: 'Core infrastructure, APIs, and backend services powering the product.',
    issueCount: 5,
    color: '#6366f1',
  },
  {
    id: 'proj-2',
    name: 'Mobile App',
    identifier: 'MOB',
    description: 'iOS and Android applications for end users.',
    issueCount: 4,
    color: '#f59e0b',
  },
  {
    id: 'proj-3',
    name: 'Design System',
    identifier: 'DS',
    description: 'Shared component library, tokens, and design guidelines.',
    issueCount: 3,
    color: '#10b981',
  },
]

export const issues: Issue[] = [
  // Platform issues
  {
    id: 'issue-1',
    identifier: 'PLAT-1',
    title: 'Migrate authentication to OAuth 2.0',
    status: 'in_progress',
    priority: 'high',
    projectId: 'proj-1',
    assigneeId: 'user-1',
    assigneeName: 'Alice',
    createdAt: '2026-03-20T10:00:00Z',
  },
  {
    id: 'issue-2',
    identifier: 'PLAT-2',
    title: 'Add rate limiting to public API endpoints',
    status: 'todo',
    priority: 'medium',
    projectId: 'proj-1',
    assigneeId: 'user-1',
    assigneeName: 'Alice',
    createdAt: '2026-03-21T09:00:00Z',
  },
  {
    id: 'issue-3',
    identifier: 'PLAT-3',
    title: 'Database connection pool exhaustion under load',
    status: 'in_review',
    priority: 'urgent',
    projectId: 'proj-1',
    assigneeId: 'user-2',
    assigneeName: 'Bob',
    createdAt: '2026-03-19T14:00:00Z',
  },
  {
    id: 'issue-4',
    identifier: 'PLAT-4',
    title: 'Set up structured logging with trace IDs',
    status: 'backlog',
    priority: 'low',
    projectId: 'proj-1',
    assigneeId: null,
    assigneeName: null,
    createdAt: '2026-03-18T11:00:00Z',
  },
  {
    id: 'issue-5',
    identifier: 'PLAT-5',
    title: 'Deploy staging environment on Railway',
    status: 'done',
    priority: 'medium',
    projectId: 'proj-1',
    assigneeId: 'user-2',
    assigneeName: 'Bob',
    createdAt: '2026-03-15T08:00:00Z',
  },
  // Mobile App issues
  {
    id: 'issue-6',
    identifier: 'MOB-1',
    title: 'Implement offline mode with local cache',
    status: 'in_progress',
    priority: 'high',
    projectId: 'proj-2',
    assigneeId: 'user-1',
    assigneeName: 'Alice',
    createdAt: '2026-03-22T10:00:00Z',
  },
  {
    id: 'issue-7',
    identifier: 'MOB-2',
    title: 'Fix crash on iOS 17 when opening notifications',
    status: 'in_review',
    priority: 'urgent',
    projectId: 'proj-2',
    assigneeId: 'user-3',
    assigneeName: 'Carol',
    createdAt: '2026-03-23T09:30:00Z',
  },
  {
    id: 'issue-8',
    identifier: 'MOB-3',
    title: 'Add biometric authentication support',
    status: 'backlog',
    priority: 'medium',
    projectId: 'proj-2',
    assigneeId: null,
    assigneeName: null,
    createdAt: '2026-03-17T13:00:00Z',
  },
  {
    id: 'issue-9',
    identifier: 'MOB-4',
    title: 'Dark mode support for all screens',
    status: 'cancelled',
    priority: 'low',
    projectId: 'proj-2',
    assigneeId: 'user-3',
    assigneeName: 'Carol',
    createdAt: '2026-03-10T10:00:00Z',
  },
  // Design System issues
  {
    id: 'issue-10',
    identifier: 'DS-1',
    title: 'Create token documentation site',
    status: 'todo',
    priority: 'medium',
    projectId: 'proj-3',
    assigneeId: 'user-1',
    assigneeName: 'Alice',
    createdAt: '2026-03-24T08:00:00Z',
  },
  {
    id: 'issue-11',
    identifier: 'DS-2',
    title: 'Audit and consolidate spacing tokens',
    status: 'in_progress',
    priority: 'high',
    projectId: 'proj-3',
    assigneeId: 'user-3',
    assigneeName: 'Carol',
    createdAt: '2026-03-22T14:00:00Z',
  },
  {
    id: 'issue-12',
    identifier: 'DS-3',
    title: 'Publish component library to npm',
    status: 'done',
    priority: 'no_priority',
    projectId: 'proj-3',
    assigneeId: 'user-2',
    assigneeName: 'Bob',
    createdAt: '2026-03-12T10:00:00Z',
  },
]
