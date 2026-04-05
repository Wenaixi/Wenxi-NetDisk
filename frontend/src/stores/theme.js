import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

const THEME_KEY = 'wenxi-theme'

export const useThemeStore = defineStore('theme', () => {
  const savedTheme = localStorage.getItem(THEME_KEY) || 'dark'
  const theme = ref(savedTheme) // 'dark' | 'light' | 'auto'

  function setTheme(newTheme) {
    theme.value = newTheme
    localStorage.setItem(THEME_KEY, newTheme)
    applyTheme()
  }

  function applyTheme() {
    let actualTheme = theme.value
    if (actualTheme === 'auto') {
      actualTheme = window.matchMedia('(prefers-color-scheme: light)').matches ? 'light' : 'dark'
    }
    const root = document.documentElement
    if (actualTheme === 'light') {
      root.setAttribute('data-theme', 'light')
      root.style.setProperty('--bg-primary', '#f5f5f5')
      root.style.setProperty('--bg-card', '#ffffff')
      root.style.setProperty('--bg-hover', '#e8e8e8')
      root.style.setProperty('--border', '#d9d9d9')
      root.style.setProperty('--text-primary', '#1a1a1a')
      root.style.setProperty('--text-secondary', '#666666')
      root.style.setProperty('--bg-table-hover', '#fafafa')
    } else {
      root.setAttribute('data-theme', 'dark')
      root.style.setProperty('--bg-primary', '#0f0f0f')
      root.style.setProperty('--bg-card', '#1a1a1a')
      root.style.setProperty('--bg-hover', '#252525')
      root.style.setProperty('--border', '#333333')
      root.style.setProperty('--text-primary', '#ffffff')
      root.style.setProperty('--text-secondary', '#888888')
      root.style.setProperty('--bg-table-hover', '#1e1e1e')
    }
  }

  function initTheme() {
    applyTheme()
    // 监听系统主题变化
    const mediaQuery = window.matchMedia('(prefers-color-scheme: light)')
    mediaQuery.addEventListener('change', () => {
      if (theme.value === 'auto') {
        applyTheme()
      }
    })
  }

  // 响应式监听 theme 变化
  watch(theme, () => {
    applyTheme()
  })

  return {
    theme,
    setTheme,
    initTheme
  }
})
