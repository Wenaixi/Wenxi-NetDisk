import { defineStore } from 'pinia'
import { ref } from 'vue'

const TASKS_KEY = 'wenxi-task-history'

export const useTaskHistoryStore = defineStore('taskHistory', () => {
  const tasks = ref([])

  // 加载历史任务
  function loadTasks() {
    try {
      const saved = localStorage.getItem(TASKS_KEY)
      if (saved) {
        tasks.value = JSON.parse(saved)
      }
    } catch {
      tasks.value = []
    }
  }

  // 保存任务记录
  function addTask(task) {
    const entry = {
      id: Date.now() + Math.random(),
      type: task.type, // 'upload' | 'download' | 'sync'
      fileName: task.fileName,
      fileSize: task.fileSize || 0,
      status: task.status || 'completed', // 'completed' | 'failed'
      errorMessage: task.errorMessage || null,
      createdAt: new Date().toISOString()
    }
    tasks.value.unshift(entry)
    saveTasks()
  }

  // 删除单个任务
  function removeTask(id) {
    tasks.value = tasks.value.filter(t => t.id !== id)
    saveTasks()
  }

  // 清空所有历史
  function clearAll() {
    tasks.value = []
    saveTasks()
  }

  // 仅清空已完成的
  function clearCompleted() {
    tasks.value = tasks.value.filter(t => t.status !== 'completed')
    saveTasks()
  }

  // 仅清空失败的
  function clearFailed() {
    tasks.value = tasks.value.filter(t => t.status !== 'failed')
    saveTasks()
  }

  function saveTasks() {
    localStorage.setItem(TASKS_KEY, JSON.stringify(tasks.value))
  }

  // 初始化
  loadTasks()

  return {
    tasks,
    addTask,
    removeTask,
    clearAll,
    clearCompleted,
    clearFailed
  }
})
