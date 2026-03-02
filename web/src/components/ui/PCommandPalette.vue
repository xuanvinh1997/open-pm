<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useCommandPalette } from '@/composables/useCommandPalette'
import { useWorkspaceStore } from '@/stores/workspace.store'
import { useProjectStore } from '@/stores/project.store'
import { searchApi } from '@/api/search.api'
import type { Issue } from '@/types/issue.types'
import type { Page } from '@/types/page.types'
import { Search, Home, Briefcase, Settings, FolderKanban, LayoutList, FileText } from 'lucide-vue-next'

const router = useRouter()
const { isOpen, close } = useCommandPalette()
const workspaceStore = useWorkspaceStore()
const projectStore = useProjectStore()

const query = ref('')
const selectedIndex = ref(0)
const searchResults = ref<{ issues: Issue[]; pages: Page[] }>({ issues: [], pages: [] })
const isSearching = ref(false)

interface CommandItem {
  id: string
  label: string
  description?: string
  icon: typeof Home
  action: () => void
}

const staticItems = computed<CommandItem[]>(() => {
  const slug = workspaceStore.currentWorkspace?.slug
  if (!slug) return []

  const all: CommandItem[] = [
    { id: 'home', label: 'Home', description: 'Go to dashboard', icon: Home, action: () => router.push(`/${slug}`) },
    { id: 'projects', label: 'Projects', description: 'View all projects', icon: Briefcase, action: () => router.push(`/${slug}/projects`) },
    { id: 'settings', label: 'Settings', description: 'Workspace settings', icon: Settings, action: () => router.push(`/${slug}/settings`) },
  ]

  for (const project of projectStore.projects) {
    all.push({
      id: `project-${project.id}`,
      label: project.name,
      description: project.identifier,
      icon: FolderKanban,
      action: () => router.push(`/${slug}/projects/${project.id}/issues`),
    })
  }

  return all
})

const items = computed<CommandItem[]>(() => {
  const slug = workspaceStore.currentWorkspace?.slug
  if (!slug) return []

  // If we have search results from the API, show those
  if (query.value.length >= 3 && (searchResults.value.issues.length > 0 || searchResults.value.pages.length > 0)) {
    const results: CommandItem[] = []
    for (const issue of searchResults.value.issues) {
      results.push({
        id: `issue-${issue.id}`,
        label: issue.name,
        description: issue.priority,
        icon: LayoutList,
        action: () => router.push(`/${slug}/projects/${issue.project_id}/issues/${issue.id}`),
      })
    }
    for (const page of searchResults.value.pages) {
      results.push({
        id: `page-${page.id}`,
        label: page.name,
        description: 'Page',
        icon: FileText,
        action: () => router.push(`/${slug}/projects/${page.project_id}/pages/${page.id}`),
      })
    }
    return results
  }

  // Otherwise filter static items
  if (!query.value) return staticItems.value

  const q = query.value.toLowerCase()
  return staticItems.value.filter(
    (item) =>
      item.label.toLowerCase().includes(q) ||
      item.description?.toLowerCase().includes(q),
  )
})

function handleSelect(item: CommandItem) {
  item.action()
  close()
  query.value = ''
}

function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'ArrowDown') {
    e.preventDefault()
    selectedIndex.value = Math.min(selectedIndex.value + 1, items.value.length - 1)
  } else if (e.key === 'ArrowUp') {
    e.preventDefault()
    selectedIndex.value = Math.max(selectedIndex.value - 1, 0)
  } else if (e.key === 'Enter') {
    e.preventDefault()
    if (items.value[selectedIndex.value]) {
      handleSelect(items.value[selectedIndex.value])
    }
  } else if (e.key === 'Escape') {
    close()
    query.value = ''
  }
}

function handleBackdropClick(e: MouseEvent) {
  if (e.target === e.currentTarget) {
    close()
    query.value = ''
  }
}

let searchTimeout: ReturnType<typeof setTimeout> | null = null
watch(query, (val) => {
  selectedIndex.value = 0

  // Debounced API search for 3+ character queries
  if (searchTimeout) clearTimeout(searchTimeout)
  if (val.length >= 3) {
    isSearching.value = true
    searchTimeout = setTimeout(async () => {
      const slug = workspaceStore.currentWorkspace?.slug
      if (!slug) return
      try {
        const { data } = await searchApi.search(slug, val)
        searchResults.value = data
      } catch {
        searchResults.value = { issues: [], pages: [] }
      } finally {
        isSearching.value = false
      }
    }, 300)
  } else {
    searchResults.value = { issues: [], pages: [] }
    isSearching.value = false
  }
})

watch(isOpen, (val) => {
  if (val) {
    selectedIndex.value = 0
    query.value = ''
  }
})
</script>

<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition duration-150 ease-out"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition duration-100 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="isOpen"
        class="fixed inset-0 z-[100] flex items-start justify-center bg-black/50 backdrop-blur-sm pt-[20vh]"
        @click="handleBackdropClick"
        @keydown="handleKeydown"
      >
        <div class="w-full max-w-xl rounded-xl border border-custom-border-200 bg-custom-background-100 shadow-custom-lg overflow-hidden">
          <!-- Search input -->
          <div class="flex items-center gap-3 border-b border-custom-border-200 px-4">
            <Search class="h-5 w-5 text-custom-text-300 flex-shrink-0" />
            <input
              v-model="query"
              type="text"
              placeholder="Search or jump to..."
              class="w-full border-0 bg-transparent py-3.5 text-sm text-custom-text-100 placeholder:text-custom-text-300 outline-none focus:ring-0"
              autofocus
            />
            <kbd class="hidden sm:inline-flex items-center rounded border border-custom-border-200 px-1.5 py-0.5 text-2xs text-custom-text-300 font-mono">
              ESC
            </kbd>
          </div>

          <!-- Results -->
          <div class="max-h-[300px] overflow-y-auto py-2">
            <div v-if="items.length === 0" class="px-4 py-8 text-center text-sm text-custom-text-300">
              No results found
            </div>
            <button
              v-for="(item, index) in items"
              :key="item.id"
              class="flex w-full items-center gap-3 px-4 py-2.5 text-left text-sm transition-colors"
              :class="index === selectedIndex ? 'bg-custom-background-80 text-custom-text-100' : 'text-custom-text-200 hover:bg-custom-background-80'"
              @click="handleSelect(item)"
              @mouseenter="selectedIndex = index"
            >
              <component :is="item.icon" class="h-4 w-4 flex-shrink-0 text-custom-text-300" />
              <span class="flex-1">{{ item.label }}</span>
              <span v-if="item.description" class="text-xs text-custom-text-300">{{ item.description }}</span>
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>
