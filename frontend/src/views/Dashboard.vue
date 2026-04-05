<template>
  <div class="min-h-screen bg-[#0f0f0f]"
       @dragover.prevent="isDragging = true"
       @dragleave="handleDragLeave"
       @drop.prevent="handleDrop">

    <!-- 拖拽上传遮罩 -->
    <div v-if="isDragging"
         class="fixed inset-0 z-50 bg-[rgba(0,0,0,0.8)] flex items-center justify-center pointer-events-none">
      <div class="border-4 border-dashed border-[#00D4FF] rounded-none p-12 text-center bg-[#1a1a1a]">
        <n-icon size="64" color="#00D4FF"><CloudUpload /></n-icon>
        <div class="text-white text-xl mt-4">拖拽文件到此处上传</div>
        <div class="text-gray-400 text-sm mt-2">上传到: {{ currentFolderName }}</div>
      </div>
    </div>

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
          <n-button
            v-if="selectedItems.length > 0"
            type="info"
            @click="openBatchMoveModal"
          >
            移动到 ({{ selectedItems.length }})
          </n-button>
          <n-button
            v-if="selectedItems.length > 0"
            type="warning"
            @click="batchDownload"
          >
            下载选中 ({{ selectedItems.length }})
          </n-button>
        </div>
        <div class="flex items-center gap-2">
          <n-checkbox
            v-if="fileStore.files.length > 0 || fileStore.folders.length > 0"
            v-model:checked="selectMode"
            label="选择模式"
          />
          <n-input
            v-if="fileStore.files.length > 0 || fileStore.folders.length > 0"
            v-model:value="searchQuery"
            placeholder="搜索..."
            clearable
            style="width: 150px"
          />
          <n-select
            v-if="fileStore.files.length > 0 || fileStore.folders.length > 0"
            v-model:value="sortBy"
            :options="sortOptions"
            style="width: 120px"
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
          <div v-if="filteredFolders.length > 0" class="mb-6">
            <div class="text-gray-400 text-sm mb-3">文件夹 ({{ filteredFolders.length }})</div>
            <div class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4">
              <div
                v-for="folder in filteredFolders"
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
          <div v-if="filteredFiles.length > 0">
            <div class="text-gray-400 text-sm mb-3">文件 ({{ filteredFiles.length }})</div>
            <div class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4">
              <div
                v-for="file in filteredFiles"
                :key="file.id"
                :class="[
                  'bg-[#1a1a1a] p-4 cursor-pointer hover:bg-[#252525]',
                  selectMode && selectedItems.includes(file.id) ? 'ring-2 ring-blue-500' : ''
                ]"
                @click="selectMode ? toggleSelect(file.id) : handleFileClick(file)"
                @contextmenu.prevent="!selectMode && openFileMenu(file)"
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
      <n-card title="加密上传文件" style="width: 600px;">
        <div class="space-y-4">
          <n-alert type="info" :show-icon="false">
            文件将在本地加密后上传，服务器和存储端均无法读取文件内容。
          </n-alert>

          <!-- 文件选择区 -->
          <n-upload
            v-if="uploadStore.uploadQueue.length === 0"
            :multiple="true"
            @change="handleFilesChange"
            accept="*/*"
            :max="10"
          >
            <n-button>选择文件 (最多10个)</n-button>
          </n-upload>

          <!-- 上传队列 -->
          <div v-else class="space-y-3">
            <div class="flex items-center justify-between">
              <span class="text-gray-400 text-sm">
                {{ uploadStore.queueUploadedCount }} / {{ uploadStore.queueTotalCount }} 已完成
              </span>
              <n-button size="small" text @click="uploadStore.clearQueue" :disabled="uploadStore.isQueueUploading">
                清空队列
              </n-button>
            </div>

            <div class="max-h-60 overflow-y-auto space-y-2">
              <div
                v-for="item in uploadStore.uploadQueue"
                :key="item.id"
                class="bg-[#1a1a1a] p-3"
              >
                <div class="flex items-center justify-between mb-2">
                  <div class="flex items-center gap-2 flex-1 min-w-0">
                    <n-icon size="16" color="#a78bfa"><Document /></n-icon>
                    <span class="text-white truncate">{{ item.file.name }}</span>
                  </div>
                  <div class="flex items-center gap-2">
                    <n-tag v-if="item.status === 'completed'" type="success" size="small">完成</n-tag>
                    <n-tag v-else-if="item.status === 'error'" type="error" size="small">失败</n-tag>
                    <n-spin v-else-if="item.status === 'encrypting'" size="small" />
                    <n-button
                      v-if="item.status !== 'completed' && item.status !== 'encrypting' && !uploadStore.isQueueUploading"
                      size="tiny"
                      text
                      @click="uploadStore.removeFromQueue(item.id)"
                    >
                      <n-icon><Close /></n-icon>
                    </n-button>
                  </div>
                </div>
                <div v-if="item.status === 'uploading'" class="text-gray-500 text-xs mb-1">
                  {{ item.progress }}%
                </div>
                <n-progress
                  v-if="item.status === 'uploading'"
                  :percentage="item.progress"
                  :show-indicator="false"
                  :height="4"
                />
                <div v-if="item.error" class="text-red-400 text-xs mt-1">{{ item.error }}</div>
              </div>
            </div>
          </div>
        </div>

        <template #footer>
          <div class="flex justify-end gap-2">
            <n-button @click="closeUploadModal" :disabled="uploadStore.isQueueUploading">取消</n-button>
            <n-button
              v-if="uploadStore.uploadQueue.length > 0 && !uploadStore.isQueueUploading && !uploadStore.hasQueueItems"
              type="primary"
              @click="startBatchUpload"
            >
              开始上传
            </n-button>
            <n-button
              v-if="uploadStore.isQueueUploading"
              type="primary"
              :loading="true"
            >
              上传中 ({{ uploadStore.queueUploadedCount }}/{{ uploadStore.queueTotalCount }})
            </n-button>
          </div>
        </template>
      </n-card>
    </n-modal>

    <n-modal v-model:show="showFileMenu" v-if="selectedFileItem">
      <n-card :title="selectedFileItem.name" style="width: 400px;">
        <div class="space-y-2">
          <n-button block @click="showDetailModal">详情</n-button>
          <n-button v-if="selectedFileItem.size !== undefined" block @click="downloadFile(selectedFileItem)">下载</n-button>
          <n-button v-if="selectedFileItem.size !== undefined" block type="info" @click="openShareModal">分享</n-button>
          <n-button block @click="openRenameModal">重命名</n-button>
          <n-button block @click="openMoveModal">移动</n-button>
          <n-button block type="error" @click="deleteFile(selectedFileItem)">删除</n-button>
          <n-button block @click="showFileMenu = false">取消</n-button>
        </div>
      </n-card>
    </n-modal>

    <FileDetailModal :item="detailItem" @update="handleUpdateDescription" />

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

    <n-modal v-model:show="showMoveModal">
      <n-card title="移动到文件夹" style="width: 400px;">
        <div class="space-y-4">
          <div v-if="moveTargetFolders.length === 0" class="text-center py-6 text-gray-500">
            暂无文件夹
          </div>
          <div v-else class="max-h-60 overflow-y-auto space-y-2">
            <div
              v-for="folder in moveTargetFolders"
              :key="folder.id"
              :class="[
                'p-3 cursor-pointer hover:bg-[#252525]',
                targetFolderId === folder.id ? 'bg-[#252525] ring-1 ring-blue-500' : 'bg-[#1a1a1a]'
              ]"
              @click="targetFolderId = folder.id"
            >
              <div class="flex items-center gap-2">
                <n-icon size="20" color="#60a5fa"><Folder /></n-icon>
                <span class="text-white">{{ folder.name }}</span>
              </div>
            </div>
          </div>
        </div>
        <template #footer>
          <div class="flex justify-end gap-2">
            <n-button @click="showMoveModal = false">取消</n-button>
            <n-button type="primary" :loading="isMoving" @click="confirmMove">移动</n-button>
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
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useFileStore } from '../stores/file'
import { useUploadStore } from '../stores/upload'
import { useDownloadStore } from '../stores/download'
import { useShareStore } from '../stores/share'
import { encryptFile, generateEncryptionKey } from '../utils/crypto'
import { fileAPI, shareAPI } from '../api'
import { useMessage } from 'naive-ui'
import AppHeader from '../components/AppHeader.vue'
import FileDetailModal from '../components/FileDetailModal.vue'
import {
  ArrowBack, CloudUpload, Refresh, Folder, Create, Document, Close
} from '@vicons/ionicons5'

const router = useRouter()
const authStore = useAuthStore()
const fileStore = useFileStore()
const uploadStore = useUploadStore()
const downloadStore = useDownloadStore()
const shareStore = useShareStore()
const message = useMessage()

const showNewFolderModal = ref(false)
const showUploadModal = ref(false)
const showFileMenu = ref(false)
const showShareModal = ref(false)
const newFolderName = ref('')
const selectedFile = ref(null)
const selectedFileItem = ref(null)
const detailItem = ref(null)

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

// 拖拽上传相关
const isDragging = ref(false)

function handleDragLeave(e) {
  // 只有当拖拽目标确实是当前元素时才关闭遮罩
  if (e.currentTarget.contains(e.relatedTarget)) return
  isDragging.value = false
}

async function handleDrop(e) {
  isDragging.value = false
  const files = Array.from(e.dataTransfer?.files || [])
  if (files.length === 0) return

  // 过滤掉文件夹
  const validFiles = files.filter(f => f.type && !f.webkitRelativePath)
  if (validFiles.length === 0) return

  uploadStore.addToQueue(validFiles)
  showUploadModal.value = true
  if (!uploadStore.isQueueUploading) {
    const key = await generateEncryptionKey()
    await uploadStore.uploadAll(key, fileStore.currentFolder)
  }
}

const currentFolderName = computed(() => {
  const crumbs = fileStore.breadcrumbs
  return crumbs.length > 0 ? crumbs[crumbs.length - 1].name : '根目录'
})

// 选择模式相关
const selectMode = ref(false)
const selectedItems = ref([])

// 搜索和排序相关
const searchQuery = ref('')
const sortBy = ref('name')
const sortOptions = [
  { label: '名称', value: 'name' },
  { label: '大小', value: 'size' },
  { label: '时间', value: 'time' }
]

// 过滤和排序后的文件夹/文件列表
const filteredFolders = computed(() => {
  let result = [...fileStore.folders]
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    result = result.filter(f => f.name.toLowerCase().includes(q))
  }
  result.sort((a, b) => {
    if (sortBy.value === 'name') return a.name.localeCompare(b.name)
    if (sortBy.value === 'time') return new Date(b.created_at || 0) - new Date(a.created_at || 0)
    return 0
  })
  return result
})

const filteredFiles = computed(() => {
  let result = [...fileStore.files]
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    result = result.filter(f => f.name.toLowerCase().includes(q))
  }
  result.sort((a, b) => {
    if (sortBy.value === 'name') return a.name.localeCompare(b.name)
    if (sortBy.value === 'size') return (b.size || 0) - (a.size || 0)
    if (sortBy.value === 'time') return new Date(b.created_at || 0) - new Date(a.created_at || 0)
    return 0
  })
  return result
})

// 批量移动相关
const showMoveModal = ref(false)
const moveTargetFolders = ref([])
const targetFolderId = ref(null)
const isMoving = ref(false)
const moveMode = ref('batch') // 'batch' or 'single'

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
      // 判断是文件还是文件夹
      const isFile = fileStore.files.some(f => f.id === id)
      const isFolder = fileStore.folders.some(f => f.id === id)
      if (isFile) {
        await fileStore.deleteFile(id)
        deletedCount++
      } else if (isFolder) {
        await fileStore.deleteFolder(id)
        deletedCount++
      }
    }
    message.success(`成功删除 ${deletedCount} 个项目`)
    selectedItems.value = []
    selectMode.value = false
    fileStore.fetchFiles(fileStore.currentFolder)
  } catch (err) {
    message.error(err.message || '删除失败')
  }
}

async function batchDownload() {
  if (selectedItems.value.length === 0) return

  try {
    // 只下载文件，不下载文件夹
    const filesToDownload = fileStore.files.filter(f => selectedItems.value.includes(f.id))

    if (filesToDownload.length === 0) {
      message.warning('请选择要下载的文件')
      return
    }

    // 添加到下载队列
    const downloadItems = filesToDownload.map(f => ({
      id: f.id,
      name: f.name,
      size: f.size,
      encryptionKey: f.encryption_key,
      encryptionNonce: f.encryption_nonce
    }))

    downloadStore.addToQueue(downloadItems)
    await downloadStore.downloadAll()

    const successCount = downloadStore.queueCompletedCount
    const errorCount = downloadStore.queueErrors.length

    if (errorCount > 0) {
      message.warning(`下载完成: ${successCount} 个成功, ${errorCount} 个失败`)
    } else {
      message.success(`成功下载 ${successCount} 个文件`)
    }

    downloadStore.clearQueue()
    selectedItems.value = []
    selectMode.value = false
  } catch (err) {
    message.error(err.message || '下载失败')
  }
}

async function openBatchMoveModal() {
  // 获取所有文件夹作为移动目标
  moveTargetFolders.value = [...fileStore.folders]
  // 默认选择当前文件夹
  targetFolderId.value = fileStore.currentFolder
  moveMode.value = 'batch'
  showMoveModal.value = true
}

async function confirmMove() {
  isMoving.value = true

  try {
    if (moveMode.value === 'batch') {
      // 批量移动
      const ids = [...selectedItems.value]
      let movedCount = 0
      for (const id of ids) {
        const isFile = fileStore.files.some(f => f.id === id)
        const isFolder = fileStore.folders.some(f => f.id === id)
        if (isFile) {
          await fileStore.moveFile(id, targetFolderId.value)
          movedCount++
        } else if (isFolder) {
          await fileStore.moveFolder(id, targetFolderId.value)
          movedCount++
        }
      }
      message.success(`成功移动 ${movedCount} 个项目`)
      selectedItems.value = []
      selectMode.value = false
    } else {
      // 单个移动
      if (!selectedFileItem.value) return
      const isFile = selectedFileItem.value.size !== undefined
      if (isFile) {
        await fileStore.moveFile(selectedFileItem.value.id, targetFolderId.value)
      } else {
        await fileStore.moveFolder(selectedFileItem.value.id, targetFolderId.value)
      }
      message.success('移动成功')
    }
    showMoveModal.value = false
    fileStore.fetchFiles(fileStore.currentFolder)
  } catch (err) {
    message.error(err.message || '移动失败')
  } finally {
    isMoving.value = false
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

function handleFilesChange(options) {
  const files = options.fileList.map(item => item.file)
  uploadStore.addToQueue(files)
}

function clearFile() {
  selectedFile.value = null
  uploadStore.reset()
}

function closeUploadModal() {
  if (uploadStore.isQueueUploading) return
  showUploadModal.value = false
  selectedFile.value = null
  uploadStore.reset()
  uploadStore.clearQueue()
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

async function startBatchUpload() {
  if (uploadStore.uploadQueue.length === 0) {
    message.warning('请选择文件')
    return
  }

  try {
    const key = await generateEncryptionKey()
    await uploadStore.uploadAll(key, fileStore.currentFolder)

    const successCount = uploadStore.queueUploadedCount
    const errorCount = uploadStore.queueErrors.length

    if (errorCount > 0) {
      message.warning(`上传完成: ${successCount} 个成功, ${errorCount} 个失败`)
    } else {
      message.success(`成功上传 ${successCount} 个文件`)
    }

    fileStore.fetchFiles(fileStore.currentFolder)
    showUploadModal.value = false
    uploadStore.clearQueue()
  } catch (err) {
    message.error(err.message || '上传失败')
  }
}

function handleFileClick(file) {
  selectedFileItem.value = file
  showFileMenu.value = true
}

function openFileMenu(file) {
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

async function handleUpdateDescription({ id, description }) {
  try {
    const isFile = detailItem.value?.size !== undefined
    if (isFile) {
      await fileStore.updateFileDescription(id, description)
    } else {
      await fileStore.updateFolderDescription(id, description)
    }
    message.success('描述已保存')
  } catch (err) {
    message.error(err.message || '保存失败')
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

function openMoveModal() {
  showFileMenu.value = false
  // 获取所有文件夹作为移动目标
  moveTargetFolders.value = [...fileStore.folders]
  // 默认选择当前文件夹
  targetFolderId.value = fileStore.currentFolder
  moveMode.value = 'single'
  showMoveModal.value = true
}

function showDetailModal() {
  showFileMenu.value = false
  detailItem.value = selectedFileItem.value
}

async function confirmSingleMove() {
  if (!selectedFileItem.value) return
  isMoving.value = true

  try {
    const isFile = selectedFileItem.value.size !== undefined
    if (isFile) {
      await fileStore.moveFile(selectedFileItem.value.id, targetFolderId.value)
    } else {
      await fileStore.moveFolder(selectedFileItem.value.id, targetFolderId.value)
    }
    message.success('移动成功')
    showMoveModal.value = false
    fileStore.fetchFiles(fileStore.currentFolder)
  } catch (err) {
    message.error(err.message || '移动失败')
  } finally {
    isMoving.value = false
  }
}

async function confirmRename() {
  if (!renameName.value.trim() || !selectedFileItem.value) return
  try {
    // 判断是文件还是文件夹
    if (selectedFileItem.value.size !== undefined) {
      await fileStore.renameFile(selectedFileItem.value.id, renameName.value)
    } else {
      await fileStore.renameFolder(selectedFileItem.value.id, renameName.value)
    }
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
