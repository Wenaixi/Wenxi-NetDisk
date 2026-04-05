import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

const CALCULATE_KEY = 'wenxi-calculate'

export const useCalculateStore = defineStore('calculate', () => {
  // 每日上传流量记录 { 'YYYY-MM-DD': bytes }
  const dailyRecords = ref({})
  // 上传流量警告阈值（字节）默认7GB
  const warningSize = ref(7 * 1024 * 1024 * 1024)
  // 是否启用流量警告
  const warningEnabled = ref(true)

  const today = new Date().toISOString().slice(0, 10)
  const todayBytes = computed(() => dailyRecords.value[today] || 0)

  // 加载记录
  function loadRecords() {
    try {
      const saved = localStorage.getItem(CALCULATE_KEY)
      if (saved) {
        const data = JSON.parse(saved)
        dailyRecords.value = data.dailyRecords || {}
        warningSize.value = data.warningSize ?? 7 * 1024 * 1024 * 1024
        warningEnabled.value = data.warningEnabled ?? true
      }
    } catch {
      dailyRecords.value = {}
      warningSize.value = 7 * 1024 * 1024 * 1024
      warningEnabled.value = true
    }
  }

  // 保存记录
  function saveRecords() {
    // 清理超过30天的记录
    const cutoff = new Date()
    cutoff.setDate(cutoff.getDate() - 30)
    const cutoffStr = cutoff.toISOString().slice(0, 10)

    const cleaned = {}
    for (const [date, bytes] of Object.entries(dailyRecords.value)) {
      if (date >= cutoffStr) {
        cleaned[date] = bytes
      }
    }
    dailyRecords.value = cleaned

    localStorage.setItem(CALCULATE_KEY, JSON.stringify({
      dailyRecords: cleaned,
      warningSize: warningSize.value,
      warningEnabled: warningEnabled.value
    }))
  }

  // 记录上传流量
  function recordUploadBytes(bytes) {
    if (!dailyRecords.value[today]) {
      dailyRecords.value[today] = 0
    }
    dailyRecords.value[today] += bytes
    saveRecords()
  }

  // 检查是否超过警告阈值
  function checkWarningSize() {
    if (!warningEnabled.value) return false
    return todayBytes.value >= warningSize.value
  }

  // 设置警告阈值
  function setWarningSize(bytes) {
    warningSize.value = bytes
    saveRecords()
  }

  // 切换警告开关
  function setWarningEnabled(enabled) {
    warningEnabled.value = enabled
    saveRecords()
  }

  // 清空历史记录
  function clearHistory() {
    dailyRecords.value = {}
    saveRecords()
  }

  loadRecords()

  return {
    dailyRecords,
    warningSize,
    warningEnabled,
    todayBytes,
    recordUploadBytes,
    checkWarningSize,
    setWarningSize,
    setWarningEnabled,
    clearHistory
  }
})
