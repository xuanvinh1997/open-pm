<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import type { IssuePriority, IssueType } from '@/types/issue.types'
import type { ProjectMember } from '@/types/project.types'
import { PRIORITY_CONFIG, ISSUE_TYPE_CONFIG } from '@/utils/issue-helpers'
import PAvatar from '@/components/ui/PAvatar.vue'
import { X } from 'lucide-vue-next'

export interface IssueFilters {
  priority: IssuePriority[]
  type: IssueType[]
  assignee: string[]
}

interface Props {
  members: ProjectMember[]
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:filters': [filters: IssueFilters]
}>()

const selectedPriorities = ref<IssuePriority[]>([])
const selectedTypes = ref<IssueType[]>([])
const selectedAssignees = ref<string[]>([])

const priorities = Object.entries(PRIORITY_CONFIG) as [IssuePriority, typeof PRIORITY_CONFIG[IssuePriority]][]
const types = Object.entries(ISSUE_TYPE_CONFIG) as [IssueType, typeof ISSUE_TYPE_CONFIG[IssueType]][]

const hasFilters = computed(() =>
  selectedPriorities.value.length > 0 || selectedTypes.value.length > 0 || selectedAssignees.value.length > 0
)

function togglePriority(p: IssuePriority) {
  const idx = selectedPriorities.value.indexOf(p)
  if (idx === -1) selectedPriorities.value.push(p)
  else selectedPriorities.value.splice(idx, 1)
}

function toggleType(t: IssueType) {
  const idx = selectedTypes.value.indexOf(t)
  if (idx === -1) selectedTypes.value.push(t)
  else selectedTypes.value.splice(idx, 1)
}

function toggleAssignee(id: string) {
  const idx = selectedAssignees.value.indexOf(id)
  if (idx === -1) selectedAssignees.value.push(id)
  else selectedAssignees.value.splice(idx, 1)
}

function clearAll() {
  selectedPriorities.value = []
  selectedTypes.value = []
  selectedAssignees.value = []
}

watch([selectedPriorities, selectedTypes, selectedAssignees], () => {
  emit('update:filters', {
    priority: [...selectedPriorities.value],
    type: [...selectedTypes.value],
    assignee: [...selectedAssignees.value],
  })
}, { deep: true })
</script>

<template>
  <div class="flex items-center gap-2 border-b border-custom-border-200 bg-custom-background-100 px-6 py-2 text-xs">
    <span class="text-custom-text-300 mr-1">Filters:</span>

    <!-- Priority filters -->
    <div class="flex items-center gap-1">
      <button
        v-for="[key, config] in priorities"
        :key="key"
        class="rounded-md px-2 py-1 transition-colors border"
        :class="selectedPriorities.includes(key)
          ? 'border-brand-500 bg-brand-50 text-brand-600'
          : 'border-transparent text-custom-text-300 hover:bg-custom-background-80 hover:text-custom-text-200'"
        @click="togglePriority(key)"
      >
        <component :is="config.icon" class="inline h-3 w-3 mr-0.5" :style="{ color: config.color }" />
        {{ config.label }}
      </button>
    </div>

    <span class="text-custom-border-200">|</span>

    <!-- Type filters -->
    <div class="flex items-center gap-1">
      <button
        v-for="[key, config] in types"
        :key="key"
        class="rounded-md px-2 py-1 transition-colors border"
        :class="selectedTypes.includes(key)
          ? 'border-brand-500 bg-brand-50 text-brand-600'
          : 'border-transparent text-custom-text-300 hover:bg-custom-background-80 hover:text-custom-text-200'"
        @click="toggleType(key)"
      >
        <component :is="config.icon" class="inline h-3 w-3 mr-0.5" :style="{ color: config.color }" />
        {{ config.label }}
      </button>
    </div>

    <span class="text-custom-border-200">|</span>

    <!-- Assignee filters -->
    <div class="flex items-center gap-1">
      <button
        v-for="member in props.members.slice(0, 5)"
        :key="member.user_id"
        class="flex items-center gap-1 rounded-md px-1.5 py-1 transition-colors border"
        :class="selectedAssignees.includes(member.user_id)
          ? 'border-brand-500 bg-brand-50 text-brand-600'
          : 'border-transparent text-custom-text-300 hover:bg-custom-background-80 hover:text-custom-text-200'"
        @click="toggleAssignee(member.user_id)"
      >
        <PAvatar
          :name="`${member.first_name} ${member.last_name}`"
          :src="member.avatar_url"
          size="xs"
        />
        <span>{{ member.display_name || member.first_name }}</span>
      </button>
    </div>

    <!-- Clear all -->
    <button
      v-if="hasFilters"
      class="ml-auto flex items-center gap-1 rounded-md px-2 py-1 text-custom-text-300 hover:bg-custom-background-80 hover:text-custom-text-200 transition-colors"
      @click="clearAll"
    >
      <X class="h-3 w-3" />
      Clear
    </button>
  </div>
</template>
