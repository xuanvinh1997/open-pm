<script setup lang="ts">
import { ChevronRight } from 'lucide-vue-next'
import type { Component } from 'vue'

interface BreadcrumbItem {
  label: string
  to?: string
  icon?: Component
}

interface Props {
  items: BreadcrumbItem[]
}

const props = defineProps<Props>()
</script>

<template>
  <nav class="flex items-center gap-1.5 text-sm">
    <template v-for="(item, index) in props.items" :key="index">
      <ChevronRight v-if="index > 0" class="h-3.5 w-3.5 text-custom-text-300 flex-shrink-0" />
      <router-link
        v-if="item.to"
        :to="item.to"
        class="flex items-center gap-1.5 text-custom-text-300 hover:text-custom-text-100 transition-colors"
      >
        <component v-if="item.icon" :is="item.icon" class="h-4 w-4" />
        <span>{{ item.label }}</span>
      </router-link>
      <span
        v-else
        class="flex items-center gap-1.5 font-medium text-custom-text-100"
      >
        <component v-if="item.icon" :is="item.icon" class="h-4 w-4" />
        <span>{{ item.label }}</span>
      </span>
    </template>
  </nav>
</template>
