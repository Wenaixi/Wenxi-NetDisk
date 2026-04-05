import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { encryptFile, generateEncryptionKey, exportKey } from '../utils/crypto'
import { fileAPI } from '../api'
import { useUploadTaskStore } from './uploadTask'

// 分块大小 2MB
const ChunkSize = 2 * 1024 * 1024

export const useUploadStore = defineStore('upload', () => {
  // State
  const isEncrypting = ref(false)
  const isUploading = ref(false)
  const progress = ref(0)
  const uploadedBytes = ref(0)
  const uploadedBytesRef = ref(0) // 用于加密上传函数内部追踪进度
  const totalBytes = ref(0)
  const speed = ref(0)
  const error = ref(null)
  const currentFile = ref(null)

  // 批量上传队列
  const uploadQueue = ref([])
  const isQueueUploading = ref(false)
  const queueCurrentIndex = ref(0)
  const queueTotalCount = ref(0)
  const queueUploadedCount = ref(0)
  const queueErrors = ref([])

  // Getters
  const canUpload = computed(() => !isEncrypting.value && !isUploading.value)
  const hasQueueItems = computed(() => uploadQueue.value.length > 0)

  // Actions
  function reset() {
    isEncrypting.value = false
    isUploading.value = false
    progress.value = 0
    uploadedBytes.value = 0
    totalBytes.value = 0
    speed.value = 0
    error.value = null
    currentFile.value = null
  }

  function addToQueue(files) {
    for (const file of files) {
      uploadQueue.value.push({
        id: Date.now() + Math.random(),
        file,
        status: 'pending', // pending, encrypting, uploading, completed, error
        progress: 0,
        error: null
      })
    }
  }

  function removeFromQueue(id) {
    const idx = uploadQueue.value.findIndex(f => f.id === id)
    if (idx !== -1) {
      uploadQueue.value.splice(idx, 1)
    }
  }

  function clearQueue() {
    uploadQueue.value = []
    queueErrors.value = []
    queueCurrentIndex.value = 0
    queueUploadedCount.value = 0
  }

  async function uploadAll(key, folderId = null) {
    if (uploadQueue.value.length === 0) return

    const uploadTaskStore = useUploadTaskStore()
    isQueueUploading.value = true
    queueTotalCount.value = uploadQueue.value.length
    queueUploadedCount.value = 0
    queueErrors.value = []

    for (let i = 0; i < uploadQueue.value.length; i++) {
      const item = uploadQueue.value[i]
      if (item.status === 'completed') continue

      queueCurrentIndex.value = i
      item.status = 'encrypting'

      // 创建任务记录
      const taskId = item.id.toString()
      uploadTaskStore.addTask({
        id: taskId,
        name: item.file.name,
        size: item.file.size,
        status: 'pending',
      })

      try {
        const file = item.file
        const totalChunks = Math.ceil(file.size / ChunkSize)

        // 加密文件(按块加密上传)
        let encryptedBlob
        if (file.size > ChunkSize) {
          // 大文件按块加密
          encryptedBlob = await encryptFileInChunks(file, key, totalChunks, (chunkProgress) => {
            item.progress = Math.round(chunkProgress * 30) // 加密占30%进度
            uploadTaskStore.updateProgress(taskId, chunkProgress * 30)
          })
        } else {
          encryptedBlob = await encryptFile(file, key)
        }

        item.status = 'uploading'
        uploadTaskStore.updateProgress(taskId, 30)

        // 初始化上传
        const sessionResponse = await fileAPI.initializeUpload({
          file_name: file.name,
          file_size: encryptedBlob.size,
          mime_type: file.type || 'application/octet-stream',
          folder_id: folderId,
        })

        const { session_id, total_chunks } = sessionResponse

        // 上传(分块)
        if (encryptedBlob.size > ChunkSize) {
          let uploadedBytes = 0
          for (let i = 0; i < total_chunks; i++) {
            const start = i * ChunkSize
            const end = Math.min(start + ChunkSize, encryptedBlob.size)
            const chunk = encryptedBlob.slice(start, end)

            const formData = new FormData()
            formData.append('file', chunk)
            formData.append('chunk_index', String(i))
            formData.append('folder_id', String(folderId || -1))

            await fileAPI.uploadChunk(session_id, formData)

            uploadedBytes += chunk.size
            const uploadProgress = 30 + Math.round((uploadedBytes / encryptedBlob.size) * 70)
            item.progress = uploadProgress
            uploadTaskStore.updateProgress(taskId, uploadProgress)
          }
        } else {
          // 小文件直接作为单个分块上传
          const formData = new FormData()
          formData.append('file', encryptedBlob, file.name)
          formData.append('chunk_index', '0')
          formData.append('folder_id', String(folderId || -1))
          await fileAPI.uploadChunk(session_id, formData)
          item.progress = 100
          uploadTaskStore.updateProgress(taskId, 100)
        }

        // 完成上传
        const exportedKey = await exportKey(key)
        await fileAPI.completeUpload(session_id, {
          encryption_key: exportedKey.key,
          encryption_nonce: exportedKey.iv,
        })

        item.status = 'completed'
        uploadTaskStore.completeTask(taskId)
        queueUploadedCount.value++
      } catch (err) {
        item.status = 'error'
        item.error = err.message || '上传失败'
        uploadTaskStore.failTask(taskId, err.message)
        queueErrors.value.push({ file: item.file.name, error: item.error })
      }
    }

    isQueueUploading.value = false
  }

  async function encryptAndUpload(file, key, folderId = null) {
    if (!file) throw new Error('请选择文件')

    error.value = null
    currentFile.value = file
    totalBytes.value = file.size

    try {
      // Step 1: 加密文件
      isEncrypting.value = true
      const encryptedBlob = await encryptFile(file, key)
      isEncrypting.value = false

      // Step 2: 初始化上传会话
      isUploading.value = true
      progress.value = 0

      const sessionResponse = await fileAPI.initializeUpload({
        file_name: file.name,
        file_size: encryptedBlob.size,
        mime_type: file.type || 'application/octet-stream',
        folder_id: folderId,
      })

      const { session_id, total_chunks } = sessionResponse

      // 上传分块
      if (encryptedBlob.size > ChunkSize) {
        let uploadedBytes = 0
        const startTime = Date.now()
        const lastLoaded = { value: 0 }

        for (let i = 0; i < total_chunks; i++) {
          const start = i * ChunkSize
          const end = Math.min(start + ChunkSize, encryptedBlob.size)
          const chunk = encryptedBlob.slice(start, end)

          const formData = new FormData()
          formData.append('file', chunk)
          formData.append('chunk_index', String(i))
          formData.append('folder_id', String(folderId || -1))

          await fileAPI.uploadChunk(session_id, formData)

          uploadedBytes += chunk.size
          uploadedBytesRef.value = uploadedBytes
          progress.value = Math.round((uploadedBytes / encryptedBlob.size) * 100)

          // Calculate speed
          const elapsed = (Date.now() - startTime) / 1000
          if (elapsed > 0) {
            const bytesPerSecond = (uploadedBytes - lastLoaded.value) / 0.5
            speed.value = bytesPerSecond
            lastLoaded.value = uploadedBytes
          }
        }
      } else {
        // 小文件直接作为单个分块上传
        const formData = new FormData()
        formData.append('file', encryptedBlob, file.name)
        formData.append('chunk_index', '0')
        formData.append('folder_id', String(folderId || -1))
        await fileAPI.uploadChunk(session_id, formData)
        progress.value = 100
      }

      // Step 4: 导出密钥并完成上传
      const exportedKey = await exportKey(key)

      await fileAPI.completeUpload(session_id, {
        encryption_key: exportedKey.key,
        encryption_nonce: exportedKey.iv,
      })

      progress.value = 100
      isUploading.value = false

      return {
        sessionId: session_id,
        fileName: file.name,
        encryptedSize: encryptedBlob.size,
      }
    } catch (err) {
      isEncrypting.value = false
      isUploading.value = false
      error.value = err.message || '上传失败'
      throw err
    }
  }

  /**
   * 按块加密大文件
   * @param {File} file 原始文件
   * @param {CryptoKey} key 加密密钥
   * @param {number} totalChunks 总分块数
   * @param {function} onProgress 进度回调 (0-1)
   * @returns {Promise<Blob>} 加密后的Blob
   */
  async function encryptFileInChunks(file, key, totalChunks, onProgress) {
    const parts = []
    for (let i = 0; i < totalChunks; i++) {
      const start = i * ChunkSize
      const end = Math.min(start + ChunkSize, file.size)
      const chunk = file.slice(start, end)
      const encryptedChunk = await encryptFile(chunk, key)
      parts.push(encryptedChunk)
      onProgress?.((i + 1) / totalChunks)
    }
    return new Blob(parts, { type: 'application/octet-stream' })
  }

  /**
   * 按块上传大文件(断点续传支持)
   * @param {number} sessionId 上传会话ID
   * @param {number} folderId 目标文件夹ID
   * @param {Blob} blob 加密后的Blob
   * @param {number} totalChunks 总分块数
   * @param {function} onProgress 进度回调 (0-1)
   */
  async function uploadInChunks(sessionId, folderId, blob, totalChunks, onProgress) {
    let uploadedBytes = 0

    for (let i = 0; i < totalChunks; i++) {
      const start = i * ChunkSize
      const end = Math.min(start + ChunkSize, blob.size)
      const chunk = blob.slice(start, end)

      const formData = new FormData()
      formData.append('file', chunk)
      formData.append('chunk_index', String(i))
      formData.append('folder_id', String(folderId || -1))

      await fileAPI.uploadChunk(sessionId, formData)

      uploadedBytes += chunk.size
      onProgress?.(uploadedBytes / blob.size)
    }
  }

  async function uploadToUrl(url, blob, onProgress) {
    return new Promise((resolve, reject) => {
      const xhr = new XMLHttpRequest()

      xhr.upload.addEventListener('progress', (event) => {
        if (event.lengthComputable) {
          onProgress?.(event.loaded)
        }
      })

      xhr.addEventListener('load', () => {
        if (xhr.status >= 200 && xhr.status < 300) {
          resolve(xhr.response)
        } else {
          reject(new Error(`上传失败: ${xhr.status}`))
        }
      })

      xhr.addEventListener('error', () => {
        reject(new Error('网络错误'))
      })

      xhr.addEventListener('abort', () => {
        reject(new Error('上传已取消'))
      })

      xhr.open('PUT', url, true)
      xhr.setRequestHeader('Content-Type', 'application/octet-stream')
      xhr.send(blob)
    })
  }

  return {
    isEncrypting,
    isUploading,
    progress,
    uploadedBytes,
    totalBytes,
    speed,
    error,
    currentFile,
    canUpload,
    // 批量上传
    uploadQueue,
    isQueueUploading,
    queueCurrentIndex,
    queueTotalCount,
    queueUploadedCount,
    queueErrors,
    hasQueueItems,
    addToQueue,
    removeFromQueue,
    clearQueue,
    uploadAll,
    // 单文件上传
    reset,
    encryptAndUpload,
  }
})
