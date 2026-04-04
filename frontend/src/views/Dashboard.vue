<template>
  <div class="min-h-screen bg-[#0f0f0f]">
    <AppHeader />

    <main class="p-6">
      <div class="flex items-center justify-between mb-4">
        <div class="flex items-center gap-2">
          <n-button @click="fileStore.navigateBack()" :disabled="fileStore.breadcrumbs.length <= 1">
            <template #icon><n-icon><ArrowBack /></n-icon></template>
            返回
          </n-button>
          <n-button @click="showUploadModal = true">
            <template #icon><n-icon><CloudUpload /></n-icon></template>
            加密上传
          </n-button>
          <n-button @click="showNewFolderModal = true">
            <template #icon><n-icon><Create /></n-icon></template>
            新建文件夹
          </n-button>
          <n-button
            v-if="selectedItems.length > 0"
            type="error"
            @click="batchDelete"
          >
            删除选中 ({{ selectedItems.length }})
          </n-button>
        </div>
        <div class="flex items-center gap-2">
          <n-checkbox
            v-if="fileStore.files.length > 0 || fileStore.folders.length > 0"
            v-model:checked="selectMode"
            label="选择模式"
          />
          <n-button @click="refresh">
            <template #icon><n-icon><Refresh /></n-icon></template>
          </n-button>
        </div>
      </div>

      <div class="flex items-center gap-1 mb-4 text-sm">
        <template v-for="(crumb, idx) in fileStore.breadcrumbs" :key="crumb.id">
          <span v-if="idx > 0" class="text-gray-500 mx-1">/</span>
          <n-button
            text
            :type="idx === fileStore.breadcrumbs.length - 1 ? 'primary' : 'default'"
            @click="navigateToBreadcrumb(idx)"
          >
            {{ crumb.name }}
          </n-button>
        </template>
      </div>

      <n-spin :show="fileStore.loading">
        <div v-if="fileStore.folders.length === 0 && fileStore.files.length === 0" class="text-center py-12">
          <n-empty description="文件夹为空" />
        </div>

        <div v-else>
          <!-- 文件夹区域 -->
          <div v-if="fileStore.folders.length > 0" class="mb-6">
            <div class="text-gray-400 text-sm mb-3">文件夹</div>
            <div class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4">
              <div
                v-for="folder in fileStore.folders"
                :key="folder.id"
                :class="[
                  'bg-[#1a1a1a] p-4 cursor-pointer hover:bg-[#252525]',
                  selectMode && selectedItems.includes(folder.id) ? 'ring-2 ring-blue-500' : ''
                ]"
                @dblclick="!selectMode && fileStore.navigateToFolder(folder)"
                @click="selectMode ? toggleSelect(folder.id) : openFolderMenu(folder)"
                @contextmenu.prevent="!selectMode && openFolderMenu(folder)"
              >
                <div class="flex items-center gap-2">
                  <n-checkbox
                    v-if="selectMode"
                    :checked="selectedItems.includes(folder.id)"
                    @click.stop
                    @update:checked="toggleSelect(folder.id)"
                  />
                  <n-icon size="24" color="#60a5fa"><Folder /></n-icon>
                  <span class="text-white truncate">{{ folder.name }}</span>
                </div>
              </div>
            </div>
          </div>

          <!-- 文件区域 -->
          <div v-if="fileStore.files.length > 0">
            <div class="text-gray-400 text-sm mb-3">文件</div>
            <div class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4">
              <div
                v-for="file in fileStore.files"
                :key="file.id"
                :class="[
                  'bg-[#1a1a1a] p-4 cursor-pointer hover:bg-[#252525]',
                  selectMode && selectedItems.includes(file.id) ? 'ring-2 ring-blue-500' : ''
                ]"
                @click="selectMode ? toggleSelect(file.id) : handleFileClick(file)"
              >
                <div class="flex items-center gap-2 mb-2">
                  <n-checkbox
                    v-if="selectMode"
                    :checked="selectedItems.includes(file.id)"
                    @click.stop
                    @update:checked="toggleSelect(file.id)"
                  />
                  <n-icon size="24" color="#a78bfa"><Document /></n-icon>
                  <span class="text-white truncate">{{ file.name }}</span>
                </div>
                <div class="text-gray-500 text-sm">{{ formatSize(file.size) }}</div>
              </div>
            </div>
          </div>
        </div>
      </n-spin>
    </main>

    <n-modal v-model:show="showNewFolderModal">
      <n-card title="新建文件夹" style="width: 400px;">
        <n-input v-model:value="newFolderName" placeholder="文件夹名称" />
        <template #footer>
          <div class="flex justify-end gap-2">
            <n-button @click="showNewFolderModal = false">取消</n-button>
            <n-button type="primary" @click="createFolder">创建</n-button>
          </div>
        </template>
      </n-card>
    </n-modal>

    <n-modal v-model:show="showUploadModal" :mask-closable="false">
      <n-card title="加密上传文件" style="width: 500px;">
        <div class="space-y-4">
          <n-alert type="info" :show-icon="false">
            文件将在本地加密后上传，服务器和存储端均无法读取文件内容。
          </n-alert>

          <n-upload
            v-if="!selectedFile"
            :max="1"
            @change="handleFileChange"
            accept="*/*"
          >
            <n-button>选择文件</n-button>
          </n-upload>

          <div v-else class="bg-[#252525] p-4 rounded-none">
            <div class="flex items-center justify-between mb-2">
              <span class="text-white">{{ selectedFile.name }}</span>
              <n-button text @click="clearFile" :disabled="uploadStore.isUploading">
                <n-icon><Close /></n-icon>
              </n-button>
            </div>
            <div class="text-gray-400 text-sm">{{ formatSize(selectedFile.size) }}</div>
          </div>

          <div v-if="uploadStore.isEncrypting" class="space-y-2">
            <div class="flex items-center gap-2 text-blue-400">
              <n-spin size="small" />
              <span>正在加密...</span>
            </div>
          </div>

          <div v-if="uploadStore.isUploading" class="space-y-2">
            <div class="flex justify-between text-sm">
              <span class="text-gray-400">上传进度</span>
              <span class="text-white">{{ uploadStore.progress }}%</span>
            </div>
            <n-progress :percentage="uploadStore.progress" :show-indicator="false" />
            <div class="text-gray-500 text-xs">
              已上传: {{ formatSize(uploadStore.uploadedBytes) }} / {{ formatSize(uploadStore.totalBytes) }}
              <span v-if="uploadStore.speed > 0">({{ formatSpeed(uploadStore.speed) }})</span>
            </div>
          </div>

          <n-alert v-if="uploadStore.error" type="error" :show-icon="false">
            {{ uploadStore.error }}
          </n-alert>
        </div>

        <template #footer>
          <div class="flex justify-end gap-2">
            <n-button @click="closeUploadModal" :disabled="uploadStore.isUploading || uploadStore.isEncrypting">取消</n-button>
            <n-button
              type="primary"
              :loading="uploadStore.isUploading || uploadStore.isEncrypting"
              :disabled="!selectedFile || uploadStore.isUploading || uploadStore.isEncrypting"
              @click="startUpload"
            >
              {{ uploadStore.isUploading ? '上传中...' : (uploadStore.isEncrypting ? '加密中...' : '开始上传') }}
            </n-button>
          </div>
        </template>
      </n-card>
    </n-modal>

    <n-modal v-model:show="showFileMenu" v-if="selectedFileItem">
      <n-card :title="selectedFileItem.name" style="width: 400px;">
        <div class="space-y-2">
          <n-button v-if="selectedFileItem.size !== undefined" block @click="downloadFile(selectedFileItem)">下载</n-button>
          <n-button v-if="selectedFileItem.size !== undefined" block type="info" @click="openShareModal">分享</n-button>
          <n-button block @click="openRenameModal">重命名</n-button>
          <n-button block type="error" @click="deleteFile(selectedFileItem)">删除</n-button>
          <n-button block @click="showFileMenu = false">取消</n-button>
        </div>
      </n-card>
    </n-modal>

    <n-modal v-model:show="showRenameModal">
      <n-card title="重命名文件" style="width: 400px;">
        <n-input v-model:value="renameName" placeholder="新名称" />
        <template #footer>
          <div class="flex justify-end gap-2">
            <n-button @click="showRenameModal = false">取消</n-button>
            <n-button type="primary" @click="confirmRename">确定</n-button>
          </div>
        </template>
      </n-card>
    </n-modal>

    <n-modal v-model:show="showShareModal" :mask-closable="false">
      <n-card title="分享文件" style="width: 450px;">
        <div class="space-y-4">
          <n-alert v-if="shareResult" type="success" :show-icon="false">
            <div class="space-y-2">
              <div class="flex items-center gap-2">
                <span class="text-gray-400">分享链接:</span>
                <n-input :value="shareResult.url" readonly />
                <n-button size="small" type="primary" @click="copyShareLink">复制</n-button>
              </div>
              <div v-if="shareResult.pwd" class="text-gray-400">
                密码: <span class="text-white font-bold">{{ shareResult.pwd }}</span>
              </div>
              <div v-if="shareResult.expires_at" class="text-gray-400">
                有效期至: {{ shareResult.expires_at }}
              </div>
            </div>
          </n-alert>

          <div v-else class="space-y-4">
            <div>
              <div class="text-gray-400 text-sm mb-2">链接有效期</div>
              <n-select
                v-model:value="shareExpires"
                :options="expiresOptions"
              />
            </div>

            <div>
              <div class="text-gray-400 text-sm mb-2">访问密码 (可选)</div>
              <n-input
                v-model:value="sharePassword"
                placeholder="留空则无需密码"
                show-password-on="click"
              />
            </div>
          </div>

          <n-alert v-if="shareError" type="error" :show-icon="false">
            {{ shareError }}
          </n-alert>
        </div>

        <template #footer>
          <div class="flex justify-end gap-2">
            <n-button @click="closeShareModal">关闭</n-button>
            <n-button
              v-if="!shareResult"
              type="primary"
              :loading="isSharing"
              @click="createShare"
            >
              创建分享
            </n-button>
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
import { useFileStore } from '../stores/file'
import { useUploadStore } from '../stores/upload'
import { useShareStore } from '../stores/share'
import { encryptFile, generateEncryptionKey } from '../utils/crypto'
import { fileAPI, shareAPI } from '../api'
import { useMessage } from 'naive-ui'
import AppHeader from '../components/AppHeader.vue'
import {
  ArrowBack, CloudUpload, Refresh, Folder, Create, Document, Close
} from '@vicons/ionicons5'

const router = useRouter()
const authStore = useAuthStore()
const fileStore = useFileStore()
const uploadStore = useUploadStore()
const shareStore = useShareStore()
const message = useMessage()

const showNewFolderModal = ref(false)
const showUploadModal = ref(false)
const showFileMenu = ref(false)
const showShareModal = ref(false)
const newFolderName = ref('')
const selectedFile = ref(null)
const selectedFileItem = ref(null)

// 分享相关
const shareResult = ref(null)
const shareError = ref(null)
const sharePassword = ref('')
const shareExpires = ref('7d')
const isSharing = ref(false)
const expiresOptions = [
  { label: '1小时', value: '1h' },
  { label: '1天', value: '1d' },
  { label: '7天', value: '7d' },
  { label: '永久', value: 'never' }
]

// 重命名相关
const showRenameModal = ref(false)
const renameName = ref('')

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
    for (const id of ids) {
      // 判断是文件还是文件夹
      const isFile = fileStore.files.some(f => f.id === id)
      if (isFile) {
        await fileStore.deleteFile(id)
      } else {
        await fileStore.deleteFolder(id)
      }
    }
    message.success(`成功删除 ${ids.length} 个项目`)
    selectedItems.value = []
    selectMode.value = false
    fileStore.fetchFiles(fileStore.currentFolder)
  } catch (err) {
    message.error(err.message || '删除失败')
  }
}

function navigateToBreadcrumb(idx) {
  const crumb = fileStore.breadcrumbs[idx]
  if (crumb.id === fileStore.currentFolder) return
  fileStore.breadcrumbs = fileStore.breadcrumbs.slice(0, idx + 1)
  fileStore.fetchFiles(crumb.id)
}

async function createFolder() {
  if (!newFolderName.value.trim()) return
  await fileStore.createFolder(newFolderName.value, fileStore.currentFolder)
  newFolderName.value = ''
  showNewFolderModal.value = false
  message.success('文件夹创建成功')
}

function handleFileChange(options) {
  selectedFile.value = options.file.file
  uploadStore.reset()
}

function clearFile() {
  selectedFile.value = null
  uploadStore.reset()
}

function closeUploadModal() {
  if (uploadStore.isUploading || uploadStore.isEncrypting) return
  showUploadModal.value = false
  selectedFile.value = null
  uploadStore.reset()
}

async function startUpload() {
  if (!selectedFile.value) {
    message.warning('请选择文件')
    return
  }

  try {
    const key = await generateEncryptionKey()
    const encryptedBlob = await uploadStore.encryptAndUpload(
      selectedFile.value,
      key,
      fileStore.currentFolder
    )

    message.success('文件上传成功')
    fileStore.fetchFiles(fileStore.currentFolder)
    showUploadModal.value = false
    selectedFile.value = null
  } catch (err) {
    message.error(err.message || '上传失败')
  }
}

function handleFileClick(file) {
  selectedFileItem.value = file
  showFileMenu.value = true
}

async function downloadFile(file) {
  showFileMenu.value = false
  try {
    const response = await fileAPI.download(file.id)
    const url = response.download_url

    const a = document.createElement('a')
    a.href = url
    a.download = file.name
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)

    message.success('开始下载')
  } catch (err) {
    message.error(err.message || '下载失败')
  }
}

async function deleteFile(file) {
  showFileMenu.value = false
  try {
    await fileStore.deleteFile(file.id)
    message.success('删除成功')
  } catch (err) {
    message.error(err.message || '删除失败')
  }
}

async function refresh() {
  await fileStore.fetchFiles(fileStore.currentFolder)
  message.success('刷新成功')
}

function openShareModal() {
  showFileMenu.value = false
  shareResult.value = null
  shareError.value = null
  sharePassword.value = ''
  shareExpires.value = '7d'
  showShareModal.value = true
}

function openRenameModal() {
  showFileMenu.value = false
  renameName.value = selectedFileItem.value?.name || ''
  showRenameModal.value = true
}

async function confirmRename() {
  if (!renameName.value.trim() || !selectedFileItem.value) return
  try {
    await fileStore.renameFile(selectedFileItem.value.id, renameName.value)
    message.success('重命名成功')
    showRenameModal.value = false
    fileStore.fetchFiles(fileStore.currentFolder)
  } catch (err) {
    message.error(err.message || '重命名失败')
  }
}

function openFolderMenu(folder) {
  selectedFileItem.value = folder
  showFileMenu.value = true
}

function closeShareModal() {
  if (isSharing.value) return
  showShareModal.value = false
  shareResult.value = null
  shareError.value = null
}

async function createShare() {
  if (!selectedFileItem.value) return
  isSharing.value = true
  shareError.value = null

  try {
    const options = {}
    if (sharePassword.value) {
      options.password = sharePassword.value
    }
    if (shareExpires.value !== 'never') {
      const now = new Date()
      let expiresAt = new Date(now)
      switch (shareExpires.value) {
        case '1h':
          expiresAt.setHours(expiresAt.getHours() + 1)
          break
        case '1d':
          expiresAt.setDate(expiresAt.getDate() + 1)
          break
        case '7d':
          expiresAt.setDate(expiresAt.getDate() + 7)
          break
      }
      options.expiresAt = expiresAt.toISOString()
    }

    const res = await shareAPI.create({
      file_id: selectedFileItem.value.id,
      ...options
    })

    shareResult.value = {
      url: res.share_url,
      pwd: res.share_pwd || null,
      expires_at: res.expires_at || null
    }
  } catch (err) {
    shareError.value = err.message || '创建分享失败'
  } finally {
    isSharing.value = false
  }
}

async function copyShareLink() {
  if (!shareResult.value?.url) return
  try {
    await navigator.clipboard.writeText(window.location.origin + shareResult.value.url)
    message.success('链接已复制')
  } catch {
    message.error('复制失败')
  }
}

function formatSize(bytes) {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

function formatSpeed(bytesPerSecond) {
  return formatSize(bytesPerSecond) + '/s'
}

onMounted(async () => {
  await authStore.fetchUser()
  await fileStore.fetchFiles()
})
</script>
