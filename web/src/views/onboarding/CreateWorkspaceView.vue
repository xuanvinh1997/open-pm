<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useWorkspaceStore } from '@/stores/workspace.store'
import PInput from '@/components/ui/PInput.vue'
import PButton from '@/components/ui/PButton.vue'

const router = useRouter()
const workspaceStore = useWorkspaceStore()

const name = ref('')
const slug = ref('')
const error = ref('')

watch(name, (val) => {
  slug.value = val
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-|-$/g, '')
})

async function handleCreate() {
  error.value = ''
  try {
    const workspace = await workspaceStore.createWorkspace(name.value, slug.value)
    router.push(`/${workspace.slug}`)
  } catch (e: any) {
    error.value = e.response?.data?.message || 'Failed to create workspace'
  }
}
</script>

<template>
  <div class="flex min-h-screen items-center justify-center bg-custom-background-80 px-4">
    <div class="w-full max-w-[440px]">
      <div class="mb-10 text-center">
        <div class="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-xl bg-brand-600">
          <span class="text-xl font-bold text-white">O</span>
        </div>
        <h1 class="text-2xl font-bold text-custom-text-100">Create your workspace</h1>
        <p class="mt-1 text-sm text-custom-text-300">A workspace is where your team collaborates on projects.</p>
      </div>

      <form
        @submit.prevent="handleCreate"
        class="space-y-5 rounded-xl border border-custom-border-200 bg-custom-background-100 p-8 shadow-custom-sm"
      >
        <div v-if="error" class="rounded-lg border border-red-200 bg-red-50 p-3 text-sm text-red-600">
          {{ error }}
        </div>

        <div>
          <label class="mb-1.5 block text-sm font-medium text-custom-text-200">Workspace name</label>
          <PInput v-model="name" type="text" placeholder="My Team" />
        </div>

        <div>
          <label class="mb-1.5 block text-sm font-medium text-custom-text-200">URL slug</label>
          <div class="flex items-center">
            <span class="rounded-l-md border border-r-0 border-custom-border-200 bg-custom-background-80 px-3 py-2 text-sm text-custom-text-300">
              open-pm.dev/
            </span>
            <input
              v-model="slug"
              type="text"
              required
              maxlength="48"
              pattern="[a-z0-9][a-z0-9-]*[a-z0-9]"
              class="flex-1 rounded-r-md border border-custom-border-200 bg-custom-background-100 px-3 py-2 text-sm text-custom-text-100 placeholder:text-custom-text-300 outline-none focus:border-brand-500 focus:ring-1 focus:ring-brand-500"
              placeholder="my-team"
            />
          </div>
        </div>

        <PButton type="submit" variant="primary" class="w-full">
          Create workspace
        </PButton>
      </form>
    </div>
  </div>
</template>
