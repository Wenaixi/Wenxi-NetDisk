import { defineStore } from 'pinia'
import { ref } from 'vue'
import { fileAPI } from '../api'
import { decryptFile, importKey } from '../utils/crypto'

export const useDownloadStore = defineStore('download', () => {
  const isDownloading = ref(false)
  const isDecrypting = ref(false)
  const progress = ref(0)
  const currentFile = ref(null)
  const error = ref(null)
  const downloadSpeed = ref(0)

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
    downloadAndDecrypt,
    reset
  }
})
