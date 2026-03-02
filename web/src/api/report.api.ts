import apiClient from './client'
import type {
  IssueAnalyticsReport,
  WorkLogReport,
  SprintReportData,
  VelocityPoint,
  ProjectHealthReport,
  Sprint,
} from '@/types/report.types'

const base = (slug: string, projectId: string) =>
  `/api/v1/workspaces/${slug}/projects/${projectId}`

export const reportApi = {
  listSprints(slug: string, projectId: string) {
    return apiClient.get<{ results: Sprint[] }>(`${base(slug, projectId)}/sprints`)
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

  getSprintReport(slug: string, projectId: string, sprintId: string) {
    return apiClient.get<SprintReportData>(`${base(slug, projectId)}/reports/sprints/${sprintId}`)
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
