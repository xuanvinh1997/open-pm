import apiClient from './client'
import type {
  IssueAnalyticsReport,
  WorkLogReport,
  CycleReportData,
  VelocityPoint,
  ProjectHealthReport,
  Cycle,
} from '@/types/report.types'

const base = (slug: string, projectId: string) =>
  `/api/v1/workspaces/${slug}/projects/${projectId}`

export const reportApi = {
  listCycles(slug: string, projectId: string) {
    return apiClient.get<{ results: Cycle[] }>(`${base(slug, projectId)}/cycles`)
  },

  getIssueAnalytics(slug: string, projectId: string) {
    return apiClient.get<IssueAnalyticsReport>(`${base(slug, projectId)}/reports/issues`)
  },

  getWorkLogs(slug: string, projectId: string, startDate?: string, endDate?: string) {
    const params = new URLSearchParams()
    if (startDate) params.set('start_date', startDate)
    if (endDate) params.set('end_date', endDate)
    const qs = params.toString()
    return apiClient.get<WorkLogReport>(`${base(slug, projectId)}/reports/work-logs${qs ? '?' + qs : ''}`)
  },

  getCycleReport(slug: string, projectId: string, cycleId: string) {
    return apiClient.get<CycleReportData>(`${base(slug, projectId)}/reports/cycles/${cycleId}`)
  },

  getVelocity(slug: string, projectId: string) {
    return apiClient.get<{ results: VelocityPoint[] }>(`${base(slug, projectId)}/reports/velocity`)
  },

  getHealth(slug: string, projectId: string) {
    return apiClient.get<ProjectHealthReport>(`${base(slug, projectId)}/reports/health`)
  },

  exportPDF(slug: string, projectId: string, reportType: string) {
    return apiClient.get(`${base(slug, projectId)}/reports/export?type=${reportType}`, {
      responseType: 'blob',
    })
  },
}
