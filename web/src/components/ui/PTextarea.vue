<script setup lang="ts">
import { ref, watch } from 'vue'

interface Props {
  modelValue?: string
  placeholder?: string
  rows?: number
  error?: string
  disabled?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  placeholder: '',
  rows: 3,
  error: '',
  disabled: false,
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const textareaRef = ref<HTMLTextAreaElement>()

function autoResize() {
  const el = textareaRef.value
  if (!el) return
  el.style.height = 'auto'
  el.style.height = el.scrollHeight + 'px'
}

watch(() => props.modelValue, () => {
  setTimeout(autoResize, 0)
})
</script>

<template>
  <div>
    <textarea
      ref="textareaRef"
      :value="props.modelValue"
      :placeholder="props.placeholder"
      :rows="props.rows"
      :disabled="props.disabled"
      @input="emit('update:modelValue', ($event.target as HTMLTextAreaElement).value); autoResize()"
      class="w-full rounded-md border bg-custom-background-100 px-3 py-2 text-sm text-custom-text-100 placeholder:text-custom-text-300 transition-colors duration-100 outline-none resize-none"
      :class="[
        props.error
          ? 'border-red-500 focus:border-red-500 focus:ring-1 focus:ring-red-500'
          : 'border-custom-border-200 focus:border-brand-500 focus:ring-1 focus:ring-brand-500',
        props.disabled ? 'opacity-50 cursor-not-allowed' : '',
      ]"
    />
    <p v-if="props.error" class="mt-1 text-xs text-red-500">{{ props.error }}</p>
  </div>
</template>
