import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useSyncStore = defineStore('sync', () => {
  const syncTasks = ref([])
  const isRunning = ref(false)

  // 状态统计
  const pendingCount = computed(() => syncTasks.value.filter(t => t.status === 'pending').length)
  const downloadingCount = computed(() => syncTasks.value.filter(t => t.status === 'downloading').length)
  const uploadingCount = computed(() => syncTasks.value.filter(t => t.status === 'uploading').length)
  const completedCount = computed(() => syncTasks.value.filter(t => t.status === 'completed').length)
  const errorCount = computed(() => syncTasks.value.filter(t => t.status === 'error').length)

  // 添加同步任务
  function addTask({ url, name, type = 'direct', pwd = '', folderId = null, trashOnFinish = false }) {
    const task = {
      id: Date.now() + Math.random(),
      name: name || extractFileName(url),
      url,
      type,        // 'direct' (直链) 或 'lanzou' (蓝奏云分享)
      pwd,
      folderId,
      trashOnFinish,
      status: 'pending',
      progress: 0,
      step: 'waiting', // 'waiting' | 'downloading' | 'uploading' | 'completed' | 'error'
      error: null,
      size: 0,
      downloadedSize: 0,
      speed: 0,
      createdAt: new Date().toISOString()
    }
    syncTasks.value.unshift(task)
    return task
  }

  // 从 URL 提取文件名
  function extractFileName(url) {
    try {
      const urlObj = new URL(url)
      return decodeURIComponent(urlObj.pathname.split('/').pop()) || '未知文件'
    } catch {
      return url.slice(-20)
    }
  }

  // 开始同步（执行所有 pending 任务）
  async function startSync(onDownloadProgress, onUploadProgress) {
    isRunning.value = true
    for (const task of syncTasks.value) {
      if (task.status !== 'pending') continue
      await executeTask(task, onDownloadProgress, onUploadProgress)
    }
    isRunning.value = false
  }

  // 执行单个任务
  async function executeTask(task, onDownloadProgress, onUploadProgress) {
    task.status = 'downloading'
    task.step = 'downloading'
    task.error = null

    try {
      // 阶段1: 下载文件到本地
      await downloadFile(task, onDownloadProgress)

      // 阶段2: 上传到蓝奏云
      task.status = 'uploading'
      task.step = 'uploading'
      task.progress = 0

      await uploadToLanzou(task, onUploadProgress)

      task.status = 'completed'
      task.step = 'completed'
      task.progress = 100
    } catch (err) {
      task.status = 'error'
      task.step = 'error'
      task.error = err.message || '同步失败'
    }
  }

  // 下载文件 (直链或蓝奏云分享)
  async function downloadFile(task, onProgress) {
    if (task.type === 'direct') {
      // 直链下载
      const response = await fetch(task.url, { method: 'GET' })
      if (!response.ok) throw new Error(`下载失败: ${response.status}`)

      const contentLength = response.headers.get('content-length')
      const total = contentLength ? parseInt(contentLength, 10) : 0
      task.size = total

      const reader = response.body.getReader()
      const chunks = []
      let received = 0

      while (true) {
        const { done, value } = await reader.read()
        if (done) break
        chunks.push(value)
        received += value.length
        task.downloadedSize = received
        task.progress = total ? Math.round((received / total) * 50) : 0
        task.speed = received / ((Date.now() - new Date(task.createdAt).getTime()) / 1000 || 1)
        onProgress?.(task)
      }

      return new Blob(chunks)
    } else {
      // 蓝奏云分享解析
      throw new Error('蓝奏云分享同步待实现')
    }
  }

  // 上传到蓝奏云 (通过后端中转)
  async function uploadToLanzou(task, onProgress) {
    // 实际上传通过后端API
    // 这里只是模拟进度
    for (let i = 0; i <= 100; i += 10) {
      task.progress = 50 + Math.round(i * 0.5)
      onProgress?.(task)
      await new Promise(r => setTimeout(r, 100))
    }
  }

  // 删除任务
  function removeTask(id) {
    syncTasks.value = syncTasks.value.filter(t => t.id !== id)
  }

  // 清除已完成的任务
  function clearCompleted() {
    syncTasks.value = syncTasks.value.filter(t => t.status !== 'completed')
  }

  // 清除所有任务
  function clearAll() {
    if (isRunning.value) return
    syncTasks.value = []
  }

  // 重试失败任务
  function retryTask(id) {
    const task = syncTasks.value.find(t => t.id === id)
    if (task && task.status === 'error') {
      task.status = 'pending'
      task.error = null
      task.progress = 0
      task.step = 'waiting'
    }
  }

  return {
    syncTasks,
    isRunning,
    pendingCount,
    downloadingCount,
    uploadingCount,
    completedCount,
    errorCount,
    addTask,
    startSync,
    removeTask,
    clearCompleted,
    clearAll,
    retryTask
  }
})
