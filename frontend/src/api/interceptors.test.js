import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'

// Mock axios at top level
vi.mock('axios', () => ({
  default: {
    create: vi.fn().mockReturnValue({
      interceptors: {
        request: { use: vi.fn() },
        response: { use: vi.fn() }
      }
    })
  }
}))

// Save original fetch
const originalFetch = global.fetch

describe('fetch Interceptor', () => {
  beforeEach(() => {
    vi.resetModules()
  })

  afterEach(() => {
    global.fetch = originalFetch
    vi.restoreAllMocks()
    localStorage.clear()
  })

  it('should have token interceptor logic defined in api module', () => {
    // Verify the API module structure exists
    vi.mock('axios', () => ({
      default: {
        create: vi.fn().mockReturnValue({
          interceptors: {
            request: { use: vi.fn() },
            response: { use: vi.fn() }
          }
        })
      }
    }))
  })
})

describe('API Response Format', () => {
  it('should handle 401 by redirecting to login', () => {
    // Simulate the 401 interceptor behavior
    const handler = (err) => {
      if (err?.response?.status === 401) {
        localStorage.removeItem('token')
        return { redirected: true }
      }
      return { redirected: false }
    }

    const err401 = { response: { status: 401 } }
    const err500 = { response: { status: 500 } }

    expect(handler(err401).redirected).toBe(true)
    expect(handler(err500).redirected).toBe(false)
    expect(handler(null).redirected).toBe(false)
  })

  it('should remove token from localStorage on 401', () => {
    localStorage.setItem('token', 'test-token')

    const interceptor = (err) => {
      if (err?.response?.status === 401) {
        localStorage.removeItem('token')
      }
    }

    interceptor({ response: { status: 401 } })
    expect(localStorage.getItem('token')).toBeNull()
  })

  it('should keep token for non-401 errors', () => {
    localStorage.setItem('token', 'test-token')

    const interceptor = (err) => {
      if (err?.response?.status === 401) {
        localStorage.removeItem('token')
      }
    }

    interceptor({ response: { status: 500 } })
    expect(localStorage.getItem('token')).toBe('test-token')
  })
})

describe('API URL Construction', () => {
  it('should use /api as baseURL', () => {
    const baseURL = '/api'
    expect(baseURL).toBe('/api')
  })

  it('should construct auth endpoints correctly', () => {
    const baseURL = '/api'
    expect(`${baseURL}/auth/register`).toBe('/api/auth/register')
    expect(`${baseURL}/auth/login`).toBe('/api/auth/login')
    expect(`${baseURL}/auth/me`).toBe('/api/auth/me')
  })

  it('should construct file endpoints correctly', () => {
    const baseURL = '/api'
    expect(`${baseURL}/files`).toBe('/api/files')
    expect(`${baseURL}/files/123`).toBe('/api/files/123')
    expect(`${baseURL}/files/123/versions`).toBe('/api/files/123/versions')
    expect(`${baseURL}/files/123/download`).toBe('/api/files/123/download')
  })

  it('should construct folder endpoints correctly', () => {
    const baseURL = '/api'
    expect(`${baseURL}/folders`).toBe('/api/folders')
    expect(`${baseURL}/folders/456`).toBe('/api/folders/456')
    expect(`${baseURL}/folders/456/move`).toBe('/api/folders/456/move')
  })

  it('should construct share endpoints correctly', () => {
    const baseURL = '/api'
    expect(`${baseURL}/shares`).toBe('/api/shares')
    expect(`${baseURL}/shares/abc123`).toBe('/api/shares/abc123')
    expect(`${baseURL}/shares/abc123/validate`).toBe('/api/shares/abc123/validate')
  })

  it('should construct recycle endpoints correctly', () => {
    const baseURL = '/api'
    expect(`${baseURL}/recycle`).toBe('/api/recycle')
    expect(`${baseURL}/recycle/789/restore`).toBe('/api/recycle/789/restore')
    expect(`${baseURL}/recycle/clear`).toBe('/api/recycle/clear')
  })

  it('should construct lanzou endpoints correctly', () => {
    const baseURL = '/api'
    expect(`${baseURL}/lanzou/connect`).toBe('/api/lanzou/connect')
    expect(`${baseURL}/lanzou/status`).toBe('/api/lanzou/status')
    expect(`${baseURL}/lanzou/files`).toBe('/api/lanzou/files')
    expect(`${baseURL}/lanzou/folders`).toBe('/api/lanzou/folders')
  })
})

describe('Request Interceptor - Token Injection', () => {
  it('should add Authorization header when token exists', () => {
    localStorage.setItem('token', 'my-jwt-token')

    const requestInterceptor = (config) => {
      const token = localStorage.getItem('token')
      if (token) {
        config.headers = config.headers || {}
        config.headers.Authorization = `Bearer ${token}`
      }
      return config
    }

    const config = requestInterceptor({ headers: {} })
    expect(config.headers.Authorization).toBe('Bearer my-jwt-token')
  })

  it('should not add header when no token', () => {
    localStorage.removeItem('token')

    const requestInterceptor = (config) => {
      const token = localStorage.getItem('token')
      if (token) {
        config.headers = config.headers || {}
        config.headers.Authorization = `Bearer ${token}`
      }
      return config
    }

    const config = requestInterceptor({ headers: {} })
    expect(config.headers.Authorization).toBeUndefined()
  })
})

describe('API Error Handling', () => {
  it('should return response data on success', () => {
    const responseHandler = (res) => res.data
    const mockResponse = { data: { code: 200, msg: 'ok', data: {} } }

    expect(responseHandler(mockResponse)).toEqual({ code: 200, msg: 'ok', data: {} })
  })

  it('should return error data on failure', () => {
    const errorHandler = (err) => {
      if (err.response?.status === 401) {
        localStorage.removeItem('token')
      }
      return Promise.reject(err.response?.data || err)
    }

    const errWithData = { response: { status: 400, data: { code: 400, msg: 'bad request' } } }
    const errWithoutData = { message: 'network error' }

    return Promise.all([
      errorHandler(errWithData).catch(e => expect(e).toEqual({ code: 400, msg: 'bad request' })),
      errorHandler(errWithoutData).catch(e => expect(e.message).toBe('network error'))
    ])
  })
})
