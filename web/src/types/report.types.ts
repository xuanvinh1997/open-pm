import type { Issue } from './issue.types'

export interface Cycle {
  id: string
  project_id: string
  workspace_id: string
  name: string
  description: string
  start_date?: string
  end_date?: string
  owned_by: string
  sort_order: number
  created_at: string
  updated_at: string
}

export interface IssuesByStateReport {
  state_id: string
  state_name: string
  state_group: string
  color: string
  count: number
}

export interface IssuesByPriorityReport {
  priority: string
  count: number
}

export interface IssuesByTypeReport {
  issue_type: string
  count: number
}

export interface IssuesByAssigneeReport {
  user_id: string
  display_name: string
  avatar_url: string
  count: number
}

export interface WorkLogSummaryByMember {
  user_id: string
  display_name: string
  avatar_url: string
  total_minutes: number
  log_count: number
}

export interface WorkLogSummaryByIssue {
  issue_id: string
  issue_name: string
  sequence_id: number
  total_minutes: number
}

export interface WorkLogSummaryByDate {
  date: string
  total_minutes: number
}

export interface BurndownPoint {
  date: string
  remaining_count: number
  completed_count: number
  ideal_remaining: number
}

export interface VelocityPoint {
  cycle_id: string
  cycle_name: string
  completed_count: number
  total_count: number
}

export interface ProjectHealthReport {
  overdue_issues: Issue[]
  blocked_issues: Issue[]
  unestimated_issues: Issue[]
  overdue_count: number
  blocked_count: number
  unestimated_count: number
}

export interface IssueAnalyticsReport {
  by_state: IssuesByStateReport[]
  by_priority: IssuesByPriorityReport[]
  by_type: IssuesByTypeReport[]
  by_assignee: IssuesByAssigneeReport[]
  total: number
  completed: number
  overdue: number
}

export interface WorkLogReport {
  by_member: WorkLogSummaryByMember[]
  by_issue: WorkLogSummaryByIssue[]
  by_date: WorkLogSummaryByDate[]
  total_minutes: number
  start_date: string
  end_date: string
}

export interface CycleReportData {
  cycle: Cycle
  total_issues: number
  completed_issues: number
  burndown_points: BurndownPoint[]
  issues_by_state: IssuesByStateReport[]
}
