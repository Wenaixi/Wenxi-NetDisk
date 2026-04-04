<template>
  <div class="min-h-screen bg-[#0f0f0f] flex items-center justify-center p-4">
    <div class="w-full max-w-md">
      <div class="text-center mb-8">
        <h1 class="text-3xl font-bold text-white mb-2">文希云盘</h1>
        <p class="text-gray-400">安全、简洁的文件存储服务</p>
      </div>

      <div class="bg-[#1a1a1a] p-8">
        <h2 class="text-xl text-white mb-6">登录</h2>

        <n-form ref="formRef" :model="form" :rules="rules">
          <n-form-item path="email" label="邮箱">
            <n-input v-model:value="form.email" placeholder="请输入邮箱" size="large" />
          </n-form-item>
          <n-form-item path="password" label="密码">
            <n-input v-model:value="form.password" type="password" placeholder="请输入密码" size="large" />
          </n-form-item>
        </n-form>

        <n-button type="primary" block size="large" :loading="loading" @click="handleLogin" class="mt-4">
          登录
        </n-button>

        <div class="text-center mt-4">
          <n-button text @click="$router.push('/register')">
            还没有账号？立即注册
          </n-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { NForm, NFormItem, NInput, NButton, useMessage } from 'naive-ui'

const router = useRouter()
const authStore = useAuthStore()
const message = useMessage()

const formRef = ref(null)
const loading = ref(false)
const form = ref({
  email: '',
  password: ''
})

const rules = {
  email: [{ required: true, message: '请输入邮箱', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

async function handleLogin() {
  try {
    loading.value = true
    await authStore.login(form.value)
    message.success('登录成功')
    router.push('/dashboard')
  } catch (err) {
    message.error(err.message || '登录失败')
  } finally {
    loading.value = false
  }
}
</script>