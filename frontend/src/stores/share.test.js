import { setActivePinia, createPinia } from 'pinia'
import { describe, it, expect, beforeEach, vi } from 'vitest'
import { useShareStore } from './share'

// Mock API
vi.mock('../api', () => ({
  shareAPI: {
    list: vi.fn(),
    create: vi.fn(),
    get: vi.fn(),
    delete: vi.fn()
  }
}))

import { shareAPI } from '../api'

describe('useShareStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('should initialize with empty state', () => {
    const store = useShareStore()
    expect(store.shares).toEqual([])
    expect(store.loading).toBe(false)
    expect(store.error).toBe(null)
  })

  it('should fetch shares successfully', async () => {
    const mockShares = [
      { id: 1, file_id: 1, file_name: 'test.txt', requires_password: false },
      { id: 2, file_id: 2, file_name: 'secret.pdf', requires_password: true }
    ]
    shareAPI.list.mockResolvedValue(mockShares)

    const store = useShareStore()
    await store.fetchShares()

    expect(shareAPI.list).toHaveBeenCalledOnce()
    expect(store.shares).toEqual(mockShares)
    expect(store.loading).toBe(false)
    expect(store.error).toBe(null)
  })

  it('should set error when fetch fails', async () => {
    shareAPI.list.mockRejectedValue(new Error('Network error'))

    const store = useShareStore()
    await expect(store.fetchShares()).rejects.toThrow('Network error')

    expect(store.error).toBe('Network error')
    expect(store.loading).toBe(false)
  })

  it('should create share successfully', async () => {
    const mockShare = {
      share_token: 'abc123',
      share_url: '/api/shares/abc123'
    }
    shareAPI.create.mockResolvedValue(mockShare)

    const store = useShareStore()
    const result = await store.createShare(1, {})

    expect(shareAPI.create).toHaveBeenCalledWith({
      file_id: 1,
      password: null,
      expires_at: null
    })
    expect(result).toEqual(mockShare)
    expect(store.error).toBe(null)
  })

  it('should create share with password', async () => {
    const mockShare = {
      share_token: 'xyz789',
      share_url: '/api/shares/xyz789'
    }
    shareAPI.create.mockResolvedValue(mockShare)

    const store = useShareStore()
    await store.createShare(1, { password: 'secret123' })

    expect(shareAPI.create).toHaveBeenCalledWith({
      file_id: 1,
      password: 'secret123',
      expires_at: null
    })
  })

  it('should set error when create fails', async () => {
    shareAPI.create.mockRejectedValue(new Error('Create failed'))

    const store = useShareStore()
    await expect(store.createShare(1, {})).rejects.toThrow('Create failed')

    expect(store.error).toBe('Create failed')
  })

  it('should get share by token', async () => {
    const mockShare = {
      id: 1,
      file_id: 1,
      file_name: 'test.txt',
      requires_password: false
    }
    shareAPI.get.mockResolvedValue(mockShare)

    const store = useShareStore()
    const result = await store.getShare('abc123')

    expect(shareAPI.get).toHaveBeenCalledWith('abc123')
    expect(result).toEqual(mockShare)
  })

  it('should delete share and update list', async () => {
    shareAPI.delete.mockResolvedValue(undefined)

    const store = useShareStore()
    store.shares = [
      { id: 1, file_id: 1, file_name: 'test1.txt' },
      { id: 2, file_id: 2, file_name: 'test2.txt' }
    ]

    await store.deleteShare(1)

    expect(shareAPI.delete).toHaveBeenCalledWith(1)
    expect(store.shares).toHaveLength(1)
    expect(store.shares[0].id).toBe(2)
  })

  it('should set error when delete fails', async () => {
    shareAPI.delete.mockRejectedValue(new Error('Delete failed'))

    const store = useShareStore()
    store.shares = [{ id: 1, file_id: 1 }]

    await expect(store.deleteShare(1)).rejects.toThrow('Delete failed')

    expect(store.error).toBe('Delete failed')
  })

  it('should reset error', async () => {
    const store = useShareStore()
    store.error = 'Some error'

    store.resetError()

    expect(store.error).toBe(null)
  })
})
