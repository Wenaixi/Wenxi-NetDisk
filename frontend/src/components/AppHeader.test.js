import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import AppHeader from './AppHeader.vue'

// Mock naive-ui
vi.mock('naive-ui', () => ({
  NDropdown: {
    template: '<div class="n-dropdown"><slot /></div>',
    props: ['options', 'onSelect']
  },
  NIcon: {
    template: '<span class="n-icon"><slot /></span>',
    props: ['size', 'component']
  },
  NButton: {
    template: '<button class="n-button" :type="type"><slot /></button>',
    props: ['type', 'text']
  }
}))

// Mock @vicons/ionicons5
vi.mock('@vicons/ionicons5', () => ({
  ChevronDown: { template: '<span data-icon="chevron-down"></span>' },
  Settings: { template: '<span data-icon="settings"></span>' },
  LogOutOutline: { template: '<span data-icon="logout"></span>' }
}))

// Mock vue-router
const mockPush = vi.fn()
vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: mockPush
  }),
  useRoute: () => ({ path: '/dashboard' })
}))

// Mock auth store - use a function so each test can override
const mockAuthStore = {
  user: { email: 'test@example.com' },
  logout: vi.fn()
}

vi.mock('../stores/auth', () => ({
  useAuthStore: vi.fn(() => mockAuthStore)
}))

import { useAuthStore } from '../stores/auth'

const mountComponent = () => {
  return mount(AppHeader, {
    global: {
      plugins: [],
      stubs: {
        RouterLink: true
      }
    }
  })
}

describe('AppHeader', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    mockPush.mockReset()
    mockAuthStore.user = { email: 'test@example.com' }
    mockAuthStore.logout = vi.fn()
  })

  it('should render header with app title', () => {
    const wrapper = mountComponent()

    expect(wrapper.text()).toContain('文希云盘')
  })

  it('should have navigation buttons', () => {
    const wrapper = mountComponent()

    expect(wrapper.text()).toContain('本地文件')
    expect(wrapper.text()).toContain('蓝奏云')
  })

  it('should display user email', () => {
    const wrapper = mountComponent()

    expect(wrapper.vm.user).toEqual({ email: 'test@example.com' })
  })

  it('should have dropdown menu with settings and logout', () => {
    const wrapper = mountComponent()

    expect(wrapper.vm.menuOptions).toHaveLength(2)
    expect(wrapper.vm.menuOptions[0].key).toBe('settings')
    expect(wrapper.vm.menuOptions[1].key).toBe('logout')
  })

  it('should navigate to settings when menu select is settings', () => {
    const wrapper = mountComponent()

    wrapper.vm.handleMenuSelect('settings')

    expect(mockPush).toHaveBeenCalledWith('/settings')
  })

  it('should call logout and redirect to login when menu select is logout', () => {
    const wrapper = mountComponent()

    wrapper.vm.handleMenuSelect('logout')

    expect(mockAuthStore.logout).toHaveBeenCalled()
    expect(mockPush).toHaveBeenCalledWith('/login')
  })

  it('should highlight current route as primary', () => {
    const wrapper = mountComponent()

    expect(wrapper.vm.isActive('/dashboard')).toBe('primary')
    expect(wrapper.vm.isActive('/lanzou')).toBe(undefined)
  })

  it('should handle missing user gracefully', () => {
    mockAuthStore.user = null

    const wrapper = mountComponent()

    expect(wrapper.vm.user).toBeNull()
  })
})
