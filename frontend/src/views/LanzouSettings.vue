<template>
  <div class="min-h-screen bg-[#0f0f0f]">
    <header class="bg-[#1a1a1a] px-6 py-4 flex items-center justify-between">
      <div class="flex items-center gap-4">
        <h1 class="text-xl font-bold text-white">文希云盘</h1>
        <nav class="flex items-center gap-2 ml-8">
          <n-button text @click="$router.push('/dashboard')">首页</n-button>
        </nav>
      </div>
      <div class="flex items-center gap-4">
        <n-button text class="text-gray-300" @click="$router.push('/settings')">
          {{ authStore.user?.email }}
        </n-button>
      </div>
    </header>

    <main class="p-6 max-w-2xl mx-auto">
      <h2 class="text-2xl font-bold text-white mb-6">蓝奏云设置</h2>

      <div class="bg-[#1a1a1a] p-6 mb-6">
        <h3 class="text-lg text-white mb-4">连接状态</h3>

        <div v-if="connected" class="space-y-4">
          <div class="flex items-center gap-3 text-green-500">
            <n-icon size="20"><CheckmarkCircle" /></n-icon>
            <span>已连接到蓝奏云</span>
          </div>
          <div class="text-gray-400 text-sm">
            用户ID: {{ statusInfo?.uid || '未知' }}
          </div>
          <n-button type="error" @click="handleDisconnect">断开连接</n-button>
        </div>

        <div v-else class="space-y-4">
          <div class="flex items-center gap-3 text-yellow-500">
            <n-icon size="20"><Warning" /></n-icon>
            <span>未连接蓝奏云</span>
          </div>
          <p class="text-gray-400 text-sm">
            连接后可以同步蓝奏云文件，享受更快的上传下载速度。
          </p>
          <n-button type="primary" @click="showConnectModal = true">连接蓝奏云</n-button>
        </div>
      </div>

      <div class="bg-[#1a1a1a] p-6">
        <h3 class="text-lg text-white mb-4">帮助</h3>
        <div class="text-gray-400 text-sm space-y-2">
          <p>1. 打开蓝奏云网页版并登录</p>
          <p>2. 按F12打开开发者工具，进入Network标签</p>
          <p>3. 刷新页面，找到任意请求，复制Cookie值</p>
          <p>4. 将Cookie粘贴到上方输入框中</p>
        </div>
      </div>
    </main>

    <!-- 连接蓝奏云弹窗 -->
    <n-modal v-model:show="showConnectModal">
      <n-card title="连接蓝奏云" style="width: 500px;">
        <n-form :model="connectForm" :rules="rules" ref="formRef">
          <n-form-item label="Cookie" path="cookie">
            <n-input
              v-model:value="connectForm.cookie"
              type="textarea"
              placeholder="请输入蓝奏云Cookie"
              :rows="4"
            />
          </n-form-item>
        </n-form>

        <template #footer>
          <div class="flex justify-end gap-2">
            <n-button @click="showConnectModal = false">取消</n-button>
            <n-button type="primary" :loading="connecting" @click="handleConnect">连接</n-button>
          </div>
        </template>
      </n-card>
    </n-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { lanzouAPI } from '../api/lanzou'
import { NButton, NIcon, NModal, NCard, NForm, NFormItem, NInput, useMessage } from 'naive-ui'
import { CheckmarkCircle, Warning } from '@vicons/ionicons5'

const authStore = useAuthStore()
const message = useMessage()

const connected = ref(false)
const statusInfo = ref(null)
const showConnectModal = ref(false)
const connecting = ref(false)
const formRef = ref(null)

const connectForm = reactive({
  cookie: ''
})

const rules = {
  cookie: [{ required: true, message: '请输入Cookie', trigger: 'blur' }]
}

async function fetchStatus() {
  try {
    const res = await lanzouAPI.status()
    connected.value = res.connected
    statusInfo.value = res.info || null
  } catch {
    connected.value = false
  }
}

async function handleConnect() {
  try {
    connecting.value = true
    await lanzouAPI.connect({ cookie: connectForm.cookie })
    message.success('连接成功')
    showConnectModal.value = false
    connectForm.cookie = ''
    await fetchStatus()
  } catch (err) {
    message.error(err.message || '连接失败')
  } finally {
    connecting.value = false
  }
}

async function handleDisconnect() {
  try {
    await lanzouAPI.disconnect()
    message.success('已断开连接')
    connected.value = false
    statusInfo.value = null
  } catch (err) {
    message.error(err.message || '断开失败')
  }
}

onMounted(() => {
  fetchStatus()
})
</script>
