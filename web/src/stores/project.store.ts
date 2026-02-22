import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { projectApi } from '@/api/project.api'
import { useAuthStore } from '@/stores/auth.store'
import { ROLE_ADMIN, ROLE_MEMBER } from '@/utils/roles'
import type { Project, State, Label, ProjectMember } from '@/types/project.types'
import type { EstimateSystem } from '@/types/issue.types'

export const useProjectStore = defineStore('project', () => {
  const projects = ref<Project[]>([])
  const currentProject = ref<Project | null>(null)
  const states = ref<State[]>([])
  const labels = ref<Label[]>([])
  const members = ref<ProjectMember[]>([])
  const estimateSystem = ref<EstimateSystem | null>(null)
  const currentProjectRole = ref<number>(0)
  const loading = ref(false)

  const isProjectAdmin = computed(() => currentProjectRole.value >= ROLE_ADMIN)
  const isProjectMember = computed(() => currentProjectRole.value >= ROLE_MEMBER)

  async function fetchProjects(slug: string) {
    loading.value = true
    try {
      const { data } = await projectApi.list(slug)
      projects.value = data.results
    } finally {
      loading.value = false
    }
  }

  async function setCurrentProject(slug: string, projectId: string) {
    const existing = projects.value.find((p) => p.id === projectId)
    if (existing) {
      currentProject.value = existing
    } else {
      const { data } = await projectApi.get(slug, projectId)
      currentProject.value = data
    }
  }

  async function fetchMembers(slug: string, projectId: string) {
    const { data } = await projectApi.listMembers(slug, projectId)
    members.value = data.results
    // Identify current user's role in this project
    const authStore = useAuthStore()
    const me = data.results.find((m: ProjectMember) => m.user_id === authStore.user?.id)
    currentProjectRole.value = me?.role ?? 0
  }

  async function fetchStates(slug: string, projectId: string) {
    const { data } = await projectApi.listStates(slug, projectId)
    states.value = data.results
  }

  async function fetchLabels(slug: string, projectId: string) {
    const { data } = await projectApi.listLabels(slug, projectId)
    labels.value = data.results
  }

  async function fetchEstimateSystem(slug: string, projectId: string) {
    try {
      const { data } = await projectApi.getEstimateSystem(slug, projectId)
      estimateSystem.value = data
    } catch {
      estimateSystem.value = null
    }
  }

  async function createProject(slug: string, name: string, identifier: string, description?: string) {
    const { data } = await projectApi.create(slug, { name, identifier, description })
    projects.value.push(data)
    return data
  }

  return {
    projects,
    currentProject,
    states,
    labels,
    members,
    estimateSystem,
    currentProjectRole,
    loading,
    isProjectAdmin,
    isProjectMember,
    fetchProjects,
    setCurrentProject,
    fetchMembers,
    fetchStates,
    fetchLabels,
    fetchEstimateSystem,
    createProject,
  }
})
