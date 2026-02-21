<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useProjectStore } from '@/stores/project.store'
import { useWorkspaceStore } from '@/stores/workspace.store'
import { useAuthStore } from '@/stores/auth.store'
import { projectApi } from '@/api/project.api'
import { ROLE_MEMBER, ROLE_LABELS } from '@/utils/roles'
import PModal from '@/components/ui/PModal.vue'
import PAvatar from '@/components/ui/PAvatar.vue'
import PBadge from '@/components/ui/PBadge.vue'
import PButton from '@/components/ui/PButton.vue'
import PRoleSelector from '@/components/ui/PRoleSelector.vue'
import { useToast } from '@/composables/useToast'
import { extractErrorMessage } from '@/utils/api-error'
import { Trash2, UserPlus } from 'lucide-vue-next'

interface Props {
  open: boolean
}

defineProps<Props>()
const emit = defineEmits<{
  'update:open': [value: boolean]
}>()

const route = useRoute()
const projectStore = useProjectStore()
const workspaceStore = useWorkspaceStore()
const authStore = useAuthStore()

const toast = useToast()

const slug = computed(() => route.params.workspaceSlug as string)
const projectId = computed(() => route.params.projectId as string)

// Add member state
const showAddMember = ref(false)
const selectedUserId = ref('')
const selectedRole = ref(ROLE_MEMBER)
const addLoading = ref(false)

// Remove state
const confirmRemoveUserId = ref<string | null>(null)
const removeLoading = ref(false)

const sortedMembers = computed(() => {
  return [...projectStore.members].sort((a, b) => b.role - a.role)
})

const availableWorkspaceMembers = computed(() => {
  const projectMemberIds = new Set(projectStore.members.map(m => m.user_id))
  return workspaceStore.members.filter(m => !projectMemberIds.has(m.user_id))
})

watch(() => selectedUserId.value, () => {
  // Reset when user changes
})

function isCurrentUser(userId: string) {
  return userId === authStore.user?.id
}

function memberDisplayName(member: { display_name: string; first_name: string; last_name: string; email?: string }) {
  if (member.display_name) return member.display_name
  if (member.first_name) return `${member.first_name} ${member.last_name || ''}`.trim()
  return member.email || ''
}

function roleBadgeColor(role: number): string {
  if (role >= 25) return '#F59E0B'
  if (role >= 20) return '#8B5CF6'
  if (role >= 15) return '#3B82F6'
  return '#6B7280'
}

async function handleRoleChange(userId: string, role: number) {
  try {
    await projectApi.updateMemberRole(slug.value, projectId.value, userId, role)
    const member = projectStore.members.find(m => m.user_id === userId)
    if (member) member.role = role
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to update member role'))
    await projectStore.fetchMembers(slug.value, projectId.value)
  }
}

async function handleAddMember() {
  if (!selectedUserId.value) return
  addLoading.value = true
  try {
    await projectApi.addMember(slug.value, projectId.value, {
      user_id: selectedUserId.value,
      role: selectedRole.value,
    })
    await projectStore.fetchMembers(slug.value, projectId.value)
    showAddMember.value = false
    selectedUserId.value = ''
    selectedRole.value = ROLE_MEMBER
    toast.success('Member added')
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to add member'))
  } finally {
    addLoading.value = false
  }
}

async function handleRemoveMember() {
  if (!confirmRemoveUserId.value) return
  removeLoading.value = true
  try {
    await projectApi.removeMember(slug.value, projectId.value, confirmRemoveUserId.value)
    projectStore.members = projectStore.members.filter(m => m.user_id !== confirmRemoveUserId.value)
    confirmRemoveUserId.value = null
    toast.success('Member removed')
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to remove member'))
  } finally {
    removeLoading.value = false
  }
}
</script>

<template>
  <PModal :open="open" @update:open="emit('update:open', $event)" title="Project Members" size="lg">
    <div class="space-y-4">
      <!-- Member list -->
      <div class="divide-y divide-custom-border-200 rounded-lg border border-custom-border-200">
        <div
          v-for="member in sortedMembers"
          :key="member.id"
          class="flex items-center gap-3 px-4 py-2.5"
        >
          <PAvatar
            :name="memberDisplayName(member)"
            :src="member.avatar_url"
            size="sm"
          />
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2">
              <span class="text-sm font-medium text-custom-text-100 truncate">
                {{ memberDisplayName(member) }}
              </span>
              <span v-if="isCurrentUser(member.user_id)" class="text-2xs text-custom-text-300">(You)</span>
            </div>
            <div class="text-xs text-custom-text-300 truncate">{{ member.email }}</div>
          </div>
          <div class="flex items-center gap-2">
            <PRoleSelector
              v-if="projectStore.isProjectAdmin && !isCurrentUser(member.user_id)"
              :model-value="member.role"
              :max-role="projectStore.currentProjectRole"
              @update:model-value="handleRoleChange(member.user_id, $event)"
            />
            <PBadge v-else :color="roleBadgeColor(member.role)" size="sm">
              {{ ROLE_LABELS[member.role] || 'Unknown' }}
            </PBadge>
            <button
              v-if="projectStore.isProjectAdmin && !isCurrentUser(member.user_id)"
              class="rounded-md p-1 text-custom-text-300 hover:bg-red-50 hover:text-red-500 transition-colors"
              title="Remove member"
              @click="confirmRemoveUserId = member.user_id"
            >
              <Trash2 class="h-3.5 w-3.5" />
            </button>
          </div>
        </div>
        <div v-if="sortedMembers.length === 0" class="px-4 py-6 text-center text-sm text-custom-text-300">
          No members yet
        </div>
      </div>

      <!-- Add member section (admins only) -->
      <div v-if="projectStore.isProjectAdmin">
        <button
          v-if="!showAddMember"
          class="flex items-center gap-2 text-sm text-brand-500 hover:text-brand-600 transition-colors"
          @click="showAddMember = true"
        >
          <UserPlus class="h-4 w-4" />
          Add workspace member
        </button>
        <div v-else class="rounded-lg border border-custom-border-200 p-4 space-y-3">
          <div>
            <label class="mb-1.5 block text-sm font-medium text-custom-text-200">Select member</label>
            <select
              v-model="selectedUserId"
              class="w-full rounded-md border border-custom-border-200 bg-custom-background-100 px-3 py-2 text-sm text-custom-text-100 focus:border-brand-500 focus:ring-1 focus:ring-brand-500 focus:outline-none"
            >
              <option value="" disabled>Choose a workspace member...</option>
              <option
                v-for="wm in availableWorkspaceMembers"
                :key="wm.user_id"
                :value="wm.user_id"
              >
                {{ memberDisplayName(wm) }} ({{ wm.email }})
              </option>
            </select>
          </div>
          <div>
            <label class="mb-1.5 block text-sm font-medium text-custom-text-200">Role</label>
            <div class="flex gap-2">
              <button
                v-for="opt in [
                  { value: 5, label: 'Guest' },
                  { value: 15, label: 'Member' },
                  { value: 20, label: 'Admin' },
                ]"
                :key="opt.value"
                type="button"
                class="flex-1 rounded-md border px-3 py-1.5 text-sm font-medium transition-colors"
                :class="selectedRole === opt.value
                  ? 'border-brand-500 bg-brand-50 text-brand-600 dark:bg-brand-500/10'
                  : 'border-custom-border-200 text-custom-text-200 hover:bg-custom-background-80'"
                @click="selectedRole = opt.value"
              >
                {{ opt.label }}
              </button>
            </div>
          </div>
          <div class="flex justify-end gap-2">
            <PButton variant="secondary" size="sm" @click="showAddMember = false; selectedUserId = ''; selectedRole = ROLE_MEMBER">
              Cancel
            </PButton>
            <PButton variant="primary" size="sm" :loading="addLoading" :disabled="!selectedUserId" @click="handleAddMember">
              Add member
            </PButton>
          </div>
        </div>
      </div>
    </div>

    <!-- Remove confirmation -->
    <PModal :open="!!confirmRemoveUserId" @update:open="confirmRemoveUserId = null" title="Remove member" size="sm">
      <p class="text-sm text-custom-text-200">
        Are you sure you want to remove this member from the project?
      </p>
      <template #footer>
        <PButton variant="secondary" size="sm" @click="confirmRemoveUserId = null">Cancel</PButton>
        <PButton variant="danger" size="sm" :loading="removeLoading" @click="handleRemoveMember">
          Remove
        </PButton>
      </template>
    </PModal>
  </PModal>
</template>
