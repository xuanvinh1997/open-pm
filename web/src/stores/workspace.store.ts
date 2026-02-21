import { defineStore } from 'pinia'
import { ref } from 'vue'
import { workspaceApi } from '@/api/workspace.api'
import type { Workspace, WorkspaceMember } from '@/types/workspace.types'

export const useWorkspaceStore = defineStore('workspace', () => {
  const workspaces = ref<Workspace[]>([])
  const currentWorkspace = ref<Workspace | null>(null)
  const members = ref<WorkspaceMember[]>([])
  const loading = ref(false)

  async function fetchWorkspaces() {
    loading.value = true
    try {
      const { data } = await workspaceApi.list()
      workspaces.value = data.results
    } finally {
      loading.value = false
    }
  }

  async function setCurrentWorkspace(slug: string) {
    const existing = workspaces.value.find((w) => w.slug === slug)
    if (existing) {
      currentWorkspace.value = existing
    } else {
      const { data } = await workspaceApi.get(slug)
      currentWorkspace.value = data
    }
  }

  async function fetchMembers(slug: string) {
    const { data } = await workspaceApi.listMembers(slug)
    members.value = data.results
  }

  async function createWorkspace(name: string, slug: string) {
    const { data } = await workspaceApi.create({ name, slug })
    workspaces.value.push(data)
    currentWorkspace.value = data
    return data
  }

  return {
    workspaces,
    currentWorkspace,
    members,
    loading,
    fetchWorkspaces,
    setCurrentWorkspace,
    fetchMembers,
    createWorkspace,
  }
})
