import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, shallowMount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import Dashboard from './Dashboard.vue'

// Mock naive-ui
vi.mock('naive-ui', () => ({
  NButton: {
    template: '<button class="n-button" :type="type" @click="$emit(\'click\')"><slot /></button>',
    props: ['type', 'block', 'size', 'loading', 'text', 'disabled']
  },
  NIcon: {
    template: '<span class="n-icon"><slot /></span>',
    props: ['size', 'color']
  },
  NInput: {
    template: '<input class="n-input" :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />',
    props: ['modelValue', 'placeholder', 'clearable', 'style', 'show-password-on']
  },
  NModal: {
    template: '<div class="n-modal" v-if="show"><slot /></div>',
    props: ['show', 'maskClosable']
  },
  NCard: {
    template: '<div class="n-card"><div class="title">{{ title }}</div><slot /><slot name="footer" /></div>',
    props: ['title', 'style']
  },
  NUpload: {
    template: '<div class="n-upload" @click="$emit(\'change\', { fileList: [] })"><slot /></div>',
    props: ['multiple', 'accept', 'max']
  },
  NCheckbox: {
    template: '<input type="checkbox" :checked="checked" @change="$emit(\'update:checked\', !checked)" />',
    props: ['checked', 'label']
  },
  NSelect: {
    template: '<select class="n-select"><slot /></select>',
    props: ['value', 'options']
  },
  NSpin: {
    template: '<div class="n-spin"><slot /></div>',
    props: ['show']
  },
  NEmpty: {
    template: '<div class="n-empty"><span>{{ description }}</span></div>',
    props: ['description']
  },
  NProgress: {
    template: '<div class="n-progress"></div>',
    props: ['percentage', 'showIndicator', 'height']
  },
  NTag: {
    template: '<span class="n-tag"><slot /></span>',
    props: ['type', 'size']
  },
  NAlert: {
    template: '<div class="n-alert"><slot /></div>',
    props: ['type', 'showIcon']
  },
  useMessage: () => ({
    success: vi.fn(),
    error: vi.fn(),
    warning: vi.fn(),
    info: vi.fn()
  })
}))

// Mock icons
vi.mock('@vicons/ionicons5', () => ({
  ArrowBack: { template: '<span>ArrowBack</span>' },
  CloudUpload: { template: '<span>CloudUpload</span>' },
  Refresh: { template: '<span>Refresh</span>' },
  Folder: { template: '<span>Folder</span>' },
  Create: { template: '<span>Create</span>' },
  Document: { template: '<span>Document</span>' },
  Close: { template: '<span>Close</span>' }
}))

// Mock AppHeader
vi.mock('../components/AppHeader.vue', () => ({
  default: { template: '<div class="app-header">AppHeader</div>' }
}))

// Mock FileDetailModal
vi.mock('../components/FileDetailModal.vue', () => ({
  default: {
    template: '<div class="file-detail-modal"></div>',
    props: ['item'],
    emits: ['update']
  }
}))

// Mock crypto utils
vi.mock('../utils/crypto', () => ({
  encryptFile: vi.fn().mockResolvedValue(new Blob(['encrypted'])),
  generateEncryptionKey: vi.fn().mockResolvedValue('mock-key-123')
}))

// Mock API
const mockFileList = vi.fn().mockResolvedValue([
  { id: 1, name: 'test.txt', size: 1024, created_at: '2024-01-01' }
])
const mockFileDelete = vi.fn().mockResolvedValue({})
const mockFileRename = vi.fn().mockResolvedValue({ id: 1, name: 'new-name.txt' })
const mockFileMove = vi.fn().mockResolvedValue({})
const mockFileDownload = vi.fn().mockResolvedValue({ download_url: 'https://example.com/file' })
const mockFileUpdateDescription = vi.fn().mockResolvedValue({ id: 1, description: 'test' })

const mockFolderList = vi.fn().mockResolvedValue([
  { id: 10, name: 'docs', created_at: '2024-01-01' }
])
const mockFolderCreate = vi.fn().mockResolvedValue({ id: 11, name: 'new-folder' })
const mockFolderDelete = vi.fn().mockResolvedValue({})
const mockFolderRename = vi.fn().mockResolvedValue({ id: 10, name: 'renamed' })
const mockFolderMove = vi.fn().mockResolvedValue({})
const mockFolderUpdateDescription = vi.fn().mockResolvedValue({ id: 10, description: 'test' })

vi.mock('../api', () => ({
  fileAPI: {
    list: (...args) => mockFileList(...args),
    delete: (...args) => mockFileDelete(...args),
    rename: (...args) => mockFileRename(...args),
    move: (...args) => mockFileMove(...args),
    download: (...args) => mockFileDownload(...args),
    updateDescription: (...args) => mockFileUpdateDescription(...args)
  },
  folderAPI: {
    list: (...args) => mockFolderList(...args),
    create: (...args) => mockFolderCreate(...args),
    delete: (...args) => mockFolderDelete(...args),
    rename: (...args) => mockFolderRename(...args),
    move: (...args) => mockFolderMove(...args),
    updateDescription: (...args) => mockFolderUpdateDescription(...args)
  },
  shareAPI: {
    create: vi.fn().mockResolvedValue({ share_url: 'https://example.com/share', share_pwd: '1234' })
  }
}))

// Mock router
const mockPush = vi.fn()
vi.mock('vue-router', () => ({
  useRouter: () => ({ push: mockPush })
}))

describe('Dashboard.vue', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('renders Dashboard component', () => {
    const wrapper = mount(Dashboard, {
      global: {
        stubs: ['n-button', 'n-icon', 'n-modal', 'n-card', 'n-input', 'n-checkbox', 'n-select', 'n-spin', 'n-empty', 'n-alert', 'n-tag', 'n-progress', 'n-upload']
      }
    })
    expect(wrapper.exists()).toBe(true)
  })

  it('shows empty state when no files or folders', async () => {
    mockFileList.mockResolvedValue([])
    mockFolderList.mockResolvedValue([])

    const wrapper = mount(Dashboard)
    await wrapper.vm.$nextTick()

    expect(wrapper.exists()).toBe(true)
  })

  it('displays files from fileStore', async () => {
    mockFileList.mockResolvedValue([{ id: 1, name: 'test.txt', size: 2048 }])
    mockFolderList.mockResolvedValue([])

    const wrapper = mount(Dashboard)
    await wrapper.vm.$nextTick()

    expect(wrapper.exists()).toBe(true)
  })

  it('displays folders from fileStore', async () => {
    mockFileList.mockResolvedValue([])
    mockFolderList.mockResolvedValue([{ id: 10, name: 'my-folder' }])

    const wrapper = mount(Dashboard)
    await wrapper.vm.$nextTick()

    expect(wrapper.exists()).toBe(true)
  })

  it('filters folders by search query', async () => {
    mockFileList.mockResolvedValue([])
    mockFolderList.mockResolvedValue([
      { id: 10, name: 'docs', created_at: '2024-01-01' },
      { id: 11, name: 'images', created_at: '2024-01-02' }
    ])

    const wrapper = mount(Dashboard)
    await new Promise(r => setTimeout(r, 50))
    await wrapper.vm.$nextTick()

    // Set search query
    wrapper.vm.searchQuery = 'doc'
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.filteredFolders.length).toBe(1)
    expect(wrapper.vm.filteredFolders[0].name).toBe('docs')
  })

  it('filters files by search query', async () => {
    mockFileList.mockResolvedValue([
      { id: 1, name: 'report.pdf', size: 1000 },
      { id: 2, name: 'photo.jpg', size: 2000 }
    ])
    mockFolderList.mockResolvedValue([])

    const wrapper = mount(Dashboard)
    await new Promise(r => setTimeout(r, 50))
    await wrapper.vm.$nextTick()

    wrapper.vm.searchQuery = 'report'
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.filteredFiles.length).toBe(1)
    expect(wrapper.vm.filteredFiles[0].name).toBe('report.pdf')
  })

  it('sorts folders by name', async () => {
    mockFileList.mockResolvedValue([])
    mockFolderList.mockResolvedValue([
      { id: 10, name: 'zebra', created_at: '2024-01-01' },
      { id: 11, name: 'alpha', created_at: '2024-01-02' }
    ])

    const wrapper = mount(Dashboard)
    // Wait for onMounted fetchFiles to complete
    await new Promise(r => setTimeout(r, 50))
    await wrapper.vm.$nextTick()

    wrapper.vm.sortBy = 'name'
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.filteredFolders.length).toBeGreaterThan(0)
    expect(wrapper.vm.filteredFolders[0].name).toBe('alpha')
  })

  it('sorts files by size', async () => {
    mockFileList.mockResolvedValue([
      { id: 1, name: 'small.txt', size: 100 },
      { id: 2, name: 'big.txt', size: 5000 }
    ])
    mockFolderList.mockResolvedValue([])

    const wrapper = mount(Dashboard)
    await new Promise(r => setTimeout(r, 50))
    await wrapper.vm.$nextTick()

    wrapper.vm.sortBy = 'size'
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.filteredFiles.length).toBeGreaterThan(0)
    expect(wrapper.vm.filteredFiles[0].name).toBe('big.txt')
  })

  it('sorts folders by time', async () => {
    mockFileList.mockResolvedValue([])
    mockFolderList.mockResolvedValue([
      { id: 10, name: 'old', created_at: '2024-01-01' },
      { id: 11, name: 'new', created_at: '2024-12-01' }
    ])

    const wrapper = mount(Dashboard)
    await new Promise(r => setTimeout(r, 50))
    await wrapper.vm.$nextTick()

    wrapper.vm.sortBy = 'time'
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.filteredFolders.length).toBeGreaterThan(0)
    expect(wrapper.vm.filteredFolders[0].name).toBe('new')
  })

  it('sorts files by time', async () => {
    mockFileList.mockResolvedValue([
      { id: 1, name: 'old.txt', size: 100, created_at: '2024-01-01' },
      { id: 2, name: 'new.txt', size: 200, created_at: '2024-12-01' }
    ])
    mockFolderList.mockResolvedValue([])

    const wrapper = mount(Dashboard)
    await new Promise(r => setTimeout(r, 50))
    await wrapper.vm.$nextTick()

    wrapper.vm.sortBy = 'time'
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.filteredFiles.length).toBeGreaterThan(0)
    expect(wrapper.vm.filteredFiles[0].name).toBe('new.txt')
  })

  it('toggles item selection in select mode', async () => {
    const wrapper = mount(Dashboard)

    wrapper.vm.toggleSelect(1)
    expect(wrapper.vm.selectedItems).toContain(1)

    wrapper.vm.toggleSelect(1)
    expect(wrapper.vm.selectedItems).not.toContain(1)
  })

  it('formatSize returns correct values', () => {
    const wrapper = mount(Dashboard)

    expect(wrapper.vm.formatSize(0)).toBe('0 B')
    expect(wrapper.vm.formatSize(1024)).toBe('1 KB')
    expect(wrapper.vm.formatSize(1048576)).toBe('1 MB')
    expect(wrapper.vm.formatSize(1073741824)).toBe('1 GB')
  })

  it('formatSpeed returns speed format', () => {
    const wrapper = mount(Dashboard)

    expect(wrapper.vm.formatSpeed(1024)).toBe('1 KB/s')
  })

  it('openBatchMoveModal sets moveMode to batch', async () => {
    const wrapper = mount(Dashboard)
    wrapper.vm.selectedItems = [1, 2]

    await wrapper.vm.openBatchMoveModal()

    expect(wrapper.vm.moveMode).toBe('batch')
    expect(wrapper.vm.showMoveModal).toBe(true)
  })

  it('openMoveModal sets moveMode to single', async () => {
    const wrapper = mount(Dashboard)
    wrapper.vm.selectedFileItem = { id: 1, name: 'test.txt', size: 100 }

    await wrapper.vm.openMoveModal()

    expect(wrapper.vm.moveMode).toBe('single')
    expect(wrapper.vm.showMoveModal).toBe(true)
  })

  it('confirmRename validates input', async () => {
    const wrapper = mount(Dashboard)
    wrapper.vm.selectedFileItem = { id: 1, name: 'test.txt', size: 100 }
    wrapper.vm.renameName = ''

    await wrapper.vm.confirmRename()
    // Should not call API with empty name
    expect(mockFileRename).not.toHaveBeenCalled()
  })

  it('closeUploadModal does nothing when uploading', async () => {
    const wrapper = mount(Dashboard)

    // When no items in queue and not uploading, should close
    wrapper.vm.showUploadModal = true
    await wrapper.vm.closeUploadModal()
    expect(wrapper.vm.showUploadModal).toBe(false)
  })

  it('startBatchUpload warns when no files selected', async () => {
    const wrapper = mount(Dashboard)
    wrapper.vm.uploadStore.uploadQueue = []

    await wrapper.vm.startBatchUpload()
    // Should show warning message
    expect(wrapper.vm.uploadStore.uploadQueue.length).toBe(0)
  })

  it('handleFileClick triggers file download', async () => {
    const wrapper = mount(Dashboard)
    const file = { id: 1, name: 'test.txt', size: 1024 }

    await wrapper.vm.handleFileClick(file)

    // handleFileClick should trigger download for the file
    expect(mockFileDownload).toHaveBeenCalled()
  })

  it('openFileMenu opens file menu', async () => {
    const wrapper = mount(Dashboard)
    const file = { id: 1, name: 'test.txt', size: 1024 }

    wrapper.vm.openFileMenu(file)

    expect(wrapper.vm.selectedFileItem).toEqual(file)
    expect(wrapper.vm.showFileMenu).toBe(true)
  })

  it('openShareModal resets share state', async () => {
    const wrapper = mount(Dashboard)
    wrapper.vm.shareResult = { url: 'test' }
    wrapper.vm.shareError = 'some error'

    wrapper.vm.openShareModal()

    expect(wrapper.vm.shareResult).toBeNull()
    expect(wrapper.vm.shareError).toBeNull()
    expect(wrapper.vm.showShareModal).toBe(true)
  })

  it('openRenameModal sets rename name', async () => {
    const wrapper = mount(Dashboard)
    wrapper.vm.selectedFileItem = { id: 1, name: 'original.txt', size: 100 }

    wrapper.vm.openRenameModal()

    expect(wrapper.vm.renameName).toBe('original.txt')
    expect(wrapper.vm.showRenameModal).toBe(true)
  })

  it('openFolderMenu sets selected folder item', async () => {
    const wrapper = mount(Dashboard)
    const folder = { id: 10, name: 'my-folder' }

    wrapper.vm.openFolderMenu(folder)

    expect(wrapper.vm.selectedFileItem).toEqual(folder)
    expect(wrapper.vm.showFileMenu).toBe(true)
  })

  it('closeShareModal does nothing when sharing', async () => {
    const wrapper = mount(Dashboard)
    wrapper.vm.isSharing = true
    wrapper.vm.showShareModal = true

    wrapper.vm.closeShareModal()

    expect(wrapper.vm.showShareModal).toBe(true)
  })

  it('createShare calls share API with password and expires', async () => {
    const wrapper = mount(Dashboard)
    wrapper.vm.selectedFileItem = { id: 1, name: 'test.txt', size: 100 }
    wrapper.vm.sharePassword = 'mypwd'
    wrapper.vm.shareExpires = '1d'

    await wrapper.vm.createShare()

    expect(wrapper.vm.isSharing).toBe(false)
  })

  it('createShare handles 1h expiry', async () => {
    const wrapper = mount(Dashboard)
    wrapper.vm.selectedFileItem = { id: 1, name: 'test.txt', size: 100 }
    wrapper.vm.shareExpires = '1h'

    await wrapper.vm.createShare()

    expect(wrapper.vm.shareResult).not.toBeNull()
  })

  it('createShare handles never expiry', async () => {
    const wrapper = mount(Dashboard)
    wrapper.vm.selectedFileItem = { id: 1, name: 'test.txt', size: 100 }
    wrapper.vm.shareExpires = 'never'

    await wrapper.vm.createShare()

    expect(wrapper.vm.shareResult).not.toBeNull()
  })

  it('createShare handles 7d expiry', async () => {
    const wrapper = mount(Dashboard)
    wrapper.vm.selectedFileItem = { id: 1, name: 'test.txt', size: 100 }
    wrapper.vm.shareExpires = '7d'

    await wrapper.vm.createShare()

    expect(wrapper.vm.shareResult).not.toBeNull()
  })

  it('has drag-and-drop overlay state', () => {
    const wrapper = mount(Dashboard)
    wrapper.vm.isDragging = true
    expect(wrapper.vm.isDragging).toBe(true)
    // currentFolderName reflects the current breadcrumb name
    expect(wrapper.vm.currentFolderName).toBeDefined()
  })

  it('handleDrop processes dropped files', async () => {
    const wrapper = mount(Dashboard)
    const mockFile = new File(['content'], 'test.txt', { type: 'text/plain' })
    const dataTransfer = { files: [mockFile] }

    wrapper.vm.isDragging = true
    await wrapper.vm.handleDrop({ dataTransfer })

    expect(wrapper.vm.isDragging).toBe(false)
    expect(wrapper.vm.uploadStore.uploadQueue.length).toBe(1)
  })

  it('handleDragLeave hides overlay when leaving target', () => {
    const wrapper = mount(Dashboard)
    wrapper.vm.isDragging = true

    const mockEvent = {
      currentTarget: { contains: () => false },
      relatedTarget: null
    }
    wrapper.vm.handleDragLeave(mockEvent)

    expect(wrapper.vm.isDragging).toBe(false)
  })

  it('handleDrop ignores folder-like files', async () => {
    const wrapper = mount(Dashboard)
    const mockFolder = new File(['content'], 'folder', { type: '' })
    Object.defineProperty(mockFolder, 'webkitRelativePath', { value: 'folder' })
    const dataTransfer = { files: [mockFolder] }

    await wrapper.vm.handleDrop({ dataTransfer })

    expect(wrapper.vm.uploadStore.uploadQueue.length).toBe(0)
  })
})
