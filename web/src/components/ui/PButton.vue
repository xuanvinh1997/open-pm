<script setup lang="ts">
import PSpinner from './PSpinner.vue'

interface Props {
  variant?: 'primary' | 'secondary' | 'outline' | 'ghost' | 'danger'
  size?: 'sm' | 'md' | 'lg'
  loading?: boolean
  disabled?: boolean
  type?: 'button' | 'submit' | 'reset'
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'primary',
  size: 'md',
  loading: false,
  disabled: false,
  type: 'button',
})

const variantClasses: Record<string, string> = {
  primary: 'bg-brand-600 text-white hover:bg-brand-700 shadow-custom-xs',
  secondary: 'bg-custom-background-80 text-custom-text-200 border border-custom-border-200 hover:bg-custom-background-90',
  outline: 'border border-custom-border-200 text-custom-text-200 hover:bg-custom-background-80',
  ghost: 'text-custom-text-200 hover:bg-custom-background-80',
  danger: 'bg-red-500 text-white hover:bg-red-600 shadow-custom-xs',
}

const sizeClasses: Record<string, string> = {
  sm: 'px-2.5 py-1.5 text-xs gap-1.5',
  md: 'px-3.5 py-2 text-sm gap-2',
  lg: 'px-4 py-2.5 text-sm gap-2',
}
</script>

<template>
  <button
    :type="props.type"
    :disabled="props.disabled || props.loading"
    class="inline-flex items-center justify-center rounded-md font-medium transition-colors duration-100 disabled:opacity-50 disabled:cursor-not-allowed"
    :class="[variantClasses[props.variant], sizeClasses[props.size]]"
  >
    <PSpinner v-if="props.loading" size="sm" />
    <slot />
  </button>
</template>
