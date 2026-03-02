<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import ViewHeader from '@/components/issues/ViewHeader.vue'
import IssueAnalyticsTab from '@/components/reports/IssueAnalyticsTab.vue'
import WorkLogTab from '@/components/reports/WorkLogTab.vue'
import SprintReportTab from '@/components/reports/SprintReportTab.vue'
import HealthTab from '@/components/reports/HealthTab.vue'

const route = useRoute()
const slug = computed(() => route.params.workspaceSlug as string)
const projectId = computed(() => route.params.projectId as string)

type TabKey = 'issues' | 'work-logs' | 'sprints' | 'health'
const activeTab = ref<TabKey>('issues')

const tabs: { key: TabKey; label: string }[] = [
  { key: 'issues', label: 'Issue Analytics' },
  { key: 'work-logs', label: 'Work Logs' },
  { key: 'sprints', label: 'Sprint Reports' },
  { key: 'health', label: 'Project Health' },
]
</script>

<template>
  <div class="flex h-full flex-col">
    <ViewHeader active-view="reports" />

    <!-- Tab bar -->
    <div class="border-b border-custom-border-200 px-6">
      <nav class="-mb-px flex gap-6">
        <button
          v-for="tab in tabs"
          :key="tab.key"
          class="whitespace-nowrap border-b-2 py-3 text-sm font-medium transition-colors"
          :class="
            activeTab === tab.key
              ? 'border-brand-500 text-brand-500'
              : 'border-transparent text-custom-text-300 hover:border-custom-border-200 hover:text-custom-text-200'
          "
          @click="activeTab = tab.key"
        >
          {{ tab.label }}
        </button>
      </nav>
    </div>

    <!-- Tab content -->
    <div class="flex-1 overflow-y-auto p-6">
      <IssueAnalyticsTab v-if="activeTab === 'issues'" :slug="slug" :project-id="projectId" />
      <WorkLogTab v-if="activeTab === 'work-logs'" :slug="slug" :project-id="projectId" />
      <SprintReportTab v-if="activeTab === 'sprints'" :slug="slug" :project-id="projectId" />
      <HealthTab v-if="activeTab === 'health'" :slug="slug" :project-id="projectId" />
    </div>
  </div>
</template>
