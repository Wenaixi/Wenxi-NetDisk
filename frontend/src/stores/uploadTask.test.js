import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useUploadTaskStore } from './uploadTask'

describe('uploadTask store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('should initialize with empty tasks', () => {
    const store = useUploadTaskStore()
    expect(store.tasks).toHaveLength(0)
    expect(store.pendingCount).toBe(0)
    expect(store.activeCount).toBe(0)
    expect(store.completedCount).toBe(0)
    expect(store.errorCount).toBe(0)
  })

  it('should add task', () => {
    const store = useUploadTaskStore()
    store.addTask({ name: 'test.zip', size: 1024 })
    expect(store.tasks).toHaveLength(1)
    expect(store.tasks[0].name).toBe('test.zip')
    expect(store.tasks[0].status).toBe('pending')
    expect(store.tasks[0].progress).toBe(0)
  })

  it('should add task with custom id', () => {
    const store = useUploadTaskStore()
    store.addTask({ id: 'custom-id', name: 'test.zip' })
    expect(store.tasks[0].id).toBe('custom-id')
  })

  it('should remove task', () => {
    const store = useUploadTaskStore()
    store.addTask({ name: 'test.zip' })
    const id = store.tasks[0].id
    store.removeTask(id)
    expect(store.tasks).toHaveLength(0)
  })

  it('should clear completed tasks', () => {
    const store = useUploadTaskStore()
    store.addTask({ name: 'file1.zip' })
    store.addTask({ name: 'file2.zip' })
    // unshift adds to beginning, so order is [file2, file1]
    store.tasks[0].status = 'completed' // file2
    store.clearCompleted()
    expect(store.tasks).toHaveLength(1)
    expect(store.tasks[0].name).toBe('file1.zip')
  })

  it('should pause task', () => {
    const store = useUploadTaskStore()
    store.addTask({ name: 'test.zip' })
    const id = store.tasks[0].id
    store.tasks[0].status = 'uploading'
    store.pauseTask(id)
    expect(store.tasks[0].status).toBe('paused')
  })

  it('should resume task', () => {
    const store = useUploadTaskStore()
    store.addTask({ name: 'test.zip' })
    const id = store.tasks[0].id
    store.tasks[0].status = 'paused'
    store.resumeTask(id)
    expect(store.tasks[0].status).toBe('uploading')
  })

  it('should pause all tasks', () => {
    const store = useUploadTaskStore()
    store.addTask({ name: 'file1.zip' })
    store.addTask({ name: 'file2.zip' })
    store.tasks[0].status = 'uploading'
    store.tasks[1].status = 'uploading'
    store.pauseAll()
    expect(store.tasks.every(t => t.status !== 'uploading')).toBe(true)
  })

  it('should resume all tasks', () => {
    const store = useUploadTaskStore()
    store.addTask({ name: 'file1.zip' })
    store.addTask({ name: 'file2.zip' })
    store.tasks[0].status = 'paused'
    store.tasks[1].status = 'pending'
    store.resumeAll()
    expect(store.tasks.every(t => t.status === 'uploading')).toBe(true)
  })

  it('should update progress', () => {
    const store = useUploadTaskStore()
    store.addTask({ name: 'test.zip' })
    const id = store.tasks[0].id
    store.updateProgress(id, 75, '2MB/s')
    expect(store.tasks[0].progress).toBe(75)
    expect(store.tasks[0].speed).toBe('2MB/s')
  })

  it('should cap progress at 100', () => {
    const store = useUploadTaskStore()
    store.addTask({ name: 'test.zip' })
    const id = store.tasks[0].id
    store.updateProgress(id, 150)
    expect(store.tasks[0].progress).toBe(100)
  })

  it('should complete task', () => {
    const store = useUploadTaskStore()
    store.addTask({ name: 'test.zip' })
    const id = store.tasks[0].id
    store.completeTask(id)
    expect(store.tasks[0].status).toBe('completed')
    expect(store.tasks[0].progress).toBe(100)
  })

  it('should fail task', () => {
    const store = useUploadTaskStore()
    store.addTask({ name: 'test.zip' })
    const id = store.tasks[0].id
    store.failTask(id, 'Network error')
    expect(store.tasks[0].status).toBe('error')
    expect(store.tasks[0].error).toBe('Network error')
  })

  it('should calculate counts correctly', () => {
    const store = useUploadTaskStore()
    store.addTask({ name: 'file1.zip' })
    store.addTask({ name: 'file2.zip' })
    store.addTask({ name: 'file3.zip' })
    store.addTask({ name: 'file4.zip' })
    store.tasks[0].status = 'uploading'
    store.tasks[1].status = 'completed'
    store.tasks[2].status = 'error'

    expect(store.pendingCount).toBe(1)
    expect(store.activeCount).toBe(1)
    expect(store.completedCount).toBe(1)
    expect(store.errorCount).toBe(1)
  })
})
