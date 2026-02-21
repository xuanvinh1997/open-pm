import { ref } from 'vue'

export type ToastType = 'success' | 'error' | 'warning' | 'info'

export interface Toast {
  id: number
  message: string
  type: ToastType
  duration: number
}

const toasts = ref<Toast[]>([])
let nextId = 0

function addToast(message: string, type: ToastType = 'info', duration?: number) {
  const id = nextId++
  const defaultDuration = type === 'error' ? 8000 : 5000
  const toast: Toast = { id, message, type, duration: duration ?? defaultDuration }
  toasts.value.push(toast)

  setTimeout(() => {
    removeToast(id)
  }, toast.duration)
}

function removeToast(id: number) {
  const idx = toasts.value.findIndex((t) => t.id === id)
  if (idx !== -1) toasts.value.splice(idx, 1)
}

export function useToast() {
  return {
    toasts,
    addToast,
    removeToast,
    success: (msg: string, duration?: number) => addToast(msg, 'success', duration),
    error: (msg: string, duration?: number) => addToast(msg, 'error', duration),
    warning: (msg: string, duration?: number) => addToast(msg, 'warning', duration),
    info: (msg: string, duration?: number) => addToast(msg, 'info', duration),
  }
}
