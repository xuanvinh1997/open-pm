import { useStorage } from '@vueuse/core'
import { computed, watchEffect } from 'vue'

type Theme = 'light' | 'dark'

const theme = useStorage<Theme>('theme', 'light')

export function useTheme() {
  watchEffect(() => {
    const root = document.documentElement
    if (theme.value === 'dark') {
      root.classList.add('dark')
    } else {
      root.classList.remove('dark')
    }
  })

  const toggleTheme = () => {
    theme.value = theme.value === 'light' ? 'dark' : 'light'
  }

  const setTheme = (t: Theme) => {
    theme.value = t
  }

  return {
    theme: computed(() => theme.value),
    toggleTheme,
    setTheme,
  }
}
