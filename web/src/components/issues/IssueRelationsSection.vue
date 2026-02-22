<script setup lang="ts">
import { ref } from 'vue'
import type { IssueRelation, RelationType } from '@/types/issue.types'
import { RELATION_TYPE_CONFIG } from '@/utils/issue-helpers'
import PDropdown from '@/components/ui/PDropdown.vue'
import { Plus, X, ChevronDown } from 'lucide-vue-next'

interface Props {
  relations: IssueRelation[]
  projectIdentifier: string
  projectId: string
  workspaceSlug: string
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'add-relation': [data: { related_issue_id: string; relation_type: RelationType }]
  'remove-relation': [relationId: string]
}>()

const showForm = ref(false)
const relationType = ref<RelationType>('relates_to')
const relatedIssueId = ref('')

const relationTypes: RelationType[] = ['relates_to', 'blocks', 'blocked_by', 'duplicate_of']

function handleSubmit() {
  if (!relatedIssueId.value.trim()) return
  emit('add-relation', {
    related_issue_id: relatedIssueId.value.trim(),
    relation_type: relationType.value,
  })
  relatedIssueId.value = ''
  showForm.value = false
}
</script>

<template>
  <div class="mb-8">
    <div class="mb-3 flex items-center justify-between">
      <h3 class="text-sm font-semibold text-custom-text-100">Relations</h3>
      <button
        @click="showForm = !showForm"
        class="flex items-center gap-1 rounded-md px-2 py-1 text-xs text-custom-primary-100 hover:bg-custom-background-80 transition-colors"
      >
        <Plus class="h-3 w-3" />
        Add relation
      </button>
    </div>

    <!-- Add relation form -->
    <div v-if="showForm" class="mb-3 rounded-md border border-custom-border-200 p-3 space-y-2">
      <div>
        <label class="mb-1 block text-xs font-medium text-custom-text-300">Relation type</label>
        <PDropdown align="left" width="12rem">
          <template #trigger>
            <button
              type="button"
              class="flex items-center gap-2 rounded-md border border-custom-border-200 bg-custom-background-100 px-3 py-1.5 text-sm text-custom-text-100 hover:bg-custom-background-90 transition-colors w-full"
            >
              <component
                :is="RELATION_TYPE_CONFIG[relationType].icon"
                class="h-3.5 w-3.5 flex-shrink-0"
                :style="{ color: RELATION_TYPE_CONFIG[relationType].color }"
              />
              <span class="flex-1 text-left">{{ RELATION_TYPE_CONFIG[relationType].label }}</span>
              <ChevronDown class="h-3.5 w-3.5 text-custom-text-300 flex-shrink-0" />
            </button>
          </template>
          <template #default="{ close }">
            <button
              v-for="rt in relationTypes"
              :key="rt"
              class="flex w-full items-center gap-2 rounded-md px-3 py-2 text-sm hover:bg-custom-background-80 transition-colors"
              :class="rt === relationType ? 'bg-custom-background-80 text-custom-text-100' : 'text-custom-text-200'"
              @click="relationType = rt; close()"
            >
              <component
                :is="RELATION_TYPE_CONFIG[rt].icon"
                class="h-3.5 w-3.5 flex-shrink-0"
                :style="{ color: RELATION_TYPE_CONFIG[rt].color }"
              />
              <span>{{ RELATION_TYPE_CONFIG[rt].label }}</span>
            </button>
          </template>
        </PDropdown>
      </div>
      <div>
        <label class="mb-1 block text-xs font-medium text-custom-text-300">Issue ID</label>
        <input
          v-model="relatedIssueId"
          type="text"
          placeholder="Paste issue UUID"
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
          :disabled="!relatedIssueId.trim()"
          class="rounded-md bg-custom-primary-100 px-3 py-1.5 text-xs text-white hover:bg-custom-primary-200 transition-colors disabled:opacity-50"
        >
          Add
        </button>
      </div>
    </div>

    <!-- Relations list -->
    <div v-if="props.relations.length > 0" class="space-y-1">
      <div
        v-for="rel in props.relations"
        :key="rel.id"
        class="group flex items-center gap-2 rounded-md border border-custom-border-200 px-3 py-2 text-sm"
      >
        <component
          :is="RELATION_TYPE_CONFIG[rel.relation_type]?.icon"
          class="h-3.5 w-3.5 flex-shrink-0"
          :style="{ color: RELATION_TYPE_CONFIG[rel.relation_type]?.color }"
        />
        <span class="text-xs text-custom-text-300">{{ RELATION_TYPE_CONFIG[rel.relation_type]?.label }}</span>
        <router-link
          :to="`/${props.workspaceSlug}/projects/${props.projectId}/issues/${rel.related_issue_id}`"
          class="flex items-center gap-1.5 hover:text-custom-primary-100 transition-colors flex-1 min-w-0"
        >
          <span class="font-mono text-xs text-custom-text-300">
            {{ props.projectIdentifier }}-{{ rel.related_issue_sequence_id }}
          </span>
          <span class="text-custom-text-100 truncate">{{ rel.related_issue_name }}</span>
        </router-link>
        <button
          @click="emit('remove-relation', rel.id)"
          class="flex-shrink-0 rounded p-1 text-custom-text-300 opacity-0 group-hover:opacity-100 hover:bg-red-50 hover:text-red-500 transition-all"
        >
          <X class="h-3 w-3" />
        </button>
      </div>
    </div>
    <p v-else-if="!showForm" class="text-sm text-custom-text-300 italic">No relations</p>
  </div>
</template>
