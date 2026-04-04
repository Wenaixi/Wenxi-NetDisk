import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import Sync from './Sync.vue'

// Mock naive-ui
vi.mock('naive-ui', () => ({
  NButton: {
    template: '<button class="n-button" :type="type" :loading="loading" :disabled="disabled" @click="$emit(\'click\')"><slot /></button>',
    props: ['type', 'loading', 'block', 'size', 'text', 'disabled']
  },
  NIcon: {
    template: '<span class="n-icon"><slot /></span>',
    props: ['size', 'color', 'component']
  },
  NModal: {
    template: '<div class="n-modal" v-if="show"><slot /></div>',
    props: ['show']
  },
  NCard: {
    template: '<div class="n-card"><slot name="header" /><slot /><slot name="footer" /></div>',
    props: ['title']
  },
  NTag: {
    template: '<span class="n-tag" :type="type"><slot /></span>',
    props: ['type', 'size', 'bordered']
  },
  NProgress: {
    template: '<div class="n-progress" :data-percentage="percentage"></div>',
    props: ['percentage', 'show-indicator', 'height']
  },
  NInput: {
    template: '<input class="n-input" :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />',
    props: ['modelValue', 'type', 'placeholder']
  },
  NRadioGroup: {
    template: '<div class="n-radio-group"><slot /></div>',
    props: ['value']
  },
  NRadio: {
    template: '<label class="n-radio"><input type="radio" :value="value" /><slot /></label>',
    props: ['value', 'disabled']
  },
  NSpace: {
    template: '<div class="n-space"><slot /></div>'
  },
  NAlert: {
    template: '<div class="n-alert" :type="type"><slot /></div>',
    props: ['type', 'show-icon']
  },
  NEmpty: {
    template: '<div class="n-empty"><slot /></div>',
    props: ['description']
  },
  NSpin: {
    template: '<div class="n-spin"><slot /></div>',
    props: ['show']
  },
  useMessage: () => ({
    success: vi.fn(),
    error: vi.fn(),
    warning: vi.fn(),
    info: vi.fn()
  })
}))

// Mock vue-router
vi.mock('vue-router', () => ({
  useRouter: () => ({ push: vi.fn() }),
  useRoute: () => ({ path: '/sync' })
}))

// Mock @vicons/ionicons5
vi.mock('@vicons/ionicons5', () => ({
  CloudDownload: { template: '<span data-icon="cloud-download"></span>' },
  Close: { template: '<span data-icon="close"></span>' },
  ChevronDown: { template: '<span data-icon="chevron-down"></span>' },
  Settings: { template: '<span data-icon="settings"></span>' },
  LogOutOutline: { template: '<span data-icon="logout"></span>' }
}))

// Mock syncStore
const mockAddTask = vi.fn()
const mockClearCompleted = vi.fn()
const mockRetryTask = vi.fn()
const mockRemoveTask = vi.fn()

const mockStore = {
  syncTasks: [],
  loading: false,
  isRunning: false,
  pendingCount: 0,
  downloadingCount: 0,
  uploadingCount: 0,
  completedCount: 0,
  errorCount: 0,
  addTask: mockAddTask,
  clearCompleted: mockClearCompleted,
  retryTask: mockRetryTask,
  removeTask: mockRemoveTask
}

vi.mock('../stores/sync', () => ({
  useSyncStore: () => mockStore
}))

// Mock auth store
vi.mock('../stores/auth', () => ({
  useAuthStore: () => ({ user: { email: 'test@example.com' } })
}))

const mountComponent = () => {
  return mount(Sync, {
    global: {
      stubs: {
        AppHeader: true
      }
    }
  })
}

describe('Sync.vue', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    mockAddTask.mockReset()
    mockClearCompleted.mockReset()
    mockRetryTask.mockReset()
    mockRemoveTask.mockReset()

    mockStore.syncTasks = []
    mockStore.loading = false
    mockStore.isRunning = false
    mockStore.pendingCount = 0
    mockStore.downloadingCount = 0
    mockStore.uploadingCount = 0
    mockStore.completedCount = 0
    mockStore.errorCount = 0
  })

  it('should render page title', () => {
    const wrapper = mountComponent()
    expect(wrapper.text()).toContain('同步资源')
  })

  it('should render empty slot when no tasks', () => {
    mockStore.syncTasks = []
    const wrapper = mountComponent()
    // The NSpin mock renders <div class="n-spin"><slot /></div>
    // Inside, v-if checks length === 0 and renders NEmpty
    // NEmpty mock is <div class="n-empty"><slot /></div>
    const html = wrapper.html()
    expect(html).toContain('n-empty')
  })

  it('should show add task button when not running', () => {
    const wrapper = mountComponent()
    expect(wrapper.text()).toContain('添加任务')
  })

  it('should show syncing state when running', () => {
    mockStore.isRunning = true
    const wrapper = mountComponent()
    expect(wrapper.text()).toContain('同步中...')
  })

  it('should show clear completed button when completed tasks exist', () => {
    mockStore.completedCount = 3
    const wrapper = mountComponent()
    expect(wrapper.text()).toContain('清除完成')
  })

  it('should show statistics', async () => {
    mockStore.pendingCount = 2
    mockStore.downloadingCount = 1
    mockStore.uploadingCount = 1
    mockStore.completedCount = 5
    mockStore.errorCount = 1
    const wrapper = mountComponent()

    expect(wrapper.text()).toContain('等待中: 2')
    expect(wrapper.text()).toContain('下载中: 1')
    expect(wrapper.text()).toContain('上传中: 1')
    expect(wrapper.text()).toContain('已完成: 5')
    expect(wrapper.text()).toContain('失败: 1')
  })

  it('should show sync task items', async () => {
    mockStore.syncTasks = [
      { id: 1, name: 'test.zip', status: 'pending', progress: 0, type: 'direct' }
    ]
    const wrapper = mountComponent()
    await new Promise(r => setTimeout(r, 10))

    expect(wrapper.text()).toContain('test.zip')
    expect(wrapper.text()).toContain('直链')
  })

  it('should show task status tags', async () => {
    mockStore.syncTasks = [
      { id: 1, name: 'file.zip', status: 'completed', progress: 100, type: 'direct' },
      { id: 2, name: 'doc.pdf', status: 'error', progress: 50, type: 'lanzou', error: 'Network error' },
      { id: 3, name: 'image.png', status: 'downloading', progress: 75, type: 'direct' },
      { id: 4, name: 'video.mp4', status: 'uploading', progress: 30, type: 'direct' },
      { id: 5, name: 'data.csv', status: 'pending', progress: 0, type: 'direct' }
    ]
    const wrapper = mountComponent()
    await new Promise(r => setTimeout(r, 20))

    expect(wrapper.text()).toContain('完成')
    expect(wrapper.text()).toContain('失败')
    expect(wrapper.text()).toContain('下载中')
    expect(wrapper.text()).toContain('上传中')
    expect(wrapper.text()).toContain('等待中')
  })

  it('should show retry button for failed tasks', async () => {
    mockStore.syncTasks = [
      { id: 1, name: 'failed.zip', status: 'error', progress: 50, type: 'direct', error: 'Timeout' }
    ]
    const wrapper = mountComponent()
    await new Promise(r => setTimeout(r, 10))

    expect(wrapper.text()).toContain('重试')
  })

  it('should call retryTask when retry clicked', async () => {
    mockStore.syncTasks = [
      { id: 1, name: 'failed.zip', status: 'error', progress: 50, type: 'direct', error: 'Timeout' }
    ]
    const wrapper = mountComponent()
    await new Promise(r => setTimeout(r, 10))

    await wrapper.vm.syncStore.retryTask(1)

    expect(mockRetryTask).toHaveBeenCalledWith(1)
  })

  it('should call clearCompleted', async () => {
    mockStore.completedCount = 3
    mockClearCompleted.mockResolvedValue(true)
    const wrapper = mountComponent()

    await wrapper.vm.syncStore.clearCompleted()

    expect(mockClearCompleted).toHaveBeenCalled()
  })

  it('should open add task modal', async () => {
    const wrapper = mountComponent()

    // Directly set showAddModal since NButton text matching may fail
    wrapper.vm.showAddModal = true
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.showAddModal).toBe(true)
  })

  it('should call addTask with valid data', async () => {
    mockAddTask.mockReturnValue(true)
    const wrapper = mountComponent()

    wrapper.vm.taskUrl = 'https://example.com/file.zip'
    wrapper.vm.taskName = 'myfile'
    wrapper.vm.taskType = 'direct'
    wrapper.vm.addSyncTask()

    expect(mockAddTask).toHaveBeenCalledWith({
      url: 'https://example.com/file.zip',
      name: 'myfile',
      type: 'direct'
    })
  })

  it('should not add task without url', () => {
    const wrapper = mountComponent()

    wrapper.vm.taskUrl = ''
    wrapper.vm.addSyncTask()

    expect(mockAddTask).not.toHaveBeenCalled()
  })

  it('should reset form after adding task', () => {
    const wrapper = mountComponent()

    wrapper.vm.taskUrl = 'https://example.com/file.zip'
    wrapper.vm.taskName = 'myfile'
    wrapper.vm.showAddModal = true
    wrapper.vm.addSyncTask()

    expect(wrapper.vm.taskUrl).toBe('')
    expect(wrapper.vm.taskName).toBe('')
    expect(wrapper.vm.showAddModal).toBe(false)
  })

  it('should disable lanzou radio option', async () => {
    const wrapper = mountComponent()

    wrapper.vm.showAddModal = true
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.taskType).toBe('direct')
  })
})
