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

  it('should have delete method', () => {
    expect(typeof lanzouAPI.delete).toBe('function')
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

  it('should have setAccess method', () => {
    expect(typeof lanzouAPI.setAccess).toBe('function')
  })

  it('should have getFileDetail method', () => {
    expect(typeof lanzouAPI.getFileDetail).toBe('function')
  })

  it('should have batchDelete method', () => {
    expect(typeof lanzouAPI.batchDelete).toBe('function')
  })

  it('should have getFileDescription method', () => {
    expect(typeof lanzouAPI.getFileDescription).toBe('function')
  })

  it('should have setFileDescription method', () => {
    expect(typeof lanzouAPI.setFileDescription).toBe('function')
  })

  it('should have batchMove method', () => {
    expect(typeof lanzouAPI.batchMove).toBe('function')
  })

  it('should have parseShare method', () => {
    expect(typeof lanzouAPI.parseShare).toBe('function')
  })

  it('should have batchParse method', () => {
    expect(typeof lanzouAPI.batchParse).toBe('function')
  })
})
