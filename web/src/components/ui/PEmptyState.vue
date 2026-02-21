<script setup lang="ts">
import type { Component } from 'vue'

interface Props {
  title: string
  description?: string
  icon?: Component
}

const props = defineProps<Props>()
</script>

<template>
  <div class="flex flex-col items-center justify-center py-20 text-center">
    <component
      v-if="props.icon"
      :is="props.icon"
      class="mb-5 h-14 w-14 text-custom-text-300"
    />
    <h3 class="mb-2 text-lg font-semibold text-custom-text-100">{{ props.title }}</h3>
    <p v-if="props.description" class="mb-6 max-w-md text-sm text-custom-text-300">
      {{ props.description }}
    </p>
    <slot />
  </div>
</template>
