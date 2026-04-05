import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useUploadTaskStore = defineStore('uploadTask', () => {
  const tasks = ref([])

  const pendingCount = computed(() => tasks.value.filter(t => t.status === 'pending').length)
  const activeCount = computed(() => tasks.value.filter(t => t.status === 'uploading').length)
  const completedCount = computed(() => tasks.value.filter(t => t.status === 'completed').length)
  const errorCount = computed(() => tasks.value.filter(t => t.status === 'error').length)

  function addTask(task) {
    tasks.value.unshift({
      id: task.id || Date.now().toString() + Math.random().toString(36).slice(2),
      name: task.name || '',
      size: task.size || 0,
      status: 'pending',
      progress: 0,
      speed: '',
      error: null,
      ...task,
    })
  }

  function removeTask(id) {
    const idx = tasks.value.findIndex(t => t.id === id)
    if (idx !== -1) {
      tasks.value.splice(idx, 1)
    }
  }

  function clearCompleted() {
    tasks.value = tasks.value.filter(t => t.status !== 'completed')
  }

  function pauseTask(id) {
    const task = tasks.value.find(t => t.id === id)
    if (task && task.status === 'uploading') {
      task.status = 'paused'
    }
  }

  function resumeTask(id) {
    const task = tasks.value.find(t => t.id === id)
    if (task && (task.status === 'paused' || task.status === 'pending')) {
      task.status = 'uploading'
    }
  }

  function pauseAll() {
    tasks.value.forEach(t => {
      if (t.status === 'uploading') t.status = 'paused'
    })
  }

  function resumeAll() {
    tasks.value.forEach(t => {
      if (t.status === 'paused' || t.status === 'pending') t.status = 'uploading'
    })
  }

  function updateProgress(id, progress, speed) {
    const task = tasks.value.find(t => t.id === id)
    if (task) {
      task.progress = Math.min(100, Math.round(progress))
      task.speed = speed || ''
    }
  }

  function completeTask(id) {
    const task = tasks.value.find(t => t.id === id)
    if (task) {
      task.status = 'completed'
      task.progress = 100
    }
  }

  function failTask(id, error) {
    const task = tasks.value.find(t => t.id === id)
    if (task) {
      task.status = 'error'
      task.error = error || '上传失败'
    }
  }

  return {
    tasks,
    pendingCount,
    activeCount,
    completedCount,
    errorCount,
    addTask,
    removeTask,
    clearCompleted,
    pauseTask,
    resumeTask,
    pauseAll,
    resumeAll,
    updateProgress,
    completeTask,
    failTask,
  }
})
