<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useWorkspaceStore } from '@/stores/workspace.store'
import { useProjectStore } from '@/stores/project.store'
import { useAuthStore } from '@/stores/auth.store'
import PBreadcrumb from '@/components/ui/PBreadcrumb.vue'
import PEmptyState from '@/components/ui/PEmptyState.vue'
import PButton from '@/components/ui/PButton.vue'
import { Home, FolderKanban, Users, LayoutDashboard, Plus } from 'lucide-vue-next'

const route = useRoute()
const workspaceStore = useWorkspaceStore()
const projectStore = useProjectStore()
const authStore = useAuthStore()

const slug = route.params.workspaceSlug as string

onMounted(async () => {
  await projectStore.fetchProjects(slug)
})
</script>

<template>
  <div class="h-full overflow-y-auto">
    <!-- Header -->
    <div class="border-b border-custom-border-200 bg-custom-background-100 px-6 py-4">
      <PBreadcrumb :items="[{ label: 'Home', icon: Home }]" />
    </div>

    <div class="p-6">
      <!-- Welcome -->
      <div class="mb-8">
        <h1 class="text-2xl font-semibold text-custom-text-100">
          Welcome back, {{ authStore.user?.first_name || 'there' }}
        </h1>
        <p class="mt-1 text-sm text-custom-text-300">
          {{ workspaceStore.currentWorkspace?.name }} dashboard
        </p>
      </div>

      <!-- Stats -->
      <div class="mb-8 grid grid-cols-3 gap-4">
        <div class="rounded-xl border border-custom-border-200 bg-custom-background-100 p-5 transition-all hover:shadow-custom-sm">
          <div class="flex items-center gap-3">
            <div class="rounded-lg bg-brand-500/10 p-2.5">
              <FolderKanban class="h-5 w-5 text-brand-600" />
            </div>
            <div>
              <p class="text-2xl font-bold text-custom-text-100">{{ projectStore.projects.length }}</p>
              <p class="text-sm text-custom-text-300">Projects</p>
            </div>
          </div>
        </div>
        <div class="rounded-xl border border-custom-border-200 bg-custom-background-100 p-5 transition-all hover:shadow-custom-sm">
          <div class="flex items-center gap-3">
            <div class="rounded-lg bg-green-500/10 p-2.5">
              <LayoutDashboard class="h-5 w-5 text-green-600" />
            </div>
            <div>
              <p class="text-2xl font-bold text-custom-text-100">0</p>
              <p class="text-sm text-custom-text-300">Active issues</p>
            </div>
          </div>
        </div>
        <div class="rounded-xl border border-custom-border-200 bg-custom-background-100 p-5 transition-all hover:shadow-custom-sm">
          <div class="flex items-center gap-3">
            <div class="rounded-lg bg-purple-500/10 p-2.5">
              <Users class="h-5 w-5 text-purple-600" />
            </div>
            <div>
              <p class="text-2xl font-bold text-custom-text-100">{{ workspaceStore.members.length || 1 }}</p>
              <p class="text-sm text-custom-text-300">Members</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Projects -->
      <div>
        <div class="mb-4 flex items-center justify-between">
          <h2 class="text-lg font-semibold text-custom-text-100">Projects</h2>
          <router-link
            :to="`/${slug}/projects`"
            class="text-sm font-medium text-brand-600 hover:text-brand-700"
          >
            View all
          </router-link>
        </div>

        <PEmptyState
          v-if="projectStore.projects.length === 0"
          title="No projects yet"
          description="Create your first project to start tracking work."
          :icon="FolderKanban"
        >
          <router-link :to="`/${slug}/projects`">
            <PButton variant="primary">
              <Plus class="h-4 w-4" />
              Create project
            </PButton>
          </router-link>
        </PEmptyState>

        <div v-else class="grid grid-cols-2 gap-4">
          <router-link
            v-for="project in projectStore.projects"
            :key="project.id"
            :to="`/${slug}/projects/${project.id}/issues`"
            class="group rounded-xl border border-custom-border-200 bg-custom-background-100 p-5 transition-all hover:shadow-custom-sm hover:border-custom-border-300"
          >
            <div class="flex items-center gap-3">
              <span class="text-2xl">{{ project.emoji || '📁' }}</span>
              <div>
                <h3 class="font-medium text-custom-text-100 group-hover:text-brand-600 transition-colors">{{ project.name }}</h3>
                <p class="text-xs text-custom-text-300 font-mono">{{ project.identifier }}</p>
              </div>
            </div>
            <p v-if="project.description" class="mt-3 line-clamp-2 text-sm text-custom-text-300">
              {{ project.description }}
            </p>
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>
