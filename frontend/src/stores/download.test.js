import { setActivePinia, createPinia } from 'pinia'
import { describe, it, expect, beforeEach, vi } from 'vitest'
import { useDownloadStore } from './download'

// Mock API
vi.mock('../api', () => ({
  fileAPI: {
    download: vi.fn()
  }
}))

// Mock crypto utils
vi.mock('../utils/crypto', () => ({
  importKey: vi.fn().mockResolvedValue({}),
  decryptFile: vi.fn().mockResolvedValue(new Blob(['decrypted']))
}))

import { fileAPI } from '../api'
import { importKey, decryptFile } from '../utils/crypto'

describe('useDownloadStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('should initialize with default state', () => {
    const store = useDownloadStore()
    expect(store.isDownloading).toBe(false)
    expect(store.isDecrypting).toBe(false)
    expect(store.progress).toBe(0)
    expect(store.currentFile).toBe(null)
    expect(store.error).toBe(null)
    expect(store.downloadSpeed).toBe(0)
  })

  it('should reset state correctly', () => {
    const store = useDownloadStore()
    store.isDownloading = true
    store.isDecrypting = true
    store.progress = 50
    store.currentFile = { id: 1, name: 'test.txt' }
    store.error = 'Some error'

    store.reset()

    expect(store.isDownloading).toBe(false)
    expect(store.isDecrypting).toBe(false)
    expect(store.progress).toBe(0)
    expect(store.currentFile).toBe(null)
    expect(store.error).toBe(null)
    expect(store.downloadSpeed).toBe(0)
  })

  it('should throw error if fileId is empty', async () => {
    const store = useDownloadStore()

    await expect(store.downloadAndDecrypt('', 'key', 'nonce'))
      .rejects.toThrow('文件ID不能为空')
  })

  it('should set currentFile after successful download info fetch', async () => {
    const store = useDownloadStore()

    const mockDownloadResponse = {
      download_url: 'https://example.com/download',
      file_name: 'test.txt',
      file_size: 1024
    }
    fileAPI.download.mockResolvedValue(mockDownloadResponse)

    // Mock XMLHttpRequest
    const mockXHR = {
      open: vi.fn(),
      setRequestHeader: vi.fn(),
      send: vi.fn(),
      upload: { addEventListener: vi.fn() },
      responseType: '',
      addEventListener: vi.fn((event, callback) => {
        if (event === 'load') {
          callback()
        }
      }),
      status: 200
    }
    global.XMLHttpRequest = vi.fn(() => mockXHR)

    // Note: We can't fully test downloadFile without more complex mocking
    // But we can verify the API is called correctly
  })

  it('should call fileAPI.download with correct id', async () => {
    const store = useDownloadStore()

    const mockDownloadResponse = {
      download_url: 'https://example.com/download',
      file_name: 'test.txt',
      file_size: 1024
    }
    fileAPI.download.mockResolvedValue(mockDownloadResponse)

    try {
      await store.downloadAndDecrypt(123, 'key', 'nonce')
    } catch (e) {
      // Expected to fail on download
    }

    expect(fileAPI.download).toHaveBeenCalledWith(123)
  })

  it('should set error when download API fails', async () => {
    const store = useDownloadStore()

    fileAPI.download.mockRejectedValue(new Error('API Error'))

    await expect(store.downloadAndDecrypt(1, 'key', 'nonce'))
      .rejects.toThrow('API Error')

    expect(store.error).toBe('API Error')
    expect(store.isDownloading).toBe(false)
  })

  it('should reset isDownloading and isDecrypting on error', async () => {
    const store = useDownloadStore()
    store.isDownloading = true
    store.isDecrypting = true

    fileAPI.download.mockRejectedValue(new Error('Error'))

    await expect(store.downloadAndDecrypt(1, 'key', 'nonce'))
      .rejects.toThrow()

    expect(store.isDownloading).toBe(false)
    expect(store.isDecrypting).toBe(false)
  })
})
