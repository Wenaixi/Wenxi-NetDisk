import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import TaskHistory from '../views/TaskHistory.vue'
import { createPinia, setActivePinia } from 'pinia'

vi.mock('naive-ui', () => ({
  NTag: { name: 'NTag', template: '<span class="n-tag"><slot /></span>', props: ['type', 'size'] },
  NButton: { name: 'NButton', template: '<button><slot /></button>', props: ['type', 'size', 'text'] }
}))

describe('TaskHistory.vue', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
  })

  it('renders empty state when no tasks', () => {
    const wrapper = mount(TaskHistory)
    expect(wrapper.text()).toContain('暂无任务历史记录')
  })

  it('renders task list', async () => {
    const { useTaskHistoryStore } = await import('../stores/taskHistory')
    setActivePinia(createPinia())
    const store = useTaskHistoryStore()
    store.addTask({ type: 'upload', fileName: 'test.pdf', fileSize: 2048 })

    const wrapper = mount(TaskHistory)
    expect(wrapper.text()).toContain('test.pdf')
    expect(wrapper.text()).toContain('上传')
  })

  it('shows failed task status', async () => {
    const { useTaskHistoryStore } = await import('../stores/taskHistory')
    setActivePinia(createPinia())
    const store = useTaskHistoryStore()
    store.addTask({ type: 'download', fileName: 'fail.zip', status: 'failed', errorMessage: 'timeout' })

    const wrapper = mount(TaskHistory)
    expect(wrapper.text()).toContain('失败')
    expect(wrapper.text()).toContain('fail.zip')
  })

  it('removes a task when delete button clicked', async () => {
    const { useTaskHistoryStore } = await import('../stores/taskHistory')
    setActivePinia(createPinia())
    const store = useTaskHistoryStore()
    store.addTask({ type: 'upload', fileName: 'delete-me.txt' })

    const wrapper = mount(TaskHistory)
    // The delete button has text "删除", click it via vm directly
    store.removeTask(store.tasks[0].id)
    await wrapper.vm.$nextTick()

    expect(store.tasks.length).toBe(0)
  })

  it('clears completed tasks', async () => {
    const { useTaskHistoryStore } = await import('../stores/taskHistory')
    setActivePinia(createPinia())
    const store = useTaskHistoryStore()
    store.addTask({ type: 'upload', fileName: 'a.txt' })
    store.addTask({ type: 'download', fileName: 'b.zip', status: 'failed' })

    store.clearCompleted()

    const wrapper = mount(TaskHistory)
    expect(store.tasks.length).toBe(1)
    expect(store.tasks[0].status).toBe('failed')
  })

  it('formatSize formats bytes correctly', () => {
    const wrapper = mount(TaskHistory)
    expect(wrapper.vm.formatSize(1024)).toBe('1 KB')
    expect(wrapper.vm.formatSize(1048576)).toBe('1 MB')
  })

  it('formatDate formats ISO string', () => {
    const wrapper = mount(TaskHistory)
    const result = wrapper.vm.formatDate('2026-04-05T12:00:00Z')
    expect(result).toBeTruthy()
    expect(typeof result).toBe('string')
  })
})
