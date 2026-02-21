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
          path: 'projects/:projectId/issues/:issueId',
          name: 'issue-detail',
          component: () => import('@/views/project/IssueDetailView.vue'),
        },
        {
          path: 'settings',
          name: 'workspace-settings',
          component: () => import('@/views/workspace/WorkspaceSettings.vue'),
        },
      ],
    },
    // Root redirect
    { path: '/', redirect: '/auth/login' },
  ],
})

// Navigation guard
router.beforeEach((to) => {
  const isAuthenticated = !!localStorage.getItem('access_token')

  if (to.meta.requiresAuth && !isAuthenticated) {
    return { name: 'login' }
  }

  if (to.path.startsWith('/auth') && isAuthenticated) {
    return '/'
  }
})

export default router
