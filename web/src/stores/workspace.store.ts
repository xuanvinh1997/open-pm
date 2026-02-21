import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { workspaceApi } from '@/api/workspace.api'
import { useAuthStore } from '@/stores/auth.store'
import { ROLE_ADMIN, ROLE_OWNER } from '@/utils/roles'
import type { Workspace, WorkspaceMember, WorkspaceInvite } from '@/types/workspace.types'

export const useWorkspaceStore = defineStore('workspace', () => {
  const workspaces = ref<Workspace[]>([])
  const currentWorkspace = ref<Workspace | null>(null)
  const members = ref<WorkspaceMember[]>([])
  const invites = ref<WorkspaceInvite[]>([])
  const currentUserRole = ref<number>(0)
  const loading = ref(false)

  const isAdmin = computed(() => currentUserRole.value >= ROLE_ADMIN)
  const isOwner = computed(() => currentUserRole.value >= ROLE_OWNER)

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
    const authStore = useAuthStore()
    const me = data.results.find((m: WorkspaceMember) => m.user_id === authStore.user?.id)
    currentUserRole.value = me?.role ?? 0
  }

  async function fetchInvites(slug: string) {
    const { data } = await workspaceApi.listInvites(slug)
    invites.value = data.results
  }

  async function inviteMember(slug: string, data: { email: string; role: number; message?: string }) {
    await workspaceApi.invite(slug, data)
    await fetchInvites(slug)
  }

  async function updateMemberRole(slug: string, userId: string, role: number) {
    await workspaceApi.updateMemberRole(slug, userId, role)
    const member = members.value.find(m => m.user_id === userId)
    if (member) member.role = role
  }

  async function removeMember(slug: string, userId: string) {
    await workspaceApi.removeMember(slug, userId)
    members.value = members.value.filter(m => m.user_id !== userId)
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
    invites,
    currentUserRole,
    loading,
    isAdmin,
    isOwner,
    fetchWorkspaces,
    setCurrentWorkspace,
    fetchMembers,
    fetchInvites,
    inviteMember,
    updateMemberRole,
    removeMember,
    createWorkspace,
  }
})
