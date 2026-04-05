import { describe, it, expect } from 'vitest'
import { recycleAPI } from './recycle'

describe('recycleAPI', () => {
  it('should have list method', () => {
    expect(typeof recycleAPI.list).toBe('function')
  })

  it('should have restoreFile method', () => {
    expect(typeof recycleAPI.restoreFile).toBe('function')
  })

  it('should have restoreFolder method', () => {
    expect(typeof recycleAPI.restoreFolder).toBe('function')
  })

  it('should have permanentDelete method', () => {
    expect(typeof recycleAPI.permanentDelete).toBe('function')
  })

  it('should have clearAll method', () => {
    expect(typeof recycleAPI.clearAll).toBe('function')
  })
})
