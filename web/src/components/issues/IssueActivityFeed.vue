<script setup lang="ts">
import { ref, nextTick } from 'vue'
import type { IssueComment } from '@/types/issue.types'
import PAvatar from '@/components/ui/PAvatar.vue'
import PButton from '@/components/ui/PButton.vue'
import { formatRelativeDate } from '@/utils/helpers'
import { Pencil, Trash2, X, Check } from 'lucide-vue-next'

interface Props {
  comments: IssueComment[]
  currentUserId?: string
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update-comment': [id: string, html: string]
  'delete-comment': [id: string]
}>()

const editingCommentId = ref<string | null>(null)
const editText = ref('')
const editInputRef = ref<HTMLTextAreaElement>()

function startEdit(comment: IssueComment) {
  editingCommentId.value = comment.id
  editText.value = comment.comment_stripped || ''
  nextTick(() => editInputRef.value?.focus())
}

function cancelEdit() {
  editingCommentId.value = null
  editText.value = ''
}

function saveEdit(commentId: string) {
  if (!editText.value.trim()) return
  const html = `<p>${editText.value.trim().replace(/\n/g, '</p><p>')}</p>`
  emit('update-comment', commentId, html)
  editingCommentId.value = null
  editText.value = ''
}

function isOwner(comment: IssueComment) {
  return props.currentUserId && comment.actor_id === props.currentUserId
}
</script>

<template>
  <div class="space-y-4">
    <h3 class="text-sm font-semibold text-custom-text-100">Activity</h3>

    <div v-if="props.comments.length === 0" class="py-6 text-center text-sm text-custom-text-300">
      No comments yet
    </div>

    <div v-for="comment in props.comments" :key="comment.id" class="group flex gap-3">
      <PAvatar
        :name="`${comment.first_name} ${comment.last_name}`"
        :src="comment.avatar_url"
        size="sm"
      />
      <div class="flex-1 min-w-0">
        <div class="flex items-center gap-2">
          <span class="text-sm font-medium text-custom-text-100">
            {{ comment.display_name || `${comment.first_name} ${comment.last_name}` }}
          </span>
          <span class="text-xs text-custom-text-300">{{ formatRelativeDate(comment.created_at) }}</span>
          <div v-if="isOwner(comment) && editingCommentId !== comment.id" class="ml-auto flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
            <button
              class="rounded p-1 text-custom-text-300 hover:bg-custom-background-80 hover:text-custom-text-100 transition-colors"
              title="Edit comment"
              @click="startEdit(comment)"
            >
              <Pencil class="h-3 w-3" />
            </button>
            <button
              class="rounded p-1 text-custom-text-300 hover:bg-red-50 hover:text-red-500 transition-colors"
              title="Delete comment"
              @click="emit('delete-comment', comment.id)"
            >
              <Trash2 class="h-3 w-3" />
            </button>
          </div>
        </div>
        <!-- Edit mode -->
        <div v-if="editingCommentId === comment.id" class="mt-1">
          <textarea
            ref="editInputRef"
            v-model="editText"
            class="w-full rounded-md border border-custom-border-200 bg-custom-background-100 px-3 py-2 text-sm text-custom-text-200 outline-none focus:border-brand-500 focus:ring-1 focus:ring-brand-500 resize-none"
            :rows="3"
            @keydown.escape="cancelEdit"
          />
          <div class="mt-1 flex items-center gap-2">
            <PButton variant="primary" size="sm" @click="saveEdit(comment.id)">
              <Check class="h-3 w-3" />
              Save
            </PButton>
            <PButton variant="ghost" size="sm" @click="cancelEdit">
              <X class="h-3 w-3" />
              Cancel
            </PButton>
          </div>
        </div>
        <!-- Display mode -->
        <div
          v-else
          class="mt-1 text-sm text-custom-text-200 prose prose-sm max-w-none"
          v-html="comment.comment_html || comment.comment_stripped"
        />
      </div>
    </div>
  </div>
</template>
