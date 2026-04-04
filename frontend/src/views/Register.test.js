import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import Register from './Register.vue'

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
  useRoute: () => ({ path: '/register' })
}))

// Mock auth store
const mockRegister = vi.fn()

vi.mock('../stores/auth', () => ({
  useAuthStore: () => ({
    user: { email: 'test@example.com' },
    register: mockRegister
  })
}))

const mountComponent = () => {
  return mount(Register, {
    global: {
      plugins: [],
      stubs: {
        RouterLink: true
      }
    }
  })
}

describe('Register.vue', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    mockPush.mockReset()
  })

  it('should render registration form elements', () => {
    const wrapper = mountComponent()

    expect(wrapper.text()).toContain('文希云盘')
    expect(wrapper.text()).toContain('注册')
    expect(wrapper.text()).toContain('确认密码')
  })

  it('should have three input fields', () => {
    const wrapper = mountComponent()

    const inputs = wrapper.findAll('.n-input')
    expect(inputs.length).toBe(3)
  })

  it('should have password confirmation validation', () => {
    const wrapper = mountComponent()

    expect(wrapper.vm.rules.confirmPassword).toBeTruthy()
    expect(wrapper.vm.rules.confirmPassword.length).toBe(2)
  })

  it('should validate password match - fail on different passwords', () => {
    const wrapper = mountComponent()

    wrapper.vm.form.password = 'password123'

    const validator = wrapper.vm.rules.confirmPassword[1].validator
    const result = validator(null, 'different')
    expect(result).toBeInstanceOf(Error)
  })

  it('should validate password match - pass on same password', () => {
    const wrapper = mountComponent()

    wrapper.vm.form.password = 'password123'

    const validator = wrapper.vm.rules.confirmPassword[1].validator
    const result = validator(null, 'password123')
    expect(result).toBe(true)
  })

  it('should call authStore.register on form submit', async () => {
    mockRegister.mockResolvedValue({ user: { email: 'test@example.com' } })
    const wrapper = mountComponent()

    const inputs = wrapper.findAll('.n-input')
    await inputs[0].setValue('test@example.com')
    await inputs[1].setValue('password123')
    await inputs[2].setValue('password123')

    const buttons = wrapper.findAll('.n-button')
    const registerButton = buttons.find(b => b.attributes('type') === 'primary')
    await registerButton.trigger('click')

    expect(mockRegister).toHaveBeenCalled()
  })

  it('should redirect to dashboard on successful registration', async () => {
    mockRegister.mockResolvedValue({ user: { email: 'test@example.com' } })
    const wrapper = mountComponent()

    const inputs = wrapper.findAll('.n-input')
    await inputs[0].setValue('test@example.com')
    await inputs[1].setValue('password123')
    await inputs[2].setValue('password123')

    const buttons = wrapper.findAll('.n-button')
    const registerButton = buttons.find(b => b.attributes('type') === 'primary')
    await registerButton.trigger('click')

    expect(mockPush).toHaveBeenCalledWith('/dashboard')
  })

  it('should have login link', () => {
    const wrapper = mountComponent()

    expect(wrapper.text()).toContain('已有账号？立即登录')
  })

  it('should handle registration failure', async () => {
    mockRegister.mockRejectedValue(new Error('Email already exists'))
    const wrapper = mountComponent()

    const inputs = wrapper.findAll('.n-input')
    await inputs[0].setValue('existing@example.com')
    await inputs[1].setValue('password123')
    await inputs[2].setValue('password123')

    const buttons = wrapper.findAll('.n-button')
    const registerButton = buttons.find(b => b.attributes('type') === 'primary')
    await registerButton.trigger('click')

    expect(mockPush).not.toHaveBeenCalledWith('/dashboard')
  })

  it('should show loading state during registration', async () => {
    let resolvePromise
    mockRegister.mockImplementation(() => new Promise(resolve => {
      resolvePromise = resolve
    }))
    const wrapper = mountComponent()

    const inputs = wrapper.findAll('.n-input')
    await inputs[0].setValue('test@example.com')
    await inputs[1].setValue('password123')
    await inputs[2].setValue('password123')

    const buttons = wrapper.findAll('.n-button')
    const registerButton = buttons.find(b => b.attributes('type') === 'primary')
    await registerButton.trigger('click')

    expect(wrapper.vm.loading).toBe(true)

    resolvePromise({ user: { email: 'test@example.com' } })
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.loading).toBe(false)
  })
})
