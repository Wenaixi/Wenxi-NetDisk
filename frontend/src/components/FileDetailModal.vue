<template>
  <n-modal v-model:show="show" :mask-closable="true">
    <n-card
      :title="item?.name || '文件详情'"
      style="width: 500px;"
    >
      <div class="space-y-4">
        <!-- 图标和名称 -->
        <div class="flex items-center gap-3">
          <n-icon :size="32" :color="isFile ? '#a78bfa' : '#60a5fa'">
            <Document v-if="isFile" />
            <Folder v-else />
          </n-icon>
          <div class="flex-1 min-w-0">
            <div class="text-white font-bold truncate">{{ item?.name }}</div>
            <n-tag size="small" :type="isFile ? 'success' : 'info'">
              {{ isFile ? '文件' : '文件夹' }}
            </n-tag>
          </div>
        </div>

        <!-- 详情列表 -->
        <div class="divide-y divide-gray-700">
          <div v-if="isFile" class="py-2 flex justify-between">
            <span class="text-gray-400">文件大小</span>
            <span class="text-white">{{ formatSize(item?.size) }}</span>
          </div>
          <div v-if="isFile" class="py-2 flex justify-between">
            <span class="text-gray-400">文件类型</span>
            <span class="text-white">{{ item?.mime_type || '未知' }}</span>
          </div>
          <div class="py-2 flex justify-between">
            <span class="text-gray-400">创建时间</span>
            <span class="text-white">{{ formatDate(item?.created_at) }}</span>
          </div>
          <div class="py-2 flex justify-between">
            <span class="text-gray-400">更新时间</span>
            <span class="text-white">{{ formatDate(item?.updated_at) }}</span>
          </div>
          <div class="py-2 flex justify-between">
            <span class="text-gray-400">ID</span>
            <span class="text-white text-sm">{{ item?.id }}</span>
          </div>
        </div>

        <!-- 描述编辑 -->
        <div>
          <div class="text-gray-400 text-sm mb-2">描述</div>
          <n-input
            v-model:value="description"
            type="textarea"
            placeholder="添加描述..."
            :autosize="{ minRows: 2, maxRows: 4 }"
            maxlength="500"
            show-count
          />
          <n-button
            size="small"
            type="primary"
            class="mt-2"
            :loading="saving"
            @click="saveDescription"
            :disabled="!description && !item?.description"
          >
            {{ item?.description ? '更新描述' : '保存描述' }}
          </n-button>
        </div>

        <!-- 蓝奏云信息 (仅文件) -->
        <div v-if="isFile && item?.lanzou_file_id" class="text-sm">
          <div class="text-gray-400 mb-1">蓝奏云信息</div>
          <div class="bg-[#1a1a1a] p-3 space-y-1">
            <div class="flex justify-between">
              <span class="text-gray-500">蓝奏云文件ID</span>
              <span class="text-gray-300 text-xs">{{ item.lanzou_file_id }}</span>
            </div>
          </div>
        </div>
      </div>

      <template #footer>
        <div class="flex justify-end">
          <n-button @click="show = false">关闭</n-button>
        </div>
      </template>
    </n-card>
  </n-modal>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { Document, Folder } from '@vicons/ionicons5'

const props = defineProps({
  item: { type: Object, default: null }
})

const emit = defineEmits(['update'])

const show = ref(false)
const description = ref('')
const saving = ref(false)

const isFile = computed(() => props.item?.size !== undefined)

function formatDate(dateStr) {
  if (!dateStr) return '未知'
  const d = new Date(dateStr)
  return d.toLocaleString('zh-CN', {
    year: 'numeric', month: '2-digit', day: '2-digit',
    hour: '2-digit', minute: '2-digit'
  })
}

function formatSize(bytes) {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

async function saveDescription() {
  saving.value = true
  try {
    emit('update', { id: props.item.id, description: description.value })
  } finally {
    saving.value = false
  }
}

watch(() => props.item, (val) => {
  if (val) {
    description.value = val.description || ''
    show.value = true
  }
}, { immediate: true })
</script>
