<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useProjectStore } from '@/stores/project.store'
import PBreadcrumb from '@/components/ui/PBreadcrumb.vue'
import PButton from '@/components/ui/PButton.vue'
import PEmptyState from '@/components/ui/PEmptyState.vue'
import CreateProjectModal from '@/components/project/CreateProjectModal.vue'
import { Plus, FolderKanban, Briefcase } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const projectStore = useProjectStore()
const slug = route.params.workspaceSlug as string

const showCreateModal = ref(false)
const error = ref('')

onMounted(async () => {
  await projectStore.fetchProjects(slug)
})

async function handleCreate(data: { name: string; identifier: string; description?: string }) {
  error.value = ''
  try {
    const project = await projectStore.createProject(slug, data.name, data.identifier, data.description)
    showCreateModal.value = false
    router.push(`/${slug}/projects/${project.id}/issues`)
  } catch (e: any) {
    error.value = e.response?.data?.message || 'Failed to create project'
  }
}
</script>

<template>
  <div class="h-full overflow-y-auto">
    <!-- Header -->
    <div class="flex items-center justify-between border-b border-custom-border-200 bg-custom-background-100 px-6 py-4">
      <PBreadcrumb :items="[{ label: 'Projects', icon: Briefcase }]" />
      <PButton variant="primary" size="sm" @click="showCreateModal = true">
        <Plus class="h-4 w-4" />
        New project
      </PButton>
    </div>

    <div class="p-6">
      <!-- Empty state -->
      <PEmptyState
        v-if="projectStore.projects.length === 0"
        title="No projects"
        description="Get started by creating your first project."
        :icon="FolderKanban"
      >
        <PButton variant="primary" @click="showCreateModal = true">
          <Plus class="h-4 w-4" />
          New project
        </PButton>
      </PEmptyState>

      <!-- Projects grid -->
      <div v-else class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
        <router-link
          v-for="project in projectStore.projects"
          :key="project.id"
          :to="`/${slug}/projects/${project.id}/issues`"
          class="group rounded-xl border border-custom-border-200 bg-custom-background-100 p-5 transition-all hover:shadow-custom-sm hover:border-custom-border-300"
        >
          <div class="flex items-start gap-3">
            <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-custom-background-80 flex-shrink-0">
              <Briefcase class="h-5 w-5 text-custom-text-300" />
            </div>
            <div class="flex-1 min-w-0">
              <h3 class="font-semibold text-custom-text-100 group-hover:text-brand-600 transition-colors">
                {{ project.name }}
              </h3>
              <p class="mt-0.5 text-xs text-custom-text-300 font-mono">{{ project.identifier }}</p>
              <p v-if="project.description" class="mt-2 line-clamp-2 text-sm text-custom-text-300">
                {{ project.description }}
              </p>
            </div>
          </div>
        </router-link>
      </div>
    </div>

    <!-- Create project modal -->
    <CreateProjectModal
      v-model:open="showCreateModal"
      @create="handleCreate"
    />
  </div>
</template>
