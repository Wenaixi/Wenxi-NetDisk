import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useDownloadTaskStore } from './downloadTask'

describe('downloadTask store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('should initialize with empty tasks', () => {
    const store = useDownloadTaskStore()
    expect(store.tasks).toHaveLength(0)
    expect(store.pendingCount).toBe(0)
    expect(store.activeCount).toBe(0)
    expect(store.completedCount).toBe(0)
    expect(store.errorCount).toBe(0)
  })

  it('should add task', () => {
    const store = useDownloadTaskStore()
    store.addTask({ name: 'test.zip', size: 2048 })
    expect(store.tasks).toHaveLength(1)
    expect(store.tasks[0].name).toBe('test.zip')
    expect(store.tasks[0].status).toBe('pending')
  })

  it('should remove task', () => {
    const store = useDownloadTaskStore()
    store.addTask({ name: 'test.zip' })
    const id = store.tasks[0].id
    store.removeTask(id)
    expect(store.tasks).toHaveLength(0)
  })

  it('should clear completed tasks', () => {
    const store = useDownloadTaskStore()
    store.addTask({ name: 'file1.zip' })
    store.addTask({ name: 'file2.zip' })
    // unshift adds to beginning, so order is [file2, file1]
    store.tasks[0].status = 'completed' // file2 is completed
    store.clearCompleted()
    expect(store.tasks).toHaveLength(1)
    expect(store.tasks[0].name).toBe('file1.zip')
  })

  it('should pause task', () => {
    const store = useDownloadTaskStore()
    store.addTask({ name: 'test.zip' })
    const id = store.tasks[0].id
    store.tasks[0].status = 'downloading'
    store.pauseTask(id)
    expect(store.tasks[0].status).toBe('paused')
  })

  it('should resume task', () => {
    const store = useDownloadTaskStore()
    store.addTask({ name: 'test.zip' })
    const id = store.tasks[0].id
    store.tasks[0].status = 'paused'
    store.resumeTask(id)
    expect(store.tasks[0].status).toBe('downloading')
  })

  it('should pause all tasks', () => {
    const store = useDownloadTaskStore()
    store.addTask({ name: 'file1.zip' })
    store.addTask({ name: 'file2.zip' })
    store.tasks[0].status = 'downloading'
    store.tasks[1].status = 'downloading'
    store.pauseAll()
    expect(store.tasks.filter(t => t.status === 'downloading')).toHaveLength(0)
  })

  it('should resume all tasks', () => {
    const store = useDownloadTaskStore()
    store.addTask({ name: 'file1.zip' })
    store.addTask({ name: 'file2.zip' })
    store.tasks[0].status = 'paused'
    store.tasks[1].status = 'pending'
    store.resumeAll()
    expect(store.tasks.every(t => t.status === 'downloading')).toBe(true)
  })

  it('should update progress', () => {
    const store = useDownloadTaskStore()
    store.addTask({ name: 'test.zip' })
    const id = store.tasks[0].id
    store.updateProgress(id, 60, '3MB/s')
    expect(store.tasks[0].progress).toBe(60)
    expect(store.tasks[0].speed).toBe('3MB/s')
  })

  it('should cap progress at 100', () => {
    const store = useDownloadTaskStore()
    store.addTask({ name: 'test.zip' })
    const id = store.tasks[0].id
    store.updateProgress(id, 200)
    expect(store.tasks[0].progress).toBe(100)
  })

  it('should complete task', () => {
    const store = useDownloadTaskStore()
    store.addTask({ name: 'test.zip' })
    const id = store.tasks[0].id
    store.completeTask(id)
    expect(store.tasks[0].status).toBe('completed')
    expect(store.tasks[0].progress).toBe(100)
  })

  it('should fail task with error message', () => {
    const store = useDownloadTaskStore()
    store.addTask({ name: 'test.zip' })
    const id = store.tasks[0].id
    store.failTask(id, 'Connection refused')
    expect(store.tasks[0].status).toBe('error')
    expect(store.tasks[0].error).toBe('Connection refused')
  })

  it('should calculate counts correctly', () => {
    const store = useDownloadTaskStore()
    store.addTask({ name: 'file1.zip' })
    store.addTask({ name: 'file2.zip' })
    store.addTask({ name: 'file3.zip' })
    store.tasks[0].status = 'downloading'
    store.tasks[1].status = 'completed'
    store.tasks[2].status = 'error'

    expect(store.activeCount).toBe(1)
    expect(store.completedCount).toBe(1)
    expect(store.errorCount).toBe(1)
  })
})
