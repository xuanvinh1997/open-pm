<script setup lang="ts">
import type { UserSummary } from '@/types/issue.types'
import PDropdown from '@/components/ui/PDropdown.vue'
import PAvatar from '@/components/ui/PAvatar.vue'
import { ChevronDown, Check } from 'lucide-vue-next'

interface Props {
  modelValue?: string[]
  members: UserSummary[]
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: () => [],
})

const emit = defineEmits<{
  'update:modelValue': [value: string[]]
}>()

function toggleMember(userId: string) {
  const current = [...props.modelValue]
  const index = current.indexOf(userId)
  if (index >= 0) {
    current.splice(index, 1)
  } else {
    current.push(userId)
  }
  emit('update:modelValue', current)
}

function isSelected(userId: string) {
  return props.modelValue.includes(userId)
}
</script>

<template>
  <PDropdown align="left" width="14rem">
    <template #trigger>
      <button
        type="button"
        class="flex items-center gap-2 rounded-md border border-custom-border-200 bg-custom-background-100 px-3 py-2 text-sm text-custom-text-100 hover:bg-custom-background-90 transition-colors w-full"
      >
        <div v-if="props.modelValue.length > 0" class="flex -space-x-1">
          <PAvatar
            v-for="id in props.modelValue.slice(0, 3)"
            :key="id"
            :name="props.members.find(m => m.id === id)?.first_name || ''"
            size="xs"
          />
        </div>
        <span v-else class="text-custom-text-300">Assignees</span>
        <span v-if="props.modelValue.length > 0" class="text-xs text-custom-text-300">{{ props.modelValue.length }}</span>
        <ChevronDown class="h-3.5 w-3.5 text-custom-text-300 flex-shrink-0 ml-auto" />
      </button>
    </template>
    <template #default>
      <button
        v-for="member in props.members"
        :key="member.id"
        class="flex w-full items-center gap-2 rounded-md px-3 py-2 text-sm hover:bg-custom-background-80 transition-colors text-custom-text-200"
        @click="toggleMember(member.id)"
      >
        <PAvatar :name="`${member.first_name} ${member.last_name}`" :src="member.avatar_url" size="xs" />
        <span class="flex-1 truncate">{{ member.first_name }} {{ member.last_name }}</span>
        <Check v-if="isSelected(member.id)" class="h-3.5 w-3.5 text-brand-600 flex-shrink-0" />
      </button>
      <div v-if="props.members.length === 0" class="px-3 py-4 text-center text-xs text-custom-text-300">
        No members
      </div>
    </template>
  </PDropdown>
</template>
