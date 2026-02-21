<script setup lang="ts">
import { computed } from 'vue'
import { ROLE_GUEST, ROLE_MEMBER, ROLE_ADMIN, ROLE_OWNER } from '@/utils/roles'
import PDropdown from './PDropdown.vue'
import { ChevronDown, Check } from 'lucide-vue-next'

interface Props {
  modelValue: number
  maxRole?: number
  disabled?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  maxRole: ROLE_ADMIN,
  disabled: false,
})

const emit = defineEmits<{
  'update:modelValue': [value: number]
}>()

const roleOptions = computed(() => {
  const options = [
    { value: ROLE_GUEST, label: 'Guest', description: 'Read-only access' },
    { value: ROLE_MEMBER, label: 'Member', description: 'Can create and edit issues' },
    { value: ROLE_ADMIN, label: 'Admin', description: 'Can manage members and settings' },
  ]
  if (props.modelValue === ROLE_OWNER) {
    options.push({ value: ROLE_OWNER, label: 'Owner', description: 'Full workspace control' })
  }
  return options
})

const currentLabel = computed(() => {
  const opt = roleOptions.value.find(o => o.value === props.modelValue)
  return opt?.label ?? 'Unknown'
})

function selectRole(value: number, close: () => void) {
  if (value !== props.modelValue && value <= props.maxRole && value !== ROLE_OWNER) {
    emit('update:modelValue', value)
  }
  close()
}
</script>

<template>
  <PDropdown v-if="!disabled" align="right" width="14rem">
    <template #trigger>
      <button
        class="inline-flex items-center gap-1 rounded-md border border-custom-border-200 px-2 py-1 text-xs font-medium text-custom-text-200 hover:bg-custom-background-80 transition-colors"
      >
        {{ currentLabel }}
        <ChevronDown class="h-3 w-3 text-custom-text-300" />
      </button>
    </template>
    <template #default="{ close }">
      <div class="p-1">
        <button
          v-for="option in roleOptions"
          :key="option.value"
          class="flex w-full items-center gap-2 rounded-md px-2.5 py-2 text-left transition-colors"
          :class="[
            option.value === modelValue
              ? 'bg-custom-background-80 text-custom-text-100'
              : option.value > maxRole || option.value === ROLE_OWNER
                ? 'text-custom-text-400 cursor-not-allowed'
                : 'text-custom-text-200 hover:bg-custom-background-80'
          ]"
          :disabled="option.value > maxRole || option.value === ROLE_OWNER"
          @click="selectRole(option.value, close)"
        >
          <div class="flex-1">
            <div class="text-sm font-medium">{{ option.label }}</div>
            <div class="text-2xs text-custom-text-300">{{ option.description }}</div>
          </div>
          <Check v-if="option.value === modelValue" class="h-3.5 w-3.5 text-brand-500 flex-shrink-0" />
        </button>
      </div>
    </template>
  </PDropdown>
  <span
    v-else
    class="inline-flex items-center rounded-md border border-custom-border-200 px-2 py-1 text-xs font-medium text-custom-text-300"
  >
    {{ currentLabel }}
  </span>
</template>
