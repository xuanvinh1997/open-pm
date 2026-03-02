<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useWorkspaceStore } from '@/stores/workspace.store'
import { useProjectStore } from '@/stores/project.store'
import { useAuthStore } from '@/stores/auth.store'
import { useSidebar } from '@/composables/useSidebar'
import { useTheme } from '@/composables/useTheme'
import { useCommandPalette } from '@/composables/useCommandPalette'
import PAvatar from '@/components/ui/PAvatar.vue'
import PTooltip from '@/components/ui/PTooltip.vue'
import PDropdown from '@/components/ui/PDropdown.vue'
import NotificationPopover from '@/components/notifications/NotificationPopover.vue'
import {
  Home,
  Briefcase,
  Settings,
  LogOut,
  Plus,
  ChevronDown,
  ChevronRight,
  ChevronsLeft,
  ChevronsRight,
  LayoutList,
  Kanban,
  BarChart3,
  FileBarChart,
  Search,
  Sun,
  Moon,
  Users,
  Repeat,
  Layers,
  FileText,
} from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const workspaceStore = useWorkspaceStore()
const projectStore = useProjectStore()
const authStore = useAuthStore()
const { isCollapsed, toggle } = useSidebar()
const { theme, toggleTheme } = useTheme()
const { open: openCommandPalette } = useCommandPalette()

const slug = computed(() => route.params.workspaceSlug as string)
const expandedProjects = ref<Set<string>>(new Set())

const navItems = computed(() => {
  const items = [
    { name: 'Home', icon: Home, to: `/${slug.value}` },
    { name: 'Projects', icon: Briefcase, to: `/${slug.value}/projects` },
  ]
  if (workspaceStore.isAdmin) {
    items.push({ name: 'Members', icon: Users, to: `/${slug.value}/members` })
    items.push({ name: 'Settings', icon: Settings, to: `/${slug.value}/settings` })
  }
  return items
})

function toggleProject(projectId: string) {
  if (expandedProjects.value.has(projectId)) {
    expandedProjects.value.delete(projectId)
  } else {
    expandedProjects.value.add(projectId)
  }
}

function isActive(to: string): boolean {
  return route.path === to
}

function handleLogout() {
  authStore.logout()
  router.push('/auth/login')
}

function switchWorkspace(wsSlug: string) {
  router.push(`/${wsSlug}`)
}

const userName = computed(() => {
  const user = authStore.user
  if (!user) return ''
  if (user.display_name) return user.display_name
  if (user.first_name) return `${user.first_name} ${user.last_name || ''}`.trim()
  return user.email || ''
})
</script>

<template>
  <div class="flex h-full flex-col bg-custom-background-80 border-r border-custom-border-200">
    <!-- Workspace selector -->
    <div class="flex h-12 items-center border-b border-custom-border-200 px-3">
      <template v-if="!isCollapsed">
        <PDropdown align="left" width="14rem">
          <template #trigger>
            <div class="flex items-center gap-2 rounded-md px-2 py-1.5 hover:bg-custom-background-90 transition-colors cursor-pointer flex-1 min-w-0">
              <div class="flex h-7 w-7 items-center justify-center rounded-md bg-brand-600 text-xs font-bold text-white flex-shrink-0">
                {{ workspaceStore.currentWorkspace?.name?.charAt(0)?.toUpperCase() || 'O' }}
              </div>
              <span class="truncate text-sm font-semibold text-custom-text-100">
                {{ workspaceStore.currentWorkspace?.name || 'Open-PM' }}
              </span>
              <ChevronDown class="h-3.5 w-3.5 text-custom-text-300 flex-shrink-0" />
            </div>
          </template>
          <template #default="{ close }">
            <div class="px-2 py-1.5 text-xs font-semibold text-custom-text-300 uppercase tracking-wider">
              Workspaces
            </div>
            <button
              v-for="ws in workspaceStore.workspaces"
              :key="ws.id"
              class="flex w-full items-center gap-2 rounded-md px-2 py-2 text-sm text-custom-text-200 hover:bg-custom-background-80 transition-colors"
              :class="ws.slug === slug ? 'bg-custom-background-80 text-custom-text-100 font-medium' : ''"
              @click="switchWorkspace(ws.slug); close()"
            >
              <div class="flex h-6 w-6 items-center justify-center rounded bg-brand-600 text-2xs font-bold text-white flex-shrink-0">
                {{ ws.name.charAt(0).toUpperCase() }}
              </div>
              <span class="truncate">{{ ws.name }}</span>
            </button>
          </template>
        </PDropdown>
      </template>
      <template v-else>
        <PTooltip content="Switch workspace" position="right">
          <div class="mx-auto flex h-8 w-8 items-center justify-center rounded-md bg-brand-600 text-xs font-bold text-white cursor-pointer">
            {{ workspaceStore.currentWorkspace?.name?.charAt(0)?.toUpperCase() || 'O' }}
          </div>
        </PTooltip>
      </template>
    </div>

    <!-- Search trigger -->
    <div class="px-3 pt-3 pb-1" v-if="!isCollapsed">
      <button
        class="flex w-full items-center gap-2.5 rounded-md border border-custom-border-200 bg-custom-background-100 px-2.5 py-1.5 text-sm text-custom-text-300 hover:bg-custom-background-90 transition-colors"
        @click="openCommandPalette"
      >
        <Search class="h-3.5 w-3.5" />
        <span class="flex-1 text-left">Search...</span>
        <kbd class="rounded border border-custom-border-200 px-1 py-0.5 text-2xs font-mono text-custom-text-300">⌘K</kbd>
      </button>
    </div>

    <!-- Navigation -->
    <nav class="flex-1 space-y-0.5 overflow-y-auto px-3 py-3">
      <template v-if="!isCollapsed">
        <router-link
          v-for="item in navItems"
          :key="item.name"
          :to="item.to"
          class="sidebar-nav-item"
          :class="{ active: isActive(item.to) }"
        >
          <component :is="item.icon" class="h-4 w-4 flex-shrink-0" />
          <span>{{ item.name }}</span>
        </router-link>
      </template>
      <template v-else>
        <PTooltip v-for="item in navItems" :key="item.name" :content="item.name" position="right">
          <router-link
            :to="item.to"
            class="flex items-center justify-center rounded-md p-2 text-custom-text-200 hover:bg-custom-background-90 transition-colors"
            :class="{ 'bg-custom-background-90 text-custom-text-100': isActive(item.to) }"
          >
            <component :is="item.icon" class="h-4.5 w-4.5" />
          </router-link>
        </PTooltip>
      </template>

      <!-- Projects section -->
      <div class="mt-4" v-if="!isCollapsed">
        <div class="mb-1 flex items-center justify-between px-2.5">
          <span class="text-xs font-semibold uppercase tracking-wider text-custom-text-300">Projects</span>
          <router-link
            :to="`/${slug}/projects`"
            class="rounded p-0.5 text-custom-text-300 hover:bg-custom-background-90 hover:text-custom-text-200 transition-colors"
          >
            <Plus class="h-3.5 w-3.5" />
          </router-link>
        </div>
        <div v-for="project in projectStore.projects.slice(0, 8)" :key="project.id" class="space-y-0.5">
          <button
            class="flex w-full items-center gap-2 rounded-md px-2.5 py-1.5 text-sm text-custom-text-200 hover:bg-custom-background-90 transition-colors"
            :class="{ 'bg-custom-background-90 text-custom-text-100': route.path.includes(project.id) }"
            @click="toggleProject(project.id)"
          >
            <ChevronRight
              class="h-3 w-3 flex-shrink-0 text-custom-text-300 transition-transform duration-100"
              :class="{ 'rotate-90': expandedProjects.has(project.id) }"
            />
            <Briefcase class="h-4 w-4 flex-shrink-0 text-custom-text-300" />
            <span class="flex-1 truncate text-left">{{ project.name }}</span>
            <span class="text-2xs text-custom-text-300 font-mono">{{ project.identifier }}</span>
          </button>
          <!-- Expanded sub-nav -->
          <div v-if="expandedProjects.has(project.id)" class="ml-5 space-y-0.5">
            <router-link
              :to="`/${slug}/projects/${project.id}/issues`"
              class="flex items-center gap-2 rounded-md px-2.5 py-1 text-xs text-custom-text-300 hover:bg-custom-background-90 hover:text-custom-text-200 transition-colors"
              :class="{ 'bg-custom-background-90 text-custom-text-100': route.path === `/${slug}/projects/${project.id}/issues` }"
            >
              <LayoutList class="h-3.5 w-3.5" />
              Issues
            </router-link>
            <router-link
              :to="`/${slug}/projects/${project.id}/board`"
              class="flex items-center gap-2 rounded-md px-2.5 py-1 text-xs text-custom-text-300 hover:bg-custom-background-90 hover:text-custom-text-200 transition-colors"
              :class="{ 'bg-custom-background-90 text-custom-text-100': route.path === `/${slug}/projects/${project.id}/board` }"
            >
              <Kanban class="h-3.5 w-3.5" />
              Board
            </router-link>
            <router-link
              :to="`/${slug}/projects/${project.id}/analytics`"
              class="flex items-center gap-2 rounded-md px-2.5 py-1 text-xs text-custom-text-300 hover:bg-custom-background-90 hover:text-custom-text-200 transition-colors"
              :class="{ 'bg-custom-background-90 text-custom-text-100': route.path === `/${slug}/projects/${project.id}/analytics` }"
            >
              <BarChart3 class="h-3.5 w-3.5" />
              Analytics
            </router-link>
            <router-link
              :to="`/${slug}/projects/${project.id}/reports`"
              class="flex items-center gap-2 rounded-md px-2.5 py-1 text-xs text-custom-text-300 hover:bg-custom-background-90 hover:text-custom-text-200 transition-colors"
              :class="{ 'bg-custom-background-90 text-custom-text-100': route.path === `/${slug}/projects/${project.id}/reports` }"
            >
              <FileBarChart class="h-3.5 w-3.5" />
              Reports
            </router-link>
            <router-link
              v-if="project.sprint_view"
              :to="`/${slug}/projects/${project.id}/sprints`"
              class="flex items-center gap-2 rounded-md px-2.5 py-1 text-xs text-custom-text-300 hover:bg-custom-background-90 hover:text-custom-text-200 transition-colors"
              :class="{ 'bg-custom-background-90 text-custom-text-100': route.path.startsWith(`/${slug}/projects/${project.id}/sprints`) }"
            >
              <Repeat class="h-3.5 w-3.5" />
              Sprints
            </router-link>
            <router-link
              v-if="project.epic_view"
              :to="`/${slug}/projects/${project.id}/epics`"
              class="flex items-center gap-2 rounded-md px-2.5 py-1 text-xs text-custom-text-300 hover:bg-custom-background-90 hover:text-custom-text-200 transition-colors"
              :class="{ 'bg-custom-background-90 text-custom-text-100': route.path.startsWith(`/${slug}/projects/${project.id}/epics`) }"
            >
              <Layers class="h-3.5 w-3.5" />
              Epics
            </router-link>
            <router-link
              v-if="project.page_view"
              :to="`/${slug}/projects/${project.id}/pages`"
              class="flex items-center gap-2 rounded-md px-2.5 py-1 text-xs text-custom-text-300 hover:bg-custom-background-90 hover:text-custom-text-200 transition-colors"
              :class="{ 'bg-custom-background-90 text-custom-text-100': route.path.startsWith(`/${slug}/projects/${project.id}/pages`) }"
            >
              <FileText class="h-3.5 w-3.5" />
              Pages
            </router-link>
          </div>
        </div>
      </div>
    </nav>

    <!-- Collapse toggle + User -->
    <div class="border-t border-custom-border-200 p-3">
      <template v-if="!isCollapsed">
        <div class="flex items-center gap-2 rounded-md px-2 py-1.5">
          <PAvatar
            :name="userName"
            :src="authStore.user?.avatar_url"
            size="sm"
          />
          <span class="flex-1 truncate text-sm text-custom-text-200">
            {{ userName }}
          </span>
          <NotificationPopover />
          <button
            @click="toggleTheme"
            class="rounded p-1 text-custom-text-300 hover:bg-custom-background-90 hover:text-custom-text-200 transition-colors"
            :title="theme === 'light' ? 'Switch to dark mode' : 'Switch to light mode'"
          >
            <Moon v-if="theme === 'light'" class="h-4 w-4" />
            <Sun v-else class="h-4 w-4" />
          </button>
          <button
            @click="toggle"
            class="rounded p-1 text-custom-text-300 hover:bg-custom-background-90 hover:text-custom-text-200 transition-colors"
            title="Collapse sidebar"
          >
            <ChevronsLeft class="h-4 w-4" />
          </button>
          <button
            @click="handleLogout"
            class="rounded p-1 text-custom-text-300 hover:bg-custom-background-90 hover:text-custom-text-200 transition-colors"
            title="Logout"
          >
            <LogOut class="h-4 w-4" />
          </button>
        </div>
      </template>
      <template v-else>
        <div class="flex flex-col items-center gap-2">
          <PTooltip :content="userName" position="right">
            <PAvatar
              :name="userName"
              :src="authStore.user?.avatar_url"
              size="sm"
            />
          </PTooltip>
          <PTooltip content="Notifications" position="right">
            <NotificationPopover />
          </PTooltip>
          <PTooltip :content="theme === 'light' ? 'Dark mode' : 'Light mode'" position="right">
            <button
              @click="toggleTheme"
              class="rounded p-1 text-custom-text-300 hover:bg-custom-background-90 hover:text-custom-text-200 transition-colors"
            >
              <Moon v-if="theme === 'light'" class="h-4 w-4" />
              <Sun v-else class="h-4 w-4" />
            </button>
          </PTooltip>
          <PTooltip content="Expand sidebar" position="right">
            <button
              @click="toggle"
              class="rounded p-1 text-custom-text-300 hover:bg-custom-background-90 hover:text-custom-text-200 transition-colors"
            >
              <ChevronsRight class="h-4 w-4" />
            </button>
          </PTooltip>
        </div>
      </template>
    </div>
  </div>
</template>
