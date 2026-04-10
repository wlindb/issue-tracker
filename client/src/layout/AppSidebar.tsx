import { useState } from 'react'
import { BuildingIcon, CheckIcon, ChevronsUpDownIcon, FolderOpenIcon, ListFilterIcon, MoonIcon, PlusCircleIcon, SunIcon, UserIcon } from 'lucide-react'
import { NavLink, useLocation, useNavigate } from 'react-router-dom'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarRail,
  SidebarSeparator,
  SidebarTrigger,
} from '@/components/ui/sidebar'
import { useTheme } from '@/context/ThemeContext'
import { useWorkspace } from '@/context/WorkspaceContext'
import { useKeycloak } from '@/auth/KeycloakProvider'

const navItems = [
  { to: '/my-issues', label: 'My Issues', icon: UserIcon },
  { to: '/all-issues', label: 'All Issues', icon: ListFilterIcon },
  { to: '/projects', label: 'Projects', icon: FolderOpenIcon },
]

function getInitials(name: string): string {
  return name
    .split(' ')
    .map((part) => part[0])
    .slice(0, 2)
    .join('')
    .toUpperCase()
}

export function AppSidebar() {
  const { theme, toggleTheme } = useTheme()
  const { keycloak } = useKeycloak()
  const { workspaces, activeWorkspace, setActiveWorkspace } = useWorkspace()
  const location = useLocation()
  const navigate = useNavigate()
  const username = (keycloak.tokenParsed?.preferred_username as string | undefined) ?? 'User'
  const [workspaceMenuOpen, setWorkspaceMenuOpen] = useState(false)

  return (
    <Sidebar collapsible="icon">
      <SidebarHeader>
        <div className="relative">
          <button
            type="button"
            onClick={() => setWorkspaceMenuOpen((prev) => !prev)}
            className="flex w-full items-center gap-2 rounded-md px-2 py-1.5 text-left hover:bg-muted group-data-[collapsible=icon]:justify-center"
          >
            <BuildingIcon className="size-4 shrink-0" />
            <span className="flex-1 truncate text-sm font-semibold tracking-tight group-data-[collapsible=icon]:hidden">
              {activeWorkspace?.name ?? 'Workspace'}
            </span>
            <ChevronsUpDownIcon className="size-3.5 shrink-0 text-muted-foreground group-data-[collapsible=icon]:hidden" />
          </button>
          <SidebarTrigger className="absolute right-0 top-1/2 -translate-y-1/2 group-data-[collapsible=icon]:relative group-data-[collapsible=icon]:right-auto group-data-[collapsible=icon]:top-auto group-data-[collapsible=icon]:translate-y-0" />

          {workspaceMenuOpen && (
            <>
              <div
                className="fixed inset-0 z-40"
                role="button"
                aria-label="Close menu"
                tabIndex={0}
                onClick={() => setWorkspaceMenuOpen(false)}
                onKeyDown={(e) => { if (e.key === 'Escape' || e.key === 'Enter') setWorkspaceMenuOpen(false) }}
              />
              <div className="absolute left-0 top-full z-50 mt-1 w-56 rounded-lg border border-border bg-popover p-1 shadow-md group-data-[collapsible=icon]:left-12">
                <div className="px-2 py-1.5 text-xs font-medium text-muted-foreground">
                  Workspaces
                </div>
                {workspaces.map((ws) => (
                  <button
                    key={ws.id}
                    type="button"
                    onClick={() => {
                      setActiveWorkspace(ws)
                      setWorkspaceMenuOpen(false)
                    }}
                    className="flex w-full items-center gap-2 rounded-md px-2 py-1.5 text-sm hover:bg-muted"
                  >
                    <BuildingIcon className="size-3.5 shrink-0" />
                    <span className="flex-1 truncate">{ws.name}</span>
                    {ws.id === activeWorkspace?.id && (
                      <CheckIcon className="size-3.5 shrink-0 text-primary" />
                    )}
                  </button>
                ))}
                <div className="my-1 h-px bg-border" />
                <button
                  type="button"
                  onClick={() => {
                    setWorkspaceMenuOpen(false)
                    navigate('/create-workspace')
                  }}
                  className="flex w-full items-center gap-2 rounded-md px-2 py-1.5 text-sm hover:bg-muted"
                >
                  <PlusCircleIcon className="size-3.5 shrink-0" />
                  <span>Create workspace</span>
                </button>
              </div>
            </>
          )}
        </div>
      </SidebarHeader>

      <SidebarSeparator />

      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupLabel>Navigation</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              {navItems.map(({ to, label, icon: Icon }) => {
                const isActive =
                  to === '/projects'
                    ? location.pathname.startsWith('/projects')
                    : location.pathname === to
                return (
                  <SidebarMenuItem key={to}>
                    <SidebarMenuButton
                      isActive={isActive}
                      tooltip={label}
                      render={<NavLink to={to} />}
                    >
                      <Icon />
                      <span>{label}</span>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                )
              })}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>

      <SidebarSeparator />

      <SidebarFooter>
        <div className="flex items-center gap-2 px-2 py-1 group-data-[collapsible=icon]:flex-col group-data-[collapsible=icon]:items-center group-data-[collapsible=icon]:gap-1">
          <Avatar size="sm">
            <AvatarFallback>{getInitials(username)}</AvatarFallback>
          </Avatar>
          <span className="flex-1 truncate text-xs text-sidebar-foreground group-data-[collapsible=icon]:hidden">
            {username}
          </span>
          <Button
            variant="ghost"
            size="icon-sm"
            onClick={toggleTheme}
            aria-label="Toggle theme"
          >
            {theme === 'dark' ? <SunIcon /> : <MoonIcon />}
          </Button>
        </div>
        <div className="px-2 group-data-[collapsible=icon]:hidden">
          <Button
            variant="outline"
            size="sm"
            className="w-full"
            onClick={() => keycloak.logout()}
          >
            Sign out
          </Button>
        </div>
      </SidebarFooter>

      <SidebarRail />
    </Sidebar>
  )
}
