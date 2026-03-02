<script setup lang="ts">
import { ref } from 'vue'
import { RichTextEditor } from '@/components/editor'
import PButton from '@/components/ui/PButton.vue'

const emit = defineEmits<{
  submit: [html: string, json?: Record<string, unknown>, stripped?: string]
}>()

const html = ref('')
const json = ref<Record<string, unknown>>({})
const stripped = ref('')
const editorKey = ref(0)

function handleSubmit() {
  if (!stripped.value.trim()) return
  emit('submit', html.value, json.value, stripped.value)
  html.value = ''
  json.value = {}
  stripped.value = ''
  editorKey.value++
}
</script>

<template>
  <div class="border-t border-custom-border-200 pt-4">
    <RichTextEditor
      :key="editorKey"
      v-model="html"
      toolbar="compact"
      placeholder="Write a comment..."
      min-height="60px"
      @update:json="(v) => json = v"
      @update:stripped="(v) => stripped = v"
      @submit="handleSubmit"
    />
    <div class="mt-2 flex justify-end">
      <PButton
        variant="primary"
        size="sm"
        :disabled="!stripped.trim()"
        @click="handleSubmit"
      >
        Comment
      </PButton>
    </div>
  </div>
</template>
