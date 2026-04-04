import { describe, it, expect } from 'vitest'
import { arrayBufferToBase64, base64ToArrayBuffer } from '../utils/crypto'

// Test pure functions that don't need crypto API
describe('crypto utils - pure functions', () => {
  describe('arrayBufferToBase64', () => {
    it('should convert ArrayBuffer to base64', () => {
      const buffer = new Uint8Array([72, 101, 108, 108, 111]) // "Hello"
      const result = arrayBufferToBase64(buffer)

      expect(result).toBe('SGVsbG8=')
    })

    it('should convert empty buffer to empty string', () => {
      const buffer = new Uint8Array([])
      const result = arrayBufferToBase64(buffer)

      expect(result).toBe('')
    })
  })

  describe('base64ToArrayBuffer', () => {
    it('should convert base64 to ArrayBuffer', () => {
      const result = base64ToArrayBuffer('SGVsbG8=')
      const bytes = new Uint8Array(result)

      expect(bytes).toEqual(new Uint8Array([72, 101, 108, 108, 111]))
    })

    it('should handle empty string', () => {
      const result = base64ToArrayBuffer('')
      const bytes = new Uint8Array(result)

      expect(bytes.length).toBe(0)
    })
  })
})
