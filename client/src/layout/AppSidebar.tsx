import { FolderOpenIcon, ListFilterIcon, MoonIcon, SunIcon, UserIcon } from 'lucide-react'
import { NavLink, useLocation } from 'react-router-dom'
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
  const location = useLocation()
  const username = (keycloak.tokenParsed?.preferred_username as string | undefined) ?? 'User'

  return (
    <Sidebar collapsible="icon">
      <SidebarHeader>
        <div className="flex items-center justify-between px-2 py-1">
          <span className="text-sm font-semibold tracking-tight group-data-[collapsible=icon]:hidden">
            IssueTracker
          </span>
          <SidebarTrigger className="-mr-1" />
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
