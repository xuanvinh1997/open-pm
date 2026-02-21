<script setup lang="ts">
import { computed } from 'vue'
import { Bar } from 'vue-chartjs'
import {
  Chart as ChartJS,
  BarElement,
  CategoryScale,
  LinearScale,
  Tooltip,
} from 'chart.js'

ChartJS.register(BarElement, CategoryScale, LinearScale, Tooltip)

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
    borderRadius: 6,
    borderSkipped: false as const,
    maxBarThickness: 40,
  }],
}))

const chartOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { display: false },
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
  scales: {
    x: {
      grid: { display: false },
      ticks: {
        color: getCssVar('--color-text-300', '#A3A3A3'),
        font: { family: 'Inter', size: 12 },
      },
      border: { display: false },
    },
    y: {
      grid: {
        color: getCssVar('--color-border-100', '#E5E5E5'),
      },
      ticks: {
        color: getCssVar('--color-text-300', '#A3A3A3'),
        font: { family: 'Inter', size: 12 },
        stepSize: 1,
      },
      border: { display: false },
      beginAtZero: true,
    },
  },
}))
</script>

<template>
  <div class="rounded-xl border border-custom-border-200 bg-custom-background-100 p-5">
    <h3 v-if="title" class="mb-4 text-sm font-semibold text-custom-text-100">{{ title }}</h3>
    <div class="h-[250px]">
      <Bar :data="chartData" :options="chartOptions" />
    </div>
  </div>
</template>
