<script setup lang="ts">
import { computed } from 'vue'
import type { Issue } from '@/types/issue.types'
import type { State } from '@/types/project.types'
import { PRIORITY_CONFIG, ISSUE_TYPE_CONFIG } from '@/utils/issue-helpers'
import { formatDate } from '@/utils/helpers'
import PAvatar from '@/components/ui/PAvatar.vue'
import PBadge from '@/components/ui/PBadge.vue'
import { Circle } from 'lucide-vue-next'

interface Props {
  issue: Issue
  identifier: string
  state?: State
}

const props = defineProps<Props>()

const emit = defineEmits<{
  click: [issue: Issue]
}>()

const priorityConfig = computed(() => PRIORITY_CONFIG[props.issue.priority] || PRIORITY_CONFIG.none)
const typeConfig = computed(() => ISSUE_TYPE_CONFIG[props.issue.issue_type] || ISSUE_TYPE_CONFIG.task)

const isOverdue = computed(() => {
  if (!props.issue.target_date || props.issue.completed_at) return false
  return new Date(props.issue.target_date) < new Date(new Date().toDateString())
})
</script>

<template>
  <div
    class="flex items-center gap-3 border-b border-custom-border-200 px-4 py-2.5 hover:bg-custom-background-80 cursor-pointer transition-colors"
    @click="emit('click', props.issue)"
  >
    <!-- Type icon -->
    <component
      :is="typeConfig.icon"
      class="h-3.5 w-3.5 flex-shrink-0"
      :style="{ color: typeConfig.color }"
    />

    <!-- Priority icon -->
    <component
      :is="priorityConfig.icon"
      class="h-3.5 w-3.5 flex-shrink-0"
      :style="{ color: priorityConfig.color }"
    />

    <!-- State dot -->
    <Circle
      v-if="props.state"
      class="h-3 w-3 flex-shrink-0"
      :style="{ color: props.state.color }"
      fill="currentColor"
    />

    <!-- Identifier -->
    <span class="flex-shrink-0 font-mono text-xs text-custom-text-300">
      {{ props.identifier }}-{{ props.issue.sequence_id }}
    </span>

    <!-- Title -->
    <span class="flex-1 truncate text-sm text-custom-text-100">
      {{ props.issue.name }}
    </span>

    <!-- Labels -->
    <div v-if="props.issue.labels && props.issue.labels.length > 0" class="flex items-center gap-1 flex-shrink-0">
      <PBadge
        v-for="label in props.issue.labels.slice(0, 2)"
        :key="label.id"
        :color="label.color"
        variant="soft"
        size="xs"
      >
        {{ label.name }}
      </PBadge>
    </div>

    <!-- Assignees -->
    <div v-if="props.issue.assignees && props.issue.assignees.length > 0" class="flex -space-x-1 flex-shrink-0">
      <PAvatar
        v-for="assignee in props.issue.assignees.slice(0, 2)"
        :key="assignee.id"
        :name="`${assignee.first_name} ${assignee.last_name}`"
        :src="assignee.avatar_url"
        size="xs"
      />
    </div>

    <!-- Due date -->
    <span
      v-if="props.issue.target_date"
      class="flex-shrink-0 text-xs"
      :class="isOverdue ? 'text-red-500 font-medium' : 'text-custom-text-300'"
    >
      {{ formatDate(props.issue.target_date) }}
    </span>
  </div>
</template>
