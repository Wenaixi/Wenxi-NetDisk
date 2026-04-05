import { describe, it, expect, beforeEach, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useCalculateStore } from './calculate'

describe('calculate store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
  })

  it('initializes with empty records', () => {
    const store = useCalculateStore()
    expect(store.warningSize).toBe(7 * 1024 * 1024 * 1024)
    expect(store.warningEnabled).toBe(true)
    expect(store.todayBytes).toBe(0)
  })

  it('records upload bytes for today', () => {
    const store = useCalculateStore()
    store.recordUploadBytes(1024 * 1024) // 1MB

    const today = new Date().toISOString().slice(0, 10)
    expect(store.dailyRecords[today]).toBe(1024 * 1024)
    expect(store.todayBytes).toBe(1024 * 1024)
  })

  it('accumulates upload bytes', () => {
    const store = useCalculateStore()
    store.recordUploadBytes(1024 * 1024)
    store.recordUploadBytes(2 * 1024 * 1024)

    expect(store.todayBytes).toBe(3 * 1024 * 1024)
  })

  it('checks warning size', () => {
    const store = useCalculateStore()
    store.setWarningSize(5 * 1024 * 1024) // 5MB
    store.recordUploadBytes(6 * 1024 * 1024) // 6MB

    expect(store.checkWarningSize()).toBe(true)
  })

  it('returns false when under warning', () => {
    const store = useCalculateStore()
    store.setWarningSize(10 * 1024 * 1024)
    store.recordUploadBytes(1 * 1024 * 1024)

    expect(store.checkWarningSize()).toBe(false)
  })

  it('respects warningEnabled flag', () => {
    const store = useCalculateStore()
    store.setWarningEnabled(false)
    store.setWarningSize(1)
    store.recordUploadBytes(100)

    expect(store.checkWarningSize()).toBe(false)
  })

  it('persists records to localStorage', () => {
    const store = useCalculateStore()
    store.recordUploadBytes(2048)

    const saved = JSON.parse(localStorage.getItem('wenxi-calculate'))
    const today = new Date().toISOString().slice(0, 10)
    expect(saved.dailyRecords[today]).toBe(2048)
  })

  it('loads records from localStorage', () => {
    const today = new Date().toISOString().slice(0, 10)
    localStorage.setItem('wenxi-calculate', JSON.stringify({
      dailyRecords: { [today]: 4096 },
      warningSize: 5 * 1024 * 1024 * 1024,
      warningEnabled: true
    }))

    setActivePinia(createPinia())
    const store = useCalculateStore()
    expect(store.todayBytes).toBe(4096)
  })

  it('clears history', () => {
    const store = useCalculateStore()
    store.recordUploadBytes(1024)
    store.clearHistory()

    expect(store.dailyRecords).toEqual({})
    expect(store.todayBytes).toBe(0)
  })

  it('sets warning size', () => {
    const store = useCalculateStore()
    store.setWarningSize(10 * 1024 * 1024 * 1024) // 10GB
    expect(store.warningSize).toBe(10 * 1024 * 1024 * 1024)
  })

  it('toggles warning enabled', () => {
    const store = useCalculateStore()
    store.setWarningEnabled(false)
    expect(store.warningEnabled).toBe(false)
  })

  it('clears records older than 30 days', () => {
    // Directly set old data in localStorage to simulate historical records
    const oldDate = new Date()
    oldDate.setDate(oldDate.getDate() - 31)
    const oldKey = oldDate.toISOString().slice(0, 10)
    const todayKey = new Date().toISOString().slice(0, 10)

    localStorage.setItem('wenxi-calculate', JSON.stringify({
      dailyRecords: { [oldKey]: 999, [todayKey]: 123 },
      warningSize: 7 * 1024 * 1024 * 1024,
      warningEnabled: true
    }))

    setActivePinia(createPinia())
    const store = useCalculateStore()

    // Trigger a save by recording new bytes, which will clean old records
    store.recordUploadBytes(1)

    expect(store.dailyRecords[oldKey]).toBeUndefined()
    // todayKey should still have data (123 + 1 = 124)
    expect(store.dailyRecords[todayKey]).toBe(124)
  })
})
