<template>
  <div class="min-h-screen bg-[#0f0f0f] flex items-center justify-center p-4">
    <div class="w-full max-w-md">
      <div class="text-center mb-8">
        <h1 class="text-3xl font-bold text-white mb-2">文希云盘</h1>
        <p class="text-gray-400">安全、简洁的文件存储服务</p>
      </div>

      <div class="bg-[#1a1a1a] p-8">
        <h2 class="text-xl text-white mb-6">注册</h2>

        <n-form ref="formRef" :model="form" :rules="rules">
          <n-form-item path="email" label="邮箱">
            <n-input v-model:value="form.email" placeholder="请输入邮箱" size="large" />
          </n-form-item>
          <n-form-item path="password" label="密码">
            <n-input v-model:value="form.password" type="password" placeholder="请输入密码" size="large" />
          </n-form-item>
          <n-form-item path="confirmPassword" label="确认密码">
            <n-input v-model:value="form.confirmPassword" type="password" placeholder="请再次输入密码" size="large" />
          </n-form-item>
        </n-form>

        <n-button type="primary" block size="large" :loading="loading" @click="handleRegister" class="mt-4">
          注册
        </n-button>

        <div class="text-center mt-4">
          <n-button text @click="$router.push('/login')">
            已有账号？立即登录
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
  password: '',
  confirmPassword: ''
})

const validateConfirmPassword = (rule, value) => {
  if (value !== form.value.password) {
    return new Error('两次输入的密码不一致')
  }
  return true
}

const rules = {
  email: [{ required: true, message: '请输入邮箱', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

async function handleRegister() {
  try {
    loading.value = true
    await authStore.register(form.value)
    message.success('注册成功')
    router.push('/dashboard')
  } catch (err) {
    message.error(err.message || '注册失败')
  } finally {
    loading.value = false
  }
}
</script>