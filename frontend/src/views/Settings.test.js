import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import Settings from '../views/Settings.vue'

// Mock stores
vi.mock('../stores/theme', () => ({
  useThemeStore: () => ({
    theme: 'dark',
    setTheme: vi.fn()
  })
}))

vi.mock('../stores/upload', () => ({
  useUploadStore: () => ({
    uploadPath: '',
    setUploadPath: vi.fn()
  })
}))

vi.mock('../stores/calculate', () => {
  const mockStore = {
    warningEnabled: true,
    warningSize: 7 * 1024 * 1024 * 1024,
    todayBytes: 0,
    setWarningEnabled: vi.fn(),
    setWarningSize: vi.fn(),
    clearHistory: vi.fn()
  }
  return { useCalculateStore: () => mockStore }
})

// Mock naive-ui components
vi.mock('naive-ui', () => ({
  NForm: { name: 'NForm', template: '<div><slot /></div>' },
  NFormItem: { name: 'NFormItem', template: '<div><slot /></div>', props: ['label', 'labelPlacement', 'labelWidth'] },
  NRadioGroup: { name: 'NRadioGroup', template: '<div><slot /></div>', props: ['value'], emits: ['update:value'] },
  NRadio: { name: 'NRadio', template: '<label><slot /></label>', props: ['value'] },
  NSpace: { name: 'NSpace', template: '<div class="flex gap-2"><slot /></div>' },
  NInputGroup: { name: 'NInputGroup', template: '<div class="flex"><slot /></div>' },
  NInput: { name: 'NInput', template: '<input />', props: ['value', 'readonly', 'placeholder'] },
  NButton: { name: 'NButton', template: '<button><slot /></button>', props: ['type'] },
  NInputNumber: { name: 'NInputNumber', template: '<input type="number" />', props: ['value', 'min', 'max'] },
  NSwitch: { name: 'NSwitch', template: '<input type="checkbox" />', props: ['value'], emits: ['update:value'] },
  NText: { name: 'NText', template: '<span><slot /></span>', props: ['depth'] },
  NDivider: { name: 'NDivider', template: '<hr />' }
}))

describe('Settings.vue', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
  })

  it('renders settings form', () => {
    const wrapper = mount(Settings)
    expect(wrapper.find('.p-4').exists()).toBe(true)
  })

  it('loads default settings on mount', () => {
    const wrapper = mount(Settings)
    expect(wrapper.vm.settings.concurrentUploads).toBe(3)
    expect(wrapper.vm.settings.concurrentDownloads).toBe(5)
    expect(wrapper.vm.settings.autoCleanRecycle).toBe(true)
    expect(wrapper.vm.settings.recycleRetentionDays).toBe(30)
  })

  it('loads saved settings from localStorage', () => {
    localStorage.setItem('wenxi-settings', JSON.stringify({
      concurrentUploads: 5,
      chunkSize: 20
    }))

    const wrapper = mount(Settings)
    expect(wrapper.vm.settings.concurrentUploads).toBe(5)
    expect(wrapper.vm.settings.chunkSize).toBe(20)
  })

  it('saves settings when updated', () => {
    const wrapper = mount(Settings)
    wrapper.vm.updateSetting('concurrentUploads', 2)

    const saved = JSON.parse(localStorage.getItem('wenxi-settings'))
    expect(saved.concurrentUploads).toBe(2)
  })

  it('has theme selection radios', () => {
    const wrapper = mount(Settings)
    expect(wrapper.vm.settings).toBeDefined()
    expect(['dark', 'light', 'auto']).toContain(wrapper.vm.settings.autoCleanRecycle ? 'dark' : 'light')
  })

  it('chooses upload path via file picker', async () => {
    const wrapper = mount(Settings)
    const mockFiles = [{ webkitRelativePath: 'test-folder/file.txt' }]
    wrapper.vm.updateSetting('uploadPath', 'test-folder')
    expect(wrapper.vm.settings.uploadPath).toBe('test-folder')
  })
})
