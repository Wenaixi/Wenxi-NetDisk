import { setActivePinia, createPinia } from 'pinia'
import { describe, it, expect, beforeEach, vi } from 'vitest'
import { useFileStore } from './file'

// Mock API
vi.mock('../api', () => ({
  fileAPI: {
    list: vi.fn(),
    upload: vi.fn(),
    delete: vi.fn(),
    rename: vi.fn()
  },
  folderAPI: {
    list: vi.fn(),
    create: vi.fn(),
    delete: vi.fn(),
    rename: vi.fn()
  }
}))

import { fileAPI, folderAPI } from '../api'

describe('useFileStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('should initialize with empty state', () => {
    const store = useFileStore()
    expect(store.files).toEqual([])
    expect(store.folders).toEqual([])
    expect(store.currentFolder).toBe(null)
    expect(store.loading).toBe(false)
    expect(store.breadcrumbs).toEqual([{ id: null, name: '首页' }])
  })

  it('should fetch files and folders', async () => {
    const store = useFileStore()
    const mockFiles = [
      { id: 1, name: 'test.txt', size: 1024 }
    ]
    const mockFolders = [
      { id: 1, name: 'Documents' }
    ]
    fileAPI.list.mockResolvedValue(mockFiles)
    folderAPI.list.mockResolvedValue(mockFolders)

    await store.fetchFiles(null)

    expect(store.files).toEqual(mockFiles)
    expect(store.folders).toEqual(mockFolders)
    expect(store.currentFolder).toBe(null)
    expect(store.loading).toBe(false)
  })

  it('should fetch files for specific folder', async () => {
    const store = useFileStore()
    fileAPI.list.mockResolvedValue([])
    folderAPI.list.mockResolvedValue([])

    await store.fetchFiles(5)

    expect(fileAPI.list).toHaveBeenCalledWith({ folder_id: 5 })
    expect(folderAPI.list).toHaveBeenCalledWith({ parent_id: 5 })
    expect(store.currentFolder).toBe(5)
  })

  it('should create folder and add to folders list', async () => {
    const store = useFileStore()
    const newFolder = { id: 1, name: 'NewFolder' }
    folderAPI.create.mockResolvedValue(newFolder)

    const result = await store.createFolder('NewFolder')

    expect(folderAPI.create).toHaveBeenCalledWith({ name: 'NewFolder', parent_id: null })
    expect(store.folders).toHaveLength(1)
    expect(store.folders[0]).toEqual(newFolder)
    expect(result).toEqual(newFolder)
  })

  it('should delete file from list', async () => {
    const store = useFileStore()
    store.files = [
      { id: 1, name: 'file1.txt' },
      { id: 2, name: 'file2.txt' }
    ]
    fileAPI.delete.mockResolvedValue({})

    await store.deleteFile(1)

    expect(fileAPI.delete).toHaveBeenCalledWith(1)
    expect(store.files).toHaveLength(1)
    expect(store.files[0].id).toBe(2)
  })

  it('should delete folder from list', async () => {
    const store = useFileStore()
    store.folders = [
      { id: 1, name: 'folder1' },
      { id: 2, name: 'folder2' }
    ]
    folderAPI.delete.mockResolvedValue({})

    await store.deleteFolder(1)

    expect(folderAPI.delete).toHaveBeenCalledWith(1)
    expect(store.folders).toHaveLength(1)
    expect(store.folders[0].id).toBe(2)
  })

  it('should rename file in list', async () => {
    const store = useFileStore()
    store.files = [{ id: 1, name: 'old.txt' }]
    const updated = { id: 1, name: 'new.txt' }
    fileAPI.rename.mockResolvedValue(updated)

    const result = await store.renameFile(1, 'new.txt')

    expect(fileAPI.rename).toHaveBeenCalledWith(1, 'new.txt')
    expect(store.files[0].name).toBe('new.txt')
    expect(result).toEqual(updated)
  })

  it('should rename folder in list', async () => {
    const store = useFileStore()
    store.folders = [{ id: 1, name: 'old' }]
    const updated = { id: 1, name: 'new' }
    folderAPI.rename.mockResolvedValue(updated)

    const result = await store.renameFolder(1, 'new')

    expect(folderAPI.rename).toHaveBeenCalledWith(1, 'new')
    expect(store.folders[0].name).toBe('new')
    expect(result).toEqual(updated)
  })

  it('should navigate to folder', async () => {
    const store = useFileStore()
    fileAPI.list.mockResolvedValue([])
    folderAPI.list.mockResolvedValue([])

    store.navigateToFolder({ id: 1, name: 'Documents' })

    expect(store.breadcrumbs).toHaveLength(2)
    expect(store.breadcrumbs[1]).toEqual({ id: 1, name: 'Documents' })
  })

  it('should navigate to root', () => {
    const store = useFileStore()
    store.breadcrumbs = [
      { id: null, name: '首页' },
      { id: 1, name: 'Documents' }
    ]

    store.navigateToRoot()

    expect(store.breadcrumbs).toHaveLength(1)
    expect(store.breadcrumbs[0]).toEqual({ id: null, name: '首页' })
  })

  it('should navigate back', () => {
    const store = useFileStore()
    store.breadcrumbs = [
      { id: null, name: '首页' },
      { id: 1, name: 'Documents' },
      { id: 2, name: 'SubFolder' }
    ]

    store.navigateBack()

    expect(store.breadcrumbs).toHaveLength(2)
    expect(store.breadcrumbs[1]).toEqual({ id: 1, name: 'Documents' })
  })

  it('should not navigate back when at root', () => {
    const store = useFileStore()
    store.breadcrumbs = [{ id: null, name: '首页' }]

    store.navigateBack()

    expect(store.breadcrumbs).toHaveLength(1)
  })

  it('should handle fetchFiles error', async () => {
    const store = useFileStore()
    fileAPI.list.mockRejectedValue(new Error('Network error'))
    folderAPI.list.mockResolvedValue([])

    await expect(store.fetchFiles()).rejects.toThrow('Network error')
    expect(store.files).toEqual([])
    expect(store.folders).toEqual([])
  })
})
