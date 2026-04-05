import { describe, it, expect } from 'vitest'
import { lanzouAPI } from './lanzou'

describe('lanzouAPI', () => {
  it('should have connect method', () => {
    expect(typeof lanzouAPI.connect).toBe('function')
  })

  it('should have status method', () => {
    expect(typeof lanzouAPI.status).toBe('function')
  })

  it('should have disconnect method', () => {
    expect(typeof lanzouAPI.disconnect).toBe('function')
  })

  it('should have listFiles method', () => {
    expect(typeof lanzouAPI.listFiles).toBe('function')
  })

  it('should have listFolders method', () => {
    expect(typeof lanzouAPI.listFolders).toBe('function')
  })

  it('should have createFolder method', () => {
    expect(typeof lanzouAPI.createFolder).toBe('function')
  })

  it('should have deleteFile method', () => {
    expect(typeof lanzouAPI.deleteFile).toBe('function')
  })

  it('should have deleteFolder method', () => {
    expect(typeof lanzouAPI.deleteFolder).toBe('function')
  })

  it('should have rename method', () => {
    expect(typeof lanzouAPI.rename).toBe('function')
  })

  it('should have move method', () => {
    expect(typeof lanzouAPI.move).toBe('function')
  })

  it('should have getDownloadUrl method', () => {
    expect(typeof lanzouAPI.getDownloadUrl).toBe('function')
  })

  it('should have createShare method', () => {
    expect(typeof lanzouAPI.createShare).toBe('function')
  })
})
