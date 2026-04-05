import { describe, it, expect, vi, beforeEach, beforeAll } from 'vitest'
import { shallowMount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'

// Mock clipboard API at module level
const mockClipboardReadText = vi.fn().mockResolvedValue('')

beforeAll(() => {
  Object.defineProperty(window, 'navigator', {
    value: {
      ...window.navigator,
      clipboard: {
        readText: mockClipboardReadText
      }
    },
    writable: true,
    configurable: true
  })
})

// Mock vue-router
vi.mock('vue-router', () => ({
  useRouter: () => ({ push: vi.fn() }),
  useRoute: () => ({ path: '/share-parse' })
}))

// Mock @vicons/ionicons5
vi.mock('@vicons/ionicons5', () => ({
  DocumentText: { name: 'DocumentText' },
  Folder: { name: 'Folder' },
  CloudDownload: { name: 'CloudDownload' }
}))

// Mock api
const mockApiPost = vi.fn()
vi.mock('../api', () => ({
  default: {
    post: (...args) => mockApiPost(...args)
  }
}))

// Mock downloadTask store
const mockAddTask = vi.fn()
vi.mock('../stores/downloadTask', () => ({
  useDownloadTaskStore: () => ({
    addTask: mockAddTask
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

import ShareParse from './ShareParse.vue'

const mountComponent = () => {
  return shallowMount(ShareParse, {
    global: {
      stubs: {
        AppHeader: true,
        'n-input': { template: '<input />' },
        'n-button': { template: '<button><slot /></button>' },
        'n-icon': { template: '<span><slot /></span>' },
        'n-alert': { template: '<div class="n-alert"><slot /></div>' },
        'n-spin': { template: '<div class="n-spin"><slot /></div>' },
        'n-empty': { template: '<div class="n-empty"><slot /></div>' },
        'n-card': { template: '<div class="n-card"><slot /></div>' },
        'n-checkbox': { template: '<label><slot /></label>' }
      }
    }
  })
}

describe('ShareParse.vue', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    mockApiPost.mockReset()
    mockAddTask.mockReset()
    mockClipboardReadText.mockResolvedValue('')
  })

  it('should render page title', () => {
    const wrapper = mountComponent()
    expect(wrapper.text()).toContain('链接解析')
  })

  it('should have shareUrl and sharePwd refs', () => {
    const wrapper = mountComponent()
    expect(wrapper.vm.shareUrl).toBe('')
    expect(wrapper.vm.sharePwd).toBe('')
  })

  it('should initialize loading and error as false/empty', () => {
    const wrapper = mountComponent()
    expect(wrapper.vm.loading).toBe(false)
    expect(wrapper.vm.error).toBe('')
  })

  it('should show autoMerge ref as false', () => {
    const wrapper = mountComponent()
    expect(wrapper.vm.autoMerge).toBe(false)
  })

  it('should compute allFiles as empty array when no shares', () => {
    const wrapper = mountComponent()
    expect(wrapper.vm.allFiles).toEqual([])
  })

  it('should compute totalFiles as 0 when no shares', () => {
    const wrapper = mountComponent()
    expect(wrapper.vm.totalFiles).toBe(0)
  })

  it('should not parse when shareUrl is empty', async () => {
    const wrapper = mountComponent()
    await wrapper.vm.parseLinks()
    expect(mockApiPost).not.toHaveBeenCalled()
  })

  it('should parse single share link', async () => {
    const mockData = {
      name: 'test',
      type: 'file',
      size: '10MB',
      list: [{ url: 'https://d.0/', name: 'test.zip', size: '10MB', time: '2026-01-01' }]
    }
    mockApiPost.mockResolvedValue({ data: mockData })

    const wrapper = mountComponent()
    wrapper.vm.shareUrl = 'https://lanzoui.com/abc'
    wrapper.vm.sharePwd = '123'
    await wrapper.vm.parseLinks()

    expect(mockApiPost).toHaveBeenCalledWith('/lanzou/share/parse', {
      url: 'https://lanzoui.com/abc',
      pwd: '123'
    })
    expect(wrapper.vm.shareFiles).toHaveLength(1)
    expect(wrapper.vm.currentShare).toEqual(mockData)
  })

  it('should parse multiple lines', async () => {
    const mockData1 = { name: 'file1', type: 'file', list: [{ url: 'u1', name: 'f1.zip', size: '1MB', time: '' }] }
    const mockData2 = { name: 'file2', type: 'file', list: [{ url: 'u2', name: 'f2.zip', size: '2MB', time: '' }] }
    mockApiPost
      .mockResolvedValueOnce({ data: mockData1 })
      .mockResolvedValueOnce({ data: mockData2 })

    const wrapper = mountComponent()
    wrapper.vm.shareUrl = 'https://lanzoui.com/abc\nhttps://lanzoui.com/def'
    await wrapper.vm.parseLinks()

    expect(mockApiPost).toHaveBeenCalledTimes(2)
    expect(wrapper.vm.shareFiles).toHaveLength(2)
    expect(wrapper.vm.totalFiles).toBe(2)
  })

  it('should skip non-URL lines', async () => {
    mockApiPost.mockResolvedValue({ data: { name: 'test', type: 'file', list: [] } })

    const wrapper = mountComponent()
    wrapper.vm.shareUrl = 'https://lanzoui.com/abc\ninvalid text\nnot a url'
    await wrapper.vm.parseLinks()

    expect(mockApiPost).toHaveBeenCalledTimes(1)
  })

  it('should handle parse error gracefully', async () => {
    mockApiPost.mockRejectedValue(new Error('Network error'))

    const wrapper = mountComponent()
    wrapper.vm.shareUrl = 'https://lanzoui.com/bad'
    await wrapper.vm.parseLinks()

    expect(wrapper.vm.error).toContain('Network error')
    expect(wrapper.vm.loading).toBe(false)
  })

  it('should set currentShare only when single result', async () => {
    const mockData = { name: 'single', type: 'file', list: [{ url: 'u', name: 'f.zip', size: '', time: '' }] }
    mockApiPost.mockResolvedValue({ data: mockData })

    const wrapper = mountComponent()
    wrapper.vm.shareUrl = 'https://lanzoui.com/single'
    await wrapper.vm.parseLinks()

    expect(wrapper.vm.currentShare).toEqual(mockData)
  })

  it('should not set currentShare when multiple results', async () => {
    mockApiPost
      .mockResolvedValue({ data: { name: 'a', type: 'file', list: [] } })
      .mockResolvedValue({ data: { name: 'b', type: 'file', list: [] } })

    const wrapper = mountComponent()
    wrapper.vm.shareUrl = 'https://lanzoui.com/a\nhttps://lanzoui.com/b'
    await wrapper.vm.parseLinks()

    expect(wrapper.vm.currentShare).toBeNull()
  })

  it('should detect split file name and auto-set autoMerge', async () => {
    const mockData = {
      name: 'file.part001of3.zip',
      type: 'folder',
      list: [{ url: 'u', name: 'f.zip', size: '', time: '' }]
    }
    mockApiPost.mockResolvedValue({ data: mockData })

    const wrapper = mountComponent()
    wrapper.vm.shareUrl = 'https://lanzoui.com/split'
    await wrapper.vm.parseLinks()

    expect(wrapper.vm.autoMerge).toBe(true)
  })

  it('should not autoMerge for non-split files', async () => {
    const mockData = { name: 'normal.zip', type: 'folder', list: [{ url: 'u', name: 'f.zip', size: '', time: '' }] }
    mockApiPost.mockResolvedValue({ data: mockData })

    const wrapper = mountComponent()
    wrapper.vm.shareUrl = 'https://lanzoui.com/normal'
    await wrapper.vm.parseLinks()

    expect(wrapper.vm.autoMerge).toBe(false)
  })

  it('should compute allFiles from shareFiles', async () => {
    const mockData = {
      name: 'test',
      type: 'file',
      list: [
        { url: 'u1', name: 'f1.zip', size: '1MB', time: '2026-01-01' },
        { url: 'u2', name: 'f2.zip', size: '2MB', time: '2026-01-02' }
      ]
    }
    mockApiPost.mockResolvedValue({ data: mockData })

    const wrapper = mountComponent()
    wrapper.vm.shareUrl = 'https://lanzoui.com/abc'
    await wrapper.vm.parseLinks()

    expect(wrapper.vm.allFiles).toHaveLength(2)
    expect(wrapper.vm.allFiles[0].shareName).toBe('test')
    expect(wrapper.vm.totalFiles).toBe(2)
  })

  it('should toggleSelect add item', () => {
    const wrapper = mountComponent()
    const item = { url: 'u1', name: 'test.zip' }
    wrapper.vm.toggleSelect(item)
    expect(wrapper.vm.selectedRows).toHaveLength(1)
  })

  it('should toggleSelect remove item if already selected', () => {
    const wrapper = mountComponent()
    const item = { url: 'u1', name: 'test.zip' }
    wrapper.vm.toggleSelect(item)
    wrapper.vm.toggleSelect(item)
    expect(wrapper.vm.selectedRows).toHaveLength(0)
  })

  it('should isSelected return true for selected item', () => {
    const wrapper = mountComponent()
    const item = { url: 'u1', name: 'test.zip' }
    wrapper.vm.selectedRows.push(item)
    expect(wrapper.vm.isSelected(item)).toBe(true)
  })

  it('should downloadSingle add task to store', () => {
    const wrapper = mountComponent()
    const item = { name: 'test.zip', url: 'https://d.0/', size: '5MB', pwd: '' }
    wrapper.vm.downloadSingle(item)

    expect(mockAddTask).toHaveBeenCalledWith({
      name: 'test.zip',
      url: 'https://d.0/',
      pwd: '',
      size: 0,
      merge: false
    })
  })

  it('should downloadAll add all files to store', async () => {
    const mockData = {
      name: 'test',
      type: 'file',
      list: [
        { url: 'u1', name: 'f1.zip', size: '1MB', time: '' },
        { url: 'u2', name: 'f2.zip', size: '2MB', time: '' }
      ]
    }
    mockApiPost.mockResolvedValue({ data: mockData })

    const wrapper = mountComponent()
    wrapper.vm.shareUrl = 'https://lanzoui.com/abc'
    await wrapper.vm.parseLinks()
    wrapper.vm.downloadAll()

    expect(mockAddTask).toHaveBeenCalledTimes(2)
  })

  it('should downloadSelected add selected items to store', () => {
    const wrapper = mountComponent()
    wrapper.vm.selectedRows = [
      { url: 'u1', name: 'f1.zip', size: '1MB', pwd: '' },
      { url: 'u2', name: 'f2.zip', size: '2MB', pwd: '' }
    ]
    wrapper.vm.downloadSelected()

    expect(mockAddTask).toHaveBeenCalledTimes(2)
    expect(wrapper.vm.selectedRows).toHaveLength(0)
  })

  it('should clearResults reset all state', async () => {
    const mockData = { name: 'test', type: 'file', list: [{ url: 'u', name: 'f.zip', size: '', time: '' }] }
    mockApiPost.mockResolvedValue({ data: mockData })

    const wrapper = mountComponent()
    wrapper.vm.shareUrl = 'https://lanzoui.com/abc'
    wrapper.vm.sharePwd = '123'
    await wrapper.vm.parseLinks()
    wrapper.vm.error = 'some error'
    wrapper.vm.selectedRows.push({ url: 'u', name: 'f.zip' })

    wrapper.vm.clearResults()

    expect(wrapper.vm.shareFiles).toHaveLength(0)
    expect(wrapper.vm.currentShare).toBeNull()
    expect(wrapper.vm.selectedRows).toHaveLength(0)
    expect(wrapper.vm.error).toBe('')
  })

  it('should downloadSelected do nothing when no selection', () => {
    const wrapper = mountComponent()
    wrapper.vm.downloadSelected()
    expect(mockAddTask).not.toHaveBeenCalled()
  })

  it('should apply sharePwd to download tasks when item has no pwd', () => {
    const wrapper = mountComponent()
    wrapper.vm.sharePwd = 'sharedPwd'
    const item = { name: 'test.zip', url: 'https://d.0/', size: '0' }
    wrapper.vm.downloadSingle(item)

    expect(mockAddTask).toHaveBeenCalledWith(
      expect.objectContaining({ pwd: 'sharedPwd' })
    )
  })

  it('should apply autoMerge to downloadAll tasks', async () => {
    mockApiPost.mockResolvedValue({ data: { name: 'test', type: 'file', list: [{ url: 'u', name: 'f.zip', size: '', time: '' }] } })

    const wrapper = mountComponent()
    wrapper.vm.shareUrl = 'https://lanzoui.com/abc'
    wrapper.vm.autoMerge = true
    await wrapper.vm.parseLinks()
    wrapper.vm.downloadAll()

    expect(mockAddTask).toHaveBeenCalledWith(
      expect.objectContaining({ merge: true })
    )
  })
})
