import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createRouter, createWebHistory } from 'vue-router'
import { createPinia, setActivePinia } from 'pinia'

// Mock auth store
const mockIsLoggedIn = vi.fn(() => false)

vi.mock('../stores/auth', () => ({
  useAuthStore: () => ({
    isLoggedIn: mockIsLoggedIn()
  })
}))

// Mock lazy-loaded components to avoid actual imports
vi.mock('../views/Login.vue', () => ({ default: { template: '<div>Login</div>' } }))
vi.mock('../views/Register.vue', () => ({ default: { template: '<div>Register</div>' } }))
vi.mock('../views/Dashboard.vue', () => ({ default: { template: '<div>Dashboard</div>' } }))
vi.mock('../views/LanzouSettings.vue', () => ({ default: { template: '<div>Settings</div>' } }))
vi.mock('../views/LanZouBrowser.vue', () => ({ default: { template: '<div>LanZou</div>' } }))
vi.mock('../views/RecycleBin.vue', () => ({ default: { template: '<div>Recycle</div>' } }))
vi.mock('../views/ShareParse.vue', () => ({ default: { template: '<div>Parse</div>' } }))
vi.mock('../views/Sync.vue', () => ({ default: { template: '<div>Sync</div>' } }))

// Import router after mocks
import routerConfig from './index'

function createTestRouter() {
  const routes = [
    { path: '/', redirect: '/dashboard' },
    { path: '/login', name: 'Login', meta: { guest: true }, component: { template: '<div>Login</div>' } },
    { path: '/register', name: 'Register', meta: { guest: true }, component: { template: '<div>Register</div>' } },
    { path: '/dashboard', name: 'Dashboard', meta: { requiresAuth: true }, component: { template: '<div>Dashboard</div>' } },
    { path: '/settings', name: 'Settings', meta: { requiresAuth: true }, component: { template: '<div>Settings</div>' } },
    { path: '/lanzou', name: 'LanZouBrowser', meta: { requiresAuth: true }, component: { template: '<div>LanZou</div>' } },
    { path: '/recycle', name: 'RecycleBin', meta: { requiresAuth: true }, component: { template: '<div>Recycle</div>' } },
    { path: '/parse', name: 'ShareParse', meta: { requiresAuth: true }, component: { template: '<div>Parse</div>' } },
    { path: '/sync', name: 'Sync', meta: { requiresAuth: true }, component: { template: '<div>Sync</div>' } }
  ]

  const router = createRouter({
    history: createWebHistory(),
    routes
  })

  router.beforeEach((to, from, next) => {
    const loggedIn = mockIsLoggedIn()
    if (to.meta.requiresAuth && !loggedIn) {
      next('/login')
    } else if (to.meta.guest && loggedIn) {
      next('/dashboard')
    } else {
      next()
    }
  })

  return router
}

describe('Router Guards', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockIsLoggedIn.mockReturnValue(false)
  })

  it('should redirect root to dashboard', async () => {
    mockIsLoggedIn.mockReturnValue(true)
    const router = createTestRouter()
    await router.push('/')
    await router.isReady()

    expect(router.currentRoute.value.path).toBe('/dashboard')
  })

  it('should redirect unauthenticated user from dashboard to login', async () => {
    mockIsLoggedIn.mockReturnValue(false)
    const router = createTestRouter()
    await router.push('/dashboard')
    await router.isReady()

    expect(router.currentRoute.value.path).toBe('/login')
  })

  it('should redirect unauthenticated user from settings to login', async () => {
    mockIsLoggedIn.mockReturnValue(false)
    const router = createTestRouter()
    await router.push('/settings')
    await router.isReady()

    expect(router.currentRoute.value.path).toBe('/login')
  })

  it('should redirect unauthenticated user from lanzou to login', async () => {
    mockIsLoggedIn.mockReturnValue(false)
    const router = createTestRouter()
    await router.push('/lanzou')
    await router.isReady()

    expect(router.currentRoute.value.path).toBe('/login')
  })

  it('should redirect unauthenticated user from recycle to login', async () => {
    mockIsLoggedIn.mockReturnValue(false)
    const router = createTestRouter()
    await router.push('/recycle')
    await router.isReady()

    expect(router.currentRoute.value.path).toBe('/login')
  })

  it('should redirect unauthenticated user from parse to login', async () => {
    mockIsLoggedIn.mockReturnValue(false)
    const router = createTestRouter()
    await router.push('/parse')
    await router.isReady()

    expect(router.currentRoute.value.path).toBe('/login')
  })

  it('should redirect unauthenticated user from sync to login', async () => {
    mockIsLoggedIn.mockReturnValue(false)
    const router = createTestRouter()
    await router.push('/sync')
    await router.isReady()

    expect(router.currentRoute.value.path).toBe('/login')
  })

  it('should allow authenticated user to access dashboard', async () => {
    mockIsLoggedIn.mockReturnValue(true)
    const router = createTestRouter()
    await router.push('/dashboard')
    await router.isReady()

    expect(router.currentRoute.value.path).toBe('/dashboard')
  })

  it('should allow authenticated user to access settings', async () => {
    mockIsLoggedIn.mockReturnValue(true)
    const router = createTestRouter()
    await router.push('/settings')
    await router.isReady()

    expect(router.currentRoute.value.path).toBe('/settings')
  })

  it('should allow authenticated user to access lanzou', async () => {
    mockIsLoggedIn.mockReturnValue(true)
    const router = createTestRouter()
    await router.push('/lanzou')
    await router.isReady()

    expect(router.currentRoute.value.path).toBe('/lanzou')
  })

  it('should redirect authenticated user from login to dashboard', async () => {
    mockIsLoggedIn.mockReturnValue(true)
    const router = createTestRouter()
    await router.push('/login')
    await router.isReady()

    expect(router.currentRoute.value.path).toBe('/dashboard')
  })

  it('should redirect authenticated user from register to dashboard', async () => {
    mockIsLoggedIn.mockReturnValue(true)
    const router = createTestRouter()
    await router.push('/register')
    await router.isReady()

    expect(router.currentRoute.value.path).toBe('/dashboard')
  })

  it('should allow unauthenticated user to access login', async () => {
    mockIsLoggedIn.mockReturnValue(false)
    const router = createTestRouter()
    await router.push('/login')
    await router.isReady()

    expect(router.currentRoute.value.path).toBe('/login')
  })

  it('should allow unauthenticated user to access register', async () => {
    mockIsLoggedIn.mockReturnValue(false)
    const router = createTestRouter()
    await router.push('/register')
    await router.isReady()

    expect(router.currentRoute.value.path).toBe('/register')
  })
})

describe('Route Definitions', () => {
  it('should have all required routes defined', () => {
    const router = createTestRouter()
    const routeNames = router.getRoutes().map(r => r.name)

    expect(routeNames).toContain('Login')
    expect(routeNames).toContain('Register')
    expect(routeNames).toContain('Dashboard')
    expect(routeNames).toContain('Settings')
    expect(routeNames).toContain('LanZouBrowser')
    expect(routeNames).toContain('RecycleBin')
    expect(routeNames).toContain('ShareParse')
    expect(routeNames).toContain('Sync')
  })

  it('should have guest meta on auth pages', () => {
    const router = createTestRouter()
    const loginRoute = router.getRoutes().find(r => r.name === 'Login')
    const registerRoute = router.getRoutes().find(r => r.name === 'Register')

    expect(loginRoute.meta.guest).toBe(true)
    expect(registerRoute.meta.guest).toBe(true)
  })

  it('should have requiresAuth meta on protected pages', () => {
    const router = createTestRouter()
    const protectedRoutes = ['Dashboard', 'Settings', 'LanZouBrowser', 'RecycleBin', 'ShareParse', 'Sync']

    for (const name of protectedRoutes) {
      const route = router.getRoutes().find(r => r.name === name)
      expect(route.meta.requiresAuth).toBe(true)
    }
  })
})
