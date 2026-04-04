import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useRecycleStore } from './recycle'

// Mock API
vi.mock('../api/recycle', () => ({
  recycleAPI: {
    list: vi.fn(),
    restoreFile: vi.fn(),
    restoreFolder: vi.fn(),
    permanentDelete: vi.fn(),
    clearAll: vi.fn(),
  }
}))

import { recycleAPI } from '../api/recycle'

describe('useRecycleStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('should initialize with default state', () => {
    const store = useRecycleStore()

    expect(store.items).toEqual([])
    expect(store.loading).toBe(false)
    expect(store.error).toBe(null)
    expect(store.selectedItems).toEqual([])
    expect(store.selectMode).toBe(false)
  })

  it('should compute counts correctly', () => {
    const store = useRecycleStore()
    store.items = [
      { id: 1, item_type: 'file', size: 1024 },
      { id: 2, item_type: 'file', size: 2048 },
      { id: 3, item_type: 'folder', size: 0 },
    ]

    expect(store.fileCount).toBe(2)
    expect(store.folderCount).toBe(1)
    expect(store.totalCount).toBe(3)
  })

  it('should compute hasSelectedItems correctly', () => {
    const store = useRecycleStore()
    expect(store.hasSelectedItems).toBe(false)

    store.selectedItems.push(1)
    expect(store.hasSelectedItems).toBe(true)
  })

  it('should toggle select correctly', () => {
    const store = useRecycleStore()

    store.toggleSelect(1)
    expect(store.selectedItems).toContain(1)

    store.toggleSelect(1)
    expect(store.selectedItems).not.toContain(1)
  })

  it('should fetch list and populate items', async () => {
    const store = useRecycleStore()
    const mockItems = [
      { id: 1, original_name: 'file1.txt', item_type: 'file', size: 1024 },
      { id: 2, original_name: 'folder1', item_type: 'folder', size: 0 },
    ]

    recycleAPI.list.mockResolvedValueOnce({ data: mockItems })
    await store.fetchList()

    expect(store.items).toEqual(mockItems)
    expect(store.loading).toBe(false)
    expect(store.error).toBe(null)
  })

  it('should set error on fetch failure', async () => {
    const store = useRecycleStore()
    recycleAPI.list.mockRejectedValueOnce(new Error('Network error'))

    await store.fetchList()

    expect(store.error).toBe('Network error')
    expect(store.loading).toBe(false)
  })

  it('should restore item and remove from list', async () => {
    const store = useRecycleStore()
    store.items = [{ id: 1, original_name: 'file.txt', item_type: 'file', size: 1024 }]
    store.selectedItems = [1]

    recycleAPI.restoreFile.mockResolvedValueOnce({})
    await store.restore(1)

    expect(store.items).toHaveLength(0)
    expect(store.selectedItems).toHaveLength(0)
  })

  it('should permanent delete item and remove from list', async () => {
    const store = useRecycleStore()
    store.items = [{ id: 1, original_name: 'file.txt', item_type: 'file', size: 1024 }]
    store.selectedItems = [1]

    recycleAPI.permanentDelete.mockResolvedValueOnce({})
    await store.permanentDelete(1)

    expect(store.items).toHaveLength(0)
    expect(store.selectedItems).toHaveLength(0)
  })

  it('should batch restore items', async () => {
    const store = useRecycleStore()
    store.items = [
      { id: 1, original_name: 'file1.txt', item_type: 'file', size: 1024 },
      { id: 2, original_name: 'file2.txt', item_type: 'file', size: 2048 },
    ]
    store.selectedItems = [1, 2]

    recycleAPI.restoreFile.mockResolvedValue({})
    await store.batchRestore()

    expect(store.selectedItems).toHaveLength(0)
    expect(recycleAPI.restoreFile).toHaveBeenCalledTimes(2)
  })

  it('should batch delete items', async () => {
    const store = useRecycleStore()
    store.items = [
      { id: 1, original_name: 'file1.txt', item_type: 'file', size: 1024 },
      { id: 2, original_name: 'file2.txt', item_type: 'file', size: 2048 },
    ]
    store.selectedItems = [1, 2]

    recycleAPI.permanentDelete.mockResolvedValue({})
    await store.batchDelete()

    expect(store.selectedItems).toHaveLength(0)
  })

  it('should clear all items', async () => {
    const store = useRecycleStore()
    store.items = [
      { id: 1, original_name: 'file.txt', item_type: 'file', size: 1024 },
    ]
    store.selectedItems = [1]

    recycleAPI.clearAll.mockResolvedValueOnce({})
    await store.clearAll()

    expect(store.items).toHaveLength(0)
    expect(store.selectedItems).toHaveLength(0)
  })

  it('should format file size correctly', () => {
    const store = useRecycleStore()

    expect(store.formatSize(0)).toBe('0 B')
    expect(store.formatSize(1024)).toBe('1.00 KB')
    expect(store.formatSize(1048576)).toBe('1.00 MB')
    expect(store.formatSize(1073741824)).toBe('1.00 GB')
  })

  it('should format date correctly', () => {
    const store = useRecycleStore()

    expect(store.formatDate('2026-04-04T10:00:00Z')).toContain('2026')
    expect(store.formatDate(null)).toBe('--')
    expect(store.formatDate('')).toBe('--')
  })

  it('should not batch restore when no items selected', async () => {
    const store = useRecycleStore()
    store.selectedItems = []

    await store.batchRestore()
    expect(recycleAPI.restoreFile).not.toHaveBeenCalled()
  })

  it('should not batch delete when no items selected', async () => {
    const store = useRecycleStore()
    store.selectedItems = []

    await store.batchDelete()
    expect(recycleAPI.permanentDelete).not.toHaveBeenCalled()
  })
})
