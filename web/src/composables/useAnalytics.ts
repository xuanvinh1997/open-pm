import { computed, type Ref } from 'vue'
import type { Issue, IssuePriority } from '@/types/issue.types'
import type { State } from '@/types/project.types'
import { PRIORITY_CONFIG, STATE_GROUP_CONFIG } from '@/utils/issue-helpers'

interface AnalyticsInput {
  issues: Ref<Issue[]>
  states: Ref<State[]>
}

export function useAnalytics({ issues, states }: AnalyticsInput) {
  const issuesByState = computed(() => {
    const counts = new Map<string, { name: string; color: string; count: number }>()
    for (const state of states.value) {
      counts.set(state.id, { name: state.name, color: state.color, count: 0 })
    }
    for (const issue of issues.value) {
      const entry = counts.get(issue.state_id || '')
      if (entry) entry.count++
    }
    return Array.from(counts.values())
  })

  const issuesByStateGroup = computed(() => {
    const groups: Record<string, number> = {}
    const stateMap = new Map(states.value.map((s) => [s.id, s.group]))
    for (const issue of issues.value) {
      const group = stateMap.get(issue.state_id || '') || 'backlog'
      groups[group] = (groups[group] || 0) + 1
    }
    return Object.entries(groups).map(([group, count]) => ({
      group,
      count,
      color: STATE_GROUP_CONFIG[group]?.color || '#A3A3A3',
      label: group.charAt(0).toUpperCase() + group.slice(1),
    }))
  })

  const issuesByPriority = computed(() => {
    const priorities: IssuePriority[] = ['urgent', 'high', 'medium', 'low', 'none']
    const counts: Record<string, number> = {}
    for (const issue of issues.value) {
      counts[issue.priority] = (counts[issue.priority] || 0) + 1
    }
    return priorities.map((p) => ({
      priority: p,
      count: counts[p] || 0,
      color: PRIORITY_CONFIG[p].color,
      label: PRIORITY_CONFIG[p].label,
    }))
  })

  const completionOverTime = computed(() => {
    const days = 30
    const now = new Date()
    const buckets: { date: string; count: number }[] = []
    for (let i = days - 1; i >= 0; i--) {
      const d = new Date(now)
      d.setDate(d.getDate() - i)
      buckets.push({ date: d.toISOString().slice(0, 10), count: 0 })
    }
    for (const issue of issues.value) {
      if (issue.completed_at) {
        const dateStr = issue.completed_at.slice(0, 10)
        const bucket = buckets.find((b) => b.date === dateStr)
        if (bucket) bucket.count++
      }
    }
    return buckets
  })

  const stats = computed(() => ({
    total: issues.value.length,
    completed: issues.value.filter((i) => i.completed_at).length,
    active: issues.value.filter((i) => {
      const state = states.value.find((s) => s.id === i.state_id)
      return state?.group === 'started'
    }).length,
    overdue: issues.value.filter((i) => {
      if (!i.target_date || i.completed_at) return false
      return new Date(i.target_date) < new Date()
    }).length,
  }))

  return {
    issuesByState,
    issuesByStateGroup,
    issuesByPriority,
    completionOverTime,
    stats,
  }
}
