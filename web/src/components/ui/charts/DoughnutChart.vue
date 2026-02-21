<script setup lang="ts">
import { computed } from 'vue'
import { Doughnut } from 'vue-chartjs'
import {
  Chart as ChartJS,
  ArcElement,
  Tooltip,
  Legend,
} from 'chart.js'

ChartJS.register(ArcElement, Tooltip, Legend)

interface Props {
  labels: string[]
  data: number[]
  colors: string[]
  title?: string
}

const props = defineProps<Props>()

function getCssVar(name: string, fallback: string): string {
  return getComputedStyle(document.documentElement).getPropertyValue(name).trim() || fallback
}

const chartData = computed(() => ({
  labels: props.labels,
  datasets: [{
    data: props.data,
    backgroundColor: props.colors,
    borderWidth: 0,
    hoverOffset: 4,
  }],
}))

const chartOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  cutout: '65%',
  plugins: {
    legend: {
      position: 'bottom' as const,
      labels: {
        color: getCssVar('--color-text-200', '#525252'),
        padding: 16,
        usePointStyle: true,
        pointStyleWidth: 8,
        font: { family: 'Inter', size: 12 },
      },
    },
    tooltip: {
      backgroundColor: getCssVar('--color-background-100', '#FFFFFF'),
      titleColor: getCssVar('--color-text-100', '#171717'),
      bodyColor: getCssVar('--color-text-200', '#525252'),
      borderColor: getCssVar('--color-border-200', '#EAEBEB'),
      borderWidth: 1,
      padding: 12,
      cornerRadius: 8,
      titleFont: { family: 'Inter', weight: '600' as const },
      bodyFont: { family: 'Inter' },
    },
  },
}))
</script>

<template>
  <div class="rounded-xl border border-custom-border-200 bg-custom-background-100 p-5">
    <h3 v-if="title" class="mb-4 text-sm font-semibold text-custom-text-100">{{ title }}</h3>
    <div class="h-[250px]">
      <Doughnut :data="chartData" :options="chartOptions" />
    </div>
  </div>
</template>
