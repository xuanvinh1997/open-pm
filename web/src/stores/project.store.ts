import { defineStore } from 'pinia'
import { ref } from 'vue'
import { projectApi } from '@/api/project.api'
import type { Project, State, Label } from '@/types/project.types'

export const useProjectStore = defineStore('project', () => {
  const projects = ref<Project[]>([])
  const currentProject = ref<Project | null>(null)
  const states = ref<State[]>([])
  const labels = ref<Label[]>([])
  const loading = ref(false)

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

  async function fetchStates(slug: string, projectId: string) {
    const { data } = await projectApi.listStates(slug, projectId)
    states.value = data.results
  }

  async function fetchLabels(slug: string, projectId: string) {
    const { data } = await projectApi.listLabels(slug, projectId)
    labels.value = data.results
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
    loading,
    fetchProjects,
    setCurrentProject,
    fetchStates,
    fetchLabels,
    createProject,
  }
})
