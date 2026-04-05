import { describe, it, expect, vi, beforeEach } from 'vitest'
import { shallowMount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'

vi.mock('vue-router', () => ({
  useRouter: () => ({ push: vi.fn() }),
  useRoute: () => ({ path: '/upload-tasks' })
}))

vi.mock('naive-ui', () => ({
  useMessage: () => ({ success: vi.fn(), error: vi.fn(), warning: vi.fn() })
}))

vi.mock('@vicons/ionicons5', () => ({
  Document: { template: '<span data-icon="document"></span>' },
  Close: { template: '<span data-icon="close"></span>' }
}))

vi.mock('../stores/uploadTask', () => ({
  useUploadTaskStore: () => ({
    tasks: [
      { id: '1', name: 'file1.zip', size: 1024, status: 'uploading', progress: 50, speed: '1MB/s', error: null },
      { id: '2', name: 'file2.zip', size: 2048, status: 'pending', progress: 0, speed: '', error: null },
      { id: '3', name: 'file3.zip', size: 512, status: 'completed', progress: 100, speed: '', error: null },
      { id: '4', name: 'file4.zip', size: 3072, status: 'error', progress: 30, speed: '', error: '上传失败' },
    ],
    pendingCount: 1,
    activeCount: 1,
    completedCount: 1,
    errorCount: 1,
    pauseTask: vi.fn(),
    resumeTask: vi.fn(),
    pauseAll: vi.fn(),
    resumeAll: vi.fn(),
    removeTask: vi.fn(),
  })
}))

vi.mock('../components/AppHeader.vue', () => ({
  default: { template: '<div data-component="app-header"></div>' }
}))

import UploadTasks from './UploadTasks.vue'

describe('UploadTasks.vue', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('should render page header', () => {
    const wrapper = shallowMount(UploadTasks, {
      global: { stubs: { 'n-button': true, 'n-icon': true, 'n-tag': true, 'n-progress': true, 'n-empty': true } }
    })
    expect(wrapper.html()).toBeTruthy()
    expect(wrapper.html()).toContain('上传任务')
  })

  it('should have formatSize function', () => {
    const wrapper = shallowMount(UploadTasks, {
      global: { stubs: { 'n-button': true, 'n-icon': true, 'n-tag': true, 'n-progress': true, 'n-empty': true } }
    })
    expect(wrapper.vm.formatSize(0)).toBe('0 B')
    expect(wrapper.vm.formatSize(1024)).toContain('KB')
    expect(wrapper.vm.formatSize(1048576)).toContain('MB')
    expect(wrapper.vm.formatSize(1073741824)).toContain('GB')
  })

  it('should have toggleTask function', () => {
    const wrapper = shallowMount(UploadTasks, {
      global: { stubs: { 'n-button': true, 'n-icon': true, 'n-tag': true, 'n-progress': true, 'n-empty': true } }
    })
    expect(typeof wrapper.vm.toggleTask).toBe('function')
  })

  it('should render task items', () => {
    const wrapper = shallowMount(UploadTasks, {
      global: { stubs: { 'n-button': true, 'n-icon': true, 'n-tag': true, 'n-progress': true, 'n-empty': true } }
    })
    expect(wrapper.html()).toContain('file1.zip')
    expect(wrapper.html()).toContain('file2.zip')
  })

  it('should render error message for failed task', () => {
    const wrapper = shallowMount(UploadTasks, {
      global: { stubs: { 'n-button': true, 'n-icon': true, 'n-tag': true, 'n-progress': true, 'n-empty': true } }
    })
    expect(wrapper.html()).toContain('上传失败')
  })
})
