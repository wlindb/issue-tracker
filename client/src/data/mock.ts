export type Status = 'backlog' | 'todo' | 'in_progress' | 'done' | 'cancelled'
export type Priority = 'none' | 'urgent' | 'high' | 'medium' | 'low'

export interface Project {
  id: string
  identifier: string
  name: string
  description?: string | null
  ownerId: string
  issueCount: number
  createdAt: string
  updatedAt: string
}

export interface Issue {
  id: string
  identifier: string
  projectId: string
  title: string
  description?: string | null
  status: Status
  priority: Priority
  labels: string[]
  assigneeId?: string | null
  reporterId: string
  createdAt: string
  updatedAt: string
}

export interface Comment {
  id: string
  issueId: string
  authorId: string
  authorName: string
  body: string
  createdAt: string
}

export const CURRENT_USER_ID = 'user-1'

export interface User {
  id: string
  name: string
}

export const users: User[] = [
  { id: 'user-1', name: 'Alice' },
  { id: 'user-2', name: 'Bob' },
  { id: 'user-3', name: 'Carol' },
]

export const projects: Project[] = [
  {
    id: 'proj-1',
    name: 'Platform',
    identifier: 'PLAT',
    description: 'Core infrastructure, APIs, and backend services powering the product.',
    ownerId: 'user-1',
    issueCount: 5,
    createdAt: '2026-03-01T00:00:00Z',
    updatedAt: '2026-03-01T00:00:00Z',
  },
  {
    id: 'proj-2',
    name: 'Mobile App',
    identifier: 'MOB',
    description: 'iOS and Android applications for end users.',
    ownerId: 'user-1',
    issueCount: 4,
    createdAt: '2026-03-01T00:00:00Z',
    updatedAt: '2026-03-01T00:00:00Z',
  },
  {
    id: 'proj-3',
    name: 'Design System',
    identifier: 'DS',
    description: 'Shared component library, tokens, and design guidelines.',
    ownerId: 'user-1',
    issueCount: 3,
    createdAt: '2026-03-01T00:00:00Z',
    updatedAt: '2026-03-01T00:00:00Z',
  },
]

export const issues: Issue[] = [
  // Platform issues
  {
    id: 'issue-1',
    identifier: 'PLAT-1',
    projectId: 'proj-1',
    title: 'Migrate authentication to OAuth 2.0',
    description:
      'Replace the current session-based auth with OAuth 2.0 + PKCE flow. This will allow third-party integrations and improve security posture. Target: Auth0 as the provider, with a migration script to move existing sessions.',
    labels: ['auth', 'security', 'backend'],
    status: 'in_progress',
    priority: 'high',
    assigneeId: 'user-1',
    reporterId: 'user-1',
    createdAt: '2026-03-20T10:00:00Z',
    updatedAt: '2026-03-20T10:00:00Z',
  },
  {
    id: 'issue-2',
    identifier: 'PLAT-2',
    projectId: 'proj-1',
    title: 'Add rate limiting to public API endpoints',
    description: null,
    labels: ['api', 'backend'],
    status: 'todo',
    priority: 'medium',
    assigneeId: 'user-1',
    reporterId: 'user-1',
    createdAt: '2026-03-21T09:00:00Z',
    updatedAt: '2026-03-21T09:00:00Z',
  },
  {
    id: 'issue-3',
    identifier: 'PLAT-3',
    projectId: 'proj-1',
    title: 'Database connection pool exhaustion under load',
    description:
      'Under sustained load (>500 rps) the pg connection pool hits its limit and requests start queuing, causing p99 latency to spike above 2s. Need to investigate pool sizing, query duration, and whether we can introduce read replicas.',
    labels: ['database', 'performance', 'urgent'],
    status: 'in_progress',
    priority: 'urgent',
    assigneeId: 'user-2',
    reporterId: 'user-1',
    createdAt: '2026-03-19T14:00:00Z',
    updatedAt: '2026-03-19T14:00:00Z',
  },
  {
    id: 'issue-4',
    identifier: 'PLAT-4',
    projectId: 'proj-1',
    title: 'Set up structured logging with trace IDs',
    description: null,
    labels: ['observability'],
    status: 'backlog',
    priority: 'low',
    assigneeId: null,
    reporterId: 'user-1',
    createdAt: '2026-03-18T11:00:00Z',
    updatedAt: '2026-03-18T11:00:00Z',
  },
  {
    id: 'issue-5',
    identifier: 'PLAT-5',
    projectId: 'proj-1',
    title: 'Deploy staging environment on Railway',
    description: null,
    labels: ['devops'],
    status: 'done',
    priority: 'medium',
    assigneeId: 'user-2',
    reporterId: 'user-1',
    createdAt: '2026-03-15T08:00:00Z',
    updatedAt: '2026-03-15T08:00:00Z',
  },
  // Mobile App issues
  {
    id: 'issue-6',
    identifier: 'MOB-1',
    projectId: 'proj-2',
    title: 'Implement offline mode with local cache',
    description: null,
    labels: ['offline', 'mobile'],
    status: 'in_progress',
    priority: 'high',
    assigneeId: 'user-1',
    reporterId: 'user-1',
    createdAt: '2026-03-22T10:00:00Z',
    updatedAt: '2026-03-22T10:00:00Z',
  },
  {
    id: 'issue-7',
    identifier: 'MOB-2',
    projectId: 'proj-2',
    title: 'Fix crash on iOS 17 when opening notifications',
    description: null,
    labels: ['bug', 'ios'],
    status: 'in_progress',
    priority: 'urgent',
    assigneeId: 'user-3',
    reporterId: 'user-1',
    createdAt: '2026-03-23T09:30:00Z',
    updatedAt: '2026-03-23T09:30:00Z',
  },
  {
    id: 'issue-8',
    identifier: 'MOB-3',
    projectId: 'proj-2',
    title: 'Add biometric authentication support',
    description: null,
    labels: ['auth', 'mobile'],
    status: 'backlog',
    priority: 'medium',
    assigneeId: null,
    reporterId: 'user-1',
    createdAt: '2026-03-17T13:00:00Z',
    updatedAt: '2026-03-17T13:00:00Z',
  },
  {
    id: 'issue-9',
    identifier: 'MOB-4',
    projectId: 'proj-2',
    title: 'Dark mode support for all screens',
    description: null,
    labels: [],
    status: 'cancelled',
    priority: 'low',
    assigneeId: 'user-3',
    reporterId: 'user-1',
    createdAt: '2026-03-10T10:00:00Z',
    updatedAt: '2026-03-10T10:00:00Z',
  },
  // Design System issues
  {
    id: 'issue-10',
    identifier: 'DS-1',
    projectId: 'proj-3',
    title: 'Create token documentation site',
    description: null,
    labels: ['docs', 'design'],
    status: 'todo',
    priority: 'medium',
    assigneeId: 'user-1',
    reporterId: 'user-1',
    createdAt: '2026-03-24T08:00:00Z',
    updatedAt: '2026-03-24T08:00:00Z',
  },
  {
    id: 'issue-11',
    identifier: 'DS-2',
    projectId: 'proj-3',
    title: 'Audit and consolidate spacing tokens',
    description: null,
    labels: ['design', 'tokens'],
    status: 'in_progress',
    priority: 'high',
    assigneeId: 'user-3',
    reporterId: 'user-1',
    createdAt: '2026-03-22T14:00:00Z',
    updatedAt: '2026-03-22T14:00:00Z',
  },
  {
    id: 'issue-12',
    identifier: 'DS-3',
    projectId: 'proj-3',
    title: 'Publish component library to npm',
    description: null,
    labels: [],
    status: 'done',
    priority: 'none',
    assigneeId: 'user-2',
    reporterId: 'user-1',
    createdAt: '2026-03-12T10:00:00Z',
    updatedAt: '2026-03-12T10:00:00Z',
  },
]

export const comments: Comment[] = [
  {
    id: 'comment-1',
    issueId: 'issue-1',
    authorId: 'user-2',
    authorName: 'Bob',
    body: 'Looked into Auth0 pricing — the Developer plan covers our current user count. Should we add the migration plan to the PR description?',
    createdAt: '2026-03-21T11:00:00Z',
  },
  {
    id: 'comment-2',
    issueId: 'issue-1',
    authorId: 'user-1',
    authorName: 'Alice',
    body: 'Yes, I will add it. Also planning to keep the old session endpoint alive for 2 weeks after rollout for a smooth transition.',
    createdAt: '2026-03-21T11:45:00Z',
  },
  {
    id: 'comment-3',
    issueId: 'issue-3',
    authorId: 'user-3',
    authorName: 'Carol',
    body: 'Reproduced locally with k6 at 600 rps. Pool size is currently 10 — bumping to 25 reduced the p99 significantly in my tests. Sharing results in the PR.',
    createdAt: '2026-03-20T16:30:00Z',
  },
]
