<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth.store'
import { useWorkspaceStore } from '@/stores/workspace.store'
import PInput from '@/components/ui/PInput.vue'
import PButton from '@/components/ui/PButton.vue'

const router = useRouter()
const authStore = useAuthStore()
const workspaceStore = useWorkspaceStore()

const email = ref('')
const password = ref('')
const error = ref('')

async function handleLogin() {
  error.value = ''
  try {
    await authStore.login({ email: email.value, password: password.value })
    await workspaceStore.fetchWorkspaces()
    if (workspaceStore.workspaces.length > 0) {
      router.push(`/${workspaceStore.workspaces[0].slug}`)
    } else {
      router.push('/onboarding/workspace')
    }
  } catch (e: any) {
    error.value = e.response?.data?.message || 'Login failed'
  }
}
</script>

<template>
  <form @submit.prevent="handleLogin" class="space-y-5">
    <h2 class="text-xl font-semibold text-custom-text-100">Sign in</h2>

    <div v-if="error" class="rounded-lg border border-red-200 bg-red-50 p-3 text-sm text-red-600">
      {{ error }}
    </div>

    <div>
      <label class="mb-1.5 block text-sm font-medium text-custom-text-200">Email</label>
      <PInput
        v-model="email"
        type="email"
        placeholder="you@example.com"
      />
    </div>

    <div>
      <label class="mb-1.5 block text-sm font-medium text-custom-text-200">Password</label>
      <PInput
        v-model="password"
        type="password"
        placeholder="Enter your password"
      />
    </div>

    <PButton
      type="submit"
      variant="primary"
      :loading="authStore.loading"
      :disabled="authStore.loading"
      class="w-full"
    >
      {{ authStore.loading ? 'Signing in...' : 'Sign in' }}
    </PButton>

    <p class="text-center text-sm text-custom-text-300">
      Don't have an account?
      <router-link to="/auth/signup" class="font-medium text-brand-600 hover:text-brand-700">
        Sign up
      </router-link>
    </p>
  </form>
</template>
