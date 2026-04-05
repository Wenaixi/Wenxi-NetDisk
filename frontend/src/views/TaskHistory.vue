<template>
  <div class="p-6">
    <div class="flex justify-between items-center mb-6">
      <h2 class="text-2xl font-bold text-white">任务历史</h2>
      <div class="flex gap-2">
        <n-button size="small" @click="taskHistoryStore.clearCompleted">清理已完成</n-button>
        <n-button size="small" type="error" @click="taskHistoryStore.clearAll">清空历史</n-button>
      </div>
    </div>

    <div v-if="taskHistoryStore.tasks.length === 0" class="text-center py-20 text-gray-400">
      <p class="text-lg">暂无任务历史记录</p>
    </div>

    <div v-else class="divide-y divide-gray-700">
      <div v-for="task in taskHistoryStore.tasks" :key="task.id"
           class="flex items-center justify-between py-3 px-4 hover:bg-[#1a1a1a] transition-colors">
        <div class="flex items-center gap-3">
          <n-tag :type="task.type === 'upload' ? 'success' : task.type === 'download' ? 'info' : 'warning'" size="small">
            {{ task.type === 'upload' ? '上传' : task.type === 'download' ? '下载' : '同步' }}
          </n-tag>
          <span class="text-white">{{ task.fileName }}</span>
          <span class="text-gray-500 text-sm">{{ formatSize(task.fileSize) }}</span>
        </div>
        <div class="flex items-center gap-3">
          <n-tag :type="task.status === 'completed' ? 'success' : 'error'" size="small">
            {{ task.status === 'completed' ? '完成' : '失败' }}
          </n-tag>
          <span class="text-gray-400 text-sm">{{ formatDate(task.createdAt) }}</span>
          <n-button text type="error" size="small" @click="taskHistoryStore.removeTask(task.id)">删除</n-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useTaskHistoryStore } from '../stores/taskHistory'

const taskHistoryStore = useTaskHistoryStore()

function formatSize(bytes) {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

function formatDate(dateStr) {
  if (!dateStr) return '未知'
  const d = new Date(dateStr)
  return d.toLocaleString('zh-CN', {
    year: 'numeric', month: '2-digit', day: '2-digit',
    hour: '2-digit', minute: '2-digit'
  })
}
</script>
