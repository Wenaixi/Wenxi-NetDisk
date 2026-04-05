<template>
  <div class="p-6">
    <h2 class="text-xl font-bold text-white mb-4">已完成任务</h2>

    <div class="flex gap-4 mb-4">
      <n-button :type="activeTab === 'download' ? 'primary' : 'default'" @click="activeTab = 'download'">
        下载完成 ({{ downloadTasks.length }})
      </n-button>
      <n-button :type="activeTab === 'upload' ? 'primary' : 'default'" @click="activeTab = 'upload'">
        上传完成 ({{ uploadTasks.length }})
      </n-button>
      <n-button :type="activeTab === 'sync' ? 'primary' : 'default'" @click="activeTab = 'sync'">
        同步完成 ({{ syncTasks.length }})
      </n-button>
    </div>

    <div class="flex justify-end mb-4">
      <n-button type="error" @click="clearAll">清除全部记录</n-button>
    </div>

    <!-- 下载完成列表 -->
    <div v-if="activeTab === 'download'" class="space-y-2">
      <div v-if="downloadTasks.length === 0" class="text-center py-8 text-gray-500">
        暂无下载完成记录
      </div>
      <div
        v-for="task in downloadTasks"
        :key="task.id"
        class="bg-[#1a1a1a] p-3 flex items-center justify-between"
      >
        <div class="flex items-center gap-3">
          <n-icon size="20" color="#a78bfa"><Document /></n-icon>
          <div>
            <div class="text-white">{{ task.name || '未知文件' }}</div>
            <div class="text-gray-500 text-sm">{{ formatSize(task.size) }}</div>
          </div>
        </div>
        <div class="flex items-center gap-2">
          <span class="text-green-400 text-sm">已完成</span>
          <n-button size="small" @click="removeDownload(task.id)">
            <n-icon><Close /></n-icon>
          </n-button>
        </div>
      </div>
    </div>

    <!-- 上传完成列表 -->
    <div v-if="activeTab === 'upload'" class="space-y-2">
      <div v-if="uploadTasks.length === 0" class="text-center py-8 text-gray-500">
        暂无上传完成记录
      </div>
      <div
        v-for="task in uploadTasks"
        :key="task.id"
        class="bg-[#1a1a1a] p-3 flex items-center justify-between"
      >
        <div class="flex items-center gap-3">
          <n-icon size="20" color="#60a5fa"><Folder /></n-icon>
          <div>
            <div class="text-white">{{ task.name || '未知文件' }}</div>
            <div class="text-gray-500 text-sm">{{ formatSize(task.size) }}</div>
          </div>
        </div>
        <div class="flex items-center gap-2">
          <span class="text-green-400 text-sm">已完成</span>
          <n-button size="small" @click="removeUpload(task.id)">
            <n-icon><Close /></n-icon>
          </n-button>
        </div>
      </div>
    </div>

    <!-- 同步完成列表 -->
    <div v-if="activeTab === 'sync'" class="space-y-2">
      <div v-if="syncTasks.length === 0" class="text-center py-8 text-gray-500">
        暂无同步完成记录
      </div>
      <div
        v-for="task in syncTasks"
        :key="task.id"
        class="bg-[#1a1a1a] p-3 flex items-center justify-between"
      >
        <div class="flex items-center gap-3">
          <n-icon size="20" color="#34d399"><Cloud /></n-icon>
          <div>
            <div class="text-white">{{ task.name || '未知任务' }}</div>
            <div class="text-gray-500 text-sm">{{ task.url || '' }}</div>
          </div>
        </div>
        <div class="flex items-center gap-2">
          <span class="text-green-400 text-sm">已完成</span>
          <n-button size="small" @click="removeSync(task.id)">
            <n-icon><Close /></n-icon>
          </n-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useDownloadTaskStore } from '../stores/downloadTask'
import { useUploadTaskStore } from '../stores/uploadTask'
import { useSyncStore } from '../stores/sync'
import { Document, Folder, Cloud, Close } from '@vicons/ionicons5'
import { formatFileSize } from '../utils/fileSplit'

const downloadTaskStore = useDownloadTaskStore()
const uploadTaskStore = useUploadTaskStore()
const syncStore = useSyncStore()

const activeTab = ref('download')

function getDownloadTasks() {
  return downloadTaskStore.tasks.filter(t => t.status === 'completed')
}
function getUploadTasks() {
  return uploadTaskStore.tasks.filter(t => t.status === 'completed')
}
function getSyncTasks() {
  return syncStore.tasks.filter(t => t.status === 'completed')
}

const downloadTasks = ref(getDownloadTasks())
const uploadTasks = ref(getUploadTasks())
const syncTasks = ref(getSyncTasks())

function formatSize(bytes) {
  if (!bytes) return '0 B'
  return formatFileSize(bytes)
}

function removeDownload(id) {
  downloadTaskStore.removeTask(id)
  downloadTasks.value = getDownloadTasks()
}

function removeUpload(id) {
  uploadTaskStore.removeTask(id)
  uploadTasks.value = getUploadTasks()
}

function removeSync(id) {
  syncStore.removeTask(id)
  syncTasks.value = getSyncTasks()
}

function clearAll() {
  if (activeTab.value === 'download') {
    downloadTaskStore.clearCompleted()
    downloadTasks.value = getDownloadTasks()
  } else if (activeTab.value === 'upload') {
    uploadTaskStore.clearCompleted()
    uploadTasks.value = getUploadTasks()
  } else {
    syncStore.clearCompleted()
    syncTasks.value = getSyncTasks()
  }
}
</script>
