<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useWorkspaceStore } from '@/stores/workspace.store'
import { useAuthStore } from '@/stores/auth.store'
import { ROLE_MEMBER, ROLE_LABELS } from '@/utils/roles'
import PBreadcrumb from '@/components/ui/PBreadcrumb.vue'
import PAvatar from '@/components/ui/PAvatar.vue'
import PBadge from '@/components/ui/PBadge.vue'
import PButton from '@/components/ui/PButton.vue'
import PModal from '@/components/ui/PModal.vue'
import PInput from '@/components/ui/PInput.vue'
import PRoleSelector from '@/components/ui/PRoleSelector.vue'
import { useToast } from '@/composables/useToast'
import { extractErrorMessage } from '@/utils/api-error'
import { Users, Mail, Trash2 } from 'lucide-vue-next'

const router = useRouter()
const route = useRoute()
const workspaceStore = useWorkspaceStore()
const authStore = useAuthStore()

const toast = useToast()

const slug = computed(() => route.params.workspaceSlug as string)

// Invite modal state
const showInviteModal = ref(false)
const inviteEmail = ref('')
const inviteRole = ref(ROLE_MEMBER)
const inviteMessage = ref('')
const inviteLoading = ref(false)
const inviteError = ref('')

// Remove confirmation
const confirmRemoveUserId = ref<string | null>(null)
const removeLoading = ref(false)

onMounted(async () => {
  if (!workspaceStore.isAdmin) {
    router.replace(`/${slug.value}`)
    return
  }
  await Promise.all([
    workspaceStore.fetchMembers(slug.value),
    workspaceStore.fetchInvites(slug.value),
  ])
})

const sortedMembers = computed(() => {
  return [...workspaceStore.members].sort((a, b) => b.role - a.role)
})

function isCurrentUser(userId: string) {
  return userId === authStore.user?.id
}

function isOwner(userId: string) {
  return userId === workspaceStore.currentWorkspace?.owner_id
}

function memberDisplayName(member: { display_name: string; first_name: string; last_name: string; email?: string }) {
  if (member.display_name) return member.display_name
  if (member.first_name) return `${member.first_name} ${member.last_name || ''}`.trim()
  return member.email || ''
}

async function handleRoleChange(userId: string, role: number) {
  try {
    await workspaceStore.updateMemberRole(slug.value, userId, role)
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to update member role'))
    await workspaceStore.fetchMembers(slug.value)
  }
}

async function handleRemoveMember() {
  if (!confirmRemoveUserId.value) return
  removeLoading.value = true
  try {
    await workspaceStore.removeMember(slug.value, confirmRemoveUserId.value)
    confirmRemoveUserId.value = null
    toast.success('Member removed')
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to remove member'))
  } finally {
    removeLoading.value = false
  }
}

async function handleInvite() {
  inviteError.value = ''
  if (!inviteEmail.value.trim()) {
    inviteError.value = 'Email is required'
    return
  }
  inviteLoading.value = true
  try {
    await workspaceStore.inviteMember(slug.value, {
      email: inviteEmail.value.trim(),
      role: inviteRole.value,
      message: inviteMessage.value.trim() || undefined,
    })
    showInviteModal.value = false
    inviteEmail.value = ''
    inviteRole.value = ROLE_MEMBER
    inviteMessage.value = ''
    toast.success('Invite sent')
  } catch (e: any) {
    inviteError.value = extractErrorMessage(e, 'Failed to send invite')
  } finally {
    inviteLoading.value = false
  }
}

function roleBadgeColor(role: number): string {
  if (role >= 25) return '#F59E0B'
  if (role >= 20) return '#8B5CF6'
  if (role >= 15) return '#3B82F6'
  return '#6B7280'
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}
</script>

<template>
  <div v-if="workspaceStore.isAdmin" class="h-full overflow-y-auto">
    <!-- Header -->
    <div class="border-b border-custom-border-200 bg-custom-background-100 px-6 py-4">
      <div class="flex items-center justify-between">
        <PBreadcrumb :items="[{ label: 'Members', icon: Users }]" />
        <PButton variant="primary" size="sm" @click="showInviteModal = true">
          <Mail class="h-4 w-4" />
          Invite member
        </PButton>
      </div>
    </div>

    <div class="p-6">
      <div class="max-w-4xl">
        <!-- Members table -->
        <div class="rounded-xl border border-custom-border-200 bg-custom-background-100 overflow-hidden">
          <div class="border-b border-custom-border-200 px-6 py-3">
            <h2 class="text-sm font-semibold text-custom-text-100">
              Members ({{ workspaceStore.members.length }})
            </h2>
          </div>
          <div class="divide-y divide-custom-border-200">
            <div
              v-for="member in sortedMembers"
              :key="member.id"
              class="flex items-center gap-4 px-6 py-3"
            >
              <PAvatar
                :name="memberDisplayName(member)"
                :src="member.avatar_url"
                size="md"
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
                  v-if="!isOwner(member.user_id) && !isCurrentUser(member.user_id)"
                  :model-value="member.role"
                  :max-role="workspaceStore.currentUserRole"
                  @update:model-value="handleRoleChange(member.user_id, $event)"
                />
                <PBadge v-else :color="roleBadgeColor(member.role)" size="sm">
                  {{ ROLE_LABELS[member.role] || 'Unknown' }}
                </PBadge>
                <button
                  v-if="!isOwner(member.user_id) && !isCurrentUser(member.user_id)"
                  class="rounded-md p-1.5 text-custom-text-300 hover:bg-red-50 hover:text-red-500 transition-colors"
                  title="Remove member"
                  @click="confirmRemoveUserId = member.user_id"
                >
                  <Trash2 class="h-4 w-4" />
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Pending Invites -->
        <div v-if="workspaceStore.invites.length > 0" class="mt-6 rounded-xl border border-custom-border-200 bg-custom-background-100 overflow-hidden">
          <div class="border-b border-custom-border-200 px-6 py-3">
            <h2 class="text-sm font-semibold text-custom-text-100">
              Pending Invites ({{ workspaceStore.invites.length }})
            </h2>
          </div>
          <div class="divide-y divide-custom-border-200">
            <div
              v-for="invite in workspaceStore.invites"
              :key="invite.id"
              class="flex items-center gap-4 px-6 py-3"
            >
              <div class="flex h-10 w-10 items-center justify-center rounded-full bg-custom-background-80">
                <Mail class="h-4 w-4 text-custom-text-300" />
              </div>
              <div class="flex-1 min-w-0">
                <div class="text-sm font-medium text-custom-text-100 truncate">{{ invite.email }}</div>
                <div class="text-xs text-custom-text-300">Invited {{ formatDate(invite.created_at) }}</div>
              </div>
              <PBadge :color="roleBadgeColor(invite.role)" size="sm">
                {{ ROLE_LABELS[invite.role] || 'Unknown' }}
              </PBadge>
              <PBadge color="#F59E0B" variant="outline" size="sm">Pending</PBadge>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Invite Modal -->
    <PModal :open="showInviteModal" @update:open="showInviteModal = $event" title="Invite member" size="sm">
      <form @submit.prevent="handleInvite" class="space-y-4">
        <div>
          <label class="mb-1.5 block text-sm font-medium text-custom-text-200">Email address</label>
          <PInput
            v-model="inviteEmail"
            type="email"
            placeholder="name@example.com"
            :error="inviteError"
          />
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
              class="flex-1 rounded-md border px-3 py-2 text-sm font-medium transition-colors"
              :class="inviteRole === opt.value
                ? 'border-brand-500 bg-brand-50 text-brand-600 dark:bg-brand-500/10'
                : 'border-custom-border-200 text-custom-text-200 hover:bg-custom-background-80'"
              @click="inviteRole = opt.value"
            >
              {{ opt.label }}
            </button>
          </div>
        </div>
        <div>
          <label class="mb-1.5 block text-sm font-medium text-custom-text-200">Message (optional)</label>
          <textarea
            v-model="inviteMessage"
            rows="3"
            placeholder="Add a personal message..."
            class="w-full rounded-md border border-custom-border-200 bg-custom-background-100 px-3 py-2 text-sm text-custom-text-100 placeholder:text-custom-text-400 focus:border-brand-500 focus:ring-1 focus:ring-brand-500 focus:outline-none"
          />
        </div>
      </form>
      <template #footer>
        <PButton variant="secondary" size="sm" @click="showInviteModal = false">Cancel</PButton>
        <PButton variant="primary" size="sm" :loading="inviteLoading" @click="handleInvite">
          Send invite
        </PButton>
      </template>
    </PModal>

    <!-- Remove confirmation modal -->
    <PModal :open="!!confirmRemoveUserId" @update:open="confirmRemoveUserId = null" title="Remove member" size="sm">
      <p class="text-sm text-custom-text-200">
        Are you sure you want to remove this member from the workspace? They will lose access to all projects.
      </p>
      <template #footer>
        <PButton variant="secondary" size="sm" @click="confirmRemoveUserId = null">Cancel</PButton>
        <PButton variant="danger" size="sm" :loading="removeLoading" @click="handleRemoveMember">
          Remove
        </PButton>
      </template>
    </PModal>
  </div>
</template>
