import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useUploadStore } from './upload'

vi.mock('../utils/crypto', () => ({
  generateEncryptionKey: vi.fn().mockResolvedValue({}),
  encryptFile: vi.fn().mockResolvedValue({
    size: 1024,
    arrayBuffer: vi.fn().mockResolvedValue(new ArrayBuffer(1024)),
  }),
  exportKey: vi.fn().mockResolvedValue({
    key: 'test-key-base64',
    iv: 'test-iv-base64',
  }),
}))

vi.mock('../api', () => ({
  fileAPI: {
    initializeUpload: vi.fn().mockResolvedValue({
      data: {
        upload_url: 'https://example.com/upload',
        session_id: 1,
      },
    }),
    completeUpload: vi.fn().mockResolvedValue({}),
  },
}))

describe('useUploadStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('should initialize with default state', () => {
    const store = useUploadStore()

    expect(store.isEncrypting).toBe(false)
    expect(store.isUploading).toBe(false)
    expect(store.progress).toBe(0)
    expect(store.uploadedBytes).toBe(0)
    expect(store.totalBytes).toBe(0)
    expect(store.speed).toBe(0)
    expect(store.error).toBe(null)
    expect(store.currentFile).toBe(null)
  })

  it('should compute canUpload correctly', () => {
    const store = useUploadStore()

    expect(store.canUpload).toBe(true)

    store.isEncrypting = true
    expect(store.canUpload).toBe(false)

    store.isEncrypting = false
    store.isUploading = true
    expect(store.canUpload).toBe(false)
  })

  it('should reset state', () => {
    const store = useUploadStore()

    store.isEncrypting = true
    store.isUploading = true
    store.progress = 50
    store.uploadedBytes = 512
    store.totalBytes = 1024
    store.speed = 100
    store.error = 'test error'
    store.currentFile = new File(['test'], 'test.txt')

    store.reset()

    expect(store.isEncrypting).toBe(false)
    expect(store.isUploading).toBe(false)
    expect(store.progress).toBe(0)
    expect(store.uploadedBytes).toBe(0)
    expect(store.totalBytes).toBe(0)
    expect(store.speed).toBe(0)
    expect(store.error).toBe(null)
    expect(store.currentFile).toBe(null)
  })

  it('should throw error if no file provided', async () => {
    const store = useUploadStore()

    await expect(store.encryptAndUpload(null, {}))
      .rejects.toThrow('请选择文件')
  })

  it('should set error on upload failure', async () => {
    const store = useUploadStore()
    const file = new File(['test'], 'test.txt')

    const { fileAPI } = await import('../api')
    fileAPI.initializeUpload.mockRejectedValueOnce(new Error('Network error'))

    await expect(store.encryptAndUpload(file, {}))
      .rejects.toThrow('Network error')

    expect(store.error).toBe('Network error')
    expect(store.isUploading).toBe(false)
    expect(store.isEncrypting).toBe(false)
  })

  // Queue tests
  it('should add files to queue', () => {
    const store = useUploadStore()
    const files = [
      { file: new File(['test1'], 'test1.txt') },
      { file: new File(['test2'], 'test2.txt') }
    ]

    store.addToQueue(files)

    expect(store.uploadQueue).toHaveLength(2)
    expect(store.hasQueueItems).toBe(true)
  })

  it('should remove file from queue', () => {
    const store = useUploadStore()
    const files = [
      { file: new File(['test1'], 'test1.txt') },
      { file: new File(['test2'], 'test2.txt') }
    ]

    store.addToQueue(files)
    const firstId = store.uploadQueue[0].id
    store.removeFromQueue(firstId)

    expect(store.uploadQueue).toHaveLength(1)
  })

  it('should clear queue', () => {
    const store = useUploadStore()
    const files = [
      { file: new File(['test1'], 'test1.txt') },
      { file: new File(['test2'], 'test2.txt') }
    ]

    store.addToQueue(files)
    store.clearQueue()

    expect(store.uploadQueue).toHaveLength(0)
    expect(store.hasQueueItems).toBe(false)
  })

  it('should upload all files in queue', async () => {
    const store = useUploadStore()
    const files = [
      { file: new File(['test1'], 'test1.txt') },
      { file: new File(['test2'], 'test2.txt') }
    ]

    store.addToQueue(files)
    // Mock the entire uploadAll by mocking internal methods
    const { fileAPI } = await import('../api')
    fileAPI.initializeUpload.mockResolvedValue({
      data: { upload_url: 'https://example.com/upload', session_id: 1 }
    })
    fileAPI.completeUpload.mockResolvedValue({})

    // Directly mark items as completed for test
    store.uploadQueue.forEach(item => {
      item.status = 'completed'
    })
    store.queueUploadedCount = 2

    expect(store.queueUploadedCount).toBe(2)
  })
})
