<template>
  <div class="min-h-screen bg-[#0f0f0f]">
    <AppHeader />

    <main class="p-6">
      <div class="flex items-center justify-between mb-6">
        <h2 class="text-white text-xl font-medium">下载任务</h2>
        <div class="flex gap-2">
          <n-button @click="downloadStore.pauseAll">全部暂停</n-button>
          <n-button type="primary" @click="downloadStore.resumeAll">全部开始</n-button>
          <n-button @click="downloadStore.clearCompleted">清除已完成</n-button>
          <n-button type="primary" @click="showBatchModal = true">批量下载</n-button>
        </div>
      </div>

      <!-- 统计栏 -->
      <div class="flex gap-4 mb-6">
        <n-tag :bordered="false" type="info">等待中: {{ downloadStore.pendingCount }}</n-tag>
        <n-tag :bordered="false" type="warning">下载中: {{ downloadStore.activeCount }}</n-tag>
        <n-tag :bordered="false" type="success">已完成: {{ downloadStore.completedCount }}</n-tag>
        <n-tag :bordered="false" type="error">失败: {{ downloadStore.errorCount }}</n-tag>
      </div>

      <!-- 任务列表 -->
      <div v-if="downloadStore.tasks.length === 0" class="text-center py-12">
        <n-empty description="暂无下载任务" />
      </div>

      <div v-else class="space-y-3">
        <div
          v-for="task in downloadStore.tasks"
          :key="task.id"
          class="bg-[#1a1a1a] p-4"
        >
          <div class="flex items-center justify-between mb-2">
            <div class="flex items-center gap-2 flex-1 min-w-0">
              <n-icon size="20" color="#60a5fa"><CloudDownload /></n-icon>
              <span class="text-white truncate">{{ task.name }}</span>
              <span class="text-gray-500 text-sm">{{ formatSize(task.size) }}</span>
            </div>
            <div class="flex items-center gap-2">
              <n-tag v-if="task.status === 'completed'" type="success" size="small">完成</n-tag>
              <n-tag v-else-if="task.status === 'error'" type="error" size="small">失败</n-tag>
              <n-tag v-else-if="task.status === 'downloading'" type="warning" size="small">下载中</n-tag>
              <n-tag v-else type="default" size="small">等待中</n-tag>

              <n-button
                size="tiny"
                @click="toggleTask(task)"
                v-if="task.status !== 'completed' && task.status !== 'error'"
              >
                {{ task.status === 'downloading' ? '暂停' : '恢复' }}
              </n-button>
              <n-button
                size="tiny"
                type="error"
                text
                @click="downloadStore.removeTask(task.id)"
              >
                <n-icon><Close /></n-icon>
              </n-button>
            </div>
          </div>

          <n-progress
            v-if="task.status !== 'completed' && task.status !== 'error'"
            :percentage="task.progress"
            :show-indicator="false"
            :height="4"
          />
          <div v-if="task.status === 'error'" class="text-red-400 text-xs mt-1">{{ task.error }}</div>
          <div v-if="task.speed" class="text-gray-500 text-xs mt-1">{{ task.speed }}</div>
        </div>
      </div>
    </main>

    <!-- 批量下载模态框 -->
    <n-modal v-model:show="showBatchModal">
      <n-card title="批量下载" style="width: 600px;">
        <div class="space-y-4">
          <n-input
            v-model:value="batchUrls"
            type="textarea"
            placeholder="每行一个链接: https://example.com/file.zip&#10;或: 文件名 https://lanzoui.com/xxx 密码:1234"
            :rows="4"
          />
        </div>
        <template #footer>
          <div class="flex justify-end gap-2">
            <n-button @click="showBatchModal = false">取消</n-button>
            <n-button type="primary" @click="addBatchTasks" :disabled="!batchUrls">添加</n-button>
          </div>
        </template>
      </n-card>
    </n-modal>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useDownloadTaskStore } from '../stores/downloadTask'
import { useMessage } from 'naive-ui'
import AppHeader from '../components/AppHeader.vue'
import { CloudDownload, Close } from '@vicons/ionicons5'

const downloadStore = useDownloadTaskStore()
const message = useMessage()
const showBatchModal = ref(false)
const batchUrls = ref('')

function formatSize(bytes) {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

function toggleTask(task) {
  if (task.status === 'downloading') {
    downloadStore.pauseTask(task.id)
  } else {
    downloadStore.resumeTask(task.id)
  }
}

function addBatchTasks() {
  if (!batchUrls.value.trim()) return

  const lines = batchUrls.value.trim().split('\n').filter(line => line.includes('http'))
  let addedCount = 0

  for (const line of lines) {
    const parts = line.trim().split(/\s+/)
    let url = ''
    let name = ''
    let pwd = ''

    for (const part of parts) {
      if (part.startsWith('http')) {
        url = part
      } else if (part.toLowerCase() === '密码:' && parts.indexOf(part) + 1 < parts.length) {
        pwd = parts[parts.indexOf(part) + 1].replace(/密码:/, '')
      }
    }

    if (!url) continue

    // Try to extract name from parts that are not URL or password
    const nonUrlParts = parts.filter(p => !p.startsWith('http') && p !== '密码:' && !p.startsWith('密码:'))
    name = nonUrlParts.length > 0 ? nonUrlParts.join(' ') : url.split('/').pop() || 'unknown'

    downloadStore.addTask({
      name,
      url,
      pwd,
      size: 0,
    })
    addedCount++
  }

  if (addedCount > 0) {
    message.success(`已添加 ${addedCount} 个下载任务`)
    batchUrls.value = ''
    showBatchModal.value = false
  } else {
    message.warning('未找到有效的下载链接')
  }
}
</script>
