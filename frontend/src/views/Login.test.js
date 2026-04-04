import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import Login from './Login.vue'

// Mock naive-ui
vi.mock('naive-ui', () => ({
  NForm: {
    template: '<form><slot /></form>',
    props: ['model', 'rules', 'ref']
  },
  NFormItem: {
    template: '<div class="n-form-item"><span class="label">{{ label }}</span><slot /></div>',
    props: ['path', 'label']
  },
  NInput: {
    template: '<input class="n-input" :type="type || \'text\'" :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />',
    props: ['modelValue', 'type', 'placeholder', 'size']
  },
  NButton: {
    template: '<button class="n-button" :type="type" @click="$emit(\'click\')"><slot /></button>',
    props: ['type', 'block', 'size', 'loading', 'text']
  },
  useMessage: () => ({
    success: vi.fn(),
    error: vi.fn(),
    warning: vi.fn(),
    info: vi.fn()
  })
}))

// Mock vue-router
const mockPush = vi.fn()
vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: mockPush
  }),
  useRoute: () => ({ path: '/login' })
}))

// Mock auth store
const mockLogin = vi.fn()
const mockLogout = vi.fn()

vi.mock('../stores/auth', () => ({
  useAuthStore: () => ({
    user: { email: 'test@example.com' },
    login: mockLogin,
    logout: mockLogout
  })
}))

const mountComponent = () => {
  return mount(Login, {
    global: {
      plugins: [],
      stubs: {
        RouterLink: true
      }
    }
  })
}

describe('Login.vue', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    mockPush.mockReset()
  })

  it('should render login form elements', () => {
    const wrapper = mountComponent()

    expect(wrapper.text()).toContain('文希云盘')
    expect(wrapper.text()).toContain('登录')
    expect(wrapper.text()).toContain('邮箱')
    expect(wrapper.text()).toContain('密码')
  })

  it('should have email and password input fields', () => {
    const wrapper = mountComponent()

    const inputs = wrapper.findAll('.n-input')
    expect(inputs.length).toBe(2)

    expect(inputs[0].attributes('type')).toBe('text')
    expect(inputs[1].attributes('type')).toBe('password')
  })

  it('should have login button', () => {
    const wrapper = mountComponent()

    const buttons = wrapper.findAll('.n-button')
    const loginButton = buttons.find(b => b.attributes('type') === 'primary')
    expect(loginButton).toBeTruthy()
  })

  it('should have register link', () => {
    const wrapper = mountComponent()

    expect(wrapper.text()).toContain('还没有账号？立即注册')
  })

  it('should call authStore.login on form submit', async () => {
    mockLogin.mockResolvedValue({ user: { email: 'test@example.com' } })
    const wrapper = mountComponent()

    const inputs = wrapper.findAll('.n-input')
    await inputs[0].setValue('test@example.com')
    await inputs[1].setValue('password123')

    const buttons = wrapper.findAll('.n-button')
    const loginButton = buttons.find(b => b.attributes('type') === 'primary')
    await loginButton.trigger('click')

    expect(mockLogin).toHaveBeenCalled()
  })

  it('should redirect to dashboard on successful login', async () => {
    mockLogin.mockResolvedValue({ user: { email: 'test@example.com' } })
    const wrapper = mountComponent()

    const inputs = wrapper.findAll('.n-input')
    await inputs[0].setValue('test@example.com')
    await inputs[1].setValue('password123')

    const buttons = wrapper.findAll('.n-button')
    const loginButton = buttons.find(b => b.attributes('type') === 'primary')
    await loginButton.trigger('click')

    expect(mockPush).toHaveBeenCalledWith('/dashboard')
  })

  it('should not redirect on login failure', async () => {
    mockLogin.mockRejectedValue(new Error('Invalid credentials'))
    const wrapper = mountComponent()

    const inputs = wrapper.findAll('.n-input')
    await inputs[0].setValue('test@example.com')
    await inputs[1].setValue('wrong')

    const buttons = wrapper.findAll('.n-button')
    const loginButton = buttons.find(b => b.attributes('type') === 'primary')
    await loginButton.trigger('click')

    expect(mockPush).not.toHaveBeenCalledWith('/dashboard')
  })

  it('should have correct form validation rules', () => {
    const wrapper = mountComponent()

    expect(wrapper.vm.rules).toBeTruthy()
    expect(wrapper.vm.rules.email).toBeTruthy()
    expect(wrapper.vm.rules.password).toBeTruthy()
  })
})
