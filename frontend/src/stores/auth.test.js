import { setActivePinia, createPinia } from 'pinia'
import { describe, it, expect, beforeEach, vi } from 'vitest'
import { useAuthStore } from './auth'

// Mock localStorage
const localStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn()
}
global.localStorage = localStorageMock

// Mock authAPI
vi.mock('../api', () => ({
  authAPI: {
    login: vi.fn(),
    register: vi.fn(),
    me: vi.fn()
  }
}))

import { authAPI } from '../api'

describe('useAuthStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorageMock.getItem.mockReturnValue(null)
    localStorageMock.setItem.mockClear()
    localStorageMock.removeItem.mockClear()
    vi.clearAllMocks()
  })

  it('should initialize with empty token and user', () => {
    const store = useAuthStore()
    expect(store.token).toBe('')
    expect(store.user).toBe(null)
    expect(store.isLoggedIn).toBe(false)
  })

  it('should login successfully', async () => {
    const store = useAuthStore()
    const mockResponse = {
      token: 'test-token-123',
      user: { id: 1, email: 'test@example.com' }
    }
    authAPI.login.mockResolvedValue(mockResponse)

    const result = await store.login({ email: 'test@example.com', password: 'password' })

    expect(store.token).toBe('test-token-123')
    expect(store.user).toEqual(mockResponse.user)
    expect(store.isLoggedIn).toBe(true)
    expect(localStorage.setItem).toHaveBeenCalledWith('token', 'test-token-123')
    expect(result).toEqual(mockResponse)
  })

  it('should register successfully', async () => {
    const store = useAuthStore()
    const mockResponse = {
      token: 'register-token-456',
      user: { id: 2, email: 'new@example.com' }
    }
    authAPI.register.mockResolvedValue(mockResponse)

    const result = await store.register({ email: 'new@example.com', password: 'password' })

    expect(store.token).toBe('register-token-456')
    expect(store.user).toEqual(mockResponse.user)
    expect(store.isLoggedIn).toBe(true)
  })

  it('should fetch user successfully', async () => {
    const store = useAuthStore()
    store.token = 'existing-token'
    const mockUser = { id: 1, email: 'test@example.com' }
    authAPI.me.mockResolvedValue(mockUser)

    const result = await store.fetchUser()

    expect(store.user).toEqual(mockUser)
    expect(result).toEqual(mockUser)
  })

  it('should return null and logout when fetchUser fails', async () => {
    const store = useAuthStore()
    store.token = 'invalid-token'
    authAPI.me.mockRejectedValue(new Error('Unauthorized'))

    const result = await store.fetchUser()

    expect(store.user).toBe(null)
    expect(store.token).toBe('')
    expect(result).toBe(null)
  })

  it('should logout and clear storage', () => {
    const store = useAuthStore()
    store.token = 'some-token'
    store.user = { id: 1 }

    store.logout()

    expect(store.token).toBe('')
    expect(store.user).toBe(null)
    expect(store.isLoggedIn).toBe(false)
    expect(localStorage.removeItem).toHaveBeenCalledWith('token')
  })

  it('should initialize token from localStorage', () => {
    localStorageMock.getItem.mockReturnValue('stored-token')
    const store = useAuthStore()
    expect(store.token).toBe('stored-token')
    expect(store.isLoggedIn).toBe(true)
  })
})
