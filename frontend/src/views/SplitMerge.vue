<template>
  <div class="p-6">
    <div class="flex gap-4 mb-6">
      <n-button :type="activeTab === 'split' ? 'primary' : 'default'" @click="activeTab = 'split'">文件分割</n-button>
      <n-button :type="activeTab === 'merge' ? 'primary' : 'default'" @click="activeTab = 'merge'">文件合并</n-button>
    </div>

    <!-- 分割页面 -->
    <div v-if="activeTab === 'split'">
      <div class="mb-4">
        <n-button @click="triggerFileInput">选择文件</n-button>
        <input ref="splitFileInput" type="file" class="hidden" @change="handleSplitFileSelect" />
      </div>

      <div v-if="selectedSplitFile" class="bg-[#1a1a1a] p-4 mb-4">
        <div class="flex justify-between mb-2">
          <span class="text-gray-400">文件名</span>
          <span class="text-white">{{ selectedSplitFile.name }}</span>
        </div>
        <div class="flex justify-between mb-2">
          <span class="text-gray-400">文件大小</span>
          <span class="text-white">{{ formatFileSize(selectedSplitFile.size) }}</span>
        </div>
        <div class="flex justify-between mb-2">
          <span class="text-gray-400">预计分块数</span>
          <span class="text-white">{{ splitChunks }} 块 (每块 {{ formatFileSize(chunkSize) }})</span>
        </div>
      </div>

      <div class="bg-[#1a1a1a] p-4 mb-4">
        <div class="flex items-center gap-4 mb-4">
          <span class="text-gray-400">分块大小</span>
          <n-input-number v-model:value="chunkSize" :min="1024 * 1024" :max="10 * 1024 * 1024" :step="1024 * 1024" />
          <span class="text-gray-500">{{ formatFileSize(chunkSize) }}</span>
        </div>
        <n-button type="primary" :loading="splitting" :disabled="!selectedSplitFile" @click="doSplit">
          {{ splitting ? '分割中...' : '开始分割' }}
        </n-button>
      </div>

      <div v-if="splitResult.length > 0" class="divide-y divide-gray-700">
        <div v-for="(chunk, idx) in splitResult" :key="idx" class="flex justify-between py-2 px-4">
          <span class="text-white">{{ chunk.name }}</span>
          <span class="text-gray-400">{{ formatFileSize(chunk.size) }}</span>
          <n-button size="small" @click="downloadChunk(chunk)">下载</n-button>
        </div>
      </div>
    </div>

    <!-- 合并页面 -->
    <div v-if="activeTab === 'merge'">
      <div class="mb-4">
        <n-button @click="triggerMergeFileInput">选择分块文件</n-button>
        <input ref="mergeFileInput" type="file" class="hidden" multiple @change="handleMergeFileSelect" />
      </div>

      <div v-if="mergeFiles.length > 0" class="bg-[#1a1a1a] p-4 mb-4">
        <div class="text-gray-400 mb-2">已选择 {{ mergeFiles.length }} 个分块</div>
        <div v-for="(f, idx) in mergeFiles" :key="idx" class="flex justify-between py-1">
          <span class="text-white text-sm">{{ f.name }}</span>
          <n-button size="small" type="error" @click="removeMergeFile(idx)">移除</n-button>
        </div>
        <n-button type="primary" class="mt-4" :loading="merging" @click="doMerge">
          {{ merging ? '合并中...' : '开始合并' }}
        </n-button>
      </div>

      <div v-if="mergeComplete" class="bg-[#1a1a1a] p-4">
        <div class="text-green-400 mb-2">合并完成!</div>
        <div class="flex justify-between mb-2">
          <span class="text-gray-400">文件名</span>
          <span class="text-white">{{ mergedFileName }}</span>
        </div>
        <div class="flex justify-between mb-2">
          <span class="text-gray-400">文件大小</span>
          <span class="text-white">{{ formatFileSize(mergedBlob?.size || 0) }}</span>
        </div>
        <n-button @click="downloadMergedFile">下载合并文件</n-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { splitFile as splitFileUtil, mergeBlobs, getChunkFileName, parseChunkFileName, formatFileSize, DEFAULT_CHUNK_SIZE } from '../utils/fileSplit'

const activeTab = ref('split')

// 分割相关
const splitFileInput = ref(null)
const selectedSplitFile = ref(null)
const chunkSize = ref(2 * 1024 * 1024)
const splitting = ref(false)
const splitResult = ref([])

const splitChunks = computed(() => {
  if (!selectedSplitFile.value) return 0
  return Math.ceil(selectedSplitFile.value.size / chunkSize.value)
})

// 合并相关
const mergeFileInput = ref(null)
const mergeFiles = ref([])
const merging = ref(false)
const mergeComplete = ref(false)
const mergedBlob = ref(null)
const mergedFileName = ref('')

function triggerFileInput() {
  splitFileInput.value?.click()
}

function triggerMergeFileInput() {
  mergeFileInput.value?.click()
}

async function handleSplitFileSelect(e) {
  const file = e.target.files?.[0]
  if (file) {
    selectedSplitFile.value = file
    splitResult.value = []
  }
}

function handleMergeFileSelect(e) {
  const files = Array.from(e.target.files || [])
  if (files.length > 0) {
    mergeFiles.value = [...mergeFiles.value, ...files]
    mergeComplete.value = false
  }
}

function removeMergeFile(idx) {
  mergeFiles.value.splice(idx, 1)
}

async function doSplit() {
  if (!selectedSplitFile.value) return
  splitting.value = true
  try {
    const { chunks } = await splitFileUtil(selectedSplitFile.value, chunkSize.value)
    splitResult.value = chunks.map((chunk, idx) => ({
      name: getChunkFileName(selectedSplitFile.value.name, idx + 1, chunks.length),
      blob: chunk.blob,
      size: chunk.size
    }))
  } finally {
    splitting.value = false
  }
}

function downloadChunk(chunk) {
  const url = URL.createObjectURL(chunk.blob)
  const a = document.createElement('a')
  a.href = url
  a.download = chunk.name
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
}

async function doMerge() {
  if (mergeFiles.value.length === 0) return
  merging.value = true
  try {
    // 尝试解析第一个文件名来获取原始文件名
    const parsed = parseChunkFileName(mergeFiles.value[0].name)
    mergedFileName.value = parsed ? parsed.name : 'merged_file'

    const blob = mergeBlobs(mergeFiles.value, mergedFileName.value)
    mergedBlob.value = blob
    mergeComplete.value = true
  } finally {
    merging.value = false
  }
}

function downloadMergedFile() {
  if (!mergedBlob.value) return
  const url = URL.createObjectURL(mergedBlob.value)
  const a = document.createElement('a')
  a.href = url
  a.download = mergedFileName.value
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
}
</script>
