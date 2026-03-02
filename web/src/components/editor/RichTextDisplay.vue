<script setup lang="ts">
import { computed } from 'vue'
import { sanitizeHtml } from '@/utils/editor'

interface Props {
  html?: string
}

const props = defineProps<Props>()

const sanitized = computed(() => (props.html ? sanitizeHtml(props.html) : ''))
</script>

<template>
  <div
    v-if="sanitized"
    class="prose prose-sm max-w-none text-custom-text-200"
    v-html="sanitized"
  />
</template>
