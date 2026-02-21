<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useWorkspaceStore } from '@/stores/workspace.store'
import { useAuthStore } from '@/stores/auth.store'
import { useSidebar } from '@/composables/useSidebar'
import { useMediaQuery } from '@vueuse/core'
import Sidebar from './Sidebar.vue'
import PCommandPalette from '@/components/ui/PCommandPalette.vue'

const route = useRoute()
const workspaceStore = useWorkspaceStore()
const authStore = useAuthStore()
const { isCollapsed, collapse } = useSidebar()
const isMobile = useMediaQuery('(max-width: 768px)')
const mobileMenuOpen = ref(false)

watch(isMobile, (mobile) => {
  if (mobile) {
    collapse()
    mobileMenuOpen.value = false
  }
})

watch(() => route.path, () => {
  if (isMobile.value) {
    mobileMenuOpen.value = false
  }
})

onMounted(async () => {
  if (!authStore.user) {
    await authStore.fetchUser()
  }
  await workspaceStore.fetchWorkspaces()
  const slug = route.params.workspaceSlug as string
  if (slug) {
    await workspaceStore.setCurrentWorkspace(slug)
    await workspaceStore.fetchMembers(slug)
  }
})
</script>

<template>
  <div class="flex h-screen overflow-hidden bg-custom-background-100">
    <!-- Desktop sidebar -->
    <aside
      class="hidden md:block flex-shrink-0 transition-all duration-200 ease-in-out"
      :class="isCollapsed ? 'w-[70px]' : 'w-[248px]'"
    >
      <Sidebar />
    </aside>

    <!-- Mobile sidebar overlay -->
    <Transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="mobileMenuOpen"
        class="fixed inset-0 z-40 bg-black/50 md:hidden"
        @click="mobileMenuOpen = false"
      />
    </Transition>
    <Transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="-translate-x-full"
      enter-to-class="translate-x-0"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="translate-x-0"
      leave-to-class="-translate-x-full"
    >
      <aside
        v-if="mobileMenuOpen"
        class="fixed inset-y-0 left-0 z-50 w-[248px] md:hidden"
      >
        <Sidebar />
      </aside>
    </Transition>

    <!-- Mobile header bar -->
    <div class="flex flex-1 flex-col overflow-hidden">
      <div class="flex h-12 items-center gap-3 border-b border-custom-border-200 px-4 md:hidden">
        <button
          @click="mobileMenuOpen = true"
          class="rounded-md p-1.5 text-custom-text-200 hover:bg-custom-background-80 transition-colors"
        >
          <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
          </svg>
        </button>
        <span class="text-sm font-semibold text-custom-text-100">{{ workspaceStore.currentWorkspace?.name || 'Open-PM' }}</span>
      </div>

      <main class="flex-1 overflow-hidden">
        <router-view v-slot="{ Component }">
          <Transition
            name="page"
            mode="out-in"
          >
            <component :is="Component" />
          </Transition>
        </router-view>
      </main>
    </div>
    <PCommandPalette />
  </div>
</template>
