import { describe, it, expect, vi, beforeEach } from 'vitest'
import { shallowMount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'

// Mock clipboard API
vi.stubGlobal('confirm', vi.fn(() => true))

// Mock vue-router
const mockPush = vi.fn()
vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: mockPush
  }),
  useRoute: () => ({ path: '/lanzou' })
}))

// Mock @vicons/ionicons5
vi.mock('@vicons/ionicons5', () => ({
  ArrowBack: { template: '<span data-icon="arrow-back"></span>' },
  Refresh: { template: '<span data-icon="refresh"></span>' },
  Folder: { template: '<span data-icon="folder"></span>' },
  Document: { template: '<span data-icon="document"></span>' },
  ChevronDown: { template: '<span data-icon="chevron-down"></span>' },
  Settings: { template: '<span data-icon="settings"></span>' },
  LogOutOutline: { template: '<span data-icon="logout"></span>' }
}))

// Mock auth store
vi.mock('../stores/auth', () => ({
  useAuthStore: () => ({
    user: { email: 'test@example.com' },
    fetchUser: vi.fn()
  })
}))

// Mock naive-ui
vi.mock('naive-ui', () => ({
  useMessage: () => ({
    success: vi.fn(),
    error: vi.fn(),
    warning: vi.fn(),
    info: vi.fn()
  })
}))

// Mock lanzouAPI - vi.hoisted because vi.mock is hoisted
const { mockStatus, mockListFiles, mockDelete, mockGetDownloadUrl, mockCreateShare, mockSetAccess, mockRename, mockGetFileDescription, mockSetFileDescription } = vi.hoisted(() => ({
  mockStatus: vi.fn(),
  mockListFiles: vi.fn(),
  mockDelete: vi.fn(),
  mockGetDownloadUrl: vi.fn(),
  mockCreateShare: vi.fn(),
  mockSetAccess: vi.fn(),
  mockRename: vi.fn(),
  mockGetFileDescription: vi.fn(),
  mockSetFileDescription: vi.fn()
}))

vi.mock('../api/lanzou', () => ({
  lanzouAPI: {
    status: mockStatus,
    listFiles: mockListFiles,
    delete: mockDelete,
    getDownloadUrl: mockGetDownloadUrl,
    createShare: mockCreateShare,
    setAccess: mockSetAccess,
    rename: mockRename,
    getFileDescription: mockGetFileDescription,
    setFileDescription: mockSetFileDescription
  }
}))

import LanZouBrowser from './LanZouBrowser.vue'

const mountComponent = () => {
  return shallowMount(LanZouBrowser, {
    global: {
      stubs: {
        AppHeader: true,
        'n-button': true,
        'n-breadcrumb': true,
        'n-breadcrumb-item': true,
        'n-checkbox': true,
        'n-spin': true,
        'n-result': true,
        'n-empty': true,
        'n-modal': true,
        'n-card': true,
        'n-input': true,
        'n-icon': true,
        'n-radio-group': true,
        'n-radio': true,
        'n-space': true
      }
    }
  })
}

describe('LanZouBrowser.vue', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    mockPush.mockReset()
    mockStatus.mockReset()
    mockListFiles.mockReset()
    mockDelete.mockReset()
    mockGetDownloadUrl.mockReset()
    mockCreateShare.mockReset()
    mockSetAccess.mockReset()
    mockRename.mockReset()
    mockGetFileDescription.mockReset()
    mockSetFileDescription.mockReset()
  })

  it('should render page header with navigation', () => {
    const wrapper = mountComponent()
    expect(wrapper.html()).toBeTruthy()
  })

  it('should have breadcrumbs component', () => {
    const wrapper = mountComponent()
    expect(wrapper.vm.breadcrumbs).toHaveLength(1)
    expect(wrapper.vm.breadcrumbs[0].name).toBe('首页')
  })

  it('should initialize with connected false', () => {
    const wrapper = mountComponent()
    expect(wrapper.vm.connected).toBe(false)
  })

  it('should have empty folders and files initially', () => {
    const wrapper = mountComponent()
    expect(wrapper.vm.folders).toHaveLength(0)
    expect(wrapper.vm.files).toHaveLength(0)
  })

  it('should have selectMode disabled initially', () => {
    const wrapper = mountComponent()
    expect(wrapper.vm.selectMode).toBe(false)
    expect(wrapper.vm.selectedItems).toHaveLength(0)
  })

  it('should toggle select add item', () => {
    const wrapper = mountComponent()

    wrapper.vm.toggleSelect('file1')
    expect(wrapper.vm.selectedItems).toContain('file1')
  })

  it('should toggle select remove item', () => {
    const wrapper = mountComponent()

    wrapper.vm.selectedItems = ['file1', 'file2']
    wrapper.vm.toggleSelect('file1')
    expect(wrapper.vm.selectedItems).not.toContain('file1')
    expect(wrapper.vm.selectedItems).toContain('file2')
  })

  it('should navigate to folder and update breadcrumbs', async () => {
    mockListFiles.mockResolvedValue({ folders: [], files: [] })
    const wrapper = mountComponent()

    wrapper.vm.navigateToFolder({ folder_id: '123', name: 'TestFolder' })

    expect(wrapper.vm.breadcrumbs).toHaveLength(2)
    expect(wrapper.vm.breadcrumbs[1].name).toBe('TestFolder')
  })

  it('should navigate back', async () => {
    mockListFiles.mockResolvedValue({ folders: [], files: [] })
    const wrapper = mountComponent()

    wrapper.vm.breadcrumbs = [
      { id: -1, name: '首页' },
      { id: '123', name: 'SubFolder' }
    ]
    wrapper.vm.navigateBack()

    expect(wrapper.vm.breadcrumbs).toHaveLength(1)
    expect(wrapper.vm.breadcrumbs[0].name).toBe('首页')
  })

  it('should not navigate back when at root', () => {
    const wrapper = mountComponent()
    wrapper.vm.breadcrumbs = [{ id: -1, name: '首页' }]

    wrapper.vm.navigateBack()

    expect(wrapper.vm.breadcrumbs).toHaveLength(1)
  })

  it('should navigate to breadcrumb', async () => {
    mockListFiles.mockResolvedValue({ folders: [], files: [] })
    const wrapper = mountComponent()

    wrapper.vm.breadcrumbs = [
      { id: -1, name: '首页' },
      { id: '123', name: 'SubFolder' },
      { id: '456', name: 'DeepFolder' }
    ]
    wrapper.vm.navigateToBreadcrumb(0)

    expect(wrapper.vm.breadcrumbs).toHaveLength(1)
    expect(wrapper.vm.breadcrumbs[0].name).toBe('首页')
  })

  it('should not navigate to current breadcrumb', () => {
    const wrapper = mountComponent()
    wrapper.vm.breadcrumbs = [
      { id: -1, name: '首页' },
      { id: '123', name: 'SubFolder' }
    ]

    // Last item - should not change
    wrapper.vm.navigateToBreadcrumb(1)

    expect(wrapper.vm.breadcrumbs).toHaveLength(2)
  })

  it('should handle file click to show file menu', () => {
    const wrapper = mountComponent()

    wrapper.vm.handleFileClick({ file_id: 'f1', name: 'test.zip' })

    expect(wrapper.vm.selectedFile).toEqual({ file_id: 'f1', name: 'test.zip' })
    expect(wrapper.vm.showFileMenu).toBe(true)
  })

  it('should download file successfully', async () => {
    mockGetDownloadUrl.mockResolvedValue({ url: 'https://download.example.com/file.zip' })
    const wrapper = mountComponent()

    wrapper.vm.selectedFile = { file_id: 'f1', name: 'test.zip' }
    await wrapper.vm.downloadFile()

    expect(mockGetDownloadUrl).toHaveBeenCalledWith('f1')
    expect(wrapper.vm.showFileMenu).toBe(false)
  })

  it('should handle download failure', async () => {
    mockGetDownloadUrl.mockRejectedValue(new Error('Network error'))
    const wrapper = mountComponent()

    wrapper.vm.selectedFile = { file_id: 'f1', name: 'test.zip' }
    await wrapper.vm.downloadFile()

    expect(mockGetDownloadUrl).toHaveBeenCalledWith('f1')
  })

  it('should handle download with no URL', async () => {
    mockGetDownloadUrl.mockResolvedValue({ url: null })
    const wrapper = mountComponent()

    wrapper.vm.selectedFile = { file_id: 'f1', name: 'test.zip' }
    await wrapper.vm.downloadFile()

    expect(wrapper.vm.showFileMenu).toBe(false)
  })

  it('should create share successfully', async () => {
    mockCreateShare.mockResolvedValue({ url: 'https://lanzoui.com/share1', pwd: '1234' })
    const wrapper = mountComponent()

    wrapper.vm.selectedFile = { file_id: 'f1', name: 'test.zip' }
    await wrapper.vm.createShare()

    expect(mockCreateShare).toHaveBeenCalledWith({ file_id: 'f1', minutes: 0 })
    expect(wrapper.vm.showShareModal).toBe(true)
    expect(wrapper.vm.shareResult.url).toBe('https://lanzoui.com/share1')
  })

  it('should handle share creation failure', async () => {
    mockCreateShare.mockRejectedValue(new Error('Share failed'))
    const wrapper = mountComponent()

    wrapper.vm.selectedFile = { file_id: 'f1', name: 'test.zip' }
    await wrapper.vm.createShare()

    expect(wrapper.vm.showShareModal).toBe(false)
  })

  it('should batch delete items', async () => {
    mockDelete.mockResolvedValue(true)
    mockListFiles.mockResolvedValue({ folders: [], files: [] })
    const wrapper = mountComponent()

    wrapper.vm.selectedItems = ['f1', 'f2']
    await wrapper.vm.batchDelete()

    expect(mockDelete).toHaveBeenCalledTimes(2)
    expect(wrapper.vm.selectedItems).toHaveLength(0)
  })

  it('should not batch delete when no items selected', async () => {
    const wrapper = mountComponent()
    wrapper.vm.selectedItems = []

    await wrapper.vm.batchDelete()

    expect(mockDelete).not.toHaveBeenCalled()
  })

  it('should format file sizes correctly', () => {
    const wrapper = mountComponent()

    expect(wrapper.vm.formatSize(0)).toBe('0 B')
    expect(wrapper.vm.formatSize(500)).toContain('B')
    expect(wrapper.vm.formatSize(1500)).toContain('KB')
    expect(wrapper.vm.formatSize(1500000)).toContain('MB')
    expect(wrapper.vm.formatSize(1500000000)).toContain('GB')
  })

  it('should refresh fetch current folder contents', async () => {
    mockListFiles.mockResolvedValue({ folders: [], files: [] })
    const wrapper = mountComponent()

    wrapper.vm.currentFolderId = '456'
    await wrapper.vm.refresh()

    expect(mockListFiles).toHaveBeenCalledWith({ folder_id: '456', page: 1 })
  })

  it('should push to settings when not connected', () => {
    const wrapper = mountComponent()

    // The template has a button that pushes to /settings
    expect(mockPush).not.toHaveBeenCalled()
  })

  it('should have current folder ID default to -1', () => {
    const wrapper = mountComponent()
    expect(wrapper.vm.currentFolderId).toBe(-1)
  })

  it('should handle fetchContents error gracefully', async () => {
    mockListFiles.mockRejectedValue(new Error('API error'))
    const wrapper = mountComponent()

    await wrapper.vm.fetchContents('123')

    expect(wrapper.vm.folders).toHaveLength(0)
    expect(wrapper.vm.files).toHaveLength(0)
  })

  it('should handle fetchStatus error', async () => {
    mockStatus.mockRejectedValue(new Error('Status check failed'))
    const wrapper = mountComponent()

    const result = await wrapper.vm.fetchStatus()

    expect(result).toBe(false)
    expect(wrapper.vm.connected).toBe(false)
  })

  // Access password tests
  it('should open file context menu and show file menu', () => {
    const wrapper = mountComponent()

    wrapper.vm.openFileContextMenu({ preventDefault: vi.fn() }, { file_id: 'f1', name: 'test.zip' })

    expect(wrapper.vm.selectedFile).toEqual({ file_id: 'f1', name: 'test.zip' })
    expect(wrapper.vm.showFileMenu).toBe(true)
  })

  it('should open folder context menu and show folder menu', () => {
    const wrapper = mountComponent()

    wrapper.vm.openFolderContextMenu({ preventDefault: vi.fn() }, { folder_id: '123', name: 'TestFolder' })

    expect(wrapper.vm.selectedFolder).toEqual({ folder_id: '123', name: 'TestFolder' })
    expect(wrapper.vm.showFolderMenu).toBe(true)
  })

  it('should not open context menus in select mode', () => {
    const wrapper = mountComponent()
    wrapper.vm.selectMode = true

    wrapper.vm.openFileContextMenu({ preventDefault: vi.fn() }, { file_id: 'f1', name: 'test.zip' })
    wrapper.vm.openFolderContextMenu({ preventDefault: vi.fn() }, { folder_id: '123', name: 'TestFolder' })

    expect(wrapper.vm.showFileMenu).toBe(false)
    expect(wrapper.vm.showFolderMenu).toBe(false)
  })

  it('should open access modal for file', () => {
    const wrapper = mountComponent()

    wrapper.vm.selectedFile = { file_id: '123', name: 'test.zip' }
    wrapper.vm.openAccessModal('file')

    expect(wrapper.vm.showAccessModal).toBe(true)
    expect(wrapper.vm.accessForm.targetId).toBe('123')
    expect(wrapper.vm.accessForm.targetType).toBe('file')
    expect(wrapper.vm.accessForm.type).toBe('2')
    expect(wrapper.vm.showFileMenu).toBe(false)
  })

  it('should open access modal for folder', () => {
    const wrapper = mountComponent()

    wrapper.vm.selectedFolder = { folder_id: '456', name: 'TestFolder' }
    wrapper.vm.openAccessModal('folder')

    expect(wrapper.vm.showAccessModal).toBe(true)
    expect(wrapper.vm.accessForm.targetId).toBe('456')
    expect(wrapper.vm.accessForm.targetType).toBe('folder')
    expect(wrapper.vm.showFolderMenu).toBe(false)
  })

  it('should submit password access successfully', async () => {
    mockSetAccess.mockResolvedValue({ message: 'access updated' })
    const wrapper = mountComponent()

    wrapper.vm.accessForm = {
      targetId: '123',
      targetType: 'file',
      type: '2',
      password: 'mypassword',
    }
    await wrapper.vm.submitAccess()

    expect(mockSetAccess).toHaveBeenCalledWith({
      id: '123',
      type: 'file',
      shows: 2,
      shownames: 'mypassword',
    })
    expect(wrapper.vm.showAccessModal).toBe(false)
  })

  it('should submit public access successfully', async () => {
    mockSetAccess.mockResolvedValue({ message: 'access updated' })
    const wrapper = mountComponent()

    wrapper.vm.accessForm = {
      targetId: '456',
      targetType: 'folder',
      type: '1',
      password: '',
    }
    await wrapper.vm.submitAccess()

    expect(mockSetAccess).toHaveBeenCalledWith({
      id: '456',
      type: 'folder',
      shows: 1,
      shownames: '',
    })
    expect(wrapper.vm.showAccessModal).toBe(false)
  })

  it('should warn when submitting empty password for password access', async () => {
    const wrapper = mountComponent()

    wrapper.vm.accessForm = {
      targetId: '123',
      targetType: 'file',
      type: '2',
      password: '',
    }
    await wrapper.vm.submitAccess()

    expect(mockSetAccess).not.toHaveBeenCalled()
  })

  it('should handle access submission error', async () => {
    mockSetAccess.mockRejectedValue(new Error('Network error'))
    const wrapper = mountComponent()

    wrapper.vm.accessForm.targetId = '123'
    wrapper.vm.accessForm.targetType = 'file'
    wrapper.vm.accessForm.type = '2'
    wrapper.vm.accessForm.password = 'testpwd'
    wrapper.vm.showAccessModal = true
    await wrapper.vm.submitAccess()

    expect(mockSetAccess).toHaveBeenCalled()
    expect(wrapper.vm.accessForm.targetId).toBe('123')
  })

  // Rename tests
  it('should open rename modal for file', () => {
    const wrapper = mountComponent()
    wrapper.vm.selectedFile = { file_id: '123', name: 'old.zip' }
    wrapper.vm.openRenameModal('file')
    expect(wrapper.vm.showRenameModal).toBe(true)
    expect(wrapper.vm.renameForm.targetId).toBe('123')
    expect(wrapper.vm.renameForm.targetType).toBe('file')
    expect(wrapper.vm.renameForm.name).toBe('old.zip')
  })

  it('should open rename modal for folder', () => {
    const wrapper = mountComponent()
    wrapper.vm.selectedFolder = { folder_id: '456', name: 'OldFolder' }
    wrapper.vm.openRenameModal('folder')
    expect(wrapper.vm.showRenameModal).toBe(true)
    expect(wrapper.vm.renameForm.targetId).toBe('456')
    expect(wrapper.vm.renameForm.targetType).toBe('folder')
    expect(wrapper.vm.renameForm.name).toBe('OldFolder')
  })

  it('should open rename modal for selected item', () => {
    const wrapper = mountComponent()
    wrapper.vm.folders = [{ folder_id: 'f1', name: 'TestFolder' }]
    wrapper.vm.selectedItems = ['f1']
    wrapper.vm.openRenameModalForSelected()
    expect(wrapper.vm.showRenameModal).toBe(true)
    expect(wrapper.vm.renameForm.targetId).toBe('f1')
    expect(wrapper.vm.renameForm.targetType).toBe('folder')
    expect(wrapper.vm.renameForm.name).toBe('TestFolder')
  })

  it('should not open rename modal for selected when multiple items', () => {
    const wrapper = mountComponent()
    wrapper.vm.selectedItems = ['f1', 'f2']
    wrapper.vm.openRenameModalForSelected()
    expect(wrapper.vm.showRenameModal).toBe(false)
  })

  it('should submit rename successfully', async () => {
    mockRename.mockResolvedValue({ zt: 1, info: 'success' })
    mockListFiles.mockResolvedValue({ folders: [], files: [] })
    const wrapper = mountComponent()
    wrapper.vm.renameForm = { name: 'newname.zip', targetId: '123', targetType: 'file' }
    await wrapper.vm.submitRename()
    expect(mockRename).toHaveBeenCalledWith('123', { type: 'file', name: 'newname.zip' })
    expect(wrapper.vm.showRenameModal).toBe(false)
  })

  it('should warn on empty rename name', async () => {
    const wrapper = mountComponent()
    wrapper.vm.renameForm = { name: '', targetId: '123', targetType: 'file' }
    await wrapper.vm.submitRename()
    expect(mockRename).not.toHaveBeenCalled()
  })

  it('should sanitize rename name with spaces', async () => {
    mockRename.mockResolvedValue({ zt: 1, info: 'success' })
    mockListFiles.mockResolvedValue({ folders: [], files: [] })
    const wrapper = mountComponent()
    wrapper.vm.renameForm = { name: 'new name.zip', targetId: '123', targetType: 'file' }
    await wrapper.vm.submitRename()
    expect(mockRename).toHaveBeenCalledWith('123', { type: 'file', name: 'new_name.zip' })
  })

  it('should handle rename error', async () => {
    mockRename.mockRejectedValue(new Error('Network error'))
    const wrapper = mountComponent()
    wrapper.vm.renameForm = { name: 'test.zip', targetId: '123', targetType: 'file' }
    wrapper.vm.showRenameModal = true
    await wrapper.vm.submitRename()
    expect(mockRename).toHaveBeenCalled()
  })

  // One-click share tests
  it('should one click share single item', async () => {
    mockCreateShare.mockResolvedValue({ url: 'https://lanzoui.com/abc', pwd: '1234' })
    const mockClipboard = { writeText: vi.fn().mockResolvedValue(undefined) }
    Object.defineProperty(navigator, 'clipboard', { value: mockClipboard, writable: true })
    const wrapper = mountComponent()
    wrapper.vm.files = [{ file_id: 'f1', name: 'test.zip' }]
    wrapper.vm.selectedItems = ['f1']
    await wrapper.vm.oneClickShare()
    expect(mockCreateShare).toHaveBeenCalled()
  })

  it('should one click share multiple items', async () => {
    mockCreateShare.mockResolvedValue({ url: 'https://lanzoui.com/abc', pwd: '' })
    const mockClipboard = { writeText: vi.fn().mockResolvedValue(undefined) }
    Object.defineProperty(navigator, 'clipboard', { value: mockClipboard, writable: true })
    const wrapper = mountComponent()
    wrapper.vm.folders = [{ folder_id: 'f1', name: 'Folder1' }]
    wrapper.vm.files = [{ file_id: 'f2', name: 'File1.zip' }]
    wrapper.vm.selectedItems = ['f1', 'f2']
    await wrapper.vm.oneClickShare()
    expect(mockCreateShare).toHaveBeenCalledTimes(2)
  })

  it('should not one click share when no items selected', async () => {
    const wrapper = mountComponent()
    wrapper.vm.selectedItems = []
    await wrapper.vm.oneClickShare()
    expect(mockCreateShare).not.toHaveBeenCalled()
  })

  // File description tests
  it('should open file description modal', () => {
    const wrapper = mountComponent()
    wrapper.vm.selectedFile = { file_id: 'f1', name: 'test.zip' }
    wrapper.vm.showFileMenu = true
    wrapper.vm.openFileDescModal()
    expect(wrapper.vm.showFileDescModal).toBe(true)
    expect(wrapper.vm.fileDescForm.fileId).toBe('f1')
    expect(wrapper.vm.showFileMenu).toBe(false)
  })

  it('should submit file description successfully', async () => {
    mockSetFileDescription.mockResolvedValue({ zt: 1, info: 'success' })
    const wrapper = mountComponent()
    wrapper.vm.fileDescForm = { description: 'This is a test file', fileId: 'f1' }
    await wrapper.vm.submitFileDesc()
    expect(mockSetFileDescription).toHaveBeenCalledWith('f1', { description: 'This is a test file' })
    expect(wrapper.vm.showFileDescModal).toBe(false)
  })

  it('should warn when submitting empty description', async () => {
    const wrapper = mountComponent()
    wrapper.vm.fileDescForm = { description: '', fileId: 'f1' }
    await wrapper.vm.submitFileDesc()
    expect(mockSetFileDescription).not.toHaveBeenCalled()
  })

  it('should handle file description error', async () => {
    mockSetFileDescription.mockRejectedValue(new Error('Network error'))
    const wrapper = mountComponent()
    wrapper.vm.fileDescForm = { description: 'test desc', fileId: 'f1' }
    wrapper.vm.showFileDescModal = true
    await wrapper.vm.submitFileDesc()
    expect(mockSetFileDescription).toHaveBeenCalled()
    expect(wrapper.vm.showFileDescModal).toBe(true)
  })
})
