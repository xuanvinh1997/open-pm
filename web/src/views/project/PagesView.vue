<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useProjectStore } from '@/stores/project.store'
import { usePageStore } from '@/stores/page.store'
import type { PageTreeNode } from '@/stores/page.store'
import CreatePageModal from '@/components/pages/CreatePageModal.vue'
import PBreadcrumb from '@/components/ui/PBreadcrumb.vue'
import PEmptyState from '@/components/ui/PEmptyState.vue'
import PSpinner from '@/components/ui/PSpinner.vue'
import PButton from '@/components/ui/PButton.vue'
import { useToast } from '@/composables/useToast'
import { extractErrorMessage } from '@/utils/api-error'
import { formatRelativeDate } from '@/utils/helpers'
import { FileText, Plus, ChevronRight, Lock } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const projectStore = useProjectStore()
const pageStore = usePageStore()
const toast = useToast()

const slug = route.params.workspaceSlug as string
const projectId = route.params.projectId as string

const loading = ref(false)
const showCreateModal = ref(false)
const expandedPages = ref<Set<string>>(new Set())

const breadcrumbs = computed(() => [
  { label: projectStore.currentProject?.name || 'Project', to: `/${slug}/projects` },
  { label: 'Pages' },
])

onMounted(async () => {
  loading.value = true
  try {
    await projectStore.setCurrentProject(slug, projectId)
    await pageStore.fetchPages(slug, projectId)
  } finally {
    loading.value = false
  }
})

async function handleCreatePage(data: { name: string; parent_id?: string; color?: string }) {
  try {
    await pageStore.createPage(slug, projectId, data)
    showCreateModal.value = false
    toast.success('Page created')
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to create page'))
  }
}

function handlePageClick(pageId: string) {
  router.push(`/${slug}/projects/${projectId}/pages/${pageId}`)
}

function toggleExpand(pageId: string) {
  if (expandedPages.value.has(pageId)) {
    expandedPages.value.delete(pageId)
  } else {
    expandedPages.value.add(pageId)
  }
}
</script>

<template>
  <div class="flex h-full flex-col">
    <!-- Header -->
    <div class="flex items-center justify-between border-b border-custom-border-200 px-4 py-3">
      <div class="flex items-center gap-3">
        <PBreadcrumb :items="breadcrumbs" />
      </div>
      <PButton variant="primary" size="sm" @click="showCreateModal = true">
        <Plus class="h-4 w-4" />
        New page
      </PButton>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex flex-1 items-center justify-center">
      <PSpinner size="lg" />
    </div>

    <!-- Empty state -->
    <div v-else-if="pageStore.pages.length === 0" class="flex-1">
      <PEmptyState
        title="No pages"
        description="Create your first page to start documenting."
        :icon="FileText"
      >
        <PButton variant="primary" @click="showCreateModal = true">
          <Plus class="h-4 w-4" />
          New page
        </PButton>
      </PEmptyState>
    </div>

    <!-- Pages tree list -->
    <div v-else class="flex-1 overflow-y-auto p-4">
      <div class="space-y-0.5">
        <template v-for="node in pageStore.treePages" :key="node.id">
          <div
            class="flex items-center gap-2 rounded-md px-3 py-2 hover:bg-custom-background-90 transition-colors cursor-pointer"
            @click="handlePageClick(node.id)"
          >
            <button
              v-if="node.children.length > 0"
              class="p-0.5"
              @click.stop="toggleExpand(node.id)"
            >
              <ChevronRight
                class="h-3.5 w-3.5 text-custom-text-300 transition-transform duration-100"
                :class="{ 'rotate-90': expandedPages.has(node.id) }"
              />
            </button>
            <span v-else class="w-5" />
            <FileText class="h-4 w-4 text-custom-text-300 flex-shrink-0" />
            <span class="flex-1 text-sm text-custom-text-100 truncate">{{ node.name }}</span>
            <Lock v-if="node.is_locked" class="h-3.5 w-3.5 text-custom-text-300" />
            <span class="text-2xs text-custom-text-300">{{ formatRelativeDate(node.updated_at) }}</span>
          </div>
          <!-- Children (one level deep for now) -->
          <template v-if="expandedPages.has(node.id)">
            <div
              v-for="child in node.children"
              :key="child.id"
              class="ml-7 flex items-center gap-2 rounded-md px-3 py-2 hover:bg-custom-background-90 transition-colors cursor-pointer"
              @click="handlePageClick(child.id)"
            >
              <FileText class="h-4 w-4 text-custom-text-300 flex-shrink-0" />
              <span class="flex-1 text-sm text-custom-text-100 truncate">{{ child.name }}</span>
              <Lock v-if="child.is_locked" class="h-3.5 w-3.5 text-custom-text-300" />
              <span class="text-2xs text-custom-text-300">{{ formatRelativeDate(child.updated_at) }}</span>
            </div>
          </template>
        </template>
      </div>
    </div>

    <CreatePageModal
      v-model:open="showCreateModal"
      :pages="pageStore.pages"
      @create="handleCreatePage"
    />
  </div>
</template>
