<script setup lang="ts">
import { ref, onMounted, toRef } from 'vue'
import { useRoute } from 'vue-router'
import { useWorkspaceStore } from '@/stores/workspace.store'
import { useProjectStore } from '@/stores/project.store'
import { useAuthStore } from '@/stores/auth.store'
import { issueApi } from '@/api/issue.api'
import { projectApi } from '@/api/project.api'
import { useAnalytics } from '@/composables/useAnalytics'
import { useTheme } from '@/composables/useTheme'
import type { Issue } from '@/types/issue.types'
import type { State } from '@/types/project.types'
import PBreadcrumb from '@/components/ui/PBreadcrumb.vue'
import PEmptyState from '@/components/ui/PEmptyState.vue'
import PButton from '@/components/ui/PButton.vue'
import DoughnutChart from '@/components/ui/charts/DoughnutChart.vue'
import BarChart from '@/components/ui/charts/BarChart.vue'
import { useToast } from '@/composables/useToast'
import { extractErrorMessage } from '@/utils/api-error'
import { Home, FolderKanban, Users, LayoutDashboard, Plus, Briefcase } from 'lucide-vue-next'

const route = useRoute()
const workspaceStore = useWorkspaceStore()
const projectStore = useProjectStore()
const authStore = useAuthStore()
const { theme } = useTheme()

const toast = useToast()

const slug = route.params.workspaceSlug as string

const allIssues = ref<Issue[]>([])
const allStates = ref<State[]>([])
const chartsLoading = ref(false)

const { issuesByStateGroup, issuesByPriority, stats } =
  useAnalytics({ issues: allIssues, states: allStates })

onMounted(async () => {
  try {
    await projectStore.fetchProjects(slug)
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to load projects'))
    return
  }

  if (projectStore.projects.length > 0) {
    chartsLoading.value = true
    try {
      const projectsToFetch = projectStore.projects.slice(0, 5)
      const results = await Promise.all(
        projectsToFetch.map(async (p) => {
          const [issueRes, stateRes] = await Promise.all([
            issueApi.list(slug, p.id, 1, 200),
            projectApi.listStates(slug, p.id),
          ])
          return { issues: issueRes.data.results, states: stateRes.data.results }
        }),
      )
      allIssues.value = results.flatMap((r) => r.issues)
      allStates.value = results.flatMap((r) => r.states)
    } catch (e) {
      toast.error(extractErrorMessage(e, 'Failed to load dashboard data'))
    } finally {
      chartsLoading.value = false
    }
  }
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
              <p class="text-2xl font-bold text-custom-text-100">{{ stats.active }}</p>
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

      <!-- Charts -->
      <div v-if="allIssues.length > 0" class="mb-8 grid grid-cols-1 gap-6 md:grid-cols-2" :key="theme">
        <DoughnutChart
          title="Issues by Status"
          :labels="issuesByStateGroup.map(s => s.label)"
          :data="issuesByStateGroup.map(s => s.count)"
          :colors="issuesByStateGroup.map(s => s.color)"
        />
        <BarChart
          title="Issues by Priority"
          :labels="issuesByPriority.map(p => p.label)"
          :data="issuesByPriority.map(p => p.count)"
          :colors="issuesByPriority.map(p => p.color)"
        />
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
              <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-custom-background-80 flex-shrink-0">
                <Briefcase class="h-5 w-5 text-custom-text-300" />
              </div>
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
