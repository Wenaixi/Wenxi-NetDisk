import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { fileAPI } from '../api'
import { decryptFile, importKey } from '../utils/crypto'

export const useDownloadStore = defineStore('download', () => {
  const isDownloading = ref(false)
  const isDecrypting = ref(false)
  const progress = ref(0)
  const currentFile = ref(null)
  const error = ref(null)
  const downloadSpeed = ref(0)

  // 批量下载队列
  const downloadQueue = ref([])
  const isQueueDownloading = ref(false)
  const queueCurrentIndex = ref(0)
  const queueTotalCount = ref(0)
  const queueCompletedCount = ref(0)
  const queueErrors = ref([])

  // Getters
  const hasQueueItems = computed(() => downloadQueue.value.length > 0)

  async function downloadAndDecrypt(fileId, encryptionKey, encryptionNonce) {
    if (!fileId) {
      throw new Error('文件ID不能为空')
    }

    error.value = null
    isDownloading.value = true
    progress.value = 0

    try {
      // Step 1: 获取下载链接和文件信息
      const response = await fileAPI.download(fileId)
      const { download_url, file_name, file_size } = response

      currentFile.value = {
        id: fileId,
        name: file_name,
        size: file_size
      }

      // Step 2: 下载加密文件
      const encryptedBlob = await downloadFile(download_url, (loaded) => {
        progress.value = Math.round((loaded / file_size) * 50) // 0-50% for download
      })

      // Step 3: 解密文件
      isDecrypting.value = true
      isDownloading.value = false
      progress.value = 50

      const key = await importKey(encryptionKey)
      const decryptedBlob = await decryptFile(encryptedBlob, key)

      isDecrypting.value = false
      progress.value = 100

      // Step 4: 保存文件
      saveBlob(decryptedBlob, file_name)

      return {
        fileName: file_name,
        size: file_size,
        blob: decryptedBlob
      }
    } catch (err) {
      error.value = err.message || '下载失败'
      isDownloading.value = false
      isDecrypting.value = false
      throw err
    }
  }

  // 批量下载功能
  function addToQueue(files) {
    // files: [{ id, name, encryptionKey, encryptionNonce, size }]
    for (const file of files) {
      downloadQueue.value.push({
        id: file.id || (Date.now() + Math.random()),
        name: file.name || `file_${Date.now()}`,
        encryptionKey: file.encryptionKey,
        encryptionNonce: file.encryptionNonce,
        size: file.size,
        status: 'pending', // pending, downloading, decrypting, completed, error
        progress: 0,
        error: null
      })
    }
  }

  function removeFromQueue(id) {
    const idx = downloadQueue.value.findIndex(f => f.id === id)
    if (idx !== -1) {
      downloadQueue.value.splice(idx, 1)
    }
  }

  function clearQueue() {
    downloadQueue.value = []
    queueErrors.value = []
    queueCurrentIndex.value = 0
    queueCompletedCount.value = 0
  }

  async function downloadAll() {
    if (downloadQueue.value.length === 0) return

    isQueueDownloading.value = true
    queueTotalCount.value = downloadQueue.value.length
    queueCompletedCount.value = 0
    queueErrors.value = []

    for (let i = 0; i < downloadQueue.value.length; i++) {
      const item = downloadQueue.value[i]
      if (item.status === 'completed') continue

      queueCurrentIndex.value = i
      item.status = 'downloading'

      try {
        // 获取下载链接
        const response = await fileAPI.download(item.id)
        const { download_url, file_name, file_size } = response
        item.name = file_name || item.name

        // 下载文件
        const encryptedBlob = await downloadFile(download_url, (loaded) => {
          if (file_size > 0) {
            item.progress = Math.round((loaded / file_size) * 50)
          }
        })

        // 解密
        item.status = 'decrypting'
        item.progress = 50

        if (item.encryptionKey) {
          const key = await importKey(item.encryptionKey)
          const decryptedBlob = await decryptFile(encryptedBlob, key)
          saveBlob(decryptedBlob, item.name)
        } else {
          // 无加密直接保存
          saveBlob(encryptedBlob, item.name)
        }

        item.status = 'completed'
        item.progress = 100
        queueCompletedCount.value++
      } catch (err) {
        item.status = 'error'
        item.error = err.message || '下载失败'
        queueErrors.value.push({ file: item.name, error: item.error })
      }
    }

    isQueueDownloading.value = false
  }

  async function downloadFile(url, onProgress) {
    return new Promise((resolve, reject) => {
      const xhr = new XMLHttpRequest()

      xhr.addEventListener('progress', (event) => {
        if (event.lengthComputable) {
          onProgress?.(event.loaded)
        }
      })

      xhr.addEventListener('load', () => {
        if (xhr.status >= 200 && xhr.status < 300) {
          const blob = new Blob([xhr.response])
          resolve(blob)
        } else {
          reject(new Error(`下载失败: ${xhr.status}`))
        }
      })

      xhr.addEventListener('error', () => {
        reject(new Error('网络错误'))
      })

      xhr.addEventListener('abort', () => {
        reject(new Error('下载已取消'))
      })

      xhr.responseType = 'blob'
      xhr.open('GET', url, true)
      xhr.send()
    })
  }

  function saveBlob(blob, fileName) {
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = fileName
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
  }

  function reset() {
    isDownloading.value = false
    isDecrypting.value = false
    progress.value = 0
    currentFile.value = null
    error.value = null
    downloadSpeed.value = 0
  }

  return {
    isDownloading,
    isDecrypting,
    progress,
    currentFile,
    error,
    downloadSpeed,
    // 批量下载
    downloadQueue,
    isQueueDownloading,
    queueCurrentIndex,
    queueTotalCount,
    queueCompletedCount,
    queueErrors,
    hasQueueItems,
    addToQueue,
    removeFromQueue,
    clearQueue,
    downloadAll,
    // 单文件下载
    downloadAndDecrypt,
    reset
  }
})
