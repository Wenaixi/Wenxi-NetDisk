/**
 * 客户端加密工具
 * 使用Web Crypto API的AES-GCM算法实现端到端加密
 * 注意: Web Crypto API不支持ChaCha20，因此使用AES-GCM作为替代
 */

/**
 * 生成加密密钥
 * @returns {Promise<CryptoKey>} 加密密钥
 */
export async function generateEncryptionKey() {
  return await window.crypto.subtle.generateKey(
    {
      name: 'AES-GCM',
      length: 256,
    },
    true, // extractable
    ['encrypt', 'decrypt']
  )
}

/**
 * 从密码派生密钥
 * @param {string} password 用户密码
 * @param {Uint8Array} salt 盐值
 * @returns {Promise<CryptoKey>} 派生密钥
 */
export async function deriveKeyFromPassword(password, salt) {
  const encoder = new TextEncoder()
  const passwordData = encoder.encode(password)

  const keyMaterial = await window.crypto.subtle.importKey(
    'raw',
    passwordData,
    { name: 'PBKDF2' },
    false,
    ['deriveKey']
  )

  return await window.crypto.subtle.deriveKey(
    {
      name: 'PBKDF2',
      salt: salt,
      iterations: 100000,
      hash: 'SHA-256',
    },
    keyMaterial,
    { name: 'AES-GCM', length: 256 },
    true,
    ['encrypt', 'decrypt']
  )
}

/**
 * 生成随机盐值
 * @param {number} length 盐值长度(字节)
 * @returns {Uint8Array} 随机盐值
 */
export function generateSalt(length = 16) {
  return window.crypto.getRandomValues(new Uint8Array(length))
}

/**
 * 生成随机IV
 * @returns {Uint8Array} 随机IV
 */
export function generateIV() {
  return window.crypto.getRandomValues(new Uint8Array(12))
}

/**
 * 加密文件
 * @param {File} file 原始文件
 * @param {CryptoKey} key 加密密钥
 * @returns {Promise<Blob>} 加密后的文件Blob
 */
export async function encryptFile(file, key) {
  const iv = generateIV()
  const fileData = await file.arrayBuffer()

  const encryptedData = await window.crypto.subtle.encrypt(
    {
      name: 'AES-GCM',
      iv: iv,
    },
    key,
    fileData
  )

  // 将IV附加到加密数据前(12字节)
  const result = new Uint8Array(iv.length + encryptedData.byteLength)
  result.set(iv, 0)
  result.set(new Uint8Array(encryptedData), iv.length)

  return new Blob([result], { type: 'application/octet-stream' })
}

/**
 * 解密文件
 * @param {Blob} encryptedBlob 加密的文件Blob
 * @param {CryptoKey} key 解密密钥
 * @returns {Promise<Blob>} 解密后的文件Blob
 */
export async function decryptFile(encryptedBlob, key) {
  const encryptedData = await encryptedBlob.arrayBuffer()
  const data = new Uint8Array(encryptedData)

  // 提取IV(前12字节)
  const iv = data.slice(0, 12)
  const ciphertext = data.slice(12)

  const decryptedData = await window.crypto.subtle.decrypt(
    {
      name: 'AES-GCM',
      iv: iv,
    },
    key,
    ciphertext
  )

  return new Blob([decryptedData])
}

/**
 * 导出密钥为可传输格式
 * @param {CryptoKey} key 加密密钥
 * @returns {Promise<{key: string, iv: string}>} 导出的密钥和IV
 */
export async function exportKey(key) {
  const rawKey = await window.crypto.subtle.exportKey('raw', key)
  const keyArray = new Uint8Array(rawKey)

  return {
    key: arrayBufferToBase64(keyArray),
    iv: arrayBufferToBase64(generateIV()),
  }
}

/**
 * 导入密钥
 * @param {string} keyBase64 Base64编码的密钥
 * @returns {Promise<CryptoKey>} 导入的密钥
 */
export async function importKey(keyBase64) {
  const keyData = base64ToArrayBuffer(keyBase64)

  return await window.crypto.subtle.importKey(
    'raw',
    keyData,
    { name: 'AES-GCM', length: 256 },
    true,
    ['encrypt', 'decrypt']
  )
}

/**
 * ArrayBuffer转Base64
 * @param {ArrayBuffer|Uint8Array} buffer
 * @returns {string}
 */
function arrayBufferToBase64(buffer) {
  const bytes = new Uint8Array(buffer)
  let binary = ''
  for (let i = 0; i < bytes.byteLength; i++) {
    binary += String.fromCharCode(bytes[i])
  }
  return btoa(binary)
}

/**
 * Base64转ArrayBuffer
 * @param {string} base64
 * @returns {ArrayBuffer}
 */
function base64ToArrayBuffer(base64) {
  const binary = atob(base64)
  const bytes = new Uint8Array(binary.length)
  for (let i = 0; i < binary.length; i++) {
    bytes[i] = binary.charCodeAt(i)
  }
  return bytes.buffer
}

/**
 * 计算文件SHA-256哈希
 * @param {File} file
 * @returns {Promise<string>} Hex编码的哈希值
 */
export async function computeFileHash(file) {
  const data = await file.arrayBuffer()
  const hashBuffer = await window.crypto.subtle.digest('SHA-256', data)
  const hashArray = new Uint8Array(hashBuffer)
  return Array.from(hashArray)
    .map(b => b.toString(16).padStart(2, '0'))
    .join('')
}
