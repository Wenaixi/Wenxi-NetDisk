<template>
  <header class="bg-[#1a1a1a] px-6 py-4 flex items-center justify-between">
    <div class="flex items-center gap-4">
      <h1 class="text-xl font-bold text-white">文希云盘</h1>
      <nav class="flex items-center gap-2 ml-8">
        <n-button text @click="$router.push('/dashboard')" :type="isActive('/dashboard')">
          本地文件
        </n-button>
        <n-button text @click="$router.push('/lanzou')" :type="isActive('/lanzou')">
          蓝奏云
        </n-button>
        <n-button text @click="$router.push('/sync')" :type="isActive('/sync')">
          同步任务
        </n-button>
        <n-button text @click="$router.push('/share-parse')" :type="isActive('/share-parse')">
          链接解析
        </n-button>
        <n-button text @click="$router.push('/split-merge')" :type="isActive('/split-merge')">
          分割/合并
        </n-button>
        <n-button text @click="$router.push('/task-history')" :type="isActive('/task-history')">
          任务历史
        </n-button>
      </nav>
    </div>
    <div class="flex items-center gap-4">
      <n-button text @click="$router.push('/upload-tasks')">
        上传任务
      </n-button>
      <n-button text @click="$router.push('/download-tasks')">
        下载任务
      </n-button>
      <n-button text @click="$router.push('/recycle')" :type="isActive('/recycle')">
        回收站
      </n-button>
      <n-dropdown :options="menuOptions" @select="handleMenuSelect">
        <n-button text class="text-gray-300">
          {{ user?.email || '用户' }}
          <n-icon><ChevronDown /></n-icon>
        </n-button>
      </n-dropdown>
    </div>
  </header>
</template>

<script setup>
import { computed, h } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { NIcon } from 'naive-ui'
import { ChevronDown, Settings, LogOutOutline } from '@vicons/ionicons5'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const user = computed(() => authStore.user)

const menuOptions = [
  {
    label: '设置',
    key: 'settings',
    icon: () => h(NIcon, null, { default: () => h(Settings) })
  },
  {
    label: '退出登录',
    key: 'logout',
    icon: () => h(NIcon, null, { default: () => h(LogOutOutline) })
  }
]

function isActive(path) {
  return route.path === path ? 'primary' : undefined
}

function handleMenuSelect(key) {
  if (key === 'settings') {
    router.push('/settings')
  } else if (key === 'logout') {
    authStore.logout()
    router.push('/login')
  }
}
</script>
