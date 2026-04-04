<template>
  <div class="min-h-screen bg-[#0f0f0f]">
    <AppHeader />

    <main class="p-6">
      <div class="flex items-center justify-between mb-4">
        <div class="flex items-center gap-2">
          <n-button @click="recycleStore.fetchList" :loading="recycleStore.loading">
            <template #icon><n-icon><Refresh /></n-icon></template>
            刷新
          </n-button>
          <n-button
            v-if="recycleStore.hasSelectedItems"
            type="success"
            @click="batchRestore"
          >
            恢复选中 ({{ recycleStore.selectedItems.length }})
          </n-button>
          <n-button
            v-if="recycleStore.hasSelectedItems"
            type="error"
            @click="batchDelete"
          >
            彻底删除 ({{ recycleStore.selectedItems.length }})
          </n-button>
          <n-button
            v-if="recycleStore.totalCount > 0"
            type="error"
            strong
            @click="clearAll"
          >
            清空回收站
          </n-button>
        </div>
        <div class="text-gray-400 text-sm">
          共 {{ recycleStore.totalCount }} 项 (文件: {{ recycleStore.fileCount }}, 文件夹: {{ recycleStore.folderCount }})
        </div>
      </div>

      <n-spin :show="recycleStore.loading">
        <n-empty
          v-if="recycleStore.items.length === 0"
          description="回收站为空"
          class="py-12"
        />

        <div v-else class="space-y-2">
          <div
            v-for="item in recycleStore.items"
            :key="item.id"
            :class="[
              'bg-[#1a1a1a] p-4 flex items-center justify-between',
              recycleStore.selectMode && recycleStore.selectedItems.includes(item.id) ? 'ring-2 ring-blue-500' : ''
            ]"
          >
            <div class="flex items-center gap-3">
              <n-checkbox
                v-if="recycleStore.selectMode"
                :checked="recycleStore.selectedItems.includes(item.id)"
                @update:checked="recycleStore.toggleSelect(item.id)"
              />
              <n-icon size="24" :component="item.item_type === 'file' ? DocumentText : Folder" />
              <div>
                <div class="text-white text-sm">{{ item.original_name }}</div>
                <div class="text-gray-500 text-xs">
                  {{ recycleStore.formatSize(item.size) }} · 删除于 {{ recycleStore.formatDate(item.deleted_at) }}
                </div>
              </div>
            </div>
            <div class="flex items-center gap-2">
              <n-button size="small" type="success" @click="restore(item.id)">
                恢复
              </n-button>
              <n-popconfirm @positive-click="permanentDelete(item.id)">
                <template #trigger>
                  <n-button size="small" type="error">
                    彻底删除
                  </n-button>
                </template>
                确定要彻底删除吗？此操作不可恢复。
              </n-popconfirm>
            </div>
          </div>
        </div>
      </n-spin>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { DocumentText, Folder, Refresh } from '@vicons/ionicons5'
import { useMessage } from 'naive-ui'
import AppHeader from '../components/AppHeader.vue'
import { useRecycleStore } from '../stores/recycle'

const message = useMessage()
const recycleStore = useRecycleStore()

onMounted(() => {
  recycleStore.fetchList()
})

async function restore(id) {
  try {
    await recycleStore.restore(id)
    message.success('已恢复')
  } catch (err) {
    message.error(err.message || '恢复失败')
  }
}

async function permanentDelete(id) {
  try {
    await recycleStore.permanentDelete(id)
    message.success('已彻底删除')
  } catch (err) {
    message.error(err.message || '删除失败')
  }
}

async function batchRestore() {
  if (recycleStore.selectedItems.length === 0) return
  try {
    await recycleStore.batchRestore()
    message.success('已批量恢复')
  } catch (err) {
    message.error(err.message || '批量恢复失败')
  }
}

async function batchDelete() {
  if (recycleStore.selectedItems.length === 0) return
  try {
    await recycleStore.batchDelete()
    message.success('已批量彻底删除')
  } catch (err) {
    message.error(err.message || '批量删除失败')
  }
}

async function clearAll() {
  try {
    await recycleStore.clearAll()
    message.success('回收站已清空')
  } catch (err) {
    message.error(err.message || '清空失败')
  }
}
</script>
