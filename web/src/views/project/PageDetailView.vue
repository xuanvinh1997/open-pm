<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useProjectStore } from '@/stores/project.store'
import { usePageStore } from '@/stores/page.store'
import PBreadcrumb from '@/components/ui/PBreadcrumb.vue'
import PSpinner from '@/components/ui/PSpinner.vue'
import PButton from '@/components/ui/PButton.vue'
import { RichTextEditor } from '@/components/editor'
import { useToast } from '@/composables/useToast'
import { extractErrorMessage } from '@/utils/api-error'
import { Lock, Unlock, Trash2 } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const projectStore = useProjectStore()
const pageStore = usePageStore()
const toast = useToast()

const slug = route.params.workspaceSlug as string
const projectId = route.params.projectId as string
const pageId = route.params.pageId as string

const loading = ref(false)
const saving = ref(false)
const editName = ref('')
const editHtml = ref('')
const editJson = ref<Record<string, unknown>>({})
const editStripped = ref('')

const breadcrumbs = computed(() => [
  { label: projectStore.currentProject?.name || 'Project', to: `/${slug}/projects` },
  { label: 'Pages', to: `/${slug}/projects/${projectId}/pages` },
  { label: pageStore.currentPage?.name || 'Page' },
])

onMounted(async () => {
  loading.value = true
  try {
    await projectStore.setCurrentProject(slug, projectId)
    const page = await pageStore.fetchPage(slug, projectId, pageId)
    editName.value = page.name
    editHtml.value = page.description_html || ''
    editJson.value = page.description_json || {}
    editStripped.value = page.description_stripped || ''
  } finally {
    loading.value = false
  }
})

async function handleSave() {
  saving.value = true
  try {
    await pageStore.updatePage(slug, projectId, pageId, {
      name: editName.value,
      description_html: editHtml.value,
      description_json: editJson.value,
      description_stripped: editStripped.value,
    })
    toast.success('Page saved')
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to save page'))
  } finally {
    saving.value = false
  }
}

async function toggleLock() {
  if (!pageStore.currentPage) return
  try {
    await pageStore.updatePage(slug, projectId, pageId, {
      is_locked: !pageStore.currentPage.is_locked,
    })
    toast.success(pageStore.currentPage.is_locked ? 'Page locked' : 'Page unlocked')
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to toggle lock'))
  }
}

async function handleDelete() {
  if (!confirm('Are you sure you want to delete this page?')) return
  try {
    await pageStore.deletePage(slug, projectId, pageId)
    toast.success('Page deleted')
    router.push(`/${slug}/projects/${projectId}/pages`)
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to delete page'))
  }
}
</script>

<template>
  <div class="flex h-full flex-col">
    <!-- Header -->
    <div class="flex items-center justify-between border-b border-custom-border-200 px-4 py-3">
      <PBreadcrumb :items="breadcrumbs" />
      <div class="flex items-center gap-2">
        <PButton variant="ghost" size="sm" @click="toggleLock" :title="pageStore.currentPage?.is_locked ? 'Unlock page' : 'Lock page'">
          <Lock v-if="pageStore.currentPage?.is_locked" class="h-4 w-4" />
          <Unlock v-else class="h-4 w-4" />
        </PButton>
        <PButton variant="ghost" size="sm" @click="handleDelete">
          <Trash2 class="h-4 w-4 text-red-500" />
        </PButton>
        <PButton variant="primary" size="sm" :loading="saving" @click="handleSave">
          Save
        </PButton>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex flex-1 items-center justify-center">
      <PSpinner size="lg" />
    </div>

    <!-- Editor -->
    <div v-else class="flex-1 overflow-y-auto p-6">
      <div class="mx-auto max-w-3xl">
        <input
          v-model="editName"
          class="mb-4 w-full border-none bg-transparent text-2xl font-bold text-custom-text-100 placeholder-custom-text-300 focus:outline-none"
          placeholder="Untitled page"
        />
        <RichTextEditor
          v-model="editHtml"
          :json="editJson"
          toolbar="full"
          placeholder="Start writing..."
          min-height="400px"
          @update:json="(v) => editJson = v"
          @update:stripped="(v) => editStripped = v"
        />
      </div>
    </div>
  </div>
</template>
