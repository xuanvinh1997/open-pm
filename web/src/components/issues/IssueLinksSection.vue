<script setup lang="ts">
import { ref } from 'vue'
import type { IssueLink } from '@/types/issue.types'
import { Plus, X, ExternalLink } from 'lucide-vue-next'

interface Props {
  links: IssueLink[]
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'add-link': [data: { title: string; url: string }]
  'remove-link': [id: string]
}>()

const showForm = ref(false)
const title = ref('')
const url = ref('')

function handleSubmit() {
  if (!title.value.trim() || !url.value.trim()) return
  emit('add-link', {
    title: title.value.trim(),
    url: url.value.trim(),
  })
  title.value = ''
  url.value = ''
  showForm.value = false
}

function getDomain(urlStr: string): string {
  try {
    return new URL(urlStr).hostname
  } catch {
    return urlStr
  }
}
</script>

<template>
  <div class="mb-8">
    <div class="mb-3 flex items-center justify-between">
      <h3 class="text-sm font-semibold text-custom-text-100">Links</h3>
      <button
        @click="showForm = !showForm"
        class="flex items-center gap-1 rounded-md px-2 py-1 text-xs text-custom-primary-100 hover:bg-custom-background-80 transition-colors"
      >
        <Plus class="h-3 w-3" />
        Add link
      </button>
    </div>

    <!-- Add link form -->
    <div v-if="showForm" class="mb-3 rounded-md border border-custom-border-200 p-3 space-y-2">
      <div>
        <label class="mb-1 block text-xs font-medium text-custom-text-300">Title</label>
        <input
          v-model="title"
          type="text"
          placeholder="Link title"
          class="w-full rounded-md border border-custom-border-200 bg-custom-background-100 px-3 py-1.5 text-sm text-custom-text-200 placeholder:text-custom-text-400 focus:border-custom-primary-100 focus:outline-none"
        />
      </div>
      <div>
        <label class="mb-1 block text-xs font-medium text-custom-text-300">URL</label>
        <input
          v-model="url"
          type="url"
          placeholder="https://..."
          class="w-full rounded-md border border-custom-border-200 bg-custom-background-100 px-3 py-1.5 text-sm text-custom-text-200 placeholder:text-custom-text-400 focus:border-custom-primary-100 focus:outline-none"
        />
      </div>
      <div class="flex justify-end gap-2">
        <button
          @click="showForm = false"
          class="rounded-md px-3 py-1.5 text-xs text-custom-text-300 hover:bg-custom-background-80 transition-colors"
        >
          Cancel
        </button>
        <button
          @click="handleSubmit"
          :disabled="!title.trim() || !url.trim()"
          class="rounded-md bg-custom-primary-100 px-3 py-1.5 text-xs text-white hover:bg-custom-primary-200 transition-colors disabled:opacity-50"
        >
          Add
        </button>
      </div>
    </div>

    <!-- Links list -->
    <div v-if="props.links.length > 0" class="space-y-1">
      <div
        v-for="link in props.links"
        :key="link.id"
        class="group flex items-center gap-2 rounded-md border border-custom-border-200 px-3 py-2 text-sm"
      >
        <ExternalLink class="h-3.5 w-3.5 flex-shrink-0 text-custom-text-300" />
        <a
          :href="link.url"
          target="_blank"
          rel="noopener noreferrer"
          class="flex-1 min-w-0 hover:text-custom-primary-100 transition-colors"
        >
          <span class="text-custom-text-100">{{ link.title }}</span>
          <span class="ml-1.5 text-xs text-custom-text-300">{{ getDomain(link.url) }}</span>
        </a>
        <button
          @click="emit('remove-link', link.id)"
          class="flex-shrink-0 rounded p-1 text-custom-text-300 opacity-0 group-hover:opacity-100 hover:bg-red-50 hover:text-red-500 transition-all"
        >
          <X class="h-3 w-3" />
        </button>
      </div>
    </div>
    <p v-else-if="!showForm" class="text-sm text-custom-text-300 italic">No links</p>
  </div>
</template>
