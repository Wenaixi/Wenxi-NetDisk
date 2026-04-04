<template>
  <div class="min-h-screen bg-[#0f0f0f]">
    <AppHeader />

    <main class="p-6 max-w-4xl mx-auto">
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
            :loading="shareParseStore.loading"
            @click="parseLinks"
            style="height: 100%"
          >
            {{ !shareUrl ? '粘贴并解析' : '解析' }}
          </n-button>
        </div>
      </div>

      <!-- 错误提示 -->
      <n-alert v-if="shareParseStore.error" type="error" class="mb-4" closable>
        {{ shareParseStore.error }}
      </n-alert>

      <!-- 解析结果 -->
      <n-spin :show="shareParseStore.loading">
        <n-empty v-if="shareParseStore.parsedShares.length === 0 && !shareParseStore.currentShare" description="暂无解析结果" class="py-12" />

        <template v-else>
          <!-- 当前分享信息 -->
          <div v-if="shareParseStore.currentShare" class="bg-[#1a1a1a] p-4 mb-4">
            <div class="flex items-center justify-between mb-2">
              <div class="flex items-center gap-2">
                <n-icon size="20" :component="shareParseStore.currentShare.type === 'folder' ? Folder : DocumentText" />
                <span class="text-white">{{ shareParseStore.currentShare.name }}</span>
              </div>
              <span v-if="shareParseStore.currentShare.size" class="text-gray-500 text-sm">
                {{ shareParseStore.currentShare.size }}
              </span>
            </div>

            <div v-if="shareParseStore.hasParsedShares" class="flex items-center gap-2 mt-2">
              <n-button size="small" @click="downloadAll">
                <template #icon><n-icon><CloudDownload /></n-icon></template>
                下载全部 ({{ shareParseStore.parsedShares.length }})
              </n-button>
            </div>
          </div>

          <!-- 文件列表 -->
          <div v-if="shareParseStore.parsedShares.length > 0" class="space-y-1">
            <div
              v-for="(item, idx) in shareParseStore.parsedShares"
              :key="idx"
              class="bg-[#1a1a1a] p-3 flex items-center justify-between"
            >
              <div class="flex items-center gap-3">
                <n-icon size="18" :component="DocumentText" />
                <div>
                  <div class="text-white text-sm">{{ item.name }}</div>
                  <div class="text-gray-500 text-xs">
                    <span v-if="item.size">{{ item.size }}</span>
                    <span v-if="item.time" class="mx-1">·</span>
                    <span v-if="item.time">{{ item.time }}</span>
                  </div>
                </div>
              </div>
              <n-button size="small" type="info" @click="downloadSingle(item)">
                下载
              </n-button>
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
          <p>• 支持多行批量解析，每行一个链接</p>
          <p>• 如分享有密码，请先填入密码再解析</p>
          <p>• 解析后可一键下载单个文件或全部文件</p>
        </div>
      </n-card>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { DocumentText, Folder, CloudDownload } from '@vicons/ionicons5'
import { useMessage } from 'naive-ui'
import { useShareParseStore } from '../stores/shareParse'
import AppHeader from '../components/AppHeader.vue'

const message = useMessage()
const shareParseStore = useShareParseStore()
const shareUrl = ref('')
const sharePwd = ref('')

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
    // 尝试从剪贴板读取
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

  try {
    // 支持多行解析
    const lines = shareUrl.value.trim().split('\n').filter(Boolean)
    shareParseStore.clearParsed()

    for (const line of lines) {
      const url = line.trim()
      if (!url.startsWith('http')) continue
      await shareParseStore.parseShare(url, sharePwd.value)
    }

    if (shareParseStore.parsedShares.length > 0) {
      message.success(`解析成功，共 ${shareParseStore.parsedShares.length} 个文件`)
    }
  } catch (err) {
    message.error(err.message || '解析失败')
  }
}

function downloadAll() {
  message.info('正在准备下载，请稍候...')
  // TODO: 集成下载逻辑
}

async function downloadSingle(item) {
  try {
    const url = await shareParseStore.downloadShareFile(item.url, sharePwd.value)
    // 触发下载
    window.open(url, '_blank')
    message.success('已开始下载')
  } catch (err) {
    message.error(err.message || '下载失败')
  }
}
</script>
