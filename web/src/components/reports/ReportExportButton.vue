<script setup lang="ts">
import { ref } from 'vue'
import { reportApi } from '@/api/report.api'
import PDropdown from '@/components/ui/PDropdown.vue'
import { Download, FileText, Table } from 'lucide-vue-next'

interface Props {
  reportType: string
  slug: string
  projectId: string
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'export-csv': []
}>()

const exporting = ref(false)

async function handleExportPDF() {
  exporting.value = true
  try {
    const response = await reportApi.exportPDF(props.slug, props.projectId, props.reportType)
    const blob = new Blob([response.data], { type: 'application/pdf' })
    const link = document.createElement('a')
    link.href = URL.createObjectURL(blob)
    link.download = `${props.reportType}-report.pdf`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(link.href)
  } catch {
    // PDF export may not be available yet
  } finally {
    exporting.value = false
  }
}
</script>

<template>
  <PDropdown align="right" width="10rem">
    <template #trigger>
      <button
        class="flex items-center gap-1.5 rounded-md border border-custom-border-200 bg-custom-background-100 px-3 py-1.5 text-sm text-custom-text-200 hover:bg-custom-background-80 transition-colors"
        :disabled="exporting"
      >
        <Download class="h-3.5 w-3.5" />
        Export
      </button>
    </template>
    <template #default="{ close }">
      <button
        class="flex w-full items-center gap-2 rounded-md px-3 py-2 text-sm text-custom-text-200 hover:bg-custom-background-80 transition-colors"
        @click="emit('export-csv'); close()"
      >
        <Table class="h-3.5 w-3.5" />
        Export CSV
      </button>
      <button
        class="flex w-full items-center gap-2 rounded-md px-3 py-2 text-sm text-custom-text-200 hover:bg-custom-background-80 transition-colors"
        @click="handleExportPDF(); close()"
      >
        <FileText class="h-3.5 w-3.5" />
        Export PDF
      </button>
    </template>
  </PDropdown>
</template>
