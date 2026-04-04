import { describe, it, expect, vi, beforeEach } from 'vitest'
import { shallowMount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'

// Mock clipboard API
const mockClipboardReadText = vi.fn().mockResolvedValue('')

// Mock vue-router
vi.mock('vue-router', () => ({
  useRouter: () => ({ push: vi.fn() }),
  useRoute: () => ({ path: '/share-parse' })
}))

// Mock @vicons/ionicons5
vi.mock('@vicons/ionicons5', () => ({
  DocumentText: { template: '<span data-icon="document-text"></span>' },
  Folder: { template: '<span data-icon="folder"></span>' },
  CloudDownload: { template: '<span data-icon="cloud-download"></span>' },
  ChevronDown: { template: '<span data-icon="chevron-down"></span>' },
  Settings: { template: '<span data-icon="settings"></span>' },
  LogOutOutline: { template: '<span data-icon="logout"></span>' }
}))

// Mock auth store
vi.mock('../stores/auth', () => ({
  useAuthStore: () => ({ user: { email: 'test@example.com' } })
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

// Mock shareParseStore
const mockParseShare = vi.fn()
const mockDownloadShareFile = vi.fn()
const mockClearParsed = vi.fn()

const mockStore = {
  loading: false,
  error: null,
  parsedShares: [],
  currentShare: null,
  hasParsedShares: false,
  parseShare: mockParseShare,
  downloadShareFile: mockDownloadShareFile,
  clearParsed: mockClearParsed
}

vi.mock('../stores/shareParse', () => ({
  useShareParseStore: () => mockStore
}))

import ShareParse from './ShareParse.vue'

const mountComponent = () => {
  return shallowMount(ShareParse, {
    global: {
      stubs: {
        AppHeader: true,
        'n-input': true,
        'n-button': true,
        'n-icon': true,
        'n-alert': true,
        'n-spin': true,
        'n-empty': true,
        'n-card': true
      }
    }
  })
}

describe('ShareParse.vue', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    mockParseShare.mockReset()
    mockDownloadShareFile.mockReset()
    mockClearParsed.mockReset()
    mockClipboardReadText.mockReset().mockResolvedValue('')

    mockStore.loading = false
    mockStore.error = null
    mockStore.parsedShares = []
    mockStore.currentShare = null
    mockStore.hasParsedShares = false
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

  it('should show parse button in stub', () => {
    const wrapper = mountComponent()
    expect(wrapper.html()).toContain('n-button')
  })

  it('should show instructions card in stub', () => {
    const wrapper = mountComponent()
    expect(wrapper.html()).toContain('n-card')
  })

  it('should show loading spin when store loading', () => {
    mockStore.loading = true
    const wrapper = mountComponent()
    expect(wrapper.html()).toContain('n-spin')
  })

  it('should show error alert when store has error', () => {
    mockStore.error = '解析失败'
    const wrapper = mountComponent()
    expect(wrapper.html()).toContain('n-alert')
  })

  it('should render stubs for naive-ui components', () => {
    const wrapper = mountComponent()
    // shallowMount renders stubs, n-empty is stubbed
    expect(wrapper.html()).toBeTruthy()
  })

  it('should access store state correctly', () => {
    mockStore.error = '测试错误'
    mockStore.loading = true
    const wrapper = mountComponent()

    expect(wrapper.vm.shareParseStore.error).toBe('测试错误')
    expect(wrapper.vm.shareParseStore.loading).toBe(true)
  })

  it('should parse single link', async () => {
    mockParseShare.mockResolvedValue(true)
    const wrapper = mountComponent()

    wrapper.vm.shareUrl = 'https://lanzoui.com/abc123'
    wrapper.vm.sharePwd = 'pwd'
    await wrapper.vm.parseLinks()

    expect(mockClearParsed).toHaveBeenCalled()
    expect(mockParseShare).toHaveBeenCalledWith('https://lanzoui.com/abc123', 'pwd')
  })

  it('should parse multiple lines', async () => {
    mockParseShare.mockResolvedValue(true)
    const wrapper = mountComponent()

    wrapper.vm.shareUrl = 'https://lanzoui.com/link1\nhttps://lanzoui.com/link2'
    wrapper.vm.sharePwd = ''
    await wrapper.vm.parseLinks()

    expect(mockParseShare).toHaveBeenCalledTimes(2)
  })

  it('should skip non-URL lines', async () => {
    mockParseShare.mockResolvedValue(true)
    const wrapper = mountComponent()

    wrapper.vm.shareUrl = 'https://lanzoui.com/link1\ninvalid text\nhttps://lanzoui.com/link2'
    wrapper.vm.sharePwd = ''
    await wrapper.vm.parseLinks()

    expect(mockParseShare).toHaveBeenCalledTimes(2)
  })

  it('should handle parse failure', async () => {
    mockParseShare.mockRejectedValue(new Error('Invalid share link'))
    const wrapper = mountComponent()

    wrapper.vm.shareUrl = 'https://lanzoui.com/invalid'
    await wrapper.vm.parseLinks()

    expect(mockParseShare).toHaveBeenCalled()
  })

  it('should download single file', async () => {
    mockDownloadShareFile.mockResolvedValue('https://download.example.com/file.zip')
    const windowOpen = vi.fn()
    const originalOpen = window.open
    window.open = windowOpen

    const wrapper = mountComponent()
    const item = { name: 'test.zip', url: 'https://lanzoui.com/abc', size: '5MB' }
    await wrapper.vm.downloadSingle(item)

    expect(mockDownloadShareFile).toHaveBeenCalledWith('https://lanzoui.com/abc', '')
    window.open = originalOpen
  })

  it('should handle download failure', async () => {
    mockDownloadShareFile.mockRejectedValue(new Error('Download failed'))
    const wrapper = mountComponent()

    const item = { name: 'test.zip', url: 'https://lanzoui.com/bad', size: '5MB' }
    await wrapper.vm.downloadSingle(item)

    expect(mockDownloadShareFile).toHaveBeenCalled()
  })

  it('should call downloadAll without error', () => {
    const wrapper = mountComponent()
    expect(() => wrapper.vm.downloadAll()).not.toThrow()
  })
})
