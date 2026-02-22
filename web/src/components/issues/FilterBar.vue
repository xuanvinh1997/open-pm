<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import type { State, Label, ProjectMember } from '@/types/project.types'
import { Search, X, SlidersHorizontal } from 'lucide-vue-next'

const props = defineProps<{
  states: State[]
  labels: Label[]
  members: ProjectMember[]
}>()

const emit = defineEmits<{
  (e: 'filter-change', filters: Record<string, string>): void
}>()

const route = useRoute()
const router = useRouter()

const searchText = ref((route.query.search as string) || '')
const selectedPriority = ref((route.query.priority as string) || '')
const selectedType = ref((route.query.type as string) || '')
const selectedState = ref((route.query.state as string) || '')
const selectedAssignee = ref((route.query.assignee as string) || '')
const selectedLabel = ref((route.query.label as string) || '')
const sortBy = ref((route.query.sort_by as string) || '')
const sortOrder = ref((route.query.sort_order as string) || '')
const showFilters = ref(false)

const priorities = [
  { value: 'urgent', label: 'Urgent' },
  { value: 'high', label: 'High' },
  { value: 'medium', label: 'Medium' },
  { value: 'low', label: 'Low' },
  { value: 'none', label: 'None' },
]

const issueTypes = [
  { value: 'story', label: 'Story' },
  { value: 'bug', label: 'Bug' },
  { value: 'task', label: 'Task' },
  { value: 'epic', label: 'Epic' },
]

const sortOptions = [
  { value: '', label: 'Default' },
  { value: 'created_at', label: 'Created' },
  { value: 'updated_at', label: 'Updated' },
  { value: 'priority', label: 'Priority' },
  { value: 'name', label: 'Name' },
]

function emitFilters() {
  const filters: Record<string, string> = {}
  if (searchText.value) filters.search = searchText.value
  if (selectedPriority.value) filters.priority = selectedPriority.value
  if (selectedType.value) filters.type = selectedType.value
  if (selectedState.value) filters.state = selectedState.value
  if (selectedAssignee.value) filters.assignee = selectedAssignee.value
  if (selectedLabel.value) filters.label = selectedLabel.value
  if (sortBy.value) filters.sort_by = sortBy.value
  if (sortOrder.value) filters.sort_order = sortOrder.value

  // Update URL query params
  router.replace({ query: { ...route.query, ...filters, page: undefined } })
  emit('filter-change', filters)
}

function clearFilters() {
  searchText.value = ''
  selectedPriority.value = ''
  selectedType.value = ''
  selectedState.value = ''
  selectedAssignee.value = ''
  selectedLabel.value = ''
  sortBy.value = ''
  sortOrder.value = ''
  router.replace({ query: {} })
  emit('filter-change', {})
}

const hasActiveFilters = ref(false)
watch(
  [searchText, selectedPriority, selectedType, selectedState, selectedAssignee, selectedLabel, sortBy, sortOrder],
  () => {
    hasActiveFilters.value = !!(
      searchText.value || selectedPriority.value || selectedType.value ||
      selectedState.value || selectedAssignee.value || selectedLabel.value ||
      sortBy.value || sortOrder.value
    )
  },
  { immediate: true },
)

let searchTimeout: ReturnType<typeof setTimeout> | null = null
function handleSearchInput() {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => emitFilters(), 300)
}
</script>

<template>
  <div class="space-y-2">
    <!-- Search + toggle row -->
    <div class="flex items-center gap-2">
      <div class="relative flex-1 max-w-xs">
        <Search class="absolute left-2.5 top-1/2 -translate-y-1/2 h-3.5 w-3.5 text-custom-text-300" />
        <input
          v-model="searchText"
          type="text"
          placeholder="Search issues..."
          class="w-full rounded-md border border-custom-border-200 bg-custom-background-100 py-1.5 pl-8 pr-3 text-sm text-custom-text-100 placeholder:text-custom-text-300 outline-none focus:border-brand-500 transition-colors"
          @input="handleSearchInput"
        />
      </div>
      <button
        class="flex items-center gap-1.5 rounded-md border border-custom-border-200 px-2.5 py-1.5 text-sm transition-colors"
        :class="showFilters || hasActiveFilters
          ? 'bg-brand-500/10 border-brand-500/30 text-brand-500'
          : 'text-custom-text-200 hover:bg-custom-background-80'"
        @click="showFilters = !showFilters"
      >
        <SlidersHorizontal class="h-3.5 w-3.5" />
        Filters
      </button>
      <button
        v-if="hasActiveFilters"
        class="flex items-center gap-1 rounded-md px-2 py-1.5 text-xs text-custom-text-300 hover:text-custom-text-200 transition-colors"
        @click="clearFilters"
      >
        <X class="h-3 w-3" />
        Clear
      </button>
    </div>

    <!-- Filter dropdowns -->
    <div v-if="showFilters" class="flex flex-wrap items-center gap-2">
      <select
        v-model="selectedPriority"
        class="rounded-md border border-custom-border-200 bg-custom-background-100 px-2 py-1 text-xs text-custom-text-200 outline-none"
        @change="emitFilters"
      >
        <option value="">All priorities</option>
        <option v-for="p in priorities" :key="p.value" :value="p.value">{{ p.label }}</option>
      </select>

      <select
        v-model="selectedType"
        class="rounded-md border border-custom-border-200 bg-custom-background-100 px-2 py-1 text-xs text-custom-text-200 outline-none"
        @change="emitFilters"
      >
        <option value="">All types</option>
        <option v-for="t in issueTypes" :key="t.value" :value="t.value">{{ t.label }}</option>
      </select>

      <select
        v-model="selectedState"
        class="rounded-md border border-custom-border-200 bg-custom-background-100 px-2 py-1 text-xs text-custom-text-200 outline-none"
        @change="emitFilters"
      >
        <option value="">All states</option>
        <option v-for="s in states" :key="s.id" :value="s.id">{{ s.name }}</option>
      </select>

      <select
        v-model="selectedAssignee"
        class="rounded-md border border-custom-border-200 bg-custom-background-100 px-2 py-1 text-xs text-custom-text-200 outline-none"
        @change="emitFilters"
      >
        <option value="">All assignees</option>
        <option v-for="m in members" :key="m.user_id" :value="m.user_id">
          {{ m.display_name || m.first_name || m.email }}
        </option>
      </select>

      <select
        v-model="selectedLabel"
        class="rounded-md border border-custom-border-200 bg-custom-background-100 px-2 py-1 text-xs text-custom-text-200 outline-none"
        @change="emitFilters"
      >
        <option value="">All labels</option>
        <option v-for="l in labels" :key="l.id" :value="l.id">{{ l.name }}</option>
      </select>

      <div class="h-4 w-px bg-custom-border-200" />

      <select
        v-model="sortBy"
        class="rounded-md border border-custom-border-200 bg-custom-background-100 px-2 py-1 text-xs text-custom-text-200 outline-none"
        @change="emitFilters"
      >
        <option v-for="s in sortOptions" :key="s.value" :value="s.value">Sort: {{ s.label }}</option>
      </select>

      <select
        v-model="sortOrder"
        class="rounded-md border border-custom-border-200 bg-custom-background-100 px-2 py-1 text-xs text-custom-text-200 outline-none"
        @change="emitFilters"
      >
        <option value="">Default</option>
        <option value="asc">Ascending</option>
        <option value="desc">Descending</option>
      </select>
    </div>
  </div>
</template>
