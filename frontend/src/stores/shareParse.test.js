import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useShareParseStore } from './shareParse'

// Mock API
vi.mock('../api', () => ({
  default: {
    post: vi.fn(),
  }
}))

import api from '../api'

describe('useShareParseStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('should initialize with default state', () => {
    const store = useShareParseStore()

    expect(store.loading).toBe(false)
    expect(store.error).toBe(null)
    expect(store.parsedShares).toEqual([])
    expect(store.currentShare).toBe(null)
    expect(store.downloadQueue).toEqual([])
  })

  it('should compute hasParsedShares correctly', () => {
    const store = useShareParseStore()
    expect(store.hasParsedShares).toBe(false)

    store.parsedShares.push({ name: 'file1.txt' })
    expect(store.hasParsedShares).toBe(true)
  })

  it('should compute hasDownloadItems correctly', () => {
    const store = useShareParseStore()
    expect(store.hasDownloadItems).toBe(false)

    store.downloadQueue.push({ name: 'file1.txt' })
    expect(store.hasDownloadItems).toBe(true)
  })

  it('should throw error for invalid URL', async () => {
    const store = useShareParseStore()

    await expect(store.parseShare('not-a-url'))
      .rejects.toThrow('请输入有效的分享链接')
  })

  it('should throw error for empty URL', async () => {
    const store = useShareParseStore()

    await expect(store.parseShare(''))
      .rejects.toThrow('请输入有效的分享链接')
  })

  it('should parse share and populate results', async () => {
    const store = useShareParseStore()
    const mockResponse = {
      data: {
        name: 'test-folder',
        type: 'folder',
        size: '10 MB',
        list: [
          { name: 'file1.txt', size: '1 MB', time: '2024-01-01', url: 'https://example.com/file1' },
          { name: 'file2.pdf', size: '5 MB', time: '2024-01-02', url: 'https://example.com/file2' },
        ]
      }
    }

    api.post.mockResolvedValueOnce(mockResponse)
    await store.parseShare('https://example.com/share')

    expect(store.currentShare).toEqual(mockResponse.data)
    expect(store.parsedShares).toEqual(mockResponse.data.list)
    expect(store.loading).toBe(false)
    expect(store.error).toBe(null)
  })

  it('should set error on parse failure', async () => {
    const store = useShareParseStore()
    api.post.mockRejectedValueOnce(new Error('解析失败'))

    await expect(store.parseShare('https://example.com/share'))
      .rejects.toThrow('解析失败')

    expect(store.error).toBe('解析失败')
    expect(store.loading).toBe(false)
  })

  it('should throw error when API returns no list', async () => {
    const store = useShareParseStore()
    api.post.mockResolvedValueOnce({ data: {} })

    await expect(store.parseShare('https://example.com/share'))
      .rejects.toThrow('解析失败，请检查链接是否正确')
  })

  it('should validate share URL', async () => {
    const store = useShareParseStore()
    api.post.mockResolvedValueOnce({ data: { valid: true } })

    const result = await store.validateShare('https://example.com/share')
    expect(result).toBe(true)
  })

  it('should return false on validation failure', async () => {
    const store = useShareParseStore()
    api.post.mockRejectedValueOnce(new Error('Invalid'))

    const result = await store.validateShare('https://example.com/share')
    expect(result).toBe(false)
  })

  it('should add items to download queue', () => {
    const store = useShareParseStore()
    const items = [
      { name: 'file1.txt', url: 'https://example.com/file1' },
      { name: 'file2.pdf', url: 'https://example.com/file2' },
    ]

    store.addToDownload(items)

    expect(store.downloadQueue).toHaveLength(2)
    expect(store.downloadQueue[0].status).toBe('pending')
    expect(store.downloadQueue[1].status).toBe('pending')
  })

  it('should add single item to download queue', () => {
    const store = useShareParseStore()
    const item = { name: 'file1.txt', url: 'https://example.com/file1' }

    store.addToDownload(item)

    expect(store.downloadQueue).toHaveLength(1)
  })

  it('should clear parsed results', () => {
    const store = useShareParseStore()
    store.currentShare = { name: 'test' }
    store.parsedShares = [{ name: 'file1' }]
    store.error = 'some error'

    store.clearParsed()

    expect(store.currentShare).toBe(null)
    expect(store.parsedShares).toEqual([])
    expect(store.error).toBe(null)
  })

  it('should clear download queue', () => {
    const store = useShareParseStore()
    store.downloadQueue = [
      { name: 'file1.txt', status: 'completed' },
      { name: 'file2.pdf', status: 'pending' },
    ]

    store.clearDownloadQueue()

    expect(store.downloadQueue).toHaveLength(0)
  })

  it('should get share download URL', async () => {
    const store = useShareParseStore()
    api.post.mockResolvedValueOnce({ data: { download_url: 'https://download.example.com/file' } })

    const url = await store.downloadShareFile('https://share.example.com/file')
    expect(url).toBe('https://download.example.com/file')
  })

  it('should set error on download URL failure', async () => {
    const store = useShareParseStore()
    api.post.mockRejectedValueOnce(new Error('Network error'))

    await expect(store.downloadShareFile('https://share.example.com/file'))
      .rejects.toThrow('Network error')

    expect(store.error).toBe('Network error')
  })
})
