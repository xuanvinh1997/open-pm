<script setup lang="ts">
import { ref } from 'vue'
import PTextarea from '@/components/ui/PTextarea.vue'
import PButton from '@/components/ui/PButton.vue'

const emit = defineEmits<{
  submit: [html: string]
}>()

const comment = ref('')
const loading = ref(false)

function handleSubmit() {
  if (!comment.value.trim()) return
  loading.value = true
  emit('submit', comment.value)
  comment.value = ''
  loading.value = false
}
</script>

<template>
  <div class="border-t border-custom-border-200 pt-4">
    <PTextarea
      v-model="comment"
      placeholder="Write a comment..."
      :rows="2"
    />
    <div class="mt-2 flex justify-end">
      <PButton
        variant="primary"
        size="sm"
        :loading="loading"
        :disabled="!comment.trim()"
        @click="handleSubmit"
      >
        Comment
      </PButton>
    </div>
  </div>
</template>
