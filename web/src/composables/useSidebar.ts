import { useStorage } from '@vueuse/core'
import { computed } from 'vue'

const isCollapsed = useStorage('sidebar-collapsed', false)

export function useSidebar() {
  const toggle = () => {
    isCollapsed.value = !isCollapsed.value
  }

  const collapse = () => {
    isCollapsed.value = true
  }

  const expand = () => {
    isCollapsed.value = false
  }

  return {
    isCollapsed: computed(() => isCollapsed.value),
    toggle,
    collapse,
    expand,
  }
}
