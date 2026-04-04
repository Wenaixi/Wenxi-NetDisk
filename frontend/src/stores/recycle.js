import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { recycleAPI } from '../api/recycle'

export const useRecycleStore = defineStore('recycle', () => {
  // State
  const items = ref([])
  const loading = ref(false)
  const error = ref(null)
  const selectedItems = ref([])
  const selectMode = ref(false)

  // Getters
  const fileCount = computed(() => items.value.filter(i => i.item_type === 'file').length)
  const folderCount = computed(() => items.value.filter(i => i.item_type === 'folder').length)
  const totalCount = computed(() => items.value.length)
  const hasSelectedItems = computed(() => selectedItems.value.length > 0)

  // Actions
  async function fetchList() {
    loading.value = true
    error.value = null
    try {
      const response = await recycleAPI.list()
      items.value = response.data || response.items || []
    } catch (err) {
      error.value = err.message || '获取回收站列表失败'
    } finally {
      loading.value = false
    }
  }

  async function restore(id) {
    loading.value = true
    try {
      await recycleAPI.restoreFile(id)
      items.value = items.value.filter(i => i.id !== id)
      selectedItems.value = selectedItems.value.filter(i => i !== id)
    } catch (err) {
      error.value = err.message || '恢复失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function permanentDelete(id) {
    loading.value = true
    try {
      await recycleAPI.permanentDelete(id)
      items.value = items.value.filter(i => i.id !== id)
      selectedItems.value = selectedItems.value.filter(i => i !== id)
    } catch (err) {
      error.value = err.message || '永久删除失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  function toggleSelect(id) {
    const idx = selectedItems.value.indexOf(id)
    if (idx === -1) {
      selectedItems.value.push(id)
    } else {
      selectedItems.value.splice(idx, 1)
    }
  }

  async function batchRestore() {
    if (selectedItems.value.length === 0) return
    const ids = [...selectedItems.value]
    for (const id of ids) {
      await restore(id)
    }
    selectedItems.value = []
  }

  async function batchDelete() {
    if (selectedItems.value.length === 0) return
    const ids = [...selectedItems.value]
    for (const id of ids) {
      try {
        await permanentDelete(id)
      } catch (e) {
        // continue
      }
    }
    selectedItems.value = []
  }

  async function clearAll() {
    loading.value = true
    try {
      await recycleAPI.clearAll()
      items.value = []
      selectedItems.value = []
    } catch (err) {
      error.value = err.message || '清空回收站失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  function formatSize(bytes) {
    if (bytes === 0) return '0 B'
    const units = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(1024))
    return (bytes / Math.pow(1024, i)).toFixed(2) + ' ' + units[i]
  }

  function formatDate(dateStr) {
    if (!dateStr) return '--'
    const d = new Date(dateStr)
    return d.toLocaleDateString('zh-CN')
  }

  return {
    items,
    loading,
    error,
    selectedItems,
    selectMode,
    fileCount,
    folderCount,
    totalCount,
    hasSelectedItems,
    fetchList,
    restore,
    permanentDelete,
    toggleSelect,
    batchRestore,
    batchDelete,
    clearAll,
    formatSize,
    formatDate
  }
})
