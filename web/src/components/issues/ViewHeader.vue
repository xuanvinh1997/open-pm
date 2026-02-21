<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useWorkspaceStore } from '@/stores/workspace.store'
import { useProjectStore } from '@/stores/project.store'
import PBreadcrumb from '@/components/ui/PBreadcrumb.vue'
import PButton from '@/components/ui/PButton.vue'
import { LayoutList, Kanban, Plus, Briefcase } from 'lucide-vue-next'

interface Props {
  activeView: 'list' | 'board'
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'create': []
}>()

const route = useRoute()
const workspaceStore = useWorkspaceStore()
const projectStore = useProjectStore()

const slug = computed(() => route.params.workspaceSlug as string)
const projectId = computed(() => route.params.projectId as string)

const breadcrumbs = computed(() => [
  { label: 'Projects', to: `/${slug.value}/projects`, icon: Briefcase },
  { label: projectStore.currentProject?.name || 'Project' },
  { label: props.activeView === 'list' ? 'Issues' : 'Board' },
])
</script>

<template>
  <div class="flex items-center justify-between border-b border-custom-border-200 bg-custom-background-100 px-6 py-3">
    <div class="flex items-center gap-4">
      <PBreadcrumb :items="breadcrumbs" />
    </div>
    <div class="flex items-center gap-2">
      <!-- View switcher -->
      <div class="flex items-center rounded-md border border-custom-border-200 p-0.5">
        <router-link
          :to="`/${slug}/projects/${projectId}/issues`"
          class="rounded px-2 py-1 transition-colors"
          :class="props.activeView === 'list' ? 'bg-custom-background-80 text-custom-text-100' : 'text-custom-text-300 hover:text-custom-text-200'"
        >
          <LayoutList class="h-4 w-4" />
        </router-link>
        <router-link
          :to="`/${slug}/projects/${projectId}/board`"
          class="rounded px-2 py-1 transition-colors"
          :class="props.activeView === 'board' ? 'bg-custom-background-80 text-custom-text-100' : 'text-custom-text-300 hover:text-custom-text-200'"
        >
          <Kanban class="h-4 w-4" />
        </router-link>
      </div>

      <PButton variant="primary" size="sm" @click="emit('create')">
        <Plus class="h-4 w-4" />
        New issue
      </PButton>
    </div>
  </div>
</template>
