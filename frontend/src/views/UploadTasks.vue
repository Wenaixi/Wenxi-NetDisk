<template>
  <div class="min-h-screen bg-[#0f0f0f]">
    <AppHeader />

    <main class="p-6">
      <div class="flex items-center justify-between mb-6">
        <h2 class="text-white text-xl font-medium">上传任务</h2>
        <div class="flex gap-2">
          <n-button @click="uploadStore.pauseAll">全部暂停</n-button>
          <n-button type="primary" @click="uploadStore.resumeAll">全部开始</n-button>
          <n-button @click="uploadStore.clearCompleted">清除已完成</n-button>
        </div>
      </div>

      <!-- 统计栏 -->
      <div class="flex gap-4 mb-6">
        <n-tag :bordered="false" type="info">等待中: {{ uploadStore.pendingCount }}</n-tag>
        <n-tag :bordered="false" type="warning">上传中: {{ uploadStore.activeCount }}</n-tag>
        <n-tag :bordered="false" type="success">已完成: {{ uploadStore.completedCount }}</n-tag>
        <n-tag :bordered="false" type="error">失败: {{ uploadStore.errorCount }}</n-tag>
      </div>

      <!-- 任务列表 -->
      <div v-if="uploadStore.tasks.length === 0" class="text-center py-12">
        <n-empty description="暂无上传任务" />
      </div>

      <div v-else class="space-y-3">
        <div
          v-for="task in uploadStore.tasks"
          :key="task.id"
          class="bg-[#1a1a1a] p-4"
        >
          <div class="flex items-center justify-between mb-2">
            <div class="flex items-center gap-2 flex-1 min-w-0">
              <n-icon size="20" color="#a78bfa"><Document /></n-icon>
              <span class="text-white truncate">{{ task.name }}</span>
              <span class="text-gray-500 text-sm">{{ formatSize(task.size) }}</span>
            </div>
            <div class="flex items-center gap-2">
              <n-tag v-if="task.status === 'completed'" type="success" size="small">完成</n-tag>
              <n-tag v-else-if="task.status === 'error'" type="error" size="small">失败</n-tag>
              <n-tag v-else-if="task.status === 'uploading'" type="warning" size="small">上传中</n-tag>
              <n-tag v-else type="default" size="small">等待中</n-tag>

              <n-button
                size="tiny"
                @click="toggleTask(task)"
                v-if="task.status !== 'completed'"
              >
                {{ task.status === 'uploading' ? '暂停' : '恢复' }}
              </n-button>
              <n-button
                size="tiny"
                type="error"
                text
                @click="uploadStore.removeTask(task.id)"
              >
                <n-icon><Close /></n-icon>
              </n-button>
            </div>
          </div>

          <!-- 进度条 -->
          <n-progress
            v-if="task.status !== 'completed' && task.status !== 'error'"
            :percentage="task.progress"
            :show-indicator="false"
            :height="4"
          />
          <div v-if="task.status === 'error'" class="text-red-400 text-xs mt-1">{{ task.error }}</div>
          <div v-if="task.speed" class="text-gray-500 text-xs mt-1">{{ task.speed }}</div>

          <!-- 子任务列表 (大文件分块上传) -->
          <div v-if="task.subtasks && task.subtasks.length > 1" class="mt-2 border-t border-gray-700 pt-2">
            <n-button size="tiny" text @click="toggleSubtasks(task.id)">
              {{ showSubtaskMap[task.id] ? '收起子任务' : `展开子任务 (${task.subtasks.length})` }}
            </n-button>
            <div v-if="showSubtaskMap[task.id]" class="mt-2 space-y-1">
              <div
                v-for="(sub, idx) in task.subtasks"
                :key="idx"
                class="flex items-center justify-between text-xs"
              >
                <span class="text-gray-400">{{ sub.name || `分块 ${idx + 1}` }}</span>
                <div class="flex items-center gap-2">
                  <span class="text-gray-500">{{ formatSize(sub.size || 0) }}</span>
                  <n-tag v-if="sub.status === 'completed'" type="success" size="small">完成</n-tag>
                  <n-tag v-else-if="sub.status === 'error'" type="error" size="small">失败</n-tag>
                  <n-tag v-else-if="sub.status === 'uploading'" type="warning" size="small">上传中</n-tag>
                  <n-tag v-else type="default" size="small">等待中</n-tag>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useUploadTaskStore } from '../stores/uploadTask'
import { useMessage } from 'naive-ui'
import AppHeader from '../components/AppHeader.vue'
import { Document, Close } from '@vicons/ionicons5'

const uploadStore = useUploadTaskStore()
const message = useMessage()
const showSubtaskMap = ref({})

function formatSize(bytes) {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

function toggleTask(task) {
  if (task.status === 'uploading') {
    uploadStore.pauseTask(task.id)
  } else {
    uploadStore.resumeTask(task.id)
  }
}

function toggleSubtasks(taskId) {
  showSubtaskMap.value[taskId] = !showSubtaskMap.value[taskId]
}
</script>
