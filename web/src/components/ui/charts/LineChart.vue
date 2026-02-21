<script setup lang="ts">
import { computed } from 'vue'
import { Line } from 'vue-chartjs'
import {
  Chart as ChartJS,
  LineElement,
  PointElement,
  CategoryScale,
  LinearScale,
  Filler,
  Tooltip,
} from 'chart.js'

ChartJS.register(LineElement, PointElement, CategoryScale, LinearScale, Filler, Tooltip)

interface Props {
  labels: string[]
  data: number[]
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
    borderColor: '#2B7ACC',
    backgroundColor: 'rgba(43, 122, 204, 0.1)',
    fill: true,
    tension: 0.3,
    pointBackgroundColor: '#2B7ACC',
    pointBorderColor: '#FFFFFF',
    pointBorderWidth: 2,
    pointRadius: 3,
    pointHoverRadius: 5,
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
        font: { family: 'Inter', size: 11 },
        maxTicksLimit: 10,
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
      <Line :data="chartData" :options="chartOptions" />
    </div>
  </div>
</template>
