<script setup lang="ts">
import type { Component } from 'vue'

interface Props {
  modelValue?: string
  type?: string
  placeholder?: string
  size?: 'sm' | 'md'
  error?: string
  disabled?: boolean
  icon?: Component
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  type: 'text',
  placeholder: '',
  size: 'md',
  error: '',
  disabled: false,
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const sizeClasses: Record<string, string> = {
  sm: 'px-2.5 py-1.5 text-xs',
  md: 'px-3 py-2 text-sm',
}
</script>

<template>
  <div class="relative">
    <div v-if="props.icon" class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
      <component :is="props.icon" class="h-4 w-4 text-custom-text-300" />
    </div>
    <input
      :type="props.type"
      :value="props.modelValue"
      :placeholder="props.placeholder"
      :disabled="props.disabled"
      @input="emit('update:modelValue', ($event.target as HTMLInputElement).value)"
      class="w-full rounded-md border bg-custom-background-100 text-custom-text-100 placeholder:text-custom-text-300 transition-colors duration-100 outline-none"
      :class="[
        sizeClasses[props.size],
        props.icon ? 'pl-9' : '',
        props.error
          ? 'border-red-500 focus:border-red-500 focus:ring-1 focus:ring-red-500'
          : 'border-custom-border-200 focus:border-brand-500 focus:ring-1 focus:ring-brand-500',
        props.disabled ? 'opacity-50 cursor-not-allowed' : '',
      ]"
    />
    <p v-if="props.error" class="mt-1 text-xs text-red-500">{{ props.error }}</p>
  </div>
</template>
