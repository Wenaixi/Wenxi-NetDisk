import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import LanzouSettings from './LanzouSettings.vue'

// Mock naive-ui
vi.mock('naive-ui', () => ({
  NButton: {
    template: '<button class="n-button" :type="type" :loading="loading" @click="$emit(\'click\')"><slot /></button>',
    props: ['type', 'loading', 'block', 'size', 'text']
  },
  NIcon: {
    template: '<span class="n-icon"><slot /></span>',
    props: ['size', 'component']
  },
  NModal: {
    template: '<div class="n-modal" v-if="show"><slot /></div>',
    props: ['show']
  },
  NCard: {
    template: '<div class="n-card"><slot name="header" /><slot /><slot name="footer" /></div>',
    props: ['title']
  },
  NForm: {
    template: '<form class="n-form"><slot /></form>',
    props: ['model', 'rules', 'ref']
  },
  NFormItem: {
    template: '<div class="n-form-item"><span class="label">{{ label }}</span><slot /></div>',
    props: ['path', 'label']
  },
  NInput: {
    template: '<textarea class="n-input" :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />',
    props: ['modelValue', 'type', 'placeholder', 'rows']
  },
  useMessage: () => ({
    success: vi.fn(),
    error: vi.fn(),
    warning: vi.fn(),
    info: vi.fn()
  })
}))

// Mock @vicons/ionicons5 - includes all icons used by AppHeader + this page
vi.mock('@vicons/ionicons5', () => ({
  ChevronDown: { template: '<span data-icon="chevron-down"></span>' },
  Settings: { template: '<span data-icon="settings"></span>' },
  LogOutOutline: { template: '<span data-icon="logout"></span>' },
  CheckmarkCircle: { template: '<span data-icon="checkmark-circle"></span>' },
  Warning: { template: '<span data-icon="warning"></span>' }
}))

// Mock vue-router
vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: vi.fn()
  }),
  useRoute: () => ({ path: '/settings' })
}))

// Mock auth store
vi.mock('../stores/auth', () => ({
  useAuthStore: () => ({
    user: { email: 'test@example.com' }
  })
}))

// Mock lanzouAPI - vi.hoisted because vi.mock is hoisted
const { mockStatus, mockConnect, mockDisconnect } = vi.hoisted(() => ({
  mockStatus: vi.fn(),
  mockConnect: vi.fn(),
  mockDisconnect: vi.fn()
}))

vi.mock('../api/lanzou', () => ({
  lanzouAPI: {
    status: mockStatus,
    connect: mockConnect,
    disconnect: mockDisconnect
  }
}))

import { lanzouAPI } from '../api/lanzou'

const mountComponent = () => {
  return mount(LanzouSettings, {
    global: {
      stubs: {
        RouterLink: true,
        AppHeader: true
      }
    }
  })
}

describe('LanzouSettings.vue', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    mockStatus.mockReset()
    mockConnect.mockReset()
    mockDisconnect.mockReset()
  })

  it('should render settings page title', () => {
    const wrapper = mountComponent()
    expect(wrapper.text()).toContain('蓝奏云设置')
  })

  it('should show connected status when connected', async () => {
    mockStatus.mockResolvedValue({ connected: true, info: { uid: '12345' } })
    const wrapper = mountComponent()
    await new Promise(r => setTimeout(r, 10))

    expect(wrapper.text()).toContain('已连接到蓝奏云')
    expect(wrapper.text()).toContain('12345')
  })

  it('should show disconnected status when not connected', async () => {
    mockStatus.mockResolvedValue({ connected: false })
    const wrapper = mountComponent()
    await wrapper.vm.$nextTick()

    expect(wrapper.text()).toContain('未连接蓝奏云')
  })

  it('should show disconnected status on fetch error', async () => {
    mockStatus.mockRejectedValue(new Error('Network error'))
    const wrapper = mountComponent()
    await wrapper.vm.$nextTick()

    expect(wrapper.text()).toContain('未连接蓝奏云')
  })

  it('should open connect modal when connect button clicked', async () => {
    mockStatus.mockResolvedValue({ connected: false })
    const wrapper = mountComponent()
    await wrapper.vm.$nextTick()

    const connectButton = wrapper.findAll('.n-button').find(b => b.text() === '连接蓝奏云')
    await connectButton.trigger('click')

    expect(wrapper.vm.showConnectModal).toBe(true)
  })

  it('should call lanzouAPI.connect with cookie on handleConnect', async () => {
    mockStatus.mockResolvedValue({ connected: false })
    mockConnect.mockResolvedValue({ success: true })
    const wrapper = mountComponent()
    await wrapper.vm.$nextTick()

    wrapper.vm.connectForm.cookie = 'test-cookie-value'
    await wrapper.vm.handleConnect()

    expect(mockConnect).toHaveBeenCalledWith({ cookie: 'test-cookie-value' })
  })

  it('should update connected status after successful connection', async () => {
    mockStatus.mockResolvedValue({ connected: false })
    mockConnect.mockResolvedValue({ success: true })
    mockStatus.mockResolvedValue({ connected: true, info: { uid: '999' } })
    const wrapper = mountComponent()
    await new Promise(r => setTimeout(r, 10))

    wrapper.vm.connectForm.cookie = 'test-cookie'
    await wrapper.vm.handleConnect()
    await new Promise(r => setTimeout(r, 10))

    expect(wrapper.vm.connected).toBe(true)
  })

  it('should show error message on connect failure', async () => {
    mockStatus.mockResolvedValue({ connected: false })
    mockConnect.mockRejectedValue(new Error('Invalid cookie'))
    const wrapper = mountComponent()
    await wrapper.vm.$nextTick()

    wrapper.vm.connectForm.cookie = 'bad-cookie'
    await wrapper.vm.handleConnect()

    expect(wrapper.vm.connecting).toBe(false)
  })

  it('should call lanzouAPI.disconnect on handleDisconnect', async () => {
    mockStatus.mockResolvedValue({ connected: true, info: { uid: '123' } })
    mockDisconnect.mockResolvedValue({ success: true })
    const wrapper = mountComponent()
    await wrapper.vm.$nextTick()

    await wrapper.vm.handleDisconnect()

    expect(mockDisconnect).toHaveBeenCalled()
    expect(wrapper.vm.connected).toBe(false)
  })

  it('should show error message on disconnect failure', async () => {
    mockStatus.mockResolvedValue({ connected: true, info: { uid: '123' } })
    mockDisconnect.mockRejectedValue(new Error('Server error'))
    const wrapper = mountComponent()
    await wrapper.vm.$nextTick()

    await wrapper.vm.handleDisconnect()

    expect(wrapper.vm.connected).toBe(true)
  })

  it('should have cookie input validation rule', () => {
    const wrapper = mountComponent()
    expect(wrapper.vm.rules.cookie[0].required).toBe(true)
  })

  it('should close modal and clear form after successful connect', async () => {
    mockStatus.mockResolvedValue({ connected: false })
    mockConnect.mockResolvedValue({ success: true })
    const wrapper = mountComponent()
    await wrapper.vm.$nextTick()

    wrapper.vm.showConnectModal = true
    wrapper.vm.connectForm.cookie = 'some-cookie'
    await wrapper.vm.handleConnect()

    expect(wrapper.vm.showConnectModal).toBe(false)
    expect(wrapper.vm.connectForm.cookie).toBe('')
  })
})
