<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth.store'
import PInput from '@/components/ui/PInput.vue'
import PButton from '@/components/ui/PButton.vue'

const router = useRouter()
const authStore = useAuthStore()

const firstName = ref('')
const lastName = ref('')
const email = ref('')
const password = ref('')
const error = ref('')

async function handleSignup() {
  error.value = ''
  try {
    await authStore.signup({
      email: email.value,
      password: password.value,
      first_name: firstName.value,
      last_name: lastName.value,
    })
    router.push('/onboarding/workspace')
  } catch (e: any) {
    error.value = e.response?.data?.message || 'Signup failed'
  }
}
</script>

<template>
  <form @submit.prevent="handleSignup" class="space-y-5">
    <h2 class="text-xl font-semibold text-custom-text-100">Create an account</h2>

    <div v-if="error" class="rounded-lg border border-red-200 bg-red-50 p-3 text-sm text-red-600">
      {{ error }}
    </div>

    <div class="grid grid-cols-2 gap-3">
      <div>
        <label class="mb-1.5 block text-sm font-medium text-custom-text-200">First name</label>
        <PInput v-model="firstName" type="text" placeholder="John" />
      </div>
      <div>
        <label class="mb-1.5 block text-sm font-medium text-custom-text-200">Last name</label>
        <PInput v-model="lastName" type="text" placeholder="Doe" />
      </div>
    </div>

    <div>
      <label class="mb-1.5 block text-sm font-medium text-custom-text-200">Email</label>
      <PInput v-model="email" type="email" placeholder="you@example.com" />
    </div>

    <div>
      <label class="mb-1.5 block text-sm font-medium text-custom-text-200">Password</label>
      <PInput v-model="password" type="password" placeholder="Min. 8 characters" />
    </div>

    <PButton
      type="submit"
      variant="primary"
      :loading="authStore.loading"
      :disabled="authStore.loading"
      class="w-full"
    >
      {{ authStore.loading ? 'Creating account...' : 'Create account' }}
    </PButton>

    <p class="text-center text-sm text-custom-text-300">
      Already have an account?
      <router-link to="/auth/login" class="font-medium text-brand-600 hover:text-brand-700">
        Sign in
      </router-link>
    </p>
  </form>
</template>
