<script setup lang="ts">
import { ref, onMounted, computed, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useProjectStore } from '@/stores/project.store'
import { useSprintStore } from '@/stores/sprint.store'
import { issueApi } from '@/api/issue.api'
import type { Issue, IssueComment, IssuePriority, IssueType, WorkLog, CreateWorkLogRequest, UpdateWorkLogRequest, IssueRelation, RelationType, IssueLink } from '@/types/issue.types'
import PBreadcrumb from '@/components/ui/PBreadcrumb.vue'
import PSpinner from '@/components/ui/PSpinner.vue'
import IssueDetailSidebar from '@/components/issues/IssueDetailSidebar.vue'
import IssueActivityFeed from '@/components/issues/IssueActivityFeed.vue'
import IssueCommentInput from '@/components/issues/IssueCommentInput.vue'
import IssueRelationsSection from '@/components/issues/IssueRelationsSection.vue'
import IssueLinksSection from '@/components/issues/IssueLinksSection.vue'
import AttachmentSection from '@/components/issues/AttachmentSection.vue'
import LogWorkModal from '@/components/issues/LogWorkModal.vue'
import { RichTextEditor, RichTextDisplay } from '@/components/editor'
import { useToast } from '@/composables/useToast'
import { useAuthStore } from '@/stores/auth.store'
import { extractErrorMessage } from '@/utils/api-error'
import PButton from '@/components/ui/PButton.vue'
import PModal from '@/components/ui/PModal.vue'
import { Briefcase, ArrowLeft, Trash2 } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const projectStore = useProjectStore()
const sprintStore = useSprintStore()

const slug = route.params.workspaceSlug as string
const projectId = route.params.projectId as string
const issueId = route.params.issueId as string

const toast = useToast()
const authStore = useAuthStore()

const issue = ref<Issue | null>(null)
const comments = ref<IssueComment[]>([])
const workLogs = ref<WorkLog[]>([])
const relations = ref<IssueRelation[]>([])
const links = ref<IssueLink[]>([])
const epics = ref<Issue[]>([])
const totalMinutes = ref(0)
const loading = ref(false)

const showLogWorkModal = ref(false)
const editingWorkLog = ref<WorkLog | null>(null)

// Inline editing
const editingTitle = ref(false)
const editingDescription = ref(false)
const titleInput = ref('')
const descriptionHtml = ref('')
const descriptionJson = ref<Record<string, unknown>>({})
const titleInputRef = ref<HTMLInputElement>()

function startEditTitle() {
  if (!issue.value) return
  titleInput.value = issue.value.name
  editingTitle.value = true
  nextTick(() => titleInputRef.value?.focus())
}

function startEditDescription() {
  if (!issue.value) return
  descriptionHtml.value = issue.value.description_html || ''
  descriptionJson.value = issue.value.description_json || {}
  editingDescription.value = true
}

async function saveTitle() {
  editingTitle.value = false
  if (!issue.value || !titleInput.value.trim() || titleInput.value.trim() === issue.value.name) return
  try {
    const { data } = await issueApi.update(slug, projectId, issueId, { name: titleInput.value.trim() })
    issue.value = data
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to update title'))
  }
}

async function saveDescription() {
  editingDescription.value = false
  if (!issue.value) return
  if (descriptionHtml.value === (issue.value.description_html || '')) return
  try {
    const { data } = await issueApi.update(slug, projectId, issueId, {
      description_html: descriptionHtml.value,
      description_json: descriptionJson.value,
    } as any)
    issue.value = data
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to update description'))
  }
}

function cancelEditTitle() {
  editingTitle.value = false
}

function cancelEditDescription() {
  editingDescription.value = false
}

// Delete issue
const showDeleteConfirm = ref(false)
const deleting = ref(false)

async function handleDeleteIssue() {
  deleting.value = true
  try {
    await issueApi.delete(slug, projectId, issueId)
    toast.success('Issue deleted')
    router.push(`/${slug}/projects/${projectId}/issues`)
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to delete issue'))
  } finally {
    deleting.value = false
    showDeleteConfirm.value = false
  }
}

const breadcrumbs = computed(() => [
  { label: 'Projects', to: `/${slug}/projects`, icon: Briefcase },
  { label: projectStore.currentProject?.name || 'Project', to: `/${slug}/projects/${projectId}/issues` },
  { label: issue.value ? `${projectStore.currentProject?.identifier}-${issue.value.sequence_id}` : 'Issue' },
])

onMounted(async () => {
  loading.value = true
  try {
    await projectStore.setCurrentProject(slug, projectId)
    await Promise.all([
      projectStore.fetchStates(slug, projectId),
      projectStore.fetchLabels(slug, projectId),
      projectStore.fetchMembers(slug, projectId),
      projectStore.fetchEstimateSystem(slug, projectId),
      sprintStore.fetchSprints(slug, projectId),
    ])

    const { data: issueData } = await issueApi.get(slug, projectId, issueId)
    issue.value = issueData

    const { data: commentData } = await issueApi.listComments(slug, projectId, issueId)
    comments.value = commentData.results

    await Promise.all([fetchWorkLogs(), fetchRelations(), fetchLinks(), fetchEpics()])
  } finally {
    loading.value = false
  }
})

async function fetchWorkLogs() {
  try {
    const { data } = await issueApi.listWorkLogs(slug, projectId, issueId)
    workLogs.value = data.results
    totalMinutes.value = data.total_minutes
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to load work logs'))
  }
}

async function fetchRelations() {
  try {
    const { data } = await issueApi.listRelations(slug, projectId, issueId)
    relations.value = data.results
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to load relations'))
  }
}

async function fetchLinks() {
  try {
    const { data } = await issueApi.listLinks(slug, projectId, issueId)
    links.value = data.results
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to load links'))
  }
}

async function fetchEpics() {
  try {
    const { data } = await issueApi.list(slug, projectId, 1, 200, { type: 'epic' })
    epics.value = data.results
  } catch (e) {
    // Non-critical — just no epics available
  }
}

async function handleUpdateState(stateId: string) {
  if (!issue.value) return
  try {
    const { data } = await issueApi.update(slug, projectId, issueId, { state_id: stateId })
    issue.value = data
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to update state'))
  }
}

async function handleUpdatePriority(priority: IssuePriority) {
  if (!issue.value) return
  try {
    const { data } = await issueApi.update(slug, projectId, issueId, { priority })
    issue.value = data
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to update priority'))
  }
}

async function handleUpdateAssignees(assigneeIds: string[]) {
  if (!issue.value) return
  try {
    await issueApi.update(slug, projectId, issueId, { assignee_ids: assigneeIds } as any)
    // Re-fetch to get enriched assignees
    const { data } = await issueApi.get(slug, projectId, issueId)
    issue.value = data
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to update assignees'))
  }
}

async function handleUpdateLabels(labelIds: string[]) {
  if (!issue.value) return
  try {
    await issueApi.update(slug, projectId, issueId, { label_ids: labelIds } as any)
    const { data } = await issueApi.get(slug, projectId, issueId)
    issue.value = data
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to update labels'))
  }
}

async function handleUpdateEstimatePoint(value: number | undefined) {
  if (!issue.value) return
  try {
    const { data } = await issueApi.update(slug, projectId, issueId, { estimate_point: value ?? null } as any)
    issue.value = data
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to update estimate'))
  }
}

async function handleUpdateType(issueType: IssueType) {
  if (!issue.value) return
  try {
    const { data } = await issueApi.update(slug, projectId, issueId, { issue_type: issueType })
    issue.value = data
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to update type'))
  }
}

async function handleUpdateStartDate(date: string) {
  if (!issue.value) return
  try {
    const { data } = await issueApi.update(slug, projectId, issueId, { start_date: date || undefined })
    issue.value = data
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to update start date'))
  }
}

async function handleUpdateTargetDate(date: string) {
  if (!issue.value) return
  try {
    const { data } = await issueApi.update(slug, projectId, issueId, { target_date: date || undefined })
    issue.value = data
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to update target date'))
  }
}

async function handleAddComment(html: string, json?: Record<string, unknown>, stripped?: string) {
  try {
    const { data } = await issueApi.createComment(slug, projectId, issueId, {
      comment_html: html,
      comment_json: json,
      comment_stripped: stripped,
    })
    comments.value.push(data)
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to add comment'))
  }
}

async function handleUpdateComment(commentId: string, html: string, json?: Record<string, unknown>, stripped?: string) {
  try {
    const { data } = await issueApi.updateComment(slug, projectId, issueId, commentId, {
      comment_html: html,
      comment_json: json,
      comment_stripped: stripped,
    })
    const idx = comments.value.findIndex((c) => c.id === commentId)
    if (idx !== -1) comments.value[idx] = data
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to update comment'))
  }
}

async function handleDeleteComment(commentId: string) {
  try {
    await issueApi.deleteComment(slug, projectId, issueId, commentId)
    comments.value = comments.value.filter((c) => c.id !== commentId)
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to delete comment'))
  }
}

async function handleAddRelation(data: { related_issue_id: string; relation_type: RelationType }) {
  try {
    await issueApi.addRelation(slug, projectId, issueId, data)
    await fetchRelations()
    toast.success('Relation added')
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to add relation'))
  }
}

async function handleRemoveRelation(relationId: string) {
  try {
    await issueApi.removeRelation(slug, projectId, issueId, relationId)
    await fetchRelations()
    toast.success('Relation removed')
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to remove relation'))
  }
}

async function handleAddLink(data: { title: string; url: string }) {
  try {
    await issueApi.createLink(slug, projectId, issueId, data)
    await fetchLinks()
    toast.success('Link added')
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to add link'))
  }
}

async function handleRemoveLink(id: string) {
  try {
    await issueApi.deleteLink(slug, projectId, issueId, id)
    await fetchLinks()
    toast.success('Link removed')
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to remove link'))
  }
}

function handleOpenLogWork() {
  editingWorkLog.value = null
  showLogWorkModal.value = true
}

function handleEditWorkLog(workLog: WorkLog) {
  editingWorkLog.value = workLog
  showLogWorkModal.value = true
}

async function handleCreateWorkLog(data: CreateWorkLogRequest) {
  try {
    await issueApi.createWorkLog(slug, projectId, issueId, data)
    await fetchWorkLogs()
    toast.success('Work logged successfully')
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to log work'))
  }
}

async function handleUpdateWorkLog(id: string, data: UpdateWorkLogRequest) {
  try {
    await issueApi.updateWorkLog(slug, projectId, issueId, id, data)
    await fetchWorkLogs()
    toast.success('Work log updated')
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to update work log'))
  }
}

async function handleDeleteWorkLog(id: string) {
  try {
    await issueApi.deleteWorkLog(slug, projectId, issueId, id)
    await fetchWorkLogs()
    toast.success('Work log deleted')
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to delete work log'))
  }
}

async function handleAddToSprint(sprintId: string) {
  if (!issue.value) return
  try {
    await sprintStore.addIssueToSprint(slug, projectId, sprintId, issueId)
    // Re-fetch issue to get updated sprints list
    const { data } = await issueApi.get(slug, projectId, issueId)
    issue.value = data
    toast.success('Issue added to sprint')
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to add issue to sprint'))
  }
}

async function handleRemoveFromSprint(sprintId: string) {
  if (!issue.value) return
  try {
    await sprintStore.removeIssueFromSprint(slug, projectId, sprintId, issueId)
    // Re-fetch issue to get updated sprints list
    const { data } = await issueApi.get(slug, projectId, issueId)
    issue.value = data
    toast.success('Issue removed from sprint')
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to remove issue from sprint'))
  }
}

async function handleUpdateParent(parentId: string | null) {
  if (!issue.value) return
  try {
    const payload: Record<string, unknown> = parentId
      ? { parent_id: parentId }
      : { clear_parent_id: true }
    const { data } = await issueApi.update(slug, projectId, issueId, payload as any)
    issue.value = data
    toast.success(parentId ? 'Epic assigned' : 'Removed from epic')
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to update epic'))
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
      <div class="ml-auto">
        <button
          @click="showDeleteConfirm = true"
          class="rounded-md p-1.5 text-custom-text-300 hover:bg-red-50 hover:text-red-500 transition-colors"
          title="Delete issue"
        >
          <Trash2 class="h-4 w-4" />
        </button>
      </div>
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
        <div class="mb-4">
          <input
            v-if="editingTitle"
            ref="titleInputRef"
            v-model="titleInput"
            class="w-full text-xl font-semibold text-custom-text-100 bg-transparent border-b-2 border-brand-500 outline-none py-1"
            @blur="saveTitle"
            @keydown.enter="saveTitle"
            @keydown.escape="cancelEditTitle"
          />
          <h1
            v-else
            class="text-xl font-semibold text-custom-text-100 cursor-pointer rounded px-1 -mx-1 hover:bg-custom-background-80 transition-colors"
            @click="startEditTitle"
          >
            {{ issue.name }}
          </h1>
        </div>

        <!-- Description -->
        <div class="mb-8">
          <RichTextEditor
            v-if="editingDescription"
            v-model="descriptionHtml"
            :json="descriptionJson"
            placeholder="Add a description..."
            toolbar="full"
            min-height="100px"
            autofocus
            @update:json="(v) => descriptionJson = v"
            @blur="saveDescription"
            @keydown.escape="cancelEditDescription"
          />
          <div
            v-else-if="issue.description_html"
            class="cursor-pointer rounded px-1 -mx-1 hover:bg-custom-background-80 transition-colors"
            @click="startEditDescription"
          >
            <RichTextDisplay :html="issue.description_html" />
          </div>
          <p
            v-else
            class="text-sm text-custom-text-300 italic cursor-pointer rounded px-1 -mx-1 hover:bg-custom-background-80 transition-colors"
            @click="startEditDescription"
          >
            Click to add a description...
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

        <!-- Relations -->
        <IssueRelationsSection
          :relations="relations"
          :project-identifier="projectStore.currentProject?.identifier || ''"
          :project-id="projectId"
          :workspace-slug="slug"
          @add-relation="handleAddRelation"
          @remove-relation="handleRemoveRelation"
        />

        <!-- Links -->
        <IssueLinksSection
          :links="links"
          @add-link="handleAddLink"
          @remove-link="handleRemoveLink"
        />

        <!-- Attachments -->
        <div class="mb-8">
          <AttachmentSection
            :slug="slug"
            entity-type="issue"
            :entity-id="issueId"
          />
        </div>

        <!-- Activity -->
        <IssueActivityFeed
          :comments="comments"
          :current-user-id="authStore.user?.id"
          @update-comment="handleUpdateComment"
          @delete-comment="handleDeleteComment"
        />
        <IssueCommentInput @submit="handleAddComment" class="mt-4" />
      </div>

      <!-- Properties sidebar -->
      <div class="w-[300px] flex-shrink-0 border-l border-custom-border-200 overflow-y-auto p-5">
        <IssueDetailSidebar
          :issue="issue"
          :states="projectStore.states"
          :labels="projectStore.labels"
          :members="projectStore.members"
          :sprints="sprintStore.sprints"
          :epics="epics"
          :work-logs="workLogs"
          :total-minutes="totalMinutes"
          :estimate-system="projectStore.estimateSystem"
          @update:state="handleUpdateState"
          @update:priority="handleUpdatePriority"
          @update:type="handleUpdateType"
          @update:estimate_point="handleUpdateEstimatePoint"
          @update:assignees="handleUpdateAssignees"
          @update:labels="handleUpdateLabels"
          @update:start_date="handleUpdateStartDate"
          @update:target_date="handleUpdateTargetDate"
          @update:sprint="handleAddToSprint"
          @remove:sprint="handleRemoveFromSprint"
          @update:parent="handleUpdateParent"
          @log-work="handleOpenLogWork"
          @edit-work-log="handleEditWorkLog"
          @delete-work-log="handleDeleteWorkLog"
        />
      </div>
    </div>

    <!-- Log Work Modal -->
    <LogWorkModal
      v-model:open="showLogWorkModal"
      :work-log="editingWorkLog"
      @create="handleCreateWorkLog"
      @update="handleUpdateWorkLog"
    />

    <!-- Delete Confirmation Modal -->
    <PModal v-model:open="showDeleteConfirm" title="Delete issue" size="sm">
      <p class="text-sm text-custom-text-200">
        Are you sure you want to delete
        <span class="font-semibold text-custom-text-100">{{ issue?.name }}</span>?
        This action cannot be undone.
      </p>
      <template #footer>
        <PButton variant="secondary" size="sm" @click="showDeleteConfirm = false">Cancel</PButton>
        <PButton variant="danger" size="sm" :loading="deleting" @click="handleDeleteIssue">Delete</PButton>
      </template>
    </PModal>
  </div>
</template>
