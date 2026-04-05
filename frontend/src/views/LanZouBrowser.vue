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
          <n-button v-if="selectedItems.length === 1" @click="openRenameModalForSelected">重命名</n-button>
          <n-button v-if="selectedItems.length > 0" @click="oneClickShare">一键分享</n-button>
          <n-button v-if="selectedItems.length > 0" @click="openMoveModal">移动到</n-button>
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
              @contextmenu.prevent="openFolderContextMenu($event, folder)"
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
              @contextmenu.prevent="openFileContextMenu($event, file)"
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
          <n-button block @click="openRenameModal('file')">重命名</n-button>
          <n-button block @click="openAccessModal('file')">设置访问密码</n-button>
          <n-button block @click="openFileDescModal()">文件描述</n-button>
          <n-button block @click="showFileMenu = false">取消</n-button>
        </div>
      </n-card>
    </n-modal>

    <n-modal v-model:show="showFolderMenu" v-if="selectedFolder">
      <n-card :title="selectedFolder.name" style="width: 400px;">
        <div class="space-y-2">
          <n-button block type="primary" @click="navigateToFolder(selectedFolder)">打开</n-button>
          <n-button block @click="openRenameModal('folder')">重命名</n-button>
          <n-button block @click="openAccessModal('folder')">设置访问密码</n-button>
          <n-button block @click="showFolderMenu = false">取消</n-button>
        </div>
      </n-card>
    </n-modal>

    <n-modal v-model:show="showRenameModal">
      <n-card title="重命名" style="width: 400px;">
        <div class="space-y-4">
          <n-input
            v-model:value="renameForm.name"
            placeholder="请输入新名称"
            maxlength="100"
            @keyup.enter="submitRename"
          />
        </div>
        <template #footer>
          <div class="flex justify-end gap-2">
            <n-button @click="showRenameModal = false">取消</n-button>
            <n-button type="primary" @click="submitRename" :disabled="!renameForm.name.trim()">确定</n-button>
          </div>
        </template>
      </n-card>
    </n-modal>

    <n-modal v-model:show="showAccessModal">
      <n-card title="设置访问密码" style="width: 400px;">
        <div class="space-y-4">
          <n-radio-group v-model:value="accessForm.type">
            <n-space>
              <n-radio value="2">密码访问</n-radio>
              <n-radio value="1">公开访问</n-radio>
            </n-space>
          </n-radio-group>
          <n-input
            v-if="accessForm.type === '2'"
            v-model:value="accessForm.password"
            placeholder="请输入访问密码"
            maxlength="20"
          />
        </div>
        <template #footer>
          <div class="flex justify-end gap-2">
            <n-button @click="showAccessModal = false">取消</n-button>
            <n-button type="primary" @click="submitAccess">确定</n-button>
          </div>
        </template>
      </n-card>
    </n-modal>

    <n-modal v-model:show="showFileDescModal">
      <n-card title="文件描述" style="width: 400px;">
        <div class="space-y-4">
          <n-input
            v-model:value="fileDescForm.description"
            type="textarea"
            placeholder="请输入文件描述"
            :autosize="{ minRows: 3, maxRows: 6 }"
            maxlength="500"
          />
        </div>
        <template #footer>
          <div class="flex justify-end gap-2">
            <n-button @click="showFileDescModal = false">取消</n-button>
            <n-button type="primary" @click="submitFileDesc" :loading="descSubmitting">确定</n-button>
          </div>
        </template>
      </n-card>
    </n-modal>

    <n-modal v-model:show="showMoveModal">
      <n-card title="移动到" style="width: 400px;">
        <div class="space-y-4">
          <n-select
            v-model:value="moveTargetId"
            :options="folderSelectOptions"
            placeholder="选择目标文件夹"
          />
        </div>
        <template #footer>
          <div class="flex justify-end gap-2">
            <n-button @click="showMoveModal = false">取消</n-button>
            <n-button type="primary" @click="submitMove" :disabled="!moveTargetId">确定</n-button>
          </div>
        </template>
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
import { ref, computed, onMounted } from 'vue'
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

const showFolderMenu = ref(false)
const selectedFolder = ref(null)

const showAccessModal = ref(false)
const accessForm = ref({
  type: '2',
  password: '',
  targetId: null,
  targetType: 'file',
})

const showRenameModal = ref(false)
const renameForm = ref({
  name: '',
  targetId: null,
  targetType: 'file',
})

const showFileDescModal = ref(false)
const descSubmitting = ref(false)
const fileDescForm = ref({
  description: '',
  fileId: null,
})

const showMoveModal = ref(false)
const moveTargetId = ref(null)

// 文件夹选择下拉选项 (从面包屑构建 + 根目录)
const folderSelectOptions = computed(() => {
  const opts = [{ label: '根目录', value: -1 }]
  for (const crumb of breadcrumbs.value) {
    if (crumb.id !== -1) {
      opts.push({ label: crumb.name, value: crumb.id })
    }
  }
  return opts
})

// 选择模式相关

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
    const items = selectedItems.value.map(id => {
      const folder = folders.value.find(f => f.folder_id === id)
      const file = files.value.find(f => f.file_id === id)
      if (folder) return { id: folder.folder_id, type: 'folder' }
      if (file) return { id: file.file_id, type: 'file' }
      return null
    }).filter(Boolean)

    const res = await lanzouAPI.batchDelete({ items })
    message.success(`成功删除 ${res.deleted} 个项目`)
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

function openFileContextMenu(event, file) {
  if (selectMode.value) return
  selectedFile.value = file
  showFileMenu.value = true
}

function openFolderContextMenu(event, folder) {
  if (selectMode.value) return
  selectedFolder.value = folder
  showFolderMenu.value = true
}

function openAccessModal(type) {
  if (type === 'file' && selectedFile.value) {
    accessForm.value.targetId = selectedFile.value.file_id
    accessForm.value.targetType = 'file'
  } else if (type === 'folder' && selectedFolder.value) {
    accessForm.value.targetId = selectedFolder.value.folder_id
    accessForm.value.targetType = 'folder'
  }
  accessForm.value.type = '2'
  accessForm.value.password = ''
  showAccessModal.value = true
  showFileMenu.value = false
  showFolderMenu.value = false
}

async function submitAccess() {
  if (accessForm.value.type === '2' && !accessForm.value.password.trim()) {
    message.warning('请输入访问密码')
    return
  }
  try {
    await lanzouAPI.setAccess({
      id: accessForm.value.targetId,
      type: accessForm.value.targetType,
      shows: parseInt(accessForm.value.type),
      shownames: accessForm.value.type === '2' ? accessForm.value.password : '',
    })
    message.success('访问密码设置成功')
    showAccessModal.value = false
    accessForm.value = { targetId: null, targetType: 'file', type: '2', password: '' }
  } catch (err) {
    message.error(err.message || '设置访问密码失败')
  }
}

function openRenameModal(type) {
  if (type === 'file' && selectedFile.value) {
    renameForm.value.targetId = selectedFile.value.file_id
    renameForm.value.targetType = 'file'
    renameForm.value.name = selectedFile.value.name
  } else if (type === 'folder' && selectedFolder.value) {
    renameForm.value.targetId = selectedFolder.value.folder_id
    renameForm.value.targetType = 'folder'
    renameForm.value.name = selectedFolder.value.name
  }
  showRenameModal.value = true
  showFileMenu.value = false
  showFolderMenu.value = false
}

function openRenameModalForSelected() {
  if (selectedItems.value.length !== 1) return
  const id = selectedItems.value[0]
  const folder = folders.value.find(f => f.folder_id === id)
  const file = files.value.find(f => f.file_id === id)

  if (folder) {
    renameForm.value.targetId = folder.folder_id
    renameForm.value.targetType = 'folder'
    renameForm.value.name = folder.name
  } else if (file) {
    renameForm.value.targetId = file.file_id
    renameForm.value.targetType = 'file'
    renameForm.value.name = file.name
  }
  showRenameModal.value = true
}

async function submitRename() {
  if (!renameForm.value.name.trim()) {
    message.warning('请输入名称')
    return
  }
  try {
    const sanitizedName = renameForm.value.name.replace(/[ ()]/g, '_')
    await lanzouAPI.rename(renameForm.value.targetId, {
      type: renameForm.value.targetType,
      name: sanitizedName,
    })
    message.success('重命名成功')
    showRenameModal.value = false
    renameForm.value = { name: '', targetId: null, targetType: 'file' }
    await refresh()
  } catch (err) {
    message.error(err.message || '重命名失败')
  }
}

async function oneClickShare() {
  if (selectedItems.value.length === 0) return

  const items = selectedItems.value.map(id => {
    const folder = folders.value.find(f => f.folder_id === id)
    const file = files.value.find(f => f.file_id === id)
    if (folder) return { id: folder.folder_id, name: folder.name, type: 'folder' }
    if (file) return { id: file.file_id, name: file.name, type: 'file' }
    return null
  }).filter(Boolean)

  let shareText = ''
  for (const item of items) {
    try {
      const res = await lanzouAPI.createShare({ file_id: item.id, minutes: 0 })
      shareText += `${item.name} ${res.url}${res.pwd ? ` 密码:${res.pwd}` : ''}\n`
    } catch (err) {
      shareText += `${item.name} 分享失败: ${err.message}\n`
    }
  }

  try {
    await navigator.clipboard.writeText(shareText.trim())
    message.success('分享链接已复制')
  } catch {
    message.warning('复制失败，请手动复制')
  }
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

async function openFileDescModal() {
  showFileMenu.value = false
  if (!selectedFile.value) return

  fileDescForm.value.fileId = selectedFile.value.file_id
  fileDescForm.value.description = ''
  showFileDescModal.value = true

  // Fetch existing description
  try {
    const res = await lanzouAPI.getFileDescription(selectedFile.value.file_id)
    if (res.info) {
      fileDescForm.value.description = res.info
    }
  } catch {
    // Ignore fetch errors
  }
}

async function submitFileDesc() {
  if (!fileDescForm.value.description.trim()) {
    message.warning('请输入文件描述')
    return
  }
  descSubmitting.value = true
  try {
    await lanzouAPI.setFileDescription(fileDescForm.value.fileId, {
      description: fileDescForm.value.description,
    })
    message.success('文件描述设置成功')
    showFileDescModal.value = false
    fileDescForm.value = { description: '', fileId: null }
  } catch (err) {
    message.error(err.message || '设置文件描述失败')
  } finally {
    descSubmitting.value = false
  }
}

function openMoveModal() {
  moveTargetId.value = currentFolderId.value
  showMoveModal.value = true
}

async function submitMove() {
  if (!moveTargetId.value) {
    message.warning('请选择目标文件夹')
    return
  }
  try {
    const items = selectedItems.value.map(id => {
      const folder = folders.value.find(f => f.folder_id === id)
      const file = files.value.find(f => f.file_id === id)
      if (folder) return { id: folder.folder_id, type: 'folder' }
      if (file) return { id: file.file_id, type: 'file' }
      return null
    }).filter(Boolean)

    const res = await lanzouAPI.batchMove({ items, target: moveTargetId.value })
    message.success(`成功移动 ${res.moved} 个项目`)
    showMoveModal.value = false
    selectedItems.value = []
    selectMode.value = false
    await refresh()
  } catch (err) {
    message.error(err.message || '移动失败')
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
