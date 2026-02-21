import type { AxiosError } from 'axios'

export function extractErrorMessage(error: unknown, fallback = 'Something went wrong'): string {
  if (error && typeof error === 'object') {
    const axiosErr = error as AxiosError<{ message?: string; error?: string }>
    const data = axiosErr.response?.data
    if (data?.message) return data.message
    if (data?.error) return data.error
    if (axiosErr.message) return axiosErr.message
  }
  if (error instanceof Error) return error.message
  return fallback
}
