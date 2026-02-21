<script setup lang="ts">
import { useWorkspaceStore } from '@/stores/workspace.store'
import PBreadcrumb from '@/components/ui/PBreadcrumb.vue'
import PInput from '@/components/ui/PInput.vue'
import { Settings } from 'lucide-vue-next'

const workspaceStore = useWorkspaceStore()
</script>

<template>
  <div class="h-full overflow-y-auto">
    <!-- Header -->
    <div class="border-b border-custom-border-200 bg-custom-background-100 px-6 py-4">
      <PBreadcrumb :items="[{ label: 'Settings', icon: Settings }]" />
    </div>

    <div class="p-6">
      <h1 class="mb-6 text-xl font-semibold text-custom-text-100">Workspace Settings</h1>

      <div class="max-w-2xl space-y-6">
        <div class="rounded-xl border border-custom-border-200 bg-custom-background-100 p-6">
          <h2 class="mb-4 text-base font-semibold text-custom-text-100">General</h2>
          <div class="space-y-4">
            <div>
              <label class="mb-1.5 block text-sm font-medium text-custom-text-200">Name</label>
              <PInput
                :model-value="workspaceStore.currentWorkspace?.name || ''"
                type="text"
                disabled
              />
            </div>
            <div>
              <label class="mb-1.5 block text-sm font-medium text-custom-text-200">Slug</label>
              <PInput
                :model-value="workspaceStore.currentWorkspace?.slug || ''"
                type="text"
                disabled
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
