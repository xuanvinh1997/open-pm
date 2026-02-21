import { ref } from 'vue'
import { useMagicKeys, whenever } from '@vueuse/core'

const isOpen = ref(false)

export function useCommandPalette() {
  const keys = useMagicKeys()

  whenever(keys['Meta+k']!, () => {
    isOpen.value = !isOpen.value
  })

  whenever(keys['Ctrl+k']!, () => {
    isOpen.value = !isOpen.value
  })

  const open = () => {
    isOpen.value = true
  }

  const close = () => {
    isOpen.value = false
  }

  const toggle = () => {
    isOpen.value = !isOpen.value
  }

  return {
    isOpen,
    open,
    close,
    toggle,
  }
}
