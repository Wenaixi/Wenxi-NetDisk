import { setActivePinia, createPinia } from 'pinia'
import { describe, it, expect, beforeEach, vi } from 'vitest'
import { useFileStore } from './file'

// Mock API
vi.mock('../api', () => ({
  fileAPI: {
    list: vi.fn(),
    upload: vi.fn(),
    delete: vi.fn(),
    rename: vi.fn(),
    move: vi.fn(),
    updateDescription: vi.fn()
  },
  folderAPI: {
    list: vi.fn(),
    create: vi.fn(),
    delete: vi.fn(),
    rename: vi.fn(),
    move: vi.fn(),
    updateDescription: vi.fn()
  }
}))

import { fileAPI, folderAPI } from '../api'

describe('useFileStore - updateDescription', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('should update file description and update in list', async () => {
    const store = useFileStore()
    store.files = [{ id: 1, name: 'file.txt', description: '' }]
    const updated = { id: 1, name: 'file.txt', description: 'new desc' }
    fileAPI.updateDescription.mockResolvedValue(updated)

    const result = await store.updateFileDescription(1, 'new desc')

    expect(fileAPI.updateDescription).toHaveBeenCalledWith(1, 'new desc')
    expect(store.files[0].description).toBe('new desc')
    expect(result).toEqual(updated)
  })

  it('should update folder description and update in list', async () => {
    const store = useFileStore()
    store.folders = [{ id: 1, name: 'folder', description: '' }]
    const updated = { id: 1, name: 'folder', description: 'folder desc' }
    folderAPI.updateDescription.mockResolvedValue(updated)

    const result = await store.updateFolderDescription(1, 'folder desc')

    expect(folderAPI.updateDescription).toHaveBeenCalledWith(1, 'folder desc')
    expect(store.folders[0].description).toBe('folder desc')
    expect(result).toEqual(updated)
  })

  it('should propagate error from updateFileDescription', async () => {
    const store = useFileStore()
    fileAPI.updateDescription.mockRejectedValue(new Error('Update failed'))

    await expect(store.updateFileDescription(1, 'desc'))
      .rejects.toThrow('Update failed')
  })

  it('should propagate error from updateFolderDescription', async () => {
    const store = useFileStore()
    folderAPI.updateDescription.mockRejectedValue(new Error('Update failed'))

    await expect(store.updateFolderDescription(1, 'desc'))
      .rejects.toThrow('Update failed')
  })
})
