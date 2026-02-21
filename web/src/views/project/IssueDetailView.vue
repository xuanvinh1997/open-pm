<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useProjectStore } from '@/stores/project.store'
import { issueApi } from '@/api/issue.api'
import type { Issue, IssueComment, IssuePriority } from '@/types/issue.types'
import PBreadcrumb from '@/components/ui/PBreadcrumb.vue'
import PSpinner from '@/components/ui/PSpinner.vue'
import PButton from '@/components/ui/PButton.vue'
import IssueDetailSidebar from '@/components/issues/IssueDetailSidebar.vue'
import IssueActivityFeed from '@/components/issues/IssueActivityFeed.vue'
import IssueCommentInput from '@/components/issues/IssueCommentInput.vue'
import { Briefcase, ArrowLeft } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const projectStore = useProjectStore()

const slug = route.params.workspaceSlug as string
const projectId = route.params.projectId as string
const issueId = route.params.issueId as string

const issue = ref<Issue | null>(null)
const comments = ref<IssueComment[]>([])
const loading = ref(false)

const breadcrumbs = computed(() => [
  { label: 'Projects', to: `/${slug}/projects`, icon: Briefcase },
  { label: projectStore.currentProject?.name || 'Project', to: `/${slug}/projects/${projectId}/issues` },
  { label: issue.value ? `${projectStore.currentProject?.identifier}-${issue.value.sequence_id}` : 'Issue' },
])

onMounted(async () => {
  loading.value = true
  try {
    await projectStore.setCurrentProject(slug, projectId)
    await projectStore.fetchStates(slug, projectId)
    await projectStore.fetchLabels(slug, projectId)

    const { data: issueData } = await issueApi.get(slug, projectId, issueId)
    issue.value = issueData

    const { data: commentData } = await issueApi.listComments(slug, projectId, issueId)
    comments.value = commentData.results
  } finally {
    loading.value = false
  }
})

async function handleUpdateState(stateId: string) {
  if (!issue.value) return
  try {
    const { data } = await issueApi.update(slug, projectId, issueId, { state_id: stateId })
    issue.value = data
  } catch (e) {
    console.error('Failed to update state', e)
  }
}

async function handleUpdatePriority(priority: IssuePriority) {
  if (!issue.value) return
  try {
    const { data } = await issueApi.update(slug, projectId, issueId, { priority })
    issue.value = data
  } catch (e) {
    console.error('Failed to update priority', e)
  }
}

async function handleAddComment(html: string) {
  try {
    const { data } = await issueApi.createComment(slug, projectId, issueId, {
      comment_html: html,
    })
    comments.value.push(data)
  } catch (e) {
    console.error('Failed to add comment', e)
  }
}
</script>

<template>
  <div class="flex h-full flex-col">
    <!-- Header -->
    <div class="flex items-center gap-3 border-b border-custom-border-200 bg-custom-background-100 px-6 py-3">
      <button
        @click="router.back()"
        class="rounded-md p-1 text-custom-text-300 hover:bg-custom-background-80 hover:text-custom-text-100 transition-colors"
      >
        <ArrowLeft class="h-4 w-4" />
      </button>
      <PBreadcrumb :items="breadcrumbs" />
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex flex-1 items-center justify-center">
      <PSpinner size="lg" />
    </div>

    <!-- Content -->
    <div v-else-if="issue" class="flex flex-1 overflow-hidden">
      <!-- Main content -->
      <div class="flex-1 overflow-y-auto p-6">
        <!-- Title -->
        <h1 class="mb-4 text-xl font-semibold text-custom-text-100">
          {{ issue.name }}
        </h1>

        <!-- Description -->
        <div class="mb-8">
          <div
            v-if="issue.description_html"
            class="prose prose-sm max-w-none text-custom-text-200"
            v-html="issue.description_html"
          />
          <p v-else class="text-sm text-custom-text-300 italic">
            No description provided
          </p>
        </div>

        <!-- Sub-issues -->
        <div v-if="issue.sub_issues && issue.sub_issues.length > 0" class="mb-8">
          <h3 class="mb-3 text-sm font-semibold text-custom-text-100">Sub-issues</h3>
          <div class="space-y-1">
            <router-link
              v-for="sub in issue.sub_issues"
              :key="sub.id"
              :to="`/${slug}/projects/${projectId}/issues/${sub.id}`"
              class="flex items-center gap-2 rounded-md border border-custom-border-200 px-3 py-2 text-sm hover:bg-custom-background-80 transition-colors"
            >
              <span class="font-mono text-xs text-custom-text-300">
                {{ projectStore.currentProject?.identifier }}-{{ sub.sequence_id }}
              </span>
              <span class="text-custom-text-100">{{ sub.name }}</span>
            </router-link>
          </div>
        </div>

        <!-- Activity -->
        <IssueActivityFeed :comments="comments" />
        <IssueCommentInput @submit="handleAddComment" class="mt-4" />
      </div>

      <!-- Properties sidebar -->
      <div class="w-[300px] flex-shrink-0 border-l border-custom-border-200 overflow-y-auto p-5">
        <IssueDetailSidebar
          :issue="issue"
          :states="projectStore.states"
          :labels="projectStore.labels"
          @update:state="handleUpdateState"
          @update:priority="handleUpdatePriority"
        />
      </div>
    </div>
  </div>
</template>
