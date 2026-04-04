/**
 * 文件分割与合并工具
 * 用于蓝奏云大文件上传（蓝奏云单文件大小有限制）
 * 支持在浏览器环境中对Blob对象进行分割和合并
 */

// 默认分块大小 2MB（蓝奏云推荐大小）
export const DEFAULT_CHUNK_SIZE = 2 * 1024 * 1024

/**
 * 获取文件信息
 * @param {File} file 文件对象
 * @returns {{ name: string, size: number, type: string, chunks: number }}
 */
export function getFileInfo(file) {
  return {
    name: file.name,
    size: file.size,
    type: file.type || 'application/octet-stream',
    chunks: Math.ceil(file.size / DEFAULT_CHUNK_SIZE)
  }
}

/**
 * 判断文件是否需要分割
 * @param {File} file 文件对象
 * @param {number} [maxSize] 最大文件大小（默认20MB，蓝奏云限制）
 * @returns {boolean}
 */
export function needsSplit(file, maxSize = 20 * 1024 * 1024) {
  return file.size > maxSize
}

/**
 * 将文件分割成多个Chunk
 * @param {File} file 原始文件
 * @param {number} [chunkSize] 每个分块的大小（字节）
 * @returns {Promise<{ chunks: Blob[], total: number }>} 分块数组
 */
export async function splitFile(file, chunkSize = DEFAULT_CHUNK_SIZE) {
  const total = Math.ceil(file.size / chunkSize)
  const chunks = []

  for (let i = 0; i < total; i++) {
    const start = i * chunkSize
    const end = Math.min(start + chunkSize, file.size)
    const chunk = file.slice(start, end)
    chunks.push({
      index: i,
      blob: chunk,
      size: chunk.size,
      start,
      end,
      isLast: i === total - 1
    })
  }

  return { chunks, total }
}

/**
 * 合并多个Blob为一个Blob
 * @param {Blob[]} blobs Blob数组
 * @param {string} [fileName] 合并后的文件名
 * @returns {Blob} 合并后的Blob
 */
export function mergeBlobs(blobs, fileName) {
  return new Blob(blobs, { type: 'application/octet-stream' })
}

/**
 * 生成带分块信息的文件名
 * @param {string} originalName 原始文件名
 * @param {number} index 分块索引 (从1开始)
 * @param {number} total 总分块数
 * @returns {string} 带分块信息的文件名
 */
export function getChunkFileName(originalName, index, total) {
  const dotIndex = originalName.lastIndexOf('.')
  if (dotIndex > 0) {
    const name = originalName.substring(0, dotIndex)
    const ext = originalName.substring(dotIndex)
    return `${name}.part${String(index).padStart(3, '0')}of${total}${ext}`
  }
  return `${originalName}.part${String(index).padStart(3, '0')}of${total}`
}

/**
 * 解析分块文件名，提取原始文件名和分块信息
 * @param {string} chunkFileName 分块文件名
 * @returns {{ name: string, index: number, total: number } | null}
 */
export function parseChunkFileName(chunkFileName) {
  const match = chunkFileName.match(/^(.+?)\.part(\d+)of(\d+)(\..+)?$/)
  if (match) {
    return {
      name: match[1] + (match[4] || ''),
      index: parseInt(match[2], 10),
      total: parseInt(match[3], 10)
    }
  }
  return null
}

/**
 * 格式化文件大小为可读字符串
 * @param {number} bytes 字节数
 * @returns {string} 格式化后的大小
 */
export function formatFileSize(bytes) {
  if (bytes === 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const k = 1024
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return (bytes / Math.pow(k, i)).toFixed(2) + ' ' + units[i]
}

/**
 * 计算上传速度
 * @param {number} bytesLoaded 已上传字节数
 * @param {number} startTime 开始时间戳(毫秒)
 * @returns {string} 格式化后的速度
 */
export function calcSpeed(bytesLoaded, startTime) {
  const elapsed = (Date.now() - startTime) / 1000
  if (elapsed === 0) return '0 B/s'
  return formatFileSize(bytesLoaded / elapsed) + '/s'
}

/**
 * 计算剩余时间
 * @param {number} bytesTotal 总字节数
 * @param {number} bytesLoaded 已加载字节数
 * @param {number} startTime 开始时间戳(毫秒)
 * @returns {string} 格式化后的剩余时间
 */
export function calcRemaining(bytesTotal, bytesLoaded, startTime) {
  const elapsed = Date.now() - startTime
  if (elapsed === 0 || bytesLoaded === 0) return '--'
  const remaining = Math.round((elapsed / bytesLoaded) * (bytesTotal - bytesLoaded))
  const seconds = Math.round(remaining / 1000)
  if (seconds < 60) return `${seconds}s`
  if (seconds < 3600) return `${Math.floor(seconds / 60)}m ${seconds % 60}s`
  return `${Math.floor(seconds / 3600)}h ${Math.floor((seconds % 3600) / 60)}m`
}
