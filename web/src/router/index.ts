import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    // Auth routes
    {
      path: '/auth',
      component: () => import('@/components/layout/AuthLayout.vue'),
      children: [
        { path: 'login', name: 'login', component: () => import('@/views/auth/LoginView.vue') },
        { path: 'signup', name: 'signup', component: () => import('@/views/auth/SignupView.vue') },
      ],
    },
    // Onboarding
    {
      path: '/onboarding',
      meta: { requiresAuth: true },
      children: [
        {
          path: 'workspace',
          name: 'create-workspace',
          component: () => import('@/views/onboarding/CreateWorkspaceView.vue'),
        },
      ],
    },
    // App routes
    {
      path: '/:workspaceSlug',
      component: () => import('@/components/layout/AppLayout.vue'),
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          name: 'workspace-dashboard',
          component: () => import('@/views/workspace/WorkspaceDashboard.vue'),
        },
        {
          path: 'projects',
          name: 'project-list',
          component: () => import('@/views/project/ProjectListView.vue'),
        },
        {
          path: 'projects/:projectId/issues',
          name: 'issues',
          component: () => import('@/views/project/IssuesView.vue'),
        },
        {
          path: 'projects/:projectId/board',
          name: 'board',
          component: () => import('@/views/project/BoardView.vue'),
        },
        {
          path: 'projects/:projectId/analytics',
          name: 'analytics',
          component: () => import('@/views/project/AnalyticsView.vue'),
        },
        {
          path: 'projects/:projectId/reports',
          name: 'reports',
          component: () => import('@/views/project/ReportsView.vue'),
        },
        {
          path: 'projects/:projectId/sprints',
          name: 'sprints',
          component: () => import('@/views/project/SprintsView.vue'),
        },
        {
          path: 'projects/:projectId/sprints/:sprintId',
          name: 'sprint-detail',
          component: () => import('@/views/project/SprintDetailView.vue'),
        },
        {
          path: 'projects/:projectId/epics',
          name: 'epics',
          component: () => import('@/views/project/EpicsView.vue'),
        },
        {
          path: 'projects/:projectId/epics/:epicId',
          name: 'epic-detail',
          component: () => import('@/views/project/EpicDetailView.vue'),
        },
        {
          path: 'projects/:projectId/pages',
          name: 'pages',
          component: () => import('@/views/project/PagesView.vue'),
        },
        {
          path: 'projects/:projectId/pages/:pageId',
          name: 'page-detail',
          component: () => import('@/views/project/PageDetailView.vue'),
        },
        {
          path: 'projects/:projectId/issues/:issueId',
          name: 'issue-detail',
          component: () => import('@/views/project/IssueDetailView.vue'),
        },
        {
          path: 'members',
          name: 'workspace-members',
          component: () => import('@/views/workspace/WorkspaceMembersView.vue'),
        },
        {
          path: 'settings',
          name: 'workspace-settings',
          component: () => import('@/views/workspace/WorkspaceSettings.vue'),
        },
      ],
    },
    // Root redirect (handled in navigation guard)
    { path: '/', name: 'root', component: () => import('@/views/auth/LoginView.vue') },
  ],
})

// Navigation guard
router.beforeEach(async (to) => {
  const isAuthenticated = !!localStorage.getItem('access_token')

  if (to.meta.requiresAuth && !isAuthenticated) {
    return { name: 'login' }
  }

  // Redirect authenticated users away from auth pages and root
  if (isAuthenticated && (to.path === '/' || to.path.startsWith('/auth'))) {
    // Try to redirect to first workspace
    try {
      const { workspaceApi } = await import('@/api/workspace.api')
      const { data } = await workspaceApi.list()
      if (data.results && data.results.length > 0) {
        return `/${data.results[0].slug}`
      }
      return { name: 'create-workspace' }
    } catch {
      // Token may be invalid, let them go to login
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
      return { name: 'login' }
    }
  }

  // Redirect unauthenticated users from root to login
  if (to.path === '/' && !isAuthenticated) {
    return { name: 'login' }
  }
})

export default router
