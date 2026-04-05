import { describe, it, expect, beforeEach, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useThemeStore } from './theme'

vi.mock('vue', async () => {
  const actual = await vi.importActual('vue')
  return {
    ...actual,
    watch: vi.fn()
  }
})

describe('theme store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
  })

  it('initializes with default dark theme', () => {
    const store = useThemeStore()
    expect(store.theme).toBe('dark')
  })

  it('loads saved theme from localStorage', () => {
    localStorage.setItem('wenxi-theme', 'light')
    setActivePinia(createPinia())
    const store = useThemeStore()
    expect(store.theme).toBe('light')
  })

  it('sets theme and saves to localStorage', () => {
    const store = useThemeStore()
    store.setTheme('light')
    expect(store.theme).toBe('light')
    expect(localStorage.getItem('wenxi-theme')).toBe('light')
  })

  it('supports auto theme mode', () => {
    const store = useThemeStore()
    store.setTheme('auto')
    expect(store.theme).toBe('auto')
    expect(localStorage.getItem('wenxi-theme')).toBe('auto')
  })
})
