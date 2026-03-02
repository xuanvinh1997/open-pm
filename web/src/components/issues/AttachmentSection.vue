<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { assetApi } from '@/api/asset.api'
import type { FileAsset } from '@/types/asset.types'
import { useToast } from '@/composables/useToast'
import { extractErrorMessage } from '@/utils/api-error'
import { Paperclip, Upload, Trash2, Download, FileIcon } from 'lucide-vue-next'

const props = defineProps<{
  slug: string
  entityType: string
  entityId: string
}>()

const toast = useToast()
const assets = ref<FileAsset[]>([])
const uploading = ref(false)
const dragOver = ref(false)

onMounted(async () => {
  try {
    const { data } = await assetApi.list(props.slug, props.entityType, props.entityId)
    assets.value = data.results
  } catch {
    // silent fail — attachments are secondary
  }
})

function formatFileSize(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}

async function handleUpload(files: FileList | null) {
  if (!files || files.length === 0) return
  uploading.value = true
  try {
    for (const file of Array.from(files)) {
      if (file.size > 10 * 1024 * 1024) {
        toast.error(`File "${file.name}" exceeds 10MB limit`)
        continue
      }
      const { data } = await assetApi.upload(props.slug, file, props.entityType, props.entityId)
      assets.value.unshift(data)
    }
    toast.success('File(s) uploaded')
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to upload file'))
  } finally {
    uploading.value = false
  }
}

async function handleDelete(asset: FileAsset) {
  try {
    await assetApi.delete(props.slug, asset.id)
    assets.value = assets.value.filter((a) => a.id !== asset.id)
    toast.success('File deleted')
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to delete file'))
  }
}

function handleDrop(e: DragEvent) {
  e.preventDefault()
  dragOver.value = false
  handleUpload(e.dataTransfer?.files ?? null)
}

function handleFileInput(e: Event) {
  const input = e.target as HTMLInputElement
  handleUpload(input.files)
  input.value = ''
}
</script>

<template>
  <div class="space-y-3">
    <div class="flex items-center gap-2">
      <Paperclip class="h-4 w-4 text-custom-text-300" />
      <h3 class="text-sm font-medium text-custom-text-100">Attachments</h3>
      <span class="text-2xs text-custom-text-300">({{ assets.length }})</span>
    </div>

    <!-- Drop zone -->
    <div
      class="relative rounded-lg border-2 border-dashed p-4 text-center transition-colors"
      :class="dragOver
        ? 'border-brand-500 bg-brand-500/5'
        : 'border-custom-border-200 hover:border-custom-border-300'"
      @dragover.prevent="dragOver = true"
      @dragleave="dragOver = false"
      @drop="handleDrop"
    >
      <input
        type="file"
        multiple
        class="absolute inset-0 cursor-pointer opacity-0"
        @change="handleFileInput"
      />
      <Upload class="mx-auto h-5 w-5 text-custom-text-300" />
      <p class="mt-1 text-xs text-custom-text-300">
        {{ uploading ? 'Uploading...' : 'Drop files or click to upload (max 10MB)' }}
      </p>
    </div>

    <!-- Attachment list -->
    <div v-if="assets.length > 0" class="space-y-1">
      <div
        v-for="asset in assets"
        :key="asset.id"
        class="flex items-center gap-2 rounded-md border border-custom-border-200 px-3 py-2 text-sm"
      >
        <FileIcon class="h-4 w-4 flex-shrink-0 text-custom-text-300" />
        <span class="flex-1 truncate text-custom-text-200">{{ asset.file_name }}</span>
        <span class="text-2xs text-custom-text-300 flex-shrink-0">{{ formatFileSize(asset.file_size) }}</span>
        <a
          v-if="asset.download_url"
          :href="asset.download_url"
          target="_blank"
          class="rounded p-1 text-custom-text-300 hover:bg-custom-background-80 hover:text-custom-text-200 transition-colors"
        >
          <Download class="h-3.5 w-3.5" />
        </a>
        <button
          class="rounded p-1 text-custom-text-300 hover:bg-red-50 hover:text-red-600 transition-colors"
          @click="handleDelete(asset)"
        >
          <Trash2 class="h-3.5 w-3.5" />
        </button>
      </div>
    </div>
  </div>
</template>
