<script setup lang="ts">
import { ref } from 'vue'
import type { IssueComment } from '@/types/issue.types'
import PAvatar from '@/components/ui/PAvatar.vue'
import PButton from '@/components/ui/PButton.vue'
import { RichTextEditor, RichTextDisplay } from '@/components/editor'
import { formatRelativeDate } from '@/utils/helpers'
import { Pencil, Trash2, X, Check } from 'lucide-vue-next'

interface Props {
  comments: IssueComment[]
  currentUserId?: string
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update-comment': [id: string, html: string, json?: Record<string, unknown>, stripped?: string]
  'delete-comment': [id: string]
}>()

const editingCommentId = ref<string | null>(null)
const editHtml = ref('')
const editJson = ref<Record<string, unknown>>({})
const editStripped = ref('')

function startEdit(comment: IssueComment) {
  editingCommentId.value = comment.id
  editHtml.value = comment.comment_html || ''
  editJson.value = comment.comment_json || {}
  editStripped.value = comment.comment_stripped || ''
}

function cancelEdit() {
  editingCommentId.value = null
  editHtml.value = ''
  editJson.value = {}
  editStripped.value = ''
}

function saveEdit(commentId: string) {
  if (!editStripped.value.trim()) return
  emit('update-comment', commentId, editHtml.value, editJson.value, editStripped.value)
  editingCommentId.value = null
  editHtml.value = ''
  editJson.value = {}
  editStripped.value = ''
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
          <RichTextEditor
            v-model="editHtml"
            :json="editJson"
            toolbar="compact"
            min-height="60px"
            placeholder="Edit comment..."
            autofocus
            @update:json="(v) => editJson = v"
            @update:stripped="(v) => editStripped = v"
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
        <div v-else class="mt-1">
          <RichTextDisplay :html="comment.comment_html || comment.comment_stripped" />
        </div>
      </div>
    </div>
  </div>
</template>
