import { defineStore } from 'pinia'
import { ref } from 'vue'
import { moduleApi } from '@/api/module.api'
import type { Module, CreateModuleRequest, UpdateModuleRequest } from '@/types/module.types'
import type { Issue } from '@/types/issue.types'

export const useModuleStore = defineStore('module', () => {
  const modules = ref<Module[]>([])
  const currentModule = ref<Module | null>(null)
  const moduleIssues = ref<Issue[]>([])
  const totalIssues = ref(0)
  const loading = ref(false)

  async function fetchModules(slug: string, projectId: string) {
    loading.value = true
    try {
      const { data } = await moduleApi.list(slug, projectId)
      modules.value = data.results
    } finally {
      loading.value = false
    }
  }

  async function fetchModule(slug: string, projectId: string, moduleId: string) {
    const { data } = await moduleApi.get(slug, projectId, moduleId)
    currentModule.value = data.module
    moduleIssues.value = data.issues
    totalIssues.value = data.total_issues
    return data
  }

  async function createModule(slug: string, projectId: string, data: CreateModuleRequest) {
    const { data: mod } = await moduleApi.create(slug, projectId, data)
    modules.value.push(mod)
    return mod
  }

  async function updateModule(slug: string, projectId: string, moduleId: string, data: UpdateModuleRequest) {
    const { data: updated } = await moduleApi.update(slug, projectId, moduleId, data)
    const idx = modules.value.findIndex((m) => m.id === moduleId)
    if (idx >= 0) modules.value[idx] = updated
    if (currentModule.value?.id === moduleId) currentModule.value = updated
    return updated
  }

  async function deleteModule(slug: string, projectId: string, moduleId: string) {
    await moduleApi.delete(slug, projectId, moduleId)
    modules.value = modules.value.filter((m) => m.id !== moduleId)
    if (currentModule.value?.id === moduleId) currentModule.value = null
  }

  async function addIssueToModule(slug: string, projectId: string, moduleId: string, issueId: string) {
    await moduleApi.addIssue(slug, projectId, moduleId, issueId)
  }

  async function removeIssueFromModule(slug: string, projectId: string, moduleId: string, issueId: string) {
    await moduleApi.removeIssue(slug, projectId, moduleId, issueId)
    moduleIssues.value = moduleIssues.value.filter((i) => i.id !== issueId)
    totalIssues.value = Math.max(0, totalIssues.value - 1)
  }

  return {
    modules,
    currentModule,
    moduleIssues,
    totalIssues,
    loading,
    fetchModules,
    fetchModule,
    createModule,
    updateModule,
    deleteModule,
    addIssueToModule,
    removeIssueFromModule,
  }
})
