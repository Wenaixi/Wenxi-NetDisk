import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '../api'

export const useShareParseStore = defineStore('shareParse', () => {
  // State
  const loading = ref(false)
  const error = ref(null)
  const parsedShares = ref([])
  const currentShare = ref(null)
  const downloadQueue = ref([])

  // Getters
  const hasParsedShares = computed(() => parsedShares.value.length > 0)
  const hasDownloadItems = computed(() => downloadQueue.value.length > 0)

  /**
   * 解析分享链接
   * @param {string} url 分享链接
   * @param {string} pwd 密码（可选）
   */
  async function parseShare(url, pwd = '') {
    if (!url || !url.startsWith('http')) {
      throw new Error('请输入有效的分享链接')
    }

    loading.value = true
    error.value = null

    try {
      const response = await api.post('/lanzou/share/parse', { url, pwd })

      if (response.data && response.data.list) {
        currentShare.value = response.data
        parsedShares.value = response.data.list
      } else {
        throw new Error('解析失败，请检查链接是否正确')
      }
    } catch (err) {
      error.value = err.message || '解析失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  /**
   * 验证分享链接
   */
  async function validateShare(url) {
    try {
      const response = await api.post('/lanzou/share/validate', { url })
      return response.data && response.data.valid
    } catch {
      return false
    }
  }

  /**
   * 添加解析结果到下载队列
   */
  function addToDownload(items) {
    const itemsToAdd = Array.isArray(items) ? items : [items]
    downloadQueue.value.push(...itemsToAdd.map(item => ({
      ...item,
      status: 'pending', // pending, downloading, completed, error
      progress: 0,
      error: null
    })))
  }

  /**
   * 下载分享链接中的文件
   */
  async function downloadShareFile(shareUrl, pwd = '') {
    try {
      const response = await api.post('/lanzou/share/download', { url: shareUrl, pwd })
      return response.data?.download_url || response.download_url
    } catch (err) {
      error.value = err.message || '获取下载链接失败'
      throw err
    }
  }

  /**
   * 清除当前解析结果
   */
  function clearParsed() {
    currentShare.value = null
    parsedShares.value = []
    error.value = null
  }

  /**
   * 清除下载队列
   */
  function clearDownloadQueue() {
    downloadQueue.value = []
  }

  return {
    loading,
    error,
    parsedShares,
    currentShare,
    downloadQueue,
    hasParsedShares,
    hasDownloadItems,
    parseShare,
    validateShare,
    addToDownload,
    downloadShareFile,
    clearParsed,
    clearDownloadQueue
  }
})
