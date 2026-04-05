import { describe, it, expect, vi, beforeEach } from 'vitest'
import { shallowMount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'

vi.mock('naive-ui', () => ({
  NButton: { name: 'NButton', template: '<button><slot /></button>', props: ['type'] },
  NIcon: { name: 'NIcon', template: '<span class="n-icon"><slot /></span>', props: ['size'] }
}))

vi.mock('@vicons/ionicons5', () => ({
  Document: { template: '<span data-icon="document"></span>' },
  Folder: { template: '<span data-icon="folder"></span>' },
  Cloud: { template: '<span data-icon="cloud"></span>' },
  Close: { template: '<span data-icon="close"></span>' }
}))

vi.mock('../utils/fileSplit', () => ({
  formatFileSize: (bytes) => bytes ? `${bytes} B` : '0 B'
}))

const mockDownloadRemove = vi.fn()
const mockDownloadClear = vi.fn()
const mockUploadRemove = vi.fn()
const mockUploadClear = vi.fn()
const mockSyncRemove = vi.fn()
const mockSyncClear = vi.fn()

vi.mock('../stores/downloadTask', () => ({
  useDownloadTaskStore: vi.fn(() => ({
    tasks: [
      { id: 'd1', name: 'file1.zip', size: 1000, status: 'completed' },
      { id: 'd2', name: 'file2.zip', size: 2000, status: 'completed' },
      { id: 'd3', name: 'file3.zip', size: 3000, status: 'uploading' }
    ],
    removeTask: mockDownloadRemove,
    clearCompleted: mockDownloadClear
  }))
}))

vi.mock('../stores/uploadTask', () => ({
  useUploadTaskStore: vi.fn(() => ({
    tasks: [
      { id: 'u1', name: 'upload1.zip', size: 500, status: 'completed' },
      { id: 'u2', name: 'upload2.zip', size: 800, status: 'uploading' }
    ],
    removeTask: mockUploadRemove,
    clearCompleted: mockUploadClear
  }))
}))

vi.mock('../stores/sync', () => ({
  useSyncStore: vi.fn(() => ({
    tasks: [
      { id: 's1', name: 'sync1', url: 'https://example.com/file.zip', status: 'completed' }
    ],
    removeTask: mockSyncRemove,
    clearCompleted: mockSyncClear
  }))
}))

import CompletedTasks from './CompletedTasks.vue'

const mountComponent = () => {
  return shallowMount(CompletedTasks, {
    global: {
      stubs: {
        NButton: true,
        NIcon: true
      }
    }
  })
}

describe('CompletedTasks.vue', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('should render with download tab active by default', () => {
    const wrapper = mountComponent()
    expect(wrapper.vm.activeTab).toBe('download')
  })

  it('should show download completed tasks count', () => {
    const wrapper = mountComponent()
    expect(wrapper.vm.downloadTasks.length).toBe(2)
  })

  it('should show upload completed tasks count', () => {
    const wrapper = mountComponent()
    expect(wrapper.vm.uploadTasks.length).toBe(1)
  })

  it('should show sync completed tasks count', () => {
    const wrapper = mountComponent()
    expect(wrapper.vm.syncTasks.length).toBe(1)
  })

  it('should switch between tabs', async () => {
    const wrapper = mountComponent()
    wrapper.vm.activeTab = 'upload'
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.activeTab).toBe('upload')

    wrapper.vm.activeTab = 'sync'
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.activeTab).toBe('sync')
  })

  it('should format file sizes correctly', () => {
    const wrapper = mountComponent()
    expect(wrapper.vm.formatSize(0)).toBe('0 B')
    expect(wrapper.vm.formatSize(1000)).toBe('1000 B')
  })

  it('should remove download task by calling store', () => {
    const wrapper = mountComponent()
    wrapper.vm.removeDownload('d1')
    expect(mockDownloadRemove).toHaveBeenCalledWith('d1')
  })

  it('should remove upload task by calling store', () => {
    const wrapper = mountComponent()
    wrapper.vm.removeUpload('u1')
    expect(mockUploadRemove).toHaveBeenCalledWith('u1')
  })

  it('should remove sync task by calling store', () => {
    const wrapper = mountComponent()
    wrapper.vm.removeSync('s1')
    expect(mockSyncRemove).toHaveBeenCalledWith('s1')
  })

  it('should clear all download tasks', () => {
    const wrapper = mountComponent()
    wrapper.vm.activeTab = 'download'
    wrapper.vm.clearAll()
    expect(mockDownloadClear).toHaveBeenCalled()
  })

  it('should clear all upload tasks', () => {
    const wrapper = mountComponent()
    wrapper.vm.activeTab = 'upload'
    wrapper.vm.clearAll()
    expect(mockUploadClear).toHaveBeenCalled()
  })

  it('should clear all sync tasks', () => {
    const wrapper = mountComponent()
    wrapper.vm.activeTab = 'sync'
    wrapper.vm.clearAll()
    expect(mockSyncClear).toHaveBeenCalled()
  })
})
