import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const routes = [
  {
    path: '/',
    redirect: '/dashboard'
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue'),
    meta: { guest: true }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('../views/Register.vue'),
    meta: { guest: true }
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import('../views/Dashboard.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/settings',
    name: 'Settings',
    component: () => import('../views/Settings.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/lanzou-settings',
    name: 'LanzouSettings',
    component: () => import('../views/LanzouSettings.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/lanzou',
    name: 'LanZouBrowser',
    component: () => import('../views/LanZouBrowser.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/recycle',
    name: 'RecycleBin',
    component: () => import('../views/RecycleBin.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/parse',
    name: 'ShareParse',
    component: () => import('../views/ShareParse.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/sync',
    name: 'Sync',
    component: () => import('../views/Sync.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/task-history',
    name: 'TaskHistory',
    component: () => import('../views/TaskHistory.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/split-merge',
    name: 'SplitMerge',
    component: () => import('../views/SplitMerge.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/upload-tasks',
    name: 'UploadTasks',
    component: () => import('../views/UploadTasks.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/download-tasks',
    name: 'DownloadTasks',
    component: () => import('../views/DownloadTasks.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/completed-tasks',
    name: 'CompletedTasks',
    component: () => import('../views/CompletedTasks.vue'),
    meta: { requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  if (to.meta.requiresAuth && !authStore.isLoggedIn) {
    next('/login')
  } else if (to.meta.guest && authStore.isLoggedIn) {
    next('/dashboard')
  } else {
    next()
  }
})

export default router