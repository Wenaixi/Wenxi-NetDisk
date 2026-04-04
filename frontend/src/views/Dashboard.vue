<template>
  <div class="min-h-screen bg-[#0f0f0f]">
    <header class="bg-[#1a1a1a] px-6 py-4 flex items-center justify-between">
      <div class="flex items-center gap-4">
        <h1 class="text-xl font-bold text-white">文希云盘</h1>
        <nav class="flex items-center gap-2 ml-8">
          <n-button text @click="fileStore.navigateToRoot()">首页</n-button>
          <n-breadcrumb>
            <n-breadcrumb-item v-for="(crumb, idx) in fileStore.breadcrumbs" :key="crumb.id">
              <n-button text @click="navigateToBreadcrumb(idx)">{{ crumb.name }}</n-button>
            </n-breadcrumb-item>
          </n-breadcrumb>
        </nav>
      </div>
      <div class="flex items-center gap-4">
        <n-dropdown :options="userMenuOptions" @select="handleUserMenu">
          <n-button text class="text-gray-300">
            {{ authStore.user?.email }}
            <n-icon><component :is="ChevronDown" /></n-icon>
          </n-button>
        </n-dropdown>
      </div>
    </header>

    <main class="p-6">
      <div class="flex items-center justify-between mb-6">
        <div class="flex items-center gap-2">
          <n-button @click="fileStore.navigateBack()" :disabled="fileStore.breadcrumbs.length <= 1">
            <template #icon><n-icon><ArrowBack /></n-icon></template>
            返回
          </n-button>
          <n-button @click="showUploadModal = true">
            <template #icon><n-icon><CloudUpload /></n-icon></template>
            上传文件
          </n-button>
          <n-button @click="showNewFolderModal = true">
            <template #icon><n-icon><Create" /></n-icon></template>
            新建文件夹
          </n-button>
        </div>
        <n-button @click="refresh">
          <template #icon><n-icon><Refresh" /></n-icon></template>
        </n-button>
      </div>

      <n-spin :show="fileStore.loading">
        <div v-if="fileStore.folders.length === 0 && fileStore.files.length === 0" class="text-center py-12">
          <n-empty description="文件夹为空" />
        </div>

        <div v-else>
          <div class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4">
            <div
              v-for="folder in fileStore.folders"
              :key="folder.id"
              class="bg-[#1a1a1a] p-4 cursor-pointer hover:bg-[#252525]"
              @dblclick="fileStore.navigateToFolder(folder)"
            >
              <div class="flex items-center gap-2">
                <n-icon size="24" color="#60a5fa"><Folder" /></n-icon>
                <span class="text-white truncate">{{ folder.name }}</span>
              </div>
            </div>

            <div
              v-for="file in fileStore.files"
              :key="file.id"
              class="bg-[#1a1a1a] p-4"
            >
              <div class="flex items-center gap-2 mb-2">
                <n-icon size="24" color="#a78bfa"><Document" /></n-icon>
                <span class="text-white truncate">{{ file.name }}</span>
              </div>
              <div class="text-gray-500 text-sm">{{ formatSize(file.size) }}</div>
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

    <n-modal v-model:show="showUploadModal">
      <n-card title="上传文件" style="width: 400px;">
        <n-upload
          :max="1"
          @change="handleFileChange"
          @remove="handleFileRemove"
        >
          <n-button>选择文件</n-button>
        </n-upload>
        <template #footer>
          <div class="flex justify-end gap-2">
            <n-button @click="showUploadModal = false">取消</n-button>
            <n-button type="primary" :loading="uploading" @click="uploadFile">上传</n-button>
          </div>
        </template>
      </n-card>
    </n-modal>
  </div>
</template>

<script setup>
import { ref, onMounted, h } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useFileStore } from '../stores/file'
import { fileAPI } from '../api'
import { NIcon, useMessage } from 'naive-ui'
import {
  ChevronDown, ArrowBack, CloudUpload, Refresh, Folder, Create, Document, LogOutOutline
} from '@vicons/ionicons5'

const router = useRouter()
const authStore = useAuthStore()
const fileStore = useFileStore()
const message = useMessage()

const showNewFolderModal = ref(false)
const showUploadModal = ref(false)
const newFolderName = ref('')
const uploading = ref(false)
const selectedFile = ref(null)

const userMenuOptions = [
  { label: '退出登录', key: 'logout', icon: () => h(NIcon, null, { default: () => h(LogOutOutline) }) }
]

function handleUserMenu(key) {
  if (key === 'logout') {
    authStore.logout()
    router.push('/login')
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
}

function handleFileRemove() {
  selectedFile.value = null
}

async function uploadFile() {
  if (!selectedFile.value) {
    message.warning('请选择文件')
    return
  }
  uploading.value = true
  try {
    const formData = new FormData()
    formData.append('file', selectedFile.value)
    formData.append('folder_id', fileStore.currentFolder || '')
    await fileAPI.upload(formData)
    message.success('上传成功')
    fileStore.fetchFiles(fileStore.currentFolder)
    showUploadModal.value = false
    selectedFile.value = null
  } catch (err) {
    message.error(err.message || '上传失败')
  } finally {
    uploading.value = false
  }
}

async function refresh() {
  await fileStore.fetchFiles(fileStore.currentFolder)
  message.success('刷新成功')
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
  await fileStore.fetchFiles()
})
</script>