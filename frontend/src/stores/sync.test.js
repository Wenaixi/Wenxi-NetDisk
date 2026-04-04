import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useSyncStore } from './sync'

describe('useSyncStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.useFakeTimers()
  })

  it('should initialize with default state', () => {
    const store = useSyncStore()

    expect(store.syncTasks).toEqual([])
    expect(store.isRunning).toBe(false)
    expect(store.pendingCount).toBe(0)
    expect(store.completedCount).toBe(0)
  })

  it('should add a sync task', () => {
    const store = useSyncStore()

    store.addTask({
      url: 'https://example.com/file.zip',
      name: 'test-file.zip',
      type: 'direct',
      folderId: null
    })

    expect(store.syncTasks).toHaveLength(1)
    expect(store.syncTasks[0].name).toBe('test-file.zip')
    expect(store.syncTasks[0].status).toBe('pending')
    expect(store.syncTasks[0].step).toBe('waiting')
  })

  it('should extract filename from URL', () => {
    const store = useSyncStore()

    store.addTask({
      url: 'https://example.com/path/to/document.pdf',
      type: 'direct'
    })

    expect(store.syncTasks[0].name).toBe('document.pdf')
  })

  it('should handle URL without clear filename', () => {
    const store = useSyncStore()

    store.addTask({
      url: 'https://example.com/abc123',
      type: 'direct'
    })

    expect(store.syncTasks[0].name).toBe('abc123')
  })

  it('should compute counts correctly', () => {
    const store = useSyncStore()

    store.addTask({ url: 'https://a.com/file1', type: 'direct' })
    store.addTask({ url: 'https://b.com/file2', type: 'direct' })
    store.addTask({ url: 'https://c.com/file3', type: 'direct' })

    // Mark tasks
    store.syncTasks[0].status = 'completed'
    store.syncTasks[1].status = 'error'
    store.syncTasks[2].status = 'pending'

    expect(store.pendingCount).toBe(1)
    expect(store.completedCount).toBe(1)
    expect(store.errorCount).toBe(1)
  })

  it('should remove a task by id', () => {
    const store = useSyncStore()

    store.addTask({ url: 'https://a.com/file1', name: 'file1', type: 'direct' })
    store.addTask({ url: 'https://b.com/file2', name: 'file2', type: 'direct' })

    // addTask uses unshift, so file2 is at index 0, file1 at index 1
    const taskId = store.syncTasks[1].id // remove file1
    store.removeTask(taskId)

    expect(store.syncTasks).toHaveLength(1)
    expect(store.syncTasks[0].name).toBe('file2')
  })

  it('should clear completed tasks', () => {
    const store = useSyncStore()

    store.addTask({ url: 'https://a.com/file1', name: 'file1', type: 'direct' })
    store.addTask({ url: 'https://b.com/file2', name: 'file2', type: 'direct' })
    store.addTask({ url: 'https://c.com/file3', name: 'file3', type: 'direct' })

    store.syncTasks[0].status = 'completed'
    store.syncTasks[1].status = 'pending'
    store.syncTasks[2].status = 'completed'

    store.clearCompleted()

    expect(store.syncTasks).toHaveLength(1)
    expect(store.syncTasks[0].name).toBe('file2')
  })

  it('should clear all tasks when not running', () => {
    const store = useSyncStore()

    store.addTask({ url: 'https://a.com/file1', name: 'file1', type: 'direct' })
    store.addTask({ url: 'https://b.com/file2', name: 'file2', type: 'direct' })

    store.clearAll()

    expect(store.syncTasks).toHaveLength(0)
  })

  it('should not clear tasks when running', () => {
    const store = useSyncStore()

    store.addTask({ url: 'https://a.com/file1', name: 'file1', type: 'direct' })
    store.isRunning = true

    store.clearAll()

    expect(store.syncTasks).toHaveLength(1)
    store.isRunning = false
  })

  it('should retry a failed task', () => {
    const store = useSyncStore()

    store.addTask({ url: 'https://a.com/file1', name: 'file1', type: 'direct' })
    store.syncTasks[0].status = 'error'
    store.syncTasks[0].error = 'download failed'
    store.syncTasks[0].progress = 30

    store.retryTask(store.syncTasks[0].id)

    expect(store.syncTasks[0].status).toBe('pending')
    expect(store.syncTasks[0].error).toBeNull()
    expect(store.syncTasks[0].progress).toBe(0)
    expect(store.syncTasks[0].step).toBe('waiting')
  })

  it('should not retry a non-error task', () => {
    const store = useSyncStore()

    store.addTask({ url: 'https://a.com/file1', name: 'file1', type: 'direct' })
    store.syncTasks[0].status = 'pending'
    store.syncTasks[0].error = null

    store.retryTask(store.syncTasks[0].id)

    expect(store.syncTasks[0].status).toBe('pending')
  })
})
