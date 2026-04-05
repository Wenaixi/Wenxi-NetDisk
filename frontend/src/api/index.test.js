import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import api, { authAPI, fileAPI, folderAPI, shareAPI } from './index'

// Mock axios
vi.mock('axios', () => ({
  default: {
    create: () => {
      const mockAxios = {
        defaults: { headers: {} },
        interceptors: {
          request: { use: vi.fn() },
          response: { use: vi.fn() }
        },
        get: vi.fn().mockResolvedValue({ data: {} }),
        post: vi.fn().mockResolvedValue({ data: {} }),
        put: vi.fn().mockResolvedValue({ data: {} }),
        delete: vi.fn().mockResolvedValue({ data: {} }),
        _get: vi.fn(),
        _post: vi.fn(),
        _put: vi.fn(),
        _delete: vi.fn()
      }
      // Store references for testing
      return {
        interceptors: mockAxios.interceptors,
        get: (...args) => { mockAxios._get?.(...args); return Promise.resolve({ data: mockAxios._result || {} }) },
        post: (...args) => { mockAxios._post?.(...args); return Promise.resolve({ data: mockAxios._result || {} }) },
        put: (...args) => { mockAxios._put?.(...args); return Promise.resolve({ data: mockAxios._result || {} }) },
        delete: (...args) => { mockAxios._delete?.(...args); return Promise.resolve({ data: mockAxios._result || {} }) },
        _setResult: (result) => { mockAxios._result = result },
        _getCalls: () => mockAxios._get?.mock?.calls || [],
        _postCalls: () => mockAxios._post?.mock?.calls || []
      }
    }
  }
}))

describe('API Configuration', () => {
  it('should create axios instance with correct config', () => {
    // Since we mock axios, we verify the module exports exist
    expect(api).toBeDefined()
    expect(authAPI).toBeDefined()
    expect(fileAPI).toBeDefined()
    expect(folderAPI).toBeDefined()
    expect(shareAPI).toBeDefined()
  })
})

describe('authAPI', () => {
  it('should have register method', () => {
    expect(typeof authAPI.register).toBe('function')
  })

  it('should have login method', () => {
    expect(typeof authAPI.login).toBe('function')
  })

  it('should have me method', () => {
    expect(typeof authAPI.me).toBe('function')
  })
})

describe('fileAPI', () => {
  it('should have list method', () => {
    expect(typeof fileAPI.list).toBe('function')
  })

  it('should have upload method', () => {
    expect(typeof fileAPI.upload).toBe('function')
  })

  it('should have getUploadUrl method', () => {
    expect(typeof fileAPI.getUploadUrl).toBe('function')
  })

  it('should have delete method', () => {
    expect(typeof fileAPI.delete).toBe('function')
  })

  it('should have rename method', () => {
    expect(typeof fileAPI.rename).toBe('function')
  })

  it('should have move method', () => {
    expect(typeof fileAPI.move).toBe('function')
  })

  it('should have updateDescription method', () => {
    expect(typeof fileAPI.updateDescription).toBe('function')
  })

  it('should have getVersions method', () => {
    expect(typeof fileAPI.getVersions).toBe('function')
  })

  it('should have restoreVersion method', () => {
    expect(typeof fileAPI.restoreVersion).toBe('function')
  })

  it('should have initializeUpload method', () => {
    expect(typeof fileAPI.initializeUpload).toBe('function')
  })

  it('should have completeUpload method', () => {
    expect(typeof fileAPI.completeUpload).toBe('function')
  })

  it('should have download method', () => {
    expect(typeof fileAPI.download).toBe('function')
  })
})

describe('folderAPI', () => {
  it('should have list method', () => {
    expect(typeof folderAPI.list).toBe('function')
  })

  it('should have create method', () => {
    expect(typeof folderAPI.create).toBe('function')
  })

  it('should have delete method', () => {
    expect(typeof folderAPI.delete).toBe('function')
  })

  it('should have rename method', () => {
    expect(typeof folderAPI.rename).toBe('function')
  })

  it('should have updateDescription method', () => {
    expect(typeof folderAPI.updateDescription).toBe('function')
  })

  it('should have move method', () => {
    expect(typeof folderAPI.move).toBe('function')
  })
})

describe('shareAPI', () => {
  it('should have create method', () => {
    expect(typeof shareAPI.create).toBe('function')
  })

  it('should have list method', () => {
    expect(typeof shareAPI.list).toBe('function')
  })

  it('should have get method', () => {
    expect(typeof shareAPI.get).toBe('function')
  })

  it('should have validate method', () => {
    expect(typeof shareAPI.validate).toBe('function')
  })

  it('should have delete method', () => {
    expect(typeof shareAPI.delete).toBe('function')
  })
})
