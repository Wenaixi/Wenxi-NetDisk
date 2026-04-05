<template>
  <n-form label-placement="left" label-width="100" class="p-4">
    <n-form-item label="主题模式">
      <n-radio-group :value="themeStore.theme" @update:value="themeStore.setTheme">
        <n-space>
          <n-radio value="dark">深色</n-radio>
          <n-radio value="light">浅色</n-radio>
          <n-radio value="auto">跟随系统</n-radio>
        </n-space>
      </n-radio-group>
    </n-form-item>
    <n-form-item label="上传目录">
      <n-input-group>
        <n-input :value="uploadStore.uploadPath" readonly placeholder="未设置" />
        <n-button type="primary" @click="chooseUploadPath">选择</n-button>
      </n-input-group>
    </n-form-item>
    <n-form-item label="并发上传">
      <n-input-number :value="settings.concurrentUploads" @update:value="updateSetting('concurrentUploads', $event)" :min="1" :max="5" />
    </n-form-item>
    <n-form-item label="并发下载">
      <n-input-number :value="settings.concurrentDownloads" @update:value="updateSetting('concurrentDownloads', $event)" :min="1" :max="10" />
    </n-form-item>
    <n-form-item label="分块大小">
      <n-input-number :value="settings.chunkSize" @update:value="updateSetting('chunkSize', $event)" :min="1" :max="50" />
      <n-text depth="3" class="ml-2">MB</n-text>
    </n-form-item>
    <n-form-item label="自动清理回收站">
      <n-switch :value="settings.autoCleanRecycle" @update:value="updateSetting('autoCleanRecycle', $event)" />
      <n-text depth="3" class="ml-2">{{ settings.autoCleanRecycle ? `超过 ${settings.recycleRetentionDays} 天自动清理` : '已关闭' }}</n-text>
    </n-form-item>
    <n-form-item label="保留天数" v-if="settings.autoCleanRecycle">
      <n-input-number :value="settings.recycleRetentionDays" @update:value="updateSetting('recycleRetentionDays', $event)" :min="1" :max="365" />
      <n-text depth="3" class="ml-2">天</n-text>
    </n-form-item>
  </n-form>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useThemeStore } from '../stores/theme'
import { useUploadStore } from '../stores/upload'

const themeStore = useThemeStore()
const uploadStore = useUploadStore()

const DEFAULT_SETTINGS = {
  uploadPath: '',
  concurrentUploads: 3,
  concurrentDownloads: 5,
  chunkSize: 10,
  autoCleanRecycle: true,
  recycleRetentionDays: 30
}

const settings = ref({ ...DEFAULT_SETTINGS })

function loadSettings() {
  const saved = localStorage.getItem('wenxi-settings')
  if (saved) {
    try {
      settings.value = { ...DEFAULT_SETTINGS, ...JSON.parse(saved) }
    } catch {
      settings.value = { ...DEFAULT_SETTINGS }
    }
  }
}

function saveSettings() {
  localStorage.setItem('wenxi-settings', JSON.stringify(settings.value))
  uploadStore.uploadPath = settings.value.uploadPath
}

function updateSetting(key, value) {
  settings.value[key] = value
  saveSettings()
}

function chooseUploadPath() {
  // In browser, we cannot directly open a folder picker with path return
  // Use a fallback: prompt for path or use input type="file" webkitdirectory
  const input = document.createElement('input')
  input.type = 'file'
  input.webkitdirectory = true
  input.onchange = (e) => {
    const files = e.target.files
    if (files.length > 0) {
      const path = files[0].webkitRelativePath.split('/')[0]
      updateSetting('uploadPath', path)
    }
  }
  input.click()
}

onMounted(() => {
  loadSettings()
  uploadStore.uploadPath = settings.value.uploadPath
})
</script>
