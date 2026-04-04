<template>
  <div class="min-h-screen bg-[#0f0f0f]">
    <AppHeader />

    <main class="p-6">
      <div class="flex items-center justify-between mb-4">
        <h2 class="text-white text-xl">同步资源</h2>
        <div class="flex gap-2">
          <n-button v-if="syncStore.isRunning" type="info" :loading="true">同步中...</n-button>
          <n-button v-else type="primary" @click="showAddModal = true">添加任务</n-button>
          <n-button v-if="syncStore.completedCount > 0" @click="syncStore.clearCompleted">清除完成</n-button>
        </div>
      </div>

      <!-- 统计栏 -->
      <div class="flex gap-4 mb-6">
        <n-tag :bordered="false" type="info">等待中: {{ syncStore.pendingCount }}</n-tag>
        <n-tag :bordered="false" type="warning">下载中: {{ syncStore.downloadingCount }}</n-tag>
        <n-tag :bordered="false" type="success">上传中: {{ syncStore.uploadingCount }}</n-tag>
        <n-tag :bordered="false" type="success">已完成: {{ syncStore.completedCount }}</n-tag>
        <n-tag :bordered="false" type="error">失败: {{ syncStore.errorCount }}</n-tag>
      </div>

      <!-- 任务列表 -->
      <div v-if="syncStore.syncTasks.length === 0" class="text-center py-12">
        <n-empty description="暂无同步任务" />
      </div>

      <div v-else class="space-y-3">
        <div
          v-for="task in syncStore.syncTasks"
          :key="task.id"
          class="bg-[#1a1a1a] p-4"
        >
          <div class="flex items-center justify-between mb-2">
            <div class="flex items-center gap-2 flex-1 min-w-0">
              <n-icon size="20" color="#60a5fa"><CloudDownload /></n-icon>
              <span class="text-white truncate">{{ task.name }}</span>
              <n-tag size="small" :type="task.type === 'lanzou' ? 'warning' : 'default'">
                {{ task.type === 'lanzou' ? '蓝奏云' : '直链' }}
              </n-tag>
            </div>
            <div class="flex items-center gap-2">
              <n-tag v-if="task.status === 'completed'" type="success" size="small">完成</n-tag>
              <n-tag v-else-if="task.status === 'error'" type="error" size="small">失败</n-tag>
              <n-tag v-else-if="task.status === 'downloading'" type="warning" size="small">下载中</n-tag>
              <n-tag v-else-if="task.status === 'uploading'" type="info" size="small">上传中</n-tag>
              <n-tag v-else type="default" size="small">等待中</n-tag>
              <n-button v-if="task.status === 'error'" size="tiny" @click="syncStore.retryTask(task.id)">重试</n-button>
              <n-button
                v-if="task.status !== 'downloading' && task.status !== 'uploading'"
                size="tiny"
                type="error"
                text
                @click="syncStore.removeTask(task.id)"
              >
                <n-icon><Close /></n-icon>
              </n-button>
            </div>
          </div>

          <!-- 进度条 -->
          <n-progress
            v-if="task.status !== 'completed'"
            :percentage="task.progress"
            :show-indicator="false"
            :height="4"
          />
          <div v-if="task.status === 'error'" class="text-red-400 text-xs mt-1">{{ task.error }}</div>
        </div>
      </div>
    </main>

    <!-- 添加任务模态框 -->
    <n-modal v-model:show="showAddModal">
      <n-card title="添加同步任务" style="width: 500px;">
        <div class="space-y-4">
          <n-alert type="info" :show-icon="false">
            从外部链接下载文件并自动上传到蓝奏云。
          </n-alert>

          <div>
            <div class="text-gray-400 text-sm mb-2">链接类型</div>
            <n-radio-group v-model:value="taskType">
              <n-space>
                <n-radio value="direct">直链下载</n-radio>
                <n-radio value="lanzou" disabled>蓝奏云分享 (待开发)</n-radio>
              </n-space>
            </n-radio-group>
          </div>

          <div>
            <div class="text-gray-400 text-sm mb-2">文件链接</div>
            <n-input
              v-model:value="taskUrl"
              placeholder="https://example.com/file.zip"
            />
          </div>

          <div>
            <div class="text-gray-400 text-sm mb-2">文件名称 (可选)</div>
            <n-input
              v-model:value="taskName"
              placeholder="留空则从链接自动提取"
            />
          </div>
        </div>

        <template #footer>
          <div class="flex justify-end gap-2">
            <n-button @click="showAddModal = false">取消</n-button>
            <n-button type="primary" @click="addSyncTask" :disabled="!taskUrl">添加</n-button>
          </div>
        </template>
      </n-card>
    </n-modal>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useSyncStore } from '../stores/sync'
import { useMessage } from 'naive-ui'
import AppHeader from '../components/AppHeader.vue'
import { CloudDownload, Close } from '@vicons/ionicons5'

const syncStore = useSyncStore()
const message = useMessage()

const showAddModal = ref(false)
const taskType = ref('direct')
const taskUrl = ref('')
const taskName = ref('')

function addSyncTask() {
  if (!taskUrl.value) return

  syncStore.addTask({
    url: taskUrl.value,
    name: taskName.value || undefined,
    type: taskType.value
  })

  message.success('同步任务已添加')
  showAddModal.value = false
  taskUrl.value = ''
  taskName.value = ''
}
</script>
