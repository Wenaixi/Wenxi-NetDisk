import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import FileDetailModal from './FileDetailModal.vue'

// Mock @vicons/ionicons5
vi.mock('@vicons/ionicons5', () => ({
  Document: { template: '<span data-icon="document"></span>' },
  Folder: { template: '<span data-icon="folder"></span>' }
}))

describe('FileDetailModal', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('File display', () => {
    it('should show file name from item prop', () => {
      const wrapper = mount(FileDetailModal, {
        props: {
          item: {
            id: 1,
            name: 'test-file.txt',
            size: 1024,
            mime_type: 'text/plain',
            lanzou_file_id: '12345',
            created_at: '2026-04-05T10:00:00Z',
            updated_at: '2026-04-05T11:00:00Z'
          }
        }
      })

      expect(wrapper.text()).toContain('test-file.txt')
    })

    it('should identify file by size property', () => {
      const wrapper = mount(FileDetailModal, {
        props: {
          item: {
            id: 1,
            name: 'test.txt',
            size: 100,
            created_at: '2026-04-05T10:00:00Z',
            updated_at: '2026-04-05T10:00:00Z'
          }
        }
      })

      // File items have size, so isFile should be true
      expect(wrapper.text()).toContain('文件')
    })

    it('should identify folder by missing size property', () => {
      const wrapper = mount(FileDetailModal, {
        props: {
          item: {
            id: 1,
            name: 'my-folder',
            created_at: '2026-04-05T10:00:00Z',
            updated_at: '2026-04-05T10:00:00Z'
          }
        }
      })

      // Folder items don't have size, so isFile should be false
      expect(wrapper.text()).toContain('文件夹')
    })

    it('should display file size as 1 KB', () => {
      const wrapper = mount(FileDetailModal, {
        props: {
          item: {
            id: 1,
            name: 'test.txt',
            size: 1024,
            created_at: '2026-04-05T10:00:00Z',
            updated_at: '2026-04-05T10:00:00Z'
          }
        }
      })

      expect(wrapper.text()).toContain('1 KB')
    })

    it('should display mime type when available', () => {
      const wrapper = mount(FileDetailModal, {
        props: {
          item: {
            id: 1,
            name: 'test.txt',
            size: 100,
            mime_type: 'text/plain',
            created_at: '2026-04-05T10:00:00Z',
            updated_at: '2026-04-05T10:00:00Z'
          }
        }
      })

      expect(wrapper.text()).toContain('text/plain')
    })

    it('should show unknown for missing mime type', () => {
      const wrapper = mount(FileDetailModal, {
        props: {
          item: {
            id: 1,
            name: 'test',
            size: 100,
            created_at: '2026-04-05T10:00:00Z',
            updated_at: '2026-04-05T10:00:00Z'
          }
        }
      })

      expect(wrapper.text()).toContain('未知')
    })

    it('should display file ID', () => {
      const wrapper = mount(FileDetailModal, {
        props: {
          item: {
            id: 42,
            name: 'test.txt',
            size: 100,
            created_at: '2026-04-05T10:00:00Z',
            updated_at: '2026-04-05T10:00:00Z'
          }
        }
      })

      expect(wrapper.text()).toContain('42')
    })

    it('should show lanzou file info when available', () => {
      const wrapper = mount(FileDetailModal, {
        props: {
          item: {
            id: 1,
            name: 'test.txt',
            size: 100,
            lanzou_file_id: 'abc123',
            created_at: '2026-04-05T10:00:00Z',
            updated_at: '2026-04-05T10:00:00Z'
          }
        }
      })

      expect(wrapper.text()).toContain('蓝奏云信息')
      expect(wrapper.text()).toContain('abc123')
    })

    it('should not show lanzou info for folders', () => {
      const wrapper = mount(FileDetailModal, {
        props: {
          item: {
            id: 1,
            name: 'folder',
            created_at: '2026-04-05T10:00:00Z',
            updated_at: '2026-04-05T10:00:00Z'
          }
        }
      })

      expect(wrapper.text()).not.toContain('蓝奏云文件ID')
    })
  })

  describe('Description', () => {
    it('should initialize description from item prop', async () => {
      const wrapper = mount(FileDetailModal, {
        props: {
          item: {
            id: 1,
            name: 'test.txt',
            size: 100,
            description: 'existing description',
            created_at: '2026-04-05T10:00:00Z',
            updated_at: '2026-04-05T10:00:00Z'
          }
        }
      })

      // Wait for watch to trigger
      await wrapper.vm.$nextTick()
      expect(wrapper.vm.description).toBe('existing description')
    })

    it('should emit update with id and description on save', async () => {
      const wrapper = mount(FileDetailModal, {
        props: {
          item: {
            id: 1,
            name: 'test.txt',
            size: 100,
            description: 'existing',
            created_at: '2026-04-05T10:00:00Z',
            updated_at: '2026-04-05T10:00:00Z'
          }
        }
      })

      await wrapper.vm.$nextTick()

      // Simulate save via component method directly since button click requires DOM rendering
      wrapper.vm.description = 'new description'
      wrapper.vm.saveDescription()
      await wrapper.vm.$nextTick()

      expect(wrapper.emitted('update')).toBeTruthy()
      expect(wrapper.emitted('update')[0]).toEqual([{ id: 1, description: 'new description' }])
    })
  })

  describe('Size formatting', () => {
    it('should format bytes correctly', () => {
      const wrapper = mount(FileDetailModal, {
        props: {
          item: {
            id: 1,
            name: 'test.txt',
            size: 500,
            created_at: '2026-04-05T10:00:00Z',
            updated_at: '2026-04-05T10:00:00Z'
          }
        }
      })

      expect(wrapper.text()).toContain('B')
    })

    it('should format MB correctly', () => {
      const wrapper = mount(FileDetailModal, {
        props: {
          item: {
            id: 1,
            name: 'test.txt',
            size: 5 * 1024 * 1024, // 5MB
            created_at: '2026-04-05T10:00:00Z',
            updated_at: '2026-04-05T10:00:00Z'
          }
        }
      })

      expect(wrapper.text()).toContain('5 MB')
    })

    it('should format GB correctly', () => {
      const wrapper = mount(FileDetailModal, {
        props: {
          item: {
            id: 1,
            name: 'test.txt',
            size: 2 * 1024 * 1024 * 1024, // 2GB
            created_at: '2026-04-05T10:00:00Z',
            updated_at: '2026-04-05T10:00:00Z'
          }
        }
      })

      expect(wrapper.text()).toContain('2 GB')
    })

    it('should handle zero size', () => {
      const wrapper = mount(FileDetailModal, {
        props: {
          item: {
            id: 1,
            name: 'test.txt',
            size: 0,
            created_at: '2026-04-05T10:00:00Z',
            updated_at: '2026-04-05T10:00:00Z'
          }
        }
      })

      expect(wrapper.text()).toContain('0 B')
    })
  })

  describe('Date formatting', () => {
    it('should format dates correctly', () => {
      const wrapper = mount(FileDetailModal, {
        props: {
          item: {
            id: 1,
            name: 'test.txt',
            size: 100,
            created_at: '2026-04-05T10:00:00Z',
            updated_at: '2026-04-05T11:30:00Z'
          }
        }
      })

      expect(wrapper.text()).toContain('2026')
      expect(wrapper.text()).toContain('4')
    })

    it('should show unknown for missing dates', () => {
      const wrapper = mount(FileDetailModal, {
        props: {
          item: {
            id: 1,
            name: 'test.txt',
            size: 100,
            created_at: null,
            updated_at: null
          }
        }
      })

      // formatDate returns '未知' for null dates
      expect(wrapper.text()).toContain('未知')
    })
  })
})
