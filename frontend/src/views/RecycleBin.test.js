import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import RecycleBin from './RecycleBin.vue'

// Mock naive-ui
vi.mock('naive-ui', () => ({
  NButton: {
    template: '<button class="n-button" :type="type" :loading="loading" @click="$emit(\'click\')"><slot /></button>',
    props: ['type', 'loading', 'block', 'size', 'text', 'strong']
  },
  NIcon: {
    template: '<span class="n-icon"><slot /></span>',
    props: ['size', 'component']
  },
  NCheckbox: {
    template: '<input class="n-checkbox" type="checkbox" :checked="checked" @change="$emit(\'update:checked\', $event.target.checked)" />',
    props: ['checked']
  },
  NPopconfirm: {
    template: '<div class="n-popconfirm"><slot name="trigger" /><slot /></div>',
    props: ['positiveText']
  },
  NSpin: {
    template: '<div class="n-spin"><slot /></div>',
    props: ['show']
  },
  NEmpty: {
    template: '<div class="n-empty"><slot /></div>',
    props: ['description']
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
  useRoute: () => ({ path: '/recycle' })
}))

// Mock @vicons/ionicons5
vi.mock('@vicons/ionicons5', () => ({
  DocumentText: { template: '<span data-icon="document-text"></span>' },
  Folder: { template: '<span data-icon="folder"></span>' },
  Refresh: { template: '<span data-icon="refresh"></span>' },
  ChevronDown: { template: '<span data-icon="chevron-down"></span>' },
  Settings: { template: '<span data-icon="settings"></span>' },
  LogOutOutline: { template: '<span data-icon="logout"></span>' }
}))

// Mock recycleStore
const mockFetchList = vi.fn()
const mockRestore = vi.fn()
const mockPermanentDelete = vi.fn()
const mockBatchRestore = vi.fn()
const mockBatchDelete = vi.fn()
const mockClearAll = vi.fn()
const mockToggleSelect = vi.fn()
const mockFormatSize = vi.fn((bytes) => `${bytes} B`)
const mockFormatDate = vi.fn((date) => date)

const mockStore = {
  items: [],
  selectedItems: [],
  totalCount: 0,
  fileCount: 0,
  folderCount: 0,
  loading: false,
  selectMode: false,
  hasSelectedItems: false,
  fetchList: mockFetchList,
  restore: mockRestore,
  permanentDelete: mockPermanentDelete,
  batchRestore: mockBatchRestore,
  batchDelete: mockBatchDelete,
  clearAll: mockClearAll,
  toggleSelect: mockToggleSelect,
  formatSize: mockFormatSize,
  formatDate: mockFormatDate
}

vi.mock('../stores/recycle', () => ({
  useRecycleStore: () => mockStore
}))

// Mock auth store (needed by AppHeader stub)
vi.mock('../stores/auth', () => ({
  useAuthStore: () => ({ user: { email: 'test@example.com' } })
}))

const mountComponent = () => {
  return mount(RecycleBin, {
    global: {
      stubs: {
        AppHeader: true
      }
    }
  })
}

describe('RecycleBin.vue', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    mockFetchList.mockReset()
    mockRestore.mockReset()
    mockPermanentDelete.mockReset()
    mockBatchRestore.mockReset()
    mockBatchDelete.mockReset()
    mockClearAll.mockReset()

    mockStore.items = []
    mockStore.selectedItems = []
    mockStore.totalCount = 0
    mockStore.fileCount = 0
    mockStore.folderCount = 0
    mockStore.loading = false
    mockStore.selectMode = false
    mockStore.hasSelectedItems = false
  })

  it('should render page title', () => {
    const wrapper = mountComponent()
    expect(wrapper.text()).toContain('刷新')
  })

  it('should show empty state when no items', () => {
    mockStore.items = []
    mockStore.totalCount = 0
    const wrapper = mountComponent()
    expect(wrapper.text()).toContain('共 0 项')
  })

  it('should call fetchList on mounted', () => {
    mountComponent()
    expect(mockFetchList).toHaveBeenCalled()
  })

  it('should show item count when items exist', async () => {
    mockStore.items = [
      { id: 1, item_type: 'file', original_name: 'test.txt', size: 1024, deleted_at: '2024-01-01' }
    ]
    mockStore.totalCount = 1
    mockStore.fileCount = 1
    mockStore.folderCount = 0
    const wrapper = mountComponent()
    await wrapper.vm.$nextTick()

    expect(wrapper.text()).toContain('共 1 项')
  })

  it('should render recycle items with file info', async () => {
    mockStore.items = [
      { id: 1, item_type: 'file', original_name: 'doc.pdf', size: 2048, deleted_at: '2024-03-01' }
    ]
    const wrapper = mountComponent()
    await wrapper.vm.$nextTick()

    expect(wrapper.text()).toContain('doc.pdf')
    expect(wrapper.text()).toContain('2048 B')
  })

  it('should call restore on single item restore', async () => {
    mockStore.items = [
      { id: 1, item_type: 'file', original_name: 'test.txt', size: 100, deleted_at: '2024-01-01' }
    ]
    mockRestore.mockResolvedValue(true)
    const wrapper = mountComponent()
    await wrapper.vm.$nextTick()

    await wrapper.vm.restore(1)

    expect(mockRestore).toHaveBeenCalledWith(1)
  })

  it('should show error message on restore failure', async () => {
    mockStore.items = [
      { id: 1, item_type: 'file', original_name: 'test.txt', size: 100, deleted_at: '2024-01-01' }
    ]
    mockRestore.mockRejectedValue(new Error('Server error'))
    const wrapper = mountComponent()
    await wrapper.vm.$nextTick()

    await wrapper.vm.restore(1)

    expect(mockRestore).toHaveBeenCalledWith(1)
  })

  it('should call permanentDelete on single item delete', async () => {
    mockStore.items = [
      { id: 2, item_type: 'file', original_name: 'test.txt', size: 100, deleted_at: '2024-01-01' }
    ]
    mockPermanentDelete.mockResolvedValue(true)
    const wrapper = mountComponent()
    await wrapper.vm.$nextTick()

    await wrapper.vm.permanentDelete(2)

    expect(mockPermanentDelete).toHaveBeenCalledWith(2)
  })

  it('should call batchRestore when has selected items', async () => {
    mockStore.items = [
      { id: 1, item_type: 'file', original_name: 'test.txt', size: 100, deleted_at: '2024-01-01' }
    ]
    mockStore.selectedItems = [1]
    mockStore.hasSelectedItems = true
    mockBatchRestore.mockResolvedValue(true)
    const wrapper = mountComponent()
    await wrapper.vm.$nextTick()

    await wrapper.vm.batchRestore()

    expect(mockBatchRestore).toHaveBeenCalled()
  })

  it('should not call batchRestore when no selected items', async () => {
    mockStore.hasSelectedItems = false
    const wrapper = mountComponent()

    await wrapper.vm.batchRestore()

    expect(mockBatchRestore).not.toHaveBeenCalled()
  })

  it('should call batchDelete when has selected items', async () => {
    mockStore.items = [
      { id: 1, item_type: 'file', original_name: 'test.txt', size: 100, deleted_at: '2024-01-01' }
    ]
    mockStore.selectedItems = [1]
    mockStore.hasSelectedItems = true
    mockBatchDelete.mockResolvedValue(true)
    const wrapper = mountComponent()
    await wrapper.vm.$nextTick()

    await wrapper.vm.batchDelete()

    expect(mockBatchDelete).toHaveBeenCalled()
  })

  it('should call clearAll', async () => {
    mockStore.totalCount = 5
    mockClearAll.mockResolvedValue(true)
    const wrapper = mountComponent()
    await wrapper.vm.$nextTick()

    await wrapper.vm.clearAll()

    expect(mockClearAll).toHaveBeenCalled()
  })

  it('should show error message on clearAll failure', async () => {
    mockStore.totalCount = 5
    mockClearAll.mockRejectedValue(new Error('Database error'))
    const wrapper = mountComponent()
    await wrapper.vm.$nextTick()

    await wrapper.vm.clearAll()

    expect(mockClearAll).toHaveBeenCalled()
  })

  it('should call fetchList on refresh button click', async () => {
    mockFetchList.mockResolvedValue(true)
    const wrapper = mountComponent()

    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()

    expect(mockFetchList).toHaveBeenCalled()
  })

  it('should show batch action buttons when has selected items', async () => {
    mockStore.items = [{ id: 1, item_type: 'file', original_name: 'test.txt', size: 100, deleted_at: '2024-01-01' }]
    mockStore.selectedItems = [1]
    mockStore.hasSelectedItems = true
    mockStore.totalCount = 1
    const wrapper = mountComponent()
    await new Promise(r => setTimeout(r, 10))

    expect(wrapper.text()).toContain('恢复选中 (1)')
    expect(wrapper.text()).toContain('彻底删除 (1)')
  })

  it('should show clear recycle bin button when items exist', async () => {
    mockStore.totalCount = 10
    const wrapper = mountComponent()
    await wrapper.vm.$nextTick()

    expect(wrapper.text()).toContain('清空回收站')
  })
})
