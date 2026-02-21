<script setup lang="ts">
import type { IssueComment } from '@/types/issue.types'
import PAvatar from '@/components/ui/PAvatar.vue'
import { formatRelativeDate } from '@/utils/helpers'

interface Props {
  comments: IssueComment[]
}

const props = defineProps<Props>()
</script>

<template>
  <div class="space-y-4">
    <h3 class="text-sm font-semibold text-custom-text-100">Activity</h3>

    <div v-if="props.comments.length === 0" class="py-6 text-center text-sm text-custom-text-300">
      No comments yet
    </div>

    <div v-for="comment in props.comments" :key="comment.id" class="flex gap-3">
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
        </div>
        <div
          class="mt-1 text-sm text-custom-text-200 prose prose-sm max-w-none"
          v-html="comment.comment_html || comment.comment_stripped"
        />
      </div>
    </div>
  </div>
</template>
