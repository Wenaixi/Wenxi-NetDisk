import { describe, it, expect, vi, beforeEach } from 'vitest'
import {
  getFileInfo,
  needsSplit,
  splitFile,
  mergeBlobs,
  getChunkFileName,
  parseChunkFileName,
  formatFileSize,
  calcSpeed,
  calcRemaining,
  DEFAULT_CHUNK_SIZE
} from './fileSplit'

describe('fileSplit utils', () => {
  describe('getFileInfo', () => {
    it('should return correct file info', () => {
      const file = new File(['test content'], 'test.txt', { type: 'text/plain' })
      const info = getFileInfo(file)

      expect(info.name).toBe('test.txt')
      expect(info.size).toBe(12)
      expect(info.type).toBe('text/plain')
      expect(info.chunks).toBe(1)
    })
  })

  describe('needsSplit', () => {
    it('should return false for small files', () => {
      const file = new File(['small'], 'small.txt')
      expect(needsSplit(file)).toBe(false)
    })

    it('should return true for large files (over 20MB)', () => {
      const size = 21 * 1024 * 1024
      const file = new File([new ArrayBuffer(size)], 'large.zip')
      expect(needsSplit(file)).toBe(true)
    })

    it('should respect custom maxSize', () => {
      const file = new File([new ArrayBuffer(1025)], 'file.bin')
      expect(needsSplit(file, 1024)).toBe(true)
    })
  })

  describe('splitFile', () => {
    it('should split file into correct number of chunks', async () => {
      const size = DEFAULT_CHUNK_SIZE * 3 + 100
      const file = new File([new ArrayBuffer(size)], 'test.bin')
      const { chunks, total } = await splitFile(file)

      expect(total).toBe(4)
      expect(chunks.length).toBe(4)
    })

    it('should have correct chunk sizes', async () => {
      const size = DEFAULT_CHUNK_SIZE * 2
      const file = new File([new ArrayBuffer(size)], 'test.bin')
      const { chunks } = await splitFile(file)

      expect(chunks[0].size).toBe(DEFAULT_CHUNK_SIZE)
      expect(chunks[1].size).toBe(DEFAULT_CHUNK_SIZE)
    })

    it('should mark last chunk correctly', async () => {
      const size = DEFAULT_CHUNK_SIZE + 100
      const file = new File([new ArrayBuffer(size)], 'test.bin')
      const { chunks } = await splitFile(file)

      expect(chunks[0].isLast).toBe(false)
      expect(chunks[1].isLast).toBe(true)
    })

    it('should return single chunk for small files', async () => {
      const file = new File(['hello'], 'hello.txt')
      const { chunks, total } = await splitFile(file)

      expect(total).toBe(1)
      expect(chunks.length).toBe(1)
    })
  })

  describe('mergeBlobs', () => {
    it('should merge multiple blobs', () => {
      const b1 = new Blob(['hello'])
      const b2 = new Blob([' '])
      const b3 = new Blob(['world'])

      const merged = mergeBlobs([b1, b2, b3], 'merged.txt')

      expect(merged.size).toBe(11)
    })
  })

  describe('getChunkFileName', () => {
    it('should generate chunk file name with extension', () => {
      expect(getChunkFileName('video.mp4', 1, 3)).toBe('video.part001of3.mp4')
    })

    it('should generate chunk file name without extension', () => {
      expect(getChunkFileName('data', 2, 5)).toBe('data.part002of5')
    })

    it('should pad index with zeros', () => {
      expect(getChunkFileName('file.txt', 12, 100)).toBe('file.part012of100.txt')
    })
  })

  describe('parseChunkFileName', () => {
    it('should parse chunk file name correctly', () => {
      const result = parseChunkFileName('video.part002of3.mp4')

      expect(result).not.toBeNull()
      expect(result.name).toBe('video.mp4')
      expect(result.index).toBe(2)
      expect(result.total).toBe(3)
    })

    it('should return null for non-chunk file names', () => {
      expect(parseChunkFileName('normal.txt')).toBeNull()
    })

    it('should parse file name without extension', () => {
      const result = parseChunkFileName('data.part001of5')

      expect(result.name).toBe('data')
      expect(result.index).toBe(1)
      expect(result.total).toBe(5)
    })
  })

  describe('formatFileSize', () => {
    it('should format bytes', () => {
      expect(formatFileSize(0)).toBe('0 B')
      expect(formatFileSize(512)).toBe('512.00 B')
    })

    it('should format KB', () => {
      expect(formatFileSize(1024)).toBe('1.00 KB')
      expect(formatFileSize(2048)).toBe('2.00 KB')
    })

    it('should format MB', () => {
      expect(formatFileSize(1048576)).toBe('1.00 MB')
    })

    it('should format GB', () => {
      expect(formatFileSize(1073741824)).toBe('1.00 GB')
    })
  })

  describe('calcSpeed', () => {
    it('should return 0 B/s when no time elapsed', () => {
      const now = Date.now()
      expect(calcSpeed(0, now)).toBe('0 B/s')
    })
  })

  describe('calcRemaining', () => {
    it('should return -- when no data loaded', () => {
      expect(calcRemaining(1000, 0, Date.now())).toBe('--')
    })

    it('should return -- when no time elapsed', () => {
      const now = Date.now()
      expect(calcRemaining(1000, 500, now)).toBe('--')
    })
  })
})
