import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { encryptFile, generateEncryptionKey, exportKey } from '../utils/crypto'
import { fileAPI } from '../api'

export const useUploadStore = defineStore('upload', () => {
  // State
  const isEncrypting = ref(false)
  const isUploading = ref(false)
  const progress = ref(0)
  const uploadedBytes = ref(0)
  const totalBytes = ref(0)
  const speed = ref(0)
  const error = ref(null)
  const currentFile = ref(null)

  // Getters
  const canUpload = computed(() => !isEncrypting.value && !isUploading.value)

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

      const { upload_url, session_id } = sessionResponse.data

      // Step 3: 上传加密后的文件
      const startTime = Date.now()
      const lastLoaded = { value: 0 }

      await uploadToUrl(
        upload_url,
        encryptedBlob,
        (loaded) => {
          uploadedBytes.value = loaded
          progress.value = Math.round((loaded / encryptedBlob.size) * 100)

          // Calculate speed
          const elapsed = (Date.now() - startTime) / 1000
          if (elapsed > 0) {
            const bytesPerSecond = (loaded - lastLoaded.value) / 0.5
            speed.value = bytesPerSecond
            lastLoaded.value = loaded
          }
        }
      )

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
    reset,
    encryptAndUpload,
  }
})
