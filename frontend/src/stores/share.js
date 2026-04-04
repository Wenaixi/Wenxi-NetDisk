import { defineStore } from 'pinia'
import { ref } from 'vue'
import { shareAPI } from '../api'

export const useShareStore = defineStore('share', () => {
  const shares = ref([])
  const loading = ref(false)
  const error = ref(null)

  async function fetchShares() {
    loading.value = true
    error.value = null
    try {
      shares.value = await shareAPI.list()
    } catch (err) {
      error.value = err.message || '获取分享列表失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function createShare(fileId, options = {}) {
    error.value = null
    try {
      const share = await shareAPI.create({
        file_id: fileId,
        password: options.password || null,
        expires_at: options.expiresAt || null
      })
      return share
    } catch (err) {
      error.value = err.message || '创建分享失败'
      throw err
    }
  }

  async function getShare(token) {
    error.value = null
    try {
      return await shareAPI.get(token)
    } catch (err) {
      error.value = err.message || '获取分享失败'
      throw err
    }
  }

  async function deleteShare(id) {
    error.value = null
    try {
      await shareAPI.delete(id)
      shares.value = shares.value.filter(s => s.id !== id)
    } catch (err) {
      error.value = err.message || '删除分享失败'
      throw err
    }
  }

  function resetError() {
    error.value = null
  }

  return {
    shares,
    loading,
    error,
    fetchShares,
    createShare,
    getShare,
    deleteShare,
    resetError
  }
})
