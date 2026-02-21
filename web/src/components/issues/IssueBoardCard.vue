<script setup lang="ts">
import { computed } from 'vue'
import type { Issue } from '@/types/issue.types'
import { PRIORITY_CONFIG } from '@/utils/issue-helpers'
import PAvatar from '@/components/ui/PAvatar.vue'
import PBadge from '@/components/ui/PBadge.vue'

interface Props {
  issue: Issue
  identifier: string
}

const props = defineProps<Props>()

const emit = defineEmits<{
  click: [issue: Issue]
}>()

const priorityConfig = computed(() => PRIORITY_CONFIG[props.issue.priority] || PRIORITY_CONFIG.none)
</script>

<template>
  <div
    class="rounded-lg border border-custom-border-200 bg-custom-background-100 p-3 shadow-custom-xs hover:shadow-custom-sm transition-shadow cursor-pointer"
    @click="emit('click', props.issue)"
  >
    <!-- Header: Identifier + Priority -->
    <div class="mb-2 flex items-center gap-2">
      <span class="font-mono text-xs text-custom-text-300">
        {{ props.identifier }}-{{ props.issue.sequence_id }}
      </span>
      <component
        :is="priorityConfig.icon"
        class="h-3 w-3 ml-auto"
        :style="{ color: priorityConfig.color }"
      />
    </div>

    <!-- Title -->
    <p class="text-sm text-custom-text-100 line-clamp-2">{{ props.issue.name }}</p>

    <!-- Footer: Labels + Assignees -->
    <div class="mt-3 flex items-center justify-between" v-if="(props.issue.labels && props.issue.labels.length > 0) || (props.issue.assignees && props.issue.assignees.length > 0)">
      <div class="flex items-center gap-1">
        <PBadge
          v-for="label in (props.issue.labels || []).slice(0, 2)"
          :key="label.id"
          :color="label.color"
          variant="soft"
          size="xs"
        >
          {{ label.name }}
        </PBadge>
      </div>
      <div class="flex -space-x-1">
        <PAvatar
          v-for="assignee in (props.issue.assignees || []).slice(0, 2)"
          :key="assignee.id"
          :name="`${assignee.first_name} ${assignee.last_name}`"
          :src="assignee.avatar_url"
          size="xs"
        />
      </div>
    </div>
  </div>
</template>
