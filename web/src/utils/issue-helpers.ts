import {
  AlertTriangle,
  SignalHigh,
  SignalMedium,
  SignalLow,
  Minus,
  Circle,
  CircleDot,
  Loader2,
  CheckCircle2,
  XCircle,
  BookOpen,
  Bug,
  CheckSquare,
  Zap,
  Link2,
  Ban,
  Lock,
  Copy,
} from 'lucide-vue-next'
import type { Component } from 'vue'
import type { IssuePriority, IssueType, RelationType } from '@/types/issue.types'

interface PriorityConfig {
  icon: Component
  color: string
  label: string
}

interface StateGroupConfig {
  icon: Component
  color: string
}

export const PRIORITY_CONFIG: Record<IssuePriority, PriorityConfig> = {
  urgent: { icon: AlertTriangle, color: '#EF4444', label: 'Urgent' },
  high: { icon: SignalHigh, color: '#F97316', label: 'High' },
  medium: { icon: SignalMedium, color: '#EAB308', label: 'Medium' },
  low: { icon: SignalLow, color: '#3B82F6', label: 'Low' },
  none: { icon: Minus, color: '#A3A3A3', label: 'None' },
}

interface IssueTypeConfig {
  icon: Component
  color: string
  label: string
}

export const ISSUE_TYPE_CONFIG: Record<IssueType, IssueTypeConfig> = {
  story: { icon: BookOpen, color: '#8B5CF6', label: 'Story' },
  bug: { icon: Bug, color: '#EF4444', label: 'Bug' },
  task: { icon: CheckSquare, color: '#3B82F6', label: 'Task' },
  epic: { icon: Zap, color: '#F97316', label: 'Epic' },
}

interface RelationTypeConfig {
  icon: Component
  color: string
  label: string
}

export const RELATION_TYPE_CONFIG: Record<RelationType, RelationTypeConfig> = {
  relates_to: { icon: Link2, color: '#6B7280', label: 'Relates to' },
  blocks: { icon: Ban, color: '#EF4444', label: 'Blocks' },
  blocked_by: { icon: Lock, color: '#F97316', label: 'Blocked by' },
  duplicate_of: { icon: Copy, color: '#8B5CF6', label: 'Duplicate of' },
}

export const STATE_GROUP_CONFIG: Record<string, StateGroupConfig> = {
  backlog: { icon: Circle, color: '#A3A3A3' },
  unstarted: { icon: CircleDot, color: '#80838D' },
  started: { icon: Loader2, color: '#F59E0B' },
  completed: { icon: CheckCircle2, color: '#16A34A' },
  cancelled: { icon: XCircle, color: '#EF4444' },
  triage: { icon: AlertTriangle, color: '#8B5CF6' },
}
