<template>
  <div class="min-h-screen bg-[#0f0f0f]">
    <AppHeader />

    <main class="p-6 max-w-5xl mx-auto">
      <div class="mb-6">
        <h2 class="text-xl text-white mb-2">链接解析</h2>
        <p class="text-gray-500 text-sm">解析蓝奏云分享链接，批量获取文件信息并下载</p>
      </div>

      <!-- 输入区域 -->
      <div class="bg-[#1a1a1a] p-4 mb-4">
        <div class="flex items-start gap-2">
          <div class="flex-1">
            <n-input
              v-model:value="shareUrl"
              type="textarea"
              placeholder="* https://... 可同时解析多行&#10;每行一个分享链接"
              :autosize="{ minRows: 2, maxRows: 6 }"
              class="mb-2"
              @keyup.enter.ctrl="parseLinks"
            />
            <n-input
              v-model:value="sharePwd"
              placeholder="提取密码（选填）"
              style="width: 150px"
              maxlength="6"
            />
          </div>
          <n-button
            type="primary"
            :loading="loading"
            @click="parseLinks"
            style="height: 100%"
          >
            {{ !shareUrl ? '粘贴并解析' : '解析' }}
          </n-button>
        </div>
      </div>

      <!-- 错误提示 -->
      <n-alert v-if="error" type="error" class="mb-4" closable @close="error = ''">
        {{ error }}
      </n-alert>

      <!-- 解析结果 -->
      <n-spin :show="loading">
        <n-empty v-if="shareFiles.length === 0 && !currentShare" description="暂无解析结果" class="py-12" />

        <template v-else>
          <!-- 当前分享信息 -->
          <div v-if="currentShare" class="bg-[#1a1a1a] p-3 mb-2 flex items-center justify-between">
            <div class="flex items-center gap-2">
              <n-icon size="20" :component="currentShare.type === 'folder' ? Folder : DocumentText" />
              <span class="text-white">{{ currentShare.name }}</span>
              <span v-if="currentShare.size" class="text-gray-500 text-sm">（{{ currentShare.size }}）</span>
            </div>
            <div class="flex items-center gap-2">
              <n-checkbox v-model:checked="autoMerge" :disabled="selectedRows.length > 0">
                自动合并
              </n-checkbox>
            </div>
          </div>

          <!-- 操作栏 -->
          <div class="bg-[#1a1a1a] p-3 mb-2 flex items-center justify-between">
            <span class="text-gray-400 text-sm">
              共 {{ totalFiles }} 个文件
              <span v-if="selectedRows.length > 0">，已选 {{ selectedRows.length }} 项</span>
            </span>
            <div class="flex gap-2">
              <n-button
                v-if="selectedRows.length > 0"
                type="primary"
                size="small"
                @click="downloadSelected"
              >
                下载 ({{ selectedRows.length }}项)
              </n-button>
              <n-button
                v-else
                size="small"
                @click="downloadAll"
                :disabled="totalFiles === 0"
              >
                下载全部
              </n-button>
              <n-button size="small" @click="clearResults">清除结果</n-button>
            </div>
          </div>

          <!-- 文件列表 -->
          <div class="bg-[#1a1a1a]">
            <!-- 表头 -->
            <div class="flex items-center px-4 py-2 text-gray-400 text-sm border-b border-gray-700">
              <div class="w-8"></div>
              <div class="flex-1">文件名</div>
              <div class="w-32">时间</div>
              <div class="w-24">大小</div>
              <div class="w-16">操作</div>
            </div>

            <!-- 文件行 -->
            <div
              v-for="(item, idx) in allFiles"
              :key="idx"
              :class="[
                'flex items-center px-4 py-2 cursor-pointer hover:bg-[#252525]',
                isSelected(item) ? 'bg-[#1e3a5f]' : ''
              ]"
              @click="toggleSelect(item)"
            >
              <div class="w-8">
                <n-icon size="16" :component="DocumentText" color="#a78bfa" />
              </div>
              <div class="flex-1 text-white text-sm truncate">{{ item.name }}</div>
              <div class="w-32 text-gray-500 text-sm">{{ item.time || '-' }}</div>
              <div class="w-24 text-gray-500 text-sm">{{ item.size || '-' }}</div>
              <div class="w-16">
                <n-button size="tiny" text @click.stop="downloadSingle(item)">
                  <n-icon><CloudDownload /></n-icon>
                </n-button>
              </div>
            </div>
          </div>
        </template>
      </n-spin>

      <!-- 使用说明 -->
      <n-card class="mt-6 bg-[#1a1a1a]" :bordered="false">
        <template #header>
          <span class="text-white">使用说明</span>
        </template>
        <div class="text-gray-400 text-sm space-y-1">
          <p>• 支持解析蓝奏云分享链接（lanzous/lanzoui/lanzoux/lanzouv/lanzouo）</p>
          <p>• 支持多行批量解析，每行一个分享链接</p>
          <p>• 如分享有密码，请先填入密码再解析</p>
          <p>• 可点击文件行进行多选，支持批量下载到下载任务列表</p>
          <p>• 自动合并：解析文件夹时自动勾选合并选项</p>
        </div>
      </n-card>
    </main>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { DocumentText, Folder, CloudDownload } from '@vicons/ionicons5'
import { useMessage } from 'naive-ui'
import { useDownloadTaskStore } from '../stores/downloadTask'
import AppHeader from '../components/AppHeader.vue'
import api from '../api'

const message = useMessage()
const downloadStore = useDownloadTaskStore()

const shareUrl = ref('')
const sharePwd = ref('')
const loading = ref(false)
const error = ref('')
const currentShare = ref(null)
const shareFiles = ref([]) // [{ name, size, type, list: [{url, name, size, time, pwd}] }]
const selectedRows = ref([])
const autoMerge = ref(false)

// 扁平化的文件列表
const allFiles = computed(() => {
  return shareFiles.value.flatMap(share =>
    (share.list || []).map(item => ({
      ...item,
      shareName: share.name,
      shareType: share.type
    }))
  )
})

const totalFiles = computed(() => {
  return shareFiles.value.reduce((sum, s) => sum + (s.list?.length || 0), 0)
})

onMounted(() => {
  // 检查剪贴板是否有链接
  navigator.clipboard.readText().then(text => {
    if (/https?:\/\//.test(text)) {
      shareUrl.value = text
    }
  }).catch(() => {})
})

async function parseLinks() {
  if (!shareUrl.value.trim()) {
    try {
      const text = await navigator.clipboard.readText()
      if (/https?:\/\//.test(text)) {
        shareUrl.value = text
      } else {
        message.info('请先粘贴蓝奏云分享链接')
        return
      }
    } catch {
      message.info('请先输入蓝奏云分享链接')
      return
    }
  }

  loading.value = true
  error.value = ''
  shareFiles.value = []
  currentShare.value = null
  selectedRows.value = []

  try {
    const lines = shareUrl.value.trim().split('\n').filter(Boolean)

    // 如果有密码且只有一行，密码应用到所有行
    const pwd = sharePwd.value

    for (const line of lines) {
      const url = line.trim()
      if (!url.startsWith('http')) continue

      try {
        const res = await api.post('/lanzou/share/parse', { url, pwd })
        if (res.data) {
          shareFiles.value.push(res.data)
        }
      } catch (err) {
        error.value += `${url}: ${err.message}\n`
      }
    }

    // 如果只有一个分享，设置当前分享
    if (shareFiles.value.length === 1) {
      const item = shareFiles.value[0]
      currentShare.value = item
      // 如果是文件夹且是文件分割格式，默认勾选自动合并
      if (item.type === 'folder' && isSplitFileName(item.name)) {
        autoMerge.value = true
      }
    }

    if (shareFiles.value.length > 0) {
      message.success(`解析成功，共 ${totalFiles.value} 个文件`)
    } else if (!error.value) {
      message.warning('未解析到有效文件')
    }
  } catch (err) {
    error.value = err.message || '解析失败'
  } finally {
    loading.value = false
  }
}

function isSplitFileName(name) {
  // 检测是否为分割文件名 (如 file.part001of3.zip)
  return /\.part\d+of\d+/i.test(name)
}

function isSelected(item) {
  return selectedRows.value.some(r => r.url === item.url)
}

function toggleSelect(item) {
  const idx = selectedRows.value.findIndex(r => r.url === item.url)
  if (idx === -1) {
    selectedRows.value.push(item)
  } else {
    selectedRows.value.splice(idx, 1)
  }
}

function downloadSingle(item) {
  downloadStore.addTask({
    name: item.name,
    url: item.url,
    pwd: item.pwd || sharePwd.value,
    size: 0,
    merge: false
  })
  message.success('已添加到下载列表')
}

function downloadAll() {
  const files = allFiles.value
  if (files.length === 0) return

  files.forEach(item => {
    downloadStore.addTask({
      name: item.name,
      url: item.url,
      pwd: item.pwd || sharePwd.value,
      size: 0,
      merge: autoMerge.value
    })
  })

  message.success(`已添加 ${files.length} 个下载任务`)
}

function downloadSelected() {
  if (selectedRows.value.length === 0) return

  selectedRows.value.forEach(item => {
    downloadStore.addTask({
      name: item.name,
      url: item.url,
      pwd: item.pwd || sharePwd.value,
      size: 0,
      merge: false
    })
  })

  message.success(`已添加 ${selectedRows.value.length} 个下载任务`)
  selectedRows.value = []
}

function clearResults() {
  shareFiles.value = []
  currentShare.value = null
  selectedRows.value = []
  error.value = ''
}
</script>
