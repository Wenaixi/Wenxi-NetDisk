<template>
  <div class="min-h-screen bg-[#0f0f0f]">
    <AppHeader />

    <main class="p-6">
      <div class="flex items-center justify-between mb-6">
        <div class="flex items-center gap-2">
          <n-button @click="navigateBack" :disabled="breadcrumbs.length <= 1">
            <template #icon><n-icon><ArrowBack /></n-icon></template>
            返回
          </n-button>
          <n-breadcrumb>
            <n-breadcrumb-item v-for="(crumb, idx) in breadcrumbs" :key="crumb.id">
              <n-button text @click="navigateToBreadcrumb(idx)">{{ crumb.name }}</n-button>
            </n-breadcrumb-item>
          </n-breadcrumb>
        </div>
        <div class="flex items-center gap-2">
          <n-button @click="refresh" :loading="loading">
            <template #icon><n-icon><Refresh /></n-icon></template>
          </n-button>
          <n-checkbox
            v-if="files.length > 0 || folders.length > 0"
            v-model:checked="selectMode"
            label="选择模式"
          />
          <n-button
            v-if="selectedItems.length > 0"
            type="error"
            @click="batchDelete"
          >
            删除选中 ({{ selectedItems.length }})
          </n-button>
        </div>
      </div>

      <n-spin :show="loading">
        <div v-if="!connected" class="text-center py-12">
          <n-result status="warning" title="未连接蓝奏云">
            <template #footer>
              <n-button type="primary" @click="$router.push('/settings')">去设置</n-button>
            </template>
          </n-result>
        </div>

        <div v-else-if="folders.length === 0 && files.length === 0" class="text-center py-12">
          <n-empty description="文件夹为空" />
        </div>

        <div v-else>
          <div class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4">
            <div
              v-for="folder in folders"
              :key="folder.folder_id"
              :class="[
                'bg-[#1a1a1a] p-4 cursor-pointer hover:bg-[#252525]',
                selectMode && selectedItems.includes(folder.folder_id) ? 'ring-2 ring-blue-500' : ''
              ]"
              @dblclick="!selectMode && navigateToFolder(folder)"
              @click="selectMode ? toggleSelect(folder.folder_id) : navigateToFolder(folder)"
            >
              <div class="flex items-center gap-2">
                <n-checkbox
                  v-if="selectMode"
                  :checked="selectedItems.includes(folder.folder_id)"
                  @click.stop
                  @update:checked="toggleSelect(folder.folder_id)"
                />
                <n-icon size="24" color="#60a5fa"><Folder /></n-icon>
                <span class="text-white truncate">{{ folder.name }}</span>
              </div>
            </div>

            <div
              v-for="file in files"
              :key="file.file_id"
              :class="[
                'bg-[#1a1a1a] p-4 cursor-pointer hover:bg-[#252525]',
                selectMode && selectedItems.includes(file.file_id) ? 'ring-2 ring-blue-500' : ''
              ]"
              @click="selectMode ? toggleSelect(file.file_id) : handleFileClick(file)"
            >
              <div class="flex items-center gap-2 mb-2">
                <n-checkbox
                  v-if="selectMode"
                  :checked="selectedItems.includes(file.file_id)"
                  @click.stop
                  @update:checked="toggleSelect(file.file_id)"
                />
                <n-icon size="24" color="#a78bfa"><Document /></n-icon>
                <span class="text-white truncate">{{ file.name }}</span>
              </div>
              <div class="text-gray-500 text-sm">{{ formatSize(file.size) }}</div>
            </div>
          </div>
        </div>
      </n-spin>
    </main>

    <n-modal v-model:show="showFileMenu" v-if="selectedFile">
      <n-card :title="selectedFile.name" style="width: 400px;">
        <div class="space-y-2">
          <n-button block type="primary" @click="downloadFile">下载</n-button>
          <n-button block @click="createShare">创建分享链接</n-button>
          <n-button block @click="showFileMenu = false">取消</n-button>
        </div>
      </n-card>
    </n-modal>

    <n-modal v-model:show="showShareModal" v-if="shareResult">
      <n-card title="分享链接" style="width: 400px;">
        <div class="space-y-4">
          <div class="bg-[#252525] p-3 rounded-none">
            <n-input :value="shareResult.url" readonly />
          </div>
          <div v-if="shareResult.pwd" class="text-gray-400">
            密码: {{ shareResult.pwd }}
          </div>
          <div class="text-gray-500 text-sm">
            链接: {{ shareResult.url }}
          </div>
        </div>
        <template #footer>
          <div class="flex justify-end gap-2">
            <n-button @click="copyShareLink">复制链接</n-button>
            <n-button type="primary" @click="showShareModal = false">关闭</n-button>
          </div>
        </template>
      </n-card>
    </n-modal>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { lanzouAPI } from '../api/lanzou'
import { useMessage } from 'naive-ui'
import AppHeader from '../components/AppHeader.vue'
import {
  ArrowBack, Refresh, Folder, Document
} from '@vicons/ionicons5'

const router = useRouter()
const authStore = useAuthStore()
const message = useMessage()

const loading = ref(false)
const connected = ref(false)
const folders = ref([])
const files = ref([])
const breadcrumbs = ref([{ id: -1, name: '首页' }])
const currentFolderId = ref(-1)

const showFileMenu = ref(false)
const selectedFile = ref(null)
const showShareModal = ref(false)
const shareResult = ref(null)

// 选择模式相关
const selectMode = ref(false)
const selectedItems = ref([])

function toggleSelect(id) {
  const idx = selectedItems.value.indexOf(id)
  if (idx === -1) {
    selectedItems.value.push(id)
  } else {
    selectedItems.value.splice(idx, 1)
  }
}

async function batchDelete() {
  if (selectedItems.value.length === 0) return
  if (!confirm(`确定删除选中的 ${selectedItems.value.length} 个项目吗?`)) return

  try {
    const ids = [...selectedItems.value]
    let deletedCount = 0
    for (const id of ids) {
      await lanzouAPI.delete(id)
      deletedCount++
    }
    message.success(`成功删除 ${deletedCount} 个项目`)
    selectedItems.value = []
    selectMode.value = false
    await fetchContents(currentFolderId.value)
  } catch (err) {
    message.error(err.message || '删除失败')
  }
}

async function fetchStatus() {
  try {
    const res = await lanzouAPI.status()
    connected.value = res.connected
    return res.connected
  } catch {
    connected.value = false
    return false
  }
}

async function fetchContents(folderId = -1) {
  loading.value = true
  try {
    const res = await lanzouAPI.listFiles({ folder_id: folderId, page: 1 })
    folders.value = res.folders || []
    files.value = res.files || []
    currentFolderId.value = folderId
  } catch (err) {
    message.error('获取文件列表失败')
    folders.value = []
    files.value = []
  } finally {
    loading.value = false
  }
}

function navigateToFolder(folder) {
  breadcrumbs.value.push({ id: folder.folder_id, name: folder.name })
  fetchContents(folder.folder_id)
}

function navigateBack() {
  if (breadcrumbs.value.length <= 1) return
  breadcrumbs.value.pop()
  const prev = breadcrumbs.value[breadcrumbs.value.length - 1]
  fetchContents(prev.id)
}

function navigateToBreadcrumb(idx) {
  if (idx === breadcrumbs.value.length - 1) return
  breadcrumbs.value = breadcrumbs.value.slice(0, idx + 1)
  const crumb = breadcrumbs.value[idx]
  fetchContents(crumb.id)
}

function handleFileClick(file) {
  selectedFile.value = file
  showFileMenu.value = true
}

async function downloadFile() {
  showFileMenu.value = false
  try {
    const res = await lanzouAPI.getDownloadUrl(selectedFile.value.file_id)
    if (res.url) {
      const a = document.createElement('a')
      a.href = res.url
      a.download = selectedFile.value.name
      document.body.appendChild(a)
      a.click()
      document.body.removeChild(a)
      message.success('开始下载')
    } else {
      message.warning('获取下载链接失败')
    }
  } catch (err) {
    message.error(err.message || '下载失败')
  }
}

async function createShare() {
  showFileMenu.value = false
  try {
    const res = await lanzouAPI.createShare({
      file_id: selectedFile.value.file_id,
      minutes: 0
    })
    shareResult.value = res
    showShareModal.value = true
  } catch (err) {
    message.error(err.message || '创建分享失败')
  }
}

async function copyShareLink() {
  try {
    await navigator.clipboard.writeText(shareResult.value.url)
    message.success('链接已复制')
  } catch {
    message.error('复制失败')
  }
}

async function refresh() {
  await fetchContents(currentFolderId.value)
}

function formatSize(bytes) {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

onMounted(async () => {
  await authStore.fetchUser()
  const isConnected = await fetchStatus()
  if (isConnected) {
    await fetchContents()
  }
})
</script>
