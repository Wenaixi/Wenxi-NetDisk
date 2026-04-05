import { describe, it, expect, beforeEach, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useTaskHistoryStore } from './taskHistory'

describe('taskHistory store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
  })

  it('initializes with empty tasks', () => {
    const store = useTaskHistoryStore()
    expect(store.tasks).toEqual([])
  })

  it('adds a task to history', () => {
    const store = useTaskHistoryStore()
    store.addTask({ type: 'upload', fileName: 'test.pdf', fileSize: 1024 })

    expect(store.tasks.length).toBe(1)
    expect(store.tasks[0].fileName).toBe('test.pdf')
    expect(store.tasks[0].type).toBe('upload')
    expect(store.tasks[0].status).toBe('completed')
  })

  it('adds a failed task', () => {
    const store = useTaskHistoryStore()
    store.addTask({ type: 'download', fileName: 'big.zip', status: 'failed', errorMessage: 'timeout' })

    expect(store.tasks[0].status).toBe('failed')
    expect(store.tasks[0].errorMessage).toBe('timeout')
  })

  it('removes a task by id', () => {
    const store = useTaskHistoryStore()
    store.addTask({ type: 'upload', fileName: 'a.txt' })
    store.addTask({ type: 'upload', fileName: 'b.txt' })
    // tasks are prepended, so [0] is 'b.txt', [1] is 'a.txt'
    const idToRemove = store.tasks[1].id // remove 'a.txt'

    store.removeTask(idToRemove)
    expect(store.tasks.length).toBe(1)
    expect(store.tasks[0].fileName).toBe('b.txt')
  })

  it('clears all tasks', () => {
    const store = useTaskHistoryStore()
    store.addTask({ type: 'upload', fileName: 'a.txt' })
    store.addTask({ type: 'download', fileName: 'b.zip' })

    store.clearAll()
    expect(store.tasks).toEqual([])
  })

  it('clears only completed tasks', () => {
    const store = useTaskHistoryStore()
    store.addTask({ type: 'upload', fileName: 'a.txt' })
    store.addTask({ type: 'download', fileName: 'b.zip', status: 'failed' })

    store.clearCompleted()
    expect(store.tasks.length).toBe(1)
    expect(store.tasks[0].status).toBe('failed')
  })

  it('clears only failed tasks', () => {
    const store = useTaskHistoryStore()
    store.addTask({ type: 'upload', fileName: 'a.txt' })
    store.addTask({ type: 'download', fileName: 'b.zip', status: 'failed' })

    store.clearFailed()
    expect(store.tasks.length).toBe(1)
    expect(store.tasks[0].status).toBe('completed')
  })

  it('persists tasks to localStorage', () => {
    const store = useTaskHistoryStore()
    store.addTask({ type: 'sync', fileName: 'sync.pdf' })

    const saved = JSON.parse(localStorage.getItem('wenxi-task-history'))
    expect(saved.length).toBe(1)
    expect(saved[0].fileName).toBe('sync.pdf')
  })

  it('loads tasks from localStorage on init', () => {
    localStorage.setItem('wenxi-task-history', JSON.stringify([
      { id: 1, type: 'upload', fileName: 'old.pdf', status: 'completed' }
    ]))

    setActivePinia(createPinia())
    const store = useTaskHistoryStore()
    expect(store.tasks.length).toBe(1)
    expect(store.tasks[0].fileName).toBe('old.pdf')
  })

  it('handles corrupted localStorage data gracefully', () => {
    localStorage.setItem('wenxi-task-history', 'invalid json')

    setActivePinia(createPinia())
    const store = useTaskHistoryStore()
    expect(store.tasks).toEqual([])
  })
})
